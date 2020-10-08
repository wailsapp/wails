package servicebus

import (
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/wailsapp/wails/v2/internal/logger"
)

type Person interface {
	FullName() string
}

type person struct {
	Firstname string
	Lastname  string
}

func newPerson(firstname string, lastname string) *person {
	result := &person{}
	result.Firstname = firstname
	result.Lastname = lastname
	return result
}

func (p *person) FullName() string {
	return p.Firstname + " " + p.Lastname
}

func TestSingleTopic(t *testing.T) {

	is := is.New(t)

	var expected string = "I am a message!"
	var actual string

	var wg sync.WaitGroup

	// Create new bus
	bus := New(logger.New())
	messageChannel, _ := bus.Subscribe("hello")

	wg.Add(1)
	go func() {
		message := <-messageChannel
		actual = message.Data().(string)
		wg.Done()
	}()

	bus.Start()
	bus.Publish("hello", "I am a message!")
	wg.Wait()
	bus.Stop()

	is.Equal(actual, expected)

}
func TestMultipleTopics(t *testing.T) {

	is := is.New(t)

	var hello string
	var world string
	var expected string = "Hello World!"

	var wg sync.WaitGroup

	// Create new bus
	bus := New(logger.New())

	// Create subscriptions
	helloChannel, _ := bus.Subscribe("hello")
	worldChannel, _ := bus.Subscribe("world")

	wg.Add(1)
	go func() {
		counter := 2
		for counter > 0 {
			select {
			case helloMessage := <-helloChannel:
				hello = helloMessage.Data().(string)
				counter--
			case worldMessage := <-worldChannel:
				world = worldMessage.Data().(string)
				counter--
			}
		}
		wg.Done()
	}()

	bus.Start()
	bus.Publish("hello", "Hello ")
	bus.Publish("world", "World!")
	wg.Wait()
	bus.Stop()

	is.Equal(hello+world, expected)
}

func TestSingleTopicWildcard(t *testing.T) {

	is := is.New(t)

	var expected string = "I am a message!"
	var actual string

	var wg sync.WaitGroup

	// Create new bus
	bus := New(logger.New())
	messageChannel, _ := bus.Subscribe("hello")

	wg.Add(1)
	go func() {
		message := <-messageChannel
		actual = message.Data().(string)
		wg.Done()
	}()

	bus.Start()
	bus.Publish("hello:wildcard:test", "I am a message!")
	wg.Wait()
	bus.Stop()

	is.Equal(actual, expected)

}
func TestMultipleTopicsWildcard(t *testing.T) {

	is := is.New(t)

	var hello string
	var world string
	var expected string = "Hello World!"

	var wg sync.WaitGroup

	// Create new bus
	bus := New(logger.New())
	helloChannel, _ := bus.Subscribe("hello")
	worldChannel, _ := bus.Subscribe("world")

	wg.Add(1)
	go func() {
		counter := 2
		for counter > 0 {
			select {
			case helloMessage := <-helloChannel:
				hello = helloMessage.Data().(string)
				counter--
			case worldMessage := <-worldChannel:
				world = worldMessage.Data().(string)
				counter--
			}
		}
		wg.Done()
	}()

	bus.Start()
	bus.Publish("hello:wildcard:test", "Hello ")
	bus.Publish("world:wildcard:test", "World!")
	wg.Wait()
	bus.Stop()

	is.Equal(hello+world, expected)
}

func TestStructData(t *testing.T) {

	is := is.New(t)

	var expected string = "Tom Jones"
	var actual string

	var wg sync.WaitGroup

	// Create new bus
	bus := New(logger.New())
	messageChannel, _ := bus.Subscribe("person")

	wg.Add(1)
	go func() {
		message := <-messageChannel
		p := message.Data().(*person)
		actual = p.FullName()
		wg.Done()
	}()

	bus.Start()
	bus.Publish("person", newPerson("Tom", "Jones"))
	wg.Wait()
	bus.Stop()

	is.Equal(actual, expected)

}

func TestErrors(t *testing.T) {

	is := is.New(t)

	// Create new bus
	bus := New(logger.New())

	_, err := bus.Subscribe("person")
	is.NoErr(err)

	err = bus.Start()
	is.NoErr(err)

	err = bus.Publish("person", newPerson("Tom", "Jones"))
	is.NoErr(err)

	err = bus.Stop()
	is.NoErr(err)

	err = bus.Stop()
	is.True(err != nil)

	err = bus.Start()
	is.True(err != nil)

	_, err = bus.Subscribe("person")
	is.True(err != nil)

	err = bus.Publish("person", newPerson("Tom", "Jones"))
	is.True(err != nil)

}
