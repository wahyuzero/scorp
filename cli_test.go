package main

import (
	"strings"
	"testing"
)

func TestStripHTML(t *testing.T) {
	input := "<b>Hello</b> &amp; <i>World</i> &lt;test&gt; <pre><code>code block</code></pre>"
	output := stripHTML(input)

	if strings.Contains(output, "<b>") || strings.Contains(output, "</b>") {
		t.Errorf("expected no HTML tags, got: %s", output)
	}
	if !strings.Contains(output, "&") {
		t.Errorf("expected &amp; to be unescaped to &, got: %s", output)
	}
	if !strings.Contains(output, "<test>") {
		t.Errorf("expected &lt;test&gt; to be unescaped, got: %s", output)
	}
	if !strings.Contains(output, "```") {
		t.Errorf("expected code block backticks in output, got: %s", output)
	}
}

func TestFormatFinalResponse(t *testing.T) {
	input := "🤖 <b>Scorp:</b>\n\nIni adalah jawaban final."
	output := formatFinalResponse(input)

	if strings.HasPrefix(output, "🤖") || strings.HasPrefix(output, "Scorp:") {
		t.Errorf("expected Scorp header prefix to be stripped, got: %s", output)
	}
	if !strings.Contains(output, "Ini adalah jawaban final.") {
		t.Errorf("expected main text preserved, got: %s", output)
	}
}
