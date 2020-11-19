package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Server interface is mandated to be implemented by any of the servers running
type Server interface {
	Start(context.Context)
	CleanUp(context.Context) error
}

// AwaitTermination makes the program wait for the signal termination
// Valid signal termination (SIGINT, SIGTERM)
func AwaitTermination() chan os.Signal {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM)
	return interruptSignal
}
