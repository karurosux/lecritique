package logger

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	
	// Set output to stdout
	log.SetOutput(os.Stdout)
	
	// Set log level
	log.SetLevel(logrus.InfoLevel)
	
	// Use text formatter for readable logs
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}


// Get returns the configured logger instance
func Get() *logrus.Logger {
	return log
}

// Info logs an info message
func Info(msg string, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{})
	if len(fields) > 0 {
		entry = log.WithFields(fields[0])
	}
	entry.Info(msg)
}

// Error logs an error message
func Error(msg string, err error, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{"error": err.Error()})
	if len(fields) > 0 {
		for k, v := range fields[0] {
			entry = entry.WithField(k, v)
		}
	}
	entry.Error(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{})
	if len(fields) > 0 {
		entry = log.WithFields(fields[0])
	}
	entry.Warn(msg)
}

// Debug logs a debug message
func Debug(msg string, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{})
	if len(fields) > 0 {
		entry = log.WithFields(fields[0])
	}
	entry.Debug(msg)
}

// LogJSON logs an object as formatted JSON for debugging
func LogJSON(msg string, obj interface{}) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.WithField("error", err).Error("Failed to marshal object to JSON")
		return
	}
	log.WithField("json", string(jsonBytes)).Info(msg)
}

// SetLevel sets the log level
func SetLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}