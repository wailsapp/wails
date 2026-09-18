---
title: "Système de liaison"
description: "Fonctionnement de la collecte, du traitement et de la génération de code JavaScript/TypeScript par le système de liaison"
slug: "contributing/architecture/bindings"
sourcePath: "contributing/architecture/bindings.md"
---

Ce guide explique le fonctionnement interne du système de liaison de Wails afin d’aider les développeurs à comprendre les mécanismes de la génération automatique de code.

## Vue d’ensemble de l’architecture

Le système de liaison de Wails comprend trois composants principaux :

1. **Collecte** : analyse le code Go pour extraire les informations sur les services, les modèles et les autres déclarations
2. **Configuration** : gère les paramètres et les options du processus de génération des liaisons
3. **Rendu** : génère du code JavaScript/TypeScript à partir des informations collectées

@filetree
- internal/generator/
  - collect/     # Analyse des packages et extraction des informations
  - config/      # Structures et interfaces de configuration
  - render/      # Génération de code JS/TS
@end

## Processus de collecte

Le processus de collecte analyse les packages Go et extrait les informations sur les services, les modèles et les autres déclarations. Cette tâche est assurée par le package `collect`.

### Composants principaux

- **Collector** : gère les informations sur les packages et met en cache les données collectées
- **Package** : représente un package Go en cours d’analyse et stocke les services, modèles et directives collectés
- **Service** : collecte les informations sur les types de services et leurs méthodes
- **Model** : collecte des informations détaillées sur les types de modèles, notamment leurs champs, leurs valeurs et leurs paramètres de type
- **Directive** : analyse et interprète les directives `//wails:` présentes dans le code source Go

### Déroulement de la collecte

1. Le collecteur analyse les packages Go spécifiés dans le projet
2. Il identifie les types de services, c’est-à-dire les structures dotées de méthodes qui seront exposées au frontend
3. Pour chaque service, il collecte les informations sur ses méthodes
4. Il identifie les types de modèles, c’est-à-dire les structures utilisées comme paramètres ou valeurs de retour dans les méthodes de service
5. Pour chaque modèle, il collecte les informations sur ses champs et ses paramètres de type
6. Il traite toutes les directives `//wails:` trouvées dans le code

## Processus de rendu

Le processus de rendu génère du code JavaScript/TypeScript à partir des informations collectées. Cette tâche est assurée par le package `render`.

### Composants principaux

- **Renderer** : orchestre le rendu des fichiers de services, de modèles et d’index
- **Module** : représente un module JavaScript/TypeScript généré
- **Templates** : modèles de texte utilisés pour générer le code

### Déroulement du rendu

1. Pour chaque service, le moteur de rendu génère un fichier JavaScript/TypeScript contenant des fonctions qui reproduisent les méthodes du service
2. Pour chaque modèle, le moteur de rendu génère une classe JavaScript/TypeScript qui reproduit la structure du modèle
3. Le moteur de rendu génère des fichiers d’index qui réexportent tous les services et modèles
4. Le moteur de rendu applique toutes les injections de code personnalisées spécifiées par les directives `//wails:inject`

## Mappage des types

La manière dont les types Go sont mappés aux types JavaScript/TypeScript constitue l’un des aspects les plus importants du système de liaison. Voici un récapitulatif de ce mappage :

| Type Go | Type JavaScript | Type TypeScript |
| --- | --- | --- |
| `bool` | `boolean` | `boolean` |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` | `number` | `number` |
| `string` | `string` | `string` |
| `[]byte` | `Uint8Array` | `Uint8Array` |
| `[]T` | `Array<T>` | `T[]` |
| `map[string]V` | `Object` | `{ [_: string]: V }` |
| `map[K]V` (`K` non chaîne) | `Object` | `{ [_ in K]?: V }` |
| `struct` | `Object` | Classe personnalisée |
| `interface{}` | `any` | `any` |
| `*T` | `T \| null` | `T \| null` |
| `func` | Non pris en charge | Non pris en charge |
| `chan` | Non pris en charge | Non pris en charge |

## Système de directives

Le système de liaisons prend en charge plusieurs directives permettant de personnaliser le code généré. Ces directives sont ajoutées sous forme de commentaires dans votre code Go.

### Directives disponibles

- `//wails:inject` : injecte du code JavaScript/TypeScript personnalisé dans les liaisons générées
- `//wails:include` : inclut des fichiers supplémentaires dans les liaisons générées
- `//wails:internal` : marque un type ou une méthode comme interne, ce qui empêche son exportation vers le frontend
- `//wails:ignore` : ignore complètement une méthode lors de la génération des liaisons
- `//wails:id` : spécifie un identifiant personnalisé pour une méthode, en remplaçant l’identifiant par défaut basé sur un hachage

### Traitement des directives

1. Pendant la phase de collecte, le collecteur identifie et analyse les directives présentes dans le code Go
2. Les directives sont stockées avec les déclarations correspondantes (services, méthodes, modèles, etc.)
3. Pendant la phase de rendu, le moteur de rendu applique les directives afin de personnaliser le code généré

## Fonctionnalités avancées

### Génération conditionnelle de code

Le système de liaisons prend en charge la génération conditionnelle de code au moyen d’un préfixe de condition à deux caractères pour les directives `include` et `inject` :

```
<language><style>:<content>
```

Où :

- `<language>` peut être :
  - `*` - JavaScript et TypeScript
  - `j` - JavaScript uniquement
  - `t` - TypeScript uniquement


- `<style>` peut être :
  - `*` - Classes et interfaces
  - `c` - Classes uniquement
  - `i` - Interfaces uniquement


Par exemple :

```go
//wails:inject j*:console.log("JavaScript only");
//wails:inject t*:console.log("TypeScript only");
```

### Identifiants de méthode personnalisés

Par défaut, les méthodes sont identifiées par un identifiant basé sur un hachage. Vous pouvez toutefois spécifier un identifiant personnalisé à l’aide de la directive `//wails:id` :

```go
//wails:id 42
func (s *Service) CustomIDMethod() {}
```

Cela peut être utile pour préserver la compatibilité lors de la refactorisation du code.

## Considérations relatives aux performances

Le générateur de liaisons est conçu pour être efficace, mais quelques points sont à prendre en compte :

1. La première exécution sera plus lente, car elle constitue un cache des paquets à analyser
2. Les exécutions suivantes seront plus rapides, car elles utiliseront les informations mises en cache
3. Le générateur traite tous les paquets du projet, ce qui peut prendre du temps pour les projets de grande taille
4. Vous pouvez utiliser l’option `-clean` pour nettoyer le répertoire de sortie avant la génération

## Débogage

Si vous rencontrez des problèmes lors de la génération des liaisons, vous pouvez utiliser l’option `-v` pour activer la sortie de débogage :

```bash
wails3 generate bindings -v
```

Vous obtiendrez ainsi des informations détaillées sur le processus de collecte et de rendu, ce qui peut vous aider à déterminer l’origine du problème.
