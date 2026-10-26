---
title: "Rédiger de la documentation avec M-Press"
description: "Téléchargez et installez M-Press, prévisualisez la documentation de Wails et proposez une contribution à la documentation."
sourcePath: "contributing/documentation.md"
---

Wails v3 utilise **M-Press** pour transformer le Markdown en ce site de documentation. Vous pouvez corriger une faute de frappe ou une traduction, améliorer une explication, ou ajouter un exemple ou un guide. Les corrections portant uniquement sur la documentation ne nécessitent ni ticket préalable ni échec d’un test de code.

Vous avez besoin de Git pour préparer une contribution et d’un compte GitHub pour soumettre une pull request. L’exécutable M-Press téléchargeable comprend le générateur de documentation et le moteur de rendu des diagrammes. La prévisualisation de la documentation ne nécessite ni Go, ni Node.js, ni la CLI Wails, ni les dépendances de compilation des applications de bureau. Python 3 est nécessaire pour les scripts de validation supplémentaires du dépôt.

## Partir d’une page de documentation

Dans la barre latérale, ouvrez **Contribution → Documentation** pour accéder à ce guide. Pour travailler sur une page précise, sélectionnez l’**icône de modification** dans son en-tête, puis **Modifier cette page localement**. Copiez la commande correspondant à votre système d’exploitation et exécutez-la dans un terminal.

La commande utilise M-Press s’il est installé. Sinon, son programme d’installation télécharge une version officielle et en vérifie la somme de contrôle. Elle prépare une copie de travail locale pour votre contribution et ouvre la page sélectionnée dans M-Press. Avant d’exécuter les vérifications du dépôt, comparez `mpress version` à la version épinglée ci-dessous : le programme d’installation automatique peut télécharger une version plus récente.

Vous pouvez modifier la source dans votre propre éditeur pendant que la prévisualisation locale est reconstruite. L’assistant propose également les actions distinctes **Examiner les modifications**, **Valider les modifications**, **Pousser la branche de contribution** et **Ouvrir une pull request en brouillon**. Avant de soumettre votre contribution, suivez les consignes ci-dessous relatives aux vérifications et aux pull requests. La génération des traductions est facultative et peut nécessiter la configuration d’un fournisseur ; ce n’est pas le cas pour la modification et la prévisualisation du Markdown.

Pour obtenir une configuration manuelle reproductible, suivez les étapes ci-dessous.

## Télécharger et installer M-Press

Le script de compilation du dépôt télécharge la **dernière version stable de M-Press**. Utilisez la même version localement : choisissez l’une des méthodes d’installation ci-dessous, exécutez `mpress version`, puis comparez le résultat à la version indiquée sur la page de la version publiée.

### Télécharger un exécutable

