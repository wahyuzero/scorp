package telegram

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"scorp-agent/config"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	maxTGFile        = 49 * 1024 * 1024  // 49MB
	partSize         = 20 * 1024 * 1024  // 20MB
	confirmThreshold = 100 * 1024 * 1024 // 100MB
)

var browseRoots = []struct {
	Path  string
	Label string
}{
	{config.HomeDir(), "🏠 Home"},
	{config.HomeDir() + "/projects", "📁 Projects"},
	{"/data/coolify", "☁️ Coolify Data"},
	{"/tmp", "📋 Temp"},
}

// Path ID mapping with bounded capacity (LRU prune after 5000 paths)
const maxTrackedPaths = 5000

var (
	pathMap     = make(map[string]string) // pid -> path
	reversePath = make(map[string]string) // path -> pid
	pathCounter int
	pathMu      sync.Mutex
)

func PathID(path string) string {
	pathMu.Lock()
	defer pathMu.Unlock()
	if pid, ok := reversePath[path]; ok {
		return pid
	}
	// Prune if capacity exceeded to prevent slow 24/7 memory leaks
	if len(pathMap) >= maxTrackedPaths {
		pathMap = make(map[string]string)
		reversePath = make(map[string]string)
		pathCounter = 0
	}
	pathCounter++
	pid := fmt.Sprintf("p%d", pathCounter)
	pathMap[pid] = path
	reversePath[path] = pid
	return pid
}

func GetPath(pid string) string {
	pathMu.Lock()
	defer pathMu.Unlock()
	return pathMap[pid]
}

func HumanSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2fGB", float64(size)/(1024*1024*1024))
}

// ──────────────────────────────────────────────
// Keyboards
// ──────────────────────────────────────────────

func RootsKeyboard() map[string]interface{} {
	var buttons []interface{}
	for _, r := range browseRoots {
		pid := PathID(r.Path)
		buttons = append(buttons, []interface{}{
			map[string]string{"text": r.Label, "callback_data": "fb:" + pid},
		})
	}
	buttons = append(buttons, []interface{}{
		map[string]string{"text": "📤 Upload File", "callback_data": "upload"},
	})
	buttons = append(buttons, []interface{}{
		map[string]string{"text": "◀️ Back to Menu", "callback_data": "menu"},
	})
	return map[string]interface{}{"inline_keyboard": buttons}
}

func DirKeyboard(path string) (string, map[string]interface{}) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "🚫 Permission denied", BackKB("files")
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	displayPath := strings.Replace(path, config.HomeDir(), "~", 1)
	lines := []string{fmt.Sprintf("📂 <b>%s</b>", displayPath)}

	var dirs, files []os.DirEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}

	if len(dirs) == 0 && len(files) == 0 {
		lines = append(lines, "  <i>(empty)</i>")
	}

	var buttons []interface{}

	// Directories with zip button
	for i, d := range dirs {
		if i >= 12 {
			lines = append(lines, fmt.Sprintf("  ... +%d more folders", len(dirs)-12))
			break
		}
		fullPath := filepath.Join(path, d.Name())
		pid := PathID(fullPath)
		buttons = append(buttons, []interface{}{
			map[string]string{"text": "📁 " + d.Name() + "/", "callback_data": "fb:" + pid},
			map[string]string{"text": "📦", "callback_data": "zp:" + pid},
		})
	}

	// Files
	for i, f := range files {
		if i >= 12 {
			lines = append(lines, fmt.Sprintf("  ... +%d more files", len(files)-12))
			break
		}
		info, err := f.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(path, f.Name())
		pid := PathID(fullPath)
		sizeStr := HumanSize(info.Size())
		buttons = append(buttons, []interface{}{
			map[string]string{"text": fmt.Sprintf("📄 %s (%s)", f.Name(), sizeStr), "callback_data": "fd:" + pid},
		})
	}

	// Navigation
	navRow := []interface{}{}
	parent := filepath.Dir(path)
	if parent != path {
		ppid := PathID(parent)
		navRow = append(navRow, map[string]string{"text": "⬆️ Up", "callback_data": "fb:" + ppid})
	}
	navRow = append(navRow, map[string]string{"text": "📂 Roots", "callback_data": "files"})
	navRow = append(navRow, map[string]string{"text": "◀️ Menu", "callback_data": "menu"})
	buttons = append(buttons, navRow)

	return strings.Join(lines, "\n"), map[string]interface{}{"inline_keyboard": buttons}
}

