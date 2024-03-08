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

# Run go-licenses and process the output
/home/ben/go/bin/go-licenses report ./... | while IFS=, read -r package url license; do
    # Skip the main module to avoid treating it as its own dependency
    if [[ $line == *"module github.com/inferret/infer"* ]]; then
        continue
    fi

    if [[ $url == http* ]]; then
        org=$(echo $package | cut -d '/' -f 2)
        repo=$(echo $package | cut -d '/' -f 3)
        version=$(echo $url | rev | cut -d '/' -f 2 | rev)

        # Correcting version naming and paths for special cases
        if [ "$version" == "blob" ] || [ "$version" == "HEAD" ]; then
            version="HEAD"
        fi

        if [[ $package == golang.org/x/* ]]; then
            org="golang"
            repo=$(echo $package | cut -d '/' -f 4)
        fi

        filename="${org}_${repo}_${version}_LICENSE"
        filepath="$LICENSE_DIR/$filename"

        # Download the license file
        curl -sL $url -o $filepath

        # Write to DEPENDENCIES.md
        echo "## $package" >> $DEP_FILE
        echo "" >> $DEP_FILE
        echo "- **License:** $license" >> $DEP_FILE
        echo "- **Version:** $version" >> $DEP_FILE
        echo "- **License Text:** [Link to License]($filepath)" >> $DEP_FILE
        echo "" >> $DEP_FILE
    fi
done
echo "DEPENDENCIES.md and licenses have been generated."
