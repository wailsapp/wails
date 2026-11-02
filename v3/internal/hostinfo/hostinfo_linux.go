package hostinfo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/internal/pcinames"
)

func Processors() ([]Processor, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readProcessors(f)
}

type processorSocket struct {
	processor Processor
	cores     map[string]bool
	threads   int
}

func (s *processorSocket) add(fields map[string]string) {
	if model := fields["model name"]; model != "" {
		s.processor.Model = model
	} else if model := fields["Processor"]; model != "" {
		s.processor.Model = model
	}
	if cores, err := strconv.Atoi(fields["cpu cores"]); err == nil && cores > s.processor.Cores {
		s.processor.Cores = cores
	}
	if core, ok := fields["core id"]; ok {
		s.cores[core] = true
	}
	s.threads++
}

func (s *processorSocket) info() Processor {
	p := s.processor
	if p.Cores == 0 {
		p.Cores = len(s.cores)
		if p.Cores == 0 {
			p.Cores = s.threads
		}
	}
	if p.Model == "" {
		p.Model = "Unknown"
	}
	return p
}

func readProcessors(r io.Reader) ([]Processor, error) {
	var sockets []*processorSocket
	byID := make(map[string]*processorSocket)
	fields := make(map[string]string)
	flush := func() {
		if _, ok := fields["processor"]; !ok {
			clear(fields)
			return
		}
		id := fields["physical id"]
		s := byID[id]
		if s == nil {
			s = &processorSocket{cores: make(map[string]bool)}
			byID[id] = s
			sockets = append(sockets, s)
		}
		s.add(fields)
		clear(fields)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if ok {
			fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	flush()
	result := make([]Processor, 0, len(sockets))
	for _, s := range sockets {
		result = append(result, s.info())
	}
	if len(result) == 0 {
		return nil, errors.New("no processors in cpuinfo")
	}
	return result, nil
}

func Memory() (uint64, error) {
	return physicalMemory(os.DirFS("/"))
}

func physicalMemory(root fs.FS) (uint64, error) {
	if size, err := memoryBlocks(root); err == nil {
		return size, nil
	}
	// Containers and kernels without memory hotplug may expose only procfs.
	data, err := fs.ReadFile(root, "proc/meminfo")
	if err != nil {
		return 0, err
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 3 && fields[0] == "MemTotal:" && fields[2] == "kB" {
			n, err := strconv.ParseUint(fields[1], 10, 64)
			if err == nil && n > 0 && n <= ^uint64(0)/1024 {
				return n * 1024, nil
			}
		}
	}
	return 0, errors.New("no physical memory information")
}

func memoryBlocks(root fs.FS) (uint64, error) {
	const directory = "sys/devices/system/memory"
	data, err := fs.ReadFile(root, directory+"/block_size_bytes")
	if err != nil {
		return 0, err
	}
	block, err := strconv.ParseUint(strings.TrimSpace(string(data)), 16, 64)
	if err != nil || block == 0 {
		return 0, errors.New("invalid memory block size")
	}
	entries, err := fs.ReadDir(root, directory)
	if err != nil {
		return 0, err
	}
	var count uint64
	for _, entry := range entries {
		suffix, ok := strings.CutPrefix(entry.Name(), "memory")
		if !ok {
			continue
		}
		if _, err := strconv.ParseUint(suffix, 10, 64); err != nil {
			continue
		}
		state, err := fs.ReadFile(root, directory+"/"+entry.Name()+"/state")
		if err == nil && strings.TrimSpace(string(state)) == "online" {
			count++
		}
	}
	if count == 0 || count > ^uint64(0)/block {
		return 0, errors.New("no online memory blocks or size overflow")
	}
	return count * block, nil
}

func Graphics() ([]GraphicsCard, error) {
	return readGraphics("/sys/bus/pci/devices", pciNames)
}

type nameLookup func(vendorID, deviceID string) (vendor, product string)

func pciNames() (nameLookup, error) {
	db, err := pcinames.Load()
	if err != nil {
		return nil, err
	}
	return db.Lookup, nil
}

func readGraphics(directory string, names func() (nameLookup, error)) ([]GraphicsCard, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var result []GraphicsCard
	var lookup nameLookup
	for _, entry := range entries {
		device := filepath.Join(directory, entry.Name())
		class, err := readHex(filepath.Join(device, "class"))
		if err != nil || class>>16 != 3 {
			continue
		}

		if lookup == nil {
			lookup, err = names()
			if err != nil {
				return nil, err
			}
		}
		card, err := pciCard(device, lookup)
		if err != nil {
			continue
		}
		result = append(result, card)
	}
	return result, nil
}

func pciCard(device string, lookup nameLookup) (GraphicsCard, error) {
	vendor, err := readHex(filepath.Join(device, "vendor"))
	if err != nil {
		return GraphicsCard{}, err
	}
	product, err := readHex(filepath.Join(device, "device"))
	if err != nil {
		return GraphicsCard{}, err
	}
	vendorName, productName := lookup(fmt.Sprintf("%04x", vendor), fmt.Sprintf("%04x", product))
	if productName == "" {
		productName = fmt.Sprintf("PCI %04x:%04x", vendor, product)
	}
	card := GraphicsCard{Name: productName, Vendor: vendorName}
	if driver, err := os.Readlink(filepath.Join(device, "driver")); err == nil {
		card.Driver = filepath.Base(driver)
	}
	return card, nil
}

func readHex(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimPrefix(strings.TrimSpace(string(data)), "0x"), 16, 32)
}
