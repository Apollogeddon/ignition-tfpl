---
title: Provider Capabilities
description: Overview of supported Ignition resources and features.
---

The Ignition Terraform Provider supports a comprehensive set of configuration resources, allowing for Infrastructure-as-Code management of an Ignition Gateway.

## Supported Resources

### Core System

| Resource | Description |
| :--- | :--- |
| `ignition_project` | Manage Ignition Projects (Vision/Perspective/Perspective Sessions). |
| `ignition_database_connection` | Configure connections to SQL databases (MariaDB, MySQL, PostgreSQL, MSSQL, Oracle). |
| `ignition_tag_provider` | Manage Realtime Tag Providers (Standard). |
| `ignition_user_source` | Configure Internal, Database, or Active Directory user sources. |
| `ignition_identity_provider` | Setup IdPs including Internal, OpenID Connect (OIDC), and SAML 2.0. |

### Connectivity & Devices

| Resource | Description |
| :--- | :--- |
| `ignition_opc_ua_connection` | Manage outgoing OPC UA Client connections. |
| `ignition_device` | Configure OPC UA Devices (Modbus, Siemens, Simulators, etc.). |
| `ignition_gan_outgoing` | Configure Gateway Network connections to other Gateways. |

### Gateway Settings

| Resource | Description |
| :--- | :--- |
| `ignition_redundancy` | **Singleton**. Configure Master/Backup redundancy roles and sync settings. |
| `ignition_gan_settings` | **Singleton**. General Gateway Network settings (SSL requirements, proxy hops). |
| `ignition_smtp_profile` | Configure Email/SMTP profiles for alarm notifications and reporting. |

### Alarming & Auditing

| Resource | Description |
| :--- | :--- |
| `ignition_alarm_journal` | Configure storage for Alarm History (Database or Remote). |
| `ignition_audit_profile` | Configure Audit Logs (Database or Internal). |
| `ignition_alarm_notification_profile` | Configure notification pipelines (Email). |

### Data Storage

| Resource | Description |
| :--- | :--- |
| `ignition_store_forward` | Configure Store-and-Forward engines to buffer data during database outages. |

## Data Sources

The provider includes **Data Sources** for most of the resources listed above. This allows you to reference existing configuration on a Gateway that was not created by Terraform.

**Example:**

```hcl
data "ignition_project" "global" {
  name = "global"
}

resource "ignition_project" "site_a" {
  name   = "site_a"
  parent = data.ignition_project.global.name
}
```

## Feature Highlights

- **Polymorphism**: Resources like `ignition_device` or `ignition_user_source` automatically adapt their validation and available fields based on the `type` selected.
- **Secure Configuration**: Built-in support for Ignition's encryption endpoints ensures passwords and secrets are handled securely during transmission.
- **Drift Detection**: Full support for `terraform plan` to detect manual changes made in the Ignition Designer or Web Config interface.
