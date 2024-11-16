package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var apiKey string

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system"`
}

type Response struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence any    `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func sendRequest(apiKey string, userPrompt string) (*Response, error) {
	time.Sleep(2 * time.Second)
	url := "https://api.anthropic.com/v1/messages"

	requestBody := RequestBody{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 4096,
		System:    `You are a highly skilled translator with expertise in many languages. Your task is to identify the language of the text I provide and accurately translate it into the specified target language, indicated by a language code e.g. 'zh' for chinese, while preserving the meaning, tone, and nuance of the original text. Please maintain proper grammar, spelling, and punctuation in the translated version. It will contain special Markdown syntax which must be preserved. You must also ensure that the translated text is culturally appropriate and sensitive to the target audience. You are responsible for the quality of the translation. If you are unable to provide a high-quality translation, please let me know.`,
		Messages: []Message{
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling request body")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error creating request")
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body")
	}

	//unmarshal the response
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response")
	}

	return &response, nil
}

func translateDirectory(rootpath string, langcode string, langname string) error {
	// Walk through the files in the directory
	// for each file, read the file, print out the full path and filesize
	// if the file is a directory, call this function recursively
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == rootpath {
			return nil
		}
		// Skip if not `.md` file or _category_.json
		if filepath.Ext(path) != ".md" || strings.Contains(path, "_category_.json") {
			return nil
		}
		if info.IsDir() {
			return translateDirectory(path, langcode, langname)
		}
		// Calculate the target filename. This is the same path, but with using "../../docs/langcode" instead of "../../docs/en"
		targetPath := strings.Replace(path, "/en/", "/"+langcode+"/", 1)
		// If target file exists, skip it.
		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("Skipping: %s\n", path)
			return nil
		}
		fmt.Printf("Translating: %s, Size: %d\n", path, info.Size())
		// Load the file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		user := fmt.Sprintf("Translate the following text from English to %s (%s):\n\n%s", langname, langcode, string(data))
		// Send the request
		response, err := sendRequest(apiKey, user)
		if err != nil {
			return err
		}
		// Create the target directory
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}
		// Write the content
		if len(response.Content) == 0 {
			fmt.Printf("No content returned for: %v\n", response)
			return nil
		}
		err = os.WriteFile(targetPath, []byte(response.Content[0].Text), 0644)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func main() {

	// Read the first argument. It is a language code. If it is not provided, show an error message
	if len(os.Args) < 2 {
		fmt.Println("Please provide the language as the first argument")
		return
	}
	lang := os.Args[1]

	if len(os.Args) < 3 {
		fmt.Println("Please provide the language code as the second argument")
		return
	}
	langcode := os.Args[2]

	// load the API key from the .env file
	data, err := os.ReadFile(".env")
	if err != nil {
		fmt.Println("Error reading .env file:", err)
		return
	}
	apiKey = strings.TrimSpace(string(data))

	rootPath := "../../docs/en"
	err = translateDirectory(rootPath, langcode, lang)
	if err != nil {
		fmt.Println("Error translating directory:", err)

	}

}
