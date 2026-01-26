---
title: Installation
description: How to install and configure the Ignition Terraform Provider.
---

## Prerequisites

Before using the Ignition Terraform Provider, ensure you have the following installed:

- **Terraform** (v1.0+)
- **Ignition Gateway** (v8.1+) running the Web Developer module (implied by API usage, though core API is usually built-in).

## Provider Configuration

To use the provider, you must configure it in your Terraform files. The provider requires the location of your Ignition Gateway and an API Token for authentication.

### Basic Configuration

Add the following to your `main.tf` or `versions.tf` file:

```hcl
terraform {
  required_providers {
    ignition = {
      source  = "apollogeddon/ignition"
      version = ">= 0.0.1"
    }
  }
}

provider "ignition" {
  host  = "http://localhost:8088"
  token = "YOUR_API_TOKEN_HERE"
  allow_insecure_tls = false # Set to true if using self-signed certs
}
```

### Environment Variables

For security best practices, avoid hardcoding sensitive tokens in your `.tf` files. The provider supports the following environment variables:

| Variable | Description |
| :--- | :--- |
| `IGNITION_HOST` | The base URL of the Ignition Gateway (e.g., `http://10.10.1.5:8088`). |
| `IGNITION_TOKEN` | The API Token generated in the Ignition Gateway Config section. |

When using environment variables, you can keep the provider block empty or minimal:

```hcl
provider "ignition" {}
```

## Generating an API Token

1. Log into your Ignition Gateway Web Interface.
2. Navigate to **Config** > **Security** > **API Tokens**.
3. Click **Create New API Token**.
4. Give it a description (e.g., "Terraform").
5. Copy the generated token immediately; it will not be shown again.

## Verification

To verify the installation, create a simple data source fetch (e.g., reading the default project) and run `terraform init` and `terraform plan`.

```hcl
data "ignition_project" "example" {
  name = "MyProject"
}

output "project_desc" {
  value = data.ignition_project.example.description
}
```
