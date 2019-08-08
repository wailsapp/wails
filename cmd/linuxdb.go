package cmd

import (
	"log"

	"github.com/leaanthony/mewn"
	"gopkg.in/yaml.v3"
)

// LinuxDB is the database for linux distribution data.
type LinuxDB struct {
	Distributions map[string]*Distribution `yaml:"distributions"`
}

// Distribution holds the os-release ID and a map of releases.
type Distribution struct {
	ID       string              `yaml:"id"`
	Releases map[string]*Release `yaml:"releases"`
}

// GetRelease attempts to return the specific Release information
// for the given release name. If there is no specific match, the
// default release data is returned.
func (d *Distribution) GetRelease(name string) *Release {
	result := d.Releases[name]
	if result == nil {
		result = d.Releases["default"]
	}
	return result
}

// Release holds the name and version of the release as given by
// os-release. Programs is a slice of dependant programs required
// to be present on the local installation for Wails to function.
// Libraries is a slice of libraries that must be present for Wails
// applications to compile.
type Release struct {
	Name              string          `yaml:"name"`
	Version           string          `yaml:"version"`
	GccVersionCommand string          `yaml:"gccversioncommand"`
	Programs          []*Prerequisite `yaml:"programs"`
	Libraries         []*Prerequisite `yaml:"libraries"`
}

// Prerequisite is a simple struct containing a program/library name
// plus the distribution specific help text indicating how to install
// it.
type Prerequisite struct {
	Name string `yaml:"name"`
	Help string `yaml:"help,omitempty"`
}

// Load will load the given filename from disk and attempt to
// import the data into the LinuxDB.
func (l *LinuxDB) Load(filename string) error {
	if fs.FileExists(filename) {
		data, err := fs.LoadAsBytes(filename)
		if err != nil {
			return err
		}
		return l.ImportData(data)
	}
	return nil
}

// ImportData will unmarshal the given YAML formatted data
// into the LinuxDB
func (l *LinuxDB) ImportData(data []byte) error {
	return yaml.Unmarshal(data, l)
}

// GetDistro returns the Distribution information for the
// given distribution name. If the distribution is not supported,
// nil is returned.
func (l *LinuxDB) GetDistro(name string) *Distribution {
	return l.Distributions[name]
}

// NewLinuxDB creates a new LinuxDB instance from the bundled
// linuxdb.yaml file.
func NewLinuxDB() *LinuxDB {
	data := mewn.Bytes("./linuxdb.yaml")
	result := LinuxDB{
		Distributions: make(map[string]*Distribution),
	}
	err := result.ImportData(data)
	if err != nil {
		log.Fatal(err)
	}
	return &result
}
