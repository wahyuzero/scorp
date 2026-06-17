package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
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
	{homeDir(), "🏠 Home"},
	{homeDir() + "/projects", "📁 Projects"},
	{"/data/coolify", "☁️ Coolify Data"},
	{"/tmp", "📋 Temp"},
}

// Path ID mapping
var (
	pathMap     = make(map[string]string) // pid -> path
	reversePath = make(map[string]string) // path -> pid
	pathCounter int
	pathMu      sync.Mutex
)

func pathID(path string) string {
	pathMu.Lock()
	defer pathMu.Unlock()
	if pid, ok := reversePath[path]; ok {
		return pid
	}
	pathCounter++
	pid := fmt.Sprintf("p%d", pathCounter)
	pathMap[pid] = path
	reversePath[path] = pid
	return pid
}

func getPath(pid string) string {
	pathMu.Lock()
	defer pathMu.Unlock()
	return pathMap[pid]
}

func humanSize(size int64) string {
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

func rootsKeyboard() map[string]interface{} {
	var buttons []interface{}
	for _, r := range browseRoots {
		pid := pathID(r.Path)
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

func dirKeyboard(path string) (string, map[string]interface{}) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "🚫 Permission denied", backKB("files")
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	displayPath := strings.Replace(path, homeDir(), "~", 1)
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
		pid := pathID(fullPath)
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
		pid := pathID(fullPath)
		sizeStr := humanSize(info.Size())
		buttons = append(buttons, []interface{}{
			map[string]string{"text": fmt.Sprintf("📄 %s (%s)", f.Name(), sizeStr), "callback_data": "fd:" + pid},
		})
	}

	// Navigation
	navRow := []interface{}{}
	parent := filepath.Dir(path)
	if parent != path {
		ppid := pathID(parent)
		navRow = append(navRow, map[string]string{"text": "⬆️ Up", "callback_data": "fb:" + ppid})
	}
	navRow = append(navRow, map[string]string{"text": "📂 Roots", "callback_data": "files"})
	navRow = append(navRow, map[string]string{"text": "◀️ Menu", "callback_data": "menu"})
	buttons = append(buttons, navRow)

	return strings.Join(lines, "\n"), map[string]interface{}{"inline_keyboard": buttons}
}

func fileDetailKeyboard(path string) (string, map[string]interface{}) {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("❌ Error: %v", err), backKB("files")
	}

	name := filepath.Base(path)
	displayPath := strings.Replace(path, homeDir(), "~", 1)
	modified := info.ModTime().Format("2006-01-02 15:04")

	lines := []string{
		fmt.Sprintf("📄 <b>%s</b>", name),
		fmt.Sprintf("📂 %s", displayPath),
		fmt.Sprintf("📏 Size: %s", humanSize(info.Size())),
		fmt.Sprintf("🕐 Modified: %s", modified),
	}

	pid := pathID(path)
	var buttons []interface{}
	if info.Size() <= maxTGFile {
		buttons = append(buttons, []interface{}{
			map[string]string{"text": "⬇️ Download", "callback_data": "dl:" + pid},
		})
	} else {
		lines = append(lines, "\n⚠️ Too large for Telegram (max 50MB)")
	}

	ppid := pathID(filepath.Dir(path))
	buttons = append(buttons, []interface{}{
		map[string]string{"text": "⬆️ Back", "callback_data": "fb:" + ppid},
		map[string]string{"text": "◀️ Menu", "callback_data": "menu"},
	})

	return strings.Join(lines, "\n"), map[string]interface{}{"inline_keyboard": buttons}
}

func folderZipInfo(path string) (string, map[string]interface{}) {
	displayPath := strings.Replace(path, homeDir(), "~", 1)

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
		fmt.Sprintf("📏 Size: ~%s (%d files)", humanSize(totalSize), fileCount),
	}

	pid := pathID(path)
	var buttons []interface{}

	if totalSize > confirmThreshold {
		parts := totalSize/partSize + 1
		lines = append(lines, fmt.Sprintf("\n⚠️ <b>Large folder (%s)</b>", humanSize(totalSize)))
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

	ppid := pathID(filepath.Dir(path))
	buttons = append(buttons, []interface{}{
		map[string]string{"text": "⬆️ Back", "callback_data": "fb:" + ppid},
		map[string]string{"text": "❌ Cancel", "callback_data": "files"},
	})

	return strings.Join(lines, "\n"), map[string]interface{}{"inline_keyboard": buttons}
}

func backKB(target string) map[string]interface{} {
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

func sendFile(chatID string, filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil || info.Size() > maxTGFile {
		return false
	}
	name := filepath.Base(filePath)
	caption := fmt.Sprintf("📄 %s\n📏 %s", name, humanSize(info.Size()))
	return sendDocument(chatID, filePath, caption)
}

// ──────────────────────────────────────────────
// Folder ZIP + split
// ──────────────────────────────────────────────

func sendFolderAsZip(chatID string, folderPath string) bool {
	folderName := filepath.Base(folderPath)

	sendMessage(fmt.Sprintf("📦 Zipping <b>%s</b>...", folderName), nil)

	// Create ZIP in temp
	tmpDir, err := os.MkdirTemp("", "vps_mon_zip_")
	if err != nil {
		sendMessage(fmt.Sprintf("❌ ZIP error: %v", err), nil)
		return false
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, folderName+".zip")
	if err := createZip(folderPath, zipPath); err != nil {
		sendMessage(fmt.Sprintf("❌ ZIP error: %v", err), nil)
		return false
	}

	zipInfo, _ := os.Stat(zipPath)
	zipSize := zipInfo.Size()
	log.Printf("[files] ZIP created: %s (%s)", zipPath, humanSize(zipSize))

	if zipSize <= maxTGFile {
		ok := sendFile(chatID, zipPath)
		if ok {
			sendMessage(fmt.Sprintf("✅ <b>%s.zip</b> sent (%s)", folderName, humanSize(zipSize)), nil)
		}
		return ok
	}

	// Split
	parts := splitFile(zipPath, partSize, tmpDir)
	sendMessage(fmt.Sprintf("📦 <b>%s.zip</b> = %s\nSplitting into %d parts of ~20MB...",
		folderName, humanSize(zipSize), len(parts)), nil)

	for i, partPath := range parts {
		pInfo, _ := os.Stat(partPath)
		caption := fmt.Sprintf("📦 Part %d/%d — %s", i+1, len(parts), humanSize(pInfo.Size()))
		if !sendDocument(chatID, partPath, caption) {
			sendMessage(fmt.Sprintf("❌ Failed to send part %d", i+1), nil)
			return false
		}
		time.Sleep(1 * time.Second)
	}

	sendMessage(fmt.Sprintf("✅ <b>%s.zip</b> sent in %d parts\n📏 Total: %s\n\n💡 To reassemble:\n<code>cat %s.zip.part* > %s.zip</code>",
		folderName, len(parts), humanSize(zipSize), folderName, folderName), nil)
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

func splitFile(filePath string, chunkSize int64, outDir string) []string {
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

// sendDocumentBytes sends raw bytes as a document
func sendDocumentBytes(chatID string, data []byte, filename string, caption string) bool {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("chat_id", chatID)
	if caption != "" {
		writer.WriteField("caption", caption)
	}

	part, err := writer.CreateFormFile("document", filename)
	if err != nil {
		return false
	}
	part.Write(data)
	writer.Close()

	client := httpLong
	resp, err := client.Post(tgBase+"/sendDocument", writer.FormDataContentType(), body)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
