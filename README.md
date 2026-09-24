<h1 align="center">s26mini-lineageos</h1>

<p align="center">
  <b>Remove the factory backdoor from the fake "S26 ULTRA Mini" and run LineageOS 21 on it.</b>
</p>

<p align="center">
  <a href="https://github.com/oauramos/s26mini-lineageos/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/oauramos/s26mini-lineageos"></a>
  <a href="LICENSE"><img alt="License: MIT" src="https://img.shields.io/github/license/oauramos/s26mini-lineageos"></a>
  <a href="https://github.com/oauramos/s26mini-lineageos/wiki"><img alt="Wiki" src="https://img.shields.io/badge/docs-wiki-blue"></a>
</p>

<p align="center">
  🇧🇷 <a href="README.pt-BR.md"><b>Leia em português</b></a>
</p>

---

A cheap **"Samsung S26 ULTRA Mini"** is sold all over Brazil. It isn't a Samsung: it's a generic **MediaTek MT6739** board whose firmware ships, **straight out of a sealed box**, with:

- a **code loader hidden in Android's core runtime** that can inject code into every app on the phone (banking, WhatsApp, browser)
- privileged tools that **rewrite the IMEI** and fake the phone's specs
- a **faked security patch** (claims 2026, really 2020)
- a Chinese "system update" service running permanently as the system user

This project **replaces the whole system** with **LineageOS 21** (Android 14, September 2026 security patches). Everything above lives in `/system`, so it goes away.

| | Stock | With this project |
|---|---|---|
| OS | Android 10 Go, disguised | **LineageOS 21** (Android 14) |
| Security patch | claims 2026-01, really 2020-08 | **2026-09-01** |
| Backdoor loader, IMEI tools, FOTA | present | **gone** |

