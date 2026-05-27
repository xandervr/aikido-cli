.PHONY: build install install-system install-completions uninstall test fmt vet tidy clean

BIN := bin/aikido
PKG := github.com/xandervr/aikido-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X $(PKG)/internal/version.Version=$(VERSION) -X $(PKG)/internal/version.Commit=$(COMMIT) -X $(PKG)/internal/version.Date=$(DATE)

# Where to install the binary. Defaults to Go's bin dir (`go env GOBIN`, or
# `$(go env GOPATH)/bin` if GOBIN is empty) so the install lands wherever Go
# already puts binaries — typically on PATH and, for asdf/goenv-managed Go,
# behind the shim that resolves `aikido`. Override with
# `make install INSTALL_DIR=/some/path`.
GOBIN := $(shell go env GOBIN)
ifeq ($(strip $(GOBIN)),)
GOBIN := $(shell go env GOPATH)/bin
endif
INSTALL_DIR ?= $(GOBIN)

# Where shell completions land. Override per-shell if you keep them elsewhere.
ZSH_COMP_DIR ?= $(HOME)/.zsh/completions
BASH_COMP_DIR ?= $(HOME)/.local/share/bash-completion/completions
FISH_COMP_DIR ?= $(HOME)/.config/fish/completions

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/aikido

install: build
	@mkdir -p $(INSTALL_DIR)
	@install -m 0755 $(BIN) $(INSTALL_DIR)/aikido
	@echo "✓ Installed: $(INSTALL_DIR)/aikido ($$($(INSTALL_DIR)/aikido --version))"
	@case ":$$PATH:" in \
	  *":$(INSTALL_DIR):"*) echo "✓ $(INSTALL_DIR) is on PATH" ;; \
	  *) printf '\n  ⚠ %s is NOT on your PATH.\n  Add this to ~/.zshrc (or ~/.bashrc):\n    export PATH="%s:$$PATH"\n\n' "$(INSTALL_DIR)" "$(INSTALL_DIR)" ;; \
	esac
	@RESOLVED=$$(command -v aikido 2>/dev/null); \
	if [ -n "$$RESOLVED" ] && [ "$$RESOLVED" != "$(INSTALL_DIR)/aikido" ]; then \
	  printf '\n  ⚠ Another aikido is earlier on PATH: %s\n    It shadows the fresh install. Either remove it or move %s ahead in PATH.\n\n' "$$RESOLVED" "$(INSTALL_DIR)"; \
	fi

install-system: build
	sudo install -m 0755 $(BIN) /usr/local/bin/aikido
	@echo "✓ Installed: /usr/local/bin/aikido"

install-completions: build
	@AIKIDO=$$(command -v aikido 2>/dev/null || echo ./$(BIN)); \
	echo "Using binary: $$AIKIDO"; \
	mkdir -p $(ZSH_COMP_DIR) $(BASH_COMP_DIR) $(FISH_COMP_DIR); \
	$$AIKIDO completion zsh  > $(ZSH_COMP_DIR)/_aikido; \
	$$AIKIDO completion bash > $(BASH_COMP_DIR)/aikido; \
	$$AIKIDO completion fish > $(FISH_COMP_DIR)/aikido.fish; \
	echo "✓ zsh:  $(ZSH_COMP_DIR)/_aikido"; \
	echo "✓ bash: $(BASH_COMP_DIR)/aikido"; \
	echo "✓ fish: $(FISH_COMP_DIR)/aikido.fish"; \
	echo ""; \
	echo "  zsh: ensure ~/.zshrc contains BEFORE 'compinit':"; \
	echo "    fpath=($(ZSH_COMP_DIR) \$$fpath)"; \
	echo "    autoload -Uz compinit && compinit"; \
	echo ""; \
	echo "  bash: requires bash-completion (e.g. 'brew install bash-completion@2')"; \
	echo "  fish: completions load automatically on next shell start."

uninstall:
	@rm -f $(INSTALL_DIR)/aikido
	@rm -f $(ZSH_COMP_DIR)/_aikido
	@rm -f $(BASH_COMP_DIR)/aikido
	@rm -f $(FISH_COMP_DIR)/aikido.fish
	@echo "✓ Uninstalled (binary + completions from $(INSTALL_DIR), $(ZSH_COMP_DIR), $(BASH_COMP_DIR), $(FISH_COMP_DIR))"
	@echo "  System-wide install at /usr/local/bin/aikido is NOT removed by this target."

test:
	go test ./... -count=1 -race -cover

fmt:
	gofmt -s -w .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/
