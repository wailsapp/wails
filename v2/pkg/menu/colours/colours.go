package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Rgb struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

type Hsl struct {
	H float64 `json:"h"`
	S float64 `json:"s"`
	L float64 `json:"l"`
}

type InputCol struct {
	Colorid   uint8  `json:"colorId"`
	Hexstring string `json:"hexString"`
	Rgb       Rgb    `json:"rgb"`
	Hsl       Hsl    `json:"hsl"`
	Name      string `json:"name"`
}

//go:embed gen.tmpl
var Template string

func main() {
	var Cols []InputCol

	resp, err := http.Get("https://jonasjacek.github.io/colors/data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &Cols)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("cols").Parse(Template)
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, Cols)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filepath.Join("..", "cols.go"), buffer.Bytes(), 0o755)
	if err != nil {
		log.Fatal(err)
	}
}
