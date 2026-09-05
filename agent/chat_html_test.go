package agent

import (
	"strings"
	"testing"
)

func TestMarkdownToTelegramHTMLEscaping(t *testing.T) {
	// Raw comparison operators must be escaped exactly once (rendered as < >).
	out := markdownToTelegramHTML("Output: `10 < 20 && 30 > 5`")
	if !strings.Contains(out, "<code>10 &lt; 20 &amp;&amp; 30 &gt; 5</code>") {
		t.Errorf("code span not single-escaped: %s", out)
	}
	if strings.Contains(out, "&amp;lt;") {
		t.Errorf("double escaping detected: %s", out)
	}

	// Stray invalid tags must be neutralized so Telegram's parser accepts.
	out = markdownToTelegramHTML("see <name> tag")
	if strings.Contains(out, "<name>") {
		t.Errorf("invalid tag not escaped: %s", out)
	}

	// Native Telegram tags the model emits stay functional.
	out = markdownToTelegramHTML("<b>bold</b> and <a href=\"https://x.co\">link</a>")
	if !strings.Contains(out, "<b>bold</b>") || !strings.Contains(out, `<a href="https://x.co">link</a>`) {
		t.Errorf("native tags broken: %s", out)
	}

	// Markdown still converts.
	out = markdownToTelegramHTML("**bold** and *ital*")
	if !strings.Contains(out, "<b>bold</b>") || !strings.Contains(out, "<i>ital</i>") {
		t.Errorf("markdown conversion broken: %s", out)
	}
}
