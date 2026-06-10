# Config Reference

All files in this directory are symlinked to `~/.config/` by GNU Stow (via `make deploy`). The repo is the source of truth.

Always work on a `user/<name>` branch (`make branch/create`). Never commit to master directly.

Note: Files directly under `home/` (`.bashrc`, `.bashrc.d/`, `.local/bin/`) are also stowed to `$HOME`.

Every file follows: **inherit** (Omarchy form/defaults) → **override** (repo dotfiles) → **apply** (stow via `make deploy`, backwards compatible with `omarchy cmd`).

No file directly sources Omarchy stock defaults — the only exception is `hypr/hyprland.conf` which orchestrates the layered sourcing.

## Override patterns

| App | Config file(s) | Pattern | omarchy integration | Notes |
|-----|---------------|---------|-------------------|-------|
| **Terminal** | `alacritty/alacritty.toml`, `kitty/kitty.conf` | Theme import | `omarchy restart terminal`, `omarchy refresh` | Theme colors at runtime; font, padding, keybinds inline |
| **hypr** | `hyprland.conf` | Orchestrator | `hyprctl reload` | Sources layered config via `source =` |
| **hypr** | `bindings.conf`, `input.conf`, `looknfeel.conf`, `autostart.conf`, `envs.conf`, `monitors.conf`, `tiling.conf` | Layer 3 override | `hyprctl reload` | Only the values differing from Omarchy defaults (tiling.conf adds mode toggle keybinding) |
| **hypr** | `hypridle.conf`, `hyprsunset.conf` | Inherit→override | `omarchy restart hypridle`, `omarchy restart hyprsunset` | Inherits Omarchy form, overrides from repo |
| **hypr** | `hyprlock.conf` | Theme import | `omarchy theme set` (applies on theme switch) | Sources theme at runtime + full config inline |
| **waybar** | `config.jsonc`, `style.css` | Inherit→override | `omarchy restart waybar`, `omarchy theme refresh` | Inherits Omarchy form, overrides from repo |
| **mako** | `config` | Inherit→override | `omarchy restart mako` | Structural rules and action bindings, no hardcoded colors |
| **fastfetch** | `config.jsonc` | Inherit→override | `omarchy branding about` | Loads ASCII logo from `~/.config/omarchy/branding/about.txt` |
| **git** | `config` | Inherit→override | — | Identity, aliases, diff/push/merge settings |
| **starship** | `starship.toml` | Inherit→override | `omarchy restart terminal` | Minimal git prompt + language stubs |
| **tmux** | `tmux/tmux.conf` | Inherit→override | `omarchy restart tmux` | Terminal multiplexer keybinds and status bar |
| **btop** | `btop.conf` | Inherit→override | `omarchy restart btop`, `omarchy theme refresh` | `color_theme = "current"` defers to Omarchy |
| **uwsm** | `default` | Inherit→override | `omarchy default editor/terminal` | Sets `$EDITOR`, `$TERMINAL` |
| **walker** | `config.toml` | Inherit→override | `omarchy restart walker` | References Omarchy theme location |
| **custom-branding** | `about.txt`, `screensaver.txt`, `src-pngs/` | Inherit→override | `omarchy branding about text`, `omarchy branding screensaver text` | ASCII art generated from PNGs; old omarchy branding removed |
| **omarchy** | `themes/` | Inherit→override | `omarchy theme list\|set\|remove`, `omarchy theme bg next\|set` | 8 themes (4 families × dark/light); stock themes removed |
| **omarchy** | `hooks/` | Inherit→override | `omarchy hook <name>` | `theme-set`, `font-set`, `post-update` |
| **omarchy** | `branding/` | Inherit→override | `omarchy branding about\|screensaver` | Real files (synced from `custom-branding/` by deploy.sh) |

## Omarchy layer sync

Every layer follows the two-direction sync pattern: **desired state → detect via omarchy → install + remove**.

| Layer | Scope | omarchy commands | Install direction | Remove direction |
|-------|-------|-----------------|-------------------|-----------------|
| **0** | System state | `omarchy pkg/webapp/tui/npx install/remove` | `omarchy pkg add`, `omarchy webapp/tui/npx install` | `omarchy pkg drop`, `omarchy webapp/tui remove`, `rm ~/.local/bin/<cmd>` |
| **1** | Compatibility | `omarchy default/config/hook/install/remove` | `omarchy default browser/editor/terminal` | `omarchy remove browser/dev env`, `omarchy config direct boot remove` |
| **2** | Rebranding | `omarchy theme/branding/screensaver` | `omarchy theme set`, `omarchy branding about text` | `omarchy theme remove <name>`, `omarchy branding about reset` |
| **3** | Config overrides | `make deploy` (stow), `omarchy restart/hook/theme refresh` | `make deploy`, `omarchy restart <component>` | `stow -D home` then `stow home` (removes stale symlinks) |
| **4** | Scripts, bins, tests, wrappers, extensions | `make deploy` (stow), `omarchy system/hyprland/brightness/capture`, `omarchy os-conf <cmd>` | `make deploy` (symlinks + chmod +x + wrapper creation + toggle/extension deploy) | Remove from `home/.local/bin/` + `stow --restow home` |

## Config 5-layer architecture

`hyprland.conf` sources configs in this order (later wins):

```
Layer 0: System state — allowlist-based two-direction sync via layer-zero.sh (omarchy pkg/webapp/tui/npx)
Layer 1: Omarchy defaults (~/.local/share/omarchy/default/hypr/*.conf)
Layer 2: Theme overrides (~/.config/omarchy/current/theme/hyprland.conf)
Layer 3: Dotfile overrides (~/.config/hypr/*.conf) — this repo
Layer 4: Dynamic toggles (~/.local/state/omarchy/toggles/hypr/*.conf)
```

## Inspecting a config

```bash
# Check if a file is managed by this repo
readlink -f ~/.config/hypr/bindings.conf
# Should show: .../dotfiles/home/.config/hypr/bindings.conf

# Check overall state
make status

# Verify Hyprland
hyprctl configerrors
```
