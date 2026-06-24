package tools

import (
	"net/http"
	"time"
)

// Common constants used across tool implementations

// HTTP clients for tool use
var HttpShort = &http.Client{
	Timeout: 15 * time.Second,
}

var HttpLong = &http.Client{
	Timeout: 5 * time.Minute,
}