Ouvrez la page de la [dernière version stable de M-Press](https://github.com/leaanthony/mpress/releases/latest) ou téléchargez l’archive correspondante ci-dessous. ARM64 inclut Apple Silicon ; AMD64 est l’architecture habituelle des PC Intel/AMD.

| Ordinateur | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

Téléchargez [checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt) depuis la même version publiée. Dans votre répertoire de téléchargement, calculez l’empreinte SHA-256 de l’archive et comparez-la à la ligne correspondant à ce nom de fichier dans `checksums.txt` :

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

Si vous avez choisi une autre architecture, utilisez le nom réel du fichier téléchargé. Une fois les hachages identiques, extrayez et installez l’exécutable :

@tabs
[macOS]
Depuis le répertoire de téléchargement, avec l’archive Apple Silicon :

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Pour les Mac Intel, remplacez `mpress-darwin-amd64.tar.gz`. Ajoutez la ligne `export PATH` à `~/.zshrc` afin que la commande reste disponible dans les nouveaux terminaux.

[Linux]
Depuis le répertoire de téléchargement :

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Pour ARM64, remplacez `mpress-linux-arm64.tar.gz`. Ajoutez la ligne `export PATH` au fichier de démarrage de votre shell, tel que `~/.bashrc`, afin que la commande reste disponible dans les nouveaux terminaux.

[Windows PowerShell]
Depuis le répertoire de téléchargement :

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

Pour ARM64, remplacez `mpress-windows-arm64.zip`. Pour que la commande reste disponible lors des prochaines sessions, ajoutez votre répertoire utilisateur `bin` à la variable **Path** de votre compte dans les **Variables d’environnement** de Windows, puis rouvrez le terminal.
@end

### Installer avec Go

Si vous utilisez déjà Go, vous pouvez plutôt compiler la dernière version stable. Consultez son fichier `go.mod` pour connaître la chaîne d’outils requise ; la sélection automatique de la chaîne d’outils de Go peut la télécharger. La version utilisée pour cette passe de documentation nécessite Go 1.26.10 ou une version ultérieure.

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go installe la commande dans `GOBIN` ou, lorsque `GOBIN` n’est pas défini, dans le répertoire `bin` sous `go env GOPATH`. Ajoutez ce répertoire à votre `PATH` et rouvrez le terminal si `mpress` est introuvable.

## Obtenir les sources de la documentation

[Créez un fork de wailsapp/wails](https://github.com/wailsapp/wails/fork), puis exécutez les commandes suivantes en remplaçant `YOUR_USERNAME` par votre nom d’utilisateur GitHub :

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

Votre branche de documentation part ainsi de la branche **master** du dépôt en amont. Remplacez `docs/my-change` par un nom de branche descriptif de votre choix.

## Lancer la prévisualisation en direct

Depuis la racine du dépôt Wails, où se trouve `mpress.yaml` :

```sh
mpress dev
```

Ouvrez [http://127.0.0.1:3000](http://127.0.0.1:3000) dans votre navigateur. Pour consulter ce guide, accédez à [http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/). Laissez la commande s’exécuter pendant vos modifications ; l’enregistrement des changements apportés aux sources reconstruit la prévisualisation. Avant d’exécuter les commandes de validation ci-dessous, appuyez sur **Ctrl+C** pour l’arrêter.

Le dépôt contient déjà la configuration M-Press. Il n’est pas nécessaire d’initialiser un nouveau site. Pour utiliser un autre port, exécutez `mpress dev --port 3001` et ouvrez l’adresse localhost correspondante.

## Modifier ou ajouter une page

Les pages sources en anglais se trouvent dans `docs/mpress/content/`. Les traductions se trouvent dans des dossiers de langue tels que `fr/` et `id/`. Les instructions de contribution accessibles depuis l’icône de modification indiquent le fichier source de la page actuelle.

Pour une page existante, modifiez son Markdown et conservez les métadonnées YAML situées en haut. Utilisez des titres, des paragraphes, des listes, des liens et des blocs de code délimités. Conservez les composants M-Press appariés, tels que `@note` / `@end` et `@tabs` / `@end`. La [documentation de création de contenu M-Press](https://m-press.me/authoring/) décrit la syntaxe prise en charge.

Pour un nouveau guide, créez par exemple le fichier `docs/mpress/content/guides/my-guide.md` :

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

Ajoutez la page à la section appropriée de `docs/mpress/content/_nav.yaml`, sauf si cette section utilise déjà `autogenerate`. Vérifiez que les lecteurs peuvent y accéder depuis la navigation et les liens des pages connexes. Les nouvelles pages en anglais doivent également disposer de copies complètes pour toutes les langues avant que les contrôles de couverture puissent réussir ; coordonnez-vous avec les responsables de maintenance ou les contributeurs aux traductions.

La sortie générée se trouve dans `docs/mpress/site/`. Validez les pages sources ainsi que toute modification prévue de la navigation ou des ressources statiques ; les fichiers générés du site et les caches locaux doivent rester en dehors de la pull request.

## Vérifier votre contribution

Depuis la racine du dépôt :

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Sous Windows, utilisez `python` si Python 3 est installé ainsi. Si [Task](https://taskfile.dev/) est déjà installé, `task docs:build` exécute les mêmes vérifications en utilisant le `mpress` présent dans votre PATH.

La vérification des traductions échoue si des fichiers de traduction sont manquants, vides ou orphelins. Par défaut, les passages traduits obsolètes ou endommagés sont signalés comme avertissements ; `--strict` fait également échouer la vérification pour ces problèmes. Vous pouvez soumettre une modification de la prose anglaise existante pendant la synchronisation des traductions. La correction d’une traduction peut ne modifier que la langue concernée. Conservez les commandes, les noms d’API, les liens, les exemples de code et les connexions des diagrammes intacts ; n’utilisez pas de texte anglais provisoire comme contenu traduit.

Lisez la page modifiée dans le navigateur et exécutez tous les exemples de code que vous avez modifiés. Si une vérification échoue sur une page que vous n’avez pas modifiée, joignez la sortie et le contexte pertinent à votre PR au lieu d’assouplir la vérification.

## Soumettre une pull request

Examinez votre diff, validez vos modifications des sources avec un commit, puis poussez votre branche de contribution vers votre fork. Ouvrez une PR ciblant **master** dans `wailsapp/wails`. Décrivez les modifications, ajoutez un lien vers la page concernée et répertoriez les vérifications que vous avez exécutées. Pour les modifications de mise en page, joignez des captures d’écran avant et après ; pour les exemples de code modifiés, indiquez la plateforme et la version.

Le site de production est statique ; la prévisualisation et la validation locales ne nécessitent ni identifiants de déploiement ni prestataire de traduction. Les responsables de la maintenance se chargent de la publication après la revue.

Pour les contributions au code et les propositions de fonctionnalités, consultez [Contribuer à Wails](/contributing/). Pour le fonctionnement interne du framework, consultez la [Présentation technique](/contributing/overview/).
