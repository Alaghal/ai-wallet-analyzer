package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OpenAIClient struct {
	apiURL string
	apiKey string
	model  string
	client *http.Client
}

func NewOpenAIClient(apiURL, apiKey, model string, timeout time.Duration) *OpenAIClient {
	return &OpenAIClient{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type openAIChatRequest struct {
	Model    string              `json:"model"`
	Messages []openAIChatMessage `json:"messages"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message openAIChatMessage `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) GenerateSummary(ctx context.Context, prompt string) (string, error) {
	payload := openAIChatRequest{
		Model: c.model,
		Messages: []openAIChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal openai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create openai request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai returned status %d", resp.StatusCode)
	}

	var parsed openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode openai response: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	summary := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if summary == "" {
		return "", fmt.Errorf("openai returned empty summary")
	}

	return summary, nil
}
