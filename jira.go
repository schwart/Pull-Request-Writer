package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetJiraLink(
	domain string,
	issueKey string,
) string {
	return fmt.Sprintf("https://%s/browse/%s", domain, issueKey)
}

func GetJiraDescription(
	email string,
	apiToken string,
	domain string,
	issueKey string,
) (string, error) {
	url := fmt.Sprintf("https://%s/rest/api/3/issue/%s", domain, issueKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	// Add query parameters: fields=description and expand=renderedFields
	q := req.URL.Query()
	q.Add("fields", "description")
	q.Add("expand", "renderedFields")
	req.URL.RawQuery = q.Encode()

	// Set up basic authentication
	req.SetBasicAuth(email, apiToken)

	// Create an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	// defer close, then log fatal if error
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Error fetching issue: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the JSON response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	var issueData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &issueData); err != nil {
		return "", fmt.Errorf("Error parsing JSON: %v", err)
	}

	// Extract the renderedFields.description
	renderedFields, ok := issueData["renderedFields"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("No 'renderedFields' object found in response.")
	}

	description, ok := renderedFields["description"].(string)
	if !ok {
		return "", fmt.Errorf("No 'description' found in renderedFields.")
	}

	return description, nil
}
