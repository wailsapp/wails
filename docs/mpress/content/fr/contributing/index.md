---
title: "Contribuer"
description: "Contribuer à Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Bienvenue aux contributeurs !

Vos contributions à Wails sont les bienvenues ! Que vous corrigiez des bogues, ajoutiez des fonctionnalités ou amélioriez la documentation, votre aide est précieuse.

## Façons de contribuer

### 1. Signaler des problèmes

Vous avez trouvé un bogue ? [Ouvrez un ticket](https://github.com/wailsapp/wails/issues/new) en indiquant :

- Une description claire
- Les étapes de reproduction
- Le comportement attendu et le comportement réel
- Les informations système
- Des exemples de code

### 2. Améliorer la documentation

Les PR de correction sont les bienvenues sans ticket préalable ni échec d’un test de code.  
Suivez [Rédiger de la documentation avec M-Press](/contributing/documentation/) pour le téléchargement, l’installation sous macOS/Linux/Windows, l’aperçu en direct et la soumission de votre modification.

Les améliorations de la documentation sont toujours les bienvenues :

- Corriger les coquilles et les erreurs
- Ajouter des exemples
- Clarifier les explications
- Traduire le contenu

### 3. Soumettre du code

Contribuez au code au moyen de pull requests :

- Corrections de bogues
- Nouvelles fonctionnalités
- Améliorations des performances
- Tests

### 4. Proposer une amélioration (WEP)

Les nouvelles fonctionnalités et les modifications du comportement public suivent le processus Wails Enhancement Proposal (WEP). Celui-ci garantit la transparence du développement des fonctionnalités et veille à ce que chaque proposition acceptée ait une personne chargée de sa mise en œuvre. N’ouvrez pas de ticket de demande de fonctionnalité.

1. Vous pouvez éventuellement présenter votre idée dans la catégorie [Idées](https://github.com/wailsapp/wails/discussions/categories/ideas) de GitHub Discussions ou sur [Discord](https://discord.gg/JDdSxwjhGf) afin d’évaluer l’intérêt qu’elle suscite.
2. Copiez [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/WEP_TEMPLATE.md) dans `v3/wep/proposals/<proposal name>/proposal.md` et remplissez chaque section.
3. Ouvrez une pull request en brouillon intitulée `[WEP] <title>` qui contient uniquement la proposition. La PR est l’espace officiel pour en discuter.
4. Recueillez des commentaires et des soutiens (commentaires et réactions « pouce levé » sur la PR). Prévoyez au moins deux semaines de discussion et convenez de la personne qui mettra en œuvre la proposition.
5. Marquez la PR comme prête pour la revue. La décision finale revient aux mainteneurs : les propositions acceptées reçoivent un numéro WEP et sont fusionnées.

Le processus complet est décrit dans [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/releases/v3.0.0/v3/wep/README.md).

## Bien démarrer

### Créer un fork et cloner le dépôt

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails

# Add upstream remote
git remote add upstream https://github.com/wailsapp/wails.git
```

### Compiler depuis les sources

```bash
# Build the v3 CLI. Go downloads any required modules automatically.
cd v3
go build -o ../wails3 ./cmd/wails3

# Confirm the built CLI runs and reports its version.
../wails3 version
```

### Exécuter les tests

Les tests font partie intégrante de la modification ; ils ne constituent pas une simple vérification finale. Placez les tests unitaires à côté du code qu’ils exercent et privilégiez les tests pilotés par table lorsqu’un même comportement est vérifié avec plusieurs entrées ou cas limites. Nommez chaque cas afin qu’un échec indique clairement le scénario concerné.

Toute logique nouvelle ou modifiée doit être entièrement couverte. Visez une couverture de 100 % des instructions Go du code ajouté ou modifié par votre PR ; ne considérez pas le pourcentage à l’échelle du dépôt comme un substitut aux tests de la modification. Une lacune peut être légitime — par exemple, un chemin d’erreur propre à un système d’exploitation ou une condition impossible à reproduire en pratique sans matériel réel — mais expliquez-la dans la description de la PR et précisez pourquoi elle ne peut raisonnablement pas être testée.

```bash
# Run all v3 tests
cd v3
go test ./...

# Run specific package tests
go test ./pkg/application

# Inspect coverage for the packages you changed
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
cd ..
```

Pour les suites d’intégration, la détection des conditions de concurrence et l’ensemble des commandes équivalentes à celles de la CI, consultez [Tests et intégration continue](/contributing/testing-ci/).

## Apporter des modifications

### Créer une branche

```bash
# Update master
git checkout master
git pull upstream master

# Create feature branch
git checkout -b feature/my-feature
```

### Effectuer vos modifications

1. **Écrivez le code** conformément aux conventions Go
2. **Ajoutez des tests** pour les nouvelles fonctionnalités
3. **Mettez à jour la documentation** si nécessaire
4. **Exécutez les tests** pour vérifier que rien ne casse
5. **Validez les modifications** avec des messages clairs

### Consignes relatives aux commits

```bash
# Good commit messages
git commit -m "fix: resolve window focus issue on macOS"
git commit -m "feat: add support for custom window chrome"
git commit -m "docs: improve bindings documentation"

# Use conventional commits:
# - feat: New feature
# - fix: Bug fix
# - docs: Documentation
# - test: Tests
# - refactor: Code refactoring
# - chore: Maintenance
```

### Soumettre une pull request

```bash
# Push to your fork
git push origin feature/my-feature

# Open pull request on GitHub
# Provide clear description
# Reference related issues
```

## Consignes relatives aux pull requests

### Bonne description de PR

```markdown
## Description
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- Tested on macOS 14
- Tested on Windows 11
- All tests passing

## Related Issues
Fixes #123
```

### Liste de contrôle de la PR

- [ ] Le code respecte les conventions Go
- [ ] Tests ajoutés ou mis à jour
- [ ] Documentation mise à jour
- [ ] Tous les tests réussissent
- [ ] Aucune modification incompatible, ou celle-ci est documentée
- [ ] Messages de commit clairs

## Consignes relatives au code

### Style du code Go

```go
// ✅ Good: Clear, documented, tested
// ProcessData processes the input data and returns the result.
// It returns an error if the data is invalid.
func ProcessData(data string) (string, error) {
    if data == "" {
        return "", errors.New("data cannot be empty")
    }
    
    result := process(data)
    return result, nil
}

// ❌ Bad: No docs, no error handling
func ProcessData(data string) string {
    return process(data)
}
```

### Tests

```go
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "processed", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessData() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Documentation

### Rédiger la documentation

Suivez [Rédiger de la documentation avec M-Press](/contributing/documentation/) pour connaître le flux de travail complet de la documentation. Ce guide couvre la version épinglée, les options d’installation de l’exécutable et de Go, les fichiers sources modifiables, la navigation, les vérifications des traductions et la soumission d’une PR vers `master`.

Les travaux portant uniquement sur la documentation ne nécessitent ni compilation d’une application Wails ni dépendances propres à une plateforme. Depuis la racine du dépôt, M-Press affiche un aperçu du site existant avec `mpress dev`.

### Style de la documentation

- Utiliser l’orthographe de l’anglais international
- Commencer par le problème
- Fournir des exemples fonctionnels
- Inclure des conseils de dépannage
- Ajouter des renvois vers le contenu connexe

## Communauté

### Obtenir de l’aide

- **Discord :** [Rejoignez notre communauté](https://discord.gg/JDdSxwjhGf)
- **GitHub Discussions :** posez vos questions
- **GitHub Issues :** signalez des bogues

### Code de conduite

Faites preuve de respect, d’inclusion et de professionnalisme. Nous sommes tous ici pour créer ensemble d’excellents logiciels. Consultez le [Code de conduite](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md) pour en savoir plus.

## Reconnaissance

Les contributeurs sont mentionnés dans :

- Les notes de version
- La liste des contributeurs
- Les statistiques GitHub

Merci de contribuer à Wails ! 🎉

## Étapes suivantes

@cards{cols="2"}
◆ Dépôt GitHub
Consultez le dépôt Wails.

[Voir sur GitHub →](https://github.com/wailsapp/wails)

---
◆ Communauté Discord
Rejoignez la communauté.

[Rejoindre Discord →](https://discord.gg/JDdSxwjhGf)

---
📖 Documentation
Consultez la documentation.

[Parcourir la documentation →](/quick-start/why-wails/)

@end
