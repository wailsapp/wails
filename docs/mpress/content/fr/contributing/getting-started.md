---
title: "Bien démarrer"
description: "Comment commencer à contribuer à Wails v3"
slug: "contributing/getting-started"
sourcePath: "contributing/getting-started.md"
---

## Bienvenue parmi les contributeurs !

Merci de l’intérêt que vous portez à Wails ! Ce guide vous aidera à apporter votre première contribution.

## Prérequis

Avant de commencer, assurez-vous de disposer des éléments suivants :

- **Go 1.26+** installé ([télécharger](https://go.dev/dl/))
- **Node.js 20+** et **npm** ([télécharger](https://nodejs.org/))
- **Git** configuré avec votre compte GitHub
- Connaissances de base de Go et de JavaScript/TypeScript

### Prérequis propres à chaque plateforme

**macOS :**

- Outils en ligne de commande Xcode : `xcode-select --install`

**Windows :**

- MSYS2 ou un environnement similaire de type Unix est recommandé
- Runtime WebView2 (généralement préinstallé sous Windows 11)

**Linux :**

- `gcc`, `pkg-config`, `libgtk-4-dev`, `libwebkitgtk-6.0-dev` (pile GTK4 par défaut)
- Installez-les avec : `sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev` (Debian/Ubuntu)
- Pour l’ancien processus de compilation `-tags gtk3`, installez également `libgtk-3-dev` et `libwebkit2gtk-4.1-dev`

## Vue d’ensemble du processus de contribution

Le processus de contribution habituel suit les étapes ci-dessous :

1. **Fork et clonage** — Créez votre propre copie du dépôt Wails
2. **Configuration** — Compilez la CLI Wails et vérifiez votre environnement
3. **Branche** — Créez une branche de fonctionnalité pour vos modifications
4. **Développement** — Apportez vos modifications conformément à nos normes de codage
5. **Tests** — Exécutez les tests pour vérifier que tout fonctionne
6. **Commit** — Créez un commit avec un message clair respectant les conventions
7. **Soumission** — Ouvrez une pull request pour révision
8. **Itération** — Répondez aux retours et apportez les ajustements nécessaires
9. **Fusion** — Une fois approuvées, vos modifications sont intégrées à Wails !

## Guide détaillé

Choisissez votre type de contribution :

@tabs
[Correction de bug]
@steps
### Rechercher ou signaler le bug
- Vérifiez si le bug a déjà été signalé dans les [issues GitHub](https://github.com/wailsapp/wails/issues)
- Dans le cas contraire, créez une nouvelle issue en indiquant les étapes permettant de reproduire le bug
- Attendez une confirmation avant de commencer le travail

### Créer un fork et le cloner
Créez un fork du dépôt à l’adresse [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Clonez votre fork :

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Compiler et vérifier
Compilez Wails et vérifiez que vous pouvez reproduire le bug :

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Reproduce the bug to understand it
```

### Créer une branche de correction de bug
Créez une branche pour votre correction :

```bash
git checkout -b fix/issue-123-window-crash
```

### Corriger le bug
- Apportez uniquement les modifications nécessaires pour corriger le bug
- Ne remaniez pas de code sans rapport avec le bug
- Ajoutez ou mettez à jour des tests pour éviter toute régression

```bash
# Make your changes
# Add tests in *_test.go files
```

### Tester votre correction
Exécutez les tests pour vérifier que la correction fonctionne :

```bash
go test ./...

# Test the specific package
go test ./pkg/application -v

# Run with race detector
go test ./... -race
```

### Créer un commit pour votre correction
Créez un commit avec un message clair :

```bash
git commit -m "fix: prevent window crash when closing during initialization

Fixes #123"
```

### Soumettre une pull request
Poussez votre branche et créez une pull request :

```bash
git push origin fix/issue-123-window-crash
```

Dans la description de votre pull request :

- Expliquez le bug et sa cause racine
- Décrivez votre correction
- Référencez l’issue : « Fixes #123 »
- Indiquez le comportement avant et après la correction

### Répondre aux retours
Traitez les commentaires de révision et mettez à jour votre pull request si nécessaire.

@end

[WEP (amélioration)]
@steps
### Rédiger une WEP
- Consultez le [processus WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md)
- Copiez le modèle de WEP dans `v3/wep/proposals/<name>/proposal.md`
- Ouvrez une pull request en brouillon intitulée `[WEP] <title>` et contenant uniquement la WEP
- Attendez la décision d’un responsable de maintenance avant de commencer l’implémentation

### Créer un fork et le cloner
Créez un fork du dépôt à l’adresse [github.com/wailsapp/wails/fork](https://github.com/wailsapp/wails/fork)

Clonez votre fork :

```bash
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
```

### Configurer l’environnement de développement
Compilez Wails et vérifiez votre environnement :

```bash
cd v3
go build -o ../wails3 ./cmd/wails3

# Run tests to ensure everything works
go test ./...
```

### Créer une branche de fonctionnalité
Créez une branche au nom descriptif :

```bash
git checkout -b feat/window-transparency-support
```

### Implémenter la fonctionnalité
- Respectez nos [normes de codage](/contributing/standards/)
- Limitez les modifications à la fonctionnalité
- Écrivez du code propre et documenté
- Ajoutez des tests complets

```bash
# Example: Adding a new window method
# 1. Add to window.go interface
# 2. Implement in platform files (darwin, windows, linux)
# 3. Add tests
# 4. Update documentation
```

### Effectuer des tests approfondis
Testez votre fonctionnalité :

```bash
# Unit tests
go test ./pkg/application -v

# Integration test - create a test app
cd ..
./wails3 init -n feature-test
cd feature-test
# Add code using your new feature
../wails3 dev
```

### Documenter votre fonctionnalité
- Ajoutez des commentaires de documentation à toutes les API publiques
- Mettez à jour la documentation pertinente dans `/docs/mpress/content/`
- Ajoutez des exemples, le cas échéant

### Respecter les conventions de commit
Utilisez les commits conventionnels :

```bash
git commit -m "feat: add window transparency support

- Add SetTransparent() method to Window API
- Implement for macOS, Windows, and Linux
- Add tests and documentation

Closes #456"
```

### Soumettre une pull request
Poussez votre branche et créez une pull request :

```bash
git push origin feat/window-transparency-support
```

Dans votre pull request :

- Décrivez la fonctionnalité et ses cas d’utilisation
- Présentez des exemples ou des captures d’écran
- Répertoriez toutes les modifications incompatibles
- Référencez la pull request de la WEP acceptée

### Itérer en fonction de la révision
Les mainteneurs peuvent demander des modifications. Faites preuve de patience et de coopération.

@end

[Documentation]
Les PR de correction sont les bienvenues sans qu’il soit nécessaire d’ouvrir d’abord une issue. Les corrections portant uniquement sur la documentation n’ont pas besoin d’un test de code en échec. Suivez [Rédiger la documentation avec M-Press](/contributing/documentation/) pour connaître les étapes d’installation de M-Press, les chemins des sources, la prévisualisation, la validation et la création d’une PR.

@end

## Trouver des issues sur lesquelles travailler

- Recherchez les labels [`good first issue`](https://github.com/wailsapp/wails/labels/good%20first%20issue)
- Consultez les issues [`help wanted`](https://github.com/wailsapp/wails/labels/help%20wanted)
- Parcourez les [issues ouvertes](https://github.com/wailsapp/wails/issues) et demandez à ce qu’une issue vous soit attribuée

## Obtenir de l’aide

- **Discord :** rejoignez le [Discord de Wails](https://discord.gg/JDdSxwjhGf)
- **Discussions :** publiez un message dans les [discussions GitHub](https://github.com/wailsapp/wails/discussions)
- **Issues :** ouvrez une issue pour un bug reproductible ; utilisez les discussions pour les questions et une PR de WEP pour les améliorations

## Code de conduite

Soyez respectueux, constructifs et accueillants. Nous bâtissons une communauté conviviale dont l’objectif est de créer ensemble d’excellents logiciels.

## Étapes suivantes

- Configurez votre [environnement de développement](/contributing/setup/)
- Consultez nos [normes de codage](/contributing/standards/)
- Explorez la [documentation technique](/contributing/overview/)
