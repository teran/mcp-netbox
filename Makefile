# mcp-netbox Makefile
#
# Go version is pinned to match go.mod (go 1.27.0). All recipes use tabs.

GO ?= go
GOARCH_LINT_VERSION := v1.19.0

.PHONY: test test-race test-e2e build lint arch

## test runs the default unit test suite (no e2e, no race detector).
test:
	$(GO) test ./...

## test-race runs the unit tests with the race detector enabled.
test-race:
	$(GO) test -race ./...

## test-e2e runs the end-to-end NetBox integration test (requires Docker and
## MCP_NETBOX_E2E=1; see e2e/netbox_e2e_test.go).
test-e2e:
	MCP_NETBOX_E2E=1 $(GO) test -tags e2e ./e2e/...

## build compiles all packages.
build:
	$(GO) build ./...

## lint runs golangci-lint with the repo configuration.
lint:
	golangci-lint run ./...

## arch runs the architecture dependency check.
arch:
	$(GO) run github.com/fe3dback/go-arch-lint@$(GOARCH_LINT_VERSION) check
