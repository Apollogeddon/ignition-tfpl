#!/bin/bash
set -e

# Ensure base directory exists
mkdir -p webpage/src/content/docs/reference

# Process all markdown files in docs/ recursively
find docs -type f -name "*.md" | while read -r file; do
  # Determine relative path and target destination
  rel_path="${file#docs/}"
  dest_path="webpage/src/content/docs/reference/$rel_path"
  
  # Create target subdirectory if it doesn't exist
  mkdir -p "$(dirname "$dest_path")"
  
  # Special title handling for index.md
  default_title="Reference"
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