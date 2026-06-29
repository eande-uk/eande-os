# E-OS — E&E OS

Simpler Arch+Hyprland experience. Planned distro.

## Status

**Planned** — repo to be created at `github.com/eande-uk/e-os`.

## Purpose

A standalone Linux distribution providing the E&E desktop experience
on Arch+Hyprland without the full erch install pipeline.

## Architecture

Follows the shared L0-L4 layer system:

- L0: Simplified install pipeline
- L1: Shell, Hyprland, env defaults
- L2: Core application configs
- L3: Curated theme subset
- L4: Essential CLI commands

## Setup

Once the repo is created:

```bash
git submodule add git@github.com:eande-uk/e-os.git E-OS
```

Then from the hub:

```bash
make E-OS/init
```
