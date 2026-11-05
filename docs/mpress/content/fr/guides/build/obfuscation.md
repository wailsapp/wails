---
title: "Builds obfusqués"
description: "Compilez votre application Wails avec Garble afin de protéger votre code source contre la rétro-ingénierie"
slug: "guides/build/obfuscation"
sourcePath: "guides/build/obfuscation.md"
---

[Garble](https://github.com/burrowers/garble) est un outil de compilation Go qui remplace `go build` afin de renommer les symboles, d’obfusquer les constantes et de supprimer les informations de débogage du binaire obtenu. Wails v3 prend directement en charge Garble grâce à deux nouvelles commandes.

## Prérequis

- **Go 1.26.2 ou version ultérieure** — requis par Garble v0.16.0
- **Garble v0.16.0**

```bash
go install mvdan.cc/garble@v0.16.0
```

@note{type="tip"}
La version minimale de Go requise par Garble varie selon les versions. Si vous utilisez une ancienne chaîne d’outils Go, consultez la [page des versions de Garble](https://github.com/burrowers/garble/releases) avant l’installation afin de trouver une version compatible avec votre chaîne d’outils.

@end

## Obligatoire : ajoutez des balises JSON aux types de vos services

Toute structure renvoyée ou acceptée par une méthode de service liée doit comporter des balises JSON explicites sur chacun de ses champs exportés :

```go
// Without tags — breaks under Garble
type OrderSummary struct {
    ID        int
    Total     float64
    LineItems []LineItem
}

// With tags — safe under Garble
type OrderSummary struct {
    ID        int       `json:"id"`
    Total     float64   `json:"total"`
    LineItems []LineItem `json:"lineItems"`
}
```

@note{type="caution"}
Garble renomme les champs exportés des structures, et Wails transmet ces structures à `json.Marshal` au moyen d’un paramètre `interface{}` que Garble ne peut pas suivre par analyse statique. Un build obfusqué sans balises JSON sera compilé correctement, mais votre frontend recevra à l’exécution des noms de champs altérés ou vides. Ajoutez les balises avant de lancer un build obfusqué.

@end

Les propres types de Wails — `Screen`, `Rect`, `Point`, `Size`, `EnvironmentInfo`, `OSInfo`, `Capabilities` — comportent déjà des balises. Vous devez uniquement ajouter des balises à vos propres types.

## Compiler avec obfuscation

@steps
### Générer le fichier d’identifiants stables
Exécutez cette commande chaque fois que vous ajoutez, renommez ou supprimez une méthode de service liée :

```bash
wails3 generate bindings -obfuscated
```

Cette commande crée `wails_obfuscated.gen.go` dans le répertoire de votre package principal — validez ce fichier dans le dépôt.

### Compiler avec Garble
```bash
wails3 build --obfuscated
```

Compile l’application en utilisant les liaisons obfusquées.

@end

## Transmettre des options supplémentaires à Garble

Utilisez `--garbleargs` pour transmettre directement des options à `garble` :

```bash
# Obfuscate string literals and reduce binary size
wails3 build --obfuscated --garbleargs "-literals -tiny"

# Reproducible output — same seed produces the same binary
wails3 build --obfuscated --garbleargs "-seed=deadbeef"
```

Consultez la [documentation de Garble](https://github.com/burrowers/garble#flags) pour connaître la liste complète des options prises en charge.

## Avancé : écrire le fichier d’identifiants dans un autre package

Par défaut, `wails_obfuscated.gen.go` est écrit à côté de votre package `main`. Si votre projet conserve les services dans un sous-package importé par `main`, vous pouvez plutôt y écrire le fichier avec `-obfuscated-output` :

```bash
wails3 generate bindings -obfuscated -obfuscated-output ./internal/services
```

@note{type="caution"}
Le package de destination doit être importé — directement ou transitivement — par votre package `main` afin que son `init()` s’exécute au démarrage. Si le package n’est pas accessible, les identifiants stables ne sont jamais enregistrés et les appels de liaisons échouent (par exemple, avec des erreurs `binding not found` à l’exécution).

@end

## Dépannage

### `garble: command not found`

Garble n’est pas installé ou `$(go env GOPATH)/bin` ne figure pas dans votre `PATH`.

```bash
go install mvdan.cc/garble@v0.16.0
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Le frontend reçoit des valeurs de champs incorrectes ou vides

Les types de retour de vos services ne comportent pas de balises `json:"..."`. Vérifiez chaque structure renvoyée par vos méthodes liées et ajoutez des balises explicites à chaque champ exporté.

### Erreur `binding not found` dans la console du navigateur

Le fichier d’identifiants stables est absent ou n’a pas été inclus dans la compilation. Vérifiez les points suivants :

- `wails_obfuscated.gen.go` existe dans le répertoire de votre package principal (ou dans le répertoire transmis à `-obfuscated-output`)
- Vous avez exécuté `wails3 build --obfuscated`, qui ajoute la balise de compilation `wails_obfuscated`
- Si vous avez utilisé `-obfuscated-output`, le package de destination est importé par `main`

### Windows Defender signale le build comme un virus

Pendant la compilation, Windows Defender signale de manière heuristique les binaires Go obfusqués avec Garble, car ils ne comportent pas de symboles de débogage et ressemblent à des exécutables compressés. La compilation échoue avec le message suivant :

```
open C:\Users\...\AppData\Local\Temp\go-build...\a.out.exe: The file contains a virus or potentially unwanted software.
```

Ajoutez votre répertoire temporaire (dans lequel Go écrit les artefacts de compilation intermédiaires) et le répertoire de votre projet à la liste d’exclusions de Defender :

```powershell
Add-MpPreference -ExclusionPath "$env:TEMP"
Add-MpPreference -ExclusionPath "C:\path\to\your\project"
```

Ces exclusions s’appliquent uniquement aux chemins indiqués et ne désactivent pas Defender globalement.

### La compilation échoue avec `unsupported Go version`

Garble v0.16.0 nécessite Go 1.26.2 ou version ultérieure. Mettez Go à niveau ou consultez la [page des versions de Garble](https://github.com/burrowers/garble/releases) afin de trouver une version compatible avec votre chaîne d’outils.
