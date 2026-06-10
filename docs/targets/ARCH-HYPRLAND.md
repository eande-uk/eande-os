# Target: Arch + Hyprland (No Omarchy)

This guide covers applying E&E OS configs to a bare Arch Linux system with
Hyprland — no omarchy installed.

## Overview

On bare Arch+Hyprland, the E&E OS repo provides what omarchy normally would:

1. **layer-zero/** manages package installation from scratch
2. **dotfiles/** provides the full config (Hyprland, Waybar, Mako, etc.)
3. **scripts/** provide helpers that omarchy tooling would normally handle

This target requires more setup because omarchy's defaults are missing, but
the final state is the same desktop experience.

## Prerequisites

- Arch Linux base install (arch-chroot or running system)
- Hyprland installed (`sudo pacman -S hyprland`)
- `git`, `stow`, `gum` installed
- A display manager or `uwsm` for Wayland session management

## Quick Start

```bash
# Clone the repo
git clone git@github.com:eande-uk/eande-os.git
cd eande-os

# (Optional) Select a profile for layer-zero
# Edit layer-zero/config.sh to enable profile categories

# Run layer-zero to install packages
make layer-zero

# Deploy dotfiles (full config — there are no omarchy defaults)
make deploy

# Reload Hyprland
hyprctl reload
```

## Pre-Install: Required Packages

Before deploying, ensure these base packages are installed:

```bash
sudo pacman -S --needed \
  hyprland hypridle hyprlock hyprsunset hyprpicker \
  waybar mako walker uwsm \
  alacritty kitty tmux \
  starship fastfetch btop \
  git github-cli lazygit \
  neovim \
  xdg-desktop-portal-hyprland xdg-desktop-portal-gtk \
  polkit-gnome gnome-keyring \
  swaybg swayosd grim slurp satty wl-clipboard \
  wireplumber pamixer playerctl \
  fcitx5 fcitx5-gtk fcitx5-qt \
  noto-fonts noto-fonts-cjk noto-fonts-emoji \
  ttf-jetbrains-mono-nerd \
  iwd avahi nss-mdns \
  cups cups-browsed system-config-printer \
  docker docker-compose
```

Or use `layer-zero/` to install everything:

```bash
# First install layer-zero dependencies
sudo pacman -S --needed omarchy  # Not available on bare Arch — see FAQ
```

If omarchy commands are not available, layer-zero cannot dispatch through
`omarchy pkg`, `omarchy webapp`, etc. In this case, use the profile package
lists directly with pacman:

```bash
# Install from a profile
sudo pacman -S --needed $(cat layer-zero/profiles/WORK/Dev.pkgs)
```

## Deployment Steps

### 1. Profile Selection

Edit `layer-zero/config.sh` to enable desired categories:

```bash
PRUNE_PACMAN_CATEGORIES=(
    # gaming
    media
    office
    # communication
    # browsers
    runtimes
    # terminals
)
```

Or use a profile:

```bash
# Copy profile packages to allowlist
cat layer-zero/profiles/WORK/Dev.pkgs >> layer-zero/allowlist.txt
```

### 2. Layer Zero Sync

```bash
# Install all profile packages
./layer-zero/layer-zero.sh --apply
```

### 3. Deploy Dotfiles

Since there are no omarchy defaults, dotfiles/ provides the complete config:

```bash
make deploy
```

### 4. Manual Steps (No Omarchy)

Some things normally handled by omarchy must be done manually:

```bash
# Set default theme
omarchy theme set "Tokyo Night" -- if omarchy is installed

# Without omarchy: manually symlink or copy theme files
ln -sf ~/.config/omarchy/themes/tokyo-night/colors.toml ~/.config/colors.toml
```

### 5. Post-Deploy

```bash
# Restart Hyprland
hyprctl reload

# Restart services
systemctl --user restart waybar
systemctl --user restart mako
systemctl --user restart wireplumber

# Verify
make status
```

## FAQ

### Q: Can I install omarchy on bare Arch to get full compatibility?

**A:** Yes. `erch/` IS omarchy (the fork). Clone it separately:

```bash
git clone git@github.com:eande-uk/erch.git
cd erch
./install.sh
```

Then deploy `eande-os` dotfiles on top. This gives the full erch experience
on bare Arch.

### Q: What doesn't work without omarchy?

**A:** These dotfiles features depend on omarchy CLI commands:
- `omarchy theme set/list/refresh` — theme switching
- `omarchy restart <component>` — restart scripts
- `omarchy refresh <config>` — config refresh
- `omarchy toggle <feature>` — dynamic toggles
- `omarchy pkg/webapp/tui/npx` — package management (layer-zero)

Workarounds are noted in the relevant config files.

### Q: How do I handle `omarchy` commands in dotfiles configs?

**A:** Configs are designed to fall back gracefully. If an `omarchy` command
is not found, the config skips the failing source/include. Check individual
configs for `omarchy cmd` guards.
