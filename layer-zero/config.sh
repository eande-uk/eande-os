# Layer Zero — Bloat categories to prune
#
# Each array maps to a package manager/handler type.
# All operations dispatch through omarchy commands:
#   PRUNE_PACMAN_CATEGORIES  → omarchy pkg add | omarchy pkg drop
#   PRUNE_WEBAPP_CATEGORIES  → omarchy webapp install | omarchy webapp remove
#   PRUNE_TUI_CATEGORIES     → omarchy tui install | omarchy tui remove
#   PRUNE_NPX_CATEGORIES     → omarchy npx install | rm ~/.local/bin/<cmd>
#
# Uncomment a category to enable two-direction sync.
# Items ON allowlist → ensured installed. Items NOT on allowlist → removed.

# ── Pacman packages ──────────────────────────
PRUNE_PACMAN_CATEGORIES=(
    gaming
    # media
    # office
    # communication
    # browsers
    # runtimes
    # terminals
)

# ── Omarchy webapp launchers ─────────────────
PRUNE_WEBAPP_CATEGORIES=(
    webapps
)

# ── Omarchy TUI launchers ────────────────────
PRUNE_TUI_CATEGORIES=(
    # tui
)

# ── Omarchy NPX wrappers ─────────────────────
PRUNE_NPX_CATEGORIES=(
    # npx
)
