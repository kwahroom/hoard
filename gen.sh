#!/bin/bash

# Convert markdown files to HTML using Pando
# Usage: ./gen.sh [src] [out]

# user-specified directory or default to current directory
SRCDIR=${1:-.}
OUTDIR=${2:-.}


# Check if the directory exists
if [ ! -d "$SRCDIR" ]; then
  echo "Directory $DIR does not exist."
  exit 1
fi

# Find all markdown files in the specified directory
find "$SRCDIR" -type f -name "*.md" | while read -r file; do
  # Get the base name of the file (without path)
  base_name=$(basename "$file" .md)
  
  # Convert markdown to HTML using Pando
  pandoc "$file" -o "${OUTDIR}/${base_name}.html"
  
  # Check if the conversion was successful
  if [ $? -eq 0 ]; then
	echo "Converted $file to ${base_name}.html"
  else
	echo "Failed to convert $file"
  fi
done
