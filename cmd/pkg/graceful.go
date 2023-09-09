package pkg

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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
		err := <-errCh
		slog.Error(err.Error())
	}
}
