package logging

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// NewJSONLogger creates a logrus.Logger that writes pretty-printed JSON
// to w at the specified level. If level is invalid, it defaults to INFO
// and logs a warning.
func NewJSONLogger(w io.Writer, level string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(w)
	levelToUse, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Warnf("invalid log level %q, defaulting to INFO", level)
		levelToUse = logrus.InfoLevel
	}
	logger.SetLevel(levelToUse)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
		PrettyPrint:     true,
	})
	return logger
}
