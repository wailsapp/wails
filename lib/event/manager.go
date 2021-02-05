package event

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/messages"
)

// Manager handles and processes events
type Manager struct {
	incomingEvents chan *messages.EventData
	listeners      map[string][]*eventListener
	log            *logger.CustomLogger
	mu             sync.Mutex
	quitChannel    chan struct{}
	renderer       interfaces.Renderer // Messages will be dispatched to the frontend
	running        bool
	wg             sync.WaitGroup
}

// NewManager creates a new event manager with a 100 event buffer
func NewManager() interfaces.EventManager {
	return &Manager{
		incomingEvents: make(chan *messages.EventData, 100),
		listeners:      make(map[string][]*eventListener),
		log:            logger.NewCustomLogger("Events"),
		quitChannel:    make(chan struct{}, 1),
		running:        false,
	}
}

// eventListener holds a callback function which is invoked when
// the event being listened for is emitted. It has a counter which
// indicates the total number of events to allow. A value of zero
// means it does not expire (default).
type eventListener struct {
	callback func(...interface{}) // Function to call with emitted event data
	counter  uint                 // Expire after counter callbacks. 0 = infinite
}

// Creates a new event listener from the given callback function
func (e *Manager) addEventListener(eventName string, callback func(...interface{}), counter uint) error {

	// Sanity check inputs
	if callback == nil {
		return fmt.Errorf("nil callback bassed to addEventListener")
	}

	// Check event has been registered before
	if e.listeners[eventName] == nil {
		e.listeners[eventName] = []*eventListener{}
	}

	// Create the callback
	listener := &eventListener{
		callback: callback,
		counter:  counter,
	}

	// Register listener
	e.listeners[eventName] = append(e.listeners[eventName], listener)

	// All good mate
	return nil
}

// Emit broadcasts the given event to the subscribed listeners
func (e *Manager) Emit(eventName string, optionalData ...interface{}) {
	e.incomingEvents <- &messages.EventData{Name: eventName, Data: optionalData}
}

// On adds a listener for the given event
func (e *Manager) On(eventName string, callback func(...interface{})) {
	// Add a persistent eventListener (counter = 0)
	err := e.addEventListener(eventName, callback, 0)
	if err != nil {
		e.log.Error(err.Error())
	}
}

// Once adds a listener that will remove after one callback
func (e *Manager) Once(eventName string, callback func(...interface{})) {
	// Add a persistent eventListener (counter = 0)
	err := e.addEventListener(eventName, callback, 1)
	if err != nil {
		e.log.Error(err.Error())
	}
}

// OnMultiple adds a listener that will trigger at most <counter> times.
func (e *Manager) OnMultiple(eventName string, callback func(...interface{}), counter uint) {
	err := e.addEventListener(eventName, callback, counter)
	if err != nil {
		e.log.Error(err.Error())
	}
}

// PushEvent places the given event on to the event queue
func (e *Manager) PushEvent(eventData *messages.EventData) {
	e.incomingEvents <- eventData
}

// Shutdown is called when exiting the Application
func (e *Manager) Shutdown() {
	e.log.Debug("Shutting Down")
	e.quitChannel <- struct{}{}
	e.log.Debug("Waiting for main loop to exit")
	e.wg.Wait()
}

// Start the event manager's queue processing
func (e *Manager) Start(renderer interfaces.Renderer) {

	e.log.Info("Starting")

	// Store renderer
	e.renderer = renderer

	// Set up waitgroup so we can wait for goroutine to quit
	e.running = true
	e.wg.Add(1)

	// Run main loop in separate goroutine
	go func() {
		e.log.Info("Listening")
		for e.running {
			// TODO: Listen for application exit
			select {
			case event := <-e.incomingEvents:
				e.log.DebugFields("Got Event", logger.Fields{
					"data": event.Data,
					"name": event.Name,
				})

				// Notify renderer
				err := e.renderer.NotifyEvent(event)
				if err != nil {
					e.log.Error(err.Error())
				}

				// Iterate listeners
				for currentIndex, listener := range e.listeners[event.Name] {

					// Call listener, perhaps with data
					if event.Data == nil {
						go listener.callback()
					} else {
						unpacked := event.Data.([]interface{})
						go listener.callback(unpacked...)
					}

					// Decrement counter if its a non-persistent listener
					if listener.counter > 0 {
						// fields about to change; enter Mutex
						e.mu.Lock()
						listener.counter--

						// expiration condition == counter WAS 1 but is NOW 0
						if listener.counter == 0 {
							// this listener has just expired; remove it NOW
							// see fast method: https://yourbasic.org/golang/delete-element-slice/
							// https://play.golang.org/p/j0YjKUN0NL1
							lastIndex := len(e.listeners[event.Name]) - 1

							// overwrite expired currentIndex listener with lastIndex listener
							e.listeners[event.Name][currentIndex] = e.listeners[event.Name][lastIndex]

							// remove the (now cloned) lastIndex listener from slice
							e.listeners[event.Name] = e.listeners[event.Name][:lastIndex]
						}
						e.mu.Unlock()
					}
				}

			case <-e.quitChannel:
				e.running = false
			}
		}
		e.wg.Done()
	}()
}
