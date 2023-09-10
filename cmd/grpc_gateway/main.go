package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/dzakaammar/event-scheduling-example/cmd/pkg"
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/app"
	"go.opentelemetry.io/otel"
)

var (
	//go:embed swagger
	swagger embed.FS

	//go:embed gen/v1/openapiv2.swagger.yaml
	openAPIFile []byte
)

func main() {
	log.Fatalf("Error running grpc gateway: %v", run())
}

func run() error {
	var grpcServerTarget string

	flag.StringVar(&grpcServerTarget, "target", "", "The address of grpc server")
	flag.Parse()

	slogHandler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(slogHandler)

	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	tp, err := pkg.InitTracerProvider(cfg.OTLPEndpoint, "grpc-gateway-server")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	if grpcServerTarget == "" {
		grpcServerTarget = cfg.GRPCAddress
	}

	grpcGatewayServer, err := app.NewGRPCGatewayServer(grpcServerTarget, swagger, openAPIFile)
	if err != nil {
		panic(err)
	}

	waitForSignal := pkg.GracefulShutdown(func() error {
		return grpcGatewayServer.Start(cfg.GRPCGatewayAddress)
	})

	waitForSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return grpcGatewayServer.Stop(ctx)
}
