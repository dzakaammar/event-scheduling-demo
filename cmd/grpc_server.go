package cmd

import (
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/dzakaammar/event-scheduling-example/internal/server"
	"github.com/dzakaammar/event-scheduling-example/internal/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func runGRPCServer(_ *cobra.Command, _ []string) error {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DbSource,
	}))
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
