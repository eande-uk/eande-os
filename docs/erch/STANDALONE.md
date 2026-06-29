# erch — Standalone System

erch is the E&E Distro — the flagship standalone Linux distribution that
ships the complete desired desktop configuration natively, from bare metal
to fully polished, without requiring any external repo.

## What erch Ships

erch bundles everything as built-in defaults:

| Layer | erch Location | Contents |
|-------|--------------|----------|
| L0 Profiles | `erch/install/` | Profile selection, hardware detection, base install |
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

## Relationship to Hub

The `eande-os` hub repo contains erch as a submodule alongside other distros
(E-OS, E-OS-AI). erch is standalone — it does not depend on the hub.

The hub provides:
- Orchestration (Makefile targets for init, deploy, test)
- Cross-distro documentation
- Verification tests
