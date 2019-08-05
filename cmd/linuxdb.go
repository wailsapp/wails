package cmd

import (
	"log"
	"path"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type LinuxDB struct {
	Distributions map[string]*Distribution `yaml:"distributions"`
	filename      string
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
		return yaml.Unmarshal([]byte(data), l)
	}
	return nil
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
	_, filename, _, _ := runtime.Caller(1)
	fullPath, err := filepath.Abs(filepath.Join(path.Dir(filename), "linuxdb.yaml"))
	if err != nil {
		log.Fatal("Unable to open linuxdb: " + fullPath)
	}
	result := LinuxDB{
		filename:      fullPath,
		Distributions: make(map[string]*Distribution),
	}
	err = result.Load(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	return &result
}