func FileDetailKeyboard(path string) (string, map[string]interface{}) {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err), BackKB("files")
	}

	name := filepath.Base(path)
	displayPath := strings.Replace(path, config.HomeDir(), "~", 1)
	modified := info.ModTime().Format("2006-01-02 15:04")

	lines := []string{
		fmt.Sprintf("📄 <b>%s</b>", name),
		fmt.Sprintf("📂 %s", displayPath),
		fmt.Sprintf("📏 Size: %s", HumanSize(info.Size())),
		fmt.Sprintf("🕐 Modified: %s", modified),
	}

	pid := PathID(path)
	var buttons []interface{}
	if info.Size() <= maxTGFile {
		buttons = append(buttons, []interface{}{
			map[string]string{"text": "⬇️ Download", "callback_data": "dl:" + pid},
		})
	} else {
		lines = append(lines, "\n⚠️ Too large for Telegram (max 50MB)")
	}

	ppid := PathID(filepath.Dir(path))
	buttons = append(buttons, []interface{}{
		map[string]string{"text": "⬆️ Back", "callback_data": "fb:" + ppid},
		map[string]string{"text": "◀️ Menu", "callback_data": "menu"},
	})

	return strings.Join(lines, "\n"), map[string]interface{}{"inline_keyboard": buttons}
}

func FolderZipInfo(path string) (string, map[string]interface{}) {
	displayPath := strings.Replace(path, config.HomeDir(), "~", 1)

	// Single walk for both size and count
	var totalSize int64
	fileCount := 0
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	lines := []string{
		"📦 <b>Download Folder as ZIP</b>",
		fmt.Sprintf("📂 %s", displayPath),
		fmt.Sprintf("📏 Size: ~%s (%d files)", HumanSize(totalSize), fileCount),
	}

	pid := PathID(path)
	var buttons []interface{}

	if totalSize > confirmThreshold {
		parts := totalSize/partSize + 1
		lines = append(lines, fmt.Sprintf("\n⚠️ <b>Large folder (%s)</b>", HumanSize(totalSize)))
		lines = append(lines, fmt.Sprintf("This will be split into %d parts of 20MB", parts))
		lines = append(lines, "Are you sure you want to download?")
		buttons = append(buttons, []interface{}{
			map[string]string{"text": "✅ Yes, Download", "callback_data": "zc:" + pid},
		})
	} else if totalSize > partSize {
		parts := totalSize/partSize + 1
		lines = append(lines, fmt.Sprintf("\n📦 Will be split into ~%d parts of 20MB", parts))
		buttons = append(buttons, []interface{}{
			map[string]string{"text": "⬇️ Download ZIP (split)", "callback_data": "zc:" + pid},
		})
	} else {
		buttons = append(buttons, []interface{}{
			map[string]string{"text": "⬇️ Download ZIP", "callback_data": "zc:" + pid},
		})
	}

	ppid := PathID(filepath.Dir(path))
	buttons = append(buttons, []interface{}{
		map[string]string{"text": "⬆️ Back", "callback_data": "fb:" + ppid},
		map[string]string{"text": "❌ Cancel", "callback_data": "files"},
	})

	return strings.Join(lines, "\n"), map[string]interface{}{"inline_keyboard": buttons}
}

func BackKB(target string) map[string]interface{} {
	return map[string]interface{}{
		"inline_keyboard": []interface{}{
			[]interface{}{
				map[string]string{"text": "◀️ Back", "callback_data": target},
			},
		},
	}
}

// ──────────────────────────────────────────────
// File Send (download to Telegram)
// ──────────────────────────────────────────────

func SendFile(chatID string, filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil || info.Size() > maxTGFile {
		return false
	}
	name := filepath.Base(filePath)
	caption := fmt.Sprintf("📄 %s\n📏 %s", name, HumanSize(info.Size()))
	ok, _ := SendMediaFile(chatID, filePath, caption, false)
	return ok
}

// ──────────────────────────────────────────────
// Folder ZIP + split
// ──────────────────────────────────────────────

