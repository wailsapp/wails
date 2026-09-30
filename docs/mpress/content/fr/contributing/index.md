---
title: "Contribuer"
description: "Contribuer à Wails"
slug: "contributing"
sourcePath: "contributing/index.md"
---

## Bienvenue aux contributeurs !

Nous accueillons avec plaisir les contributions à Wails ! Que vous corrigiez des bugs, ajoutiez des fonctionnalités ou amélioriez la documentation, votre aide est appréciée.

## Façons de contribuer

### 1. Signaler des problèmes

Vous avez trouvé un bug ? [Ouvrez un ticket](https://github.com/wailsapp/wails/issues/new) en fournissant :

- Une description claire
- Les étapes de reproduction
- Le comportement attendu et le comportement observé
- Les informations système
- Des exemples de code

### 2. Améliorer la documentation

Les PR de correction sont les bienvenues sans ticket préalable ni test de code en échec.  
Suivez la procédure [Corriger la documentation](/contributing/documentation/) pour prévisualiser  
et valider une modification avec M-Press.

Les améliorations de la documentation sont toujours les bienvenues :

- Corriger les fautes de frappe et les erreurs
- Ajouter des exemples
- Clarifier les explications
- Traduire le contenu

### 3. Soumettre du code

Contribuez au code au moyen de pull requests :

- Corrections de bugs
- Nouvelles fonctionnalités
- Améliorations des performances
- Tests

### 4. Proposer une amélioration (WEP)

Les nouvelles fonctionnalités et les modifications du comportement public suivent le processus de proposition d’amélioration de Wails (Wails Enhancement Proposal, WEP). Ce processus garantit la transparence du développement des fonctionnalités et qu’un responsable de l’implémentation est désigné pour chaque proposition acceptée. N’ouvrez pas de ticket de demande de fonctionnalité.

1. Vous pouvez éventuellement présenter votre idée dans la catégorie [Idées](https://github.com/wailsapp/wails/discussions/categories/ideas) des discussions GitHub ou sur [Discord](https://discord.gg/JDdSxwjhGf) afin d’évaluer l’intérêt qu’elle suscite.
2. Copiez [`v3/wep/WEP_TEMPLATE.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/WEP_TEMPLATE.md) dans `v3/wep/proposals/<proposal name>/proposal.md` et remplissez chaque section.
3. Ouvrez une pull request en brouillon intitulée `[WEP] <title>` et contenant uniquement la proposition. Cette PR constitue l’espace officiel pour en discuter.
4. Recueillez des retours et des marques de soutien (commentaires et réactions 👍 sur la PR). Prévoyez au moins deux semaines de discussion et convenez de la personne qui implémentera la proposition.
5. Indiquez que la PR est prête à être examinée. Les responsables de la maintenance prennent la décision finale : les propositions acceptées reçoivent un numéro WEP et sont fusionnées.

Le processus complet est documenté dans [`v3/wep/README.md`](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md).

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

Les tests font partie intégrante de la modification et ne constituent pas une simple vérification finale. Placez les tests unitaires à côté du code qu’ils testent et privilégiez les tests pilotés par table lorsqu’un même comportement est vérifié avec plusieurs entrées ou cas limites. Nommez chaque cas afin que tout échec indique clairement le scénario concerné.

La logique nouvelle ou modifiée devrait être entièrement couverte. Visez une couverture de 100 % des instructions Go pour le code ajouté ou modifié par votre PR ; ne considérez pas le pourcentage de couverture de l’ensemble du dépôt comme un substitut aux tests de la modification. Une lacune peut être justifiée, par exemple pour un chemin d’erreur propre à un système d’exploitation ou une condition impossible à reproduire en pratique sans matériel réel, mais expliquez cette lacune et pourquoi elle ne peut raisonnablement pas être testée dans la description de la PR.

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

Pour les suites d’intégration, la détection des courses aux données et toutes les commandes équivalentes à celles de la CI, consultez [Tests et intégration continue](/contributing/testing-ci/).

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

1. **Écrivez du code** conforme aux conventions Go
2. **Ajoutez des tests** pour les nouvelles fonctionnalités
3. **Mettez à jour la documentation** si nécessaire
4. **Exécutez les tests** pour vérifier que rien ne cesse de fonctionner
5. **Validez les modifications** avec des messages clairs

### Règles relatives aux commits

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

## Règles relatives aux pull requests

### Rédiger une bonne description de PR

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
- [ ] Les tests ont été ajoutés ou mis à jour
- [ ] La documentation a été mise à jour
- [ ] Tous les tests réussissent
- [ ] Aucune rupture de compatibilité, ou toute rupture est documentée
- [ ] Les messages de commit sont clairs

## Règles relatives au code

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

La documentation utilise M-Press. Modifiez les fichiers `.md` situés dans  
`docs/mpress/content/`, puis prévisualisez-les et validez-les depuis la racine du dépôt :

```bash
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
mpress dev
mpress build --strict --no-purge-css
mpress check
# Use `python` on Windows if `python3` is unavailable.
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

### Style de la documentation

- Utilisez l’orthographe de l’anglais international
- Commencez par le problème
- Fournissez des exemples fonctionnels
- Incluez une section de dépannage
- Ajoutez des renvois vers le contenu associé

## Communauté

### Obtenir de l’aide

- **Discord :** [Rejoignez notre communauté](https://discord.gg/JDdSxwjhGf)
- **Discussions GitHub :** posez vos questions
- **Tickets GitHub :** signalez les bogues

### Code de conduite

Soyez respectueux, inclusifs et professionnels. Nous sommes tous ici pour créer ensemble d’excellents logiciels. Consultez le [Code de conduite](https://github.com/wailsapp/wails/blob/master/CODE_OF_CONDUCT.md) pour plus de détails.

## Reconnaissance

Les contributeurs sont mentionnés dans :

- Notes de version
- Liste des contributeurs
- Statistiques GitHub

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
Lisez la documentation.

[Parcourir la documentation →](/quick-start/why-wails/)

@end
