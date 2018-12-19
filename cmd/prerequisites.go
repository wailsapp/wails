package cmd

import (
	"fmt"
	"runtime"
)

// binaryPrerequisite defines a binaryPrerequisite
type binaryPrerequisite struct {
	Name string
	Help string
	Path string
}

func newBinaryPrerequisite(name, help string) *binaryPrerequisite {
	return &binaryPrerequisite{Name: name, Help: help}
}

// binaryPrerequisites is a list of binaryPrerequisites
type binaryPrerequisites []*binaryPrerequisite

// Add given prereq object to list
func (p *binaryPrerequisites) Add(prereq *binaryPrerequisite) {
	*p = append(*p, prereq)
}

func (p *binaryPrerequisites) check() (success *binaryPrerequisites, failed *binaryPrerequisites) {
	success = &binaryPrerequisites{}
	failed = &binaryPrerequisites{}
	programHelper := NewProgramHelper()
	for _, prereq := range *p {
		bin := programHelper.FindProgram(prereq.Name)
		if bin == nil {
			failed.Add(prereq)
		} else {
			path, err := bin.GetFullPathToBinary()
			if err != nil {
				failed.Add(prereq)
			} else {
				prereq.Path = path
				success.Add(prereq)
			}
		}
	}

	return success, failed
}

var platformbinaryPrerequisites = make(map[string]*binaryPrerequisites)

func init() {
	platformbinaryPrerequisites["darwin"] = &binaryPrerequisites{}
	newDarwinbinaryPrerequisite("clang", "Please install with `xcode-select --install` and try again")
}

func newDarwinbinaryPrerequisite(name, help string) {
	prereq := newBinaryPrerequisite(name, help)
	platformbinaryPrerequisites["darwin"].Add(prereq)
}

func CheckBinaryPrerequisites() (*binaryPrerequisites, *binaryPrerequisites, error) {
	platformPreReqs := platformbinaryPrerequisites[runtime.GOOS]
	if platformPreReqs == nil {
		return nil, nil, fmt.Errorf("Platform '%s' is not supported at this time", runtime.GOOS)
	}
	success, failed := platformPreReqs.check()
	return success, failed, nil
}

func CheckNonBinaryPrerequisites() error {

	var err error

	// Check non-binaries
	if runtime.GOOS == "linux" {

	}
	return err
}
