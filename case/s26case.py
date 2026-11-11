"""S26 Ultra Mini - capinha. Gera duas versões com encaixes diferentes, a partir das referências de iPhone:

  v8  encaixe por cima: bandeja + moldura. A parede da bandeja tem um canal; a moldura tem uma haste em T que entra
      nesse canal e segura por atrito no contorno inteiro (como o Iphone_14pro_snapcase).
  v9  deslize: a base tem a traseira inteira, com a borda em ponta de flecha (rabo de andorinha a 45°). A tampa de
      cima corre por essa borda até encostar na base (como o Jelly Star slider).

Referencial (visto pela FRENTE, tela para cima):
  +x = direita, +y = topo do celular (lado das câmeras), -y = base (USB-C)
  z = 0 é a face externa da traseira da capa; a tela aponta para +z.
As peças saem já na orientação de impressão. A moldura da v8 sai virada, com a face da frente na mesa.
Versões anteriores (v1–v7, com a língua e a casca sobrepostas) ficam em antigos/.
"""
import math, sys
from types import SimpleNamespace

V_MOLDURA = 8
V_TRILHO = 9
import numpy as np
import trimesh
from manifold3d import Manifold, CrossSection, JoinType

# ---------------- medidas do celular (paquímetro) ----------------
PW, PH, PT = 48.3, 99.6, 11.8      # largura, altura (99.6 remedido no paquímetro), espessura
R_PHONE = 9.4                      # raio dos cantos: anel v3 de 9.4 ficou justo (8.7 entrou folgado)
R_EDGE = 1.2                       # arredondamento da borda frente/trás do celular (ESTIMADO)

# ---------------- capa (comum às duas) ----------------
C_XY = 0.15        # folga por lado em x
C_Y = 0.3          # folga por lado em y: v2–v6 0.1 faltou ~1 mm com a trava da v6; v7 0.6; v8/v9 0.3 sem trava no comprimento
C_Z = 0.20         # folga na espessura
R_EDGE_CAV = 0.6   # arredondamento do fundo da cavidade: menor que o do celular, para a quina dele nunca calçar
T_BACK = 1.6       # traseira
T_LIP = 1.2        # espessura da borda sobre a tela (na v8, o topo da moldura)
LIP_SIDE = 1.8     # quanto a borda avança sobre a tela nas laterais
LIP_END = 2.6      # ... no topo e na base

# ---------------- v8: bandeja + moldura ----------------
W_IN = 0.8         # parede interna da bandeja (encosta no celular)
W_GR = 1.3         # canal na parede, onde entra a haste da moldura
W_OUT = 1.0        # parede externa da bandeja
GR_D = 5.0         # profundidade do canal, a partir do topo da parede
FIT_M = 0.05       # folga por lado entre haste e canal (snapcase usa 0; negativo = interferência)
STEM_UP = 0.3      # a haste é mais curta que o canal, para a moldura assentar no topo das paredes
NOTCH = 0.5        # recuo da haste em volta dos furos (botões e USB), para não aparecer dentro deles

# ---------------- v9: deslize (trilho rabo de andorinha) ----------------
T_WALL9 = 2.4      # parede lateral
T_END9 = 2.0       # parede do topo e da base
Y_CUT = -16.0      # onde as paredes da base terminam e começam as da tampa (y do celular)
GAP_CUT = 0.15     # folga entre as paredes da base e da tampa, na emenda
DT_S = 0.4         # quanto o gancho da tampa entra por baixo da traseira (a partir da face interna da parede)
DT_H = 0.9         # altura da face de baixo da flecha (45°): a flecha avança DT_H - DT_S além da face interna
DT_TIP = 0.2       # ponta reta da flecha
DT_GAP = 0.15      # folga do trilho (medida na perpendicular). O Jelly segura por atrito, com ~0.1 de aperto
G_END = 0.3        # folga entre a ponta da traseira e a parede do topo, para a emenda G11 ser o batente

