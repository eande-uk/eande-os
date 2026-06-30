# E&E OS — Agent Guide

## Repo Structure

```
eande-os/        Distro hub — manages multiple E&E Linux distributions
├── erch/        Git submodule — E&E Distro (standalone, full L0-L4)
├── E-OS/        Git submodule — simpler Arch+Hyprland distro (4 profiles)
├── E-OS-AI/     Planned — agent-focused minimal OS (repo pending)
├── docs/        Architecture and reference documentation
└── tests/       Go verification tests
```

## Architecture

E&E OS is a **distro hub**. Each submodule is a standalone Linux distribution
that manages everything from kernel to userspace via a modular layer system
(L0-L4). The hub provides orchestration and shared documentation.

### Three Distros

| Distro | Purpose | Status |
|--------|---------|--------|
| **erch** | E&E Distro — full-featured, flagship | Active |
| **E-OS** | E&E OS — simpler Arch+Hyprland experience with 4 profiles | Active |
| **E-OS-AI** | E&E OS AI — agent-focused minimal OS | Planned |

### Layer System (Shared)

All distros share the same modular layer architecture:

| Layer | Name | Phase | Description |
|-------|------|-------|-------------|
| L0 | System | install | Hardware detection, base packages, profile selection |
| L1 | Defaults | install | Shell, env, display server core configs |
| L2 | Configs | install | Application configs (terminal, bar, launcher, editor) |
| L3 | Theme | first-run | Visual branding, color schemes, ASCII art |
| L4 | Polish | post-install | Hooks, migrations, optimizations, toggles |

## Key Rules

- **Each submodule is standalone** — work out-of-box independently
- **erch is the flagship** — full L0-L4, the standard other distros follow
- **Edit submodules in their own repos** — not directly in the submodule path
- **Never commit directly to master** — use `user/<name>` branch, PR to contribute
- **Hub orchestrates** — Makefile targets manage submodule lifecycle

## Make Targets

| Target | Description |
|--------|-------------|
| **Lifecycle** | |
| `init` | Create branch `user/$USER` from master + init all submodules |
| `setup` | Full bootstrap: init + erch deploy |
| **Deploy** | |
| `deploy` | Deploy erch to `~/.local/share/erch/` |
| **Inspect** | |
| `status` | Show branch, submodules, uncommitted changes |
| **Submodules** | |
| `erch/init` | Init erch submodule |
| `E-OS/init` | Init E-OS submodule |
| **Tests** | |
| `test` | Run verification tests (verbose) |
| `test/quiet` | Run verification tests (quiet) |
| **Git** | |
| `diff` | Show uncommitted changes |
| `log` | Recent commits (15) |
| `commit TYPE=t SCOPE=s DESC=d` | Stage all + commit with convention |
| **Contributing** | |
| `branch/create` | Create `user/$USER` branch from master |
| `pr` | Open PR from current branch → master |
