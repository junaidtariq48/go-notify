package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	// Set the output to standard out
	Logger.Out = os.Stdout

	// Set the log level to Info by default, change it to Debug if needed
	Logger.SetLevel(logrus.InfoLevel)

	// Set the formatter to JSON (this makes logs easier to analyze in systems like ELK)
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
