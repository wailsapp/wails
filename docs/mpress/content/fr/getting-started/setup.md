---
title: "Configuration"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Expérimental"}
L’assistant de configuration est nouveau et a principalement été testé sous Linux. Si vous rencontrez des problèmes, veuillez [les signaler](https://github.com/wailsapp/wails/issues/4904) et suivre plutôt la [procédure d’installation manuelle](/getting-started/installation/#platform-specific-dependencies).

@end

## Démarrage rapide

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

L’assistant s’ouvre dans votre navigateur et vous guide dans la vérification des dépendances, la définition des valeurs par défaut du projet et la configuration facultative de la compilation multiplateforme. Les instructions ci-dessus utilisent le parcours de configuration par défaut ; le parcours expérimental pour les projets Wake est décrit ci-dessous.

Vous pouvez ensuite créer votre premier projet :

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Fonctionnement

- **Vérification des dépendances** — Vérifie Go, npm et les outils de la plateforme
- **Configuration des valeurs par défaut** — Informations sur l’auteur, préfixe de l’identifiant de bundle et modèles préférés
- **Compilations multiplateformes** — Configuration facultative de Docker pour compiler depuis n’importe quel système hôte
- **Signature du code** — Configuration facultative pour macOS, Windows et Linux

La configuration est enregistrée dans `$XDG_CONFIG_HOME/wails/defaults.yaml`, avec repli sur `~/.config/wails/defaults.yaml`, et utilisée par `wails3 init`. Elle est distincte du fichier `build/config.yml` d’un projet.

## Configurer un projet Wake

Pour créer une première application complète, suivez le [tutoriel de prise en main de Wake](/experimental/wake/).
Il vous guide tout au long de la création du projet et de chaque écran de configuration, explique ce que vous pouvez attendre de cette fonctionnalité expérimentale et aborde la migration depuis une configuration de build existante.

Depuis le répertoire de votre projet, lancez la configuration en ayant [activé la fonctionnalité expérimentale Wake](/experimental/yaml-builds/#enable-yaml-builds) :

```bash
export WAILS_EXP_USE_WAKE=1
wails3 setup
```

Dans PowerShell :

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
wails3 setup
```

Cela ouvre un parcours centré sur le projet :

1. **Projet** — saisissez les métadonnées du projet lorsqu’aucun manifeste natif n’existe.
2. **Application** — choisissez une interface graphique, une CLI ou un serveur, un moteur de compilation ou des exécutables déjà compilés. L’utilisation de Go pur est un choix expérimental explicite.
3. **Cibles** — sélectionnez les systèmes d’exploitation et les architectures. iOS et Android restent expérimentaux ; configurez leurs paramètres de signature dans le fichier YAML.
4. **Paquets** — choisissez les formats de paquet et indiquez s’ils doivent être signés.
5. **Outils** — examinez les prérequis du parcours de compilation et de création de paquets sélectionné.
6. **Signature** — configurez la signature lorsqu’elle est sélectionnée, à l’aide des contrôles de signature existants.
7. **Enregistrement** — enregistrez la configuration du projet et affichez la commande à exécuter.

Les outils requis dépendent de vos choix. La compilation d’applications de bureau en Go pur ne nécessite ni Docker ni compilateur C. Les compilations croisées avec CGo nécessitent un compilateur compatible et les bibliothèques de la plateforme cible, ou une chaîne d’outils conteneurisée prise en charge. NSIS nécessite toujours `makensis` ; les AppImages avec interface graphique compilées depuis un autre système hôte nécessitent une arborescence d’exécution Linux fournie. Les identités de signature exportées peuvent utiliser les outils de signature Go de Wails, tandis que les identités stockées dans un trousseau natif ou un magasin de certificats conservent les outils de la plateforme. Consultez les [prérequis de création et de signature des paquets](/experimental/yaml-builds/#package-and-sign-for-distribution).

Les exécutables déjà compilés évitent la compilation Go, l’installation et la compilation du frontend, la génération des bindings ainsi que l’utilisation d’UPX. Ils nécessitent néanmoins les outils de création de paquets et de signature sélectionnés, ainsi que des fichiers binaires d’entrée valides. Consultez les [chemins des entrées précompilées](/experimental/yaml-builds/#use-prebuilt-desktop-executables). Ce choix s’applique à la commande affichée après l’enregistrement ; il ne constitue pas un paramètre de compilation YAML permanent.

L’option Enregistrer écrit `wails.yaml`, met à jour le type d’application, l’activation du frontend et les paramètres sélectionnés du backend, puis remplace la sélection dédiée de cible, de paquet et de signature dans `profiles.setup`. Les autres profils, les commandes du frontend, les exécutables compagnons, les hooks et les paramètres de projet sans rapport sont conservés. Les paramètres du backend s’appliquent à l’ensemble du projet ; les autres profils qui en héritent utilisent donc également le backend sélectionné. Les nouvelles métadonnées du projet restent dans l’assistant jusqu’à ce que vous sélectionniez Enregistrer. Les outils manquants sont signalés afin que vous puissiez les installer avant la compilation ; l’enregistrement d’une configuration valide n’installe aucun outil et ne lance aucune compilation.

Pour les compilations signées, les paramètres de signature manquants sont repris depuis vos valeurs par défaut utilisateur, notamment les fichiers de certificat situés hors du projet. Les paramètres explicites du projet sont prioritaires. La configuration ne copie ni clés privées, ni mots de passe, ni chemins globaux d’identifiants dans `wails.yaml`. Pour les clés Windows et Linux configurées globalement, fournissez le mot de passe par l’intermédiaire du magasin d’identifiants existant ou de `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. Pour une identité Apple chiffrée, définissez `darwin.signing.credential` sur le nom d’une variable d’environnement contenant son mot de passe.

Le profil enregistré s’exécute avec `wails3 build setup` ; l’option pour les exécutables déjà compilés ajoute `--skip-build`.

Laissez Wake activé lorsque vous exécutez la commande affichée. Un fichier `wails.yaml` actif ne se rabat jamais sur les Taskfiles lorsque l’indicateur n’est pas défini. La configuration ne convertit ni les Taskfiles existants ni les anciennes métadonnées : suivez le [processus de migration vérifié](/experimental/yaml-builds/#migrate-an-existing-project) avant de configurer ces projets avec Wake.

## Sous-commandes

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Vous rencontrez des problèmes ?

1. Exécutez `wails3 doctor` pour diagnostiquer les problèmes
2. Suivez la [procédure d’installation manuelle](/getting-started/installation/#platform-specific-dependencies)
3. [Signalez le problème](https://github.com/wailsapp/wails/issues/4904) en joignant la sortie de `wails3 doctor`
