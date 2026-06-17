# erch — Standalone System

erch is the E&E OS fork of omarchy. It ships the complete desired desktop
configuration natively — from bare metal to fully polished — without
requiring any external repo.

## What erch Ships

erch bundles everything as built-in defaults:

| Layer | erch Location | Contents |
|-------|--------------|----------|
| L0 Profiles | `erch/install/` + `layer-zero/` definitions | Profile selection, hardware detection, base install |
| L1 Defaults | `erch/default/` | Shell init, env vars, input, display server |
| L2 Configs | `erch/config/` | Terminal, bar, launcher, editor, git, etc. |
| L3 Theme | `erch/themes/` + `erch/default/branding/` | Color schemes, ASCII art, fonts |
| L4 Polish | `erch/bin/` + `erch/migrations/` + `erch/hooks/` | Scripts, migrations, post-install hooks |

## Install Flow

A fresh erch install progresses through these phases:

```
boot.sh / install.sh
        │
        ▼
┌── L0: System ──────────────────────────────────┐
│  • Hardware detection                           │
│  • Base packages (kernel, drivers, firmware)    │
│  • Profile selection (see PROFILES.md)          │
│  • Install pre-selected packages per profile    │
│  • Display manager, networking, audio           │
└─────────────────────────────────────────────────┘
        │
        ▼
┌── L1: Defaults ────────────────────────────────┐
│  • Copy erch/default/ to ~/.local/share/erch/   │
│  • Shell: .bashrc, aliases, completions, env    │
│  • Display server: Hyprland base configs        │
│  • Input: keyboard, touchpad, fcitx5            │
│  • Fonts, icons, cursors                        │
└─────────────────────────────────────────────────┘
        │
        ▼
┌── L2: Configs ─────────────────────────────────┐
│  • Copy erch/config/ to ~/.config/              │
│  • Terminal (alacritty, kitty)                  │
│  • Bar (waybar), launcher (walker)              │
│  • Editor (neovim), git, tmux                   │
│  • Notifications (mako), OSD (swayosd)          │
│  • Portal, auth, file manager                   │
└─────────────────────────────────────────────────┘
        │
        ▼
┌── L3: Theme ───────────────────────────────────┐
│  • Set default theme from selected profile      │
│  • Deploy branding (ASCII art, logo)            │
│  • Apply color scheme across all components     │
│  • First-run theme set                          │
└─────────────────────────────────────────────────┘
        │
        ▼
┌── L4: Polish ──────────────────────────────────┐
│  • Post-install hooks                           │
│  • Runtime optimizations                        │
│  • Dynamic toggle setup                         │
│  • First-run migrations                         │
│  • Wiki & notes: erch wiki (mdbook, pandoc)     │
└─────────────────────────────────────────────────┘
        │
        ▼
   Ready
```

## Relationship to Parent Repo

The `eande-os` parent repo contains:

- **erch/** — the submodule (this system)
- **dotfiles/** — a hard copy mirror of erch's L1-L4 configs for non-erch
  targets (upstream omarchy, Arch+Hyprland)
- **layer-zero/** — the shared profile system (authoritative definitions
  live in erch; the parent repo's copy deploys to non-erch targets)
- **scripts/** — deployment tooling for non-erch targets
- **tests/** — cross-platform verification

erch does not depend on the parent repo. The parent repo exists to ship
erch's configs to systems that don't run erch.

## dotfiles/ Compatibility

dotfiles/ works on erch as a stow overlay. This is useful for:

- Development: edit `~/.config/` and auto-sync to repo
- Testing: verify configs work on erch before shipping upstream
- Override: temporarily override erch's built-in defaults

On a stock erch install, dotfiles/ is not needed — everything is already
in place natively.
