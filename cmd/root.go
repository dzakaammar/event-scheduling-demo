package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: runGRPCServer,
}

// Execute :nodoc:
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func gracefulShutdown(fn func() error) func() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		if err := fn(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		<-sigCh
		errCh <- fmt.Errorf("interrupted")
	}()

	return func() {
		defer close(errCh)
		log.Error(<-errCh)
	}
}
