# E&E OS — Desktop Configuration Ecosystem

**Repo:** `github.com/eande-uk/eande-os`

Three targets, one config. **erch** is the source of truth; **dotfiles/**
mirrors it for non-erch systems; **layer-zero/** manages packages across
all platforms.

## Architecture

```
erch/         Git submodule — forked omarchy (Source of Truth)
             ships L0-L4 natively on standalone install
dotfiles/    Hard copy mirror of erch's L1-L4 configs (stow-deployed)
layer-zero/  Cross-platform profile system (WORK, EDUCATION, GAME)
scripts/     Deployment tooling for non-erch targets
tests/       Go verification tests
```

See **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** for full details.

## Quick Start

### Target: erch (Standalone)

```bash
git clone git@github.com:eande-uk/erch.git
cd erch
./install.sh   # Guided install: profile → L0 → L1 → L2 → L3 → L4
```

### Target: Upstream Omarchy

```bash
git clone git@github.com:eande-uk/eande-os.git
cd eande-os
make layer-zero   # Install profile packages
make deploy       # Deploy config overlay
hyprctl reload
```

### Target: Arch + Hyprland (No Omarchy)

```bash
git clone git@github.com:eande-uk/eande-os.git
cd eande-os
make layer-zero   # Install everything from scratch
make deploy       # Full config deployment
hyprctl reload
```

## Documentation

| Document | Purpose |
|----------|---------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Full system architecture |
| [docs/erch/STANDALONE.md](docs/erch/STANDALONE.md) | erch standalone system |
| [docs/erch/LEVELS.md](docs/erch/LEVELS.md) | Level progression (L0-L4) |
| [docs/erch/PROFILES.md](docs/erch/PROFILES.md) | Profile system + managed mode |
| [docs/dotfiles/MIRROR.md](docs/dotfiles/MIRROR.md) | Mirror relationship |
| [docs/dotfiles/STRUCTURE.md](docs/dotfiles/STRUCTURE.md) | Stow directory layout |
| [docs/layer-zero/OVERVIEW.md](docs/layer-zero/OVERVIEW.md) | Cross-platform layer-zero |
| [docs/layer-zero/PROFILES.md](docs/layer-zero/PROFILES.md) | Profile categories |
| [docs/targets/ERCH.md](docs/targets/ERCH.md) | erch install guide |
| [docs/targets/ORIGINAL-OMARCHY.md](docs/targets/ORIGINAL-OMARCHY.md) | Omarchy deployment |
| [docs/targets/ARCH-HYPRLAND.md](docs/targets/ARCH-HYPRLAND.md) | Arch+Hyprland setup |
| [docs/REPO.md](docs/REPO.md) | Make targets, scripts, tests |

## Make Targets

| Target | Description |
|--------|-------------|
| `make deploy` | Stow dotfiles to $HOME |
| `make test` | Run Go tests |
| `make layer-zero` | Interactive package sync |
| `make commit` | Stage + commit |
| `make pr` | Open PR to master |
