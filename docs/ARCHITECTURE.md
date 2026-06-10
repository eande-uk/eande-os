# E&E OS — System Architecture

## Overview

E&E OS manages a layered desktop configuration across three target platforms,
with **erch** as the single source of truth. Every config, default, and
package definition originates in erch and is mirrored outward to the rest of
the repo for deployment on other systems.

## Three Targets

```
                    ┌──────────────────────┐
                    │   erch (Standalone)  │
                    │  Source of Truth L0  │
                    │  ships L0-L4 natively│
                    └──────────┬───────────┘
                               │ mirrors
                    ┌──────────▼───────────┐
                    │    dotfiles/         │
                    │  hard copy of erch's │
                    │  L1-L4 configs       │
                    └──────────┬───────────┘
                               │ deploys to
          ┌────────────────────┼────────────────────┐
          ▼                    ▼                    ▼
   ┌──────────────┐   ┌────────────────┐   ┌──────────────┐
   │    erch      │   │   Upstream     │   │Arch+Hyprland │
   │  (native)    │   │   Omarchy      │   │  (from bare) │
   │              │   │                │   │              │
   │ L0-L4 built  │   │ layer-zero +   │   │ layer-zero + │
   │ in. dotfiles │   │ dotfiles as    │   │ dotfiles as  │
   │ compatible . │   │ overlay .      │   │ full config .│
   └──────────────┘   └────────────────┘   └──────────────┘
```

### Target 1: erch (Standalone)

The forked omarchy distribution. Ships everything:

- **L0** (System): Hardware detection, base packages, layer-zero profiles
- **L1** (Defaults): Shell, env, display server, input defaults
- **L2** (Configs): Application configs (terminal, bar, launcher, editor)
- **L3** (Theme): Visual branding, color schemes, ASCII art
- **L4** (Polish): Hooks, migrations, optimizations, dynamic toggles

A fresh erch install progresses automatically through L0 → L4 with profile
selection (WORK, EDUCATION, GAME) at install time.

**dotfiles/** is compatible with erch but not required — erch ships its own
copies natively. dotfiles/ is a mirror for non-erch deployments.

### Target 2: Upstream Omarchy

A stock omarchy installation (not the erch fork). Gets E&E OS configs via:

1. **layer-zero/** — profile-based package management (same system as erch)
2. **dotfiles/** — stow-deployed config overlay (mirrors erch's L1-L4)

No erch fork needed. The repo acts as an add-on config layer.

### Target 3: Arch + Hyprland (No Omarchy)

A bare Arch Linux system with Hyprland (no omarchy at all). Gets:

1. **layer-zero/** — profile-based package management
2. **dotfiles/** — full config deployment (replaces what omarchy would provide)
3. **scripts/** — install helpers for missing omarchy tooling

This target requires the most setup but ends at the same final state.

## Layer Model

| Layer | Name | Lives In | Description |
|-------|------|----------|-------------|
| **L0** | System | `erch/install/` + `layer-zero/` | Hardware detection, base packages, kernel, DM, profile-based pre-installs |
| **L1** | Defaults | `erch/default/` | Core configs: shell init, env vars, input, display server |
| **L2** | Configs | `erch/config/` | Application configs: terminals, bars, launchers, editors, git |
| **L3** | Theme | `erch/themes/` | Visual identity: color schemes, branding, ASCII art, fonts |
| **L4** | Polish | `erch/bin/` + `erch/migrations/` + `erch/hooks/` | Post-install hooks, migrations, optimizations, runtime toggles |

On erch systems, these layers are applied progressively during install.
On non-erch systems, dotfiles/ bundles the equivalent configs and applies
them at once.

## Component Ownership

| Component | erch (SoT) | dotfiles/ (Mirror) | layer-zero/ (Shared) |
|-----------|-----------|-------------------|---------------------|
| L0 Profiles | Authoritative definitions | Not mirrored | Shared by all targets |
| L1 Defaults | `erch/default/` | `dotfiles/home/` copies | N/A |
| L2 Configs | `erch/config/` | `dotfiles/home/.config/` copies | N/A |
| L3 Theme | `erch/themes/` + branding | `dotfiles/home/.config/omarchy/themes/` | N/A |
| L4 Scripts | `erch/bin/` | `dotfiles/home/.local/bin/` copies | N/A |
| Tests | `erch/test/` | `tests/` (cross-platform) | N/A |
| Deployment | erch/setup.sh | scripts/deploy.sh | layer-zero/layer-zero.sh |

### dotfiles/ as Universal Mirror

dotfiles/ is a **hard copy** of erch's L1-L4 configs, structured as a
standalone stow package. It is system-agnostic — the same files deploy to
erch, upstream omarchy, and plain Arch+Hyprland.

Sync direction: **erch → dotfiles/** (manual or tool-assisted copy).
dotfiles/ can fetch the current erch config state on demand.

### layer-zero/ as Cross-Platform Profile System

layer-zero/ controls what gets **pre-installed** on all three targets.
It is the only component shared identically across erch, upstream omarchy,
and Arch+Hyprland.

- **erch**: layer-zero profiles are authoritative definitions
- **Upstream omarchy**: layer-zero runs on top of existing omarchy install
- **Arch+Hyprland**: layer-zero runs on bare Arch, installs everything

Profiles are use-case categories (WORK:Office/Dev/AI-ML, EDUCATION:School/Uni,
GAME) that determine pre-installed packages. After setup, users can freely
install and remove anything. A managed mode exists for organizations that
need policy enforcement (opt-in, documented as desired state).

## Design Principles

1. **erch is the source of truth** — everything originates there
2. **dotfiles/ mirrors erch** — same configs, packaged for non-erch targets
3. **layer-zero is shared** — same profile system across all three targets
4. **System-agnostic configs** — dotfiles work the same regardless of profile
5. **User freedom** — profiles are initial state; users modify freely after
6. **Progressive on erch, flat elsewhere** — erch installs in levels; other
   targets get the full config at once
