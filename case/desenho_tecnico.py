"""Desenho técnico da capinha: vistas, cortes e a legenda com o código de cada parte.

Gera uma folha para cada versão (v8 bandeja + moldura, v9 deslize), montando a geometria direto do s26case.py.
Basta rodar de novo depois de mudar qualquer medida.
"""
import numpy as np
import trimesh
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.gridspec import GridSpec

import s26case as P

ORANGE, BLUE, DARK = "#d9661f", "#2f6db5", "0.15"
Y_TOP = P.PH / 2            # topo do celular
PHONE = P.to_tm(P.phone())


# ---------------- helpers de desenho ----------------
def outline(ax, mesh, origin, normal, cols, color, mirror=False, lw=1.1, ls="-"):
    s = mesh.section(plane_origin=origin, plane_normal=normal)
    if s is None:
        return
    for e in s.discrete:
        u, v = e[:, cols[0]], e[:, cols[1]]
        ax.plot(-u if mirror else u, v, color=color, lw=lw, ls=ls)


def phone(ax, origin, normal, cols):
    outline(ax, PHONE, origin, normal, cols, "0.55", lw=0.7, ls=(0, (3, 2)))


def dim(ax, p1, p2, off, text, horiz=True, fs=7):
    (x1, y1), (x2, y2) = p1, p2
    if horiz:
        yy = y1 + off
        ax.plot([x1, x1], [y1, yy], lw=0.4, c="0.4"); ax.plot([x2, x2], [y2, yy], lw=0.4, c="0.4")
        ax.annotate("", (x1, yy), (x2, yy), arrowprops=dict(arrowstyle="<->", lw=0.6, color="0.2", shrinkA=0, shrinkB=0))
        ax.text((x1 + x2) / 2, yy + (0.08 * abs(off) + 0.3 if off >= 0 else -0.08 * abs(off) - 0.3), text, ha="center",
                va="bottom" if off >= 0 else "top", fontsize=fs)
    else:
        xx = x1 + off
        ax.plot([x1, xx], [y1, y1], lw=0.4, c="0.4"); ax.plot([x2, xx], [y2, y2], lw=0.4, c="0.4")
        ax.annotate("", (xx, y1), (xx, y2), arrowprops=dict(arrowstyle="<->", lw=0.6, color="0.2", shrinkA=0, shrinkB=0))
        ax.text(xx + (0.08 * abs(off) + 0.3 if off >= 0 else -0.08 * abs(off) - 0.3), (y1 + y2) / 2, text,
                ha="left" if off >= 0 else "right", va="center", fontsize=fs, rotation=90)


def setup(ax, title, xl, yl):
    ax.set_aspect("equal"); ax.set_xlim(*xl); ax.set_ylim(*yl)
    ax.set_title(title, fontsize=10, fontweight="bold", loc="left")
    ax.tick_params(labelsize=6); ax.grid(alpha=0.15)


def make_tag(col):
    def tag(ax, code, target, at, fs=7.5):
        c = col.get(code[0], DARK)
        ax.annotate(code, target, at, fontsize=fs, fontweight="bold", color="white", ha="center", va="center",
                    bbox=dict(boxstyle="circle,pad=0.25" if len(code) < 3 else "round,pad=0.25", fc=c, ec="none"),
                    arrowprops=dict(arrowstyle="-", lw=0.7, color=c, shrinkA=6, shrinkB=0))
        ax.plot(*target, "o", ms=2.2, color=c)
    return tag


def cam_fmt(i):
    (a, b), (c, d), r = P.CAM_WINDOWS[i]
    return f"{a}–{b} / {c}–{d} / R{r}"


def common_legend(d):
    return [
        ("Geral", None, None),
        ("G1", "Raio do canto (visto de frente)", f"{P.R_PHONE} mm"),
        ("G2", "Folga lateral, por lado", f"{P.C_XY} mm"),
        ("G3", "Folga topo/base, por lado", f"{P.C_Y} mm"),
        ("G4", "Folga na espessura", f"{P.C_Z} mm"),
        ("G5", "Parede lateral (total)", f"{d.wall:.1f} mm"),
        ("G6", "Parede do topo e da base (total)", f"{d.end:.1f} mm"),
        ("G7", "Traseira (espessura)", f"{P.T_BACK} mm"),
        ("G8", "Borda sobre a tela (espessura)", f"{P.T_LIP} mm"),
        ("G9", "Borda sobre a tela, avanço nas laterais", f"{P.LIP_SIDE} mm"),
        ("G10", "Borda sobre a tela, avanço topo/base", f"{P.LIP_END} mm"),
    ]


