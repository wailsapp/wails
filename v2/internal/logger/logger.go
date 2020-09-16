package logger

const (
	// TRACE level
	TRACE uint8 = 0

	// DEBUG level logging
	DEBUG uint8 = 1

	// INFO level logging
	INFO uint8 = 2

	// WARNING level logging
	WARNING uint8 = 4

	// ERROR level logging
	ERROR uint8 = 8

	// FATAL level logging
	FATAL uint8 = 16

	// BYPASS level logging - does not use a log level
	BYPASS uint8 = 255
)

var mapLogLevel = map[uint8]string{
	TRACE:   "TRACE | ",
	DEBUG:   "DEBUG | ",
	INFO:    "INFO  | ",
	WARNING: "WARN  | ",
	ERROR:   "ERROR | ",
	FATAL:   "FATAL | ",
	BYPASS:  "",
}
