#!/bin/bash

# Set the output file
OUTPUT_FILE="all_go_files.txt"

# Create or clear the output file
: > $OUTPUT_FILE

# Perform tree -L 3 and save the output
tree -L 3 > directory_structure.txt

# Find all .go files and concatenate their content into the output file
find . -name "*.go" | while read file; do
    echo "Appending $file to $OUTPUT_FILE"
    echo "File: $file" >> $OUTPUT_FILE
    echo "=====================================" >> $OUTPUT_FILE
    cat "$file" >> $OUTPUT_FILE
    echo -e "\n\n" >> $OUTPUT_FILE
done

echo "Concatenation complete. The output is saved in $OUTPUT_FILE"
