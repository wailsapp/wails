---
title: "Création de paquets Linux"
description: "Créez des paquets de votre application Wails pour la distribuer sous Linux"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## Formats de paquets

Créez des paquets de votre application pour la distribuer sous Linux :

```bash
wails3 package GOOS=linux
```

Cette opération crée plusieurs formats dans le répertoire `bin/` :

- **AppImage** : format portable, exécutable sur toutes les distributions Linux
- **DEB** : pour Debian, Ubuntu et leurs dérivées
- **RPM** : pour Fedora, RHEL et leurs dérivées
- **Arch** : pour Arch Linux et ses dérivées

### Formats individuels

Créez des formats spécifiques :

```bash
wails3 task linux:create:appimage
wails3 task linux:create:deb
wails3 task linux:create:rpm
wails3 task linux:create:aur
```

## Personnalisation des paquets

### Entrée de bureau

Le fichier `.desktop` détermine la manière dont votre application apparaît dans les menus d’applications. Il est généré à partir des valeurs définies dans `build/linux/Taskfile.yml` :

```yaml
vars:
  APP_NAME: 'MyApp'
  EXEC: 'MyApp'
  ICON: 'MyApp'
  CATEGORIES: 'Development;'
```

### Métadonnées des paquets

Modifiez `build/linux/nfpm/nfpm.yaml` pour personnaliser les paquets DEB et RPM :

```yaml
name: myapp
version: 1.0.0
maintainer: Your Name <you@example.com>
description: My awesome Wails application
homepage: https://example.com
license: MIT
```

### AppImage

La configuration d’AppImage se trouve dans `build/linux/appimage/`. L’icône de l’application provient de `build/appicon.png`.

## Signature des paquets

Signez les paquets DEB et RPM avec une clé PGP :

```bash
# Using the wrapper (auto-detects platform)
wails3 sign GOOS=linux

# Or using tasks directly
wails3 task linux:sign:deb
wails3 task linux:sign:rpm
wails3 task linux:sign:packages  # Both
```

Configurez la signature dans `build/linux/Taskfile.yml` :

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  SIGN_ROLE: "builder"  # origin, maint, archive, or builder
```

Enregistrez le mot de passe de votre clé :

```bash
wails3 setup signing
```

Pour plus de détails, consultez [Signature des applications](/guides/build/signing/).

## Compilation pour ARM

```bash
wails3 build GOOS=linux GOARCH=arm64
wails3 package GOOS=linux GOARCH=arm64
```

@note{type="note"}
Les compilations ARM64 effectuées depuis des hôtes x86_64 utilisent Docker pour la compilation croisée avec CGO.

@end

## Prise en charge de l’ancienne version GTK3

Par défaut, Wails v3 repose sur **GTK4 avec WebKitGTK 6.0**. Une voie héritée utilisant GTK3 / WebKit2GTK 4.1 reste disponible pour les distributions qui ne fournissent pas encore WebKitGTK 6.0 (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x). Cette voie héritée doit être activée explicitement au moyen d’une balise de compilation et sa suppression est prévue dans la version v3.1.

@note{type="caution" title="Voie héritée"}
La voie GTK3 / WebKit2GTK 4.1 est prise en charge pendant toute la série v3.0.x. Planifiez la migration vers GTK4 en fonction de la disponibilité de GTK4 / WebKitGTK 6.0 dans votre distribution cible : `-tags gtk3` sera supprimé dans la version v3.1.

@end

### Dépendances

Installez les bibliothèques de développement GTK3 et WebKit2GTK 4.1 :

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk4.1-devel

# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

Les paquets pkg-config requis sont `gtk+-3.0` et `webkit2gtk-4.1`.

### Compilation avec GTK3

Utilisez l’option `-tags gtk3` :

```bash
wails3 build -tags gtk3
```

Ou directement avec Go :

```bash
go build -tags gtk3 -o myapp .
```

### Différences connues par rapport à GTK4

- **Boîtes de dialogue de fichiers** : GTK4 utilise par défaut `xdg-desktop-portal` pour les boîtes de dialogue de fichiers. Certaines options, telles que le répertoire par défaut et l’affichage des filtres personnalisés, se comportent donc différemment sous GTK3. Pour plus de détails, consultez [Référence des boîtes de dialogue — Comportement des boîtes de dialogue sous Linux](/reference/dialogs/#linux-dialog-behavior).
- **Style du menu** : GTK4 prend en charge une option `LinuxMenuStylePrimaryMenu` qui affiche un bouton de menu hamburger (☰) dans la barre d’en-tête, conformément aux recommandations GNOME HIG. Cette option n’a aucun effet sur les compilations `-tags gtk3`. Consultez [API Window — MenuStyle sous Linux](/reference/window/#linux).
- **Mise à l’échelle selon le DPI** : GTK4 utilise `gdk_monitor_get_scale` (GTK 4.14+) pour prendre en charge la mise à l’échelle fractionnaire.

### Vérification de votre compilation

Exécutez `wails3 doctor` pour vérifier votre configuration. Sans aucune option, cette commande vérifie la présence de GTK4 / WebKitGTK 6.0, qui sont utilisés par défaut. Les paquets hérités GTK3 / WebKit2GTK 4.1 sont indiqués comme facultatifs.

## Dépannage

### AppImage ne s’exécute pas

Rendez le fichier exécutable :

```bash
chmod +x MyApp-x86_64.AppImage
```

### Dépendances manquantes

Si l’application ne démarre pas, recherchez les dépendances WebKit manquantes :

```bash
# Debian/Ubuntu
sudo apt install libwebkit2gtk-4.1-0