CUTOUT_LEGEND = [
    ("Recortes", None, None),
    ("C1", "Câmera, 3 lentes (b/t/R)  [era B3]", cam_fmt(0)),
    ("C2", "Câmera, lentes peq. (b/t/R)  [era B4]", cam_fmt(1)),
    ("F1", "Furo lateral direita (do topo / alt.)  [era B5]", f"{P.BTN_RIGHT[0]}–{P.BTN_RIGHT[1]} / {P.BTN_H}"),
    ("F2", "Furo lateral esquerda (do topo / alt.)  [era B6]", f"{P.BTN_LEFT[0]}–{P.BTN_LEFT[1]} / {P.BTN_H}"),
    ("F3", "USB-C / alto-falante (larg. × alt.)  [era A3]", f"{P.USB_W} × {P.USB_H} mm"),
]


# ---------------- vistas comuns ----------------
def front_view(ax, parts, d, tag, y_cut=None):
    hw, hh = d.OW / 2, d.OH / 2
    for mesh, c in parts:
        outline(ax, mesh, (0, 0, P.ZC), (0, 0, 1), (0, 1), c, lw=0.6)
        outline(ax, mesh, (0, 0, P.Z_TOP - 0.6), (0, 0, 1), (0, 1), c)
    phone(ax, (0, 0, P.ZC), (0, 0, 1), (0, 1))
    dim(ax, (-hw, hh), (hw, hh), 4, f"{d.OW:.1f}")
    dim(ax, (hw, -hh), (hw, hh), 9, f"{d.OH:.1f}", horiz=False)
    k = 1 - 1 / np.sqrt(2)
    tag(ax, "G1", (-P.XI + P.R_CAV * k, P.CH / 2 - P.R_CAV * k), (-hw - 8, hh + 8))
    tag(ax, "G3", (-15, P.PH / 2 + P.C_Y / 2), (-24, hh + 7))
    tag(ax, "G9", (-P.PW / 2 + P.LIP_SIDE / 2, 20), (-12, 26))
    tag(ax, "G10", (8, P.PH / 2 - P.LIP_END / 2), (10, 40))
    tag(ax, "G5", (hw - d.wall / 2, 30), (hw + 7, 36))
    tag(ax, "G6", (0, -hh + d.end / 2), (-8, -hh - 7))
    if y_cut is not None:
        ax.axhline(y_cut, color="0.3", ls="-.", lw=0.7)
        tag(ax, "G11", (-hw, y_cut), (-hw - 8, y_cut - 8))
    setup(ax, "1  Vista frontal (lado da tela)", (-hw - 14, hw + 14), (-hh - 12, hh + 12))


def back_view(ax, parts, d, tag, y_cut=None):
    hw, hh = d.OW / 2, d.OH / 2
    for mesh, c in parts:
        outline(ax, mesh, (0, 0, 0.4), (0, 0, 1), (0, 1), c, mirror=True)
    if y_cut is not None:
        ax.axhline(y_cut, color="0.3", ls="-.", lw=0.7)
    (e0, e1), (t0, t1), _ = P.CAM_WINDOWS[0]
    (f0, f1), (u0, u1), _ = P.CAM_WINDOWS[1]
    L = -P.PW / 2  # borda esquerda do celular na vista traseira
    dim(ax, (L, Y_TOP), (L + e0, Y_TOP), 5, f"{e0}", fs=6)
    dim(ax, (L, Y_TOP), (L + e1, Y_TOP), 9, f"{e1}", fs=6)
    dim(ax, (L, Y_TOP), (L + f1, Y_TOP), 13, f"{f1}", fs=6)
    dim(ax, (L, Y_TOP - t1), (L, Y_TOP), -5, f"{t1}", horiz=False, fs=6)
    dim(ax, (L + f1, Y_TOP - u1), (L + f1, Y_TOP), 14, f"{u1}", horiz=False, fs=6)
    tag(ax, "C1", (L + (e0 + e1) / 2, Y_TOP - (t0 + t1) / 2), (L - 8, Y_TOP - 30))
    tag(ax, "C2", (L + (f0 + f1) / 2, Y_TOP - (u0 + u1) / 2), (hw + 7, Y_TOP - 32))
    tag(ax, "G7", (5, -30), (hw + 7, -30))
    ax.text(0, -hh - 8, "b/t = distância da borda esquerda / do topo do celular, R = raio", fontsize=6, color="0.35", ha="center")
    setup(ax, "2  Vista traseira", (-hw - 14, hw + 14), (-hh - 12, hh + 18))


