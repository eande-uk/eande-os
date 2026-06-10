# Target: erch — Fresh Install Guide

This guide covers a fresh erch install — the full standalone experience
from bare metal to polished desktop.

## Prerequisites

- Arch Linux base install (or running the erch ISO when available)
- Internet connection
- UEFI boot (recommended)

## Install

```bash
# Clone erch
git clone git@github.com:eande-uk/erch.git
cd erch

# Run the installer
./install.sh
```

The installer will guide you through:

### 1. Hardware Detection (L0)

erch detects your hardware and applies the appropriate configuration:

- **ASUS ROG**: ROG-specific drivers, fan control, anime display
- **Dell XPS**: Touchpad haptics, display color profiles
- **Framework 16**: Expansion card config, AMD GPU switch
- **Generic**: Standard Hyprland configuration

### 2. Profile Selection (L0)

Choose your use case profile(s):

```
 Select your profile:

 ┌────────────────────────────────────────────┐
 │ ● WORK ──── ○ Office  ● Dev  ○ AI/ML     │
 │ ○ EDUCATION ── ○ School  ○ Uni           │
 │ ○ GAME                                    │
 │                                           │
 │ [  Confirm  ]  [  Skip (minimal)  ]       │
 └────────────────────────────────────────────┘
```

Multiple profiles can be selected (e.g., Dev + GAME).

### 3. System Install (L0)

The installer:

1. Installs base system (kernel, firmware, drivers)
2. Installs profile packages
3. Configures display manager (SDDM)
4. Sets up audio, networking, printing
5. Installs fonts, icons, themes

### 4. Defaults + Configs (L1 + L2)

1. Deploys core defaults (`erch/default/`)
2. Deploys application configs (`erch/config/`)
3. Configures shell, env, Hyprland, applications

### 5. Theme Activation (L3)

1. Sets default theme based on profile
2. Deploys branding assets
3. Applies color scheme across all components

### 6. Polish (L4)

1. Runs post-install hooks
2. Applies runtime optimizations
3. Sets up dynamic toggles
4. Runs any pending migrations

### 7. Reboot

```bash
sudo reboot
```

## Post-Install

After reboot, your erch system is ready:

- **SUPER+SPACE**: Application menu
- **SUPER+SHIFT+V**: Vim mode toggle
- **SUPER+SHIFT+T**: Tiling mode toggle
- **omarchy theme set "Theme Name"**: Change theme
- **omarchy update system**: System update

## Using dotfiles/ on erch

dotfiles/ is compatible but not required:

```bash
# Clone the parent repo (contains dotfiles/)
git clone git@github.com:eande-uk/eande-os.git
cd eande-os
make deploy  # Optional: overlay dotfiles on erch
```

This is useful for:

- Development: edit `~/.config/` and auto-sync to repo
- Override: temporarily override erch defaults
