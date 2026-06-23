package context

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewAppContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM,
	)
}
