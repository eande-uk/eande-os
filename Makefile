REPO_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: help init setup deploy status \
        erch/init E-OS/init \
        test test/quiet \
        diff log commit branch/create pr

help:
	@echo "E&E OS — Distro Hub"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "── Lifecycle ──"
	@echo "  init              Create branch user/$$USER from master + init submodules"
	@echo "  setup             Full bootstrap: init + erch deploy"
	@echo ""
	@echo "── Deploy ──"
	@echo "  deploy            Deploy erch to ~/.local/share/erch/"
	@echo ""
	@echo "── Inspect ──"
	@echo "  status            Show branch, submodules, uncommitted changes"
	@echo ""
	@echo "── Submodules ──"
	@echo "  erch/init         Init erch submodule"
	@echo "  E-OS/init         Init E-OS submodule"
	@echo ""
	@echo "── Tests ──"
	@echo "  test              Run verification tests (verbose)"
	@echo "  test/quiet        Run verification tests (quiet)"
	@echo ""
	@echo "── Git / Commit ──"
	@echo "  diff              Show uncommitted changes"
	@echo "  log               Recent commits (15)"
	@echo "  commit TYPE=t SCOPE=s DESC=d  Stage all + commit with convention"
	@echo ""
	@echo "── Contributing ──"
	@echo "  branch/create     Create user/$$USER branch from master"
	@echo "  pr                Open PR from current branch → master"
	@echo ""

init:
	$(MAKE) branch/create
	$(MAKE) erch/init
	$(MAKE) E-OS/init

setup:
	$(MAKE) init
	$(MAKE) deploy

deploy:
	@if [ ! -d "erch/.git" ]; then \
		echo "erch submodule not initialized. Run: make erch/init"; \
		exit 1; \
	fi
	@echo "Deploying erch..."
	cd erch && ./install.sh

status:
	@echo "=== E&E OS Hub Status ==="
	@echo "Branch: $$(git rev-parse --abbrev-ref HEAD)"
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" = "main" ] || [ "$$(git rev-parse --abbrev-ref HEAD)" = "master" ]; then \
		echo "  WARNING: ON ROOT BRANCH — create a user branch: make init"; \
	fi
	@echo ""
	@echo "Submodules:"
	@echo "  erch:    $$(git -C erch rev-parse --short HEAD 2>/dev/null || echo '(not initialized)')"
	@echo "  E-OS:    $$(git -C E-OS rev-parse --short HEAD 2>/dev/null || echo '(not initialized)')"
	@echo "  E-OS-AI: (planned — repo not yet created)"
	@echo ""
	@echo "Uncommitted changes:"
	@git status --short

erch/init:
	git submodule update --init erch/
	@echo "erch submodule initialized."

E-OS/init:
	git submodule update --init E-OS/
	@echo "E-OS submodule initialized."

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
		echo "  SCOPE: optional (e.g. erch, e-os, e-os-ai)"; \
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

pr:
	@CURRENT=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$CURRENT" = "master" ]; then \
		echo "ERROR: On master. Switch to your user branch first."; \
		exit 1; \
	fi; \
	gh pr create --base master --head "$$CURRENT" --fill
