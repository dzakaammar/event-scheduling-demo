package cmd

import (
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	grpcGatewayCmd = &cobra.Command{
		Use:  "grpc-gateway",
		RunE: runGRPCGateway,
	}
	grpcServerTarget *string
)

func init() {
	grpcServerTarget = grpcGatewayCmd.PersistentFlags().StringP("target", "t", "localhost:8000", "The address of grpc server")
	rootCmd.AddCommand(grpcGatewayCmd)
}

func runGRPCGateway(cmd *cobra.Command, args []string) error {
	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	if grpcServerTarget == nil {
		grpcServerTarget = &cfg.GRPCAddress
	}

	grpcGatewayServer, err := server.NewGRPCGatewayServer(*grpcServerTarget)
	if err != nil {
		panic(err)
	}

	waitForSignal := gracefulShutdown(func() error {
		return grpcGatewayServer.Start(cfg.GRPCGatewayAddress)
	})

	waitForSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return grpcGatewayServer.Stop(ctx)
}
