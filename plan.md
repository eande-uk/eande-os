# eande-os Repo Plan

## Goal

Replace the fragile `os-conf` repo (overrides + wrappers + toggles scattered across layers) with a clean architecture: omarchy fork as submodule + portable dotfiles.

## Changes from os-conf

### Removed (moved to erch fork)
- `dotfiles/omarchy-default/toggles/tiling-mode.conf` → `erch/default/hypr/toggles/tiling-mode.conf`
- `dotfiles/home/.config/omarchy/extensions/menu.sh` → baked into `erch/bin/omarchy-menu`
- `dotfiles/home/.local/bin/omarchy-os-conf-idle` → `erch/bin/omarchy-os-conf-idle`
- `dotfiles/home/.local/bin/idle-resume` → `erch/bin/idle-resume`
- `dotfiles/home/.local/bin/omarchy-os-conf-scaling-cycle` → `erch/bin/omarchy-hyprland-monitor-scaling-cycle`
- `dotfiles/home/.local/bin/omarchy-os-conf-brightness-ddc` → `erch/bin/omarchy-os-conf-brightness-ddc`
- `dotfiles/home/.local/bin/omarchy-os-conf-source-ddc` → `erch/bin/omarchy-os-conf-source-ddc`
- `dotfiles/home/.config/hypr/bindings.conf` (SUPER+SPACE line removed — handled by fork)

### Simplified deploy.sh
- Removed toggle copy to `~/.local/share/omarchy/default/hypr/toggles/`
- Removed wrapper creation in `~/.local/share/omarchy/bin/`
- Removed scaling-cycle binary override

### Simplified post-update hook
- Removed toggle restore
- Removed wrapper restore
- Removed scaling-cycle override
- Kept stock theme hiding (still needed)

## Verification

1. Clone fresh: `git clone git@github.com:eande-uk/eande-os.git && cd eande-os && git submodule update --init erch/`
2. Deploy fork: `erch/setup.sh` (or manual copy to `~/.local/share/omarchy/`)
3. Deploy dotfiles: `make deploy`
4. Verify: `hyprctl reload`, SUPER+SPACE → menu, SUPER+SHIFT+V → toggle, menu → Trigger → Toggle → Vim Mode
5. Verify scripts: `omarchy os-conf idle`, `omarchy os-conf scaling-cycle`
