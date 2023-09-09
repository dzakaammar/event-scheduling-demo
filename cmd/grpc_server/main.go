package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/dzakaammar/event-scheduling-example/cmd/pkg"
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/dzakaammar/event-scheduling-example/internal/scheduling"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

func main() {
	log.Fatalf("Error running grpc server: %v", run())
}

func run() error {
	slogHandler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(slogHandler)

	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	tp, err := pkg.InitTracerProvider(cfg.OTLPEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	dbConn, err := sqlx.Open("pgx", cfg.DbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = dbConn.Close()
	}()

	var repo core.EventRepository
	{
		repo = postgresql.NewEventRepository(dbConn)
		repo = postgresql.NewInstrumentation(repo)
	}

	var svc core.SchedulingService
	{
		svc = scheduling.NewService(repo)
		svc = scheduling.NewInstrumentation(svc)
	}

	grpcServer := pkg.NewGRPCServer(svc)
	waitForSignal := pkg.GracefulShutdown(func() error {
		return grpcServer.Start(cfg.GRPCAddress)
	})

	waitForSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return grpcServer.Stop(ctx)
}