def side_views(fig, gs, parts, d, tag, y_cut=None):
    hw, hh = d.OW / 2, d.OH / 2
    xs = hw - 0.35
    for col, sg, code, btn, title in ((2, 1, "F1", P.BTN_RIGHT, "3  Lateral direita"), (3, -1, "F2", P.BTN_LEFT, "4  Lateral esquerda")):
        ax = fig.add_subplot(gs[0, col])
        for mesh, c in parts:
            outline(ax, mesh, (sg * xs, 0, 0), (1, 0, 0), (2, 1), c)
        if y_cut is not None:
            ax.axhline(y_cut, color="0.3", ls="-.", lw=0.7)
        by0, by1 = Y_TOP - btn[1], Y_TOP - btn[0]
        dim(ax, (P.Z_TOP, Y_TOP), (P.Z_TOP, by1), 4, f"{btn[0]}", horiz=False, fs=6)
        dim(ax, (P.Z_TOP, by0), (P.Z_TOP, by1), 9, f"{by1 - by0:.1f}", horiz=False, fs=6)
        dim(ax, (P.ZC - P.BTN_H / 2, by0), (P.ZC + P.BTN_H / 2, by0), -5, f"{P.BTN_H}", fs=6)
        tag(ax, code, (P.ZC, (by0 + by1) / 2), (-9, (by0 + by1) / 2))
        ax.text(0, hh + 3, "traseira", fontsize=6, ha="center"); ax.text(P.Z_TOP, hh + 3, "tela", fontsize=6, ha="center")
        setup(ax, title, (-14, P.Z_TOP + 16), (-hh - 12, hh + 12))


def base_view(ax, parts, d, tag):
    hw, hh = d.OW / 2, d.OH / 2
    for mesh, c in parts:
        outline(ax, mesh, (0, -hh + 0.6, 0), (0, 1, 0), (0, 2), c)
    dim(ax, (-P.USB_W / 2, P.ZC + P.USB_H / 2), (P.USB_W / 2, P.ZC + P.USB_H / 2), 6, f"{P.USB_W}", fs=6)
    dim(ax, (P.USB_W / 2, P.ZC - P.USB_H / 2), (P.USB_W / 2, P.ZC + P.USB_H / 2), 6, f"{P.USB_H}", horiz=False, fs=6)
    dim(ax, (-hw, 0), (-hw, P.Z_TOP), -3, f"{P.Z_TOP:.1f}", horiz=False, fs=6)
    tag(ax, "F3", (0, P.ZC), (0, -8))
    setup(ax, "5  Base (USB-C)", (-hw - 8, hw + 12), (-14, P.Z_TOP + 12))


def legend(ax, rows, col, note):
    ax.axis("off")
    ax.set_title("Legenda  ·  código → medida atual", fontsize=10, fontweight="bold", loc="left")
    y = 0.98
    for code, name, val in rows:
        if name is None:
            y -= 0.008
            ax.text(0.0, y, code, fontsize=10, fontweight="bold", color=DARK, transform=ax.transAxes)
            y -= 0.03
            continue
        c = col.get(code[:1], DARK)
        if code:
            ax.text(0.02, y, code, fontsize=8, fontweight="bold", color="white", transform=ax.transAxes, ha="center", va="center",
                    bbox=dict(boxstyle="round,pad=0.25", fc=c, ec="none"))
        ax.text(0.075, y, name, fontsize=7.5, transform=ax.transAxes, va="center")
        ax.text(0.99, y, val, fontsize=8, fontweight="bold", transform=ax.transAxes, va="center", ha="right")
        y -= 0.029
    ax.text(0.0, y - 0.01, note, fontsize=7, color="0.35", transform=ax.transAxes, va="top")


def new_sheet(version, subtitle, codes):
    fig = plt.figure(figsize=(16.5, 11.7))
    gs = GridSpec(2, 6, figure=fig, width_ratios=[1, 1, 0.5, 0.5, 0.12, 1.25], height_ratios=[1.25, 1],
                  left=0.03, right=0.99, top=0.9, bottom=0.04, wspace=0.18, hspace=0.18)
    fig.suptitle(f"S26 Ultra Mini — capinha v{version}  ·  desenho técnico (mm)", fontsize=14, fontweight="bold", x=0.03, ha="left")
    fig.text(0.03, 0.94, subtitle + f"   ·   use os códigos ({codes}...) para pedir ajustes", fontsize=8, color="0.3")
    return fig, gs


