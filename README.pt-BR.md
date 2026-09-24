<h1 align="center">s26mini-lineageos</h1>

<p align="center">
  <b>Remova o backdoor de fábrica do "S26 ULTRA Mini" falso e rode LineageOS 21 nele.</b>
</p>

<p align="center">
  <a href="https://github.com/oauramos/s26mini-lineageos/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/oauramos/s26mini-lineageos"></a>
  <a href="LICENSE"><img alt="Licença: MIT" src="https://img.shields.io/github/license/oauramos/s26mini-lineageos"></a>
  <a href="https://github.com/oauramos/s26mini-lineageos/wiki"><img alt="Wiki" src="https://img.shields.io/badge/docs-wiki-blue"></a>
</p>

<p align="center">
  🇺🇸 <a href="README.md"><b>Read in English</b></a>
</p>

---

Um **"Samsung S26 ULTRA Mini"** barato é vendido no Brasil inteiro. Ele não é Samsung: é uma placa genérica **MediaTek MT6739**, e o sistema já vem, **com a caixa ainda lacrada**, com:

- um **carregador de código escondido no núcleo do Android**, capaz de injetar código em qualquer app do celular (banco, WhatsApp, navegador)
- ferramentas com privilégio de sistema que **alteram o IMEI** e mentem as especificações
- **patch de segurança falso** (diz 2026, é de 2020)
- um serviço chinês de "atualização" rodando o tempo todo como sistema

Este projeto **troca o sistema inteiro** pelo **LineageOS 21** (Android 14, patches de segurança de setembro de 2026). Tudo o que está listado acima fica em `/system`, então sai junto.

| | De fábrica | Com este projeto |
|---|---|---|
| Sistema | Android 10 Go disfarçado | **LineageOS 21** (Android 14) |
| Patch de segurança | diz 2026-01, é de 2020-08 | **2026-09-01** |
| Backdoor, ferramentas de IMEI, FOTA | presentes | **removidos** |

