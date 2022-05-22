package cmd

import (
	"context"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	"github.com/dzakaammar/event-scheduling-example/internal/instrumentation"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/dzakaammar/event-scheduling-example/internal/server"
	"github.com/dzakaammar/event-scheduling-example/internal/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func runGRPCServer(_ *cobra.Command, _ []string) error {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	tp, err := tracerProvider(cfg.OTLPEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	dbConn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DbSource,
	}))
	if err != nil {
		log.Fatal(err)
	}

	var repo internal.EventRepository
	{
		repo = postgresql.NewEventRepository(dbConn)
		repo = instrumentation.NewEventRepository(repo)
	}

	var svc internal.EventService
	{
		svc = service.NewEventService(repo)
		svc = instrumentation.NewEventService(svc)
	}

	grpcEndpoint := endpoint.NewGRPCEndpoint(svc)
	grpcServer := server.NewGRPCServer(grpcEndpoint)

	waitForSignal := gracefulShutdown(func() error {
		return grpcServer.Start(cfg.GRPCAddress)
	})

	waitForSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return grpcServer.Stop(ctx)
}

func tracerProvider(agentAddr string) (*tracesdk.TracerProvider, error) {
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(agentAddr),
	)

	traceExp, err := otlptrace.New(context.Background(), traceClient)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("event-scheduling-demo"),
		),
	)
	if err != nil {
		return nil, err
	}

	bsp := tracesdk.NewBatchSpanProcessor(traceExp)
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithSpanProcessor(bsp),
		tracesdk.WithResource(res),
	)

	return tp, nil
}
