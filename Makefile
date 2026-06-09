DOTFILES_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: help init setup deploy deploy/restock deploy/dry-run status \
        adopt \
        layer-zero layer-zero/apply layer-zero/dry-run \
        theme/list theme/set test test/quiet \
        diff log commit branch/create pr

help:
	@echo "E&E UK — Omarchy Dotfiles"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "── Lifecycle ──"
	@echo "  init              Create branch user/$$USER from master + deploy"
	@echo "  setup             Full bootstrap: init + layer-zero sync"
	@echo ""
	@echo "── Deploy (symlinks: repo ↔ $$HOME) ──"
	@echo "  deploy            Link configs via stow (with backup, errors on master)"
	@echo "  deploy/dry-run    Preview what deploy would change (stow -n -v)"
	@echo "  deploy/restock    Re-apply master defaults (checkout → deploy --force → return)"
	@echo "  adopt             Adopt existing ~/.config/ as your branch defaults"
	@echo ""
	@echo "── Inspect ──"
	@echo "  status            Show branch, uncommitted changes, stow state"
	@echo ""
	@echo "── Layer 0: System state ──"
	@echo "  layer-zero              Interactive two-direction sync"
	@echo "  layer-zero/apply        Apply without confirm"
	@echo "  layer-zero/dry-run      Preview only"
	@echo ""
	@echo "── Layer 2: Rebranding ──"
	@echo "  theme/list        omarchy theme list"
	@echo "  theme/set NAME=n  omarchy theme set"
	@echo ""
	@echo "── Layer 4: Tests ──"
	@echo "  test              Run verification tests (verbose)"
	@echo "  test/quiet        Run verification tests (quiet)"
	@echo ""
	@echo "── Git / Commit ──"
	@echo "  diff              Show uncommitted changes"
	@echo "  log               Recent commits"
	@echo "  commit TYPE=t SCOPE=s DESC=d  Stage all + commit with convention"
	@echo ""
	@echo "── erch (omarchy fork) ──"
	@echo "  erch/init         Init submodule + deploy to ~/.local/share/omarchy/"
	@echo ""
	@echo "── Contributing ──"
	@echo "  branch/create     Create user/$$USER branch from master"
	@echo "  pr                Open PR from current branch → master"
	@echo ""

init:
	$(MAKE) branch/create
	./scripts/deploy.sh --adopt

setup:
	$(MAKE) init
	$(MAKE) layer-zero/apply

deploy:
	./scripts/deploy.sh

deploy/dry-run:
	./scripts/deploy.sh --dry-run

deploy/restock:
	@CURRENT_BRANCH=$$(git rev-parse --abbrev-ref HEAD); \
	git stash 2>/dev/null || true; \
	git checkout master; \
	./scripts/deploy.sh --force; \
	git checkout "$$CURRENT_BRANCH"; \
	git stash pop 2>/dev/null || true; \
	echo "Restocked from master, returned to $$CURRENT_BRANCH"

adopt:
	@BRANCH=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$BRANCH" = "main" ] || [ "$$BRANCH" = "master" ]; then \
		$(MAKE) branch/create; \
	fi; \
	./scripts/deploy.sh --adopt

status:
	@echo "=== Status ==="
	@echo "Branch: $$(git rev-parse --abbrev-ref HEAD)"
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" = "main" ] || [ "$$(git rev-parse --abbrev-ref HEAD)" = "master" ]; then \
		echo "  ⚠️  ON ROOT BRANCH — create a user branch: make init"; \
	fi
	@echo "erch:  $$(git -C erch rev-parse --short HEAD 2>/dev/null || echo '(not initialized)')"
	@echo "Theme: $$(omarchy theme current 2>/dev/null || echo '(omarchy not available)')"
	@echo ""
	@echo "Uncommitted changes (edit ~/.config/ = edit repo via symlinks):"
	@git status --short
	@echo ""
	@echo "Stow check (run 'make deploy' to fix):"
	@stow --no-folding -t $$HOME -n -v home 2>&1 || true

layer-zero:
	./layer-zero/layer-zero.sh

layer-zero/apply:
	./layer-zero/layer-zero.sh --apply

layer-zero/dry-run:
	./layer-zero/layer-zero.sh --dry-run

theme/list:
	omarchy theme list

theme/set:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make theme/set NAME=\"Theme Name\""; \
		exit 1; \
	fi
	omarchy theme set "$(NAME)"

test:
	cd tests && go test ./... -v -count=1

test/quiet:
	cd tests && go test ./... -count=1

diff:
	git diff

log:
	git log --oneline -15

commit:
	@if [ -z "$(TYPE)" ] || [ -z "$(DESC)" ]; then \
		echo "Usage: make commit TYPE=<type> SCOPE=<scope> DESC=\"<description>\""; \
		echo ""; \
		echo "  TYPE: feat|fix|docs|refactor|reconcile|chore|test"; \
		echo "  SCOPE: optional (e.g. tiling, waybar, layer-zero)"; \
		echo "  DESC: required, imperative, no period"; \
		exit 1; \
	fi
	git add -A
	git commit -m "$(TYPE)$(if $(SCOPE),($(SCOPE))): $(DESC)"

branch/create:
	@BRANCH="user/$$USER"; \
	if git show-ref --verify --quiet "refs/heads/$$BRANCH"; then \
		echo "Branch $$BRANCH already exists on $$(git rev-parse --abbrev-ref HEAD)"; \
	else \
		git checkout -b "$$BRANCH" master; \
		echo "Created branch $$BRANCH from master"; \
	fi

erch/init:
	git submodule update --init erch/
	@echo "erch submodule initialized. Run erch/setup.sh to deploy."

pr:
	@CURRENT=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$CURRENT" = "master" ]; then \
		echo "ERROR: On master. Switch to your user branch first."; \
		exit 1; \
	fi; \
	gh pr create --base master --head "$$CURRENT" --fill