Análise completa, com indicadores (em inglês): [Stock Firmware Analysis](https://github.com/oauramos/s26mini-lineageos/wiki/Stock-Firmware-Analysis)

> ⚠️ **Isso apaga tudo no celular** e desbloqueia o bootloader. Só foi testado na placa `d39g_4m_bml_s26ultra_mini_pt`. O risco é seu.
>
> 🚫 Este projeto **não** fornece nem apoia a troca de IMEI, que é **crime no Brasil**. O sistema limpo remove as ferramentas de IMEI que vieram no aparelho.

## Início rápido

```bash
git clone https://github.com/oauramos/s26mini-lineageos && cd s26mini-lineageos
scripts/flash-gsi.sh lineage-21.0-<data>-UNOFFICIAL-arm64_bvN.img   # apaga e instala
scripts/post-install.sh                                              # correções do aparelho
scripts/profile-termux.sh    # versão Termux  (ou: scripts/install-apps.sh pra versão Completa)
```

## Meu celular é esse modelo?

Ative a depuração USB e rode:

```bash
adb shell getprop ro.product.vendor.device
```

Tem que aparecer `d39g_4m_bml_s26ultra_mini_pt`. Outros aparelhos vendidos com o mesmo nome podem ter placa diferente. O instalador avisa e pede confirmação antes de continuar.

## O que você precisa

- Computador com `adb`, `fastboot`, `python3` e `curl` (Linux ou macOS)
- Cabo **USB-A → USB-C**. **USB-C ↔ USB-C não funciona** nesse celular, nem pra dados nem pra carregar (faltam os resistores CC do USB-C na placa). Em computador só com USB-C, use **um** adaptador C→A, sem hub.
- Uma imagem GSI **arm64, A/B, vanilla**:
  - **Recomendada:** LineageOS 21 (AndyYan), arquivo `...-arm64_bvN.img.gz`: <https://sourceforge.net/projects/andyyan-gsi/files/lineage-21-pre-qpr2-td/> (download lento? acrescente `?use_mirror=cfhcable` no link)
  - Alternativa: TrebleDroid `system-td-arm64-ab-vanilla.img.xz`: <https://github.com/TrebleDroid/treble_experimentations/releases>

## Instalação

1. **No celular:** Configurações → Sobre o telefone → toque 7× em **Número da versão**. Depois, em Opções do desenvolvedor, ative **Desbloqueio de OEM** e **Depuração USB**.
2. Plugue no computador e aceite o aviso "Permitir depuração USB?".
3. Descompacte a GSI (`gunzip arquivo.img.gz` ou `xz -d arquivo.img.xz`).
4. Rode `scripts/flash-gsi.sh caminho/para/system.img`.
5. Quando o celular mostrar o aviso de desbloqueio, aperte **Volume +** pra confirmar.

O primeiro boot demora alguns minutos. O aviso de "orange state" em todo boot é normal com o bootloader desbloqueado.

## Escolha a versão

Depois do primeiro boot, ative a depuração USB de novo, rode `scripts/post-install.sh` (corrige o Bluetooth) e depois **uma** das opções:

| | **Completa** | **Versão Termux** |
|---|---|---|
| Comando | `scripts/install-apps.sh` | `scripts/profile-termux.sh` |
| Pra quê | celular comum | terminal SSH dedicado |
| Configurações, Launcher, Câmera, Relógio, Calculadora, Arquivos, Teclado | ✅ | ✅ |
| Navegador do LineageOS (Jelly) | ✅ | substituído pelo Firefox |
| Galeria, Música, Gravador, Agenda | ✅ | desativados |
| Telefone, Contatos, SMS | ✅ | desativados (edite o `profile-termux.sh` se for fazer ligações) |
| Protetores de tela, impressão, papéis de parede animados, Seedvault, AudioFX | ✅ | desativados |
| Animações 0,5x, 160 dpi, apps pré-compilados | – | ✅ |
| **Pacote de apps** (abaixo) | ✅ | ✅ |
| **[Firefox](https://f-droid.org/packages/org.mozilla.fennec_fdroid/)** (Fennec F-Droid), navegador padrão | – | ✅ |
| **[Primitive FTPd](https://f-droid.org/packages/org.primftpd/)**: servidor FTP/SFTP no celular | – | ✅ |
| RAM livre (testado) | ~1,3 GB | **~1,9 GB** |

**Pacote de apps**, do repositório oficial do F-Droid, com cada APK conferido por SHA-256:

- [F-Droid](https://f-droid.org): loja de apps open source, mantém tudo aqui atualizado
- [Termux](https://termux.dev): terminal com `ssh`, `sftp`, `scp`, `git`...
- [LocalSend](https://localsend.org): envia arquivos entre computadores e celulares na mesma rede, sem internet
- [Material Files](https://github.com/zhanghai/MaterialFiles): gerenciador de arquivos com cliente **SFTP**, SMB, FTP e WebDAV (o celular acessando servidores)

O **Primitive FTPd** faz o caminho contrário: é um servidor SFTP/FTP pro computador acessar os arquivos do celular. O **Fennec F-Droid** é o Firefox compilado do código-fonte da Mozilla pelo próprio F-Droid, sem telemetria. Builds de terceiros não podem usar o nome "Firefox".

Os apps desativados não são apagados: `adb shell pm enable <pacote>` reativa qualquer um. Não há serviços do Google, então apps do Google (Maps, Play Store...) não funcionam.

## O que funciona

| Funciona | Observação |
|---|---|
| Tela, toque, brilho, Wi-Fi 2,4 e 5 GHz, ligações, dados móveis, SMS, câmeras, áudio, GPS, sensores, bateria | o GPS não tem A-GPS do Google, então a primeira localização demora mais |
| Bluetooth | precisa do `post-install.sh` |

**Problemas conhecidos:**

- **Fone Bluetooth pareia mas não conecta** quando é o celular que inicia a conexão. Depois de parear, desligue e ligue o fone uma vez. Daí em diante ele reconecta sozinho.
- **USB-C ↔ USB-C** não funciona. Use USB-A → C.

## Depois de instalar

- **Faça backup das partições de calibração** (IMEI, MAC, calibração de rádio) e do preloader: [Recovery and Backups](https://github.com/oauramos/s26mini-lineageos/wiki/Recovery-and-Backups).
- Abra o **F-Droid** uma vez pra ele atualizar a lista de apps.
- **Desligue a Depuração USB** no final.

## O que continua não confiável

O vendor, o kernel, o bootloader e o modem continuam sendo os de fábrica, e não há como substituí-los (o código-fonte do kernel nunca foi publicado). Nenhum indicador de backdoor foi encontrado no `vendor`, mas trate o aparelho como **semi-confiável**: nada de conta Google principal nem app de banco.

## Documentação

A **[Wiki](https://github.com/oauramos/s26mini-lineageos/wiki)** (em inglês) tem tudo em detalhe: hardware, análise do firmware, decisões de projeto, guia de instalação, lista completa de apps, recuperação, modelo de segurança e FAQ.

## Contribuindo

Testou outra GSI ou outro aparelho vendido como "S26 ULTRA Mini"? Abra uma issue. **Nunca publique** IMEI, número de série ou dumps de `nvram`/`nvdata`.

## Licença

Scripts e documentação: MIT. Veja [LICENSE](LICENSE).
