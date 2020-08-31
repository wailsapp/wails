package main

import (
	"math/rand"

	"github.com/wailsapp/wails"
)

// Counter is what we use for counting
type Counter struct {
	r     *wails.Runtime
	store *wails.Store
}

// WailsInit is called when the component is being initialised
func (c *Counter) WailsInit(runtime *wails.Runtime) error {
	c.r = runtime
	c.store = runtime.Store.New("Counter", 0)
	return nil
}

// RandomValue sets the counter to a random value
func (c *Counter) RandomValue() {
	c.store.Set(rand.Intn(1000))
}

// Increment will increment the counter
func (c *Counter) Increment() {

	increment := func(data int) int {
		return data + 1
	}

	// Update the store using the increment function
	c.store.Update(increment)
}

// Decrement will decrement the counter
func (c *Counter) Decrement() {

	decrement := func(data int) int {
		return data - 1
	}
	// Update the store using the decrement function
	c.store.Update(decrement)
}
