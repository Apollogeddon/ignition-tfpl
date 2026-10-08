# GitHub Workflows Documentation

This repository uses GitHub Actions to automate testing, quality assurance, documentation deployment, and the release process for the Ignition Terraform Provider.

## 🏗️ Orchestration: The Index Workflow

The [`.index.yaml`](./workflows/.index.yaml) workflow is the primary entry point for changes to the `main` branch. It orchestrates the execution of other workflows in a specific order:

1. **Testing & Quality**: Runs the unit testing and quality suites in parallel.
2. **Ignition Acceptance Tests**: Launches a real Ignition Gateway via Docker Compose and runs Terraform acceptance tests against it.
3. **Release**: Triggered only after all previous checks pass.

The documentation site has its own workflow (below), which runs on every push and pull request.

---

## 🔍 Quality Assurance

The [`quality.yaml`](./workflows/quality.yaml) workflow focuses on static analysis and security scanning:

- **GolangCI-Lint**: Runs a suite of Go linters to ensure code consistency and catch common errors.
- **Trivy**: Scans the repository filesystem for known vulnerabilities and configuration issues.
- **Govulncheck**: Scans the Go dependencies for known vulnerabilities using the Go vulnerability database.

## 🧪 Testing

The provider uses a two-tier testing strategy:

### Unit Testing

The [`testing.yaml`](./workflows/testing.yaml) workflow runs standard Go unit tests with the race detector enabled to ensure internal logic is sound and thread-safe.

### Acceptance Testing

The [`ignition.yaml`](./workflows/ignition.yaml) workflow performs "real-world" validation:

- **Environment**: Spins up an Ignition 8.3 Gateway using `docker-compose.yml`.
- **Initialization**: Waits for the Gateway to be healthy and accessible.
- **Execution**: Runs `go test -v ./internal/provider/...` with `TF_ACC=1` to execute the full Terraform resource lifecycle (Create, Read, Update, Delete) against the live API.

## 🚀 Release Process

The [`release.yaml`](./workflows/release.yaml) workflow handles versioning and distribution:

- **Release Please**: Automatically manages version bumps and `CHANGELOG.md` updates based on conventional commits.
- **GoReleaser**: Packages the provider for multiple platforms, signs the binaries with GPG, and publishes them to GitHub Releases.

## 📖 Documentation

The [`webpage.yaml`](./workflows/webpage.yaml) workflow manages the [Astro](https://astro.build/)-based documentation site:

- **Quality**: Calls [forgejs](https://github.com/apollogeddon/forgejs)'s `quality.yml`: Gitleaks over the whole repository, OSV-Scanner on the site's dependencies, Biome and the type check.
- **Markdown**: Lints every Markdown file in the repository with `markdownlint-cli2`.
- **Build**: Uses `tfplugindocs` to generate the provider docs from the schema and examples, copies them into the site with `migrate-docs.sh`, and builds the static site located in the `webpage/` directory, on pull requests too, so a broken site fails the pull request.
- **Deploy**: On `main`, publishes the build artifacts to **GitHub Pages**.

---

## 🛠️ Configuration & Maintenance

- **`release.json`**: Configures `release-please` behavior.
- **`dependabot.yml`**: Automatically keeps GitHub Actions, Go modules, and NPM dependencies up to date.
