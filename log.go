package wails

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Global logger reference
var logger = log.New()

// Fields is used by the customLogger object to output
// fields along with a message
type Fields map[string]interface{}

// Default options for the global logger
func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)
}

// Sets the log level to the given level
func setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	case "panic":
		logger.SetLevel(log.PanicLevel)
	default:
		logger.SetLevel(log.DebugLevel)
		logger.Warnf("Log level '%s' not recognised. Setting to Debug.", level)
	}
}
