<h1 align="center">packages.go</h1>


`packages` is a CLI tool written in Go, designed to assist in managing and categorizing packages for Arch Linux installations. Instead of directly installing software, this tool acts as a registry manager. It verifies if a package exists in the official **Pacman** repositories or the **AUR**, and then categorizes it into predefined installation profiles: `server`, `minimal`, or `desktop`.

This categorized list is saved to a local registry file named `drxboot.packages`. This file is formatted as a valid Bash script containing arrays of packages, which is subsequently read by my automated bootstrap script (`drxboot.sh`) to perform fresh, unattended Arch Linux installations.

## Configuration

By default, the tool looks for or creates the `drxboot.packages` file in your home directory (`~/drxboot.packages`).

If you want to store this file in a custom location (for instance, inside a dotfiles repository so you can track it with Git), you can set the `PACKAGES_PATH` environment variable in your shell configuration (`~/.bashrc` or `~/.zshrc`):

```bash
# Set a custom directory for the packages registry file
export PACKAGES_PATH="$HOME/Workspace/dotfiles"
```

## Registry Format

The generated `drxboot.packages` file is formatted as valid Bash code. This allows deployment scripts to simply `source` the file and iterate over the arrays. The structure looks like this:

```bash
#!/bin/bash

server_pacman_packages=(
	base
	git
	neovim
)

desktop_aur_packages=(
	polybar
	spotify-bin
)
```

## Commands & Usage

Below is a detailed breakdown of all available commands, what they do, and how to use them.

### `add`

Adds a new package to your registry. When you run this command, the tool will silently query both `pacman` and the AUR to verify the package actually exists. If it does, it will interactively prompt you to choose which category (`server`, `minimal`, or `desktop`) it belongs to, and finally save it to the list.

**Examples:**

```bash
# Verifies if 'neovim' exists and prompts you to select a category for it
packages add neovim

# Tries to add a package. If 'fake-package-123' doesn't exist in pacman 
# or the AUR, the tool will throw an error and abort the operation.
packages add fake-package-123
```

### `remove`

Searches for the specified package across all your lists (all categories and repositories) and completely removes it from the `drxboot.packages` registry file.

**Examples:**

```bash
# Finds 'firefox' in your saved lists and removes it entirely
packages remove firefox

# Removes a package that was saved with a specific suffix ('-git' or '-bin')
# This package was actually saved as `spotify-bin`, detects it and removes it
packages remove spotify
```

### `search`

Quickly checks your local registry to see if a package is already saved. If it finds the package, it will tell you exactly which category and repository it belongs to.

**Examples:**

```bash
# Checks if 'htop' is registered.
# Output example: [SUCCESS] Package 'htop' found in list: MINIMAL (PACMAN)
packages search htop
```

### `list`

Displays the contents of your `drxboot.packages` registry in a clean, column-formatted layout. You can run it without arguments to see absolutely everything, or pass an optional filter to narrow down the results.

**Available Filters:** `server`, `minimal`, `desktop`, `pacman`, `aur`.

**Examples:**

```bash
# Displays all packages saved across all categories and repositories
packages list

# Displays ONLY the packages categorized under 'desktop'
packages list desktop

# Displays ONLY the packages that belong to the 'aur' repository across all categories
packages list aur
```