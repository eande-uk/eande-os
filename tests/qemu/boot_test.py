#!/usr/bin/env python3
"""QEMU boot test — verify ISO boots to login prompt.

Usage:
    python3 boot_test.py --iso <path-to-iso> [--timeout 120] [--uefi]

Exit codes:
    0 — boot succeeded (login prompt reached)
    1 — boot failed (timeout or error)
"""

import argparse
import pexpect
import sys
import time

QEMU_BASE = [
    "qemu-system-x86_64",
    "-m", "4G",
    "-smp", "2",
    "-enable-kvm",
    "-nographic",
    "-serial", "mon:stdio",
    "-no-reboot",
    "-machine", "type=q35,accel=kvm",
    "-cpu", "host",
]


def boot_test(iso_path: str, timeout: int = 120, uefi: bool = True) -> bool:
    """Boot ISO and wait for login prompt."""
    cmd = QEMU_BASE.copy()
    cmd.extend(["-cdrom", iso_path, "-boot", "d"])

    if uefi:
        cmd.extend([
            "-drive", "if=pflash,format=raw,unit=0,file=/usr/share/edk2/x64/OVMF_CODE.4m.fd,readonly=on",
            "-drive", "if=pflash,format=raw,unit=1,file=/tmp/OVMF_VARS.4m.fd",
        ])

    print(f"[*] Booting ISO: {iso_path}")
    print(f"[*] UEFI: {uefi}, Timeout: {timeout}s")
    print(f"[*] QEMU command: {' '.join(cmd)}")
    print()

    child = pexpect.spawn(
        " ".join(cmd),
        encoding="utf-8",
        timeout=timeout,
    )
    child.logfile = sys.stdout

    try:
        # Wait for login prompt or shell prompt
        index = child.expect([
            r"login:",
            r"root@archiso",
            r"# ",
            r"~ #",
        ], timeout=timeout)

        print()
        print(f"[+] BOOT SUCCESS — reached prompt (match index: {index})")

        # Login if we got a login prompt
        if index == 0:
            child.sendline("root")
            child.expect([r"# ", r"~ #"], timeout=10)

        # Run a quick check
        child.sendline("uname -a")
        child.expect([r"# ", r"~ #"], timeout=10)
        print(f"[+] Kernel: {child.before.strip()}")

        child.sendline("poweroff")
        child.expect(pexpect.EOF, timeout=30)
        return True

    except pexpect.TIMEOUT:
        print()
        print("[-] BOOT FAILED — timeout waiting for login prompt")
        child.close(force=True)
        return False
    except pexpect.EOF:
        print()
        print("[-] BOOT FAILED — unexpected end of output")
        child.close(force=True)
        return False


def main():
    parser = argparse.ArgumentParser(description="QEMU boot test for E&E OS ISOs")
    parser.add_argument("--iso", required=True, help="Path to ISO file")
    parser.add_argument("--timeout", type=int, default=120, help="Boot timeout in seconds")
    parser.add_argument("--uefi", action="store_true", default=True, help="Boot in UEFI mode")
    parser.add_argument("--no-uefi", action="store_true", help="Boot in BIOS mode")
    args = parser.parse_args()

    uefi = not args.no_uefi
    success = boot_test(args.iso, args.timeout, uefi)
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
