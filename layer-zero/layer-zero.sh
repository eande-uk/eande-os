#!/bin/bash
# Layer Zero — Sync system to desired state using erch commands.
#
# Two-direction sync:
#   1. Install: items in active categories ON allowlist but not installed → install
#   2. Remove: items in active categories NOT on allowlist but installed → remove
#
# All operations dispatch through erch CLI.
#
# Usage:
#   ./layer-zero.sh                     # Interactive sync
#   ./layer-zero.sh --dry-run            # Preview only
#   ./layer-zero.sh --apply              # Skip confirm, sync immediately
#   ./layer-zero.sh --help               # Show this message

set -euo pipefail

LAYER_ZERO_DIR="$(cd "$(dirname "$0")" && pwd)"

DRY_RUN=false
SKIP_CONFIRM=false

for arg in "$@"; do
    case "$arg" in
        --dry-run) DRY_RUN=true ;;
        --apply) SKIP_CONFIRM=true ;;
        --help)
            echo "Layer Zero — Sync system to desired state"
            echo ""
            echo "Usage:"
            echo "  $0                Interactive sync"
            echo "  $0 --dry-run      Preview only"
            echo "  $0 --apply        Skip confirm, sync immediately"
            echo "  $0 --help         Show this message"
            echo ""
            if [[ -f "$LAYER_ZERO_DIR/config.sh" ]]; then
                source "$LAYER_ZERO_DIR/config.sh"
                echo "Active categories:"
                for cat in "${PRUNE_PACMAN_CATEGORIES[@]}"; do [[ -n "$cat" ]] && echo "  * $cat (pacman)"; done
                for cat in "${PRUNE_WEBAPP_CATEGORIES[@]}"; do [[ -n "$cat" ]] && echo "  * $cat (webapp)"; done
                for cat in "${PRUNE_TUI_CATEGORIES[@]}"; do [[ -n "$cat" ]] && echo "  * $cat (tui)"; done
                for cat in "${PRUNE_NPX_CATEGORIES[@]}"; do [[ -n "$cat" ]] && echo "  * $cat (npx)"; done
            fi
            exit 0
            ;;
    esac
done

source "$LAYER_ZERO_DIR/config.sh" 2>/dev/null || { echo "Error: config.sh not found"; exit 1; }
[[ -f "$LAYER_ZERO_DIR/allowlist.txt" ]] || { echo "Error: allowlist.txt not found"; exit 1; }

# ─── Build allowlist set ─────────────────────
declare -A ALLOWLIST
while IFS= read -r line; do
    line="${line%%#*}"; line="${line// /}"; [[ -z "$line" ]] && continue
    ALLOWLIST["$line"]=1
done < "$LAYER_ZERO_DIR/allowlist.txt"

# ─── Handler functions ────────────────────────

# Detection helpers
is_pacman_installed() { erch-pkg-present "$1" &>/dev/null; }
is_webapp_installed() { [[ -f "$HOME/.local/share/applications/$1.desktop" ]]; }
is_tui_installed()    { [[ -f "$HOME/.local/share/applications/$1.desktop" ]]; }
is_npx_installed()    { [[ -f "$HOME/.local/bin/$1" ]]; }

# Install helpers
pacman_install()      { $DRY_RUN && echo "[DRY-RUN] erch-pkg-add $*" || erch-pkg-add "$@"; }
webapp_install()      { local name="$1" url="$2" icon="$3" exec="${4:-}" mime="${5:-}"
                         if $DRY_RUN; then echo "[DRY-RUN] erch-webapp-install $name $url $icon $exec $mime"
                         else
                           if [[ -n "$exec" ]]; then
                             erch-webapp-install "$name" "$url" "$icon" "$exec" "$mime"
                           else
                             erch-webapp-install "$name" "$url" "$icon"
                           fi
                         fi; }
tui_install()         { local name="$1" cmd="$2" style="$3" icon="$4"
                         $DRY_RUN && echo "[DRY-RUN] erch-tui-install $name $cmd $style $icon" ||
                         erch-tui-install "$name" "$cmd" "$style" "$icon"; }
npx_install()         { local cmd="$1" pkg="$2"
                         $DRY_RUN && echo "[DRY-RUN] erch-npx-install $pkg $cmd" ||
                         erch-npx-install "$pkg" "$cmd"; }

# Remove helpers
pacman_remove()       { $DRY_RUN && echo "[DRY-RUN] erch-pkg-drop $*" || erch-pkg-drop "$@"; }
webapp_remove()       { if $DRY_RUN; then echo "[DRY-RUN] erch-webapp-remove $*"
                         else erch-webapp-remove "$@" &>/dev/null || true; fi; }
