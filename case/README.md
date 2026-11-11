# 3D-printed case

A two-part snap-on case for the S26 ULTRA Mini (48.3 × 99.6 × 11.8 mm, 9.4 mm corner radius). Tested on the real phone.

| File | Part |
|---|---|
| [`s26-mini-case_v8_bandeja.stl`](s26-mini-case_v8_bandeja.stl) | **Tray**: back and walls. The phone drops in from the top. |
| [`s26-mini-case_v8_moldura.stl`](s26-mini-case_v8_moldura.stl) | **Bezel**: the lip over the screen. Its stem slides into a channel in the tray walls and holds by friction. |

Outside size: 54.8 × 106.4 × 14.8 mm. Openings for the camera, both side buttons and USB-C.

## Printing

Both files are already oriented for printing (the bezel lies face down). Print them separately, not assembled.

- PLA, 0.4 mm nozzle, 0.16 mm layers
- No supports, 4 walls, 40% gyroid infill
- Bed 60 °C, aux fan off
- About 42 min / 19 g for the tray and 14 min / 4 g for the bezel (Anycubic Kobra S1)

## Fitting

Put the phone in the tray, then press the bezel down all the way around until it sits flat on the walls. The fit is tight on purpose: 0.05 mm per side between the stem and the channel.

## Source

The case is generated in Python (comments and labels are in Portuguese). All measurements are constants at the top of `s26case.py`.

| File | What it does |
|---|---|
| [`s26case.py`](s26case.py) | Builds two versions: **v8** (tray + bezel, the STLs above) and **v9** (a slide-on version with a dovetail rail, not tested) |
| [`desenho_tecnico.py`](desenho_tecnico.py) | Draws the technical sheets (views, sections and a legend with a code for each dimension) as PDF and PNG |
| [`s26-mini-case_v8_desenho-tecnico.pdf`](s26-mini-case_v8_desenho-tecnico.pdf) | Technical drawing of v8 (the printed case) |

Tested with Python 3.12 (`manifold3d` has no wheel for Python 3.15 yet).

```sh
python3.12 -m venv venv
./venv/bin/pip install trimesh manifold3d numpy scipy networkx shapely matplotlib
cd case
../venv/bin/python s26case.py .      # writes the v8 and v9 STLs here
../venv/bin/python desenho_tecnico.py
```
