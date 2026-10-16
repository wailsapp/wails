---
title: "Débogage"
description: "Analyse des problèmes et profilage des performances de l’application"
slug: "guides/dev/debugging"
sourcePath: "guides/dev/debugging.md"
---

Ce guide présente différents outils permettant d’examiner et d’analyser les éventuels problèmes de performances de votre application Wails à l’aide de

- [`runtime/trace`](https://pkg.go.dev/runtime/trace) pour créer des graphiques de performances et les examiner dans un navigateur

## Création de traces de performances

@steps
### Préparer votre application au traçage
Près du point d’entrée de votre programme, vérifiez que du code semblable au suivant est présent

```go
// Create the file to store our trace data within
traceFile, err := os.Create("trace.out")
if err != nil {
  log.Fatalf("trace.out could not be created: %v", err)
}

// Start the trace
if err := trace.Start(traceFile); err != nil {
  _ = traceFile.Close()
  log.Fatalf("trace.start could not start: %v", err)
}

// Trace cleanup on exit. Alternatively,
defer func() {
  trace.Stop()
  _ = traceFile.Close()
}()

...Start your wails app here...
```

La sortie sera enregistrée dans `trace.out`, sous votre chemin de travail, avec des métriques couvrant toute la durée d’exécution de votre application. La trace par défaut est assez brute ; il est vivement conseillé de consulter davantage de documentation sur le traçage afin d’ajouter plus d’informations contextuelles et de limiter ce que vous enregistrez réellement. Par exemple : `WithRegion, NewTask, Log`

### Exécuter votre application pour créer la trace
Pendant l’exécution de votre application, une trace continue sera produite jusqu’à sa fermeture. Effectuez donc quelques actions dans votre application, puis fermez-la lorsque vous avez terminé

### Installer les outils auxiliaires de visualisation
Certaines vues nécessitent [Graphviz](https://graphviz.org/). Pour vérifier qu’il est installé, exécutez `dot -V`

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Arch
sudo pacman -S graphviz
```

### Visualiser votre trace
Une fois les données de trace disponibles, vous pouvez démarrer l’interface web

```bash
go tool trace trace.out
```

Cela devrait ouvrir votre navigateur par défaut — il est recommandé d’utiliser un navigateur basé sur Chrome — sur la page d’accueil du visualiseur d’événements de trace.

Pour une première utilisation, l’écran `Syscall profile` est probablement le plus utile. Il présente en détail ce que fait exactement le programme, ainsi que les durées correspondantes

@end
