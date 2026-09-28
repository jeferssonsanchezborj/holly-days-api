// Package logger centralizes the construction of the application's
// structured logger (logrus), so every component shares the same
// configuration (format, level, output).
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New builds a configured *logrus.Logger. The log level can be controlled
// via the LOG_LEVEL environment variable (debug, info, warn, error);
// it defaults to "info". Logs are written as JSON to stdout, which is the
// recommended practice for containerized services.
func New() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	return log
}

