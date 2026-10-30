// Package pcinames reads vendor and device names from the host's pci.ids data.
package pcinames

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
)

type Database struct {
	vendors map[uint16]*vendor
}

type vendor struct {
	name    string
	devices map[uint16]string
}

func (d *Database) Lookup(vendorID, deviceID string) (string, string) {
	v, err := strconv.ParseUint(vendorID, 16, 16)
	if err != nil {
		return "", ""
	}
	entry := d.vendors[uint16(v)]
	if entry == nil {
		return "", ""
	}
	p, err := strconv.ParseUint(deviceID, 16, 16)
	if err != nil {
		return entry.name, ""
	}
	return entry.name, entry.devices[uint16(p)]
}

func read(r io.Reader) (*Database, error) {
	buffer := bufio.NewReader(r)
	if magic, _ := buffer.Peek(2); bytes.Equal(magic, []byte{0x1f, 0x8b}) {
		compressed, err := gzip.NewReader(buffer)
		if err != nil {
			return nil, err
		}
		defer compressed.Close()
		return boundedParse(compressed)
	}
	return boundedParse(buffer)
}

func boundedParse(r io.Reader) (*Database, error) {
	limited := &io.LimitedReader{R: r, N: maxDownloadBytes + 1}
	db, err := parse(limited)
	if limited.N == 0 {
		return nil, fmt.Errorf("PCI names database exceeds size limit")
	}
	return db, err
}

func parse(r io.Reader) (*Database, error) {
	db := &Database{vendors: make(map[uint16]*vendor)}
	var current *vendor
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if line[0] != '\t' {
			current = nil
			if id, name, ok := record(line); ok {
				current = &vendor{name: name, devices: make(map[uint16]string)}
				db.vendors[id] = current
			}
			continue
		}
		// Subsystems, class hierarchies and future record types are not names
		// of the PCI vendor/device pair displayed by doctor.
		if current != nil && len(line) > 1 && line[1] != '\t' {
			if id, name, ok := record(line[1:]); ok {
				current.devices[id] = name
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read PCI names: %w", err)
	}
	return db, nil
}

func record(line []byte) (uint16, string, bool) {
	if len(line) < 6 || (line[4] != ' ' && line[4] != '\t') {
		return 0, "", false
	}
	id, err := strconv.ParseUint(string(line[:4]), 16, 16)
	if err != nil {
		return 0, "", false
	}
	name := bytes.TrimSpace(line[5:])
	if len(name) == 0 {
		return 0, "", false
	}
	return uint16(id), string(name), true
}
