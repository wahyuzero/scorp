//go:build !linux

package collectors

// StartCPUSampler is a no-op fallback on non-Linux platforms
func StartCPUSampler(done <-chan struct{}) {}

// GetTopProcesses returns an empty list on non-Linux platforms
func GetTopProcesses(n int) []TopProcess {
	return nil
}

// CollectSystem returns an empty SystemData fallback on non-Linux platforms
func CollectSystem() SystemData {
	return SystemData{}
}