# ---------------- recortes (ESTIMADOS a partir da foto) ----------------
# câmera: medido na foto da VISTA TRASEIRA a partir da borda esquerda / topo, com ~1.2 mm de margem
# (da borda: início, fim), (do topo: início, fim), raio
CAM_WINDOWS = [
    ((2.9, 18.6), (1.8, 43.7), 7.8),    # C1: coluna com as 3 lentes grandes
    ((17.3, 27.7), (4.8, 27.8), 5.0),   # C2: 2 lentes pequenas + flash
]
# botões: janelas nas laterais, medidas a partir do topo (lados como na vista frontal)
BTN_RIGHT = (18.25, 45.25)    # F1: direito (+x), 27 mm
BTN_LEFT = (22.25, 51.25)     # F2: esquerdo (-x), 29 mm
BTN_H = 7.0
# base: USB-C + alto-falante + microfone
USB_W, USB_H = 24.0, 7.5      # F3

# ---------------- orelhas anti-empeno (só para imprimir) ----------------
# discos finos nas pontas livres de uma peça em U, que levantavam da mesa na v6. Cortar rente com estilete.
EAR_R = 5.0
EAR_T = 0.32       # 1ª camada de 0.2 + mais uma

# ---------------- derivadas comuns ----------------
CW, CH = PW + 2 * C_XY, PH + 2 * C_Y       # cavidade
CT = PT + C_Z
R_CAV = R_PHONE + C_XY
XI = CW / 2                                # face interna da parede
ZC = T_BACK + CT / 2                       # meio da altura do celular
Z_W = T_BACK + CT                          # altura da frente do celular (+ folga): topo das paredes da bandeja na v8
Z_TOP = Z_W + T_LIP
EPS = 0.01
SEG = 40


def dims(version):
    """medidas externas de cada versão."""
    wall = W_IN + W_GR + W_OUT if version == V_MOLDURA else T_WALL9
    end = W_IN + W_GR + W_OUT if version == V_MOLDURA else T_END9
    ow, oh = CW + 2 * wall, CH + 2 * end
    return SimpleNamespace(wall=wall, end=end, OW=ow, OH=oh, R_OUT=R_CAV + wall, XO=ow / 2)


