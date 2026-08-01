---
page_title: "Bootstrapping a Gateway with Docker"
subcategory: ""
description: |-
  Provisioning and configuring an Ignition Gateway from a single Terraform apply.
---

# Bootstrapping a Gateway with Docker

The `ignition` provider manages resources on a Gateway that is already
running — it has no way to install Ignition itself, since there is no REST
API to call until the Gateway exists. To get a fully bootstrapped Gateway
(built from a seed backup, restored on first boot, and configured) from a
single `terraform apply`, pair this provider with
[`kreuzwerker/docker`](https://registry.terraform.io/providers/kreuzwerker/docker/latest)
in the same root module.

This mirrors the pattern this repository's own `docker-compose.yml` uses for
acceptance testing: bake a seed `.gwbk` backup into a derived image, restore
it against the running container, then restart to apply the restore.

## The pattern

```terraform
terraform {
  required_providers {
    docker   = { source = "kreuzwerker/docker" }
    ignition = { source = "registry.terraform.io/apollogeddon/ignition" }
  }
}

resource "docker_image" "gateway" {
  name = "my-ignition-gateway"
  build {
    context    = path.module
    dockerfile = "Dockerfile" # FROM inductiveautomation/ignition:...  +  COPY seed.gwbk /restore.gwbk
  }
}

resource "docker_container" "gateway" {
  name  = "ignition-gateway"
  image = docker_image.gateway.image_id

  env = [
    "ACCEPT_IGNITION_EULA=Y",
    "GATEWAY_ADMIN_USERNAME=admin",
    "GATEWAY_ADMIN_PASSWORD=password123",
    "IGNITION_EDITION=standard", # see "Three things to know" below
  ]

  ports {
    internal = 8088
    external = 8088 # fixed — see "Three things to know" below
  }

  healthcheck {
    test     = ["CMD", "curl", "-f", "http://localhost:8088/system/gwinfo"]
    interval = "10s"
    timeout  = "5s"
    retries  = 12
  }

  wait         = true
  wait_timeout = 300

  # Restore live against the running Gateway, then restart to apply it —
  # see "Three things to know" for why this is preferred over the simpler
  # `-r /restore.gwbk` boot flag.
  provisioner "local-exec" {
    command = <<-EOT
      set -e
      docker exec ${self.name} /usr/local/bin/ignition/gwcmd.sh -s /restore.gwbk -y
      docker exec ${self.name} /usr/local/bin/ignition/gwcmd.sh -r
      timeout 120 bash -c 'until curl -sf http://localhost:8088/system/gwinfo >/dev/null; do sleep 3; done'
    EOT
  }
}

provider "ignition" {
  host = "http://localhost:8088"
}

resource "ignition_tag_provider" "example" {
  name       = "MyTags"
  type       = "STANDARD"
  depends_on = [docker_container.gateway]
}
```

A complete, runnable version of this — tested end-to-end with a real
`terraform apply` / `terraform destroy` cycle — is in
[`examples/bootstrap-docker`](https://github.com/apollogeddon/ignition-tfpl/tree/main/examples/bootstrap-docker).

## Three things to know

**Restore live, not via the `-r` boot flag — at least not without checking
first.** Ignition's Docker image supports passing `-r /restore.gwbk` as the
container's `command` to restore on first boot, which is simpler than the
provisioner above. It's worth trying first. But for a backup containing
custom API keys with non-default Security Level grants, that path can
silently fail to carry the grants over: the restored token still
authenticates, but every API call 403s, because the boot-time restore runs
before the Gateway's own default security initialization has fully settled.
Restoring live against an already-running Gateway (as above) doesn't have
this problem — it's the same sequence `docker-compose.yml` uses for this
repository's own acceptance tests, verified across the full test suite.
If your seed backup only has default security configuration, the simpler
`-r` flag is worth trying — just verify a real API call succeeds with your
token before relying on it.

**The container's external port is fixed, not random.** A `provider` block's
configuration can't reference another resource's computed attributes within
the same apply — Terraform needs to know how to reach the provider's API
before it can plan resources that use it. Pinning `external = 8088` lets the
`ignition` provider's `host` be a plain string instead of a reference,
sidestepping that limitation entirely.

**Every `ignition_*` resource needs `depends_on = [docker_container.gateway]`.**
`provider` blocks don't support `depends_on`, and because `host` above is a
static string rather than a reference to the container, Terraform has no
other way to infer that the Gateway must exist first — without it, Terraform
may try to create the container and configure the Gateway in parallel. For a
real configuration with many resources, put them in a child module and set
`depends_on` on the `module` call once instead of repeating it everywhere:

```terraform
module "gateway_config" {
  source     = "./config"
  depends_on = [docker_container.gateway]
}
```

`docker_container`'s own `wait = true` blocks *its* apply until the
healthcheck passes, and the `local-exec` provisioner above blocks until the
post-restore restart finishes — so by the time anything depending on this
resource runs, the Gateway is genuinely ready to accept API calls. The
provider's built-in retry logic (`client.NewClient`, up to 10 retries) is
then just a safety margin, not the primary mechanism keeping this reliable.

## Non-Docker targets

The same shape applies to a VM or bare-metal install: use a `null_resource`
with a `remote-exec` or `local-exec` provisioner (or cloud-init/user-data)
to run the platform installer and an equivalent restore step, then chain
into the `ignition` provider the same way — a static, known-in-advance host
plus explicit `depends_on`. That path is highly specific to your cloud and
OS, so it isn't included as a worked example here.
