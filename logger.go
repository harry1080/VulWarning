package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func initLogger(filename string, level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	if level == logrus.DebugLevel || level == logrus.InfoLevel {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: false,
			TimestampFormat:        "2006-01-02 15:04:05",
		})
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(logFile)
		}
	}
	return logger
}
