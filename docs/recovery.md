# Recovery (bricks, bootloops)

## Soft brick: GSI bootloops

The bootloader stays unlocked, so fastboot still works:

1. Power off (hold Power ~15 s).
2. Hold **Volume Up + Power**. In the boot menu, pick **Fastboot**. Menus vary per unit; on some it's Volume Down + Power.
3. Flash another GSI with `scripts/flash-gsi.sh <other.img>`. Known-good fallback: TrebleDroid `ci-20230905` (Android 13).

## Hard brick: no fastboot

The MT6739 BootROM on this board is **unprotected** (SBC, SLA and DAA disabled), so [mtkclient](https://github.com/bkerler/mtkclient) can read and write flash over USB even when nothing boots.

**Enter BROM mode:** phone fully off, hold **Volume Up + Volume Down**, plug the USB cable, keep holding ~5 s. The screen stays black.

**DRAM init needs this board's preloader.** Without it (or with a preloader from another MT6739 board), mtkclient initialises DRAM but fails at DA stage 2 (`Stage was't executed. Maybe dram issue ?`). Dump the preloader from your own phone **before** you need it (needs root, which the TrebleDroid/Lineage userdebug GSIs provide via `adb root`):

```bash
adb root
adb exec-out "dd if=/dev/block/mmcblk0boot0 bs=1M" > preloader_boot0.img
```

Then:

```bash
python mtk.py printgpt --preloader preloader_boot0.img
python mtk.py rl backup_dir --preloader preloader_boot0.img   # full backup
```

### macOS note

Recent mtkclient hard-fails on macOS without macFUSE (`OSError: Unable to find libfuse`), because `mfusepy` raises `OSError`, not `ImportError`. If you don't need `fs` mounting, widen the two `except ImportError:` around the `mfusepy` imports in `mtkclient/Library/Filesystem/mtkdafs.py` and `mtkclient/Library/DA/mtk_da_handler.py` to `except (ImportError, OSError):`. Also point pyusb to a libusb (for example `DYLD_FALLBACK_LIBRARY_PATH=$(brew --prefix libusb)/lib`).

Preloader mode (phone off, no buttons) never enumerated on the unit tested. Use BROM.

## Back up your calibration data

`nvram`, `nvdata`, `nvcfg` and `proinfo` hold **your** IMEI, Wi-Fi/BT MAC and RF calibration. If they get corrupted you lose mobile data and Wi-Fi. Back them up once and **never share them**:

```bash
adb root
for p in nvram nvdata nvcfg proinfo; do
  adb exec-out "dd if=/dev/block/by-name/$p bs=1M" > "$p.img"
done
```
