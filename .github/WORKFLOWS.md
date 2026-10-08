# GitHub Workflows Documentation

This repository uses GitHub Actions to automate testing, quality assurance, documentation deployment, and the release process for the Ignition Terraform Provider.

## 🏗️ Orchestration: The Index Workflow

The [`.index.yaml`](./workflows/.index.yaml) workflow runs on every pull request and on pushes to `main` that change the provider's code, its tooling or the workflows that test and release it. It runs the checks in this order:

1. **Changes**: On a pull request, works out whether the provider, its tooling or those workflows changed, and whether any shell script changed.
2. **shellcheck**: Lints every shell script, when a script changed.
3. **Testing & Quality**: forgego's reusable `testing.yml` and Trivy, in parallel, when the provider changed.
4. **Ignition Acceptance Tests**: Launches a real Ignition Gateway via Docker Compose and runs Terraform acceptance tests against it.
5. **Release**: Only on `main`, after all previous checks pass.
6. **Webpage**: On a pull request, the documentation site's checks and build ([below](#-documentation)).
7. **Auto-merge**: Merges a Dependabot pull request through [forgejs](https://github.com/apollogeddon/forgejs)'s `merge.yml`, once every job above has passed or been skipped.

A new push to a pull request cancels its previous run. On `main`, the documentation site runs on its own push trigger.

---

## 🔍 Quality and Testing

The provider's tooling comes from [forgego](https://github.com/apollogeddon/forgego): each tool is pinned in its own module under `.forgego/`, and `Taskfile.yml` runs them (`task lint`, `task test`, `task build`). Run `task hooks` once to install the git hooks: they format staged Go files and check that forgego's files are current before a commit, lint before a push, and check commit messages against Conventional Commits.

forgego's [`testing.yml`](https://github.com/apollogeddon/forgego/blob/main/.github/workflows/testing.yml) runs:

- **Quality**: `go mod tidy -diff`, the golangci-lint format check and lint (forgego's base config merged with `.golangci.local.yml` into `.golangci.yml`), a check that forgego's files are current, govulncheck and OSV-Scanner.
- **Unit tests**: `task test`, with the race detector and coverage.
- **Build**: `task build`.
- **Patch**: on `main`, upgrades modules that govulncheck reports as vulnerable and commits the fix.

**Trivy** scans the repository for known vulnerabilities and fails on high or critical ones.

## 🧪 Acceptance Testing

The [`ignition.yaml`](./workflows/ignition.yaml) workflow performs "real-world" validation:

- **Environment**: Spins up an Ignition 8.3 Gateway using `docker-compose.yml`.
- **Initialization**: Waits for the Gateway to be healthy and accessible.
- **Execution**: Runs `go test -v ./internal/provider/...` with `TF_ACC=1` to execute the full Terraform resource lifecycle (Create, Read, Update, Delete) against the live API.

## 🚀 Release Process

The [`release.yaml`](./workflows/release.yaml) workflow handles versioning and distribution:

- **Release Please**: Automatically manages version bumps and `CHANGELOG.md` updates based on conventional commits.
- **GoReleaser**: Builds the provider in the layout the Terraform Registry expects (`.goreleaser.yaml`): a zip per platform, a `SHA256SUMS` file signed with the GPG key, and the registry manifest (`terraform-registry-manifest.json`).

release-please creates each release as a draft, and GoReleaser attaches the files and then publishes it, so the release flow also works with immutable releases. Try the build locally with `task release:snapshot`, which skips the signing.

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
