package server

import (
	"context"
	"fmt"
	"net"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	srv *grpc.Server
}

func NewGRPCServer(endpoint v1.EventServiceServer) *GRPCServer {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	)
	v1.RegisterEventServiceServer(srv, endpoint)
	reflection.Register(srv)

	return &GRPCServer{
		srv: srv,
	}
}

func (g *GRPCServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

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
