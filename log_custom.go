package wails

// CustomLogger is a wrapper object to logrus
type CustomLogger struct {
	prefix    string
	errorOnly bool
}

func newCustomLogger(prefix string) *CustomLogger {
	return &CustomLogger{
		prefix: "[" + prefix + "] ",
	}
}

// Info level message
func (c *CustomLogger) Info(message string) {
	logger.Info(c.prefix + message)
}

// Infof - formatted message
func (c *CustomLogger) Infof(message string, args ...interface{}) {
	logger.Infof(c.prefix+message, args...)
}

// InfoFields - message with fields
func (c *CustomLogger) InfoFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Info(c.prefix + message)
}

// Debug level message
func (c *CustomLogger) Debug(message string) {
	logger.Debug(c.prefix + message)
}

// Debugf - formatted message
func (c *CustomLogger) Debugf(message string, args ...interface{}) {
	logger.Debugf(c.prefix+message, args...)
}

// DebugFields - message with fields
func (c *CustomLogger) DebugFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Debug(c.prefix + message)
}

// Warn level message
func (c *CustomLogger) Warn(message string) {
	logger.Warn(c.prefix + message)
}

// Warnf - formatted message
func (c *CustomLogger) Warnf(message string, args ...interface{}) {
	logger.Warnf(c.prefix+message, args...)
}

// WarnFields - message with fields
func (c *CustomLogger) WarnFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Warn(c.prefix + message)
}

// Error level message
func (c *CustomLogger) Error(message string) {
	logger.Error(c.prefix + message)
}

// Errorf - formatted message
func (c *CustomLogger) Errorf(message string, args ...interface{}) {
	logger.Errorf(c.prefix+message, args...)
}

// ErrorFields - message with fields
func (c *CustomLogger) ErrorFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Error(c.prefix + message)
}

// Fatal level message
func (c *CustomLogger) Fatal(message string) {
	logger.Fatal(c.prefix + message)
}

// Fatalf - formatted message
func (c *CustomLogger) Fatalf(message string, args ...interface{}) {
	logger.Fatalf(c.prefix+message, args...)
}

// FatalFields - message with fields
func (c *CustomLogger) FatalFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Fatal(c.prefix + message)
}

// Panic level message
func (c *CustomLogger) Panic(message string) {
	logger.Panic(c.prefix + message)
}

// Panicf - formatted message
func (c *CustomLogger) Panicf(message string, args ...interface{}) {
	logger.Panicf(c.prefix+message, args...)
}

// PanicFields - message with fields
func (c *CustomLogger) PanicFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Panic(c.prefix + message)
}
