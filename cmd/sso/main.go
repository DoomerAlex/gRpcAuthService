package main

import (
	"gRpcAuthService/internal/app"
	"gRpcAuthService/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()
	// Инициализация логгера
	log := setupLogger(cfg.Env)
	// Инициализация приложения
	log.Info("starting application", slog.String("env", cfg.Env))
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	// Запуск gRPC сервера
	go application.GRPCSrv.MustRun()
	// Начинаем отслеживать сигналы на остановку приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	// Остановка gRPC сервера
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("Undefined env: " + env)
	}
}
