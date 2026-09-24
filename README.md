# Fake "S26 ULTRA Mini": remove the factory backdoor

🇧🇷 [Leia em português](README.pt-BR.md)

A cheap "Samsung S26 ULTRA Mini" is sold all over Brazil. It isn't a Samsung: it's a generic MediaTek MT6739 board whose firmware ships, **straight out of a sealed box**, with:

- a **code loader hidden in Android's core runtime** that can inject code into every app on the phone (banking, WhatsApp, browser)
- privileged tools that **rewrite the IMEI** and fake the phone's specs
- a **faked security patch** (claims 2026, really 2020)
- a Chinese "system update" service running permanently as the system user

Full write-up with indicators: [docs/backdoor-analysis.md](docs/backdoor-analysis.md).

This repo shows how to **replace the whole system partition** with a clean, open-source Android (a Treble GSI such as LineageOS or TrebleDroid). All of the above lives in `/system`, so it goes away.

> ⚠️ **This erases everything on the phone.** It unlocks the bootloader and flashes unofficial software. Only tested on board `d39g_4m_bml_s26ultra_mini_pt`. You do this at your own risk.
>
> 🚫 This project does **not** provide or support changing IMEIs. Doing so is a crime in Brazil. The clean system removes the IMEI tools that came with the phone.

## Is my phone this model?

Enable USB debugging and run:

```bash
adb shell getprop ro.product.vendor.device
```

It should print `d39g_4m_bml_s26ultra_mini_pt`. Other boards sold under the same name may differ. The installer warns you and asks before continuing.

## What you need

- A computer with `adb`, `fastboot` and `python3` (Linux or macOS; Windows via WSL untested)
- A data USB cable, plugged **directly** into the computer (no hub)
- A GSI image, **arm64, A/B (system-as-root), vanilla**:
  - **Recommended:** LineageOS 21 by AndyYan, file `...-arm64_bvN.img.gz`: <https://sourceforge.net/projects/andyyan-gsi/files/lineage-21-pre-qpr2-td/>
  - Alternative: TrebleDroid `system-td-arm64-ab-vanilla.img.xz`: <https://github.com/TrebleDroid/treble_experimentations/releases> (`ci-20240226` = Android 14, `ci-20230905` = Android 13)
  - Avoid EROFS builds and GSIs that require Android 11+ vendors. This phone has an Android 10 vendor (VNDK 29) and kernel 4.14.

## Install

1. **On the phone:** Settings → About phone → tap **Build number** 7×. Then Developer options → enable **OEM unlocking** and **USB debugging**.
2. Plug it in, accept the "Allow USB debugging?" prompt.
3. Decompress the GSI (`gunzip file.img.gz` or `xz -d file.img.xz`).
4. Run:
   ```bash
   scripts/flash-gsi.sh path/to/system.img
   ```
5. When the phone shows the unlock warning, press **Volume Up** to confirm.

The script:

1. checks the board and that OEM unlocking is allowed
2. reboots to fastboot and unlocks the bootloader (**wipes the phone**)
3. flashes an empty vbmeta with verification disabled to `vbmeta`, `vbmeta_system` and `vbmeta_vendor` ([tools/make_vbmeta_disabled.py](tools/make_vbmeta_disabled.py))
4. reboots to fastbootd and deletes the stock `product` partition (Android 10 apps that break newer GSIs)
5. flashes the GSI to `system`, wipes userdata and reboots

The first boot takes a few minutes. The phone shows an "orange state" warning on every boot. That's expected with an unlocked bootloader.

### After installing

- Once it has booted, enable USB debugging again and run `scripts/post-install.sh`. It applies the device fixes (Bluetooth crash loop on Android 14+). Optional: `scripts/optimize.sh` for a lighter, snappier system.
- **Back up your calibration partitions** (IMEI, Wi-Fi/BT MAC, radio calibration) and dump the preloader. You'll need both if something breaks: [docs/recovery.md](docs/recovery.md).
- Turn **USB debugging off** when you're done. The GSIs are `userdebug` builds, which allow root over ADB.
- Hardware status and fixes: [docs/device-fixes.md](docs/device-fixes.md).

## What stays untrusted

Vendor, kernel, bootloader and modem firmware are still the manufacturer's, and no clean replacement exists (the kernel source was never published). No backdoor indicators were found in `vendor`, but treat the phone as **semi-trusted**: no main Google account, no banking apps.

## Docs

| | |
|---|---|
| [docs/backdoor-analysis.md](docs/backdoor-analysis.md) | what the stock firmware does, with hashes and package names |
| [docs/hardware.md](docs/hardware.md) | real hardware: SoC, panel, cameras, sensors, partition map |
| [docs/recovery.md](docs/recovery.md) | bootloops, hard bricks via mtkclient/BROM, calibration backup |
| [docs/device-fixes.md](docs/device-fixes.md) | per-GSI hardware status and tweaks |

## Contributing

Tested another GSI or another board sold as "S26 ULTRA Mini"? Open an issue with the output of:

```bash
adb shell getprop | grep -E 'ro.product.vendor.device|ro.build.display.id|ro.vendor.build.security_patch|ro.board.platform'
```

**Never post** your IMEI, serial number or `nvram`/`nvdata` dumps.

## License

Scripts and docs: MIT. See [LICENSE](LICENSE).
