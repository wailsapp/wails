---
title: "Configuration"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Expérimental"}
L’assistant de configuration est nouveau et a principalement été testé sous Linux. Si vous rencontrez des problèmes, veuillez les [signaler](https://github.com/wailsapp/wails/issues/4904) et suivre plutôt les [étapes d’installation manuelle](/getting-started/installation/#platform-specific-dependencies).

@end

## Démarrage rapide

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

L’assistant s’ouvre dans votre navigateur et vous guide dans la vérification des dépendances, la définition des valeurs par défaut du projet et la configuration facultative de la compilation multiplateforme.

Vous êtes alors prêt à créer votre premier projet :

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Fonctionnement

- **Vérification des dépendances** - Vérifie Go, npm et les outils propres à la plateforme
- **Configuration des valeurs par défaut** - Informations sur l’auteur, préfixe de l’identifiant de bundle et modèles préférés
- **Compilations multiplateformes** - Configuration Docker facultative pour compiler depuis n’importe quel système hôte
- **Signature du code** - Configuration facultative pour macOS, Windows et Linux

La configuration est enregistrée dans `$XDG_CONFIG_HOME/wails/defaults.yaml`, avec repli sur `~/.config/wails/defaults.yaml`, et utilisée par `wails3 init`. Elle est distincte du fichier `build/config.yml` d’un projet.

## Sous-commandes

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Vous rencontrez des problèmes ?

1. Exécutez `wails3 doctor` pour diagnostiquer les problèmes
2. Suivez les [étapes d’installation manuelle](/getting-started/installation/#platform-specific-dependencies)
3. [Signalez le problème](https://github.com/wailsapp/wails/issues/4904) en joignant la sortie de `wails3 doctor`
