# Security Policy

This MCP server is developed and tested largely by AI agents with minimal human
review, so security is enforced automatically by CI (lint, coverage, mutation,
gosec, govulncheck, gitleaks). The security baseline is described in `SPEC.md`
(S01–S14) — the single source of truth for the security contracts this
repository must honour.

## Reporting a Vulnerability

Please **do not** open a public issue for a security vulnerability if you
believe it is a high-severity, actively exploitable issue.

Report vulnerabilities to the maintainers via the
[project issue tracker](https://github.com/teran/mcp-netbox/issues), marking the
report as security-related. If the report must stay private, open an issue and
request a private channel before including any sensitive details.

Please include:

- the affected version(s) and how you detected the issue;
- a minimal, reproducible description of the vulnerability;
- the potential impact and, if known, a suggested fix;
- the relevant configuration context (transport, environment variables) in
  which it occurs.

## What happens next

1. The maintainers acknowledge the report within 5 business days.
2. The issue is triaged: severity, affected versions, and a fix plan.
3. A fix is developed **test-first** and lands through the normal CI gates.
4. Once a fix is available, the vulnerability is disclosed and (if appropriate)
   a `v*` release is published.

We ask that you **do not** disclose the vulnerability publicly until a fix has
been released or the maintainers explicitly agree otherwise (responsible
disclosure).

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 0.x.x   | ✅ |

Only the latest `0.x` release line is supported. Older versions are not
patched; upgrade to a supported version.

## Security Baseline

- **TLS** is never implemented inside the server — it is terminated at the
  reverse proxy.
- **Secrets** (e.g. the NetBox API token) are never leaked into tool output or
  logs — annotation-based redaction (`secret:true`), token `String`/`GoString`/
  `MarshalJSON` redaction, outbound-log path-only URL recording (never the
  authority/userinfo), and auth-header exclusion.
- **Stateless** — the server holds no state and persists nothing; the token is
  never stored server-side (per-request in HTTP, `NETBOX_TOKEN` env in STDIO).
- **NetBox target address** — `NETBOX_URL` is set by the operator from the
  environment and never takes a value from an incoming MCP request, so the
  classic client-side SSRF vector does not apply. The network and addresses
  that the server may reach are the operator's responsibility; the server does
  not restrict destination addresses (private/self-hosted NetBox is supported).
- **Input/output hygiene** — CRLF/control characters are stripped from write
  payloads and ANSI/C0 escapes are stripped from tool output.
- **Destructive tools** (`delete_*`) require human confirmation (HITL) and are
  flagged `destructiveHint: true`.

See `SPEC.md` for the full domain-specific security audit (S14).
