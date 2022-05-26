package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func GracefulShutdown(fn func() error) func() {
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
