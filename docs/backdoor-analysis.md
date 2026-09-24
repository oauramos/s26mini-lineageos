# Stock firmware analysis

Unit analysed: sealed, brand-new "S26 ULTRA Mini" bought in Brazil (September 2026). It was powered on straight out of the box and never had a user app installed. Every unit running the same firmware is expected to ship with the same code.

## TL;DR

- The phone is a **counterfeit**: Samsung never made an "S26 Ultra Mini". It's a generic MediaTek MT6739 reference board (`alps`).
- The core Android runtime library carries a **hidden code loader** that can inject code into every app on the phone.
- It ships **IMEI-rewriting and spec-spoofing tools** with system privileges.
- The advertised security patch is **faked**.
- A Chinese FOTA ("system update") service runs permanently as the system user.

## 1. Backdoor loader in `libandroid_runtime.so`

`libandroid_runtime.so` is loaded by Zygote, so it runs inside **every** app process (banking, WhatsApp, browser...). The stock build of this library contains:

```
dalvik/system/DexClassLoader
classOfDexClassLoader
dexClassLoaderInitArgRtnT
/data/local/tmp/.SystemConfig
/data/local/tmp/.SystemData
/sdcard/.SystemConfig
/sdcard/.SystemData
/sdcard/Android/media/JpgConfig
/sdcard/Android/media/JpgData
```

A stock AOSP `libandroid_runtime.so` never references `DexClassLoader` or hidden files on shared storage. The pattern (a native loader in the runtime library that pulls DEX payloads from hidden paths into every process) matches the factory-firmware backdoor families publicly documented as Triada and Keenadu. We did not attribute it to a specific family.

On the unit analysed, the payload files were **not present yet**. The loader was in place and could be activated at any time (for example by a payload dropped through the FOTA service).

Present in both the 32-bit (`/system/lib`) and 64-bit (`/system/lib64`) copies.

## 2. Pre-installed privileged tools

| Package | Path | Runs as | What it does |
|---|---|---|---|
| `com.example.artificialswitch` | `/system/app/ArtificialSwitchQmx` | system (uid 1000) | Receivers `CHANGEDIMEI`, `CHANGESIMEI`, `LYSCHANGESYSTEMPRO`: rewrites IMEI and system properties (how the phone lies about its model/specs) |
| `com.lys.writeimei` | `/system/app/WriteIMEI` | system | IMEI writer |
| `com.abfota.systemUpdate` (Rsota) | `/system/app/Rsota` | system, **persistent** | FOTA client with location, camera, calls, package install and more |
| `com.andromeda.androbench2` | `/system/app/Androbench` | app | Storage benchmark (suspicious on a phone that fakes its specs) |

Changing a phone's IMEI is a crime in Brazil. **The replacement system in this repo removes all of the above.**

## 3. Faked specs and patch level

| Claim | Reality |
|---|---|
| "S26 ULTRA Mini" | board `d39g_4m_bml_s26ultra_mini_pt`, brand `alps` |
| Security patch 2026-01-05 | vendor security patch **2020-08-05**; kernel 4.14.141 |
| Android 10 | Android 10 **Go** (low-RAM edition) |
| System fingerprint | copied from a different device: `alps/full_f202_f22_p10/...:10.0/MRA58K/...` (`MRA58K` is an Android 6 build ID) |
| Display | 384×854 |
| Cameras | 1.7 MP back (OV9760), 2 MP front (GC2355) |

## 4. Network

At first check, the system uid had 4 open TLS connections to a Cloudflare-fronted IP (`104.21.77.91`), which isn't Google. Without root we couldn't tie them to a specific process.

## Indicators of compromise

SHA-256 of the stock files (build `VK-D39G-4M-XJ9RQ0C4-3+16-BML-S26 ULTRA Mini-PT-20260227`):

```
f3eed959fab6b12cba0505567d81d3d6f16ce62202cd73ca53118c93aa423feb  /system/lib64/libandroid_runtime.so
aaf6f3fa95fccb888b320feed1479f0e552c80834577f816c40f03e281c74f1e  /system/lib/libandroid_runtime.so
3f65d528ce42e939ab46d1a7a8d2cab18e91ee34b742f5201f440e14112ec978  /system/framework/framework.jar
2d55940e2933ef9fe0b6ffaa2ede9ab96a754f7aa342845201bf0eb94fa39fb6  /system/framework/services.jar
5e1f993880c812c86de20a2f937b78cde847e486437af45c7903e2f70909fdd6  /system/app/ArtificialSwitchQmx/ArtificialSwitchQmx.apk
1137bcbd37e87500771d190178ff2e852a459886cc14f44f7a7f67c7bef08de9  /system/app/WriteIMEI/WriteIMEI.apk
dd37f33677d268e26bcffd0bcd886cb9d46ead87d7c8f90c49961007353bde52  /system/app/Rsota/Rsota.apk
```

Package names: `com.example.artificialswitch`, `com.lys.writeimei`, `com.abfota.systemUpdate`.

Samples are not distributed in this repository.

## What survives a system replacement

Everything above lives in `/system`, which the GSI fully replaces. What stays stock: `vendor`, kernel (`boot`), preloader/LK, TEE and modem firmware. A string scan of `vendor` for the indicators above came back clean, but that doesn't prove it is trustworthy. Treat the phone as semi-trusted.
