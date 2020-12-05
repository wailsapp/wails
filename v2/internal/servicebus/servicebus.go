package servicebus

import (
	"fmt"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
)

// ServiceBus is a messaging bus for Wails applications
type ServiceBus struct {
	listeners    map[string][]chan *Message
	messageQueue chan *Message
	quitChannel  chan struct{}
	wg           sync.WaitGroup
	lock         sync.RWMutex
	closed       bool
	debug        bool
	logger       logger.CustomLogger
}

// New creates a new ServiceBus
// The internal message queue is set to 100 messages
// Listener queues are set to 10
func New(logger *logger.Logger) *ServiceBus {
	return &ServiceBus{
		listeners:    make(map[string][]chan *Message),
		messageQueue: make(chan *Message, 100),
		quitChannel:  make(chan struct{}, 1),
		logger:       logger.CustomLogger("Service Bus"),
	}
}

// dispatch the given message to the listeners
func (s *ServiceBus) dispatchMessage(message *Message) {

	// Lock to prevent additions to the listeners
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Iterate over listener's topics
	for topic := range s.listeners {

		// If the topic matches
		if strings.HasPrefix(message.Topic(), topic) {

			// Iterate over the listeners
			for _, callback := range s.listeners[topic] {

				// Process the message
				callback <- message
			}
		}
	}
}

// Debug puts the service bus into debug mode.
func (s *ServiceBus) Debug() {
	s.debug = true
}

// Start the service bus
func (s *ServiceBus) Start() error {

	s.logger.Trace("Starting")

	// Prevent starting when closed
	if s.closed {
		return fmt.Errorf("cannot call start on closed servicebus")
	}

	// We run in a different thread
	go func() {

		quit := false
		s.wg.Add(1)

		// Loop until we get a quit message
		for !quit {

			select {

			// Listen for messages
			case message := <-s.messageQueue:

				// Log message if in debug mode
				if s.debug {
					s.logger.Trace("Got message: { Topic: %s, Interface: %#v }", message.Topic(), message.Data())
				}
				// Dispatch message
				s.dispatchMessage(message)

			// Listen for quit messages
			case <-s.quitChannel:
				quit = true
			}
		}

		// Indicate we have shut down
		s.wg.Done()

	}()

	return nil
}

// Stop the service bus
func (s *ServiceBus) Stop() error {

	// Prevent subscribing when closed
	if s.closed {
		return fmt.Errorf("cannot call stop on closed servicebus")
	}

	s.closed = true

	// Send quit message
	s.quitChannel <- struct{}{}

	// Wait for dispatcher to stop
	s.wg.Wait()

	// Close down subscriber channels
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, subscribers := range s.listeners {
		for _, channel := range subscribers {
			close(channel)
		}
	}

	// Close message queue
	close(s.messageQueue)

	s.logger.Trace("Stopped")
	return nil
}

// UnSubscribe removes the listeners for the given topic (Use with caution!)
func (s *ServiceBus) UnSubscribe(topic string) {
	// Prevent any reads or writes to the listeners whilst
	// we create a new one
	s.lock.Lock()
	defer s.lock.Unlock()
	s.listeners[topic] = nil
}

// Subscribe is used to register a listener's interest in a topic
func (s *ServiceBus) Subscribe(topic string) (<-chan *Message, error) {

	// Prevent subscribing when closed
	if s.closed {
		return nil, fmt.Errorf("cannot call subscribe on closed servicebus")
	}

	// Prevent any reads or writes to the listeners whilst
	// we create a new one
	s.lock.Lock()
	defer s.lock.Unlock()

	// Append the new listener
	listener := make(chan *Message, 10)
	s.listeners[topic] = append(s.listeners[topic], listener)
	return (<-chan *Message)(listener), nil

}

// Publish sends the given message on the service bus
func (s *ServiceBus) Publish(topic string, data interface{}) {
	// Prevent publish when closed
	if s.closed {
		s.logger.Fatal("cannot call publish on closed servicebus")
		return
	}

	message := NewMessage(topic, data)
	s.messageQueue <- message
}

// PublishForTarget sends the given message on the service bus for the given target
func (s *ServiceBus) PublishForTarget(topic string, data interface{}, target string) {
	// Prevent publish when closed
	if s.closed {
		s.logger.Fatal("cannot call publish on closed servicebus")
		return
	}
	message := NewMessageForTarget(topic, data, target)
	s.messageQueue <- message
}
