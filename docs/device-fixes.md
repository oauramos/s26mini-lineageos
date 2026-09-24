# Device-specific fixes

Status of hardware on replacement GSIs, and the tweaks that fix it. Apply the fixes with `scripts/post-install.sh` after the first boot.

## Tested GSIs

| GSI | Android | Security patch | Boots | Notes |
|---|---|---|---|---|
| **LineageOS 21 td `20260918`** (AndyYan, `arm64_bvN`) | 14 | **2026-09-01** | ✅ | recommended; installed with `scripts/flash-gsi.sh` |
| TrebleDroid `ci-20240226` (vanilla `arm64-ab`) | 14 | 2024-02-05 | ✅ | heavy load for the first minutes, one spontaneous reboot seen |
| TrebleDroid `ci-20230905` | 13 | 2023 | ⏳ | untested fallback |

## Hardware checklist (LineageOS 21 td)

| Feature | Status | Notes |
|---|---|---|
| Display / touch | ✅ | |
| Brightness | ✅ | backlight follows the setting (`/sys/class/leds/lcd-backlight`) |
| Wi-Fi 2.4 GHz | ✅ | scan OK |
| Wi-Fi 5 GHz | ✅ | scan sees 5 GHz APs (ch 44, 5220 MHz) |
| Bluetooth | ✅ with fix | crashes without the fix below; first connection after pairing may fail (see below) |
| Sensors (accel, light, proximity) | ✅ | detected |
| Cameras | ✅ | back + front |
| Audio | ✅ | |
| GPS | ✅ | no Google A-GPS on vanilla builds, first fix can take a few minutes. Test with GPSTest (F-Droid) |
| Mobile data / calls / SMS | ✅ | |
| Battery / charging | ✅ | reported correctly |

## Fixes

### Bluetooth crash loop (Android 14+)

Symptom: Bluetooth toggles back off. Logcat shows

```
E bluetooth: hci_layer.cc: Received UNEXPECTED command status:UNKNOWN_HCI_COMMAND opcode:0xc5a (READ_DEFAULT_ERRONEOUS_DATA_REPORTING)
F bt_shim_hci: assertion 'was_validated_' failed
```

The MT6739 controller advertises `READ_DEFAULT_ERRONEOUS_DATA_REPORTING` as supported but rejects it. Mask it out (supported-commands octet 18, bit 2):

```bash
adb shell setprop persist.sys.bt.unsupported.commands 182
```

This needs a TrebleDroid-based GSI (the property comes from TrebleDroid's Bluetooth patches). It persists across reboots.

### Bluetooth: first connection after pairing fails

Pairing works, but the first connection right after pairing may not come up. Power-cycle the accessory (headphones, speaker...) once. After that it connects normally.

### Google apps

Vanilla (`bvN`) builds have no Google Play Services, so Google Maps and other Google apps won't run. Use F-Droid apps (Organic Maps, OsmAnd), or flash the `bgN` (GApps) variant if you need them.

### Screen density

The panel is 384×854 and the GSI defaults to 190 dpi. If the UI looks cramped, try `adb shell wm density 160`.
