# Security Policy

## Supported Versions

Only the latest release of the provider receives security fixes.

## Reporting a Vulnerability

Please report security vulnerabilities privately using [GitHub's private vulnerability reporting](https://github.com/Apollogeddon/ignition-tfpl/security/advisories/new) rather than opening a public issue.

You should expect an initial response within a few days. If the issue is confirmed, a fix is released as a patch version and credited in the advisory unless you request otherwise.

## Automated Security Tooling

This repository runs the following on every change:

- **Gitleaks**: scans the repository for committed secrets
- **Trivy** and **govulncheck**: scan the repository and the Go modules for known vulnerabilities
- **OSV-Scanner**: scans the documentation site's dependencies for known vulnerabilities
- **Dependabot**: with a 3-day cooldown before new dependency versions are proposed, giving time for a compromised release to be caught and yanked upstream
