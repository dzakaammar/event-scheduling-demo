package pkg

import (
	"context"
	"fmt"
	"net"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

var logrusEntry = logrus.NewEntry(logrus.StandardLogger())

type GRPCServer struct {
	srv *grpc.Server
}

func NewGRPCServer(schedulingSvc core.SchedulingService) *GRPCServer {
	endpoint := endpoint.NewGRPCEndpoint(schedulingSvc)

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

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			otelgrpc.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		),
	)
	v1.RegisterAPIServer(srv, endpoint)
	reflection.Register(srv)

	return &GRPCServer{
		srv: srv,
	}
}

func (g *GRPCServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	fmt.Println("grpc server is running on ", address)
	return g.srv.Serve(lis)
}

func (g *GRPCServer) Stop(ctx context.Context) error {
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
