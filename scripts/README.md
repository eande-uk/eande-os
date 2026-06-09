# Scripts

## Deployment scripts

These stay in the repo and are run from `scripts/`:

### `deploy.sh`

Stow-based deploy for fresh installs or bulk updates. Creates symlinks from
`$HOME` back to the repo — editing `~/.config/` auto-syncs to your user branch.

```bash
./scripts/deploy.sh                  # Deploy (errors on master)
./scripts/deploy.sh --adopt          # Adopt existing configs into repo
./scripts/deploy.sh --force          # Bypass master guard (restock)
./scripts/deploy.sh --dry-run        # Preview only
./scripts/deploy.sh --help           # Show usage
```

Steps: branch guard → conflict scan → backup (prompted) →
`stow --no-folding -t $HOME home` → branding sync → make hooks executable →
hide stock themes → set default theme → make scripts executable.

### `setup.sh`

New-machine bootstrap orchestrator. Runs Layer Zero sync first (if config.sh
has active categories), then delegates to `deploy.sh`.

```bash
./scripts/setup.sh                  # Full setup (prune + deploy)
./scripts/setup.sh --deploy-only    # Only link configs (skip prune)
./scripts/setup.sh --prune-only     # Only prune bloat (skip deploy)
./scripts/setup.sh --help           # Show usage
```

### `../layer-zero/layer-zero.sh`

Layer Zero — two-direction sync (install + remove). Lives under `layer-zero/`, not `scripts/`.
See `dotfiles/README.md` for full docs.

```bash
./layer-zero/layer-zero.sh              # Interactive sync
./layer-zero/layer-zero.sh --dry-run    # Preview only
./layer-zero/layer-zero.sh --apply      # Skip confirm, sync immediately
```

## User-facing scripts

These are deployed to `~/.local/bin/` via `make deploy` and available on PATH.
Edit them at `dotfiles/home/.local/bin/<script>`.

| Script | Purpose |
|--------|---------|
| `fd-rename` | Batch file renamer using `fd`. Dry-run by default; `--confirm` to apply. |
| `ddcutil-brightness` | External monitor brightness via DDC/CI. Detects focused monitor. |
| `ddcutil-source` | Interactive external monitor input source picker via DDC/CI. |
| `md-to-html` | Markdown → standalone HTML with dark GitHub CSS via pandoc. |
| `mmd` | Mermaid diagram → SVG/PNG via Docker (`minlag/mermaid-cli`). |
| `monitor-scaling-toggle` | Cycle display scaling 1.00 → 2.00 → 3.13 across all monitors. |
| `system-idle` | Lock + DPMS off after 10s delay. Called by hypridle. |
| `system-idle-resume` | Restore displays and backlight after idle. |
