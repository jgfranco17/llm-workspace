package logging

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

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
