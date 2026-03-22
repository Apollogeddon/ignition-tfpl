#!/bin/bash
set -e

TF_PLUGIN_DOCS_VERSION="v0.20.1"

echo "Installing tfplugindocs $TF_PLUGIN_DOCS_VERSION..."
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@$TF_PLUGIN_DOCS_VERSION

echo "Generating documentation..."
go generate ./...
