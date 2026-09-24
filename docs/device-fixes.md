# Device-specific fixes

Status of hardware on replacement GSIs, and the tweaks that fix it. Work in progress: fill in as things get tested.

## Tested GSIs

| GSI | Boots | Notes |
|---|---|---|
| TrebleDroid `ci-20240226` (Android 14, vanilla `arm64-ab`) | ✅ | first boot OK, heavy load for the first minutes, one spontaneous reboot seen |
| LineageOS 21 td `20260918` (AndyYan, `arm64_bvN`) | ⏳ | pending |
| TrebleDroid `ci-20230905` (Android 13) | ⏳ | fallback, untested |

## Hardware checklist

| Feature | TrebleDroid A14 | LineageOS 21 |
|---|---|---|
| Display / touch | ✅ | |
| Wi-Fi | ✅ radio up, association untested | |
| Bluetooth | ⚠️ was OFF after boot, untested | |
| Mobile data / calls / SMS | ⏳ | |
| Audio (speaker, earpiece, headset) | ⏳ | |
| Cameras | ⏳ (HAL exposes 2 cameras) | |
| Brightness control | ⏳ | |
| GPS | ⏳ | |
| Sensors (accel, light, proximity) | ✅ detected | |
| Battery / charging | ✅ reported correctly | |

## Known tweaks

- **Screen density:** the panel is 384×854 and the GSI defaults to 190 dpi. If UI elements look cramped, try `adb shell wm density 160`.
- MediaTek-specific toggles (brightness, audio, BT) live in the GSI's **Phh Treble Settings** app. Document here whichever toggle fixes what.
