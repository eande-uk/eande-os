# Layer Zero — Bloat categories to prune
#
# Each array maps to a package manager/handler type.
# All operations dispatch through erch commands:
#   PRUNE_PACMAN_CATEGORIES  → erch pkg add | erch pkg drop
#   PRUNE_WEBAPP_CATEGORIES  → erch webapp install | erch webapp remove
#   PRUNE_TUI_CATEGORIES     → erch tui install | erch tui remove
#   PRUNE_NPX_CATEGORIES     → erch npx install | rm ~/.local/bin/<cmd>
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
