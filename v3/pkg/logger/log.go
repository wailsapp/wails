package logger

import (
	"fmt"
)

type Logger struct {
	output []Output
}

func New(outputs ...Output) *Logger {
	result := &Logger{}
	if outputs != nil {
		result.output = outputs
	}
	return result
}

func (l *Logger) AddOutput(output Output) {
	l.output = append(l.output, output)
}

func (l *Logger) Log(message *Message) {
	for _, o := range l.output {
		go o.Log(message)
	}
}

func (l *Logger) Flush() {
	for _, o := range l.output {
		if err := o.Flush(); err != nil {
			fmt.Printf("Error flushing '%s' Logger: %s\n", o.Name(), err.Error())
		}
	}
}
