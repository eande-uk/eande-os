# E&E OS — ISO Build Guide

## Overview

E&E OS provides custom ISOs for each distro using archiso. Each ISO boots into a TUI installer that partitions the disk, installs Arch, clones the distro repo, and runs the distro's installer.

## Installation Methods

### Method 1: Custom ISO (Recommended)

Download the pre-built ISO and write to USB:

```bash
# Write ISO to USB
sudo dd if=iso/out/erch-2026.06.30-x86_64.iso of=/dev/sdX bs=4M status=progress oflag=sync
```

Boot from USB → TUI installer runs automatically.

### Method 2: Official Arch ISO + boot.sh

Install vanilla Arch from the official ISO, then:

```bash
# For erch:
curl -fsSL https://raw.githubusercontent.com/eande-uk/erch/master/boot.sh | bash

# For E-OS:
curl -fsSL https://raw.githubusercontent.com/eande-uk/e-os/dev/boot.sh | bash
```

## Building ISOs

### Prerequisites

```bash
sudo pacman -S archiso
```

### Build All ISOs

```bash
make iso/build
```

### Build Specific ISOs

```bash
make iso/build/erch          # erch ISO
make iso/build/e-os          # All 4 E-OS ISOs
make iso/build/e-os-console  # E-OS Console only
make iso/build/e-os-school   # E-OS School only
make iso/build/e-os-uni      # E-OS Uni only
make iso/build/e-os-org      # E-OS Org only
```

### Output

ISOs are built to `iso/out/`:

```
iso/out/
├── erch-2026.06.30-x86_64.iso
├── e-os-console-2026.06.30-x86_64.iso
├── e-os-school-2026.06.30-x86_64.iso
├── e-os-uni-2026.06.30-x86_64.iso
└── e-os-org-2026.06.30-x86_64.iso
```

### Clean Build Artifacts

```bash
make iso/clean
```

## Testing ISOs

Test with QEMU before writing to USB:

```bash
make iso/test
```

Or manually:

```bash
# UEFI mode
sudo run_archiso -u -i iso/out/erch-2026.06.30-x86_64.iso

# BIOS mode
sudo run_archiso -i iso/out/erch-2026.06.30-x86_64.iso
```

## ISO Structure

```
iso/
├── erch/                    erch ISO profile
│   ├── profiledef.sh        ISO metadata
│   ├── packages.x86_64      Packages for live environment
│   ├── pacman.conf          Pacman config
│   └── airootfs/
│       ├── etc/systemd/system/getty@tty1.service.d/autologin.conf
│       └── root/installer.sh    TUI installer
├── e-os-console/            E-OS Console ISO
├── e-os-school/             E-OS School ISO
├── e-os-uni/                E-OS Uni ISO
└── e-os-org/                E-OS Org ISO
```

## Installer Flow

1. Boot ISO → Auto-login as root (TTY)
2. TUI installer runs (gum-based)
3. User selects disk, enters username/password, timezone, keymap
4. Installer partitions disk (Btrfs + Limine)
5. pacstrap base system
6. Clone distro repo → run distro installer
7. Reboot into installed system

## Requirements

- UEFI mode (no BIOS support)
- Secure Boot disabled
- x86_64 CPU
- Internet connection during install
- 8GB+ USB drive (for writing ISO)

## Troubleshooting

- **archiso not installed**: `sudo pacman -S archiso`
- **Build fails**: Check `/tmp/archiso-work-*` for errors
- **ISO won't boot**: Verify UEFI mode, Secure Boot disabled
- **Installer fails**: Check `/var/log/installer.log` in live environment
