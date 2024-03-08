#!/bin/bash

# Define the directory to store the licenses
LICENSE_DIR="./licenses"
mkdir -p $LICENSE_DIR

# Generate DEPENDENCIES.md
DEP_FILE="DEPENDENCIES.md"
echo "# Notices for Third-Party Software" > $DEP_FILE
echo "" >> $DEP_FILE
echo "This software includes and uses the following third-party libraries:" >> $DEP_FILE
echo "" >> $DEP_FILE

# Function to create a filename-safe version of the package name
sanitize_package_name() {
    echo "$1" | sed 's|/|_|g' | sed 's|:|_|g'
}

# Run go-licenses and process the output
/home/ben/go/bin/go-licenses report ./... | while IFS=, read -r package url license; do
    # Skip the main module to avoid treating it as its own dependency
    if [[ "$package" == "github.com/inferret/infer" ]]; then
        continue
    fi

    if [[ $url == http* ]]; then
        # Sanitize the package name to be used in the filename
        safe_package_name=$(sanitize_package_name "$package")

        # Use the full, sanitized package name for the filename
        filename="${safe_package_name}_LICENSE"

        filepath="$LICENSE_DIR/$filename"

        # Download the license file
        curl -sL $url -o $filepath

        # Write to DEPENDENCIES.md
        echo "## $package" >> $DEP_FILE
        echo "" >> $DEP_FILE
        echo "- **Repository:** [Link to License]($url)" >> $DEP_FILE
        echo "- **License:** $license" >> $DEP_FILE
        echo "- **License Text:** [Copy of License Text]($filepath)" >> $DEP_FILE
        echo "" >> $DEP_FILE
    fi
done
echo "DEPENDENCIES.md and licenses have been generated."