Full analysis with indicators: [Stock Firmware Analysis](https://github.com/oauramos/s26mini-lineageos/wiki/Stock-Firmware-Analysis) · [docs/backdoor-analysis.md](docs/backdoor-analysis.md)

> ⚠️ **This erases everything on the phone** and unlocks the bootloader. Only tested on board `d39g_4m_bml_s26ultra_mini_pt`. You do this at your own risk.
>
> 🚫 This project does **not** provide or support changing IMEIs. Doing so is a crime in Brazil. The clean system removes the IMEI tools that came with the phone.

## Quick start

```bash
git clone https://github.com/oauramos/s26mini-lineageos && cd s26mini-lineageos
scripts/flash-gsi.sh lineage-21.0-<date>-UNOFFICIAL-arm64_bvN.img   # wipes + installs
scripts/post-install.sh                                              # device fixes
scripts/profile-termux.sh    # Termux edition  (or: scripts/install-apps.sh for Full)
```

Step-by-step, with troubleshooting: [Installation Guide](https://github.com/oauramos/s26mini-lineageos/wiki/Installation-Guide).

## Is my phone this model?

Enable USB debugging and run:

```bash
adb shell getprop ro.product.vendor.device
```

It should print `d39g_4m_bml_s26ultra_mini_pt`. Other boards sold under the same name may differ. The installer warns you and asks before continuing.

## What you need

- A computer with `adb`, `fastboot`, `python3` and `curl` (Linux or macOS; Windows via WSL untested)
- A **USB-A → USB-C** cable. **USB-C ↔ USB-C doesn't work** on this phone, for data or charging (the board lacks the Type-C CC resistors). On USB-C-only computers, use a single C-to-A adapter, not a hub.
- A GSI image, **arm64, A/B (system-as-root), vanilla**:
  - **Recommended:** LineageOS 21 by AndyYan, file `...-arm64_bvN.img.gz`: <https://sourceforge.net/projects/andyyan-gsi/files/lineage-21-pre-qpr2-td/> (slow download? add `?use_mirror=cfhcable` to the URL)
  - Alternative: TrebleDroid `system-td-arm64-ab-vanilla.img.xz`: <https://github.com/TrebleDroid/treble_experimentations/releases> (`ci-20240226` = Android 14, `ci-20230905` = Android 13)
  - Avoid EROFS builds and GSIs that require Android 11+ vendors. This phone has an Android 10 vendor (VNDK 29) and kernel 4.14.

## Install

1. **On the phone:** Settings → About phone → tap **Build number** 7×. Then Developer options → enable **OEM unlocking** and **USB debugging**.
2. Plug it in, accept the "Allow USB debugging?" prompt.
3. Decompress the GSI (`gunzip file.img.gz` or `xz -d file.img.xz`).
4. Run `scripts/flash-gsi.sh path/to/system.img`.
5. When the phone shows the unlock warning, press **Volume Up** to confirm.

The script checks the board and OEM unlock, unlocks the bootloader (**wipes the phone**), disables AVB with a generated vbmeta ([tools/make_vbmeta_disabled.py](tools/make_vbmeta_disabled.py)), deletes the stock `product` partition, flashes the GSI and wipes userdata.

The first boot takes a few minutes. The "orange state" warning on every boot is expected with an unlocked bootloader.

## Pick an edition

After the first boot, enable USB debugging again and run `scripts/post-install.sh` (fixes the Bluetooth crash loop), then **one** of:

| | **Full** | **Termux edition** |
|---|---|---|
| Command | `scripts/install-apps.sh` | `scripts/profile-termux.sh` |
| For | a regular phone | a dedicated SSH/terminal device |
| Settings, Launcher, Camera, Clock, Calculator, Files, Keyboard | ✅ | ✅ |
| LineageOS browser (Jelly) | ✅ | replaced by Firefox |
| Gallery, Music, Recorder, Calendar | ✅ | disabled |
| Phone, Contacts, SMS | ✅ | disabled (edit `profile-termux.sh` if you need calls) |
| Screensavers, printing, live wallpapers, Seedvault, AudioFX | ✅ | disabled |
| Animations 0.5×, 160 dpi, apps precompiled | – | ✅ |
| **App bundle** (below) | ✅ | ✅ |
| **[Firefox](https://f-droid.org/packages/org.mozilla.fennec_fdroid/)** (Fennec F-Droid), default browser | – | ✅ |
| **[Primitive FTPd](https://f-droid.org/packages/org.primftpd/)**: FTP/SFTP server on the phone | – | ✅ |
| Free RAM (tested) | ~1.3 GB | **~1.9 GB** |

**App bundle**, from the official F-Droid repo with each APK pinned by SHA-256:

- [F-Droid](https://f-droid.org): open-source app store, keeps everything here updated
- [Termux](https://termux.dev): terminal with `ssh`, `sftp`, `scp`, `git`...
- [LocalSend](https://localsend.org): send files between computers and phones on the same network, no internet needed
- [Material Files](https://github.com/zhanghai/MaterialFiles): file manager with a built-in **SFTP**, SMB, FTP and WebDAV client (the phone reaching servers)

**Primitive FTPd** is the other direction: an SFTP/FTP server so a computer can reach the phone's files. **Fennec F-Droid** is Firefox built from Mozilla's source by F-Droid, without telemetry. Third-party builds can't use the "Firefox" name.

Disabled apps aren't removed: `adb shell pm enable <package>` brings any of them back. There are no Google services, so Google apps (Maps, Play Store...) won't work.

Full app list with versions and hashes: [Editions and Apps](https://github.com/oauramos/s26mini-lineageos/wiki/Editions-and-Apps).

## What works

| Works | Notes |
|---|---|
| Display, touch, brightness, Wi-Fi 2.4 + 5 GHz, calls, mobile data, SMS, cameras, audio, GPS, sensors, battery | GPS has no Google A-GPS, so the first fix is slower |
| Bluetooth | needs `post-install.sh` |

**Known issues:**

- **Bluetooth headphones pair but don't connect** when the phone starts the connection. Power-cycle the accessory once after pairing, and it reconnects by itself from then on. [Bluetooth Deep Dive](https://github.com/oauramos/s26mini-lineageos/wiki/Bluetooth-Deep-Dive)
- **USB-C ↔ USB-C** doesn't work. Use USB-A → C.

Details: [Device Fixes and Known Issues](https://github.com/oauramos/s26mini-lineageos/wiki/Device-Fixes-and-Known-Issues).

## After installing

- **Back up your calibration partitions** (IMEI, MAC, radio calibration) and the preloader: [Recovery and Backups](https://github.com/oauramos/s26mini-lineageos/wiki/Recovery-and-Backups).
- Open **F-Droid** once so it refreshes and offers updates.
- Turn **USB debugging off**. The GSIs are `userdebug` builds, which allow root over ADB.

## What stays untrusted

Vendor, kernel, bootloader and modem firmware are still the manufacturer's, and nothing replaces them (the kernel source was never published). No backdoor indicators were found in `vendor`, but treat the phone as **semi-trusted**: no main Google account, no banking apps. [Security Model](https://github.com/oauramos/s26mini-lineageos/wiki/Security-Model)

## Documentation

The **[wiki](https://github.com/oauramos/s26mini-lineageos/wiki)** has the full story:

| | |
|---|---|
| [Hardware](https://github.com/oauramos/s26mini-lineageos/wiki/Hardware) | real SoC, panel, cameras, I2C map, partitions, USB-C quirk |
| [Stock Firmware Analysis](https://github.com/oauramos/s26mini-lineageos/wiki/Stock-Firmware-Analysis) | the backdoor, IMEI tools, fake specs, IOCs |
| [Design Decisions](https://github.com/oauramos/s26mini-lineageos/wiki/Design-Decisions) | why a GSI, why LineageOS 21, why these apps... |
| [Installation Guide](https://github.com/oauramos/s26mini-lineageos/wiki/Installation-Guide) | step-by-step and troubleshooting |
| [Editions and Apps](https://github.com/oauramos/s26mini-lineageos/wiki/Editions-and-Apps) | every kept, disabled and added app |
| [Recovery and Backups](https://github.com/oauramos/s26mini-lineageos/wiki/Recovery-and-Backups) | BROM/mtkclient, calibration |
| [Security Model](https://github.com/oauramos/s26mini-lineageos/wiki/Security-Model) | what to trust, SSH hardening |
| [FAQ](https://github.com/oauramos/s26mini-lineageos/wiki/FAQ) | quick answers |

The same material lives in [`docs/`](docs/) for offline reading.

## Contributing

Tested another GSI or another board sold as "S26 ULTRA Mini"? Open an issue with the output of:

```bash
adb shell getprop | grep -E 'ro.product.vendor.device|ro.build.display.id|ro.vendor.build.security_patch|ro.board.platform'
```

**Never post** your IMEI, serial number or `nvram`/`nvdata` dumps.

## License

Scripts and docs: MIT. See [LICENSE](LICENSE).
