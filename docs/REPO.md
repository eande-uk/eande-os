# Repository Reference

This document covers the top-level structure and tooling of the
`eande-os` repository.

## Directory Structure

```
eande-os/
├── docs/              # System architecture documentation
│   ├── ARCHITECTURE.md
│   ├── erch/          # erch standalone docs
│   ├── dotfiles/      # Mirror documentation
│   ├── layer-zero/    # Profile system docs
│   └── targets/       # Per-target guides
├── erch/              # Git submodule — forked omarchy (SoT)
├── dotfiles/          # Hard copy mirror of erch L1-L4 configs
│   └── home/          # Stow package (symlinked to $HOME)
│       ├── .bashrc
│       ├── .bashrc.d/
│       ├── .config/   # Application configs
│       └── .local/    # User scripts
├── layer-zero/        # Cross-platform profile system
│   ├── allowlist.txt  # Core packages (always kept)
│   ├── config.sh      # Active bloat categories
│   ├── layer-zero.sh  # Sync engine
│   ├── bloat/         # Category definitions
│   ├── profiles/      # Profile package lists (WORK, EDUCATION, GAME)
├── scripts/           # Deployment tooling
│   ├── deploy.sh      # Stow-based deploy
│   └── setup.sh       # Bootstrap orchestrator
├── tests/             # Go verification tests
│   ├── verify_test.go
│   ├── e2e_test.go
│   ├── plan_test.go
│   ├── erch_test.go
│   ├── errors_test.go
│   ├── env_test.go
│   └── testutil/      # Test helpers
├── Makefile           # Build targets
├── README.md          # Repo overview
└── AGENTS.md          # Agent guide
```

## Make Targets

| Target | Description |
|--------|-------------|
| **Lifecycle** | |
| `init` | Create branch `user/$USER` from master + deploy |
| `setup` | Full bootstrap: init + layer-zero sync |
| **Deploy** | |
| `deploy` | Link configs via stow (with backup, errors on master) |
| `deploy/dry-run` | Preview what deploy would change |
| `deploy/restock` | Re-apply master defaults |
| `adopt` | Adopt existing `~/.config/` as branch defaults |
| **Inspect** | |
| `status` | Show branch, uncommitted changes, stow state |
| **Layer 0** | |
| `layer-zero` | Interactive two-direction sync |
| `layer-zero/apply` | Apply without confirm |
| `layer-zero/dry-run` | Preview only |
| **Theme** | |
| `theme/list` | `erch theme list` |
| `theme/set NAME=n` | `erch theme set` |
| **Tests** | |
| `test` | Run verification tests (verbose) |
| `test/quiet` | Run verification tests (quiet) |
| **Git** | |
| `diff` | Show uncommitted changes |
| `log` | Recent commits (15) |
| `commit TYPE=t SCOPE=s DESC=d` | Stage all + commit with convention |
| **erch** | |
| `erch/init` | Init submodule + deploy to `~/.local/share/erch/` |
| **Contributing** | |
| `branch/create` | Create `user/$USER` branch from master |
| `pr` | Open PR from current branch → master |

## Scripts

### Deployment Scripts

| Script | Purpose | When to Use |
|--------|---------|-------------|
| `scripts/deploy.sh` | Stow dotfiles to `$HOME` | Fresh installs, bulk updates, CI |
| `scripts/setup.sh` | New-machine bootstrap (prune + deploy) | First-time setup on a new machine |
| `layer-zero/layer-zero.sh` | Two-direction package sync | Package management across all targets |

### User-Facing Scripts

Deployed to `~/.local/bin/` via `make deploy`:

| Script | Purpose |
|--------|---------|
| `erch-commit` | Wrapper → `make commit` |
| `erch-deploy` | Wrapper → `make deploy` |
| `erch-docs-verify` | Docs consistency checker |
| `erch-layer-zero` | Wrapper → `make layer-zero` |
| `erch-md-to-html` | Markdown → HTML via pandoc (shim → `erch wiki html`) |
| `erch-wiki` | Wiki and notes (init, serve, html, pdf, build, notes) |
| `erch-mmd` | Mermaid → SVG/PNG via Docker |
| `erch-pr` | Wrapper → `make pr` |
| `erch-rename` | Batch file renamer via `fd` |
| `erch-status` | Wrapper → `make status` |
| `erch-test` | Wrapper → `go test` |

## Tests

Go verification tests in `tests/`:

| Test File | What It Verifies |
|-----------|-----------------|
| `verify_test.go` | Repo structure, content, deployment integrity |
| `e2e_test.go` | `make` targets, repo structure, erch CLI |
| `plan_test.go` | Docs coherence (AGENTS.md, Makefile, plan.md) |
| `erch_test.go` | Omarchy commands discovery, script metadata |
| `errors_test.go` | Error paths (deploy on master, missing deps) |
| `env_test.go` | Environment variables |

```bash
make test        # Verbose
make test/quiet  # Quiet
```

## Branch Model

- **master** — Production-ready, PR-only
- **user/<name>** — Personal development branches
- Never commit directly to master

```bash
# Create a user branch
make branch/create  # user/$USER from master

# Commit
make commit TYPE=feat SCOPE=waybar DESC="add cpu temperature module"

# Open PR
make pr
```

## Contributing

1. Create a user branch: `make branch/create`
2. Edit configs directly in `~/.config/` (symlinked to repo)
3. Commit: `make commit TYPE=<type> SCOPE=<scope> DESC="<description>"`
4. Open PR: `make pr`
5. Never commit to master directly
