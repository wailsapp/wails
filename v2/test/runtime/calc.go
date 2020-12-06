package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2"
)

// Calc is a calculator
type Calc struct {
	name    string
	runtime *wails.Runtime
}

func newCalc(name string) *Calc {
	return &Calc{
		name: name,
	}
}

// Name will return the name of the calculator
func (c *Calc) Name() string {
	return c.name
}

// Add will add the 2 given integers and return the result
func (c *Calc) Add(a, b int) int {
	return a + b
}

func (c *Calc) unexported() int {
	return 1
}

func (c *Calc) Mult(a, b int) int {
	return a * b
}

func (c *Calc) Divide(a, b int) (int, error) {
	if b == 0 {
		return -1, fmt.Errorf("divide by zero")
	}
	return a / b, nil
}
