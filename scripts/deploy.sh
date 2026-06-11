#!/bin/bash
# deploy.sh — Copy repo configs to $HOME (one-way push)
# Replaces stow-erch.sh. No symlinks — plain copy with backup.
# Requires: gum

set -euo pipefail

DOTFILES_DIR="$(cd "$(dirname "$0")/../dotfiles" && pwd)"
STOW_TARGET="$HOME"
BACKUP_BASE="$HOME/.config/erch/backups"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="$BACKUP_BASE/$TIMESTAMP"

# --- Parse flags ---
ADOPT=false
FORCE=false
DRY_RUN=false
for arg in "$@"; do
  case "$arg" in
    --adopt) ADOPT=true ;;
    --force) FORCE=true ;;
    --dry-run) DRY_RUN=true ;;
    --help)
      echo "Usage: $0 [--adopt] [--force] [--dry-run]"
      echo "  (none)    Link configs via stow (errors on master branch)"
      echo "  --adopt   Adopt existing ~/.config/ files into repo"
      echo "  --force   Bypass master branch guard + skip backup prompt"
      echo "  --dry-run Preview what would change (stow -n -v)"
      exit 0
      ;;
  esac
done

# --- Branch guard (skip for dry-run — preview is safe) ---
if [ "$DRY_RUN" = false ]; then
  BRANCH=$(git -C "$DOTFILES_DIR" rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
  if { [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; } && [ "$FORCE" = false ]; then
    echo "ERROR: On $BRANCH branch. Create a user branch first:"
    echo "  make init"
    echo "Or use --force to deploy (for restock)."
    exit 1
  fi
fi

# --- Prerequisites ---
if ! command -v stow &>/dev/null; then
  echo "stow is required. Install: sudo pacman -S stow"
  exit 1
fi

if ! command -v gum &>/dev/null; then
  echo "gum is required. Install: sudo pacman -S gum"
  exit 1
fi

if [[ ! -d "$DOTFILES_DIR/home/.config" ]]; then
  gum style --border normal --padding "1 2" --foreground "#BF616A" \
    "Invalid repo structure." \
    "Expected home/.config/ directory not found in:" "$DOTFILES_DIR"
  exit 1
fi

cd "$DOTFILES_DIR"

# --- Header ---
ADOPT_LABEL=""
if [ "$ADOPT" = true ]; then ADOPT_LABEL=" (adopt mode)"; fi
if [ "$DRY_RUN" = true ]; then ADOPT_LABEL=" (dry run)"; fi
gum style --border double --padding "1 2" --margin "0 0 1 0" \
  "E&E UK — Omarchy Dotfiles" \
  "" \
  "Linking configs to \$HOME via stow$ADOPT_LABEL"

# --- Dry run ---
if [ "$DRY_RUN" = true ]; then
  gum style --foreground "#88C0D0" "Dry run — showing what would change:"
  (cd "$DOTFILES_DIR" && stow --no-folding -t "$STOW_TARGET" -n -v home 2>&1)
  exit $?
fi

# --- Conflict scan ---
gum style --foreground "#88C0D0" "Scanning for existing files…"

declare -a conflicts=()

while IFS= read -r -d '' file; do
  relative="${file#$DOTFILES_DIR/home/}"
  target="$STOW_TARGET/$relative"
  if [[ -f "$target" && ! -L "$target" ]]; then
    conflicts+=("$target")
  fi
done < <(find "$DOTFILES_DIR/home" -type f -print0)

# --- Backup ---
if ((${#conflicts[@]} > 0)); then
  echo
  gum style --foreground "#EBCB8B" "Found ${#conflicts[@]} existing file(s):"
  for f in "${conflicts[@]}"; do
    echo "  • $f"
  done
  echo

  if [[ -t 0 ]] && [ "$FORCE" = false ]; then
    if gum confirm "Backup these files before overwriting?"; then
      DO_BACKUP=true
    else
      gum style --foreground "#BF616A" "Deploy cancelled. Remove conflicting files or use --force."
      exit 1
    fi
  else
    gum style --foreground "#88C0D0" "Non-interactive — auto-backing up and removing conflicts..."
    DO_BACKUP=true
  fi

  if [ "$DO_BACKUP" = true ]; then
    mkdir -p "$BACKUP_DIR"
    for file in "${conflicts[@]}"; do
      relative="${file#$STOW_TARGET/}"
      dest="$BACKUP_DIR/$relative"
      mkdir -p "$(dirname "$dest")"
      cp -a "$file" "$dest"
      rm "$file"
    done
    echo
    gum style --foreground "#A3BE8C" "✓ Backed up to: $BACKUP_DIR"
  fi
else
  echo
  gum style --foreground "#A3BE8C" "✓ No existing files found."
fi

# --- Stow ---
STOW_CMD="stow --no-folding -t $STOW_TARGET"
if [ "$ADOPT" = true ]; then
  STOW_CMD="$STOW_CMD --adopt"
fi

gum spin --spinner dot --title "Linking configs…" -- bash -c "
  cd \"$DOTFILES_DIR\" && $STOW_CMD home
"

gum style --foreground "#A3BE8C" "✓ Configs linked"

# --- Sync custom-branding to erch/branding ---
CUSTOM_BRANDING_SRC="$DOTFILES_DIR/home/.config/custom-branding"
OMARCHY_BRANDING_DST="$HOME/.config/erch/branding"

if [[ -d "$CUSTOM_BRANDING_SRC" ]]; then
  mkdir -p "$OMARCHY_BRANDING_DST"
  gum spin --spinner dot --title "Syncing branding to erch…" -- \
    cp -a "$CUSTOM_BRANDING_SRC/"* "$OMARCHY_BRANDING_DST/"
  gum style --foreground "#A3BE8C" "✓ Branding synced to erch/branding"
fi

# --- Make scripts executable ---
if [[ -d "$HOME/.local/bin" ]]; then
  gum spin --spinner dot --title "Making scripts executable…" -- \
    chmod +x "$HOME/.local/bin/"*
  gum style --foreground "#A3BE8C" "✓ Scripts ready"
fi

# --- Hooks ---
if ls "$HOME/.config/erch/hooks/"* &>/dev/null 2>&1; then
  gum spin --spinner dot --title "Making hooks executable…" -- \
    chmod +x "$HOME/.config/erch/hooks/"*
  gum style --foreground "#A3BE8C" "✓ Hooks ready"
fi

# --- Hide stock themes (repeatable: survives erch update) ---
STOCK_THEMES_SRC="$HOME/.local/share/erch/themes"
STOCK_THEMES_DST="$HOME/.config/erch/stock-themes"

if [[ -d "$STOCK_THEMES_SRC" ]]; then
  mkdir -p "$STOCK_THEMES_DST"
  for theme_dir in "$STOCK_THEMES_SRC"/*/; do
    theme="$(basename "$theme_dir")"
    if [[ -d "$theme_dir" && ! -d "$HOME/.config/erch/themes/$theme" ]]; then
      rm -rf "$STOCK_THEMES_DST/$theme"
      mv "$theme_dir" "$STOCK_THEMES_DST/"
    fi
  done
  gum style --foreground "#A3BE8C" "✓ Stock themes hidden (custom themes stay visible)"
fi



# --- Summary ---
echo
gum style --border normal --padding "1 2" "$(
  gum join --vertical \
    "$(gum style --bold "✓ Deployed")" \
    "" \
    "  Configs:  linked via stow from $DOTFILES_DIR/home/" \
    "  Vim Mode:  in menu Trigger → Toggle and via SUPER+SHIFT+V" \
    "  erch:      erch fork at erch/ submodule" \
    "  Hooks:     theme-set, font-set, post-update" \
    "  Backup:   $BACKUP_DIR" \
    "" \
    "$(gum style --foreground "#81A1C1" "  Symlinks: editing ~/.config/ edits the repo.")" \
    "$(gum style --foreground "#81A1C1" "  Always be on a user branch (never master).")" \
    "" \
    "$(gum style --foreground "#81A1C1" "  To save changes:")" \
    "  • git diff (changes appear automatically)" \
    "  • make commit TYPE=feat SCOPE=<s> DESC=\"<d>\"" \
    "  • make pr"
)"
