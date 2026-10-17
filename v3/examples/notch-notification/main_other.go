//go:build !darwin

package main

import "fmt"

func main() {
	fmt.Println("The notch-notification example requires macOS.")
}
