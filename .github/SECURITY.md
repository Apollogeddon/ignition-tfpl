# Security policy

This policy covers the `ignition` provider and its documentation site.

## Supported versions

Only the latest release of the provider receives security fixes.

## Reporting a vulnerability

Report security vulnerabilities privately through [GitHub private vulnerability reporting](https://github.com/Apollogeddon/ignition-tfpl/security/advisories/new). Do not open a public issue.

Expect an initial response within a few days. If the issue is confirmed, the fix is released as a patch version and you are credited in the advisory unless you ask not to be.

## Automated security tooling

The repository's CI runs these checks on pull requests and on pushes to `main`. The Go checks run when the provider's code or tooling changes.

| Tool | What it checks |
| :--- | :--- |
| Gitleaks | Secrets committed to the repository. |
| Trivy | Known high and critical vulnerabilities in the repository's Go modules. |
| govulncheck | Known vulnerabilities in Go code the provider calls. |
| OSV-Scanner | Known vulnerabilities in the Go modules and in the documentation site's npm dependencies. |
| Dependabot | Outdated Go modules, npm packages and GitHub Actions. New versions are proposed after a 3-day cooldown, which gives time for a compromised release to be found and withdrawn upstream. |
