---
title: "Déboguer les plantages"
description: "Collecter des rapports de diagnostic et des minidumps depuis une application Wails"
slug: "guides/debugging-crashes"
sourcePath: "guides/debugging-crashes.md"
---

Lorsque quelque chose se passe mal en production, les journaux suffisent rarement. Ce guide présente le package `github.com/wailsapp/wails/v3/pkg/debug`, qui produit des rapports de diagnostic structurés et, sous Windows, des minidumps complets pouvant être ouverts dans WinDbg ou Visual Studio.

Le package est volontairement petit et possède deux points d’entrée :

- `debug.Dump(...)` — écrire un minidump Windows du processus courant.
- `debug.Report(...)` — collecter un instantané de diagnostic détaillé, avec un minidump facultatif.

Les deux s’exécutent avec le jeton de l’utilisateur appelant. Aucune élévation, aucun `SeDebugPrivilege`, aucun `OpenProcess` : le dump du processus courant utilise le pseudo-handle `GetCurrentProcess`, qui contourne les vérifications d’ACL pour l’accès à son propre processus.

## Référence rapide

```go
import "github.com/wailsapp/wails/v3/pkg/debug"

// 1. Just a minidump (Windows only).
path, err := debug.Dump()

// 2. Dump to a specific path.
path, err := debug.Dump(debug.WithPath("C:\\crashes\\my.dmp"))

// 3. Full-memory minidump (large, gigabytes).
path, err := debug.Dump(debug.WithFullMemory())

// 4. Full diagnostic report, no minidump.
r, err := debug.Report()

// 5. Report + minidump.
r, err := debug.Report(debug.WithDump())

// 6. Report + minidump at a specific path, full-memory.
r, err := debug.Report(
    debug.WithDumpPath("C:\\crashes\\my.dmp"),
    debug.WithDumpFullMemory(),
)
```

`debug.Report` renvoie toujours un `*CrashReport`, même en cas d’échec partiel : l’appelant obtient tout ce qui a été collecté avant l’erreur, et le rapport n’est donc jamais `nil` lorsque `err != nil`.

## La structure `CrashReport`

```go
type CrashReport struct {
    Timestamp   time.Time          // when the report was generated
    System      SystemInfo         // OS, hardware, CPU, GPU from doctor
    Build       BuildInfo          // Go version, buildmode, compiler, CGO flag
    Crash       *CrashInfo         // process, memory, modules, env vars
    Diagnostics []DiagnosticResult // doctor's health-check results
    DumpPath    string             // set only if WithDump() was passed
}
```

Sérialisez avec `encoding/json` pour les services de signalement de plantages, ou parcourez directement les champs si seul un sous-ensemble vous intéresse.

## Intégration avec `PanicHandler`

Le schéma le plus utile consiste à collecter un rapport, et éventuellement un minidump, dès que Wails intercepte une panique, puis à transmettre les données combinées à votre propre chaîne de signalement :

```go
app := application.New(application.Options{
    PanicHandler: func(pd *application.PanicDetails) {
        // Always: capture a diagnostic snapshot plus a minidump.
        report, reportErr := debug.Report(debug.WithDump())
        if reportErr != nil {
            log.Printf("debug.Report: %v", reportErr)
        }

        // Now you have:
        //   pd.Error, pd.StackTrace, pd.FullStackTrace  — from wails
        //   report.DumpPath                              — minidump (Windows)
        //   report.System / report.Crash / report.Build — context
        //
        // Ship it off (Sentry, S3, support ticket, local log...).
        mycrashservice.Upload(pd, report)
    },
})
```

Wails **n’appelle pas** automatiquement `debug.Report` : le signalement de plantages implique souvent le consentement de l’utilisateur ou la suppression de données personnelles, donc cette décision relève de votre gestionnaire. Consultez le [guide de gestion des paniques](/guides/panic-handling/) pour savoir quand Wails appelle `PanicHandler` et comment couvrir aussi les goroutines créées par votre code.

## Débogage manuel

`debug.Dump` et `debug.Report` sont aussi utiles en dehors des paniques :

- Détection de blocages : lorsqu’un mécanisme de surveillance se déclenche, créez un dump du processus afin de l’ouvrir dans WinDbg et de voir précisément ce qui bloquait chaque thread.
- Rapports d’état anormal : par exemple, un nombre de goroutines qui augmente sans limite ; capturez un rapport, laissez l’application continuer et analysez-le plus tard.
- Dossiers d’assistance à la demande : reliez une entrée de menu « Signaler un bug » à `debug.Report(debug.WithDump())` et joignez le résultat au ticket d’assistance de l’utilisateur.

## Les minidumps Windows en pratique

Par défaut, les minidumps utilisent un ensemble d’indicateurs détaillé mais compact :

- `MiniDumpNormal`
- `MiniDumpWithThreadInfo`
- `MiniDumpWithHandleData`
- `MiniDumpWithUnloadedModules`

Cela produit des dumps d’environ 5 à 50 Mo qui conservent toutes les piles des threads, les tables de handles et les listes de modules, de quoi reconstituer où se trouvait chaque goroutine au moment du dump.

`debug.WithFullMemory()` (ou `debug.WithDumpFullMemory()` sur `Report`) ajoute `MiniDumpWithFullMemory`, qui capture tout l’espace d’adressage du processus. C’est utile pour comprendre pourquoi un pointeur précis a été corrompu, mais la taille des fichiers atteint des centaines de Mo ou des Go. Réservez cette option aux analyses approfondies.

Ouvrez le fichier `.dmp` obtenu dans WinDbg avec `windbg -z C:\path\to\your.dmp`, ou avec Fichier → Ouvrir un dump de plantage dans Visual Studio. Les symboles sont généralement supprimés des builds de production ; associez le dump au `.pdb` du même build, ou au binaire Go lui-même s’il contient les informations de débogage, pour obtenir des traces de pile utiles.

## Plateformes autres que Windows

`debug.Dump` renvoie une erreur « non implémenté » sur les autres plateformes. Linux et macOS disposent d’autres outils d’analyse après incident, comme les core dumps avec `GOTRACEBACK=crash`, `rr` ou des rapporteurs de plantages propres à chaque plateforme, qui ne correspondent pas directement à l’API `MiniDumpWriteDump`. `debug.Report` fonctionne toujours : vous obtenez tout sauf un fichier de dump.

L’erreur de `Dump` est conçue pour être vérifiable afin de permettre un fonctionnement dégradé maîtrisé :

```go
path, err := debug.Dump()
if err != nil {
    log.Printf("minidump unavailable on %s: %v", runtime.GOOS, err)
    // Continue with just debug.Report() or pd.StackTrace from wails.
}
```

## Remarques de sécurité

- **Aucune élévation n’est nécessaire.** Le package évite volontairement `SeDebugPrivilege` et l’approche `OpenProcess` / PID distant utilisée par les outils de collecte d’identifiants. Il ne crée de dump que du processus appelant.
- **Les dumps contiennent tout ce qui se trouve en mémoire.** Cela inclut les identifiants en cours d’utilisation, jetons d’authentification, données utilisateur et clés de chiffrement. Traitez un dump comme une image mémoire complète : protégez-le pendant son transfert, nettoyez les données stockées et envisagez leur suppression avant tout envoi à des tiers.
- **Le chemin par défaut est `os.TempDir()`.** Sous Windows, il correspond à `%LOCALAPPDATA%\Temp` : propre à l’utilisateur et non lisible par tous, mais persistant. Déplacez les dumps vers un emplacement sécurisé si vous les conservez.
