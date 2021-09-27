package crypto

import (
	"crypto/rand"
	"fmt"
	"log"
)

// RandomID returns a random ID as a string
func RandomID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", b)
}
