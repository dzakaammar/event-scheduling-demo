package cmd

import (
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	grpGatewayServer, err := server.NewGRPCGatewayServer(*grpcServerTarget)
	if err != nil {
		panic(err)
	}

	// TODO: handle graceful shutdown
	return grpGatewayServer.Start(cfg.GRPCGatewayAddress)
}