tui_remove()          { if $DRY_RUN; then echo "[DRY-RUN] erch-tui-remove $*"
                         else erch-tui-remove "$@" &>/dev/null || true; fi; }
npx_remove()          { for cmd in "$@"; do
                           $DRY_RUN && echo "[DRY-RUN] rm -f $HOME/.local/bin/$cmd" ||
                           rm -f "$HOME/.local/bin/$cmd"
                         done; }

# ─── Scan arrays ──────────────────────────────
TO_INSTALL_PACMAN=()
TO_INSTALL_WEBAPP=()
TO_INSTALL_TUI=()
TO_INSTALL_NPX=()
TO_REMOVE_PACMAN=()
TO_REMOVE_WEBAPP=()
TO_REMOVE_TUI=()
TO_REMOVE_NPX=()

# ─── Pacman: install + remove scan ────────────
for cat in "${PRUNE_PACMAN_CATEGORIES[@]:-}"; do
    # Install: items ON allowlist but NOT installed
    mf="$LAYER_ZERO_DIR/bloat/$cat-install.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS= read -r line; do
            line="${line%%#*}"; line="${line// /}"; [[ -z "$line" ]] && continue
            [[ -n "${ALLOWLIST[$line]:-}" ]] || continue
            is_pacman_installed "$line" && continue
            TO_INSTALL_PACMAN+=("$line")
        done < "$mf"
    fi
    # Remove: items NOT on allowlist but installed
    mf="$LAYER_ZERO_DIR/bloat/$cat.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS= read -r line; do
            line="${line%%#*}"; line="${line// /}"; [[ -z "$line" ]] && continue
            [[ -z "${ALLOWLIST[$line]:-}" ]] || continue
            is_pacman_installed "$line" || continue
            TO_REMOVE_PACMAN+=("$line")
        done < "$mf"
    fi
done

