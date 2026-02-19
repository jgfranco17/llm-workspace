package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func FetchURL(ctx context.Context, params map[string]string) (string, error) {
	url, err := getRequiredParam(ctx, params, "url")
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(url, "https://") {
		return "", fmt.Errorf("only HTTPS URLs are allowed")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("failed to fetch URL: status code %d", response.StatusCode)
	}

	reader := io.LimitReader(response.Body, MaxOutputSize+1)
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return limitOutput(string(body), MaxOutputSize), nil
}
