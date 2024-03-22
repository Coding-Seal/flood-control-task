package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"strconv"
	"task/floodControl"
	"task/floodControl/redisDB"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load .env", slog.Any("error", err))
		os.Exit(1)
	}
	timeLimit, err := time.ParseDuration(os.Getenv("TIME_LIMIT"))
	if err != nil {
		slog.Error("failed to parse TIME_LIMIT", slog.Any("error", err))
		os.Exit(1)
	}
	requestLimit, err := strconv.ParseInt(os.Getenv("REQUEST_LIMIT"), 10, 64)
	if err != nil {
		slog.Error("failed to parse REQUEST_LIMIT", slog.Any("error", err))
		os.Exit(1)
	}
	rdbAddr := os.Getenv("REDIS_DB_ADDR")

	slog.Info("Connecting to db", slog.String("addres", rdbAddr))
	client := redis.NewClient(&redis.Options{
		Addr:     rdbAddr,
		Password: "", // no password set
		DB:       0,  // use default DB})
	})
	err = client.Ping(context.Background()).Err()
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
	}
	var fc FloodControl = floodControl.New(redisDB.New(client), timeLimit, requestLimit)
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
	slog.Info("stopping program")
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
