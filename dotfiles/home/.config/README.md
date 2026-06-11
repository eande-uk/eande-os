# Config Reference

All files in this directory are symlinked to `~/.config/` by GNU Stow (via `make deploy`). The repo is the source of truth.

Always work on a `user/<name>` branch (`make branch/create`). Never commit to master directly.

Note: Files directly under `home/` (`.bashrc`, `.bashrc.d/`, `.local/bin/`) are also stowed to `$HOME`.

Every file follows: **inherit** (erch/upstream form/defaults) → **override** (repo dotfiles) → **apply** (stow via `make deploy`).

No file directly sources erch/upstream stock defaults — the only exception is `hypr/hyprland.conf` which orchestrates the layered sourcing.

## Override patterns

| App | Config file(s) | Pattern | erch integration | Notes |
|-----|---------------|---------|-------------------|-------|
| **Terminal** | `alacritty/alacritty.toml`, `kitty/kitty.conf` | Theme import | `erch restart terminal`, `erch refresh` | Theme colors at runtime; font, padding, keybinds inline |
| **hypr** | `hyprland.conf` | Orchestrator | `hyprctl reload` | Sources layered config via `source =` |
| **hypr** | `bindings.conf`, `input.conf`, `looknfeel.conf`, `autostart.conf`, `envs.conf`, `monitors.conf`, `tiling.conf` | Layer 3 override | `hyprctl reload` | Only the values differing from erch/upstream defaults (tiling.conf adds mode toggle keybinding) |
| **hypr** | `hypridle.conf`, `hyprsunset.conf` | Inherit→override | `erch restart hypridle`, `erch restart hyprsunset` | Inherits erch/upstream form, overrides from repo |
| **hypr** | `hyprlock.conf` | Theme import | `erch theme set` (applies on theme switch) | Sources theme at runtime + full config inline |
| **waybar** | `config.jsonc`, `style.css` | Inherit→override | `erch restart waybar`, `erch theme refresh` | Inherits erch/upstream form, overrides from repo |
| **mako** | `config` | Inherit→override | `erch restart mako` | Structural rules and action bindings, no hardcoded colors |
| **fastfetch** | `config.jsonc` | Inherit→override | `erch branding about` | Loads ASCII logo from `~/.config/erch/branding/about.txt` |
| **git** | `config` | Inherit→override | — | Identity, aliases, diff/push/merge settings |
| **starship** | `starship.toml` | Inherit→override | `erch restart terminal` | Minimal git prompt + language stubs |
| **tmux** | `tmux/tmux.conf` | Inherit→override | `erch restart tmux` | Terminal multiplexer keybinds and status bar |
| **btop** | `btop.conf` | Inherit→override | `erch restart btop`, `erch theme refresh` | `color_theme = "current"` defers to erch/upstream |
| **uwsm** | `default` | Inherit→override | `erch default editor/terminal` | Sets `$EDITOR`, `$TERMINAL` |
| **walker** | `config.toml` | Inherit→override | `erch restart walker` | References erch theme location |
| **custom-branding** | `about.txt`, `screensaver.txt`, `src-pngs/` | Inherit→override | `erch branding about text`, `erch branding screensaver text` | ASCII art generated from PNGs; old erch branding removed |
| **erch** | `themes/` | Inherit→override | `erch theme list\|set\|remove`, `erch theme bg next\|set` | 8 themes (4 families × dark/light); stock themes removed |
| **erch** | `hooks/` | Inherit→override | `erch hook <name>` | `theme-set`, `font-set`, `post-update` |
| **erch** | `branding/` | Inherit→override | `erch branding about\|screensaver` | Real files (synced from `custom-branding/` by deploy.sh) |

## erch layer sync

Every layer follows the two-direction sync pattern: **desired state → detect via erch → install + remove**.

| Layer | Scope | erch commands | Install direction | Remove direction |
|-------|-------|---------------|-------------------|-----------------|
| **0** | System state | `erch pkg/webapp/tui/npx install/remove` | `erch pkg add`, `erch webapp/tui/npx install` | `erch pkg drop`, `erch webapp/tui remove`, `rm ~/.local/bin/<cmd>` |
| **1** | Compatibility | `erch default/config/hook/install/remove` | `erch default browser/editor/terminal` | `erch remove browser/dev env`, `erch config direct boot remove` |
| **2** | Rebranding | `erch theme/branding/screensaver` | `erch theme set`, `erch branding about text` | `erch theme remove <name>`, `erch branding about reset` |
| **3** | Config overrides | `make deploy` (stow), `erch restart/hook/theme refresh` | `make deploy`, `erch restart <component>` | `stow -D home` then `stow home` (removes stale symlinks) |
| **4** | Scripts, bins, tests, wrappers, extensions | `make deploy` (stow), `erch system/hyprland/brightness/capture`, `erch <cmd>` | `make deploy` (symlinks + chmod +x + wrapper creation + toggle/extension deploy) | Remove from `home/.local/bin/` + `stow --restow home` |

## Config 5-layer architecture

`hyprland.conf` sources configs in this order (later wins):

```
Layer 0: System state — allowlist-based two-direction sync via layer-zero.sh (erch pkg/webapp/tui/npx)
Layer 1: erch/upstream defaults (~/.local/share/erch/default/hypr/*.conf)
Layer 2: Theme overrides (~/.config/erch/current/theme/hyprland.conf)
Layer 3: Dotfile overrides (~/.config/hypr/*.conf) — this repo
Layer 4: Dynamic toggles (~/.local/state/erch/toggles/hypr/*.conf)
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
