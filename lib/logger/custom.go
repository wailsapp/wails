package logger

// CustomLogger is a wrapper object to logrus
type CustomLogger struct {
	prefix    string
	errorOnly bool
}

// NewCustomLogger creates a new custom logger with the given prefix
func NewCustomLogger(prefix string) *CustomLogger {
	return &CustomLogger{
		prefix: "[" + prefix + "] ",
	}
}

// Info level message
func (c *CustomLogger) Info(message string) {
	GlobalLogger.Info(c.prefix + message)
}

// Infof - formatted message
func (c *CustomLogger) Infof(message string, args ...interface{}) {
	GlobalLogger.Infof(c.prefix+message, args...)
}

// InfoFields - message with fields
func (c *CustomLogger) InfoFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Info(c.prefix + message)
}

// Debug level message
func (c *CustomLogger) Debug(message string) {
	GlobalLogger.Debug(c.prefix + message)
}

// Debugf - formatted message
func (c *CustomLogger) Debugf(message string, args ...interface{}) {
	GlobalLogger.Debugf(c.prefix+message, args...)
}

// DebugFields - message with fields
func (c *CustomLogger) DebugFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Debug(c.prefix + message)
}

// Warn level message
func (c *CustomLogger) Warn(message string) {
	GlobalLogger.Warn(c.prefix + message)
}

// Warnf - formatted message
func (c *CustomLogger) Warnf(message string, args ...interface{}) {
	GlobalLogger.Warnf(c.prefix+message, args...)
}

// WarnFields - message with fields
func (c *CustomLogger) WarnFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Warn(c.prefix + message)
}

// Error level message
func (c *CustomLogger) Error(message string) {
	GlobalLogger.Error(c.prefix + message)
}

// Errorf - formatted message
func (c *CustomLogger) Errorf(message string, args ...interface{}) {
	GlobalLogger.Errorf(c.prefix+message, args...)
}

// ErrorFields - message with fields
func (c *CustomLogger) ErrorFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Error(c.prefix + message)
}

// Fatal level message
func (c *CustomLogger) Fatal(message string) {
	GlobalLogger.Fatal(c.prefix + message)
}

// Fatalf - formatted message
func (c *CustomLogger) Fatalf(message string, args ...interface{}) {
	GlobalLogger.Fatalf(c.prefix+message, args...)
}

// FatalFields - message with fields
func (c *CustomLogger) FatalFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Fatal(c.prefix + message)
}

// Panic level message
func (c *CustomLogger) Panic(message string) {
	GlobalLogger.Panic(c.prefix + message)
}

// Panicf - formatted message
func (c *CustomLogger) Panicf(message string, args ...interface{}) {
	GlobalLogger.Panicf(c.prefix+message, args...)
}

// PanicFields - message with fields
func (c *CustomLogger) PanicFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Panic(c.prefix + message)
}