func SendFolderAsZip(chatID string, folderPath string) bool {
	folderName := filepath.Base(folderPath)

	SendMessage(fmt.Sprintf("📦 Mengzip <b>%s</b>...", folderName), nil)

	// Create ZIP in temp
	tmpDir, err := os.MkdirTemp("", "scorp_zip_")
	if err != nil {
		SendMessage(fmt.Sprintf("❌ ZIP error: %v", err), nil)
		return false
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, folderName+".zip")
	if err := createZip(folderPath, zipPath); err != nil {
		SendMessage(fmt.Sprintf("❌ ZIP error: %v", err), nil)
		return false
	}

	zipInfo, _ := os.Stat(zipPath)
	zipSize := zipInfo.Size()
	log.Printf("[files] ZIP created: %s (%s)", zipPath, HumanSize(zipSize))

	if zipSize <= maxTGFile {
		ok := SendFile(chatID, zipPath)
		if ok {
			SendMessage(fmt.Sprintf("✅ <b>%s.zip</b> sent (%s)", folderName, HumanSize(zipSize)), nil)
		}
		return ok
	}

	// Split
	parts := SplitFile(zipPath, partSize, tmpDir)
	SendMessage(fmt.Sprintf("📦 <b>%s.zip</b> = %s\nSplit into %d parts ~20MB...", folderName, HumanSize(zipSize), len(parts)), nil)

	for i, partPath := range parts {
		pInfo, _ := os.Stat(partPath)
		caption := fmt.Sprintf("📦 Part %d/%d — %s", i+1, len(parts), HumanSize(pInfo.Size()))
		if !SendDocument(chatID, partPath, caption) {
			SendMessage(fmt.Sprintf("❌ Failed to send part %d", i+1), nil)
			return false
		}
		time.Sleep(1 * time.Second)
	}

	SendMessage(fmt.Sprintf("✅ <b>%s.zip</b> sent %d parts\n📏 Total: %s\n\n💡 Combine:\n<code>cat %s.zip.part* > %s.zip</code>", folderName, len(parts), HumanSize(zipSize), folderName, folderName), nil)
	return true
}

func createZip(folderPath, zipPath string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	folderName := filepath.Base(folderPath)
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(folderPath, path)
		arcName := filepath.Join(folderName, relPath)

		fw, err := w.Create(arcName)
		if err != nil {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()
		io.Copy(fw, file)
		return nil
	})
}

func SplitFile(filePath string, chunkSize int64, outDir string) []string {
	var parts []string
	f, err := os.Open(filePath)
	if err != nil {
		return parts
	}
	defer f.Close()

	baseName := filepath.Base(filePath)
	buf := make([]byte, chunkSize)
	partNum := 1

	for {
		n, err := f.Read(buf)
		if n == 0 {
			break
		}

		partPath := filepath.Join(outDir, fmt.Sprintf("%s.part%02d", baseName, partNum))
		os.WriteFile(partPath, buf[:n], 0644)
		parts = append(parts, partPath)
		partNum++

		if err == io.EOF {
			break
		}
	}
	return parts
}

// sendTelegramMultipart sends a local file to a specific Telegram Bot API endpoint via multipart/form-data.
func sendTelegramMultipart(endpoint string, fieldName string, chatID string, filePath string, caption string) (bool, string) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Sprintf("failed to open file: %v", err)
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("chat_id", chatID); err != nil {
		return false, err.Error()
	}
	if caption != "" {
		_ = writer.WriteField("caption", caption)
	}

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return false, err.Error()
	}
	if _, err := io.Copy(part, f); err != nil {
		return false, err.Error()
	}
	writer.Close()

	client := HttpLong
	resp, err := client.Post(TgBase+"/"+endpoint, writer.FormDataContentType(), body)
	if err != nil {
		return false, fmt.Sprintf("network error: %v", err)
	}
	defer resp.Body.Close()

	var tgResp TgResponse
	if err := json.NewDecoder(resp.Body).Decode(&tgResp); err != nil {
		return false, fmt.Sprintf("decode error: %v (HTTP %d)", err, resp.StatusCode)
	}
	if !tgResp.OK {
		return false, tgResp.Description
	}
	return true, ""
}

