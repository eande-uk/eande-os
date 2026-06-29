# E-OS-AI — E&E OS AI

Agent-focused minimal OS. Planned distro.

## Status

**Planned** — repo to be created at `github.com/eande-uk/e-os-ai`.

## Purpose

A minimal, predictable Linux environment designed for autonomous agents.
Provides a clean base with agent-first defaults.

## Architecture

Follows the shared L0-L4 layer system:

- L0: Minimal install pipeline
- L1: Agent-first defaults (no desktop, headless)
- L2: Minimal configs (SSH, tmux, basic tools)
- L3: No theme (headless)
- L4: Agent hooks only

## Setup

Once the repo is created:

```bash
git submodule add git@github.com:eande-uk/e-os-ai.git E-OS-AI
```

Then from the hub:

```bash
make E-OS-AI/init
```
