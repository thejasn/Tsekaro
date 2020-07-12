package log

import (
	"context"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

var (
	// G is an alias for GetLogger.
	//
	// We may want to define this locally to a package to get package tagged log
	// messages.
	G = GetLogger

	// L is an alias for the standard logger.
	L = logrus.NewEntry(logrus.New())
)

func init() {
	L.Logger.Out = os.Stdout
	L.Logger.Level = logrus.DebugLevel
	L.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: RFC3339NanoFixed,
	})
}

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

// WithLogger returns a new context with the provided logger. Using with ctx logrus
// to append required metadata from interceptors
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return ctxlogrus.ToContext(ctx, logger)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctxlogrus.Extract(ctx)

	if logger == nil {
		return L
	}

	return logger
}
