# Scorp Agent Build Targets
#
#   make              → production (12M, browser included)
#   make minimal      → no browser (9.7M)
#   make debug        → debug build (17M)

GO ?= /usr/local/go/bin/go
BIN := scorp-agent

LDFLAGS := -ldflags="-s -w -buildid=" -trimpath

default: $(BIN)

$(BIN): *.go
	$(GO) build $(LDFLAGS) -o $(BIN) .

minimal:
	$(GO) build -tags nobrowser $(LDFLAGS) -o $(BIN)-minimal .

debug:
	$(GO) build -o $(BIN)-debug .

deploy: $(BIN)
	sudo cp $(BIN) /usr/local/bin/$(BIN)
	sudo systemctl restart scorp-agent

deploy-minimal: minimal
	sudo cp $(BIN)-minimal /usr/local/bin/$(BIN)
	sudo systemctl restart scorp-agent

clean:
	rm -f $(BIN) $(BIN)-* ./*.bak

.PHONY: default minimal debug deploy deploy-minimal clean
