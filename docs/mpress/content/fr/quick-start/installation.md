---
title: "Installation"
description: "Installez Wails et préparez votre environnement pour créer des applications"
slug: "quick-start/installation"
sourcePath: "quick-start/installation.md"
---

## Installation rapide (5 minutes)

@note{type="tip" title="En bref — Développeurs expérimentés"}
```bash
# Install Go 1.25+, then:
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 setup   # Interactive setup wizard (experimental)
```

Vous pouvez également effectuer une vérification manuelle avec `wails3 doctor`. [Passer à la première application →](/quick-start/first-app/)

@end

## Installation pas à pas

@steps
### Installer Go (obligatoire)
Wails nécessite Go 1.25 ou une version ultérieure.

@tabs{sync-key="os"}
[Windows]
Téléchargez le programme d’installation pour Windows depuis **[go.dev/dl](https://go.dev/dl/)**, puis exécutez-le.

**Vérifiez l’installation :**

```powershell
go version  # Should show 1.25 or later
```

**Vérifiez le PATH :**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

S’il est vide, ajoutez `C:\Users\YourName\go\bin` à votre PATH.

[macOS]
**Option 1 : programme d’installation officiel**

Téléchargez le programme d’installation pour macOS (fichier .pkg) depuis **[go.dev/dl](https://go.dev/dl/)**, puis exécutez-le.

**Option 2 : Homebrew**

```bash
brew install go
```

**Vérifiez l’installation :**

```bash
go version  # Should show 1.25 or later
echo $PATH | grep go/bin  # Should show ~/go/bin
```

Si `~/go/bin` ne figure pas dans le PATH, ajoutez-le à `~/.zshrc` ou à `~/.bash_profile` :

```bash
export PATH=$PATH:~/go/bin
```

[Linux]
**Option 1 : archive tar officielle**

Téléchargez l’archive tar pour Linux depuis **[go.dev/dl](https://go.dev/dl/)**, puis :

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

**Option 2 : gestionnaire de paquets**

```bash
# Ubuntu/Debian
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**Ajoutez au PATH** (ajoutez à `~/.bashrc` ou à `~/.zshrc`) :

```bash
export PATH=$PATH:/usr/local/go/bin:~/go/bin
source ~/.bashrc  # Reload
```

**Vérifiez :**

```bash
go version
echo $PATH | grep go/bin
```

@end

### Installer les dépendances de la plateforme
@tabs{sync-key="os"}
[Windows]
**Runtime WebView2** (généralement préinstallé)

Windows 10/11 inclut WebView2 par défaut. S’il est absent :

- Téléchargez-le depuis [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/)
- Vous pouvez aussi exécuter `wails3 doctor` ultérieurement : il vous guidera

**C’est tout !** Aucune autre dépendance n’est nécessaire.

@note{type="tip" title="Conseil de performances pour Windows 11"}
Envisagez d’utiliser [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/) pour stocker vos projets. Les Dev Drives sont optimisés pour les charges de travail de développement et peuvent améliorer considérablement les temps de compilation et les vitesses d’accès au disque, jusqu’à 30 %.

@end

[macOS]
**Outils de ligne de commande Xcode** (obligatoires)

```bash
xcode-select --install
```

Cliquez sur « Installer » dans la boîte de dialogue qui s’affiche.

**Vérifiez :**

```bash
xcode-select -p  # Should show /Library/Developer/CommandLineTools
```

**C’est tout !** macOS inclut WebKit par défaut.

[Linux]
**Outils de compilation et WebKit**

@note{type="caution" title="Versions minimales des distributions"}
Par défaut, Wails v3 nécessite **WebKitGTK 6.0**. Sur les distributions qui fournissent uniquement WebKit2GTK 4.1 — Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x — vous devez compiler les applications Wails en activant explicitement l’option héritée `-tags gtk3`. Les versions plus anciennes qui fournissent uniquement WebKit2GTK 4.0 (Ubuntu 20.04, Debian 11, RHEL 8) ne sont pas prises en charge.

@end

@tabs{sync-key="distro"}
[Ubuntu/Debian]
La pile GTK4 par défaut nécessite Ubuntu 24.04 ou version ultérieure, ou Debian 13 ou version ultérieure.

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
Ajoutez ceci à votre `shell.nix` ou à votre `devShell` :

```nix
buildInputs = with pkgs; [ webkitgtk_6_0 gtk4 pkg-config gcc ];
```

[Autre]
Exécutez `wails3 doctor` après avoir installé Wails : cette commande indiquera précisément les paquets requis pour votre distribution.

@end

@note{type="info" title="Pile GTK3 héritée"}
Si votre distribution cible ne fournit pas encore WebKitGTK 6.0 (par exemple Ubuntu 22.04 LTS ou Debian 12), installez plutôt les bibliothèques de développement GTK3 et WebKit2GTK 4.1 (`libgtk-3-dev libwebkit2gtk-4.1-dev` sous Debian/Ubuntu ; leurs équivalents sous les autres distributions), puis compilez avec `wails3 build -tags gtk3`. La voie héritée est prise en charge jusqu’à la branche v3.0.x et sera supprimée dans v3.1. Pour en savoir plus, consultez [Empaquetage sous Linux — prise en charge de l’ancienne pile GTK3](/guides/build/linux/#legacy-gtk3-support).

@end

@end

### Installer l’interface en ligne de commande de Wails
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

Cette opération installe la commande `wails3` dans `~/go/bin` (ou dans `%USERPROFILE%\go\bin` sous Windows).

### Exécuter l’assistant de configuration (recommandé)
```bash
wails3 setup
```

L’assistant de configuration vérifiera vos dépendances, vous aidera à installer celles qui manquent et configurera les valeurs par défaut des projets.

@note{type="caution" title="Expérimental"}
L’assistant de configuration est récent et a principalement été testé sous Linux. Si vous rencontrez des problèmes, [signalez-les](https://github.com/wailsapp/wails/issues/4904) et utilisez plutôt `wails3 doctor`.

@end

### Vérifier l’installation
```bash
wails3 doctor
```

**Sortie attendue (ou similaire) :**

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

@note{type="info" title="Si la commande `wails3` est introuvable"}
Votre `~/go/bin` ne figure pas dans le PATH. Pour corriger ce problème, consultez l’étape 1 ci-dessus, puis redémarrez votre terminal.

@end

### Installer npm (facultatif, mais recommandé)
La plupart des modèles Wails utilisent npm pour les outils de développement frontend.

@tabs{sync-key="os"}
[Windows]
Téléchargez-le depuis [nodejs.org](https://nodejs.org/), puis exécutez le programme d’installation.

**Vérifiez :**

```powershell
npm --version
```

[macOS]
**Option 1 : programme d’installation officiel** Téléchargez-le depuis [nodejs.org](https://nodejs.org/)

**Option 2 : Homebrew**

```bash
brew install node
```

**Vérifiez :**

```bash
npm --version
```

[Linux]
**Option 1 : NodeSource**

```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs  # Ubuntu/Debian
```

**Option 2 : gestionnaire de paquets**

```bash
sudo dnf install nodejs  # Fedora
sudo pacman -S nodejs npm  # Arch
```

**Vérifiez :**

```bash
npm --version
```

@end

@note{type="tip" title="Autres gestionnaires de paquets"}
Vous préférez `pnpm`, `yarn` ou `bun` ? Aucun problème ! Il vous suffit de modifier le fichier `Taskfile.yml` de votre projet pour utiliser l’outil de votre choix.

@end

@end

## Dépannage

### Commande `wails3` introuvable

**Cause :** `~/go/bin` (ou `%USERPROFILE%\go\bin`) ne figure pas dans votre PATH.

**Solution :**

@tabs{sync-key="os"}
[Windows]
1. Ouvrez « Variables d’environnement » (recherchez ce terme dans le menu Démarrer).
2. Sous « Variables utilisateur », recherchez `Path`.
3. Cliquez sur « Modifier » → « Nouveau ».
4. Ajoutez : `C:\Users\YourName\go\bin` (remplacez `YourName`).
5. Cliquez sur « OK » dans toutes les boîtes de dialogue.
6. **Redémarrez votre terminal**

**Vérifiez :**

```powershell
$env:PATH -split ';' | Where-Object { $_ -like '*\go\bin' }
```

[macOS/Linux]
Ajoutez ceci à `~/.zshrc` (macOS) ou à `~/.bashrc` (Linux) :

```bash
export PATH=$PATH:~/go/bin
```

Rechargez :

```bash
source ~/.zshrc  # or ~/.bashrc
```

**Vérifiez :**

```bash
echo $PATH | grep go/bin
wails3 version
```

@end

---

#### `wails3 doctor` signale des dépendances manquantes

**Linux :** la sortie indique précisément les paquets à installer. Exemple :

```
❌ webkit2gtk not found
   Install with: sudo apt install libwebkit2gtk-4.1-dev
```

**Windows :** si WebView2 est manquant :

- Téléchargez-le depuis [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/).
- Sinon, il sera installé automatiquement lors de la première exécution de votre application.

**macOS :** si les outils Xcode sont manquants :

```bash
xcode-select --install
```

---

#### Version de Go trop ancienne

Wails v3 nécessite Go 1.25 ou version ultérieure. Si vous disposez d’une version plus ancienne :

@tabs{sync-key="os"}
[Windows/macOS]
Téléchargez la dernière version depuis [go.dev/dl](https://go.dev/dl/), puis réinstallez-la.

[Linux]
Téléchargez la dernière archive tar depuis [go.dev/dl](https://go.dev/dl/), puis procédez comme suit :

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
```

@end

## Version de développement (à la pointe)

Vous souhaitez utiliser le code le plus récent de la branche principale de développement ? Vous accéderez ainsi aux nouvelles fonctionnalités et aux correctifs avant leur publication, mais vous vous exposez à des bogues et à des modifications incompatibles. Cette version est recommandée uniquement aux contributeurs et aux personnes qui doivent tester les fonctionnalités à venir.

```bash
git clone https://github.com/wailsapp/wails.git
cd wails
git checkout v3
cd v3/cmd/wails3
go install
```

@note{type="caution" title="Version de développement"}
- Peut contenir des bogues ou des modifications incompatibles
- Les projets créés utiliseront la directive `replace` pour référencer la copie locale de Wails
- Recommandée uniquement pour contribuer ou tester de nouvelles fonctionnalités

@end

## Étapes suivantes

**Installation terminée !** Votre système est prêt pour le développement avec Wails.

@cards{cols="1"}
🚀 Créez votre première application
Créez une application fonctionnelle en 10 minutes.

[Tutoriel : votre première application →](/quick-start/first-app/)

@end

@cards{cols="1"}
📖 Découvrez les modèles
Découvrez ce qui est disponible immédiatement.

```bash
wails3 init -l  # List templates
```

@end

---

**Vous rencontrez des problèmes ?** Posez votre question sur [Discord](https://discord.gg/JDdSxwjhGf) ou [ouvrez un ticket](https://github.com/wailsapp/wails/issues).
