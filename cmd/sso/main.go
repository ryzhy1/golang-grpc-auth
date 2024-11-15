package main

import (
	"AuthService/internal/app"
	"AuthService/internal/config"
	"github.com/joho/godotenv"
	l "log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	envDev   = "dev"
	envProd  = "prod"
	envLocal = "local"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		l.Fatalf("Error loading .env file")
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting account service", "env", cfg.Env)

	application := app.New(log, strconv.Itoa(cfg.GRPC.AuthPort), cfg.Storage, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("Application stopped", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
