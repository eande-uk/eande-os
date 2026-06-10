# Layer Zero — Cross-Platform Overview

Layer Zero controls what gets **pre-installed** on all three targets (erch,
upstream omarchy, Arch+Hyprland). It is the only component shared identically
across the ecosystem.

## What Layer Zero Does

Layer Zero is a two-direction sync engine:

- **Items on the allowlist** → ensure installed
- **Items NOT on the allowlist** → remove (prune)

It operates on packages, webapp launchers, TUI launchers, and NPX wrappers.
All operations dispatch through omarchy CLI commands (`omarchy pkg add`,
`omarchy pkg drop`, `omarchy webapp install`, etc.).

## How It Works on Each Target

| Target | layer-zero Behavior |
|--------|---------------------|
| **erch** | Integrated at install time (L0). User selects a profile, profile packages are installed. After install, layer-zero runs to sync the allowlist (keeps profile packages, prunes active bloat categories not on allowlist). |
| **Upstream omarchy** | Runs on top of existing omarchy install. layer-zero syncs packages per the selected profile. dotfiles/ provides config overlay. |
| **Arch+Hyprland** | Runs on bare Arch. Installs everything from scratch. dotfiles/ provides the config layer that omarchy would normally provide. |

## Component Location

| File | Purpose |
|------|---------|
| `layer-zero/allowlist.txt` | Core packages to keep (always installed) |
| `layer-zero/profiles/` | Profile package lists (WORK, EDUCATION, GAME) |
| `layer-zero/config.sh` | Active bloat categories (which categories to prune) |
| `layer-zero/layer-zero.sh` | Two-direction sync engine |
| `layer-zero/bloat/` | Package category definitions for prune mode |

## erch as Source of Truth

The authoritative definitions for layer-zero live in **erch/**. The
`layer-zero/` directory in the parent repo is a copy that ships to non-erch
targets.

When a new profile or package definition is added, it originates in erch
and is mirrored to the parent repo's `layer-zero/`.

## Future: Profile Selector

Desired state — a unified installer that works across all three targets:

```
# On any target (erch, omarchy, Arch+Hyprland):
layer-zero/install.sh --profile WORK:Dev
```

This script would:

1. Detect the target (erch/omarchy/bare)
2. Install the profile's packages using the appropriate method
3. Sync the allowlist
4. Report what was installed and what was pruned
