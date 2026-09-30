---
title: "Instalação"
description: "Instale o Wails e deixe-o pronto para criar aplicativos"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## Instalação rápida (5 minutos)

@note{type="tip" title="Resumo — desenvolvedores experientes"}
```bash
# Install Go 1.25+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 setup   # Interactive setup wizard (experimental)
```

Ou verifique manualmente com `wails3 doctor`. [Ir para o primeiro aplicativo →](/quick-start/first-app/)

@end

## Instalação passo a passo

@steps
### Instalar o Go (obrigatório)
O Wails requer o Go 1.25 ou posterior.

@tabs{sync-key="os"}
[Windows]
Baixe o instalador do Windows em **[go.dev/dl](https://go.dev/dl/)** e execute-o.

**Verifique a instalação:**

```powershell
go version  # Should show 1.25 or later
```

**Verifique o PATH:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

Se estiver vazio, adicione `C:\Users\YourName\go\bin` ao PATH.

[macOS]
**Opção 1: instalador oficial**

Baixe o instalador do macOS (arquivo .pkg) em **[go.dev/dl](https://go.dev/dl/)** e execute-o.

**Opção 2: Homebrew**

```bash
brew install go
```

**Verifique a instalação:**

```bash
go version  # Should show 1.25 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

Se `~/go/bin` não estiver no PATH, adicione-o a `~/.zshrc` ou `~/.bash_profile`:

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**Opção 1: arquivo tar oficial**

Baixe o arquivo tar para Linux em **[go.dev/dl](https://go.dev/dl/)** e, em seguida:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

**Opção 2: gerenciador de pacotes**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**Adicione ao PATH** (adicione a `~/.bashrc` ou `~/.zshrc`):

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**Verifique:**

```bash
go version
echo $PATH | grep go/bin
```

@end

### Instalar as dependências da plataforma
@tabs{sync-key="os"}
[Windows]
**Runtime do WebView2** (geralmente pré-instalado)

O Windows 10/11 inclui o WebView2 por padrão. Se ele não estiver instalado:

- Baixe-o da [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Ou execute `wails3 doctor` mais tarde — ele orientará você

**Pronto!** Nenhuma outra dependência é necessária.

@note{type="tip" title="Dica de desempenho para Windows 11"}
Considere usar o [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) para armazenar seus projetos. Os Dev Drives são otimizados para cargas de trabalho de desenvolvimento e podem melhorar significativamente os tempos de compilação e as velocidades de acesso ao disco em até 30%.

@end

[macOS]
**Ferramentas de linha de comando do Xcode** (obrigatórias)

```bash
xcode-select --install
```

Clique em "Instalar" na caixa de diálogo exibida.

**Verifique:**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**Pronto!** O macOS inclui o WebKit por padrão.

[Linux]
**Ferramentas de compilação e WebKit**

@note{type="caution" title="Versões mínimas das distribuições"}
Por padrão, o Wails v3 requer o **WebKitGTK 6.0**. Nas distribuições que fornecem somente o WebKit2GTK 4.1 — Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x — é necessário compilar os aplicativos Wails com a opção `-tags gtk3` para ativar explicitamente o suporte legado. Versões mais antigas que fornecem somente o WebKit2GTK 4.0 (Ubuntu 20.04, Debian 11, RHEL 8) não são compatíveis.

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
Requer Ubuntu 24.04+ ou Debian 13+ para a pilha GTK4 padrão.

```bash
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

[Fedora]
```bash
sudo dnf install gcc pkg-config gtk4-devel webkitgtk6.0-devel
```

[Arch]
```bash
sudo pacman -S base-devel gtk4 webkitgtk-6.0
```

[openSUSE]
```bash
sudo zypper install gcc pkg-config gtk4-devel webkitgtk-6_0-devel
```

[Gentoo]
```bash
sudo emerge --ask net-libs/webkit-gtk:6
```

[NixOS]
Adicione a `shell.nix` ou `devShell`:

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[Outras]
Execute `wails3 doctor` após instalar o Wails — ele mostrará os pacotes exatos necessários para sua distribuição.

@end

@note{type="info" title="Pilha GTK3 legada"}
Se a distribuição de destino ainda não fornecer o WebKitGTK 6.0 (por exemplo, Ubuntu 22.04 LTS, Debian 12), instale as bibliotecas de desenvolvimento do GTK3 + WebKit2GTK 4.1 (`libgtk-3-dev libwebkit2gtk-4.1-dev` no Debian/Ubuntu; equivalentes em outras distribuições) e compile com `wails3 build -tags gtk3`. O caminho legado será compatível durante toda a série v3.0.x e será removido na v3.1. Consulte [Empacotamento para Linux — compatibilidade com GTK3 legado](/guides/build/linux/#legacy-gtk3-support) para obter detalhes.

@end

@end

### Instalar a CLI do Wails
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Isso instala o comando `wails3` em `~/go/bin` (ou `%USERPROFILE%\go\bin` no Windows).

### Executar o assistente de configuração (recomendado)
```bash
wails3 setup
```

O assistente de configuração verificará suas dependências, ajudará a instalar as que estiverem ausentes e configurará os padrões do projeto.

@note{type="caution" title="Experimental"}
O assistente de configuração é novo e foi testado principalmente no Linux. Se você encontrar problemas, [relate-os](https://github.com/wailsapp/wails/issues/4904) e use `wails3 doctor` em vez dele.

@end

### Verificar a instalação
```bash
wails3 doctor
```

**Saída esperada (ou semelhante):**

```
Wails (v3.0.0-dev)  Wails Doctor

# System

┌──────────────────────────────────────────────────┐
| Name          | MacOS                            |
| Version       | 26.0                             |
| ID            | 25A354                           |
| Branding      | MacOS 26.0                       |
| Platform      | darwin                           |
| Architecture  | arm64                            |
| Apple Silicon | true                             |
| CPU           | Apple M2 Pro                     |
| CPU 1         | Apple M2 Pro                     |
| CPU 2         | Apple M2 Pro                     |
| GPU           | 16 cores, Metal Support: Metal 4 |
| Memory        | 16 GB                            |
└──────────────────────────────────────────────────┘

# Build Environment

┌─────────────┬─────────────────┐
| Wails CLI   | Your installed version |
| Go Version  | go1.25.0        |
└─────────────┴─────────────────┘

# Dependencies

┌─────────────────┬─────────────────────────────────────────────────┐
| npm             | 11.6.2                                          |
| *NSIS           | Not Installed. Install with `brew install...`.  |
| Xcode cli tools | 2412                                            |
└─────────────────┴─────────────────────────────────────────────────┘

# Checking for issues

SUCCESS No issues found

# Diagnosis

SUCCESS Your system is ready for Wails development!
```

@note{type="info" title="Se o comando `wails3` não for encontrado"}
Seu `~/go/bin` não está no PATH. Consulte a etapa 1 acima para corrigir isso e reinicie o terminal.

@end

### Instalar o npm (opcional, mas recomendado)
A maioria dos modelos do Wails usa o npm como ferramenta de frontend.

@tabs{sync-key="os"}
[Windows]
Baixe em [nodejs.org](https://nodejs.org/) e execute o instalador.

**Verifique:**

```powershell
npm --version
```

[macOS]
**Opção 1: instalador oficial** Baixe em [nodejs.org](https://nodejs.org/)

**Opção 2: Homebrew**

```bash
brew install node
```

**Verifique:**

```bash
npm --version
```

[Linux]
**Opção 1: NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**Opção 2: gerenciador de pacotes**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**Verifique:**

```bash
npm --version
```

@end

@note{type="tip" title="Gerenciadores de pacotes alternativos"}
Prefere `pnpm`, `yarn` ou `bun`? Sem problemas! Basta atualizar o `Taskfile.yml` do seu projeto para usar a ferramenta de sua preferência.

@end

@end

## Solução de problemas

### Comando `wails3` não encontrado

**Causa:** `~/go/bin` (ou `%USERPROFILE%\go\bin`) não está no seu PATH.

**Solução:**

@tabs{sync-key="os"}
[Windows]
1. Abra "Variáveis de Ambiente" (pesquise no menu Iniciar)
2. Em "Variáveis de usuário", localize `Path`
3. Clique em "Editar" → "Novo"
4. Adicione: `C:\Users\YourName\go\bin` (substitua `YourName`)
5. Clique em "OK" em todas as caixas de diálogo
6. **Reinicie o terminal**

**Verifique:**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
Adicione a `~/.zshrc` (macOS) ou `~/.bashrc` (Linux):

```bash
export PATH=$PATH:~/go/bin
```

Recarregue:

```bash
source ~/.zshrc  # or ~/.bashrc
```

**Verifique:**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor` informa que há dependências ausentes

**Linux:** a saída informa exatamente quais pacotes devem ser instalados. Exemplo:

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

**Windows:** se o WebView2 estiver ausente:

- Baixe da [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Ou ele será instalado automaticamente quando você executar seu primeiro aplicativo

**macOS:** se as ferramentas do Xcode estiverem ausentes:

```bash
xcode-select --install
```

---

#### Versão do Go muito antiga

O Wails v3 requer o Go 1.25 ou posterior. Se você tiver uma versão mais antiga:

@tabs{sync-key="os"}
[Windows/macOS]
Baixe a versão mais recente em [go.dev/dl](https://go.dev/dl/) e reinstale.

[Linux]
Baixe o tarball mais recente em [go.dev/dl](https://go.dev/dl/) e depois:

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

@end

## Versão de desenvolvimento (mais recente)

Quer usar o código mais recente da ramificação principal de desenvolvimento? Isso dá acesso a novos recursos e correções antes do lançamento, mas traz o risco de bugs e alterações incompatíveis. Recomendado apenas para colaboradores ou para quem precisa testar recursos futuros.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3
cd v3/cmd/wails3
go install
```

@note{type="caution" title="Versão de desenvolvimento"}
- Pode conter bugs ou alterações incompatíveis
- Os projetos criados usarão a diretiva `replace` para apontar para o Wails local
- Recomendado apenas para colaboradores ou para testar novos recursos

@end

## Próximas etapas

**Instalação concluída!** Seu sistema está pronto para desenvolver com o Wails.

@cards{cols="1"}
🚀 Crie seu primeiro aplicativo
Crie um aplicativo funcional em 10 minutos.

[Tutorial do primeiro aplicativo →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 Explore os modelos
Veja o que está disponível por padrão.

```bash
wails3 init -l  # List templates
```

@end

---

**Está com problemas?** Pergunte no [Discord](https://discord.gg/JDdSxwjhGf) ou [abra uma issue](https://github.com/wailsapp/wails/issues).
