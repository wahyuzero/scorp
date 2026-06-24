# Scorp Agent Build Targets
#
#   make              → production (stripped)
#   make minimal      → no browser
#   make debug        → debug build
#   make deploy       → build + install + restart
#
# Version injection:
#   make VERSION=v1.0.0   or   make deploy VERSION=v1.0.0
#   Defaults to git short hash if not specified.

GO ?= /usr/local/go/bin/go
BIN := scorp

# Auto-detect version from git tag, or use VERSION env, or fallback to "dev"
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags="-s -w -buildid= -X scorp-agent/updater.Version=$(VERSION)" -trimpath

default: $(BIN)

$(BIN): *.go */*.go
	$(GO) build $(LDFLAGS) -o $(BIN) .

minimal:
	$(GO) build -tags nobrowser $(LDFLAGS) -o $(BIN)-minimal .

debug:
	$(GO) build -o $(BIN)-debug .

deploy: $(BIN)
	sudo cp $(BIN) /usr/local/bin/$(BIN)
	sudo systemctl restart $(BIN)

deploy-minimal: minimal
	sudo cp $(BIN)-minimal /usr/local/bin/$(BIN)
	sudo systemctl restart $(BIN)

clean:
	rm -f $(BIN) $(BIN)-* ./*.bak

.PHONY: default minimal debug deploy deploy-minimal clean
