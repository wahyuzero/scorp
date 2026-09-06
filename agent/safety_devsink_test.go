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

func TestDevTcpPseudoDeviceAllowed(t *testing.T) {
	safe := []string{
		`timeout 2 bash -c 'echo > /dev/tcp/127.0.0.1/22'`,
		`cat < /dev/null > /dev/tcp/google.com/443`,
		`exec 3<>/dev/tcp/localhost/8080`,
		`echo ping > /dev/udp/1.2.3.4/53`,
	}
	for _, cmd := range safe {
		if IsDangerousCommand(cmd) {
			t.Errorf("network pseudo-device should be allowed: %q", cmd)
		}
	}
	// Real device overwrites must stay gated.
	if !IsDangerousCommand("echo x > /dev/sda") {
		t.Error("overwrite of /dev/sda must stay gated")
	}
}
