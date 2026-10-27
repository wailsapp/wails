---
title: "Empacotamento para Linux"
description: "Empacote seu aplicativo Wails para distribuição no Linux"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## Formatos de pacote

Empacote seu aplicativo para distribuição no Linux:

```bash
wails3 package GOOS=linux
```

Isso cria vários formatos no diretório `bin/`:

- **AppImage**: Portátil, é executado em qualquer distribuição Linux
- **DEB**: Para Debian, Ubuntu e derivados
- **RPM**: Para Fedora, RHEL e derivados
- **Arch**: Para Arch Linux e derivados

### Formatos individuais

Compile formatos específicos:

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## Personalização de pacotes

### Entrada da área de trabalho

O arquivo `.desktop` controla como seu aplicativo aparece nos menus de aplicativos. Ele é gerado com base nos valores de `build/linux/Taskfile.yml`:

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### Metadados do pacote

Edite `build/linux/nfpm/nfpm.yaml` para personalizar os pacotes DEB e RPM:

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

A configuração do AppImage fica em `build/linux/appimage/`. O ícone do aplicativo vem de `build/appicon.png`.

O AppImage inclui o WebKit e seus processos auxiliares, além de aplicar um patch à biblioteca incluída para iniciá-los de dentro do AppImage. Seu `AppRun` inicia o aplicativo a partir de `usr/lib` dentro do AppImage, onde o caminho corrigido é resolvido tanto fora quanto dentro do sandbox bubblewrap do WebKit, de modo que o sandbox permaneça ativado. Forneça seu próprio `AppRun` somente se ele mantiver esse diretório de trabalho.

## Assinatura de pacotes

Assine pacotes DEB e RPM com uma chave PGP:

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

Configure a assinatura em `build/linux/Taskfile.yml`:

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

Armazene a senha da sua chave:

```bash
wails3 signing
```

Consulte [Assinatura de aplicativos](/guides/build/signing/) para obter detalhes.

## Compilação para ARM

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
Compilações ARM64 em hosts x86_64 usam Docker para compilação cruzada com CGO.

@end

## Suporte legado ao GTK3

Por padrão, o Wails v3 é compilado com **GTK4 e WebKitGTK 6.0**. Um caminho legado com GTK3/WebKit2GTK 4.1 ainda está disponível para distribuições que ainda não fornecem WebKitGTK 6.0 (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x). O caminho legado deve ser habilitado explicitamente por uma tag de compilação e está programado para ser removido no v3.1.

@note{type="caution" title="Caminho legado"}
O caminho GTK3/WebKit2GTK 4.1 terá suporte durante toda a série v3.0.x. Planeje a migração para GTK4 de acordo com a disponibilidade de GTK4/WebKitGTK 6.0 na sua distribuição de destino — `-tags gtk3` será removido no v3.1.

@end

### Dependências

Instale as bibliotecas de desenvolvimento do GTK3 e WebKit2GTK 4.1:

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Os pacotes pkg-config necessários são `gtk+-3.0` e `webkit2gtk-4.1`.

### Compilação com GTK3

Use a opção `-tags gtk3`:

```bash
wails3 build -tags gtk3
```

Ou diretamente com o Go:

```bash
go build -tags gtk3 -o myapp .
```

### Diferenças conhecidas em relação ao GTK4

- **Caixas de diálogo de arquivos**: o GTK4 usa `xdg-desktop-portal` para caixas de diálogo de arquivos (o padrão), o que significa que algumas opções dessas caixas (como o diretório padrão e a exibição de filtros personalizados) funcionam de maneira diferente do GTK3. Consulte [Referência de caixas de diálogo — comportamento das caixas de diálogo no Linux](/reference/dialogs/#linux-dialog-behavior) para obter detalhes.
- **Estilo do menu**: o GTK4 oferece uma opção `LinuxMenuStylePrimaryMenu` que exibe um botão de menu hambúrguer (☰) na barra de cabeçalho, de acordo com as HIG do GNOME. Essa opção não tem efeito em compilações `-tags gtk3`. Consulte [API de janelas — MenuStyle no Linux](/reference/window/#linux).
- **Escalonamento de DPI**: o GTK4 usa `gdk_monitor_get_scale` (GTK 4.14+) para oferecer suporte a escalonamento fracionário.

### Verificação da compilação

Execute `wails3 doctor` para verificar sua configuração. Sem opções, ele verifica o GTK4/WebKitGTK 6.0 (o padrão). Os pacotes legados GTK3/WebKit2GTK 4.1 são listados como opcionais.

## Solução de problemas

### O AppImage não é executado

Torne-o executável:

```bash
chmod +x MyApp-x86_64.AppImage
```

### Dependências ausentes

Use as bibliotecas correspondentes à variante de compilação do binário. Para uma compilação padrão com GTK4/WebKitGTK 6.0:

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

Para um binário compilado explicitamente com `-tags gtk3`, instale as bibliotecas legadas GTK3/WebKit2GTK 4.1 em vez disso:

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

A disponibilidade e os nomes dos pacotes dependem da versão da distribuição. Use uma versão compatível que forneça as bibliotecas necessárias ou o caminho documentado de compilação cruzada/contêiner; instalar as bibliotecas legadas não satisfaz os requisitos de uma compilação padrão. A portabilidade de um AppImage depende do que o empacotamento realmente inclui e de quais bibliotecas do host ele ainda requer.

### Nenhum compilador C encontrado

O sistema de compilação precisa do GCC ou Clang para o CGO:

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Como alternativa, execute `wails3 task setup:docker`, e o sistema de compilação usará o Docker automaticamente.

### Janela vazia ou branca em GPU NVIDIA

No Linux com drivers proprietários da NVIDIA, os aplicativos Wails podem exibir uma janela vazia ou branca ao iniciar. Isso é causado por um bug do WebKitGTK no qual o renderizador DMA-BUF falha com `gbm_bo_map()` ao usar o driver proprietário da NVIDIA (afeta X11 e Wayland, versões do driver 377–580+, GPUs da série 10 e modelos GT 710 mais antigos).

**O Wails aplica `WEBKIT_DISABLE_DMABUF_RENDERER=1` automaticamente** ao detectar o módulo do kernel da NVIDIA (`/sys/module/nvidia`), portanto a maioria dos usuários não precisa fazer nada.

Se você ainda vir uma janela vazia (por exemplo, em um contêiner no qual o caminho do módulo não está visível), defina manualmente a variável de ambiente antes de iniciar seu aplicativo:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Bugs relacionados nos projetos upstream: [WebKit nº 262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit nº 180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### Compatibilidade da remoção de símbolos do AppImage

Nas distribuições Linux modernas (Arch Linux, Fedora 39+, Ubuntu 24.04+), as bibliotecas do sistema são compiladas com seções ELF `.relr.dyn` para tornar as realocações mais eficientes. A ferramenta `linuxdeploy` usada para criar AppImages inclui uma versão mais antiga do binário `strip`, que não consegue processar essas seções modernas.

O Wails detecta automaticamente essa situação verificando as bibliotecas GTK do sistema antes de compilar o AppImage. Quando ela é detectada, a remoção de símbolos é desativada (`NO_STRIP=1`) para garantir a compatibilidade.

**O que isso significa:**

- Os AppImages serão ligeiramente maiores (~20-40%) nos sistemas afetados
- A funcionalidade do aplicativo não é afetada
- Isso é tratado automaticamente — nenhuma ação é necessária

Se você precisar de AppImages menores em sistemas modernos, poderá instalar um binário `strip` mais recente e configurar o `linuxdeploy` para usá-lo em vez da versão incluída.