def save(fig, version):
    out = f"s26-mini-case_v{version}_desenho-tecnico"
    fig.savefig(out + ".pdf")
    fig.savefig(out + ".png", dpi=110)
    plt.close(fig)
    print(out)


NOTE = ("Ainda estimados (não medidos no celular): F3 (USB-C) e C1/C2\n"
        f"(câmera, por foto); borda frente/trás do celular R_EDGE = {P.R_EDGE} mm\n"
        f"(a cavidade usa {P.R_EDGE_CAV} para a quina do celular não calçar).")


# =====================================================================
# v8: bandeja + moldura
# =====================================================================
def sheet_moldura():
    m = P.build_moldura(); d = m.d
    tray, frame = P.to_tm(m.tray), P.to_tm(m.frame)
    col = {"B": ORANGE, "M": BLUE}
    tag = make_tag(col)
    parts = ((tray, ORANGE), (frame, BLUE))
    fig, gs = new_sheet(P.V_MOLDURA, "laranja = bandeja (traseira + paredes, o celular entra por cima)   ·   "
                                     "azul = moldura (borda sobre a tela, encaixa por cima no canal)", "G3, B2, M2, C1")
    hw, hh = d.OW / 2, d.OH / 2
    front_view(fig.add_subplot(gs[0, 0]), parts, d, tag)
    back_view(fig.add_subplot(gs[0, 1]), parts, d, tag)
    side_views(fig, gs, parts, d, tag)
    base_view(fig.add_subplot(gs[1, 2:4]), parts, d, tag)

    y_c = -30.0
    x_in, x_gr, x_out = P.XI + P.W_IN, P.XI + P.W_IN + P.W_GR, hw
    z_gr = P.Z_W - P.GR_D
    ax = fig.add_subplot(gs[1, 0])
    for mesh, c in parts:
        outline(ax, mesh, (0, y_c, 0), (0, 1, 0), (0, 2), c)
    phone(ax, (0, y_c, 0), (0, 1, 0), (0, 2))
    ax.plot([-hw - 4, hw + 4], [P.Z_W] * 2, color="0.5", lw=0.4, ls=":")
    tag(ax, "M1", (x_in + P.W_GR / 2, z_gr + 2), (hw + 7, z_gr))
    tag(ax, "B2", (x_in + P.W_GR / 2, z_gr + 0.15), (hw + 7, z_gr - 5))
    tag(ax, "G8", (P.XI - 1, P.Z_TOP - P.T_LIP / 2), (P.XI - 8, P.Z_TOP + 6))
    tag(ax, "G7", (0, P.T_BACK / 2), (0, -6))
    tag(ax, "G4", (-10, P.T_BACK + P.C_Z / 2 + P.PT + P.C_Z / 4), (-16, P.Z_TOP + 6))
    ax.text(-10, P.ZC, "celular", fontsize=7, ha="center", va="center", color="0.45")
    setup(ax, f"6  Corte transversal (y = {y_c:.0f})", (-hw - 8, hw + 12), (-10, P.Z_TOP + 10))

    ax = fig.add_subplot(gs[1, 1])
    for mesh, c in parts:
        outline(ax, mesh, (0, y_c, 0), (0, 1, 0), (0, 2), c, lw=1.4)
    phone(ax, (0, y_c, 0), (0, 1, 0), (0, 2))
    zt = P.Z_TOP + 1.2
    dim(ax, (P.XI, zt), (x_in, zt), 0.0, "", fs=6); ax.text((P.XI + x_in) / 2, zt + 0.25, f"{P.W_IN}", fontsize=6, ha="center")
    dim(ax, (x_in, zt + 1.6), (x_gr, zt + 1.6), 0.0, "", fs=6); ax.text((x_in + x_gr) / 2, zt + 1.85, f"{P.W_GR}", fontsize=6, ha="center")
    dim(ax, (x_gr, zt), (x_out, zt), 0.0, "", fs=6); ax.text((x_gr + x_out) / 2, zt + 0.25, f"{P.W_OUT}", fontsize=6, ha="center")
    dim(ax, (x_out, z_gr), (x_out, P.Z_W), 1.2, f"{P.GR_D}", horiz=False, fs=6)
    dim(ax, (x_in, z_gr + P.STEM_UP), (x_in, P.Z_W), -3.2, f"{P.GR_D - P.STEM_UP:.1f}", horiz=False, fs=6)
    tag(ax, "B1", (P.XI + P.W_IN / 2, 6), (P.XI - 2.5, 4))
    tag(ax, "B2", (x_in + P.W_GR / 2, z_gr + 0.1), (x_in + 0.6, z_gr - 3))
    tag(ax, "B3", (x_gr + P.W_OUT / 2, 6), (x_out + 2.2, 4))
    tag(ax, "M1", (x_in + P.W_GR / 2, P.Z_W - 1.5), (x_out + 2.4, P.Z_W - 3))
    tag(ax, "M2", (x_in + P.FIT_M / 2, P.Z_W - 3.5), (P.XI - 2.5, P.Z_W - 5))
    ax.annotate("", (hw + 1.2, P.Z_TOP + 0.3), (hw + 1.2, P.Z_TOP + 3.0), arrowprops=dict(arrowstyle="->", lw=1, color=BLUE))
    ax.text(hw + 1.6, P.Z_TOP + 1.8, "moldura\ndesce", fontsize=6, color=BLUE, va="center")
    setup(ax, "7  Detalhe do encaixe (canal + haste)", (P.XI - 4, hw + 4.5), (-0.5, P.Z_TOP + 5.5))

    rows = common_legend(d) + [
        ("Bandeja", None, None),
        ("B1", "Parede interna (encosta no celular)", f"{P.W_IN} mm"),
        ("B2", "Canal: largura × profundidade", f"{P.W_GR} × {P.GR_D} mm"),
        ("B3", "Parede externa", f"{P.W_OUT} mm"),
        ("Moldura", None, None),
        ("M1", "Haste: largura × altura", f"{P.W_GR - 2 * P.FIT_M:.2f} × {P.GR_D - P.STEM_UP:.1f} mm"),
        ("M2", "Folga haste ↔ canal, por lado (neg. = aperto)", f"{P.FIT_M} mm"),
        ("M3", "Recuo da haste em volta dos furos", f"{P.NOTCH} mm"),
    ] + CUTOUT_LEGEND
    lax = fig.add_subplot(gs[:, 5])
    legend(lax, rows, col, NOTE + "\n\nA moldura segura só por atrito, como o snapcase do iPhone.\n"
                                  "Se ficar solta, diminuir M2; se não entrar, aumentar.\n"
                                  "A moldura sai virada para imprimir, com a frente na mesa.")
    for t in lax.texts:
        if t.get_text() in ("Bandeja", "Moldura"):
            t.set_color(ORANGE if t.get_text() == "Bandeja" else BLUE)
    save(fig, P.V_MOLDURA)


