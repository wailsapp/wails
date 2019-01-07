package wails

type CustomLogger struct {
	prefix string
}

func newCustomLogger(prefix string) *CustomLogger {
	return &CustomLogger{
		prefix: "[" + prefix + "] ",
	}
}

func (c *CustomLogger) Info(message string) {
	logger.Info(c.prefix + message)
}

func (c *CustomLogger) Infof(message string, args ...interface{}) {
	logger.Infof(c.prefix+message, args...)
}

func (c *CustomLogger) InfoFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Info(c.prefix + message)
}

func (c *CustomLogger) Debug(message string) {
	logger.Debug(c.prefix + message)
}

func (c *CustomLogger) Debugf(message string, args ...interface{}) {
	logger.Debugf(c.prefix+message, args...)
}

func (c *CustomLogger) DebugFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Debug(c.prefix + message)
}

func (c *CustomLogger) Warn(message string) {
	logger.Warn(c.prefix + message)
}

func (c *CustomLogger) Warnf(message string, args ...interface{}) {
	logger.Warnf(c.prefix+message, args...)
}

func (c *CustomLogger) WarnFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Warn(c.prefix + message)
}

func (c *CustomLogger) Error(message string) {
	logger.Error(c.prefix + message)
}

func (c *CustomLogger) Errorf(message string, args ...interface{}) {
	logger.Errorf(c.prefix+message, args...)
}

func (c *CustomLogger) ErrorFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Error(c.prefix + message)
}

func (c *CustomLogger) Fatal(message string) {
	logger.Fatal(c.prefix + message)
}

func (c *CustomLogger) Fatalf(message string, args ...interface{}) {
	logger.Fatalf(c.prefix+message, args...)
}

func (c *CustomLogger) FatalFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Fatal(c.prefix + message)
}
func (c *CustomLogger) Panic(message string) {
	logger.Panic(c.prefix + message)
}

func (c *CustomLogger) Panicf(message string, args ...interface{}) {
	logger.Panicf(c.prefix+message, args...)
}

func (c *CustomLogger) PanicFields(message string, fields Fields) {
	logger.WithFields(map[string]interface{}(fields)).Panic(c.prefix + message)
}