# Fedora
sudo dnf install webkit2gtk4.1

# Arch
sudo pacman -S webkit2gtk-4.1
```

### Aucun compilateur C trouvé

Le système de compilation nécessite GCC ou Clang pour CGO :

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Vous pouvez également exécuter `wails3 task setup:docker` : le système de compilation utilisera alors Docker automatiquement.

### Fenêtre vide ou blanche avec un GPU NVIDIA

Sous Linux avec les pilotes propriétaires NVIDIA, les applications Wails peuvent afficher une fenêtre vide ou blanche au démarrage. Ce problème est dû à un bogue de WebKitGTK : le moteur de rendu DMA-BUF échoue avec `gbm_bo_map()` lorsqu’il utilise le pilote propriétaire NVIDIA. Ce bogue touche X11 et Wayland, les versions de pilote 377–580+, ainsi que les GPU de la série 10 et les modèles GT 710 plus anciens.

**Wails applique automatiquement `WEBKIT_DISABLE_DMABUF_RENDERER=1`** lorsqu’il détecte le module du noyau NVIDIA (`/sys/module/nvidia`). La plupart des utilisateurs n’ont donc aucune intervention à effectuer.

Si une fenêtre vide s’affiche encore, par exemple dans un conteneur où le chemin du module n’est pas visible, définissez manuellement la variable d’environnement avant de lancer votre application :

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Bogues connexes en amont : [WebKit nº 262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit nº 180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### Compatibilité du nettoyage des symboles d’AppImage

Sur les distributions Linux modernes (Arch Linux, Fedora 39+, Ubuntu 24.04+), les bibliothèques système sont compilées avec des sections ELF `.relr.dyn` afin d’accélérer les réadressages. L’outil `linuxdeploy` utilisé pour créer les AppImages intègre un ancien exécutable `strip` qui ne peut pas traiter ces sections modernes.

Wails détecte automatiquement cette situation en vérifiant les bibliothèques GTK du système avant de générer l’AppImage. Si elle est détectée, la suppression des symboles est désactivée (`NO_STRIP=1`) afin d’assurer la compatibilité.

**Ce que cela implique :**

- Les AppImages seront légèrement plus volumineuses (~20-40 %) sur les systèmes concernés
- Les fonctionnalités de l’application ne sont pas affectées
- Cette opération est gérée automatiquement : aucune action n’est requise

Si vous avez besoin d’AppImages moins volumineuses sur les systèmes modernes, vous pouvez installer un binaire `strip` plus récent et configurer `linuxdeploy` pour qu’il l’utilise à la place de la version fournie avec celui-ci.
