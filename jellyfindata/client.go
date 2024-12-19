package jellyfindata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FetchUserActivity(baseURL, userID, apiKey string, extraHeaders map[string]string) ([]byte, error) {
	// Define query parameters
	params := url.Values{}
	params.Add("userId", userID)
	params.Add("isPlayed", "true")
	params.Add("recursive", "true")
	params.Add("IncludeItemTypes", "Movie,Episode")

	// Build the request URL
	endpoint := fmt.Sprintf("%s/Items", baseURL)
	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	// Create request headers
	headers := map[string]string{
		"X-Emby-Token": apiKey,
	}
	for key, value := range extraHeaders {
		headers[key] = value
	}

	// Make the HTTP GET request
	client := &http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ParseActivityResponse(data []byte) (*ActivityResponse, error) {
	var activity ActivityResponse
	err := json.Unmarshal(data, &activity)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}
