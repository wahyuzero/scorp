package tools

import (
	"os"
	"strings"
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	// 1. Test OpenAI sk- key
	input1 := "Error with key sk-1234567890abcdefghijklmnopqrstuvwxyz and user_1234567890abcdefghijklmnopqrstuvwxyz1234567890"
	redacted1 := RedactSecrets(input1)
	if strings.Contains(redacted1, "sk-1234567890") {
		t.Errorf("Expected sk- key to be redacted, got: %s", redacted1)
	}
	if strings.Contains(redacted1, "user_1234567890") {
		t.Errorf("Expected user_ key to be redacted, got: %s", redacted1)
	}

	// 2. Test GitHub token
	input2 := "Clone failed with ghp_1234567890abcdefghijklmnopqrstuvwxyz"
	redacted2 := RedactSecrets(input2)
	if strings.Contains(redacted2, "ghp_1234567890") {
		t.Errorf("Expected ghp_ token to be redacted, got: %s", redacted2)
	}

	// 3. Test Private Key
	input3 := `Connecting with:
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0Y1234567890
-----END RSA PRIVATE KEY-----
Connection done.`
	redacted3 := RedactSecrets(input3)
	if strings.Contains(redacted3, "MIIEowIBAAKCAQEA0Y1234567890") {
		t.Errorf("Expected private key to be redacted, got: %s", redacted3)
	}

	// 4. Test Key-Value secret assignment
	input4 := `DATABASE_URL=postgres://user:secretpassword123@localhost/db
PASSWORD="mypassword999"`
	redacted4 := RedactSecrets(input4)
	if strings.Contains(redacted4, "mypassword999") {
		t.Errorf("Expected password assignment to be redacted, got: %s", redacted4)
	}

	// 5. Test active env secret redaction
	os.Setenv("COMMAND_CODE_API_KEY", "super_secret_production_key_xyz123")
	input5 := "Executing request with super_secret_production_key_xyz123 inside headers"
	redacted5 := RedactSecrets(input5)
	if strings.Contains(redacted5, "super_secret_production_key_xyz123") {
		t.Errorf("Expected active env var to be redacted, got: %s", redacted5)
	}
}
