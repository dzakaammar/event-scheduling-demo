package server

import (
	"context"
	"embed"
	"net/http"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	//go:embed swagger
	swagger embed.FS

	//go:embed gen/v1/openapiv2.swagger.yaml
	openAPIFile []byte
)

type GRPCGatewayServer struct {
	srv *http.Server
}

func NewGRPCGatewayServer(grpcTarget string) (*GRPCGatewayServer, error) {
	gatewayHandler, err := grpcGatewayHandler(grpcTarget)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Mount("/api", gatewayHandler)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(swagger))))
	r.Get("/openapiv2.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openAPIFile)
	})

	return &GRPCGatewayServer{
		srv: &http.Server{
			Handler: r,
		},
	}, nil
}

func (g *GRPCGatewayServer) Start(address string) error {
	g.srv.Addr = address
	return g.srv.ListenAndServe()
}

func (g *GRPCGatewayServer) Stop(ctx context.Context) error {
	return g.srv.Shutdown(ctx)
}

func grpcGatewayHandler(grpcServerTarget string) (http.Handler, error) {
	handler := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	err := v1.RegisterAPIHandlerFromEndpoint(context.Background(), handler, grpcServerTarget, opts)
	if err != nil {
		return nil, err
	}

	return handler, nil
}