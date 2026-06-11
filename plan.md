# eande-os Repo Plan

## Goal

Replace the fragile `os-conf` repo (overrides + wrappers + toggles scattered across layers) with a clean architecture: omarchy fork as submodule + portable dotfiles.

## Changes from os-conf

### Removed (moved to erch fork)
- `dotfiles/omarchy-default/toggles/tiling-mode.conf` → `erch/default/hypr/toggles/tiling-mode.conf`
- `dotfiles/home/.config/erch/extensions/menu.sh` → baked into `erch/bin/omarchy-menu`
- `dotfiles/home/.local/bin/erch-idle` → `erch/bin/erch-idle`
- `dotfiles/home/.local/bin/idle-resume` → `erch/bin/idle-resume`
- `dotfiles/home/.local/bin/erch-scaling-cycle` → `erch/bin/omarchy-hyprland-monitor-scaling-cycle`
- `dotfiles/home/.local/bin/erch-brightness-ddc` → `erch/bin/erch-brightness-ddc`
- `dotfiles/home/.local/bin/erch-source-ddc` → `erch/bin/erch-source-ddc`
- `dotfiles/home/.config/hypr/bindings.conf` (SUPER+SPACE line removed — handled by fork)

### Simplified deploy.sh
- Removed toggle copy to `~/.local/share/erch/default/hypr/toggles/`
- Removed wrapper creation in `~/.local/share/erch/bin/`
- Removed scaling-cycle binary override

### Simplified post-update hook
- Removed toggle restore
- Removed wrapper restore
- Removed scaling-cycle override
- Kept stock theme hiding (still needed)

## Verification

1. Clone fresh: `git clone git@github.com:eande-uk/eande-os.git && cd eande-os && git submodule update --init erch/`
2. Deploy fork: `erch/setup.sh` (or manual copy to `~/.local/share/erch/`)
3. Deploy dotfiles: `make deploy`
4. Verify: `hyprctl reload`, SUPER+SPACE → menu, SUPER+SHIFT+V → toggle, menu → Trigger → Toggle → Vim Mode
5. Verify scripts: `omarchy os-conf idle`, `omarchy os-conf scaling-cycle`