def rrect_pts(w, h, r, n=SEG // 4):
    r = max(min(r, w / 2 - EPS, h / 2 - EPS), 0.05)
    pts = []
    for cx, cy, a0 in ((w / 2 - r, h / 2 - r, 0), (-w / 2 + r, h / 2 - r, 90),
                       (-w / 2 + r, -h / 2 + r, 180), (w / 2 - r, -h / 2 + r, 270)):
        for i in range(n + 1):
            a = math.radians(a0 + 90 * i / n)
            pts.append((cx + r * math.cos(a), cy + r * math.sin(a)))
    return pts


def rounded_box(w, h, r, z0, z1, rb, rt, n=8):
    pts = []
    for i in range(n + 1):
        phi = math.pi / 2 * i / n
        for (zz, rr, sgn) in ((z0, rb, 1), (z1, rt, -1)):
            ins = rr * (1 - math.sin(phi))
            z = zz + sgn * rr * (1 - math.cos(phi))
            pts += [(x, y, z) for x, y in rrect_pts(w - 2 * ins, h - 2 * ins, r - ins)]
    return Manifold.hull_points(pts)


def prism(pts2d, z0, z1):
    return CrossSection([pts2d]).extrude(z1 - z0).translate((0, 0, z0))


def ring(off0, off1, z0, z1):
    """faixa em volta da cavidade, entre os afastamentos off0 e off1 da face interna."""
    return (prism(rrect_pts(CW + 2 * off1, CH + 2 * off1, R_CAV + off1), z0, z1)
            - prism(rrect_pts(CW + 2 * off0, CH + 2 * off0, R_CAV + off0), z0 - 1, z1 + 1))


def box(x0, x1, y0, y1, z0, z1):
    return Manifold.cube((x1 - x0, y1 - y0, z1 - z0)).translate((x0, y0, z0))


def stadium_x(y0, y1, zc, h, x0, x1):
    """furo alongado atravessando a parede na direção x."""
    r = h / 2
    c = lambda y: Manifold.cylinder(x1 - x0, r, circular_segments=SEG).rotate((0, 90, 0)).translate((x0, y, zc))
    return Manifold.batch_hull([c(y0 + r), c(y1 - r)])


def stadium_y(xw, zc, h, y0, y1):
    r = h / 2
    c = lambda x: Manifold.cylinder(y1 - y0, r, circular_segments=SEG).rotate((-90, 0, 0)).translate((x, y0, zc))
    return Manifold.batch_hull([c(-xw / 2 + r), c(xw / 2 - r)])


def phone():
    return rounded_box(PW, PH, R_PHONE, T_BACK + C_Z / 2, T_BACK + C_Z / 2 + PT, R_EDGE, R_EDGE)


def camera(extend_to=None):
    """janelas da câmera na traseira. extend_to abre a janela até esse y (v9: o calombo entra deslizando)."""
    cam = Manifold()
    for (e0, e1), (t0, t1), r in CAM_WINDOWS:
        cx0, cx1 = PW / 2 - e1, PW / 2 - e0
        cy0, cy1 = PH / 2 - t1, PH / 2 - t0
        cam += prism(rrect_pts(cx1 - cx0, cy1 - cy0, r), -1, T_BACK + 1).translate(((cx0 + cx1) / 2, (cy0 + cy1) / 2, 0))
        if extend_to is not None:
            cam += box(cx0, cx1, (cy0 + cy1) / 2, extend_to, -1, T_BACK + 1)
    return cam


def holes(d, grow=0.0):
    """botões nas laterais e USB-C na base. grow aumenta tudo (para recuar a haste da moldura)."""
    btn = (stadium_x(PH / 2 - BTN_RIGHT[1] - grow, PH / 2 - BTN_RIGHT[0] + grow, ZC, BTN_H + 2 * grow, 0, d.OW / 2 + 1)
           + stadium_x(PH / 2 - BTN_LEFT[1] - grow, PH / 2 - BTN_LEFT[0] + grow, ZC, BTN_H + 2 * grow, -d.OW / 2 - 1, 0))
    usb = stadium_y(USB_W + 2 * grow, ZC, USB_H + 2 * grow, -d.OH / 2 - 1, 0)
    return btn + usb


def ear(x, y, keep_out):
    return Manifold.cylinder(EAR_T, EAR_R, circular_segments=SEG).translate((x, y, 0)) - keep_out


# =====================================================================
# v8: bandeja + moldura
# =====================================================================
def build_moldura():
    d = dims(V_MOLDURA)
    cavity = rounded_box(CW, CH, R_CAV, T_BACK, Z_W + 5, R_EDGE_CAV, 0.0)
    tray = rounded_box(d.OW, d.OH, d.R_OUT, 0, Z_W, 0.8, 0.0) - cavity - ring(W_IN, W_IN + W_GR, Z_W - GR_D, Z_W + 1)
    tray = tray - camera() - holes(d)

    cap = rounded_box(d.OW, d.OH, d.R_OUT, Z_W, Z_TOP, 0.0, 1.2)
    cap = cap - prism(rrect_pts(PW - 2 * LIP_SIDE, PH - 2 * LIP_END, R_PHONE - LIP_SIDE), Z_W - 1, Z_TOP + 1)
    stem = ring(W_IN + FIT_M, W_IN + W_GR - FIT_M, Z_W - GR_D + STEM_UP, Z_W + EPS) - holes(d, NOTCH)
    frame = cap + stem
    return SimpleNamespace(d=d, tray=tray, frame=frame, phone=phone())


def frame_for_print(frame):
    """moldura virada (giro de 180° em y, sem espelhar), com a face da frente na mesa."""
    f = frame.rotate((0, 180, 0))
    return f.translate((0, 0, -f.bounding_box()[2]))


# =====================================================================
# v9: deslize (trilho rabo de andorinha na borda da traseira)
# =====================================================================
def dovetail_section(gap=0.0):
    """perfil da traseira da base (plano x-z): borda em ponta de flecha a 45° nos dois lados."""
    xa = XI - DT_S                   # borda da traseira na mesa
    xt = xa + DT_H                   # ponta da flecha
    xtop = xt - (T_BACK - DT_H - DT_TIP)
    half = [(xa - 1, -1), (xt, DT_H), (xt, DT_H + DT_TIP), (xtop, T_BACK)]
    pts = half + [(-x, z) for x, z in reversed(half)]
    cs = CrossSection([pts])
    return cs.offset(gap, JoinType.Miter) if gap else cs


def along_y(cs, y0, y1):
    """extruda um perfil do plano x-z ao longo de y, de y0 a y1."""
    m = cs.extrude(y1 - y0).rotate((90, 0, 0))      # (u, v, w) -> (u, -w, v)
    return m.translate((0, y1, 0))


def build_trilho(with_ears=True):
    d = dims(V_TRILHO)
    outer = rounded_box(d.OW, d.OH, d.R_OUT, 0, Z_TOP, 0.8, 1.2)
    cavity = rounded_box(CW, CH, R_CAV, T_BACK, T_BACK + CT, R_EDGE_CAV, R_EDGE_CAV)
    win = prism(rrect_pts(PW - 2 * LIP_SIDE, PH - 2 * LIP_END, R_PHONE - LIP_SIDE), T_BACK + 1, Z_TOP + 1)
    case = outer - cavity - win - camera() - holes(d)

    big = 200
    y_plate = CH / 2                                   # a traseira da base vai até a ponta da cavidade
    plate = along_y(dovetail_section(), Y_CUT - EPS, y_plate)
    base = (case ^ box(-big, big, -big, Y_CUT, -1, big)) + (case ^ plate)
    base = base - camera(extend_to=y_plate + 1)         # janela aberta em cima: o calombo entra deslizando

    plate_gap = along_y(dovetail_section(DT_GAP), Y_CUT - 1, y_plate + G_END)
    top = (case ^ box(-big, big, Y_CUT + GAP_CUT, big, -1, big)) - plate_gap
    top_print = top
    if with_ears:
        keep = box(-XI - 1.0, XI + 1.0, -big, big, -1, 1)
        top_print = top + ear(-d.OW / 2, Y_CUT + GAP_CUT, keep) + ear(d.OW / 2, Y_CUT + GAP_CUT, keep)
    return SimpleNamespace(d=d, base=base, top=top, top_print=top_print, phone=phone())


def to_tm(m):
    mesh = m.to_mesh()
    return trimesh.Trimesh(np.asarray(mesh.vert_properties)[:, :3], np.asarray(mesh.tri_verts), process=True)


def vol(m):
    return round(m.volume(), 3)


if __name__ == "__main__":
    out = sys.argv[1] if len(sys.argv) > 1 else "."

    m = build_moldura()
    d = m.d
    print(f"v{V_MOLDURA} moldura: externo {d.OW:.1f} x {d.OH:.1f} x {Z_TOP:.1f} mm")
    print("  bandeja vol", vol(m.tray), "genus", m.tray.genus(), " moldura vol", vol(m.frame), "genus", m.frame.genus())
    print("  bandeja∩moldura", vol(m.tray ^ m.frame), " bandeja∩celular", vol(m.tray ^ m.phone),
          " moldura∩celular", vol(m.frame ^ m.phone))
    print("  moldura descendo (dz 0.5/2/4):", [vol(m.tray ^ m.frame.translate((0, 0, dz))) for dz in (0.5, 2, 4)],
          " celular entrando por cima (dz 3/8/13):", [vol(m.tray ^ m.phone.translate((0, 0, dz))) for dz in (3, 8, 13)])
    tag = f"s26-mini-case_v{V_MOLDURA}"
    to_tm(m.tray).export(f"{out}/{tag}_bandeja.stl")
    to_tm(frame_for_print(m.frame)).export(f"{out}/{tag}_moldura.stl")

    t = build_trilho()
    d = t.d
    print(f"v{V_TRILHO} trilho: externo {d.OW:.1f} x {d.OH:.1f} x {Z_TOP:.1f} mm")
    print("  base vol", vol(t.base), "genus", t.base.genus(), " tampa vol", vol(t.top), "genus", t.top.genus())
    print("  base∩tampa", vol(t.base ^ t.top), " base∩celular", vol(t.base ^ t.phone), " tampa∩celular", vol(t.top ^ t.phone))
    print("  tampa deslizando (dy 1/10/30/60):", [vol(t.base ^ t.top.translate((0, dy, 0))) for dy in (1, 10, 30, 60)],
          " tampa sobre o celular:", [vol(t.phone ^ t.top.translate((0, dy, 0))) for dy in (1, 10, 30, 60)])
    print("  celular entrando na base (dy 5/20/40):", [vol(t.base ^ t.phone.translate((0, dy, 0))) for dy in (5, 20, 40)])
    tag = f"s26-mini-case_v{V_TRILHO}"
    to_tm(t.base).export(f"{out}/{tag}_A_base.stl")
    to_tm(t.top_print).export(f"{out}/{tag}_B_topo.stl")
