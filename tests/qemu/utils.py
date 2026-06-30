#!/usr/bin/env python3
"""Shared utilities for QEMU e2e tests."""

import os
import subprocess
import sys
import time


def find_iso(profile="erch"):
    """Find ISO file in iso/out/ directory."""
    iso_dir = os.path.join(os.path.dirname(__file__), "..", "..", "iso", "out")
    if not os.path.isdir(iso_dir):
        return None
    for f in os.listdir(iso_dir):
        if f.startswith(profile) and f.endswith(".iso"):
            return os.path.join(iso_dir, f)
    return None


def create_disk(path, size_gb=20):
    """Create a qcow2 disk image."""
    if os.path.exists(path):
        os.remove(path)
    subprocess.run(
        ["qemu-img", "create", "-f", "qcow2", path, f"{size_gb}G"],
        check=True,
        capture_output=True,
    )
    return path


def check_qemu():
    """Check if QEMU is available."""
    result = subprocess.run(
        ["qemu-system-x86_64", "--version"],
        capture_output=True,
        text=True,
    )
    return result.returncode == 0


def check_kvm():
    """Check if KVM is available."""
    return os.path.exists("/dev/kvm")


def check_pexpect():
    """Check if pexpect is available."""
    try:
        import pexpect
        return True
    except ImportError:
        return False


def run_qemu_boot(iso_path, timeout=120, uefi=True, disk_path=None):
    """Run QEMU boot test and return (success, log_text)."""
    cmd = [
        "qemu-system-x86_64",
        "-m", "4G",
        "-smp", "2",
        "-enable-kvm",
        "-nographic",
        "-serial", "mon:stdio",
        "-no-reboot",
        "-machine", "type=q35,accel=kvm",
        "-cpu", "host",
        "-cdrom", iso_path,
        "-boot", "d",
    ]

    if disk_path:
        cmd.extend(["-drive", f"file={disk_path},if=virtio,format=qcow2"])

    if uefi:
        cmd.extend([
            "-drive", "if=pflash,format=raw,unit=0,"
            "file=/usr/share/edk2/x64/OVMF_CODE.4m.fd,readonly=on",
            "-drive", "if=pflash,format=raw,unit=1,"
            "file=/tmp/OVMF_VARS.4m.fd",
        ])

    try:
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            timeout=timeout,
        )
        output = result.stdout + result.stderr
        success = "login:" in output or "root@" in output
        return success, output
    except subprocess.TimeoutExpired:
        return False, "TIMEOUT"
    except Exception as e:
        return False, str(e)
