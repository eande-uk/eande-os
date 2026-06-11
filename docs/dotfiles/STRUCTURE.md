# Dotfiles Structure

dotfiles/ is a GNU Stow package. Files under `dotfiles/home/` are symlinked
to `$HOME/` by `make deploy`.

## Directory Layout

```
dotfiles/
├── .gitignore                   # Git ignore rules
├── .stow-local-ignore           # Stow skip rules (e.g., hyprland.conf)
└── home/
    ├── .bashrc                  # Shell init (sources erch defaults + .bashrc.d/)
    ├── .bashrc.d/               # Shell overrides (sourced by .bashrc)
    │   ├── 10-env.sh            # EDITOR, TERMINAL, HISTFILE, HISTSIZE
    │   ├── 50-aliases.sh        # Aliases (ls, grep, cat, git, ..)
    │   └── 60-functions.sh      # Functions (mkcd, extract, path)
    └── .config/
        ├── alacritty/           # Terminal emulator
        ├── btop/                # System monitor
        ├── custom-branding/     # Logo source files (PNGs)
        │   ├── about.txt        # ASCII art display logo
        │   ├── screensaver.txt  # ASCII art animation
        │   └── src-pngs/        # Source PNGs for ASCII generation
        ├── fastfetch/           # System info display
        ├── git/                 # Git configuration
        ├── hypr/                # Hyprland window manager
        │   ├── hyprland.conf    # Orchestrator (sources all layers)
        │   ├── bindings.conf    # App keybindings
        │   ├── autostart.conf   # Startup applications
        │   ├── envs.conf        # Environment variables
        │   ├── hypridle.conf    # Idle management
        │   ├── hyprlock.conf    # Screen locker
        │   ├── hyprsunset.conf  # Night light
        │   ├── input.conf       # Keyboard, mouse, touchpad
        │   ├── looknfeel.conf   # Visual settings
        │   ├── monitors.conf    # Display configuration
        │   └── tiling.conf      # Tiling mode settings
        ├── kitty/               # Terminal emulator
        ├── mako/                # Notification daemon
        ├── erch/             # erch-specific overrides
        │   ├── branding/        # ASCII art (about.txt, screensaver.txt)
        │   ├── hooks/           # Post-update, theme-set, font-set
        │   └── themes/          # Curated theme set (8 themes)
        ├── starship.toml        # Shell prompt
        ├── tmux/                # Terminal multiplexer
        ├── uwsm/                # Wayland session manager
        ├── walker/              # Application launcher
        └── waybar/              # Status bar
            ├── config.jsonc
            └── style.css
    └── .local/
        └── bin/                 # User-facing scripts
            ├── erch-commit
            ├── erch-deploy
            ├── erch-docs-verify
            ├── erch-layer-zero
            ├── erch-md-to-html
            ├── erch-mmd
            ├── erch-pr
            ├── erch-rename
            ├── erch-status
            └── erch-test
```

## Stow Mechanism

Deployment uses GNU Stow:

```bash
cd dotfiles/
stow --no-folding -t $HOME home/
```

This creates symlinks:
- `~/.bashrc` → `dotfiles/home/.bashrc`
- `~/.config/hypr/bindings.conf` → `dotfiles/home/.config/hypr/bindings.conf`
- `~/.local/bin/erch-status` → `dotfiles/home/.local/bin/erch-status`

Editing `~/.config/` automatically syncs back to the repo.

### Stow Ignore Rules

`.stow-local-ignore` prevents certain files from being stowed:

```
.config/hypr/hyprland.conf
config/hypr/hyprland.conf
```

`hyprland.conf` is the orchestrator file. On upstream omarchy, it sources
omarchy defaults. On Arch+Hyprland, it may need manual configuration. The
file is provided as a reference but not force-symlinked.

## Config Override Pattern

Every config follows: **inherit defaults → override specific values**.

On **erch**: defaults come from `erch/default/` + `erch/config/`.
On **upstream omarchy**: defaults come from stock omarchy.
On **Arch+Hyprland**: defaults come from the application itself (or are
provided entirely by dotfiles).

dotfiles/ provides only the values that differ from defaults (or the full
config for systems with no upstream defaults).

## Verification

```bash
# Check symlink targets
readlink -f ~/.config/hypr/bindings.conf
# → .../dotfiles/home/.config/hypr/bindings.conf

# Check overall state
make status

# Dry-run stow to see pending changes
make deploy/dry-run
```
