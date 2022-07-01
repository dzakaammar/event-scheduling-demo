package pkg

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var ErrInterrupted = errors.New("interrupted")

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
		errCh <- ErrInterrupted
	}()

	return func() {
		defer close(errCh)
		log.Error(<-errCh)
	}
}
