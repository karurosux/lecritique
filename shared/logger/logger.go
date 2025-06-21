package logger

import (
	"encoding/json"
	"fmt"
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
	
	// Custom JSON formatter for better readability
	log.SetFormatter(&PrettyJSONFormatter{})
}

// PrettyJSONFormatter formats logs as indented JSON
type PrettyJSONFormatter struct{}

func (f *PrettyJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	
	// Copy entry data
	for k, v := range entry.Data {
		data[k] = v
	}
	
	// Add standard fields
	data["timestamp"] = entry.Time.Format("2006-01-02T15:04:05.000Z07:00")
	data["level"] = entry.Level.String()
	data["message"] = entry.Message
	
	// Marshal with indentation for readability
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON: %w", err)
	}
	
	return append(jsonBytes, '\n'), nil
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