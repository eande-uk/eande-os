# erch Profile System

During a fresh erch install, the user selects a **profile** that determines
which packages are pre-installed at L0.

## Profiles

| Category | Profile | Use Case |
|----------|---------|----------|
| **WORK** | Office | Non-developer employees: office suite, communication, basic tools |
| **WORK** | Dev | Software developers: compilers, runtimes, containers, editors |
| **WORK** | AI/ML | AI/ML engineers: CUDA, Jupyter, ML frameworks, data tools |
| **EDUCATION** | School | Primary/secondary: educational software, content filtering |
| **EDUCATION** | Uni | University: research tools, scientific computing, LaTeX |
| **GAME** | - | Gaming: Steam, Heroic, Lutris, gaming-oriented optimizations |

### Profile Composition

Each profile ships a curated set of packages at install time. Packages are
organized in categories following the same structure as `layer-zero/bloat/`:

```
install/packages/
├── common.pkgs         # Always installed
├── work.pkgs           # WORK profile packages
├── school.pkgs         # EDUCATION:School packages
├── game.pkgs           # GAME profile packages
└── all.pkgs            # Everything
```

Profiles can be composed. A user might select both `WORK:Dev` and `GAME`,
getting the union of both package sets.

### Base Allowlist (Always Installed)

Regardless of profile, every erch install includes the core allowlist:

- Desktop environment (Hyprland, Waybar, Mako, Walker)
- Terminal (Alacritty, Kitty, Foot)
- Shell tooling (bash-completion, bat, fd, fzf, ripgrep, zoxide)
- Git + GitHub CLI
- Neovim
- Fonts (JetBrains Mono Nerd, Noto, Font Awesome)
- Desktop portals, auth, secrets
- Display tools (grim, slurp, satty, wl-clipboard)
- Audio (WirePlumber, pamixer)
- Input (fcitx5)
- Network (iwd, avahi)
- Printing (CUPS)
- Docker + Docker Compose

## Install-Time Selection

During erch install, the user is prompted to select their profile(s):

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

Selection is optional. "Skip" installs only the base allowlist.

## User Freedom After Install

Profiles are a **one-time initial state**, not a restriction. After install:

- Users can install any package with pacman/yay
- Users can remove any pre-installed package
- No profile enforcement in normal mode
- The system is fully open — same as any Arch Linux install

## Managed Mode (Organizations)

For organizations, schools, and institutions that need control, erch offers
an optional **managed mode** (opt-in, not default).

### How Managed Mode Works

1. An administrator defines a **policy** (which profiles are allowed, which
   packages are pinned, which are blocked)
2. The policy is deployed to machines (via config management or USB install
   image)
3. Managed mode enforces:
   - **Package allowlist**: packages outside the list are removed on sync
   - **Profile locking**: users cannot switch to a non-approved profile
   - **Update control**: system updates are approved by admin

### Policy Definition

```
# /etc/erch/policy.toml (desired format — not yet implemented)

[managed]
enabled = true

[profiles]
allowed = ["WORK:Office", "WORK:Dev"]
locked = true

[packages]
blocked = ["steam", "lutris", "heroic"]
pinned = ["crowdstrike-falcon-sensor", "sophos-av"]

[updates]
auto = false
approval_required = true
approval_server = "https://manage.eande.uk/pending"
```

### Managed Mode vs User Freedom

| Aspect | Normal Mode | Managed Mode |
|--------|-------------|--------------|
| Install packages | Any | Allowlist-restricted |
| Remove packages | Any | Cannot remove pinned packages |
| Switch profile | Freely | Locked by policy |
| System updates | Any time | Requires approval |
| Enforceable | No | Yes (opt-in) |

Managed mode is **documented as desired state**. The implementation will be
built incrementally as a later feature.
