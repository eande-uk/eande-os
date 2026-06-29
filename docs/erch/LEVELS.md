# erch Level System

erch applies configuration in five progressive levels during a fresh install.
Each level builds on the previous one, from bare system to fully polished.

## Level 0: System

**Phase:** `install/preflight/` + `install/packaging/`

The foundation. Hardware detection, base system packages, and profile-based
software selection.

| Component | What ships | Details |
|-----------|-----------|---------|
| Hardware detection | `erch/bin/erch-hw-*` | ASUS ROG, Dell XPS, Framework 16, etc. |
| Base packages | `install/packaging/base.sh` | Kernel, firmware, drivers, essential CLI |
| Display manager | SDDM | Login manager + Plymouth boot splash |
| Networking | iwd, avahi, NetworkManager | Wireless, mDNS, firewall (ufw) |
| Audio | WirePlumber, PipeWire | Sound server, OSD, volume control |
| Fonts | `install/packaging/fonts.sh` | JetBrains Mono, Noto, Font Awesome |
| Icons | Yaru, Gnome themes | Icon theme, GTK theme, Qt theme |
| **Profile packages** | Per profile package lists | WORK:Office/Dev/AI-ML, EDUCATION:School/Uni, GAME |

Profile selection happens at this level. See [PROFILES.md](PROFILES.md) for
package composition per profile.

## Level 1: Defaults

**Phase:** `install/defaults/`

Core operating environment. Shell, environment variables, input, and display
server defaults.

| Component | erch path | Ships |
|-----------|-----------|-------|
| Shell init | `default/bash/` | `.bashrc`, aliases, completions, env vars, functions |
| Display server | `default/hypr/` | Hyprland core: envs, input, binds, monitors, window rules |
| Input | `default/hypr/input.conf` | Keyboard layout, touchpad gestures, fcitx5 |
| Display management | `default/hypr/hypridle.conf`, `default/hypr/hyprlock.conf` | Idle/lock, DPMS, screen locker |
| Night light | `default/hypr/hyprsunset.conf` | Blue light filter |
| Compositor utilities | hyprctl, hyprpicker, hyprland-preview-share-picker | Screen sharing, color picker |
| OpenSSH | GPG, SSH agent | Key management, agent startup |

## Level 2: Configs

**Phase:** `install/config/`

Application-level configuration. Terminals, bars, launchers, editors, and
tools deployed to `~/.config/`.

| Component | erch config path | Details |
|-----------|-----------------|---------|
| Terminal | `config/alacritty/`, `config/kitty/`, `config/foot/`, `config/ghostty/` | Multiple terminal options, theme-based colors |
| Bar | `config/waybar/` | Workspaces, clock, weather, audio, network, battery, tray |
| Launcher | `config/walker/` | Application launcher with theme integration |
| Notifications | `config/mako/` | Notification daemon with actions and keybinds |
| OSD | `config/swayosd/` | On-screen display for volume, brightness |
| Editor | `config/nvim/` | Neovim config |
| File manager | `config/nautilus-python/` | Nautilus extensions |
| Git | `config/git/config` | Aliases, diff settings, merge strategy |
| Tmux | `config/tmux/tmux.conf` | Terminal multiplexer with erch theme |
| Starship | `config/starship.toml` | Prompt with git status and language info |
| Fastfetch | `config/fastfetch/config.jsonc` | System info with branding logo |
| Btop | `config/btop/btop.conf` | System monitor with theme colors |
| GPG | `config/gpg/` | GPG agent config |
| Chromium | `config/chromium/` | Browser policies and extensions |
| UWSM | `config/uwsm/default` | Wayland session manager config |
| XDG | `config/xdg-terminals.list` | Default terminal definitions |

## Level 3: Theme

**Phase:** `install/first-run/` + `install/post-install/`

Visual identity. Color schemes, branding assets, font selection, and
first-run theme activation.

| Component | erch path | Details |
|-----------|-----------|---------|
| Color themes | `themes/` (21 themes) | Catppuccin, Nord, Tokyo Night, Gruvbox, etc. |
| Branding | `default/branding/` | ASCII art logos (`about.txt`, `screensaver.txt`) |
| Branding source | PNG source files | `ee-logo.png`, `ee-mark.png`, `water-mark.png` |
| Theme hooks | `default/hooks/` | `theme-set`, `font-set`, `post-update` |
| Default theme | Set per profile | Varies by use case (dark theme for dev, etc.) |

## Level 4: Polish

**Phase:** `install/post-install/` + `migrations/`

Final touches. Post-install optimizations, migrations, runtime state,
and dynamic toggles.

| Component | erch path | Details |
|-----------|-----------|---------|
| Binary scripts | `bin/erch-*` (317 commands) | All erch CLI commands |
| Custom scripts | `bin/erch-os-conf-*` | DDC brightness, idle, scaling cycle |
| Migrations | `migrations/` (343 scripts) | Upgrade path between versions |
| Hooks | `default/hooks/` | Post-update hook, theme-set hook, font-set hook |
| Dynamic toggles | `default/hypr/toggles/` | Vim mode, tiling mode, night light, touchpad |
| Refresh scripts | `bin/erch-refresh-*` | Copy defaults → user config on demand |
| Restart scripts | `bin/erch-restart-*` | Restart individual components |
| Snapshot | `bin/erch-snapshot` | System state snapshots |
