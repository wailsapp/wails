package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	BaseURL = "https://api.crowdin.com/api/v2"
)

func doRequest(client *http.Client, method, url, token string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Request Error: %s", data["error"].(map[string]interface{})["message"])

	}

	return body, nil
}

func getBranchID(client *http.Client, projectID, branchName, token string) (int, error) {
	body, err := doRequest(client, "GET", fmt.Sprintf("%s/projects/%s/branches", BaseURL, projectID), token)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	for _, branch := range data["data"].([]interface{}) {
		if branch.(map[string]interface{})["data"].(map[string]interface{})["name"] == branchName {
			return int(branch.(map[string]interface{})["data"].(map[string]interface{})["id"].(float64)), nil
		}
	}

	return 0, fmt.Errorf("branch not found")
}

func getLanguageProgress(client *http.Client, projectID string, branchID int, token string) (map[string]float64, error) {
	body, err := doRequest(client, "GET", fmt.Sprintf("%s/projects/%s/branches/%d/languages/progress", BaseURL, projectID, branchID), token)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	progress := make(map[string]float64)

	for _, languageProgress := range data["data"].([]interface{}) {
		languageID := languageProgress.(map[string]interface{})["data"].(map[string]interface{})["languageId"].(string)
		approvalProgress := languageProgress.(map[string]interface{})["data"].(map[string]interface{})["approvalProgress"].(float64)
		progress[languageID] = approvalProgress
	}

	return progress, nil
}

func main() {
	projectID := "531392"
	branchName := "v2"
	token := os.Getenv("CROWDIN_PERSONAL_TOKEN")

	client := &http.Client{}
	branchID, err := getBranchID(client, projectID, branchName, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	progress, err := getLanguageProgress(client, projectID, branchID, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	progressYaml, err := yaml.Marshal(progress)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(filepath.Join("../../data/en/translation-progress.yaml"), progressYaml, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
