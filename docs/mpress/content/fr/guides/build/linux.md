---
title: "Création de paquets Linux"
description: "Préparez votre application Wails pour sa distribution sous Linux"
slug: "guides/build/linux"
sourcePath: "guides/build/linux.md"
---

## Formats de paquets

Préparez votre application pour sa distribution sous Linux :

```bash
wails3 package GOOS=linux
```

Cette opération crée plusieurs formats dans le répertoire `bin/` :

- **AppImage** : portable, s’exécute sur toutes les distributions Linux
- **DEB** : pour Debian, Ubuntu et leurs dérivés
- **RPM** : pour Fedora, RHEL et leurs dérivés
- **Arch** : pour Arch Linux et ses dérivés

### Formats individuels

Générez des formats spécifiques :

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

L’AppImage inclut WebKit et ses processus auxiliaires, puis modifie la bibliothèque incluse afin de les démarrer depuis l’intérieur de l’AppImage. Son fichier `AppRun` démarre l’application depuis `usr/lib` dans l’AppImage, où le chemin modifié est résolu à la fois à l’extérieur et à l’intérieur du bac à sable bubblewrap de WebKit, de sorte que le bac à sable reste activé. Ne fournissez votre propre fichier `AppRun` que s’il conserve ce répertoire de travail.

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

Stockez le mot de passe de votre clé :

```bash
wails3 signing
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

Par défaut, Wails v3 repose sur **GTK4 avec WebKitGTK 6.0**. Une voie héritée GTK3 / WebKit2GTK 4.1 reste disponible pour les distributions qui ne fournissent pas encore WebKitGTK 6.0 (Ubuntu 22.04 LTS, Debian 12, Fedora ≤ 39, RHEL 9.x). Cette voie héritée doit être activée explicitement au moyen d’une balise de compilation et sa suppression est prévue dans v3.1.

@note{type="caution" title="Voie héritée"}
La voie GTK3 / WebKit2GTK 4.1 est prise en charge pendant toute la série v3.0.x. Planifiez la migration vers GTK4 en fonction de la disponibilité de GTK4 / WebKitGTK 6.0 dans votre distribution cible : `-tags gtk3` sera supprimé dans v3.1.

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

- **Boîtes de dialogue de fichiers** : GTK4 utilise `xdg-desktop-portal` par défaut pour les boîtes de dialogue de fichiers. Certaines options de ces boîtes, comme le répertoire par défaut et l’affichage des filtres personnalisés, se comportent donc différemment par rapport à GTK3. Pour plus de détails, consultez [Référence des boîtes de dialogue — comportement des boîtes de dialogue sous Linux](/reference/dialogs/#linux-dialog-behavior).
- **Style de menu** : GTK4 prend en charge une option `LinuxMenuStylePrimaryMenu` qui affiche un bouton de menu hamburger (☰) dans la barre d’en-tête, conformément aux recommandations GNOME HIG. Cette option n’a aucun effet sur les compilations `-tags gtk3`. Consultez [API Window — Linux MenuStyle](/reference/window/#linux).
- **Mise à l’échelle selon la résolution** : GTK4 utilise `gdk_monitor_get_scale` (GTK 4.14+) pour prendre en charge la mise à l’échelle fractionnaire.

### Vérification de votre compilation

Exécutez `wails3 doctor` pour vérifier votre configuration. Sans option, cette commande recherche GTK4 / WebKitGTK 6.0, utilisés par défaut. Les paquets hérités GTK3 / WebKit2GTK 4.1 sont indiqués comme facultatifs.

## Dépannage

### L’AppImage ne s’exécute pas

Rendez-la exécutable :

```bash
chmod +x MyApp-x86_64.AppImage
```

### Dépendances manquantes

Faites correspondre les bibliothèques à la variante de compilation du fichier binaire. Pour une compilation GTK4/WebKitGTK 6.0 par défaut :

```bash
# Debian/Ubuntu versions that provide WebKitGTK 6.0
sudo apt install libgtk-4-1 libwebkitgtk-6.0-4
# Fedora
sudo dnf install gtk4 webkitgtk6.0
# Arch
sudo pacman -S gtk4 webkitgtk-6.0
```

Pour un fichier binaire explicitement compilé avec `-tags gtk3`, installez plutôt les bibliothèques héritées GTK3/WebKit2GTK 4.1 :

```bash
# Debian/Ubuntu
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
# Fedora
sudo dnf install gtk3 webkit2gtk4.1
# Arch
sudo pacman -S gtk3 webkit2gtk-4.1
```

La disponibilité et le nom des paquets dépendent de la version de la distribution. Utilisez une version prise en charge qui fournit les bibliothèques requises, ou la méthode documentée de compilation croisée ou par conteneur ; l’installation des bibliothèques héritées ne satisfait pas les exigences d’une compilation par défaut. La portabilité d’une AppImage dépend des éléments réellement inclus dans son paquet et des bibliothèques de l’hôte dont elle a encore besoin.

### Aucun compilateur C trouvé

Le système de compilation a besoin de GCC ou de Clang pour CGO :

```bash
# Debian/Ubuntu
sudo apt install build-essential

# Fedora
sudo dnf install gcc

# Arch
sudo pacman -S base-devel
```

Vous pouvez également exécuter `wails3 task setup:docker` ; le système de compilation utilisera alors Docker automatiquement.

### Fenêtre vide ou blanche avec un processeur graphique NVIDIA

Sous Linux avec les pilotes propriétaires NVIDIA, les applications Wails peuvent afficher une fenêtre vide ou blanche au démarrage. Cela est dû à un bogue de WebKitGTK : le moteur de rendu DMA-BUF échoue avec `gbm_bo_map()` face au pilote propriétaire NVIDIA. Le problème affecte X11 et Wayland, les versions de pilote 377–580+, ainsi que les processeurs graphiques de la série 10 et les anciens modèles GT 710.

**Wails applique automatiquement `WEBKIT_DISABLE_DMABUF_RENDERER=1`** lorsqu’il détecte le module du noyau NVIDIA (`/sys/module/nvidia`). La plupart des utilisateurs n’ont donc rien à faire.

Si une fenêtre vide s’affiche encore, par exemple dans un conteneur où le chemin du module n’est pas visible, définissez manuellement la variable d’environnement avant de lancer votre application :

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./myapp
```

Bogues connexes en amont : [WebKit nº 262607](https://bugs.webkit.org/show_bug.cgi?id=262607), [WebKit nº 180739](https://bugs.webkit.org/show_bug.cgi?id=180739).

### Compatibilité du dépouillement des AppImages

Sur les distributions Linux modernes (Arch Linux, Fedora 39+, Ubuntu 24.04+), les bibliothèques système sont compilées avec des sections ELF `.relr.dyn` afin d’optimiser les réadressages. L’outil `linuxdeploy` utilisé pour créer les AppImages inclut une ancienne version du fichier binaire `strip`, qui ne peut pas traiter ces sections modernes.

Wails détecte automatiquement cette situation en vérifiant les bibliothèques GTK du système avant de compiler l’AppImage. Lorsqu’elle est détectée, le dépouillement est désactivé (`NO_STRIP=1`) afin de garantir la compatibilité.

**Conséquences :**

- Les AppImages seront légèrement plus volumineuses (~20-40 %) sur les systèmes concernés
- Le fonctionnement de l’application n’est pas affecté
- Tout est géré automatiquement : aucune intervention n’est requise

Si vous avez besoin d’AppImages plus petites sur les systèmes modernes, vous pouvez installer une version plus récente du fichier binaire `strip` et configurer `linuxdeploy` afin qu’il l’utilise à la place de la version qu’il inclut.
