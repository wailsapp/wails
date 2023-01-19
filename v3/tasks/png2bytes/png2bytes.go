package main

import (
	"bytes"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) != 2 {
		println("Please provide a filename")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		println("Error reading file:", err.Error())
		os.Exit(1)
	}

	var buffer bytes.Buffer
	buffer.WriteString("var image = []byte{")
	// Iterate over the bytes and print them out in decimal
	for _, b := range data {
		// convert byte to decimal
		buffer.WriteString(strconv.Itoa(int(b)) + ", ")
	}
	buffer.WriteString("}\n")

	// write to file
	err = os.WriteFile(os.Args[1]+".go", buffer.Bytes(), 0644)
	if err != nil {
		println("Error writing file:", err.Error())
		os.Exit(1)
	}

}
