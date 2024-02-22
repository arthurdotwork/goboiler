#!/bin/bash

set -e

function replace_string() {
  local file=$1
  local search=$2
  local replace=$3

  # Replace the string
  if [ "$os" = "Darwin" ]; then
    sed -i '' -e "s#$search#$replace#g" "$file"
  else
    sed -i -e "s#$search#$replace#g" "$file"
  fi
}

# Prompt for GitHub name
echo -n "Enter GitHub name: "
read github

# Check if GitHub name is empty
if [ -z "$github" ]; then
  echo "GitHub name cannot be empty"
  exit 1
fi

# Prompt for project name
echo -n "Enter project name: "
read project

# Check if project name is empty
if [ -z "$project" ]; then
  echo "Project name cannot be empty"
  exit 1
fi

# Find all files containing the string "github.com/arthureichelberger/goboiler"
files=$(grep -rl "github.com/arthureichelberger/goboiler" .)

# Determine OS
os=$(uname)

# Iterate over the files
for file in $files
do
  # Skip if file does not exist
  if [ ! -f "$file" ]; then
    continue
  fi

  # Skip current file
  if [ "$file" = "setup.sh" ]; then
    continue
  fi

  replace_string "$file" "github.com/arthureichelberger/goboiler" "github.com/$github/$project"
done

replace_string ".github/workflows/golang.yaml" "ghcr.io/arthureichelberger/goboiler" "ghcr.io/$github/$project"
replace_string "Taskfile.yaml" "goboiler:local" "$project:local"
replace_string "pkg/prom/prom.go" "goboiler" "$project"

# Remove setup.sh
rm setup.sh

# Remove git history
rm -rf .git

# Reset the README file
echo "# $project" > README.md

# Init the git repository
git init

# Make first commit
git add . && git commit -m "feat: initial commit"
