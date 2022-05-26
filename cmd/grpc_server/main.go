package main

import (
	"context"
	"fmt"
	"net"
	"time"

	internalCmd "github.com/dzakaammar/event-scheduling-example/cmd/internal"
	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/dzakaammar/event-scheduling-example/internal/scheduling"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Fatalf("Error running grpc server: %v", run())
}

func run() error {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	tp, err := internalCmd.InitTracerProvider(cfg.OTLPEndpoint)
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

	var repo core.EventRepository
	{
		repo = postgresql.NewEventRepository(dbConn)
		repo = postgresql.NewInstrumentation(repo)
	}

	var svc core.SchedulingService
	{
		svc = scheduling.NewEventService(repo)
		svc = scheduling.NewInstrumentation(svc)
	}

	grpcEndpoint := endpoint.NewGRPCEndpoint(svc)
	grpcServer := newGRPCServer(grpcEndpoint)

	waitForSignal := internalCmd.GracefulShutdown(func() error {
		return grpcServer.Start(cfg.GRPCAddress)
	})

	waitForSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return grpcServer.Stop(ctx)
}

type grpcServer struct {
	srv *grpc.Server
}

func newGRPCServer(endpoint v1.APIServer) *grpcServer {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Level = logrus.ErrorLevel

	logrusEntry := logrus.NewEntry(logger)
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(func(code codes.Code) logrus.Level {
			switch code {
			case codes.DeadlineExceeded,
				codes.Unimplemented,
				codes.Unknown,
				codes.ResourceExhausted,
				codes.Unavailable,
				codes.Internal:
				return logrus.ErrorLevel
			default:
				return logrus.DebugLevel
			}
		}),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			otelgrpc.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		),
	)
	v1.RegisterAPIServer(srv, endpoint)
	reflection.Register(srv)

	return &grpcServer{
		srv: srv,
	}
}

func (g *grpcServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Println("grpc server is running on ", address)
	return g.srv.Serve(lis)
}

func (g *grpcServer) Stop(ctx context.Context) error {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		g.srv.GracefulStop()
	}()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		g.srv.Stop()
		return ctx.Err()
	}
}