// SendMediaFile sends ANY file to Telegram, choosing the best representation:
// - Photos (.png, .jpg, .jpeg, .webp, .bmp) -> sendPhoto (with fallback to sendDocument)
// - Videos (.mp4, .mov, .webm, .mkv, .avi) -> sendVideo (with fallback to sendDocument)
// - Audio (.mp3, .ogg, .wav, .m4a, .flac) -> sendAudio (with fallback to sendDocument)
// - Animations (.gif) -> sendAnimation (with fallback to sendDocument)
// - Documents/code/archives/all other files -> sendDocument
// If asDocument is true, forces sending as an uncompressed document attachment.
func SendMediaFile(chatID string, filePath string, caption string, asDocument bool) (bool, string) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Sprintf("file not found: %v", err)
	}
	if info.IsDir() {
		return false, fmt.Sprintf("'%s' is a directory, not a file", filePath)
	}
	if info.Size() > maxTGFile {
		return false, fmt.Sprintf("file too large (%s) — Telegram bot limit is 50MB", HumanSize(info.Size()))
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	// If forced as document or non-media extension
	if asDocument {
		ok, desc := sendTelegramMultipart("sendDocument", "document", chatID, filePath, caption)
		if !ok {
			return false, desc
		}
		return true, "document"
	}

	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".bmp":
		ok, desc := sendTelegramMultipart("sendPhoto", "photo", chatID, filePath, caption)
		if ok {
			return true, "photo"
		}
		log.Printf("[telegram] sendPhoto failed (%s) — falling back to sendDocument", desc)
		okDoc, descDoc := sendTelegramMultipart("sendDocument", "document", chatID, filePath, caption)
		if okDoc {
			return true, "document (photo fallback)"
		}
		return false, descDoc

	case ".mp4", ".mov", ".webm", ".mkv", ".avi":
		ok, desc := sendTelegramMultipart("sendVideo", "video", chatID, filePath, caption)
		if ok {
			return true, "video"
		}
		log.Printf("[telegram] sendVideo failed (%s) — falling back to sendDocument", desc)
		okDoc, descDoc := sendTelegramMultipart("sendDocument", "document", chatID, filePath, caption)
		if okDoc {
			return true, "document (video fallback)"
		}
		return false, descDoc

	case ".mp3", ".ogg", ".wav", ".m4a", ".flac", ".aac":
		ok, desc := sendTelegramMultipart("sendAudio", "audio", chatID, filePath, caption)
		if ok {
			return true, "audio"
		}
		log.Printf("[telegram] sendAudio failed (%s) — falling back to sendDocument", desc)
		okDoc, descDoc := sendTelegramMultipart("sendDocument", "document", chatID, filePath, caption)
		if okDoc {
			return true, "document (audio fallback)"
		}
		return false, descDoc

	case ".gif":
		ok, desc := sendTelegramMultipart("sendAnimation", "animation", chatID, filePath, caption)
		if ok {
			return true, "animation"
		}
		log.Printf("[telegram] sendAnimation failed (%s) — falling back to sendDocument", desc)
		okDoc, descDoc := sendTelegramMultipart("sendDocument", "document", chatID, filePath, caption)
		if okDoc {
			return true, "document (animation fallback)"
		}
		return false, descDoc

	default:
		// Any document, archive, code file, PDF, ZIP, etc.
		ok, desc := sendTelegramMultipart("sendDocument", "document", chatID, filePath, caption)
		if !ok {
			return false, desc
		}
		return true, "document"
	}
}

// SendMediaBytes sends raw in-memory bytes as a media file to Telegram.
func SendMediaBytes(chatID string, data []byte, filename string, caption string, asDocument bool) (bool, string) {
	tmpDir, err := os.MkdirTemp("", "scorp_send_")
	if err != nil {
		return false, err.Error()
	}
	defer os.RemoveAll(tmpDir)

	tmpPath := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return false, err.Error()
	}

	return SendMediaFile(chatID, tmpPath, caption, asDocument)
}

// SendDocumentBytes sends raw bytes as a document
func SendDocumentBytes(chatID string, data []byte, filename string, caption string) bool {
	ok, _ := SendMediaBytes(chatID, data, filename, caption, true)
	return ok
}
