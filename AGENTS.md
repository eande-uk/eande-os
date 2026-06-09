# E&E OS — Agent Guide

## Repo Structure

```
erch/         Git submodule — forked omarchy (github.com/eande-uk/erch)
dotfiles/     Portable Layer 3 config overlays (stow-deployed)
scripts/      Deploy tooling
tests/        Go verification tests
layer-zero/   System state sync
```

## Key Rules

- **`erch/` is a submodule** — edit the fork at `github.com/eande-uk/erch`, not the submodule directly
- **Edit `~/.config/<app>/<file>` directly** — changes auto-sync via stow symlinks
- **`make deploy`** creates symlinks (repo ↔ $HOME)
- **Never commit directly to master** — use `user/<name>` branch, PR to contribute

## Make Targets

| Target | Description |
|--------|-------------|
| `make deploy` | Stow dotfiles to $HOME |
| `make test` | Run Go tests |
| `make commit` | Stage + commit |
| `make pr` | Open PR to master |
| `make erch/init` | Init erch submodule |

## Architecture

`erch/` is the omarchy fork with our customizations baked in (Vim Mode, menu keybinding, custom scripts). `dotfiles/` are portable Layer 3 overrides that work on any system (stock omarchy, clean Arch, etc.).
