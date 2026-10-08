<br />
<div align="center">
  <a href="https://apollogeddon.github.io/ignition-tfpl">
    <img src="webpage/public/favicon.png" alt="Logo" width="100" height="100">
  </a>
  <h3 align="center">Ignition Terraform Provider</h3>
  <p align="center">
    Manage your Inductive Automation Ignition Gateway infrastructure as code.
    <br />
    <a href="https://apollogeddon.github.io/ignition-tfpl"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/apollogeddon/ignition-tfpl/issues">Report Bug</a>
    ·
    <a href="https://github.com/apollogeddon/ignition-tfpl/issues">Request Feature</a>
  </p>
</div>

## 🚀 Overview

The **Ignition Terraform Provider** allows you to manage Inductive Automation's Ignition Perspective 8.3 infrastructure with [OpenTofu](https://opentofu.org/) or Terraform. Configure Projects, Database Connections, Tag Providers, and Enterprise settings (Redundancy, GAN) alongside your cloud infrastructure.

## ✨ Features

- **Infrastructure as Code**: Version control your Gateway configuration. Manage Projects, Database Connections, and Tag Providers.
- **Secure by Design**: Sensitive credentials are encrypted in-flight using Ignition's native encryption endpoints before being stored.
- **Drift Detection**: Automatically detect and reconcile manual changes made in the Designer or Gateway Web Interface.
- **Enterprise Ready**: Support for complex architectures including Redundancy, Gateway Networks, and Identity Providers (SAML/OIDC).

## 📦 Installation

### Prerequisites

- **OpenTofu** (v1.6+) or **Terraform** (v1.0+). The provider is tested with OpenTofu.
- **Ignition Gateway** (v8.3+)

### Configuration

Add the provider to your configuration:

```hcl
terraform {
  required_providers {
    ignition = {
      source  = "apollogeddon/ignition"
      version = ">= 0.0.1"
    }
  }
}
```

| Variable | Description |
| :--- | :--- |
| `IGNITION_HOST` | The base URL of the Ignition Gateway (e.g., `http://10.10.1.5:8088`). |
| `IGNITION_TOKEN` | The API Token generated in the Ignition Gateway Config section. |

## 🛠️ Usage

### Quick Start

Configure your provider and manage a project in seconds:

```hcl
provider "ignition" {
  host  = "http://localhost:8088"
  token = var.ignition_token
}

resource "ignition_project" "example" {
  name        = "MyEnterpriseProject"
  title       = "Enterprise Dashboard"
  description = "Managed via OpenTofu"
  enabled     = true
}
```

## 🧩 Supported Resources

The provider supports a comprehensive set of Ignition resources:

- **Core System**: Projects, Database Connections, Tag Providers, User Sources, Identity Providers.
- **Connectivity**: OPC UA Connections, Devices, Gateway Network (GAN).
- **Settings**: Redundancy, SMTP Profiles, Alarm Journals, Audit Profiles.
- **Data Storage**: Store-and-Forward engines.

See the [Documentation](https://apollogeddon.github.io/ignition-tfpl) for the full list and detailed usage.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

Run `task hooks` once per clone to install the git hooks: commit messages follow [Conventional Commits](https://www.conventionalcommits.org/), which release-please uses for versioning. `task test` runs the unit tests and `task test:acc` the acceptance tests against a gateway started with `docker compose up`; both download the pinned OpenTofu release into `.bin/` (`.github/scripts/install-tofu.sh`) and run against it. To preview the reference docs locally, generate them with `.github/scripts/ci/generate-docs.sh` and `migrate-docs.sh`; git ignores the output.

## 📄 License

Released under the [MIT License](LICENSE).