# =====================================================================
# v9: deslize (trilho rabo de andorinha)
# =====================================================================
def sheet_trilho():
    t = P.build_trilho(with_ears=False); d = t.d
    base, top = P.to_tm(t.base), P.to_tm(t.top)
    col = {"A": ORANGE, "B": BLUE}
    tag = make_tag(col)
    parts = ((base, ORANGE), (top, BLUE))
    fig, gs = new_sheet(P.V_TRILHO, "laranja = parte A (base, com a traseira inteira)   ·   "
                                    "azul = parte B (topo, desliza pela borda da traseira até a emenda G11)", "G3, A1, B2, C1")
    hw, hh = d.OW / 2, d.OH / 2
    ax = fig.add_subplot(gs[0, 0])
    front_view(ax, parts, d, tag, y_cut=P.Y_CUT)
    ax.text(0, -38, "A", fontsize=22, color=ORANGE, ha="center", alpha=0.35, fontweight="bold")
    ax.text(0, 12, "B", fontsize=22, color=BLUE, ha="center", alpha=0.35, fontweight="bold")
    ax = fig.add_subplot(gs[0, 1])
    back_view(ax, parts, d, tag, y_cut=P.Y_CUT)
    ax.text(-hw + 1, 0, "trilho", fontsize=6, color=BLUE, rotation=90, va="center")
    side_views(fig, gs, parts, d, tag, y_cut=P.Y_CUT)
    base_view(fig.add_subplot(gs[1, 2:4]), parts, d, tag)

    y_c = -8.0
    xa, xt = P.XI - P.DT_S, P.XI - P.DT_S + P.DT_H
    ax = fig.add_subplot(gs[1, 0])
    for mesh, c in parts:
        outline(ax, mesh, (0, y_c, 0), (0, 1, 0), (0, 2), c)
    phone(ax, (0, y_c, 0), (0, 1, 0), (0, 2))
    tag(ax, "A1", (-xt + 0.1, P.DT_H + P.DT_TIP / 2), (-hw - 4, -5))
    tag(ax, "B1", (xa + 0.3, 0.15), (hw + 6, -4))
    tag(ax, "G8", (P.XI - 1, P.Z_TOP - P.T_LIP / 2), (P.XI - 8, P.Z_TOP + 6))
    tag(ax, "G7", (0, P.T_BACK / 2), (0, -6))
    tag(ax, "G4", (-10, P.T_BACK + P.C_Z / 2 + P.PT + P.C_Z / 4), (-16, P.Z_TOP + 6))
    ax.text(-10, P.ZC, "celular", fontsize=7, ha="center", va="center", color="0.45")
    setup(ax, f"6  Corte transversal (y = {y_c:.0f}, na parte B)", (-hw - 8, hw + 12), (-10, P.Z_TOP + 10))

    ax = fig.add_subplot(gs[1, 1])
    for mesh, c in parts:
        outline(ax, mesh, (0, y_c, 0), (0, 1, 0), (0, 2), c, lw=1.5)
    phone(ax, (0, y_c, 0), (0, 1, 0), (0, 2))
    dim(ax, (xa, -0.35), (P.XI, -0.35), -0.35, f"{P.DT_S}", fs=6)
    dim(ax, (P.XI, -0.35), (xt, -0.35), -1.2, f"{P.DT_H - P.DT_S:.1f}", fs=6)
    dim(ax, (xt + 0.25, 0), (xt + 0.25, P.T_BACK), 1.6, f"{P.T_BACK}", horiz=False, fs=6)
    ax.text(xa + 0.1, P.DT_H * 0.35, "45°", fontsize=7, ha="left")
    tag(ax, "A1", (xt - 0.05, P.DT_H + P.DT_TIP / 2), (xt + 1.8, P.DT_H + 2.2))
    tag(ax, "B1", (xa + 0.35, 0.12), (xa - 1.2, 1.9))
    tag(ax, "B2", (xt - 0.3, P.DT_H + P.DT_TIP + 0.25), (P.XI - 1.2, 3.9))
    tag(ax, "G7", (P.XI - 1.6, P.T_BACK / 2), (P.XI - 2.2, -1.9))
    ax.text(hw - 0.6, 3.6, "B corre ao longo\ndo trilho (em y)", fontsize=6, color=BLUE, ha="right", va="center")
    setup(ax, "7  Detalhe do trilho (lado direito)", (P.XI - 2.7, hw + 0.8), (-2.8, 4.8))

    rows = common_legend(d) + [
        ("G11", "Emenda A | B nas paredes (da base do celular)", f"{P.PH / 2 + P.Y_CUT:.1f} mm"),
        ("Parte A (base)", None, None),
        ("A1", "Flecha da traseira: avanço além da parede / ponta", f"{P.DT_H - P.DT_S:.1f} / {P.DT_TIP} mm"),
        ("Parte B (topo)", None, None),
        ("B1", "Gancho por baixo da traseira (a 45°)", f"{P.DT_S} mm"),
        ("B2", "Folga do trilho (perpendicular)", f"{P.DT_GAP} mm"),
        ("B3", "Folga na emenda G11 / na ponta da traseira", f"{P.GAP_CUT} / {P.G_END} mm"),
    ] + CUTOUT_LEGEND
    lax = fig.add_subplot(gs[:, 5])
    legend(lax, rows, col, NOTE + "\n\nB segura por atrito no trilho, como o Jelly Star.\n"
                                  "Se ficar solta, diminuir B2; se não correr, aumentar.\n"
                                  "A janela da câmera é aberta em cima, para o calombo\n"
                                  "das lentes entrar deslizando junto com o celular.\n"
                                  "B sai com orelhas nas pontas da emenda: cortar antes de montar.")
    for tx in lax.texts:
        if tx.get_text().startswith("Parte"):
            tx.set_color(ORANGE if "A (" in tx.get_text() else BLUE)
    save(fig, P.V_TRILHO)


if __name__ == "__main__":
    sheet_moldura()
    sheet_trilho()
