#!/bin/bash
set -e

# Ensure base directory exists (updated to sources)
mkdir -p webpage/src/content/docs/sources

# Process all markdown files in docs/ recursively
find docs -type f -name "*.md" | while read -r file; do
  # Determine relative path components
  rel_path="${file#docs/}"
  dir_part=$(dirname "$rel_path")
  base_part=$(basename "$rel_path")
  
  # Prefix with 'ignition_' if not index.md and not already prefixed
  target_name="$base_part"
  if [[ "$base_part" != "index.md" && "$base_part" != ignition_* ]]; then
    target_name="ignition_$base_part"
  fi
  
  # Determine final destination path in the sources/ directory
  if [[ "$dir_part" == "." ]]; then
    dest_path="webpage/src/content/docs/sources/$target_name"
  else
    dest_path="webpage/src/content/docs/sources/$dir_part/$target_name"
  fi
  
  # Create target subdirectory if it doesn't exist
  mkdir -p "$(dirname "$dest_path")"
  
  # Special title handling for index.md
  default_title="Sources"
  if [[ "$rel_path" == "index.md" ]]; then
    default_title="Ignition Provider"
  fi
  
  echo "Migrating $file to $dest_path"
  
  # Extract title and description using awk
  awk -v def_title="$default_title" ' 
    BEGIN { in_fm = 0; print_fm = 0; title = ""; desc = ""; }
    /^---$/ {
      if (in_fm == 0) { in_fm = 1; next; }
      if (in_fm == 1) { in_fm = 0; print_fm = 1; next; }
    }
    in_fm == 1 {
      if ($0 ~ /^page_title:/) {
        gsub(/^page_title: "/, "", $0);
        gsub(/"$/, "", $0);
        gsub(/ - ignition$/, "", $0);
        title = $0;
      }
      if ($0 ~ /^description: |-/) {
        next;
      }
      if ($0 ~ /^  /) {
        gsub(/^  /, "", $0);
        desc = (desc ? desc " " : "") $0;
      }
    }
    print_fm == 1 {
      print "---";
      print "title: \"" (title ? title : def_title) "\"";
      print "description: \"" desc "\"";
      print "---";
      print_fm = 2;
    }
    print_fm == 2 { print; }
  ' "$file" > "$dest_path"
done
