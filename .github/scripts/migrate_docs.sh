#!/bin/bash
set -e

# Function to transform frontmatter
transform_frontmatter() {
  local src_file="$1"
  local dest_file="$2"
  
  # Extract title and description using awk
  awk '
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
      print "title: \"" (title ? title : "Reference") "\"";
      print "description: \"" desc "\"";
      print "---";
      print_fm = 2;
    }
    print_fm == 2 { print; }
  ' "$src_file" > "$dest_file"
}

# Ensure base directory exists
mkdir -p webpage/src/content/docs/reference

# Process all markdown files in docs/ recursively
find docs -name "*.md" | while read -r file; do
  # Determine relative path and target destination
  rel_path="${file#docs/}"
  dest_path="webpage/src/content/docs/reference/$rel_path"
  
  # Create target subdirectory if it doesn't exist
  mkdir -p "$(dirname "$dest_path")"
  
  echo "Migrating $file to $dest_path"
  transform_frontmatter "$file" "$dest_path"
done
