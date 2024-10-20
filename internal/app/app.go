package app

import (
	appgrpc "gRpcAuthService/internal/app/grpc"
	"gRpcAuthService/internal/services/auth"
	"gRpcAuthService/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *appgrpc.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := appgrpc.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
