#!/bin/bash

# E-OS Installer — TUI-based Arch installer with gum
# Boots from custom ISO, partitions disk, installs Arch, clones E-OS, runs boot.sh

set -eEo pipefail

# ── Colors ──────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
RESET='\033[0m'

# ── Globals ─────────────────────────────────────────────
DISK=""
USERNAME=""
PASSWORD=""
TIMEZONE=""
KEYMAP="us"
LOCALE="en_US"
HOSTNAME="e-os"
BOOT_PART=""
ROOT_PART=""
ROOT_TARGET=""
ENCRYPTION="false"
DISTRO="e-os"
REPO_URL="https://github.com/eande-uk/e-os.git"
REPO_BRANCH="dev"
EOS_PROFILE=""
LOG_FILE="/var/log/installer.log"

# ── Profile Detection ───────────────────────────────────
# ISO label contains profile: EOS_CONSOLE, EOS_SCHOOL, EOS_UNI, EOS_ORG
detect_profile() {
  local label
  label=$(blkid -s LABEL -o value /dev/disk/by-label/* 2>/dev/null | grep "EOS_" | head -1 || true)

  case "$label" in
    EOS_CONSOLE*) EOS_PROFILE="console" ;;
    EOS_SCHOOL*)  EOS_PROFILE="school" ;;
    EOS_UNI*)     EOS_PROFILE="uni" ;;
    EOS_ORG*)     EOS_PROFILE="org" ;;
    *)
      # Fallback: detect from ISO filename or prompt user
      EOS_PROFILE=""
      ;;
  esac
}

# ── Logging ─────────────────────────────────────────────
log() {
  echo "[$(date '+%H:%M:%S')] $1" >> "$LOG_FILE"
}

log_error() {
  echo "[$(date '+%H:%M:%S')] ERROR: $1" >> "$LOG_FILE"
}

# ── UI Helpers ──────────────────────────────────────────
header() {
  clear
  gum style \
    --border double \
    --align center \
    --width 60 \
    --padding "1 2" \
    "E-OS Installer" \
    "" \
    "This will install E-OS (${EOS_PROFILE}) to your hard drive." \
    "" \
    "Requirements:" \
    "  - Backed up data" \
    "  - Internet connection" \
    "  - Secure Boot disabled" \
    "  - UEFI mode"
}

section() {
  gum style \
    --border rounded \
    --align left \
    --width 60 \
    --padding "0 2" \
    "$1"
}

gum_property() {
  local label="$1" value="$2"
  gum join --horizontal \
    "$(gum style --foreground 240 "  ${label}:")" \
    "$(gum style --bold " ${value}")"
}

confirm_or_exit() {
  gum confirm --affirmative="Continue" --negative="Cancel" "$1" || {
    gum style --foreground 1 "Installation cancelled."
    sleep 1
    exit 0
  }
}

# ── Disk Detection ──────────────────────────────────────
select_disk() {
  section "Select Installation Disk"

  local disks=()
  while IFS= read -r line; do
    local name=$(echo "$line" | awk '{print $1}')
    local size=$(echo "$line" | awk '{print $2}')
    disks+=("${name} — ${size}")
  done < <(lsblk -dno NAME,SIZE,TYPE | grep disk)

  if [[ ${#disks[@]} -eq 0 ]]; then
    gum style --foreground 1 "No disks found!"
    sleep 2
    exit 1
  fi

  local selected
  selected=$(gum choose --header "Choose target disk:" "${disks[@]}")
  DISK="/dev/$(echo "$selected" | awk '{print $1}')"

  gum_property "Disk" "$DISK"
  confirm_or_exit "Wipe and partition $DISK?"
}

# ── User Setup ──────────────────────────────────────────
setup_user() {
  section "User Setup"

  USERNAME=$(gum input --header "Enter username:" --placeholder "user")
  if [[ -z "$USERNAME" ]]; then
    gum style --foreground 1 "Username cannot be empty!"
    sleep 1
    setup_user
    return
  fi

  while true; do
    PASSWORD=$(gum input --password --header "Enter password:")
    [[ -n "$PASSWORD" ]] && break
    gum style --foreground 1 "Password cannot be empty!"
  done

  local password_check
  password_check=$(gum input --password --header "Confirm password:")
  if [[ "$PASSWORD" != "$password_check" ]]; then
    gum style --foreground 1 "Passwords do not match!"
    sleep 1
    setup_user
    return
  fi

  gum_property "Username" "$USERNAME"
}

# ── System Config ───────────────────────────────────────
setup_system() {
  section "System Configuration"

  # Timezone
  local timezones=()
  while IFS= read -r tz; do
    timezones+=("$tz")
  done < <(find /usr/share/zoneinfo -type f -name "*" | sed 's|/usr/share/zoneinfo/||' | grep -v "^posix/" | grep -v "^right/" | sort)
  TIMEZONE=$(gum filter --header "Select timezone:" "${timezones[@]}" --limit 10 <<< "$(printf '%s\n' "${timezones[@]}" | head -20)")

  if [[ -z "$TIMEZONE" ]]; then
    TIMEZONE="UTC"
  fi

  # Keymap
  local keymaps=()
  while IFS= read -r km; do
    keymaps+=("$km")
  done < <(localectl list-keymaps 2>/dev/null || ls /usr/share/kbd/keymaps/ | sed 's|/.*||' | sort -u)
  KEYMAP=$(gum filter --header "Select keyboard layout:" "${keymaps[@]}" --limit 10 <<< "$(printf '%s\n' "${keymaps[@]}" | head -20)")

  if [[ -z "$KEYMAP" ]]; then
    KEYMAP="us"
  fi

  # Locale
  local locales=()
  while IFS= read -r loc; do
    locales+=("$loc")
  done < <(grep -v "^#" /etc/locale.gen 2>/dev/null | sed 's/ UTF-8//' | sort -u)
  LOCALE=$(gum filter --header "Select locale:" "${locales[@]}" --limit 10 <<< "$(printf '%s\n' "${locales[@]}" | head -20)")

  if [[ -z "$LOCALE" ]]; then
    LOCALE="en_US"
  fi

  # Hostname
  HOSTNAME=$(gum input --header "Enter hostname:" --value "e-os")
  [[ -z "$HOSTNAME" ]] && HOSTNAME="e-os"

  gum_property "Timezone" "$TIMEZONE"
  gum_property "Keymap" "$KEYMAP"
  gum_property "Locale" "$LOCALE"
  gum_property "Hostname" "$HOSTNAME"
}

# ── Summary ─────────────────────────────────────────────
show_summary() {
  section "Installation Summary"

  echo ""
  gum_property "Disk" "$DISK"
  gum_property "Username" "$USERNAME"
  gum_property "Timezone" "$TIMEZONE"
  gum_property "Keymap" "$KEYMAP"
  gum_property "Locale" "$LOCALE"
  gum_property "Hostname" "$HOSTNAME"
  gum_property "Profile" "$EOS_PROFILE"
  echo ""

  confirm_or_exit "Start installation?"
}

# ── Partitioning ────────────────────────────────────────
partition_disk() {
  log "Partitioning $DISK"

  gum spin --spinner dot --title "Wiping disk..." -- \
    bash -c "wipefs -af '$DISK' && sgdisk --zap-all '$DISK'"

  gum spin --spinner dot --title "Creating partitions..." -- \
    bash -c "
      sgdisk -o '$DISK'
      sgdisk -n 1:0:+1G -t 1:ef00 -c 1:boot '$DISK'
      sgdisk -n 2:0:0 -t 2:8300 -c 2:root '$DISK'
      partprobe '$DISK'
    "

  if [[ "$DISK" == /dev/nvme* ]]; then
    BOOT_PART="${DISK}p1"
    ROOT_PART="${DISK}p2"
  else
    BOOT_PART="${DISK}1"
    ROOT_PART="${DISK}2"
  fi

  ROOT_TARGET="$ROOT_PART"
  log "Boot: $BOOT_PART, Root: $ROOT_PART"
}

# ── Formatting ──────────────────────────────────────────
format_disk() {
  log "Formatting partitions"

  gum spin --spinner dot --title "Formatting boot partition (FAT32)..." -- \
    mkfs.fat -F 32 -n BOOT "$BOOT_PART"

  gum spin --spinner dot --title "Formatting root partition (Btrfs)..." -- \
    mkfs.btrfs -f -L ROOT "$ROOT_TARGET"
}

# ── Btrfs Subvolumes ────────────────────────────────────
create_subvolumes() {
  log "Creating Btrfs subvolumes"

  gum spin --spinner dot --title "Creating Btrfs subvolumes..." -- \
    bash -c "
      mount '$ROOT_TARGET' /mnt
      btrfs subvolume create /mnt/@
      btrfs subvolume create /mnt/@home
      btrfs subvolume create /mnt/@snapshots
      umount -R /mnt
    "
}

# ── Mounting ────────────────────────────────────────────
mount_filesystem() {
  log "Mounting filesystem"

  local mount_opts="defaults,noatime,compress=zstd"

  gum spin --spinner dot --title "Mounting subvolumes..." -- \
    bash -c "
      mount --mkdir -t btrfs -o ${mount_opts},subvol=@ '$ROOT_TARGET' /mnt
      mount --mkdir -t btrfs -o ${mount_opts},subvol=@home '$ROOT_TARGET' /mnt/home
      mount --mkdir -t btrfs -o ${mount_opts},subvol=@snapshots '$ROOT_TARGET' /mnt/.snapshots
      mount --mkdir '$BOOT_PART' /mnt/boot
      mkdir -p /mnt/var/lib/portables
      mkdir -p /mnt/var/lib/machines
    "
}

# ── pacstrap ────────────────────────────────────────────
install_base() {
  log "Installing base system"

  local packages=("base" "linux" "linux-firmware" "base-devel" "btrfs-progs" "limine" "efibootmgr" "networkmanager")

  if grep -q "GenuineIntel" /proc/cpuinfo; then
    packages+=("intel-ucode")
  elif grep -q "AuthenticAMD" /proc/cpuinfo; then
    packages+=("amd-ucode")
  fi

  gum spin --spinner dot --title "Installing base packages (this may take a while)..." -- \
    pacstrap -K /mnt "${packages[@]}"
}

# ── System Configuration ────────────────────────────────
configure_system() {
  log "Configuring system"

  genfstab -U /mnt >> /mnt/etc/fstab

  arch-chroot /mnt ln -sf "/usr/share/zoneinfo/${TIMEZONE}" /etc/localtime
  arch-chroot /mnt hwclock --systohc

  echo "LANG=${LOCALE}.UTF-8" > /mnt/etc/locale.conf
  sed -i "s/^#${LOCALE}.UTF-8 UTF-8/${LOCALE}.UTF-8 UTF-8/" /mnt/etc/locale.gen 2>/dev/null || true
  sed -i "s/^#en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/" /mnt/etc/locale.gen 2>/dev/null || true
  arch-chroot /mnt locale-gen

  echo "KEYMAP=${KEYMAP}" > /mnt/etc/vconsole.conf

  echo "${HOSTNAME}" > /mnt/etc/hostname
  cat > /mnt/etc/hosts << EOF
127.0.0.1  localhost.localdomain  localhost
::1        localhost.localdomain  localhost
127.0.1.1  ${HOSTNAME}.localdomain  ${HOSTNAME}
EOF

  arch-chroot /mnt mkinitcpio -P
}

# ── Bootloader ──────────────────────────────────────────
setup_bootloader() {
  log "Setting up Limine bootloader"

  local root_uuid
  root_uuid=$(blkid -s UUID -o value "$ROOT_TARGET")

  gum spin --spinner dot --title "Installing Limine..." -- \
    bash -c "
      arch-chroot /mnt bash -c '
        mkdir -p /boot/EFI/limine
        cp /usr/share/limine/BOOTX64.EFI /boot/EFI/limine/

        cat > /boot/limine.conf << LIMINEEOF
timeout: 3

/E-OS
    protocol: linux
    path: boot():/vmlinuz-linux
    cmdline: quiet root=UUID=${root_uuid} rw rootflags=subvol=@ rootfstype=btrfs rw init=/usr/lib/systemd/systemd
    module_path: boot():/initramfs-linux.img

/E-OS (fallback)
    protocol: linux
    path: boot():/vmlinuz-linux
    cmdline: quiet root=UUID=${root_uuid} rw rootflags=subvol=@ rootfstype=btrfs rw init=/usr/lib/systemd/systemd
    module_path: boot():/initramfs-linux-fallback.img
LIMINEEOF

        efibootmgr --create \
          --disk $(echo $DISK | sed "s|/dev/||") \
          --part 1 \
          --label \"E-OS (Limine)\" \
          --loader \"\\\\EFI\\\\limine\\\\BOOTX64.EFI\" \
          --unicode
      '
    "
}

# ── User Creation ───────────────────────────────────────
create_user() {
  log "Creating user $USERNAME"

  arch-chroot /mnt useradd -m -G wheel -s /bin/bash "$USERNAME"
  echo "${USERNAME}:${PASSWORD}" | arch-chroot /mnt chpasswd
  echo "root:${PASSWORD}" | arch-chroot /mnt chpasswd

  sed -i 's^# %wheel ALL=(ALL:ALL) ALL^%wheel ALL=(ALL:ALL) ALL^g' /mnt/etc/sudoers
}

# ── Services ────────────────────────────────────────────
enable_services() {
  log "Enabling services"

  arch-chroot /mnt systemctl enable NetworkManager
  arch-chroot /mnt systemctl enable fstrim.timer
}

# ── Clone Distro ────────────────────────────────────────
clone_distro() {
  log "Cloning $DISTRO"

  gum spin --spinner dot --title "Cloning E-OS repo..." -- \
    bash -c "
      arch-chroot /mnt bash -c '
        git clone --branch ${REPO_BRANCH} ${REPO_URL} /root/${DISTRO}
      '
    "
}

# ── Run Distro Installer ────────────────────────────────
run_distro_installer() {
  log "Running E-OS installer with profile=$EOS_PROFILE"

  gum spin --spinner dot --title "Running E-OS installer (profile: $EOS_PROFILE)..." -- \
    bash -c "
      arch-chroot /mnt bash -c '
        export EOS_PROFILE=\"${EOS_PROFILE}\"
        export EOS_PATH=\"/root/${DISTRO}\"
        cd /root/${DISTRO}
        chmod +x install.sh
        ./install.sh
      '
    "
}

# ── Cleanup ─────────────────────────────────────────────
cleanup() {
  log "Cleaning up"
  gum spin --spinner dot --title "Unmounting filesystems..." -- \
    umount -R /mnt 2>/dev/null || true
}

# ── Error Handler ───────────────────────────────────────
trap_handler() {
  local exit_code=$?
  if (( exit_code != 0 )); then
    log_error "Installation failed with exit code $exit_code"
    gum style --foreground 1 "Installation failed! Check $LOG_FILE for details."
    gum style "Press Enter to reboot, or Ctrl+C to exit."
    read -r
    cleanup
    reboot
  fi
}

trap trap_handler ERR INT TERM

# ── Main ────────────────────────────────────────────────
main() {
  mkdir -p "$(dirname "$LOG_FILE")"
  touch "$LOG_FILE"
  log "Installer started"

  # Detect profile from ISO label
  detect_profile

  # If no profile detected, prompt user
  if [[ -z "$EOS_PROFILE" ]]; then
    section "Select E-OS Profile"
    EOS_PROFILE=$(gum choose --header "Choose profile:" "console" "school" "uni" "org")
    [[ -z "$EOS_PROFILE" ]] && EOS_PROFILE="console"
  fi

  HOSTNAME="e-os-${EOS_PROFILE}"

  # Check internet
  if ! ping -c 1 -W 3 1.1.1.1 &>/dev/null; then
    gum style --foreground 1 "No internet connection!"
    gum style "Please connect to a network first."
    gum style "Press Enter to retry, or Ctrl+C to exit."
    read -r
    main
    return
  fi

  header
  gum spin --spinner dot --title "Checking environment..." -- sleep 1

  select_disk
  setup_user
  setup_system
  show_summary

  gum style --foreground 2 "Starting installation..."
  sleep 1

  partition_disk
  format_disk
  create_subvolumes
  mount_filesystem
  install_base
  configure_system
  setup_bootloader
  create_user
  enable_services
  clone_distro
  run_distro_installer

  gum style \
    --border double \
    --align center \
    --width 60 \
    --padding "1 2" \
    --foreground 2 \
    "Installation Complete!" \
    "" \
    "Remove the USB drive and press Enter to reboot."

  read -r
  cleanup
  reboot
}

main "$@"
