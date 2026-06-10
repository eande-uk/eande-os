#!/bin/bash
# Bootstrap script for setting up Omarchy dotfiles on a new PC.
# Thin orchestrator — runs Layer Zero prune first, then deploys configs.
#
# Usage:
#   ./setup.sh                    # Full setup (prune + deploy)
#   ./setup.sh --deploy-only      # Only copy configs (skip prune)
#   ./setup.sh --prune-only       # Only prune bloat packages
#   ./setup.sh --help             # Show this message

set -euo pipefail

DOTFILES_DIR="$(cd "$(dirname "$0")/.." && pwd)"

DO_DEPLOY=true
DO_PRUNE=true

for arg in "$@"; do
  case "$arg" in
    --deploy-only) DO_PRUNE=false ;;
    --prune-only) DO_DEPLOY=false ;;
    --help)
      echo "Options:"
      echo "  (no args)     Full setup: prune + deploy"
      echo "  --deploy-only  Only copy configs"
      echo "  --prune-only   Only prune bloat packages"
      echo "  --help        Show this message"
      exit 0
      ;;
    *)
      echo "Unknown option: $arg"
      echo "Use --help for usage."
      exit 1
      ;;
  esac
done

# --- Step 1: Check prerequisites ---

if [[ ! -d "$DOTFILES_DIR/home/.config" ]]; then
  echo "Error: Expected directory structure not found."
  echo "Make sure you're running this from the repo root."
  exit 1
fi

# --- Step 2: Layer Zero — prune bloat ---

if $DO_PRUNE; then
  if [[ -f "$DOTFILES_DIR/layer-zero/layer-zero.sh" ]]; then
    echo ">>> Layer Zero: pruning unwanted packages..."
    bash "$DOTFILES_DIR/layer-zero/layer-zero.sh" --apply
    echo "  Done."
    echo ""
  else
    echo "  Skipping (layer-zero/layer-zero.sh not found)"
  fi
fi

# --- Step 3: Deploy configs ---

if $DO_DEPLOY; then
  echo ">>> Deploying configs..."
  bash "$DOTFILES_DIR/scripts/deploy.sh"
  echo "  Done."
  echo ""
fi

# --- Step 4: Reminders ---

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Next steps:"
echo "  1. Switch theme if desired:"
echo "     make theme/set NAME=\"Catppuccin\""
echo ""
echo "  2. Restart services:"
echo "     omarchy restart waybar"
echo "     hyprctl reload"
echo "     omarchy restart walker"
echo ""
echo "  3. Edit ~/.config/ freely — changes apply immediately."
echo "     To save them to the repo:"
echo "     git pull"
echo "     make commit TYPE=feat SCOPE=<s> DESC=\"<desc>\""
echo ""
echo "  4. Verify git identity:"
echo "     git config --global user.name"
echo "     git config --global user.email"
echo ""
echo "  5. Reboot for all changes to take effect."
