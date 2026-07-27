package collectors

import (
	"net/http"
	"time"
)

// httpShort is a local HTTP client for the collectors package.
var httpShort = &http.Client{
	Timeout: 15 * time.Second,
}
