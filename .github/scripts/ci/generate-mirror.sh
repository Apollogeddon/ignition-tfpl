#!/bin/bash
# Writes a provider network mirror into webpage/public/providers/, so the site serves
# every release to `tofu init` and `terraform init` through a network_mirror block.
# The mirror lists each published release that has a SHA256SUMS file; its archive
# URLs point at the release's zips, and their hashes come from SHA256SUMS.
# Protocol: https://opentofu.org/docs/internals/provider-network-mirror-protocol/
set -euo pipefail

repo="${GITHUB_REPOSITORY:-apollogeddon/ignition-tofu}"
downloads="${MIRROR_DOWNLOAD_URL:-https://github.com/$repo/releases/download}"
namespace=apollogeddon
type=ignition
# the source address apollogeddon/ignition means registry.opentofu.org to OpenTofu
# and registry.terraform.io to Terraform, so the mirror answers for both
hosts=(registry.opentofu.org registry.terraform.io)

root=$(git rev-parse --show-toplevel)
out="$root/webpage/public/providers"
rm -rf "$out"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

echo '{}' >"$work/versions.json"
for tag in $(gh release list --repo "$repo" --exclude-drafts --limit 1000 --json tagName --jq '.[].tagName'); do
  version="${tag#v}"
  sums="$work/$version.sums"
  if ! gh release download "$tag" --repo "$repo" --pattern "terraform-provider-${type}_${version}_SHA256SUMS" \
    --output "$sums" 2>/dev/null; then
    echo "Skipping $tag: no SHA256SUMS"
    continue
  fi
  # "<sha256>  terraform-provider-ignition_<version>_<os>_<arch>.zip"
  jq -Rn --arg base "$downloads/$tag/" --arg prefix "terraform-provider-${type}_${version}_" '
    { archives: ([inputs | select(test("\\.zip$")) | capture("^(?<hash>[0-9a-f]{64})\\s+(?<file>\\S+)$")
      | { key: (.file | ltrimstr($prefix) | rtrimstr(".zip")),
          value: { url: ($base + .file), hashes: ["zh:" + .hash] } }] | from_entries) }
  ' <"$sums" >"$work/$version.json"
  if [ "$(jq '.archives | length' "$work/$version.json")" -eq 0 ]; then
    echo "Skipping $tag: SHA256SUMS lists no zips"
    continue
  fi
  jq --arg v "$version" '.[$v] = {}' "$work/versions.json" >"$work/versions.tmp" && mv "$work/versions.tmp" "$work/versions.json"
  echo "Mirrored $tag: $(jq -r '.archives | keys | join(", ")' "$work/$version.json")"
done

for host in "${hosts[@]}"; do
  dir="$out/$host/$namespace/$type"
  mkdir -p "$dir"
  jq '{ versions: . }' "$work/versions.json" >"$dir/index.json"
  for version in $(jq -r 'keys[]' "$work/versions.json"); do
    cp "$work/$version.json" "$dir/$version.json"
  done
done
echo "Wrote the mirror for $(jq -r 'keys | join(", ")' "$work/versions.json") to $out"
