# E&E OS — Omarchy Fork + Dotfiles

**Repo:** `github.com/eande-uk/eande-os`

## Structure

```
erch/         Git submodule — forked omarchy (github.com/eande-uk/erch)
dotfiles/     Portable Layer 3 config overlays (stow-deployed)
scripts/      Deployment tooling
tests/        Go verification tests
layer-zero/   System state sync (allowlist-based)
```

## Quick Start

```bash
# Clone with submodule
git clone git@github.com:eande-uk/eande-os.git
cd eande-os
git submodule update --init erch/

# Deploy erch fork to ~/.local/share/omarchy/
erch/setup.sh

# Deploy dotfiles to $HOME
make deploy

# Reload Hyprland
hyprctl reload
```

## For non-fork systems

Clone without the submodule and just run `make deploy` for dotfiles only.

## Make Targets

| Target | Description |
|--------|-------------|
| `make deploy` | Stow dotfiles to $HOME |
| `make test` | Run Go tests |
| `make commit` | Stage + commit |
| `make pr` | Open PR to master |
