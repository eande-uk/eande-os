#!/usr/bin/env python3
"""QEMU verify test — mount disk image and verify installed system.

Usage:
    python3 verify_test.py --disk <path-to-disk> [--profile console|school|uni|org]

Exit codes:
    0 - all checks passed
    1 - one or more checks failed
"""

import argparse
import os
import subprocess
import sys
import time


def run(cmd, check=True, capture=True):
    print(f"[*] {' '.join(cmd)}")
    return subprocess.run(cmd, check=check, capture_output=capture, text=True)


def setup_nbd(disk_path):
    run(["sudo", "modprobe", "nbd", "max_part=8"], check=False)
    nbd_dev = "/dev/nbd0"
    run(["sudo", "qemu-nbd", "--disconnect", nbd_dev], check=False)
    time.sleep(1)
    run(["sudo", "qemu-nbd", "--connect", nbd_dev, disk_path])
    time.sleep(2)
    run(["sudo", "partprobe", nbd_dev], check=False)
    time.sleep(1)
    return nbd_dev


def mount_root(nbd_dev, mount_point):
    root_part = f"{nbd_dev}p2"
    if not os.path.exists(root_part):
        root_part = f"{nbd_dev}p3"
    if not os.path.exists(root_part):
        root_part = f"{nbd_dev}p1"

    os.makedirs(mount_point, exist_ok=True)
    run(["sudo", "mount", root_part, mount_point])

    boot_part = f"{nbd_dev}p1"
    boot_mount = os.path.join(mount_point, "boot")
    if os.path.exists(boot_part):
        os.makedirs(boot_mount, exist_ok=True)
        run(["sudo", "mount", boot_part, boot_mount], check=False)

    return root_part


def verify(mount_point, profile):
    checks = []

    def check(name, condition, detail=""):
        status = "PASS" if condition else "FAIL"
        msg = f"[{status}] {name}"
        if detail:
            msg += f" — {detail}"
        print(msg)
        checks.append(condition)
        return condition

    # Filesystem checks
    check("/etc/fstab exists",
          os.path.exists(f"{mount_point}/etc/fstab"))

    check("/etc/hostname exists",
          os.path.exists(f"{mount_point}/etc/hostname"))

    # Boot checks
    boot_dir = f"{mount_point}/boot"
    has_kernel = False
    has_initramfs = False
    if os.path.isdir(boot_dir):
        for f in os.listdir(boot_dir):
            if f.startswith("vmlinuz"):
                has_kernel = True
            if f.startswith("initramfs"):
                has_initramfs = True

    check("Kernel in /boot", has_kernel)
    check("Initramfs in /boot", has_initramfs)

    # Bootloader
    limine_conf = False
    for root, dirs, files in os.walk(boot_dir):
        for f in files:
            if f == "limine.conf":
                limine_conf = True
    check("Limine bootloader configured", limine_conf)

    # User exists
    passwd_path = f"{mount_point}/etc/passwd"
    user_found = False
    if os.path.exists(passwd_path):
        with open(passwd_path) as f:
            for line in f:
                if "testuser" in line or "erch" in line:
                    user_found = True
    check("Test user exists in /etc/passwd", user_found)

    # Btrfs subvolumes (check fstab or mount options)
    fstab_path = f"{mount_point}/etc/fstab"
    has_btrfs = False
    if os.path.exists(fstab_path):
        with open(fstab_path) as f:
            content = f.read()
            has_btrfs = "btrfs" in content
    check("Btrfs filesystem in fstab", has_btrfs)

    # Distro repo cloned
    erch_cloned = os.path.isdir(f"{mount_point}/root/erch") or \
                  os.path.isdir(f"{mount_point}/home/*/erch")
    eos_cloned = os.path.isdir(f"{mount_point}/root/e-os") or \
                 os.path.isdir(f"{mount_point}/home/*/e-os")
    check("Distro repo cloned (erch or e-os)", erch_cloned or eos_cloned)

    # Package manager
    pacman_db = os.path.isdir(f"{mount_point}/var/lib/pacman/local")
    check("Pacman database exists", pacman_db)

    # Systemd
    systemd_dir = f"{mount_point}/etc/systemd/system"
    check("Systemd system directory exists", os.path.isdir(systemd_dir))

    # os-release
    os_release = f"{mount_point}/etc/os-release"
    check("/etc/os-release exists", os.path.exists(os_release))

    # Summary
    passed = sum(1 for c in checks if c)
    total = len(checks)
    print()
    print(f"{'='*50}")
    print(f"Results: {passed}/{total} checks passed")
    print(f"{'='*50}")

    return all(checks)


def cleanup(nbd_dev, mount_point):
    run(["sudo", "umount", "-R", mount_point], check=False)
    run(["sudo", "qemu-nbd", "--disconnect", nbd_dev], check=False)


def main():
    parser = argparse.ArgumentParser(description="QEMU verify test")
    parser.add_argument("--disk", required=True, help="Path to qcow2 disk image")
    parser.add_argument("--profile", default="console",
                        choices=["console", "school", "uni", "org"])
    args = parser.parse_args()

    nbd_dev = setup_nbd(args.disk)
    mount_point = "/tmp/eos-verify"

    try:
        root_part = mount_root(nbd_dev, mount_point)
        success = verify(mount_point, args.profile)
        sys.exit(0 if success else 1)
    finally:
        cleanup(nbd_dev, mount_point)


if __name__ == "__main__":
    main()
