# Layer Zero — Profiles

Profiles define what gets **pre-installed** on an E&E OS system. They are
shared across all three targets (erch, upstream omarchy, Arch+Hyprland).

## Profile Catalog

### WORK:Office

Non-developer employees. Office productivity, communication, and basic tools.

| Category | Packages |
|----------|---------|
| Office | libreoffice-fresh, evince, xournalpp |
| Communication | thunderbird, signal-desktop, zoom |
| Browser | chromium, firefox |
| Media | mpv, imv, gnome-calculator |
| Printing | cups, cups-browsed, system-config-printer |
| Cloud | nextcloud-client, gvfs (mtp, nfs, smb) |

### WORK:Dev

Software developers. All runtimes, containers, and development tooling.

| Category | Packages |
|----------|---------|
| Runtimes | python, nodejs, go, rust, dotnet-runtime-9.0 |
| Containers | docker, docker-compose, docker-buildx |
| Languages | clang, llvm, ruby, tree-sitter-cli |
| Editor | neovim, vscode (via erch webapp) |
| Tools | git, gh, lazygit, lazydocker, mise, jq, yq |
| Databases | postgresql-libs, mariadb-libs, sqlite |

### WORK:AI/ML

AI/ML engineers and data scientists. GPU compute, ML frameworks, data tools.

| Category | Packages |
|----------|---------|
| GPU/Compute | cuda, cudnn, tensorrt, opencl-nvidia |
| ML Frameworks | pytorch, tensorflow, jupyter-lab |
| Data | python-numpy, python-pandas, python-scipy, python-scikit-learn |
| Visualization | python-matplotlib, python-seaborn, kdenlive |
| Tools | ollama, llama.cpp, stable-diffusion (via webapp) |

### EDUCATION:School

Primary and secondary education. Learning software, content filtering.

| Category | Packages |
|----------|---------|
| Learning | gcompris-qt, ktouch, kstars, marble |
| Office | libreoffice-fresh (educational templates) |
| Browser | chromium (with content filtering) |
| Creativity | krita, inkscape, blender |
| Printing | cups, cups-browsed |

### EDUCATION:Uni

University students and researchers. Scientific computing, writing tools.

| Category | Packages |
|----------|---------|
| Writing | texlive-most, pandoc, zotero, jabref |
| Research | r, rstudio, octave, matlab (via webapp) |
| Scientific | python-numpy, python-pandas, python-matplotlib |
| Data | postgresql, sqlite, sqlitebrowser |
| Reference | evince, xournalpp, obsidian |
| Printing | cups, cups-pdf |

### GAME

Gaming-focused setup. Game stores, emulators, gaming optimizations.

| Category | Packages |
|----------|---------|
| Stores | steam, heroic-games-launcher, lutris |
| Emulation | retroarch, dolphin-emu, pcsx2, yuzu (if available) |
| Cloud/Remote | moonlight-qt, geforce-now (via webapp) |
| Optimization | gamemode, mangohud, goverlay, proton-ge |
| Wine | wine, winetricks, protonup-qt |
| Recording | obs-studio, gpu-screen-recorder |

## Profile Composition

Profiles are additive. Selecting `WORK:Dev + GAME` installs the union of
both package sets.

All profiles inherit the **base allowlist** (core desktop, shell, terminal,
fonts, git, Docker, etc.).

## Package Categories

Packages are organized in category files under `layer-zero/bloat/`:

```
layer-zero/bloat/
├── browsers.pkgs              # Items to prune (not on allowlist)
├── browsers-install.pkgs      # Items to install (on allowlist)
├── communication.pkgs
├── communication-install.pkgs
├── gaming.pkgs
├── gaming-install.pkgs
├── media.pkgs
├── media-install.pkgs
├── npx.pkgs
├── npx-install.pkgs
├── office.pkgs
├── office-install.pkgs
├── runtimes.pkgs
├── runtimes-install.pkgs
├── terminals.pkgs
├── terminals-install.pkgs
├── tui.pkgs
├── tui-install.pkgs
├── webapps.pkgs
└── webapps-install.pkgs
```

Each category has a dual-file pattern:

| File | Purpose |
|------|---------|
| `{category}.pkgs` | Items to **remove** — pruned if not on the allowlist |
| `{category}-install.pkgs` | Items to **install** — ensured present if on the allowlist |

The two-direction sync uses both files in tandem: items ON the allowlist in
`-install.pkgs` get installed; items NOT on the allowlist in `.pkgs` get removed.
Profile definitions reference these categories and add specific packages
within each.

## User Freedom

Profiles are a **one-time initial state**. After install:

- Users install any package freely
- Users remove any pre-installed package
- No restrictions in normal mode

The two-direction sync (layer-zero.sh) only prunes packages in **active
bloat categories** that are not on the allowlist. Packages outside active
categories are never touched.

## Managed Mode

For organizations that need control, an opt-in managed mode extends the
profile system with:

- **Package pinning**: certain packages cannot be removed
- **Profile locking**: users cannot change to a non-approved profile
- **Blocklist**: certain packages cannot be installed
- **Approval queue**: updates require admin approval

See [erch/PROFILES.md](../erch/PROFILES.md#managed-mode-organizations) for
details. Managed mode is documented as desired state and will be built
incrementally.
