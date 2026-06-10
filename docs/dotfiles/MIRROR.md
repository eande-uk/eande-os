# Dotfiles — Universal Mirror

dotfiles/ is a **hard copy** of erch's L1-L4 configs, structured as a
standalone stow package. It exists to deploy the same E&E OS desktop
configuration on systems that don't run erch natively.

## Mirror Relationship

```
erch (Source of Truth)          dotfiles/ (Mirror)
─────────────────────           ────────────────────
erch/default/bash/         →    dotfiles/home/.bashrc, dotfiles/home/.bashrc.d/
erch/config/alacritty/     →    dotfiles/home/.config/alacritty/
erch/config/kitty/         →    dotfiles/home/.config/kitty/
erch/config/waybar/        →    dotfiles/home/.config/waybar/
erch/config/walker/        →    dotfiles/home/.config/walker/
erch/config/mako/          →    dotfiles/home/.config/mako/
erch/config/hypr/          →    dotfiles/home/.config/hypr/
erch/config/tmux/          →    dotfiles/home/.config/tmux/
erch/config/starship.toml  →    dotfiles/home/.config/starship.toml
erch/config/fastfetch/     →    dotfiles/home/.config/fastfetch/
erch/themes/               →    dotfiles/home/.config/omarchy/themes/
erch/default/branding/     →    dotfiles/home/.config/omarchy/branding/
erch/default/hooks/        →    dotfiles/home/.config/omarchy/hooks/
erch/bin/omarchy-os-conf-* →    dotfiles/home/.local/bin/
```

Not all erch content is mirrored. Only the user-facing configs and scripts
that make sense on non-erch systems. Omarchy-internal machinery (install
scripts, migrations, 280+ omarchy CLI commands) stays in erch.

## Sync Direction

**Always erch → dotfiles/**. Changes originate in erch (the fork) and are
copied to dotfiles/ for the parent repo.

dotfiles/ can fetch the current erch config state:

```
./scripts/sync-from-erch.sh    # Copy erch configs → dotfiles/
```

This script (desired state — not yet implemented) would:

1. Copy `erch/config/<app>/` → `dotfiles/home/.config/<app>/`
2. Copy `erch/default/bash/` → `dotfiles/home/.bashrc.d/` (with merging)
3. Copy `erch/themes/` → `dotfiles/home/.config/omarchy/themes/`
4. Copy `erch/bin/omarchy-os-conf-*` → `dotfiles/home/.local/bin/`
5. Verify no configs are missing in dotfiles/

## System Agnosticism

dotfiles/ is designed to work on **any** system:

- **erch**: works as a stow overlay on top of erch's native defaults
- **Upstream omarchy**: replaces/configures omarchy defaults via stow
- **Arch + Hyprland**: provides the complete config where no omarchy defaults
  exist

The same files deploy everywhere. layer-zero/ handles package management;
dotfiles/ handles config files.

## What Lives Only in erch (Not Mirrored)

| Content | Reason |
|---------|--------|
| `erch/install/` | erch-specific install scripts |
| `erch/migrations/` | Upgrade path between erch versions |
| `erch/bin/omarchy-*` (most CLI) | Omarchy-internal commands |
| `erch/config/omarchy/` | Omarchy's own config (not user-facing) |
| `erch/test/` | erch-specific tests |
