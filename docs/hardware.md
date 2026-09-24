# Hardware

Read from a running unit (dmesg, `/proc/hw_*_info`, I2C bus, camera HAL, sensor service). Other batches may use different panels or sensors.

## Identity

| | |
|---|---|
| Marketed as | "S26 ULTRA Mini" |
| Board | `d39g_4m_bml_s26ultra_mini_pt` (`alps` = MediaTek reference) |
| Stock firmware | `VK-D39G-4M-XJ9RQ0C4-3+16-BML-S26 ULTRA Mini-PT-20260227`, MTK `alps-mp-q0.mp1-V9.122.1` |
| Treble | yes, VNDK 29, dynamic partitions (`super` 4 GiB), A-only, system-as-root GSI format (`arm64_ab` / `bvN`) |
| USB | VID `0x0E8D` |

## Platform

| | |
|---|---|
| SoC | MediaTek **MT6739WA** (hw code `0x699`), 4× Cortex-A53, arm64 |
| GPU | PowerVR Rogue **GE8100**, GLES 3.2 |
| RAM | 3 GB |
| Storage | 16 GB eMMC (14.9 GB raw, ~10 GB userdata) |
| Kernel | 4.14.141, prebuilt. **No source released** by the manufacturer (GPL). |

## Display and input

| | |
|---|---|
| Panel | `axs15260_384x854_dsi_vdo_truly_ips_incell`: 384×854 IPS in-cell, MIPI DSI video mode, 60 Hz |
| Touch | `axs_ts` (AXS15260 in-cell) |
| Keys | `mtk-kpd` (power, volume) |

## Cameras

| | Sensor | Resolution |
|---|---|---|
| Back | OmniVision **OV9760** | 1496×1126 (~1.7 MP), AF VCM, LED flash |
| Front | GalaxyCore **GC2355** | 1600×1200 (2 MP) |

## Connectivity

| | |
|---|---|
| Wi-Fi / BT / GPS / FM | MT6739 integrated CONNSYS (WMT), FM `mt6627` |
| Modem | `MOLY.LR12A.R3.MP.V123.8`, LTE, dual SIM |
| NFC | no |

## Sensors

Accelerometer (`gsensor` / Sensortek `gsensor_stk`), ambient light + proximity (`alsps`). **No gyroscope, no compass.**

A fingerprint HAL is registered, but no fingerprint device node exists. Probably no real sensor.

## Power / audio

Li-ion, reports 2946 mAh. MTK `mt-snd-card` codec with ACCDET headset detection. Charger IC at I2C `3-006b`.

## Partition map

| Partition | Size | Notes |
|---|---|---|
| preloader (eMMC boot0) | 4 MB | needed for mtkclient DRAM init, see [recovery.md](recovery.md) |
| boot / recovery | 24 MB each | |
| dtbo | 8 MB | |
| vbmeta, vbmeta_system, vbmeta_vendor | 8 MB each | |
| **nvram / nvdata / nvcfg / proinfo** | 5 / 32 / 8 / 3 MB | IMEI + RF calibration: **back up, never share** |
| md1img / md1dsp | 64 / 16 MB | modem |
| lk, lk2, tee1, tee2, gz1, gz2, spmfw, mcupmfw, seccfg | small | boot chain |
| super | 4096 MB | logical: `system`, `vendor` (473 MB), `product` (stock, removed by the installer) |
| cache | 432 MB | |
| userdata | ~10 GB | |
