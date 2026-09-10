package netdns

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Public fallback DNS servers
var PublicDNSServers = []string{
	"1.1.1.1:53", // Cloudflare primary
	"8.8.8.8:53", // Google primary
	"9.9.9.9:53", // Quad9 primary
	"1.0.0.1:53", // Cloudflare secondary
	"8.8.4.4:53", // Google secondary
}

var (
	defaultResolver *net.Resolver
	resolverOnce    sync.Once

	cachedPlatformDNS []string
	platformDNSOnce   sync.Once

	loopbackChecked int32
	loopbackWorking int32
)

func init() {
	Init()
}

// Init configures net.DefaultResolver and http.DefaultTransport early at startup
func Init() {
	r := GetResolver()
	net.DefaultResolver = r

	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t.DialContext = NewDialer().DialContext
		t.ForceAttemptHTTP2 = true
	}
}

// GetResolver returns the singleton resilient net.Resolver
func GetResolver() *net.Resolver {
	resolverOnce.Do(func() {
		defaultResolver = &net.Resolver{
			PreferGo: true,
			Dial:     DialDNS,
		}
	})
	return defaultResolver
}

// NewDialer returns a net.Dialer wired with the resilient DNS resolver
func NewDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver:  GetResolver(),
	}
}

// NewDialContext returns a DialContext func wired with resilient DNS resolution
func NewDialContext() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return NewDialer().DialContext
}

