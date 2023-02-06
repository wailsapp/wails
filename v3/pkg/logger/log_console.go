package logger

import "fmt"

type Console struct{}

func (l *Console) Name() string {
	return "Console"
}

func (l *Console) Log(message *Message) {
	msg := fmt.Sprintf(message.Message+"\n", message.Data...)
	level := ""
	if message.Level != "" {
		level = fmt.Sprintf("[%s] ", message.Level)
	}
	sender := ""
	if message.Sender != "" {
		sender = fmt.Sprintf("%s: ", message.Sender)
	}

	fmt.Printf("%s%s%s", level, sender, msg)
}

func (l *Console) Flush() error {
	return nil
}
