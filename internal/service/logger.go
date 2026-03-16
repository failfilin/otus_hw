package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisLogsKey = "api:logs"

type LogEntry struct {
	Handler    string `json:"handler"`
	IP         string `json:"ip"`
	Status     int    `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	Timestamp  string `json:"timestamp"`
}

type RedisLogger struct {
	client *redis.Client
}

func NewRedisLogger() *RedisLogger {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis not available at %s: %v", addr, err)
	} else {
		log.Printf("Connected to Redis: %s", addr)
	}

	return &RedisLogger{client: client}
}

func (l *RedisLogger) Log(handler, ip string, status int, duration time.Duration) {
	entry := LogEntry{
		Handler:    handler,
		IP:         ip,
		Status:     status,
		DurationMs: duration.Milliseconds(),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("logger: failed to marshal log entry: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := l.client.LPush(ctx, redisLogsKey, string(data)).Err(); err != nil {
		log.Printf("logger: failed to write to Redis: %v", err)
	}
}

func (l *RedisLogger) Close() error {
	return l.client.Close()
}

func (l *RedisLogger) GetLogs(ctx context.Context, limit int64) ([]LogEntry, error) {
	results, err := l.client.LRange(ctx, redisLogsKey, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read logs from Redis: %w", err)
	}

	entries := make([]LogEntry, 0, len(results))
	for _, r := range results {
		var e LogEntry
		if err := json.Unmarshal([]byte(r), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}