// IsLoopback reports whether the host in addr is a loopback address
func IsLoopback(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	host = strings.Trim(host, "[]")
	if host == "localhost" || host == "127.0.0.1" || host == "::1" || strings.HasPrefix(host, "127.") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// isIPv6LinkLocal checks if address is an IPv6 link-local address (e.g. fe80::1)
func isIPv6LinkLocal(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	host = strings.Trim(host, "[]")
	if idx := strings.Index(host, "%"); idx != -1 {
		host = host[:idx]
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLinkLocalUnicast()
}

// ResetLoopbackProbe resets cached loopback status (for unit testing)
func ResetLoopbackProbe() {
	atomic.StoreInt32(&loopbackChecked, 0)
	atomic.StoreInt32(&loopbackWorking, 0)
}

// SetLoopbackState sets cached loopback status (for unit testing)
func SetLoopbackState(working bool) {
	if working {
		atomic.StoreInt32(&loopbackWorking, 1)
	} else {
		atomic.StoreInt32(&loopbackWorking, 0)
	}
	atomic.StoreInt32(&loopbackChecked, 1)
}

// IsLoopbackHealthy probes loopback port 53 to verify if a local DNS resolver is active
func IsLoopbackHealthy(timeout time.Duration) bool {
	if atomic.LoadInt32(&loopbackChecked) == 1 {
		return atomic.LoadInt32(&loopbackWorking) == 1
	}

	if timeout <= 0 {
		timeout = 40 * time.Millisecond
	}

	working := probeDNS("127.0.0.1:53", timeout)
	if !working {
		working = probeDNS("[::1]:53", timeout)
	}

	if working {
		atomic.StoreInt32(&loopbackWorking, 1)
	} else {
		atomic.StoreInt32(&loopbackWorking, 0)
	}
	atomic.StoreInt32(&loopbackChecked, 1)
	return working
}

// probeDNS sends a lightweight DNS query to check if port 53 is responding
func probeDNS(server string, timeout time.Duration) bool {
	d := net.Dialer{Timeout: timeout}
	c, err := d.Dial("udp", server)
	if err != nil {
		return false
	}
	defer c.Close()

	_ = c.SetDeadline(time.Now().Add(timeout))
	query := []byte{
		0x12, 0x34, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
	}
	if _, err := c.Write(query); err != nil {
		return false
	}
	buf := make([]byte, 512)
	n, err := c.Read(buf)
	return err == nil && n > 0
}

// DiscoverPlatformDNS attempts to find DNS servers configured by Android or Termux
func DiscoverPlatformDNS() []string {
	platformDNSOnce.Do(func() {
		var servers []string

		// 1. Termux resolv.conf
		termuxResolv := "/data/data/com.termux/files/usr/etc/resolv.conf"
		if data, err := os.ReadFile(termuxResolv); err == nil {
			servers = append(servers, parseResolvConf(string(data))...)
		}

		// 2. Android system properties via getprop net.dns1 / net.dns2
		for _, prop := range []string{"net.dns1", "net.dns2"} {
			if ip := getprop(prop); ip != "" {
				if parsed := net.ParseIP(ip); parsed != nil {
					servers = append(servers, net.JoinHostPort(ip, "53"))
				}
			}
		}

		// 3. System /etc/resolv.conf (only non-loopback, non link-local servers)
		if data, err := os.ReadFile("/etc/resolv.conf"); err == nil {
			for _, s := range parseResolvConf(string(data)) {
				if !IsLoopback(s) && !isIPv6LinkLocal(s) {
					servers = append(servers, s)
				}
			}
		}

		cachedPlatformDNS = deduplicate(servers)
	})

	return cachedPlatformDNS
}

// parseResolvConf extracts nameserver IP addresses from resolv.conf contents
func parseResolvConf(content string) []string {
	var servers []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "nameserver" {
			ip := fields[1]
			if !strings.Contains(ip, ":") || strings.Count(ip, ":") == 1 {
				// IPv4 or host:port
				if _, _, err := net.SplitHostPort(ip); err != nil {
					ip = net.JoinHostPort(ip, "53")
				}
			} else {
				// IPv6
				if !strings.HasPrefix(ip, "[") {
					ip = "[" + ip + "]:53"
				}
			}
			servers = append(servers, ip)
		}
	}
	return servers
}

// getprop retrieves an Android system property
func getprop(prop string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cmdPaths := []string{"/system/bin/getprop", "getprop"}
	for _, p := range cmdPaths {
		out, err := exec.CommandContext(ctx, p, prop).Output()
		if err == nil {
			val := strings.TrimSpace(string(out))
			if val != "" {
				return val
			}
		}
	}
	return ""
}

func deduplicate(list []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range list {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// BuildCandidateList compiles the ordered list of DNS servers to try
func BuildCandidateList(address string) []string {
	var candidates []string
	var fallbacks []string

	// Platform discovered DNS (Android / Termux / valid resolv.conf)
	discovered := DiscoverPlatformDNS()
	fallbacks = append(fallbacks, discovered...)

	// Standard public DNS fallbacks
	fallbacks = append(fallbacks, PublicDNSServers...)
	fallbacks = deduplicate(fallbacks)

	if IsLoopback(address) {
		if IsLoopbackHealthy(40 * time.Millisecond) {
			candidates = append(candidates, address)
			candidates = append(candidates, fallbacks...)
		} else {
			// Loopback is dead (typical Android/Termux case): bypass directly to working fallbacks
			candidates = append(candidates, fallbacks...)
			candidates = append(candidates, address) // keep at end as last resort
		}
	} else if isIPv6LinkLocal(address) {
		// Link-local IPv6 DNS servers often blackhole; prioritize public fallbacks
		candidates = append(candidates, fallbacks...)
		candidates = append(candidates, address)
	} else {
		candidates = append(candidates, address)
		candidates = append(candidates, fallbacks...)
	}

	return deduplicate(candidates)
}

// DialDNS dials DNS servers with transparent failover to platform and public DNS servers
func DialDNS(ctx context.Context, network, address string) (net.Conn, error) {
	candidates := BuildCandidateList(address)
	dialer := &net.Dialer{Timeout: 3 * time.Second}

	var firstConn net.Conn
	var firstErr error

	for _, cand := range candidates {
		dctx := ctx
		if dctx == nil || dctx.Err() != nil {
			dctx = context.Background()
		}
		c, err := dialer.DialContext(dctx, network, cand)
		if err == nil {
			firstConn = c
			break
		}
		firstErr = err
	}

	if firstConn == nil {
		return nil, firstErr
	}

	return firstConn, nil
}

// ResilientDNSConn wraps a net.Conn with automatic packet-replay failover on read/write failure
type ResilientDNSConn struct {
	net.Conn
	ctx           context.Context
	network       string
	dialer        *net.Dialer
	fallbackAddrs []string
	lastWritten   []byte
	deadline      time.Time
	mu            sync.Mutex
}

func (c *ResilientDNSConn) SetDeadline(t time.Time) error {
	c.mu.Lock()
	c.deadline = t
	c.mu.Unlock()
	return c.Conn.SetDeadline(t)
}

func (c *ResilientDNSConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *ResilientDNSConn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

func (c *ResilientDNSConn) Write(b []byte) (int, error) {
	c.mu.Lock()
	c.lastWritten = make([]byte, len(b))
	copy(c.lastWritten, b)
	conn := c.Conn
	c.mu.Unlock()

	n, err := conn.Write(b)
	if err != nil {
		if c.tryNextFallback() {
			c.mu.Lock()
			active := c.Conn
			c.mu.Unlock()
			return active.Write(b)
		}
	}
	return n, err
}

func (c *ResilientDNSConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if err != nil {
		// Read error on UDP (e.g. connection refused / ICMP port unreachable or timeout)
		for c.tryNextFallback() {
			c.mu.Lock()
			written := c.lastWritten
			active := c.Conn
			c.mu.Unlock()

			if len(written) > 0 {
				if _, wErr := active.Write(written); wErr == nil {
					n2, rErr := active.Read(b)
					if rErr == nil {
						return n2, nil
					}
				}
			}
		}
	}
	return n, err
}

func (c *ResilientDNSConn) tryNextFallback() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for len(c.fallbackAddrs) > 0 {
		nextAddr := c.fallbackAddrs[0]
		c.fallbackAddrs = c.fallbackAddrs[1:]

		// Fallback must use fresh background context to avoid canceled context from dead server
		dctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
		newConn, err := c.dialer.DialContext(dctx, c.network, nextAddr)
		cancel()

		if err == nil {
			if c.Conn != nil {
				_ = c.Conn.Close()
			}
			c.Conn = newConn
			_ = newConn.SetDeadline(time.Now().Add(2500 * time.Millisecond))
			return true
		}
	}
	return false
}
