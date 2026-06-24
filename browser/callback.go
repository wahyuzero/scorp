package browser

// Callback for sending files from browser package (set by root)
var SendFile func(chatID, filePath string) bool
