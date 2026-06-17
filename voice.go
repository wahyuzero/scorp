package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Voice / Audio Support — STT + TTS
// ──────────────────────────────────────────────

// Python binary paths (resolved dynamically via config_paths.go)
var (
	sttPythonBin = sttPythonPath()
	sttScript    = sttScriptPath()
	ttsBin       = ttsBinPath()
	uploadDir    = uploadsDir()
)

// Per-chat voice reply toggle
var (
	voiceReplyMu   sync.RWMutex
	voiceReplyChat = make(map[int64]bool) // chatID → voice replies enabled
)

// isVoiceReplyEnabled checks if a chat has voice replies enabled.
func isVoiceReplyEnabled(chatID int64) bool {
	voiceReplyMu.RLock()
	defer voiceReplyMu.RUnlock()
	return voiceReplyChat[chatID]
}

// toggleVoiceReply toggles voice replies for a chat.
func toggleVoiceReply(chatID int64) bool {
	voiceReplyMu.Lock()
	defer voiceReplyMu.Unlock()
	voiceReplyChat[chatID] = !voiceReplyChat[chatID]
	return voiceReplyChat[chatID]
}

// transcribeAudio converts an audio file to text using faster-whisper.
// audioPath is expected to be a .ogg file from Telegram.
func transcribeAudio(audioPath, language string) (string, error) {
	// Convert .ogg to .wav (16kHz mono, required by whisper)
	wavPath := strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + ".wav"

	// ffmpeg: convert to 16kHz mono WAV
	cmd := exec.Command("ffmpeg", "-y", "-i", audioPath,
		"-ar", "16000", "-ac", "1", "-f", "wav", wavPath)
	cmd.Stderr = nil
	cmd.Stdout = nil
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg convert: %w", err)
	}
	defer os.Remove(wavPath)

	// Run STT
	args := []string{sttScript, wavPath}
	if language != "" {
		args = append(args, language)
	}

	sttCmd := exec.Command(sttPythonBin, args...)
	output, err := sttCmd.Output()
	if err != nil {
		// Get stderr for debugging
		stderr := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = string(exitErr.Stderr)
		}
		return "", fmt.Errorf("STT error: %w %s", err, stderr)
	}

	text := strings.TrimSpace(string(output))
	return text, nil
}

// synthesizeSpeech converts text to speech using edge-tts.
// Returns the path to the generated .ogg file.
func synthesizeSpeech(text, voice string) (string, error) {
	if voice == "" {
		voice = "id-ID-ArdiNeural" // Default: Indonesian male voice
	}

	// edge-tts outputs MP3, then we convert to OGG with opus codec for Telegram
	mp3Path := filepath.Join(uploadDir, fmt.Sprintf("tts_%d.mp3", time.Now().UnixNano()))
	oggPath := strings.TrimSuffix(mp3Path, ".mp3") + ".ogg"

	// Run edge-tts
	cmd := exec.Command(ttsBin, "--voice", voice, "--text", text, "--write-media", mp3Path)
	cmd.Stderr = nil
	cmd.Stdout = nil
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("edge-tts: %w", err)
	}
	defer os.Remove(mp3Path)

	// Convert MP3 to OGG/Opus (Telegram voice message format)
	convCmd := exec.Command("ffmpeg", "-y", "-i", mp3Path,
		"-c:a", "libopus", "-b:a", "32k", "-vbr", "on",
		"-frame_duration", "60", oggPath)
	convCmd.Stderr = nil
	convCmd.Stdout = nil
	if err := convCmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg ogg: %w", err)
	}

	return oggPath, nil
}

// downloadTelegramFile is already defined in telegram.go
// We use the existing: downloadTelegramFile(fileID, destPath string) bool

// sendChatAction sends a chat action indicator (typing, upload_audio, etc.)
func sendChatAction(chatID int64, action string) {
	url := fmt.Sprintf("%s/sendChatAction?chat_id=%d&action=%s", tgBase, chatID, action)
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}
	resp.Body.Close()
}

// uploadFile uploads a file to Telegram API using multipart/form-data.
// fieldName is the API field (e.g., "voice", "audio", "document").
func uploadFile(apiURL, fieldName, filename string, data []byte, extra map[string]string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create the file part
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	part.Write(data)

	// Add extra fields
	for key, val := range extra {
		writer.WriteField(key, val)
	}
	writer.Close()

	resp, err := httpClient.Post(apiURL, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// sendVoiceMessage sends a .ogg file as a Telegram voice message.
func sendVoiceMessage(chatID int64, oggPath string) error {
	data, err := os.ReadFile(oggPath)
	if err != nil {
		return fmt.Errorf("read ogg: %w", err)
	}
	apiURL := fmt.Sprintf("%s/sendVoice?chat_id=%d", tgBase, chatID)
	return uploadFile(apiURL, "voice", oggPath, data, nil)
}

// sendVoiceReply synthesizes text to speech and sends as voice message.
// Cleans up the temporary file after sending.
func sendVoiceReply(chatID int64, text string) {
	oggPath, err := synthesizeSpeech(text, "")
	if err != nil {
		log.Printf("[voice] TTS error: %v", err)
		return
	}
	defer os.Remove(oggPath)

	if err := sendVoiceMessage(chatID, oggPath); err != nil {
		log.Printf("[voice] Send error: %v", err)
	}
}

// handleVoiceMessage processes an incoming voice message:
// download → STT → process as text command/agent query.
func handleVoiceMessage(doc TGDocument) {
	// Download voice file using existing telegram.go function
	os.MkdirAll(uploadDir, 0755)
	localPath := filepath.Join(uploadDir, fmt.Sprintf("voice_%d.ogg", doc.MsgID))
	if !downloadTelegramFile(doc.FileID, localPath) {
		log.Printf("[voice] Download failed for file_id=%s", doc.FileID)
		sendMessage("❌ Voice download failed", mainMenuKeyboard())
		return
	}
	defer os.Remove(localPath)

	// Show typing indicator
	sendChatAction(doc.ChatID, "typing")

	// Transcribe
	text, err := transcribeAudio(localPath, "") // auto-detect language
	if err != nil {
		log.Printf("[voice] STT failed: %v", err)
		sendMessage("❌ Transcription failed: "+err.Error(), mainMenuKeyboard())
		return
	}

	if text == "" {
		sendMessage("🔇 Empty transcription (no speech detected)", mainMenuKeyboard())
		return
	}

	log.Printf("[voice] Transcribed: %s", truncateStr(text, 80))

	// Process transcribed text as if it was typed
	cid := fmt.Sprintf("%d", doc.ChatID)
	if isAgentMode(cid) {
		// Show what was heard before processing
		sendMessage(fmt.Sprintf("🎤 <i>%s</i>", text), nil)
		handleAction("/agent "+text, doc.ChatID, 0, "")
	} else {
		// Non-agent mode: treat as command
		sendMessage(fmt.Sprintf("🎤 <i>%s</i>", text), nil)
		handleAction(text, doc.ChatID, 0, "")
	}

	// If voice replies are enabled, also send a voice response
	// (handled in the response path via isVoiceReplyEnabled check)
}
