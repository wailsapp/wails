package cmd

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
)

// FrameworkMetadata contains information about a given framework
type FrameworkMetadata struct {
	Name        string `json:"name"`
	BuildTag    string `json:"buildtag"`
	Description string `json:"description"`
}

// Utility function for creating new FrameworkMetadata structs
func loadFrameworkMetadata(pathToMetadataJSON string) (*FrameworkMetadata, error) {
	result := &FrameworkMetadata{}
	configData, err := ioutil.ReadFile(pathToMetadataJSON)
	if err != nil {
		return nil, err
	}
	// Load and unmarshall!
	err = json.Unmarshal(configData, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetFrameworks returns information about all the available frameworks
func GetFrameworks() ([]*FrameworkMetadata, error) {

	var err error

	// Calculate framework base dir
	_, filename, _, _ := runtime.Caller(1)
	frameworksBaseDir := filepath.Join(path.Dir(filename), "frameworks")

	// Get the subdirectories
	fs := NewFSHelper()
	frameworkDirs, err := fs.GetSubdirs(frameworksBaseDir)
	if err != nil {
		return nil, err
	}

	// Prepare result
	result := []*FrameworkMetadata{}

	// Iterate framework directories, looking for metadata.json files
	for _, frameworkDir := range frameworkDirs {
		var frameworkMetadata FrameworkMetadata
		metadataFile := filepath.Join(frameworkDir, "metadata.json")
		jsonData, err := ioutil.ReadFile(metadataFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonData, &frameworkMetadata)
		if err != nil {
			return nil, err
		}
		result = append(result, &frameworkMetadata)
	}

	// Read in framework metadata
	return result, nil
}
