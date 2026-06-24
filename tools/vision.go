package tools

import (
	"scorp-agent/models"
	"scorp-agent/internal/helpers"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// ──────────────────────────────────────────────
// Vision / Image Analysis Tool
// ──────────────────────────────────────────────

// executeAnalyzeImage reads an image file and sends it to a vision-capable model
// for analysis. Returns the model's text description of the image.
func ExecuteAnalyzeImage(args map[string]interface{}) (string, bool) {
	path := helpers.GetStringArg(args, "path", "")
	if path == "" {
		return "Error: path is required", false
	}

	question := helpers.GetStringArg(args, "question", "Describe this image in detail.")

	// Read image file
	imgData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading image: %v", err), false
	}

	// Check file size (limit to 10MB to avoid huge API calls)
	if len(imgData) > 10*1024*1024 {
		return fmt.Sprintf("Image too large: %d MB (max 10 MB)", len(imgData)/1024/1024), false
	}

	// Determine MIME type from extension
	mime := "image/png"
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		mime = "image/jpeg"
	} else if strings.HasSuffix(lower, ".webp") {
		mime = "image/webp"
	} else if strings.HasSuffix(lower, ".gif") {
		mime = "image/gif"
	}

	// Encode as base64
	b64Data := base64.StdEncoding.EncodeToString(imgData)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mime, b64Data)

	// Call vision model via 9router
	result, err := callVisionModel(question, dataURL)
	if err != nil {
		log.Printf("[vision] Error: %v", err)
		return fmt.Sprintf("Vision analysis error: %v", err), false
	}

	return helpers.TruncOutput(result, helpers.MaxToolOutput), false
}

// callVisionModel sends an image + question to a vision-capable model via OpenAI-compatible API
func callVisionModel(question, dataURL string) (string, error) {
	// Get the configured model or use default vision model
	modelCfg := models.RouteModel("chat")
	if modelCfg == nil {
		return "", fmt.Errorf("no model configured")
	}

	// Use a vision-capable model
	// kr/claude-sonnet-4 is the only model on 9router that reliably supports image input
	apiKey := models.ResolveAPIKey(modelCfg)
	if apiKey == "" {
		return "", fmt.Errorf("no API key for vision model — %s", models.KeySourceLabel(modelCfg))
	}
	baseURL := strings.TrimRight(modelCfg.BaseURL, "/")
	endpoint := baseURL + "/chat/completions"

	// Build request with image content
	reqBody := map[string]interface{}{
		"model": "kr/claude-sonnet-4", // Vision-capable model on 9router
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": question,
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": dataURL,
						},
					},
				},
			},
		},
		"max_tokens": 1000,
		"stream":     false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := models.GetAIClient(baseURL).Do(req)
	if err != nil {
		return "", fmt.Errorf("vision API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("vision API error (HTTP %d): %s", resp.StatusCode, helpers.TruncateStr(string(body), 300))
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("parse error: %s", helpers.TruncateStr(string(body), 200))
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from vision model")
	}

	return chatResp.Choices[0].Message.Content, nil
}
