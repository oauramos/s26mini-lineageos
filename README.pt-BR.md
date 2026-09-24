# "S26 ULTRA Mini" falso: remova o backdoor de fábrica

🇺🇸 [Read in English](README.md)

Um "Samsung S26 ULTRA Mini" barato é vendido no Brasil inteiro. Ele não é Samsung: é uma placa genérica MediaTek MT6739, e o sistema já vem, **com a caixa ainda lacrada**, com:

- um **carregador de código escondido no núcleo do Android**, capaz de injetar código em qualquer app do celular (banco, WhatsApp, navegador)
- ferramentas com privilégio de sistema que **alteram o IMEI** e mentem as especificações
- **patch de segurança falso** (diz 2026, é de 2020)
- um serviço chinês de "atualização" rodando o tempo todo como sistema

Análise completa, com indicadores: [docs/backdoor-analysis.md](docs/backdoor-analysis.md) (em inglês).

Este repositório mostra como **substituir a partição de sistema inteira** por um Android limpo e open source (uma GSI Treble, como LineageOS ou TrebleDroid). Tudo o que está listado acima fica em `/system`, então sai junto.

> ⚠️ **Isso apaga tudo no celular.** O processo desbloqueia o bootloader e instala software não oficial. Só foi testado na placa `d39g_4m_bml_s26ultra_mini_pt`. O risco é seu.
>
> 🚫 Este projeto **não** fornece nem apoia a troca de IMEI, que é crime no Brasil. O sistema limpo remove as ferramentas de IMEI que vieram no aparelho.

## Meu celular é esse modelo?

Ative a depuração USB e rode:

```bash
adb shell getprop ro.product.vendor.device
```

Tem que aparecer `d39g_4m_bml_s26ultra_mini_pt`. Outros aparelhos vendidos com o mesmo nome podem ter placa diferente. O instalador avisa e pede confirmação antes de continuar.

## O que você precisa

- Computador com `adb`, `fastboot` e `python3` (Linux ou macOS)
- Cabo USB de dados, ligado **direto** no computador (sem hub)
- Uma imagem GSI **arm64, A/B, vanilla**:
  - **Recomendada:** LineageOS 21 (AndyYan), arquivo `...-arm64_bvN.img.gz`: <https://sourceforge.net/projects/andyyan-gsi/files/lineage-21-pre-qpr2-td/>
  - Alternativa: TrebleDroid `system-td-arm64-ab-vanilla.img.xz`: <https://github.com/TrebleDroid/treble_experimentations/releases>

## Instalação

1. **No celular:** Configurações → Sobre o telefone → toque 7× em **Número da versão**. Depois, em Opções do desenvolvedor, ative **Desbloqueio de OEM** e **Depuração USB**.
2. Plugue no computador e aceite o aviso "Permitir depuração USB?".
3. Descompacte a GSI (`gunzip arquivo.img.gz` ou `xz -d arquivo.img.xz`).
4. Rode:
   ```bash
   scripts/flash-gsi.sh caminho/para/system.img
   ```
5. Quando o celular mostrar o aviso de desbloqueio, aperte **Volume +** pra confirmar.

O primeiro boot demora alguns minutos. O aviso de "orange state" em todo boot é normal com o bootloader desbloqueado.

### Depois de instalar

- Depois do primeiro boot, ative a depuração USB de novo e rode `scripts/post-install.sh`. Ele aplica as correções do aparelho (o Bluetooth não liga no Android 14+ sem elas). Opcional: `scripts/optimize.sh` deixa o sistema mais leve e rápido.
- **Faça backup das partições de calibração** (IMEI, MAC, calibração de rádio) e do preloader: [docs/recovery.md](docs/recovery.md).
- **Desligue a Depuração USB** no final.
- O vendor, o kernel e o modem continuam sendo os de fábrica. Trate o aparelho como **semi-confiável**: nada de conta Google principal nem app de banco.

## Contribuindo

Testou outra GSI ou outro aparelho vendido como "S26 ULTRA Mini"? Abra uma issue. **Nunca publique** IMEI, número de série ou dumps de `nvram`/`nvdata`.
