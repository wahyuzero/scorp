package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"scorp-agent/internal/helpers"
	"scorp-agent/models"
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

	return helpers.TruncOutput(result, helpers.MaxToolOutput), true
}

// callVisionModel sends an image + question to the configured vision model
func callVisionModel(question, dataURL string) (string, error) {
	// Build multimodal content parts
	parts := []map[string]interface{}{
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
	}
	partsJSON, _ := json.Marshal(parts)

	chatMsgs := []models.ChatMessage{
		{
			Role:    "user",
			Content: string(partsJSON),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	reply, _, err := models.CallModelWithFallback(ctx, "vision", chatMsgs)
	if err != nil {
		return "", err
	}
	return reply, nil
}
