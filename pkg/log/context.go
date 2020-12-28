package log

import (
	"context"
	"os"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
)

var doOnce sync.Once
var entry *logrus.Entry

func Init() *logrus.Entry {
	doOnce.Do(func() {
		l := logrus.New()
		entry = logrus.NewEntry(l)
		entry.Logger.Out = os.Stdout
		entry.Logger.Level = logrus.DebugLevel
		entry.Logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: RFC3339NanoFixed,
		})

	})
	return entry
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
		return Init()
	}
	fields := logrus.Fields{
		RequestIDHeader: GetReqID(ctx),
	}
	return logger.WithFields(fields)
}
