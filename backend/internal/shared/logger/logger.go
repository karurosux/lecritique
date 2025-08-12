package logger

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	
	log.SetOutput(os.Stdout)
	
	log.SetLevel(logrus.InfoLevel)
	
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}


func Get() *logrus.Logger {
	return log
}

func Info(msg string, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{})
	if len(fields) > 0 {
		entry = log.WithFields(fields[0])
	}
	entry.Info(msg)
}

func Error(msg string, err error, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{"error": err.Error()})
	if len(fields) > 0 {
		for k, v := range fields[0] {
			entry = entry.WithField(k, v)
		}
	}
	entry.Error(msg)
}

func Warn(msg string, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{})
	if len(fields) > 0 {
		entry = log.WithFields(fields[0])
	}
	entry.Warn(msg)
}

func Debug(msg string, fields ...logrus.Fields) {
	entry := log.WithFields(logrus.Fields{})
	if len(fields) > 0 {
		entry = log.WithFields(fields[0])
	}
	entry.Debug(msg)
}

func LogJSON(msg string, obj interface{}) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.WithField("error", err).Error("Failed to marshal object to JSON")
		return
	}
	log.WithField("json", string(jsonBytes)).Info(msg)
}

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