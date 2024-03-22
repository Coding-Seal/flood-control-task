package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"task/floodControl"
	"task/floodControl/redisDB"
	"time"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.NewJSONHandler(os.Stderr, opts)
	slog.SetDefault(slog.New(logger))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("Connecting to db")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB})
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
	}
	var fc FloodControl = floodControl.New(redisDB.New(client), time.Second*5, 3)
	var userID int64 = 1

	for j := 0; j < 20; j++ {
		res, err := fc.Check(context.Background(), userID)
		if err != nil {
			slog.Error("internal error occurred", slog.Any("error", err))
		}
		if res {
			slog.Debug("access granted", slog.Int64("userID", userID))
		} else {
			slog.Debug("access denied", slog.Int64("userID", userID))
		}
		time.Sleep(time.Second)

	}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