# ─── Webapp: install + remove scan ────────────
for cat in "${PRUNE_WEBAPP_CATEGORIES[@]:-}"; do
    # Install: pipe-delimited (name|url|icon|custom-exec|mime-types)
    mf="$LAYER_ZERO_DIR/bloat/$cat-install.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS='|' read -r name url icon exec mime; do
            [[ "$name" =~ ^#.*$ || -z "$name" ]] && continue
            raw_name="$name"
            name_clean="${name// /}"
            [[ -n "${ALLOWLIST[$name_clean]:-}" ]] || continue
            is_webapp_installed "$raw_name" && continue
            TO_INSTALL_WEBAPP+=("$raw_name|$url|$icon|${exec:-}|${mime:-}")
        done < "$mf"
    fi
    # Remove: names NOT on allowlist, desktop file exists
    mf="$LAYER_ZERO_DIR/bloat/$cat.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS= read -r line; do
            line="${line%%#*}"; line="${line// /}"; [[ -z "$line" ]] && continue
            [[ -z "${ALLOWLIST[$line]:-}" ]] || continue
            is_webapp_installed "$line" || continue
            TO_REMOVE_WEBAPP+=("$line")
        done < "$mf"
    fi
done

# ─── TUI: install + remove scan ───────────────
for cat in "${PRUNE_TUI_CATEGORIES[@]:-}"; do
    # Install: pipe-delimited (name|command|window-style|icon-url)
    mf="$LAYER_ZERO_DIR/bloat/$cat-install.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS='|' read -r name cmd style icon; do
            [[ "$name" =~ ^#.*$ || -z "$name" ]] && continue
            raw_name="$name"
            name_clean="${name// /}"
            [[ -n "${ALLOWLIST[$name_clean]:-}" ]] || continue
            is_tui_installed "$raw_name" && continue
            TO_INSTALL_TUI+=("$raw_name|$cmd|$style|$icon")
        done < "$mf"
    fi
    # Remove: names NOT on allowlist, desktop file exists
    mf="$LAYER_ZERO_DIR/bloat/$cat.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS= read -r line; do
            line="${line%%#*}"; line="${line// /}"; [[ -z "$line" ]] && continue
            [[ -z "${ALLOWLIST[$line]:-}" ]] || continue
            is_tui_installed "$line" || continue
            TO_REMOVE_TUI+=("$line")
        done < "$mf"
    fi
done

# ─── NPX: install + remove scan ───────────────
for cat in "${PRUNE_NPX_CATEGORIES[@]:-}"; do
    # Install: pipe-delimited (command|npm-package)
    mf="$LAYER_ZERO_DIR/bloat/$cat-install.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS='|' read -r cmd pkg; do
            [[ "$cmd" =~ ^#.*$ || -z "$cmd" ]] && continue
            cmd_clean="${cmd// /}"
            [[ -n "${ALLOWLIST[$cmd_clean]:-}" ]] || continue
            is_npx_installed "$cmd_clean" && continue
            TO_INSTALL_NPX+=("$cmd_clean|$pkg")
        done < "$mf"
    fi
    # Remove: names NOT on allowlist, binary exists
    mf="$LAYER_ZERO_DIR/bloat/$cat.pkgs"
    if [[ -f "$mf" ]]; then
        while IFS= read -r line; do
            line="${line%%#*}"; line="${line// /}"; [[ -z "$line" ]] && continue
            [[ -z "${ALLOWLIST[$line]:-}" ]] || continue
            is_npx_installed "$line" || continue
            TO_REMOVE_NPX+=("$line")
        done < "$mf"
    fi
done

# ─── Summary ──────────────────────────────────
TOTAL_INSTALL=$(( ${#TO_INSTALL_PACMAN[@]} + ${#TO_INSTALL_WEBAPP[@]} + ${#TO_INSTALL_TUI[@]} + ${#TO_INSTALL_NPX[@]} ))
TOTAL_REMOVE=$(( ${#TO_REMOVE_PACMAN[@]} + ${#TO_REMOVE_WEBAPP[@]} + ${#TO_REMOVE_TUI[@]} + ${#TO_REMOVE_NPX[@]} ))

if [[ $TOTAL_INSTALL -eq 0 && $TOTAL_REMOVE -eq 0 ]]; then
    echo "Nothing to do - system is already in desired state."
    exit 0
fi

echo "Layer Zero - Sync plan:"
echo ""

if [[ $TOTAL_INSTALL -gt 0 ]]; then
    echo "  Install ($TOTAL_INSTALL):"
    if [[ ${#TO_INSTALL_PACMAN[@]} -gt 0 ]]; then
        echo "    Pacman packages (${#TO_INSTALL_PACMAN[@]}):"
        for pkg in "${TO_INSTALL_PACMAN[@]}"; do echo "      - $pkg"; done
    fi
    if [[ ${#TO_INSTALL_WEBAPP[@]} -gt 0 ]]; then
        echo "    Webapp launchers (${#TO_INSTALL_WEBAPP[@]}):"
        for entry in "${TO_INSTALL_WEBAPP[@]}"; do
            IFS='|' read -r name _ <<< "$entry"
            echo "      - $name"
        done
    fi
    if [[ ${#TO_INSTALL_TUI[@]} -gt 0 ]]; then
        echo "    TUI launchers (${#TO_INSTALL_TUI[@]}):"
        for entry in "${TO_INSTALL_TUI[@]}"; do
            IFS='|' read -r name _ <<< "$entry"
            echo "      - $name"
        done
    fi
    if [[ ${#TO_INSTALL_NPX[@]} -gt 0 ]]; then
        echo "    NPX wrappers (${#TO_INSTALL_NPX[@]}):"
        for entry in "${TO_INSTALL_NPX[@]}"; do
            IFS='|' read -r cmd _ <<< "$entry"
            echo "      - $cmd"
        done
    fi
    echo ""
fi

if [[ $TOTAL_REMOVE -gt 0 ]]; then
    echo "  Remove ($TOTAL_REMOVE):"
    if [[ ${#TO_REMOVE_PACMAN[@]} -gt 0 ]]; then
        echo "    Pacman packages (${#TO_REMOVE_PACMAN[@]}):"
        for pkg in "${TO_REMOVE_PACMAN[@]}"; do echo "      - $pkg"; done
    fi
    if [[ ${#TO_REMOVE_WEBAPP[@]} -gt 0 ]]; then
        echo "    Webapp launchers (${#TO_REMOVE_WEBAPP[@]}):"
        for name in "${TO_REMOVE_WEBAPP[@]}"; do echo "      - $name"; done
    fi
    if [[ ${#TO_REMOVE_TUI[@]} -gt 0 ]]; then
        echo "    TUI launchers (${#TO_REMOVE_TUI[@]}):"
        for name in "${TO_REMOVE_TUI[@]}"; do echo "      - $name"; done
    fi
    if [[ ${#TO_REMOVE_NPX[@]} -gt 0 ]]; then
        echo "    NPX wrappers (${#TO_REMOVE_NPX[@]}):"
        for cmd in "${TO_REMOVE_NPX[@]}"; do echo "      - $cmd"; done
    fi
    echo ""
fi

# ─── Dry run ──────────────────────────────────
if $DRY_RUN; then
    echo "[DRY RUN] Would sync $TOTAL_INSTALL install(s) and $TOTAL_REMOVE removal(s)."
    exit 0
fi

# ─── Confirm ──────────────────────────────────
if ! $SKIP_CONFIRM; then
    echo "These changes will be applied."
    read -rp "Continue? [y/N] " reply
    case "$reply" in
        [yY]|[yY][eE][sS]) ;;
        *) echo "Cancelled."; exit 0 ;;
    esac
fi

# ─── Execute installs ─────────────────────────
INSTALLED=0
REMOVED=0

if [[ ${#TO_INSTALL_PACMAN[@]} -gt 0 ]]; then
    echo "Installing pacman packages..."
    if pacman_install "${TO_INSTALL_PACMAN[@]}"; then
        echo "  + Installed ${#TO_INSTALL_PACMAN[@]} package(s)"
        INSTALLED=$((INSTALLED + ${#TO_INSTALL_PACMAN[@]}))
    else
        echo "  Warning: some packages could not be installed"
    fi
fi

if [[ ${#TO_INSTALL_WEBAPP[@]} -gt 0 ]]; then
    echo "Installing webapp launchers..."
    for entry in "${TO_INSTALL_WEBAPP[@]}"; do
        IFS='|' read -r name url icon exec mime <<< "$entry"
        if webapp_install "$name" "$url" "$icon" "$exec" "$mime"; then
            echo "  + Installed $name"
            INSTALLED=$((INSTALLED + 1))
        else
            echo "  Warning: could not install webapp $name"
        fi
    done
fi

if [[ ${#TO_INSTALL_TUI[@]} -gt 0 ]]; then
    echo "Installing TUI launchers..."
    for entry in "${TO_INSTALL_TUI[@]}"; do
        IFS='|' read -r name cmd style icon <<< "$entry"
        if tui_install "$name" "$cmd" "$style" "$icon"; then
            echo "  + Installed $name"
            INSTALLED=$((INSTALLED + 1))
        else
            echo "  Warning: could not install TUI $name"
        fi
    done
fi

if [[ ${#TO_INSTALL_NPX[@]} -gt 0 ]]; then
    echo "Installing NPX wrappers..."
    for entry in "${TO_INSTALL_NPX[@]}"; do
        IFS='|' read -r cmd pkg <<< "$entry"
        if npx_install "$cmd" "$pkg"; then
            echo "  + Installed $cmd"
            INSTALLED=$((INSTALLED + 1))
        else
            echo "  Warning: could not install NPX $cmd"
        fi
    done
fi

# ─── Execute removals ─────────────────────────

if [[ ${#TO_REMOVE_PACMAN[@]} -gt 0 ]]; then
    echo "Removing pacman packages..."
    if pacman_remove "${TO_REMOVE_PACMAN[@]}"; then
        echo "  - Removed ${#TO_REMOVE_PACMAN[@]} package(s)"
        REMOVED=$((REMOVED + ${#TO_REMOVE_PACMAN[@]}))
    else
        echo "  Warning: some packages could not be removed"
    fi
fi

if [[ ${#TO_REMOVE_WEBAPP[@]} -gt 0 ]]; then
    echo "Removing webapp launchers..."
    if webapp_remove "${TO_REMOVE_WEBAPP[@]}"; then
        echo "  - Removed ${#TO_REMOVE_WEBAPP[@]} webapp(s)"
        REMOVED=$((REMOVED + ${#TO_REMOVE_WEBAPP[@]}))
    else
        echo "  Warning: some webapps could not be removed"
    fi
fi

if [[ ${#TO_REMOVE_TUI[@]} -gt 0 ]]; then
    echo "Removing TUI launchers..."
    if tui_remove "${TO_REMOVE_TUI[@]}"; then
        echo "  - Removed ${#TO_REMOVE_TUI[@]} TUI(s)"
        REMOVED=$((REMOVED + ${#TO_REMOVE_TUI[@]}))
    else
        echo "  Warning: some TUIs could not be removed"
    fi
fi

if [[ ${#TO_REMOVE_NPX[@]} -gt 0 ]]; then
    echo "Removing NPX wrappers..."
    for cmd in "${TO_REMOVE_NPX[@]}"; do
        if npx_remove "$cmd"; then
            echo "  - Removed $cmd"
            REMOVED=$((REMOVED + 1))
        fi
    done
fi

# ─── Post-removal theme re-sync ───────────────
if command -v erch &>/dev/null; then
    CURRENT_THEME=$(erch theme current 2>/dev/null || true)
    if [[ -n "$CURRENT_THEME" ]]; then
        echo "Re-syncing theme: $CURRENT_THEME"
        erch theme set "$CURRENT_THEME" &>/dev/null || true
    fi
fi

echo ""
echo "Layer Zero - Sync complete ($INSTALLED installed, $REMOVED removed)."
