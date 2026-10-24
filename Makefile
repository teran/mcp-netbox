# mcp-netbox Makefile
#
# Standard R07 make interface: CI and local dev call `make <target>` rather than
# raw language commands. Go version is pinned to match go.mod (go 1.27.3).
# All recipes use tabs.

GO ?= go
GOARCH_LINT_VERSION := v1.19.0
COVERAGE_THRESHOLD ?= 95

# --- Image build variables (used by the container-image target) --------------
IMAGE ?= ghcr.io/teran/mcp-netbox
VERSION ?= local
PLATFORMS ?= linux/amd64,linux/arm64
# TAGS may be a space- or comma-separated list of image tags.
TAGS ?= $(IMAGE):$(VERSION)
comma := ,
TAGS_SPACE := $(subst $(comma), ,$(TAGS))

# SNAPSHOT=false produces a real goreleaser release (used by release.yml);
# the default (true) produces a local snapshot build into dist/.
SNAPSHOT ?= true

.PHONY: test lint lint-golangci lint-vuln arch build e2e container-image mutation

## test runs the unit suite with the race detector and enforces the >=95% coverage gate.
test:
	$(GO) test -race -coverprofile=coverage.out -count=1 ./...
	@COVERAGE=$$($(GO) tool cover -func=coverage.out | grep '^total:' | awk '{print $$3}' | sed 's/%//'); \
	if awk "BEGIN {exit !($${COVERAGE} < $(COVERAGE_THRESHOLD))}"; then \
		echo "ERROR: coverage $${COVERAGE}% is below threshold of $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi; \
	echo "OK: coverage $${COVERAGE}% meets threshold of $(COVERAGE_THRESHOLD)%"

## lint runs the full static-analysis suite (golangci + govulncheck + arch).
## Note: the standalone `gosec` step is intentionally omitted — gosec's built-in
## analyzer cannot read Go 1.27 export data (v5 > v4). gosec is still enforced as
## a linter inside golangci-lint (.golangci.yml), so security scanning is retained.
lint: lint-golangci lint-vuln arch

## lint-golangci runs golangci-lint with the repo configuration.
lint-golangci:
	golangci-lint run --config .golangci.yml ./...

## lint-vuln runs the Go vulnerability scanner (govulncheck).
lint-vuln:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...

## arch runs the architecture dependency check.
arch:
	$(GO) run github.com/fe3dback/go-arch-lint@$(GOARCH_LINT_VERSION) check

## build produces a unified dist/ artifact via goreleaser (snapshot by default;
## use `make build SNAPSHOT=false` for a real release).
build:
	goreleaser release --clean $(if $(filter false,$(SNAPSHOT)),,--snapshot)

## e2e runs the end-to-end NetBox integration test (requires Docker and
## MCP_NETBOX_E2E=1; see e2e/netbox_e2e_test.go).
e2e:
	MCP_NETBOX_E2E=1 $(GO) test -tags e2e ./e2e/...

## mutation runs the gremlins mutation-testing hard gate (C02/C08GO). CI binds
## to this target (R07/N31) instead of calling gremlins directly. Thresholds
## (efficacy / mutant-coverage) live in .gremlins.yaml at the repo root, NOT as
## --threshold-* CLI flags: a gremlins/viper bug binds float CLI flags as
## strings, so they are silently ignored by assess(). The config file is parsed
## as real floats and DOES gate the exit code. The positional must be "."
## (not "./...") — gremlins v0.6.0 does not expand "./...".
mutation:
	gremlins unleash .

## container-image builds the Docker image from the unified goreleaser dist/
## artifact. Usage:
##   make container-image                 # build only
##   make container-image push=true       # build and push
## Variables: IMAGE, TAGS (space/comma list), PLATFORMS.
container-image:
	cp dist/mcp-netbox_linux_amd64_v1/mcp-netbox mcp-netbox-linux-amd64
	cp dist/mcp-netbox_linux_arm64_v8.0/mcp-netbox mcp-netbox-linux-arm64
	docker buildx build --platform $(PLATFORMS) $(if $(filter true,$(push)),--push,) $(foreach t,$(TAGS_SPACE),-t $(t)) .
