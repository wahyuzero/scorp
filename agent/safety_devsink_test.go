package agent

import "testing"

func TestDevNullRedirectionNotDangerous(t *testing.T) {
	safe := []string{
		"ss -tlnp 2>/dev/null | tail -n +2",
		"journalctl -u scorp.service -n 300 2>/dev/null > /tmp/x.txt",
		"ps aux --sort=-%mem 2>/dev/null || true",
		"command 2>/dev/null 1>/dev/null",
		"curl -s https://x >/dev/null",
		"dd if=/root/f of=/dev/null bs=1M",
	}
	for _, cmd := range safe {
		if IsDangerousCommand(cmd) {
			t.Errorf("benign command flagged dangerous: %s", cmd)
		}
	}
	unsafe := []string{
		"dd if=/dev/zero of=/dev/sda",
		"echo x > /dev/sda",
		"echo x >/dev/sdb1",
	}
	for _, cmd := range unsafe {
		if !IsDangerousCommand(cmd) {
			t.Errorf("dangerous device write not flagged: %s", cmd)
		}
	}
}
