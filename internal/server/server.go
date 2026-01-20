package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}


func (s *Server) HTTPServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
}

// Полный запуск + graceful shutdown
func (s *Server) Run(ctx context.Context, addr string) error {
	httpServer := s.HTTPServer(addr)

	// Запуск в отдельной горутине
	go func() {
		log.Println("Server started on", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Ждём отмены контекста
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Таймаут на graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("Server exited gracefully")
	return nil
}
