package cmd

import (
	"log"

	"github.com/leaanthony/mewn"
	"gopkg.in/yaml.v3"
)

type LinuxDB struct {
	Distributions map[string]*Distribution `yaml:"distributions"`
}

type Distribution struct {
	ID       string              `yaml:"id"`
	Releases map[string]*Release `yaml:"releases"`
}

func (d *Distribution) GetRelease(name string) *Release {
	result := d.Releases[name]
	if result == nil {
		result = d.Releases["default"]
	}
	return result
}

type Release struct {
	Name      string          `yaml:"name"`
	Version   string          `yaml:"version"`
	Programs  []*Prerequisite `yaml:"programs"`
	Libraries []*Prerequisite `yaml:"libraries"`
}

type Prerequisite struct {
	Name string `yaml:"name"`
	Help string `yaml:"help,omitempty"`
}

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

func (l *LinuxDB) ImportData(data []byte) error {
	return yaml.Unmarshal(data, l)
}

func (l *LinuxDB) GetDistro(name string) *Distribution {
	return l.Distributions[name]
}

// func (l *LinuxDB) GetDistribution(name string) *Distribution {
// 	return l.Distributions[name]
// }

// func (l *LinuxDB) Save() error {
// 	return fs.SaveAsJSON(l, l.filename)
// }

// func (l *LinuxDB) GetPreRequisistes(name string, version string) ([]*DistributionPrerequisite, error) {
// 	distro := l.GetDistribution(name)
// 	if distro == nil {
// 		return nil, fmt.Errorf("distribution %s unsupported at this time", name)
// 	}

// 	releaseInfo := distro.GetReleaseOrDefault(version)
// 	return releaseInfo.Prerequisites, nil

// }

// func (l *LinuxDBDistribution) GetReleaseOrDefault(release string) *DistributionRelease {
// 	result := l.Releases[release]
// 	if result == nil {
// 		result = l.Releases["default"]
// 	}
// 	return result
// }

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
