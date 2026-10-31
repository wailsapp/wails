package setupwizard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type WailsConfigInfo struct {
	CompanyName       string `json:"companyName" yaml:"companyName"`
	ProductName       string `json:"productName" yaml:"productName"`
	ProductIdentifier string `json:"productIdentifier" yaml:"productIdentifier"`
	Description       string `json:"description" yaml:"description"`
	Copyright         string `json:"copyright" yaml:"copyright"`
	Comments          string `json:"comments,omitempty" yaml:"comments,omitempty"`
	Version           string `json:"version" yaml:"version"`
}

// WailsConfig represents the wails.yaml configuration
type WailsConfig struct {
	Info WailsConfigInfo `json:"info" yaml:"info"`
}

func (w *Wizard) handleLegacyWailsConfig(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// Find wails.yaml in current directory or parent directories
	configPath := findLegacyWailsConfig()

	switch r.Method {
	case http.MethodGet:
		if configPath == "" {
			json.NewEncoder(rw).Encode(nil)
			return
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			json.NewEncoder(rw).Encode(nil)
			return
		}

		var config WailsConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			json.NewEncoder(rw).Encode(nil)
			return
		}

		json.NewEncoder(rw).Encode(config)

	case http.MethodPost:
		var config WailsConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if configPath == "" {
			configPath = "wails.yaml"
		}

		data, err := mergeLegacyProjectInfo(configPath, config.Info)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(map[string]string{"status": "saved", "path": configPath})

	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func findLegacyWailsConfig() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		configPath := filepath.Join(dir, "wails.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// Update only the wizard-owned metadata, retaining custom YAML and comments.
func mergeLegacyProjectInfo(path string, info WailsConfigInfo) ([]byte, error) {
	var document yaml.Node
	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, &document); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	var prefix []byte
	if len(document.Content) == 0 {
		// yaml.v3 has no node to attach comments to in comment-only files.
		prefix = data
		if len(prefix) > 0 && prefix[len(prefix)-1] != '\n' {
			prefix = append(prefix, '\n')
		}
		document = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}}
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("legacy configuration must be a YAML mapping")
	}
	root := document.Content[0]
	var existing *yaml.Node
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value == "info" {
			existing = root.Content[i+1]
			break
		}
	}
	if existing == nil {
		existing = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		root.Content = append(root.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "info"}, existing)
	}
	if existing.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("legacy info must be a YAML mapping")
	}
	var edited yaml.Node
	if err := edited.Encode(info); err != nil {
		return nil, err
	}
	// Empty comments are still an explicit edit, even though the legacy struct
	// omits the field when marshalled.
	if info.Comments == "" {
		edited.Content = append(edited.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "comments"}, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: ""})
	}
	for i := 0; i < len(edited.Content); i += 2 {
		found := false
		for j := 0; j < len(existing.Content); j += 2 {
			if existing.Content[j].Value == edited.Content[i].Value {
				previous := existing.Content[j+1]
				next := *edited.Content[i+1]
				next.HeadComment, next.LineComment, next.FootComment = previous.HeadComment, previous.LineComment, previous.FootComment
				existing.Content[j+1] = &next
				found = true
				break
			}
		}
		if !found {
			existing.Content = append(existing.Content, edited.Content[i], edited.Content[i+1])
		}
	}
	encoded, err := yaml.Marshal(&document)
	if err != nil {
		return nil, err
	}
	return append(prefix, encoded...), nil
}
