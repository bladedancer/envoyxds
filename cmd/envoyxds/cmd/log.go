package cmd

import (
	"errors"
	"time"

	"github.com/bladedancer/envoyxds/pkg/envoyxds"
	"github.com/sirupsen/logrus"
)

const (
	lineFormat = "line"
	jsonFormat = "json"
	logPackage = "package"
)

var log logrus.FieldLogger = logrus.StandardLogger()

func getFormatter(format string) (logrus.Formatter, error) {
	switch format {
	case lineFormat:
		return &logrus.TextFormatter{TimestampFormat: time.RFC3339}, nil
	case jsonFormat:
		return &logrus.JSONFormatter{TimestampFormat: time.RFC3339}, nil
	default:
		return nil, errors.New("[sma] invalid log format")
	}
}

// setupLogging sets up logging for each used package
func setupLogging(level string, format string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	formatter, err := getFormatter(format)

	if err != nil {
		return err
	}

	logger := logrus.New()

	logger.Level = lvl
	logger.Formatter = formatter

	log = logger.WithField(logPackage, "cmd")

	envoyxds.SetLog(logger.WithField(logPackage, "envoyxds"))

	return nil
}
