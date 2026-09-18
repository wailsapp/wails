---
title: "Mettre à niveau depuis une alpha de v3"
description: "Passer un projet Wails v3 alpha existant à une version bêta fixée"
slug: "migration/alpha-to-beta"
sourcePath: "migration/alpha-to-beta.md"
---

Ce guide s’adresse aux projets v3 alpha existants. Pour Wails v2, utilisez le [guide de v2 vers v3](/migration/v2-to-v3/).

## Avant la mise à niveau

Enregistrez un commit ou sauvegardez votre projet. Lisez le [journal des modifications](/changelog/) entre votre version alpha et la bêta choisie : le code source, les API ou la configuration de compilation peuvent nécessiter des changements. Vérifiez la [politique de compatibilité des applications de bureau](/status/) et les prérequis de votre plateforme.

Les commandes ci-dessous utilisent la version publiée `v3.0.0-beta.23` comme exemple de version exacte, sans recommander de suivre systématiquement la dernière version. Si vous choisissez une autre version, vérifiez les versions de sa CLI, de son module Go et de son runtime npm, puis adaptez toutes les commandes. Dans cet exemple, la version npm correspond à la version Go sans son préfixe `v`.

## 1. Mettre à jour la CLI

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.23
wails3 version
```

Vérifiez que `wails3 version` indique la version installée. Un ancien exécutable placé plus tôt dans le `PATH` peut masquer la nouvelle CLI.

## 2. Mettre à jour le module Go

Exécutez les commandes à la racine du projet. Examinez les changements de dépendances ; ne mettez pas globalement à jour les modules sans rapport avec cette opération.

```sh
go get github.com/wailsapp/wails/v3@v3.0.0-beta.23
go mod tidy
```

## 3. Mettre à jour le runtime frontend

Pour les projets utilisant npm et un répertoire `frontend` :

```sh
cd frontend
npm install --save-exact @wailsio/runtime@3.0.0-beta.23
cd ..
```

Conservez le fichier de verrouillage et examinez ses modifications. Si votre frontend utilise un autre gestionnaire de paquets ou répertoire, adaptez cette étape en conservant une version exacte du runtime.

## 4. Régénérer, compiler et tester

À la racine du projet, régénérez les liaisons à partir de vos services Go et compilez :

```sh
wails3 generate bindings
wails3 build
```

Lancez l’application compilée et testez vos parcours sur chaque plateforme prise en charge que vous distribuez. Examinez et enregistrez ensemble dans un commit le code source, les liaisons générées, les fichiers du module et les modifications du fichier de verrouillage frontend.

## En cas d’échec de la mise à niveau

Vérifiez la CLI dans le `PATH`, la version du module avec `go list -m github.com/wailsapp/wails/v3` et le runtime installé avec `npm --prefix frontend ls @wailsio/runtime`. Régénérez les liaisons après avoir corrigé les incompatibilités de versions. Ne supposez pas que toutes les versions alpha peuvent être mises à niveau sans modifier le code.

Si le problème persiste, [signalez un problème reproductible](https://github.com/wailsapp/wails/issues/new/choose) avec les anciennes et nouvelles versions, l’erreur exacte et la sortie de `wails3 doctor`. Suivez la [politique de sécurité](https://github.com/wailsapp/wails/blob/master/SECURITY.md) pour les vulnérabilités.
