package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	internalCmd "github.com/dzakaammar/event-scheduling-example/cmd/internal"
	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal"
)

func main() {
	log.Fatalf("Error running grpc gateway: %v", run())
}

func run() error {
	var grpcServerTarget string

	flag.StringVar(&grpcServerTarget, "target", "", "The address of grpc server")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	if grpcServerTarget == "" {
		grpcServerTarget = cfg.GRPCAddress
	}

	grpcGatewayServer, err := newGRPCGatewayServer(grpcServerTarget)
	if err != nil {
		panic(err)
	}

	waitForSignal := internalCmd.GracefulShutdown(func() error {
		return grpcGatewayServer.Start(cfg.GRPCGatewayAddress)
	})

	waitForSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return grpcGatewayServer.Stop(ctx)
}

var (
	//go:embed swagger
	swagger embed.FS

	//go:embed gen/v1/openapiv2.swagger.yaml
	openAPIFile []byte
)

type gRPCGatewayServer struct {
	srv *http.Server
}

func newGRPCGatewayServer(grpcTarget string) (*gRPCGatewayServer, error) {
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

	return &gRPCGatewayServer{
		srv: &http.Server{
			Handler: r,
		},
	}, nil
}

func (g *gRPCGatewayServer) Start(address string) error {
	g.srv.Addr = address
	fmt.Println("grpc gateway running on ", address)
	return g.srv.ListenAndServe()
}

func (g *gRPCGatewayServer) Stop(ctx context.Context) error {
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
