# E&E OS — System Architecture

## Overview

E&E OS is a **distro hub** — a single repository that orchestrates multiple
standalone Linux distributions. Each distro lives in its own submodule with
a complete install pipeline from kernel to userspace.

## Three Distros

```
                    ┌──────────────────┐
                    │    E&E OS Hub    │
                    │   (eande-os)     │
                    │  orchestration + │
                    │  shared docs     │
                    └────────┬─────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
       ┌─────────────┐ ┌──────────┐ ┌─────────────┐
       │    erch     │ │   E-OS   │ │  E-OS-AI    │
       │  Flagship   │ │ Simpler  │ │ Agent OS    │
       │  Full L0-L4 │ │ Arch+Hypr│ │ Minimal     │
       └─────────────┘ └──────────┘ └─────────────┘
```

### erch — E&E Distro (Flagship)

Full-featured standalone distro. Ships everything natively:

- **L0** (System): Hardware detection, base packages, profile selection
- **L1** (Defaults): Shell, env, display server, input defaults
- **L2** (Configs): Application configs (terminal, bar, launcher, editor)
- **L3** (Theme): Visual branding, color schemes, ASCII art
- **L4** (Polish): Hooks, migrations, optimizations, dynamic toggles

Clone `github.com/eande-uk/erch`, run `./install.sh`.

### E-OS — E&E OS (Planned)

Simpler Arch+Hyprland experience. Standalone distro for users who want
the E&E desktop without the full erch install pipeline.

### E-OS-AI — E&E OS AI (Planned)

Agent-focused minimal OS. Designed for autonomous agents that need a
clean, predictable Linux environment.

## Layer System

All distros share the same modular layer architecture:

| Layer | Name | Phase | Description |
|-------|------|-------|-------------|
| **L0** | System | install | Hardware detection, base packages, profile selection |
| **L1** | Defaults | install | Shell, env, display server core configs |
| **L2** | Configs | install | Application configs (terminal, bar, launcher, editor) |
| **L3** | Theme | first-run | Visual branding, color schemes, ASCII art |
| **L4** | Polish | post-install | Hooks, migrations, optimizations, toggles |

### Layer Ownership

| Layer | erch | E-OS | E-OS-AI |
|-------|------|------|---------|
| L0 System | Full install pipeline | Simplified pipeline | Minimal pipeline |
| L1 Defaults | Shell, Hyprland, env | Shell, Hyprland, env | Agent-first defaults |
| L2 Configs | Full app configs | Core app configs | Minimal configs |
| L3 Theme | 21 built-in themes | Curated subset | No theme (headless) |
| L4 Polish | 300+ CLI commands, migrations | Subset | Agent hooks only |

## Design Principles

1. **Each submodule is standalone** — works out-of-box independently
2. **erch is the standard** — other distros follow its layer architecture
3. **Hub orchestrates** — Makefile targets manage submodule lifecycle
4. **Shared layer model** — L0-L4 pattern applies to all distros
5. **No shared code between submodules** — each is self-contained

## Component Ownership

| Component | Location | Notes |
|-----------|----------|-------|
| erch distro | `erch/` submodule | Standalone, full L0-L4 |
| E-OS distro | `E-OS/` submodule | Standalone, simplified |
| E-OS-AI distro | `E-OS-AI/` submodule | Standalone, agent-focused |
| Hub orchestration | `Makefile` | Submodule lifecycle |
| Hub tests | `tests/` | Cross-distro verification |
| Hub docs | `docs/` | Architecture, erch reference |
