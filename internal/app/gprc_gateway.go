package app

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCGatewayServer struct {
	srv *http.Server
}

func NewGRPCGatewayServer(grpcTarget string, swagger fs.FS, openAPIYAMLFile []byte) (*GRPCGatewayServer, error) {
	gatewayHandler, err := grpcGatewayHandler(grpcTarget)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(cors.AllowAll().Handler)

	r.Group(func(api chi.Router) {
		api.Use(otelhttp.NewMiddleware("http-server"))
		api.Mount("/api", gatewayHandler)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(swagger))))
	r.Get("/openapiv2.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openAPIYAMLFile)
	})

	return &GRPCGatewayServer{
		srv: &http.Server{
			Handler:           r,
			ReadHeaderTimeout: 200 * time.Millisecond,
		},
	}, nil
}

func (g *GRPCGatewayServer) Start(address string) error {
	g.srv.Addr = address
	fmt.Println("grpc gateway running on ", address)
	return g.srv.ListenAndServe()
}

func (g *GRPCGatewayServer) Stop(ctx context.Context) error {
	return g.srv.Shutdown(ctx)
}

func grpcGatewayHandler(grpcServerTarget string) (http.Handler, error) {
	handler := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				otelgrpc.UnaryClientInterceptor(),
			),
		),
		grpc.WithBlock(),
	}

	err := v1.RegisterAPIHandlerFromEndpoint(context.Background(), handler, grpcServerTarget, opts)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
