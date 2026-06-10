# Target: Upstream Omarchy

This guide covers applying E&E OS configs to a stock omarchy installation.

## Overview

On upstream omarchy, the E&E OS repo acts as a config overlay:

1. **layer-zero/** manages packages per your selected profile
2. **dotfiles/** provides config overrides via stow

No erch fork needed. The same configs that erch ships natively are deployed
as overlays on top of omarchy.

## Prerequisites

- Stock omarchy installation (any version)
- `git`, `stow` installed
- `gum` installed (for interactive scripts)

## Quick Start

```bash
# Clone the repo
git clone git@github.com:eande-uk/eande-os.git
cd eande-os

# (Optional) Select a profile for layer-zero
# Edit layer-zero/config.sh to enable profile categories

# Run layer-zero sync
make layer-zero  # Interactive: choose what to install/remove

# Deploy dotfiles
make deploy

# Reload Hyprland
hyprctl reload
```

## Profile Selection

layer-zero on upstream omarchy works the same as on erch:

1. Enable profile categories in `layer-zero/config.sh`
2. Add profile packages to `layer-zero/allowlist.txt`
3. Run `make layer-zero/apply` to sync

Example — enable Dev tooling:

```bash
# In layer-zero/config.sh:
PRUNE_PACMAN_CATEGORIES=(runtimes)

# In layer-zero/allowlist.txt, add:
# docker, docker-compose, python, nodejs, go, rust, ...
```

## Deployment Steps

### 1. Clone and Init

```bash
git clone git@github.com:eande-uk/eande-os.git
cd eande-os
# No submodule needed — erch is not used on upstream omarchy
```

### 2. Layer Zero

```bash
# Preview what layer-zero would change
make layer-zero/dry-run

# Run interactive sync
make layer-zero

# Or apply directly
make layer-zero/apply
```

### 3. Deploy Configs

```bash
# Preview
make deploy/dry-run

# Deploy
make deploy
```

### 4. Post-Deploy

```bash
# Restart components to pick up new configs
hyprctl reload
omarchy restart waybar
omarchy restart mako
omarchy restart walker
omarchy restart terminal

# Verify deployment
make status
```

## Verification

```bash
# Check symlinks point to repo
readlink -f ~/.config/hypr/bindings.conf
# → .../eande-os/dotfiles/home/.config/hypr/bindings.conf

# Run tests
make test/quiet
```

## Updating

```bash
git pull
make deploy          # Re-apply configs
make layer-zero      # Re-sync packages
```

## Uninstalling

```bash
# Remove symlinks
cd eande-os/dotfiles
stow --delete -t $HOME home/

# Config files remain in ~/.config/ (now unmanaged)
# Remove them manually or restore from omarchy defaults
```
