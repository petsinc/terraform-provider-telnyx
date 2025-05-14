#!/bin/bash

# Note: Must be run through direnv source

if ! has asdf; then
  echo "asdf not found. Please install asdf first: https://asdf-vm.com/guide/getting-started.html"
  exit 1
fi

# Install asdf plugins
plugin_list=$(asdf plugin list)
plugins=("golang" "terraform" "nodejs" "pnpm")

for plugin in "${plugins[@]}"; do
  if ! echo $plugin_list | grep -q " $plugin "; then
    asdf plugin add $plugin
    echo "Added asdf plugin $plugin"
  fi
done

# Install necessary versions
asdf install
