package collectors

import (
	"net/http"
	"time"
)

// httpShort is a local HTTP client for collectors package.
// In a future phase, this will be shared via a common transport package.
var httpShort = &http.Client{
	Timeout: 15 * time.Second,
}
