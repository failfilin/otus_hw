package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestNew_ReturnsServer(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNew_MuxIsInitialized(t *testing.T) {
	s := New()
	if s.mux == nil {
		t.Fatal("New() did not initialize mux")
	}
}

func TestNew_RoutesAreRegistered(t *testing.T) {
	s := New()

	// Все маршруты смонтированы под /api/ через StripPrefix,
	// поэтому проверяем именно эти пути на s.mux.
	knownRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/restaurants"},
		{http.MethodPost, "/api/restaurants"},
		{http.MethodGet, "/api/restaurants/some-id"},
		{http.MethodPut, "/api/restaurants/some-id"},
		{http.MethodDelete, "/api/restaurants/some-id"},
		{http.MethodGet, "/api/menus"},
		{http.MethodGet, "/api/dishes"},
		{http.MethodGet, "/swagger/"},
	}

	for _, route := range knownRoutes {
		req, err := http.NewRequest(route.method, route.path, nil)
		if err != nil {
			t.Fatalf("failed to create request for %s %s: %v", route.method, route.path, err)
		}

		_, pattern := s.mux.Handler(req)
		if pattern == "" {
			t.Errorf("route %s %s is not registered", route.method, route.path)
		}
	}
}

// ── HTTPServer() ──────────────────────────────────────────────────────────────

func TestHTTPServer_ReturnsHTTPServer(t *testing.T) {
	s := New()
	srv := s.HTTPServer(":0")
	if srv == nil {
		t.Fatal("HTTPServer() returned nil")
	}
}

func TestHTTPServer_AddressIsSet(t *testing.T) {
	s := New()
	addr := ":8080"
	srv := s.HTTPServer(addr)
	if srv.Addr != addr {
		t.Errorf("expected Addr %q, got %q", addr, srv.Addr)
	}
}

func TestHTTPServer_HandlerIsMux(t *testing.T) {
	s := New()
	srv := s.HTTPServer(":0")
	if srv.Handler != s.mux {
		t.Error("HTTPServer handler should be the server's mux")
	}
}

func TestHTTPServer_CalledTwice_ReturnsDifferentInstances(t *testing.T) {
	s := New()
	srv1 := s.HTTPServer(":8081")
	srv2 := s.HTTPServer(":8082")
	if srv1 == srv2 {
		t.Error("each call to HTTPServer should return a new *http.Server instance")
	}
}

// freePort получает свободный порт у ОС, чтобы не было конфликтов между тестами.
func freePort(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()
	return addr
}

func TestRun_StartsAndAcceptsConnections(t *testing.T) {
	addr := freePort(t)
	s := New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Run(ctx, addr)
	}()

	// Ждём пока сервер поднимется
	if err := waitForServer(addr, 2*time.Second); err != nil {
		t.Fatalf("server did not start in time: %v", err)
	}

	resp, err := http.Get("http://" + addr + "/restaurants")
	if err != nil {
		t.Fatalf("failed to reach server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 0 {
		t.Error("expected a valid HTTP status code")
	}
}

func TestRun_GracefulShutdown(t *testing.T) {
	addr := freePort(t)
	s := New()

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Run(ctx, addr)
	}()

	if err := waitForServer(addr, 2*time.Second); err != nil {
		t.Fatalf("server did not start: %v", err)
	}

	// Отменяем контекст — ожидаем graceful shutdown
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Run() returned error on graceful shutdown: %v", err)
		}
	case <-time.After(7 * time.Second):
		t.Error("server did not shut down within timeout")
	}
}

func TestRun_ShutdownCompletesBeforeTimeout(t *testing.T) {
	addr := freePort(t)
	s := New()

	ctx, cancel := context.WithCancel(context.Background())

	go s.Run(ctx, addr) //nolint:errcheck

	if err := waitForServer(addr, 2*time.Second); err != nil {
		t.Fatalf("server did not start: %v", err)
	}

	start := time.Now()
	cancel()
	time.Sleep(500 * time.Millisecond) // даём горутине завершиться

	elapsed := time.Since(start)
	// Graceful shutdown не должен занимать дольше 5 секунд (таймаут из Run)
	if elapsed > 5*time.Second {
		t.Errorf("shutdown took too long: %v", elapsed)
	}
}

func TestRun_ContextAlreadyCancelled(t *testing.T) {
	addr := freePort(t)
	s := New()

	// Передаём уже отменённый контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Run(ctx, addr)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected error with pre-cancelled context: %v", err)
		}
	case <-time.After(7 * time.Second):
		t.Error("Run() did not return after pre-cancelled context")
	}
}

// waitForServer пытается подключиться к addr пока не истечёт timeout.
func waitForServer(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return &net.OpError{Op: "dial", Addr: nil, Err: context.DeadlineExceeded}
}
