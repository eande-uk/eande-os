# E&E OS — Agent Guide

## Repo Structure

```
erch/         Git submodule — forked omarchy (github.com/eande-uk/erch)
             Source of Truth for all configs and layer-zero definitions
dotfiles/    Hard copy mirror of erch's L1-L4 configs (stow-deployed)
scripts/     Deployment tooling
tests/       Go verification tests
layer-zero/  Cross-platform profile system (shared by all targets)
docs/        Architecture and reference documentation
```

## Architecture

The ecosystem has **three targets**:

1. **erch** — Standalone forked omarchy. Ships everything natively (L0-L4).
   Clone `github.com/eande-uk/erch`, run `./install.sh`.

2. **Upstream omarchy** — Stock omarchy. Repo provides layer-zero (packages)
   + dotfiles (config overlay). No erch submodule needed.

3. **Arch + Hyprland (no omarchy)** — Bare Arch. Repo provides everything:
   layer-zero for packages, dotfiles for full config.

**erch is the source of truth.** Changes originate in the erch fork and are
mirrored to dotfiles/ for non-erch targets. dotfiles/ is a hard copy that
can fetch erch's current state.

**layer-zero/ is shared identically** across all three targets. It controls
pre-installed packages via use-case profiles (WORK:Office/Dev/AI-ML,
EDUCATION:School/Uni, GAME).

## Key Rules

- **`erch/` is a submodule** — edit the fork at `github.com/eande-uk/erch`,
  not the submodule directly
- **Edit `~/.config/<app>/<file>` directly** — changes auto-sync via stow
  symlinks
- **`make deploy`** creates symlinks (repo ↔ $HOME)
- **Never commit directly to master** — use `user/<name>` branch, PR to
  contribute
- **dotfiles/ mirrors erch** — when adding a new config, add it to erch/
  first, then copy to dotfiles/

## Make Targets

| Target | Description |
|--------|-------------|
| `make deploy` | Stow dotfiles to $HOME |
| `make test` | Run Go tests |
| `make layer-zero` | Interactive package sync |
| `make commit` | Stage + commit |
| `make pr` | Open PR to master |
| `make erch/init` | Init erch submodule |

## Layer System (erch)

| Layer | Name | Phase | Description |
|-------|------|-------|-------------|
| L0 | System | install | Hardware detection, base packages, profile selection |
| L1 | Defaults | install | Shell, env, display server core configs |
| L2 | Configs | install | Application configs (terminal, bar, launcher, editor) |
| L3 | Theme | first-run | Visual branding, color schemes, ASCII art |
| L4 | Polish | post-install | Hooks, migrations, optimizations, toggles |
