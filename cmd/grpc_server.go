package cmd

import (
	"log"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/dzakaammar/event-scheduling-example/internal/server"
	"github.com/dzakaammar/event-scheduling-example/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func runGRPCServer(_ *cobra.Command, _ []string) error {
	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sqlx.Open(cfg.DbDriver, cfg.DbSource)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgresql.NewEventRepository(dbConn)
	svc := service.NewEventService(repo)

	grpcEndpoint := endpoint.NewGRPCEndpoint(svc)
	grpcServer := server.NewGRPCServer(grpcEndpoint)

	if err := grpcServer.Start(cfg.GRPCAddress); err != nil {
		return err
	}

	return nil
}
