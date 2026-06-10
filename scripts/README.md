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
See `docs/layer-zero/OVERVIEW.md` for full docs.

```bash
./layer-zero/layer-zero.sh              # Interactive sync
./layer-zero/layer-zero.sh --dry-run    # Preview only
./layer-zero/layer-zero.sh --apply      # Skip confirm, sync immediately
```
