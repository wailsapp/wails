package logger

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestByteBufferLogger(t *testing.T) {

	is := is.New(t)

	// Create new byte buffer logger
	var buf bytes.Buffer

	myLogger := New(&buf)
	myLogger.SetLogLevel(TRACE)

	tests := map[uint8]string{
		TRACE:   "TRACE | I am a message!\n",
		DEBUG:   "DEBUG | I am a message!\n",
		WARNING: "WARN  | I am a message!\n",
		INFO:    "INFO  | I am a message!\n",
		ERROR:   "ERROR | I am a message!\n",
	}

	methods := map[uint8]func(string, ...interface{}) error{
		TRACE:   myLogger.Trace,
		DEBUG:   myLogger.Debug,
		WARNING: myLogger.Warning,
		INFO:    myLogger.Info,
		ERROR:   myLogger.Error,
	}

	for level, expected := range tests {

		buf.Reset()

		method := methods[level]

		// Write message
		err := method("I am a message!")
		if err != nil {
			panic(err)
		}
		actual := buf.String()

		is.Equal(actual, expected)
	}

}
func TestCustomLogger(t *testing.T) {

	is := is.New(t)

	// Create new byte buffer logger
	var buf bytes.Buffer

	myLogger := New(&buf)
	myLogger.SetLogLevel(TRACE)
	customLogger := myLogger.CustomLogger("Test")

	tests := map[uint8]string{
		TRACE:   "TRACE | Test | I am a message!\n",
		DEBUG:   "DEBUG | Test | I am a message!\n",
		WARNING: "WARN  | Test | I am a message!\n",
		INFO:    "INFO  | Test | I am a message!\n",
		ERROR:   "ERROR | Test | I am a message!\n",
	}

	methods := map[uint8]func(string, ...interface{}) error{
		TRACE:   customLogger.Trace,
		DEBUG:   customLogger.Debug,
		WARNING: customLogger.Warning,
		INFO:    customLogger.Info,
		ERROR:   customLogger.Error,
	}

	for level, expected := range tests {

		buf.Reset()

		method := methods[level]

		// Write message
		err := method("I am a message!")
		if err != nil {
			panic(err)
		}
		actual := buf.String()

		is.Equal(actual, expected)
	}

}
func TestWriteln(t *testing.T) {

	is := is.New(t)

	// Create new byte buffer logger
	var buf bytes.Buffer

	myLogger := New(&buf)
	myLogger.SetLogLevel(DEBUG)

	buf.Reset()

	// Write message
	err := myLogger.Writeln("I am a message!")
	if err != nil {
		panic(err)
	}
	actual := buf.String()

	is.Equal(actual, "I am a message!\n")

	buf.Reset()

	// Write message
	err = myLogger.Write("I am a message!")
	if err != nil {
		panic(err)
	}
	actual = buf.String()

	is.Equal(actual, "I am a message!")

}

func TestLogLevel(t *testing.T) {

	is := is.New(t)

	// Create new byte buffer logger
	var buf bytes.Buffer

	myLogger := New(&buf)
	myLogger.SetLogLevel(ERROR)

	tests := map[uint8]string{
		TRACE:   "",
		DEBUG:   "",
		WARNING: "",
		INFO:    "",
		ERROR:   "ERROR | I am a message!\n",
	}

	methods := map[uint8]func(string, ...interface{}) error{
		TRACE:   myLogger.Trace,
		DEBUG:   myLogger.Debug,
		WARNING: myLogger.Warning,
		INFO:    myLogger.Info,
		ERROR:   myLogger.Error,
	}

	for level := range tests {

		method := methods[level]

		// Write message
		err := method("I am a message!")
		if err != nil {
			panic(err)
		}

	}
	actual := buf.String()

	is.Equal(actual, "ERROR | I am a message!\n")
}

func TestFileLogger(t *testing.T) {

	is := is.New(t)

	// Create new byte buffer logger
	file, err := ioutil.TempFile(".", "wailsv2test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	myLogger := New(file)
	myLogger.SetLogLevel(DEBUG)

	// Write message
	err = myLogger.Info("I am a message!")
	if err != nil {
		panic(err)
	}
	actual, err := ioutil.ReadFile(file.Name())
	if err != nil {
		panic(err)
	}

	is.Equal(string(actual), "INFO  | I am a message!\n")

}
