package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Global logger reference
var logger = logrus.New()

// Fields is used by the customLogger object to output
// fields along with a message
type Fields map[string]interface{}

// Default options for the global logger
func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
}

// SetLogLevel sets the log level to the given level
func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
		logger.Warnf("Log level '%s' not recognised. Setting to Debug.", level)
	}
}
