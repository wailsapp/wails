---
title: "Configuration"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Expérimental"}
L’assistant de configuration est nouveau et a été principalement testé sous Linux. Si vous rencontrez des problèmes, veuillez les [signaler](https://github.com/wailsapp/wails/issues/4904) et suivre plutôt les [étapes d’installation manuelle](/getting-started/installation/#platform-specific-dependencies).

@end

## Démarrage rapide

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

L’assistant s’ouvre dans votre navigateur et vous guide dans la vérification des dépendances, la définition des valeurs par défaut globales des projets, la configuration des outils multiplateformes facultatifs et la saisie des identifiants de signature. Ces préférences peuvent être réutilisées dans plusieurs projets.

Vous pouvez ensuite créer votre premier projet :

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Fonctionnement

- **Vérification des dépendances** - Vérifie Go, npm et les outils de la plateforme
- **Configuration des valeurs par défaut** - Informations sur l’auteur, préfixe de l’ID de bundle et modèles préférés
- **Compilations multiplateformes** - Configuration facultative de Docker pour compiler depuis n’importe quel hôte
- **Signature du code** - Configuration facultative pour macOS, Windows et Linux

La configuration est enregistrée dans `$XDG_CONFIG_HOME/wails/defaults.yaml`, avec repli sur `~/.config/wails/defaults.yaml`, puis utilisée par `wails3 init`. Elle est distincte du fichier `build/config.yml` d’un projet.

## Créer un projet Wake

Lorsque l’expérience Wake est activée, `wails3 init -ui` recueille vos choix de compilation propres au projet et les enregistre dans `wails.yaml` si la création aboutit. Suivez le [tutoriel de prise en main de Wake](/experimental/wake/) pour choisir un moteur de compilation, des cibles et des formats de paquet lors de la création d’un projet d’application de bureau. `wails3 setup` conserve son rôle de gestion des préférences globales.

Les projets Wake existants sont configurés dans `wails.yaml`. Utilisez le [guide de compilation YAML](/experimental/yaml-builds/) pour modifier leurs profils et leurs paramètres de création de paquets. Avant de pouvoir utiliser Wake, les projets Taskfile existants doivent faire l’objet de la [migration irréversible contrôlée](/experimental/wake/#migrate-an-existing-project) ; l’assistant de configuration ne traduit ni n’exécute leurs Taskfiles.

## Réutiliser les préférences de signature

Pour les compilations Wake signées, les paramètres de signature manquants sont repris des valeurs par défaut de votre compte utilisateur, y compris les fichiers de certificats situés hors du projet. Les paramètres explicites du projet sont prioritaires. Les clés privées, les mots de passe et les chemins globaux vers les identifiants ne sont pas copiés dans `wails.yaml`.

Pour les clés Windows et Linux configurées globalement, fournissez le mot de passe au moyen du magasin d’identifiants existant ou de `WAILS_WINDOWS_CERT_PASSWORD` / `WAILS_PGP_PASSWORD`. Pour une identité Apple chiffrée, définissez `darwin.signing.credential` sur le nom d’une variable d’environnement contenant son mot de passe.

## Sous-commandes

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Vous rencontrez des problèmes ?

1. Exécutez `wails3 doctor` pour diagnostiquer les problèmes
2. Suivez les [étapes d’installation manuelle](/getting-started/installation/#platform-specific-dependencies)
3. [Signalez le problème](https://github.com/wailsapp/wails/issues/4904) en joignant la sortie de `wails3 doctor`
