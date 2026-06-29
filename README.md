# E&E OS — Distro Hub

**Repo:** `github.com/eande-uk/eande-os`

E&E OS manages multiple E&E Linux distributions from a single hub. Each
distro is a standalone submodule with its own full install pipeline.

**Always work on a `user/<name>` branch** — never commit directly to master.

## Distros

| Distro | Description | Status |
|--------|-------------|--------|
| **erch** | E&E Distro — full-featured flagship with L0-L4 layer system | Active |
| **E-OS** | E&E OS — simpler Arch+Hyprland experience | Planned |
| **E-OS-AI** | E&E OS AI — agent-focused minimal OS | Planned |

## Quick Start

### erch (Full Distro)

```bash
git clone git@github.com:eande-uk/erch.git
cd erch
./install.sh   # Guided install: profile → L0 → L1 → L2 → L3 → L4
```

### From the Hub

```bash
git clone git@github.com:eande-uk/eande-os.git
cd eande-os
make erch/init   # Init erch submodule
make deploy      # Deploy erch
```

## Documentation

| Document | Purpose |
|----------|---------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Hub architecture and layer system |
| [docs/erch/STANDALONE.md](docs/erch/STANDALONE.md) | erch standalone system |
| [docs/erch/LEVELS.md](docs/erch/LEVELS.md) | Level progression (L0-L4) |
| [docs/erch/PROFILES.md](docs/erch/PROFILES.md) | Profile system + managed mode |
| [erch/AGENTS.md](erch/AGENTS.md) | erch dev conventions and command reference |
| [erch/docs/](erch/docs/) | erch architecture, vision, roadmap, features |

## Make Targets

| Target | Description |
|--------|-------------|
| `make init` | Create user branch + init all submodules |
| `make setup` | Full bootstrap: init + erch deploy |
| `make deploy` | Deploy erch to `~/.local/share/erch/` |
| `make status` | Show branch, submodules, changes |
| `make test` | Run Go verification tests |
| `make pr` | Open PR to master |

## Branch Model

- **master** — Production-ready, PR-only
- **user/<name>** — Personal development branches
- Never commit directly to master

```bash
make branch/create   # user/$USER from master
make commit TYPE=feat SCOPE=erch DESC="add new theme"
make pr
```
