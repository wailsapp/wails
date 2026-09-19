---
title: "Programme de mise à jour"
description: "Mise à jour des applications Wails v3 depuis l’application elle-même — fournisseurs interchangeables, vérification cryptographique, remplacement atomique et interface utilisateur par défaut que vous pouvez personnaliser ou remplacer."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

Le programme de mise à jour déploie les mises à jour logicielles depuis l’application sans vous obliger à créer votre propre chaîne de téléchargement, de vérification et de remplacement. Il repose sur `app.Updater`, accepte un ou plusieurs `Provider`s interchangeables (GitHub Releases, keygen.sh, Sparkle AppCast, le protocole ouvert Wails Update Manifest ou votre propre fournisseur), authentifie les téléchargements à l’aide d’une clé publique configurée, remplace de manière sûre le binaire en cours d’exécution et signale chaque transition par le bus d’événements Wails standard.

![Fenêtre par défaut du programme de mise à jour dans l’état Mise à jour prête — icône adaptée à l’état, pastille de version (v1.0.0 → v2.0.1 · 8.8 Mo), notes de version rendues en Markdown comprenant un tableau GFM et une seule action principale.](/assets/updater/default-window-ready.png)

## Démarrage rapide

```go {title="main.go"}
package main

import (
    "context"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

func main() {
    app := application.New(application.Options{Name: "Demo"})

    gh, _ := github.New(github.Config{Repository: "myorg/myapp"})
    if err := app.Updater.Init(updater.Config{
        CurrentVersion: "1.0.0",
        Providers:      []updater.Provider{gh},
    }); err != nil {
        log.Fatal(err)
    }

    if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
        log.Printf("update: %v", err)
    }

    _ = app.Run()
}
```

Cela ouvre la fenêtre de mise à jour du framework, recherche une mise à jour sur GitHub, télécharge l’artefact correspondant à la plateforme, le vérifie, remplace le binaire et attend que l’utilisateur redémarre l’application.

## Le cycle de vie

`app.Updater` est une machine à états qui comprend les états suivants (`updater.State`) :

| État | Situation |
| --- | --- |
| `unconfigured` | Avant l’appel de `Init` |
| `idle` | Après `Init`, avant toute vérification |
| `checking` | Une opération `Check` est en cours |
| `up-to-date` | La dernière réponse du fournisseur a indiqué que l’appelant est à jour |
| `available` | Une nouvelle version a été trouvée, mais son téléchargement n’a pas encore commencé |
| `downloading` | Les octets sont reçus en continu depuis le fournisseur |
| `verifying` | Le téléchargement est terminé ; la signature ou le condensat est en cours de vérification |
| `installing` | Les octets vérifiés sont décompressés et renommés dans le répertoire de préinstallation |
| `ready` | La mise à jour est préinstallée ; appelez `Restart` pour l’appliquer |
| `error` | Une étape précédente a échoué |

Vous pouvez consulter l’état actuel à tout moment avec `app.Updater.State()`. Chaque transition émet également un événement Wails (voir [Événements](#vnements)).

`Restart` attend que l’assistant atteigne `application.New` avant de demander à l’application en cours de quitter. Le délai de démarrage par défaut est de 30 secondes. Si votre application effectue une longue initialisation avant `application.New`, définissez `Config.HelperReadyTimeout` sur une durée plus longue, par exemple `time.Minute`. Zéro sélectionne la valeur par défaut ; les durées négatives sont rejetées. Si le délai de démarrage expire, `Restart` renvoie `updater.ErrHelperNotReady` et laisse l’application en cours ouverte.

La fenêtre par défaut reflète automatiquement l’état actuel. Par exemple, lorsque `Check` ne trouve aucune mise à niveau, l’utilisateur en est informé et ferme la fenêtre avec **Fermer** :

![Fenêtre par défaut du programme de mise à jour dans l’état À jour — coche verte, titre « Vous êtes à jour » et un seul bouton Fermer.](/assets/updater/default-window-up-to-date.png)

## Fournisseurs

Un `Provider` est tout élément qui satisfait cette interface :

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

Quatre implémentations sont incluses dans l’arborescence du projet.

### GitHub Releases — `updater/providers/github`

```go {title="github provider"}
gh, err := github.New(github.Config{
    Repository:    "myorg/myapp",     // your owner/repo (required)
    Token:         "ghp_…",           // optional; raises rate limit + private repos
    Prerelease:    false,             // include pre-releases in latest lookup
    ChecksumAsset: "SHA256SUMS",      // optional sibling asset for digest verification
    BaseURL:       "",                // optional override (e.g. GitHub Enterprise)
    AssetMatcher:  nil,               // optional custom asset-picker; nil uses DefaultAssetMatcher
    HTTPClient:    nil,               // optional client override
})
```

Le sélecteur d’artefact par défaut recherche dans le nom de fichier les sous-chaînes `GOOS` + `GOARCH` et reconnaît les alias courants (`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`). Pour les conventions de nommage personnalisées :

```go
gh, _ := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, a := range assets {
            if strings.Contains(a.Name, "my-naming-convention") &&
               strings.Contains(a.Name, req.Platform) {
                return i
            }
        }
        return -1 // no match
    },
})
```

`ChecksumAsset` est le nom d’un artefact associé de la même version, dont le contenu se compose de lignes `<sha256>  <filename>` (le format produit par `sha256sum` et `shasum -a 256`). Le fournisseur le récupère pendant `Check`, trouve la ligne correspondant à l’artefact sélectionné et renseigne `Release.Verification.Digest` afin que le framework vérifie le téléchargement.

### keygen.sh — `updater/providers/keygen`

```go {title="keygen provider"}
kg, err := keygen.New(keygen.Config{
    Account:    "your-account-slug", // required
    Product:    "product-uuid",      // optional but recommended when account has multiple products
    Package:    "",                  // optional further narrowing
    Channel:    "stable",            // "stable" / "rc" / "beta" / "alpha" / "dev"
    Filetype:   "",                  // optional artifact filetype filter ("dmg", "exe", …)
    Token:      "prod-…",            // product / environment / user / admin token; wins over LicenseKey
    LicenseKey: "",                  // license key auth (used only when Token is empty)
    BaseURL:    "",                  // optional API base override
    HTTPClient: nil,                 // optional client override; redirect-strip wrapper still applied
})
```

Le fournisseur associe automatiquement la somme de contrôle SHA-512 et la signature Ed25519ph propres à chaque artefact de keygen.sh au bloc `Release.Verification` du framework ; aucune configuration supplémentaire n’est requise.

**Formats des jetons :** les jetons keygen.sh comportent un préfixe de rôle (`admi-` / `prod-` / `envi-` / `user-`). L’UUID brut affiché dans le tableau de bord est l’*identifiant* du jeton, et non sa valeur secrète ; celle-ci n’est visible qu’au moment de la création du jeton. Pour plus de détails, consultez la [documentation sur l’authentification](https://keygen.sh/docs/api/authentication/) de keygen.sh.

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

S’intègre sans modification à une infrastructure Sparkle / WinSparkle existante. Lit `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>` et `sparkle:channel` dans le flux.

Les signatures DSA de Sparkle 1 (`sparkle:dsaSignature`) ne sont pas prises en charge ; les projets qui utilisent ce schéma de signature devraient passer à EdDSA (Sparkle 2).

### Manifeste de mise à jour Wails — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

Utilise le protocole ouvert [Wails Update Manifest](/reference/update-manifest/) : un document JSON unique qui décrit la dernière version et ses artefacts propres à chaque plateforme, avec les sommes de contrôle et les signatures intégrées. Le même document fonctionne depuis un hébergeur de fichiers statiques (S3, GitHub Pages ou n’importe quel CDN — publiez un manifeste par canal répertoriant toutes les plateformes) ou depuis un serveur de mise à jour dynamique (le fournisseur envoie `platform`, `arch`, `version` et `channel` à chaque vérification, ce qui permet au serveur de renvoyer exactement un artefact ou de conditionner l’accès à une licence).

Grâce aux paramètres substituables dans l’URL, une seule ligne de configuration suffit pour les structures statiques :

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

Les en-têtes configurés sont envoyés avec chaque requête de manifeste. Les téléchargements d’artefacts ne les réutilisent que sur l’hôte du manifeste lui-même, sans rétrogradation de `https` vers `http`. L’en-tête `Authorization` est supprimé lors de toute redirection interorigine ou entraînant une rétrogradation.

La CLI gère la publication : `wails3 updater manifest` calcule les condensats de vos fichiers de version, les signe et les décrit en une seule commande, puis `wails3 updater verify` vérifie à nouveau le résultat avant son téléversement. Consultez [Publier avec la CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

### Chaîne de repli

`Config.Providers` est ordonné. Le programme de mise à jour le parcourt séquentiellement : le premier fournisseur qui renvoie une version l’emporte ; le premier qui indique que l’application est « à jour » interrompt immédiatement la chaîne (le repli sert lorsque « le fournisseur principal est inaccessible », et non lorsque « les fournisseurs sont en désaccord »). Une erreur fait passer au fournisseur suivant.

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### Écrire votre propre fournisseur

Trois méthodes, environ 150 lignes pour une implémentation classique. Updater se charge de la vérification, de la préparation atomique, du remplacement et de la fenêtre ; le code du fournisseur détermine la prochaine version et transmet les octets en flux :

```go
type CustomProvider struct { /* config */ }

func (p *CustomProvider) Name() string { return "custom" }

func (p *CustomProvider) Check(ctx context.Context, req updater.CheckRequest) (*updater.Release, error) {
    // Hit your update endpoint, decide whether req.CurrentVersion is current,
    // and return either nil (no upgrade) or a *Release with Artifact + optional
    // Verification populated.
    // Errors here drop through to the next provider in Config.Providers.
}

func (p *CustomProvider) Download(ctx context.Context, r *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
    // Stream the artifact's bytes to dst. Call onProgress(written, total) as
    // bytes flow past; the Updater debounces emits to ~10 Hz on the event bus.
}
```

Utilisez les fournisseurs inclus dans l’arborescence du projet comme références : chacun tient dans un seul fichier Go.

## Vérification cryptographique

Les versions sont authentifiées par le vérificateur du framework, qui utilise `Config.PublicKey` comme racine de confiance :

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

Algorithmes pris en charge (`Release.Verification.SignatureAlgo`) :

| Algorithme | Données signées | Remarques |
| --- | --- | --- |
| `ed25519` | L’empreinte SHA-256 de l’artefact | Utilisé par Sparkle EdDSA |
| `ed25519ph` | L’artefact complet via le préhachage d’Ed25519ph (SHA-512 en interne) | Utilisé par keygen.sh |
| `ecdsa-p256` | L’empreinte SHA-256 de l’artefact | Les signatures `r∥s` brutes et DER sont acceptées |

La vérification de l’empreinte seule (`DigestAlgo` : `sha256` / `sha512`) est également possible lorsqu’une version fournit une empreinte, mais aucune signature.

`Config.PublicKey` est l’UNIQUE ancre de confiance utilisée pour vérifier les signatures : la source de la version ne peut en aucun cas lui substituer sa propre clé. Toute version contenant une `Signature` alors qu’aucune `Config.PublicKey` n’est configurée est rejetée par sécurité. Le vérificateur calcule l’empreinte en flux pendant le téléchargement ; même pour des mises à jour de plusieurs Go, la vérification ne nécessite donc aucune lecture supplémentaire du disque.

@note{type="caution" title="Empreinte seule ≠ vérification cryptographique"}
Une version qui ne contient que `Digest` est authentifiée au moyen du protocole TLS du registre et des garanties d’intégrité fournies par le registre lui-même, et non au moyen d’une racine cryptographique que vous contrôlez. Utilisez l’empreinte seule pour détecter l’altération accidentelle des données ; utilisez des signatures pour résister aux falsifications en cas de compromission de la chaîne de publication.

@end

### Génération d’une clé de signature

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

La clé privée est au format PEM PKCS#8 et la clé publique au format PEM PKIX. `Config.PublicKey` accepte directement le fichier `.pub` ; il accepte également la clé brute de 32 octets ou son encodage base64, que `genkey` affiche afin de permettre son intégration directe. Signez les versions avec `wails3 updater manifest -key updater.key ...` ou `wails3 updater sign` ; consultez [Publication avec la CLI wails3](/reference/update-manifest/#publishing-with-the-wails3-cli).

Ou en Go :

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## Formats des artefacts

Les fournisseurs transmettent en flux les octets du fichier que vous publiez ; le framework le décompresse ensuite avant le remplacement :

- **Binaire unique** (par exemple `myapp-linux-amd64`) — utilisé tel quel. Courant sous Linux.
- **`.zip`** — extrait sur place. L’archive doit contenir exactement une entrée de premier niveau, généralement un paquet `.app` macOS ou un binaire unique. Format de paquet recommandé pour macOS.
- **`.tar.gz`** / **`.tgz`** — extrait sur place selon la même règle imposant une seule entrée de premier niveau. Utile pour les distributions Linux qui fournissent une arborescence d’exécution avec le binaire.

Les archives contenant plusieurs entrées de premier niveau sont rejetées : le framework remplace une seule cible sur le disque ; l’instruction « remplacer la cible par cette archive » est donc ambiguë lorsque l’archive contient plusieurs éléments. Les formats `.dmg` et `.pkg` (macOS), ainsi que `.msi` (Windows), ne sont pas pris en charge dans la v1 : distribuez plutôt une archive `.zip` du paquet. L’extraction protège contre les attaques zip-slip, rejette les liens symboliques qui sortent de la racine de l’archive et limite la taille totale décompressée à 2 Gio ainsi que le nombre d’entrées à 50 000.

## La fenêtre par défaut

`app.Updater.CheckAndInstall(ctx)` ouvre une fenêtre de 520 × 540 gérée par le framework et comprenant :

- Une icône principale déterminée par l’état (↓ bleu si la mise à jour est disponible ou en cours de téléchargement, ✓ vert si elle est prête ou si l’application est à jour, ! rouge en cas d’erreur)
- Une pastille de version : `v1.0.0 → v2.0.1 · 8.8 MB`
- Un panneau de notes de version à défilement avec **rendu Markdown** (paragraphes, gras/italique, listes, tableaux GFM, code en ligne, blocs de code délimités, h1 à h3, liens)
- Une seule action principale par état (Installer / Redémarrer et appliquer / Réessayer)
- Des actions secondaires présentées sous forme de boutons fantômes (Ignorer cette version / Me le rappeler plus tard)
- Un mode sombre ou clair via `prefers-color-scheme`
- Un effet de scintillement pour la progression indéterminée lorsque la taille totale est inconnue

Elle écoute les événements `updater:*` sur le bus d’événements Wails et renvoie les actions `updater:user:*` à Go.

### Thème au moyen de variables CSS

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

La feuille de style par défaut expose les variables suivantes ; vous pouvez redéfinir celles de votre choix :

| Variable | Valeur par défaut (thème clair) | Valeur par défaut (thème sombre) |
| --- | --- | --- |
| `--bg` | `#f8f8fa` | `#1a1a1c` |
| `--surface` | `#ffffff` | `#232326` |
| `--surface-2` | `#f0f0f3` | `#2c2c30` |
| `--fg` | `#1d1d1f` | `#f5f5f7` |
| `--fg-dim` | `#6b6b73` | `#b0b0b8` |
| `--fg-faint` | `#99999f` | `#7a7a82` |
| `--border` | `#d6d6dc` | `#3a3a3e` |
| `--accent` | `#0a84ff` | `#0a84ff` |
| `--accent-fg` | `#ffffff` | — |
| `--success` | `#34c759` | — |
| `--error` | `#ff3b30` | — |
| `--radius` | `10px` | — |
| `--font` | pile de polices système | — |

### Remplacer le modèle

Fournissez votre propre HTML ; il lui suffit d’écouter les événements `wails:updater:*` et d’émettre les actions `wails:updater:user:*` :

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

Les fenêtres InitialHTML sont chargées sans origine provenant du serveur de ressources et ne peuvent donc pas récupérer `/wails/runtime.js` dynamiquement. Deux méthodes permettent à une telle fenêtre de communiquer avec l’hôte :

1. **Écrivez simplement du HTML.** Le framework injecte automatiquement une couche de compatibilité `window.wails.Events` minimale dans toute fenêtre ouverte avec `WebviewWindowOptions.AllowSimpleEventEmit = true` et `HTML` définis, ce qui correspond exactement aux chemins intégré et BYO du programme de mise à jour. Aucune étape de compilation n’est requise. L’exemple ci-dessous utilise cette méthode.
2. **Regroupez `@wailsio/runtime` avec l’outil de votre choix** (Vite, esbuild, Rollup), puis importez-le dans votre HTML personnalisé lors de la compilation. `Events.On` fonctionne immédiatement, car il s’exécute entièrement côté client ; `Events.Emit` passe par le transport fetch du runtime, que l’origine nulle rend inopérant. Installez donc un petit transport postMessage au moyen du hook [`setTransport`](https://wails.io/wails/runtime.js) du runtime, qui effectue le routage via `window._wails.invoke("wails:event:emit:<name>")`. L’injection du framework ne fait rien si `window.wails.Events` est déjà dans la portée, afin que les deux approches n’entrent pas en conflit.

Dans les deux cas, le JavaScript de votre HTML personnalisé appelle la même API `Events.On` / `Events.Emit` :

```html
<script>
const { On, Emit } = window.wails.Events;

On("wails:updater:update-available", (e) => {
    const rel = e.data ?? e;
    document.getElementById("ver").textContent = rel.version;
});

document.getElementById("install").addEventListener("click",
    () => Emit("wails:updater:user:install"));

// Ask the host to replay the current state so we paint correctly on (re)open.
Emit("wails:updater:window:ready");
</script>
```

La couche de compatibilité expose le sous-ensemble du runtime moderne nécessaire aux événements à nom simple : `Events.On(name, cb)` renvoie une fonction de désabonnement et `Events.Emit(nameOrEventObject)` effectue le routage vers l’hôte par le chemin postMessage `wails:event:emit:` soumis à autorisation. Elle est installée une seule fois au chargement de la page, avant l’exécution de vos scripts en ligne.

Si vous *souhaitez* remplacer la couche de compatibilité (ou si vous chargez le runtime complet d’une autre manière), définissez `window.wails.Events` avant l’exécution de la première balise `<script>` de la page ; l’injection sera alors ignorée.

### Habillage de la fenêtre

Remplacez les options de la fenêtre (taille, sans cadre, toujours au premier plan) sans modifier le HTML :

```go
Window: &updater.BuiltinWindow{
    Options: updater.WindowOptions{
        Title:         "My App Updater",
        Width:         640,
        Height:        480,
        Frameless:     true,
        AlwaysOnTop:   true,
        DisableResize: false,
    },
},
```

### Utiliser votre propre fenêtre

Exécutez le flux de mise à jour dans une fenêtre `*application.WebviewWindow` que vous créez vous-même. Le programme de mise à jour appelle `Show()` / `Close()` / `EmitEvent()` sur votre fenêtre ; votre HTML détermine le rendu :

![Fenêtre personnalisée du programme de mise à jour, avec un arrière-plan en dégradé rose-orange, une seule carte blanche aux coins arrondis, une typographie personnalisée et les mêmes événements de mise à jour pilotant l’état visible. Elle montre à quel point l’interface utilisateur par défaut peut être intégralement remplacée.](/assets/updater/byo-custom-window.png)

```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My Updater",
    Width:                520, Height: 460,
    HTML:                 myCustomHTML,
    AllowSimpleEventEmit: true,           // see security note below
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Votre HTML utilise `window.wails.Events.On` / `Events.Emit` comme le modèle intégré : la couche de compatibilité injectée automatiquement par le framework est intégrée à toute fenêtre comportant `AllowSimpleEventEmit: true`, que cette fenêtre appartienne au framework ou à votre application. Consultez [Remplacer le modèle](#remplacer-le-modle) pour connaître l’API.

@note{type="caution" title="`AllowSimpleEventEmit` est requis pour les fenêtres BYO du programme de mise à jour"}
Pour des raisons de sécurité, le framework utilise ce champ pour contrôler l’accès au raccourci postMessage `wails:event:emit:` : une fenêtre dans laquelle il n’est pas défini ne peut pas générer d’événements personnalisés côté hôte. La couche de compatibilité du programme de mise à jour destinée au HTML personnalisé émet les événements `updater:user:*` par ce raccourci. Si ce champ est omis dans une fenêtre BYO, chaque clic sur un bouton est donc silencieusement ignoré : l’utilisateur clique sur Installer et rien ne se passe.

Laissez `AllowSimpleEventEmit` **désactivé** pour toute fenêtre qui charge du HTML que vous ne contrôlez pas entièrement (URL distantes, contenu fourni par l’utilisateur). Lorsqu’il est activé, tout code JavaScript de la page, y compris les points d’injection XSS, peut déclencher n’importe quel gestionnaire `app.Event.On(name, …)`. Le raccourci ne transmet que des noms simples (sans charge utile) et ne peut pas atteindre le chemin binding/Call, mais il peut néanmoins déclencher des gestionnaires privilégiés d’événements personnalisés dans votre code Go si ceux-ci agissent uniquement en fonction du nom de l’événement.

Ce champ est défini en interne pour la fenêtre *intégrée* du programme de mise à jour du framework ; seuls les appelants BYO doivent penser à le définir.

@end

### Sans interface graphique

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

Aucune fenêtre n’est jamais ouverte. Abonnez votre propre interface utilisateur (ou votre fenêtre principale existante) aux événements `updater:*` et appelez `app.Updater.CheckAndInstall(ctx)` depuis un gestionnaire de bouton. Cette méthode est utile pour les vérifications périodiques en arrière-plan qui ne doivent apparaître que lorsqu’une mise à jour est trouvée, ou pour les applications qui intègrent le flux de mise à jour à un panneau de paramètres personnalisé.

## Événements

Go et JavaScript s’abonnent tous deux par le bus d’événements standard de Wails. **Ne saisissez pas manuellement les chaînes de transport** : utilisez les constantes exportées par le package updater (Go) ou le package runtime (JS). Les deux couches partagent le même ensemble de noms, dont la synchronisation est garantie par un test de non-régression.

### Depuis Go

Les constantes se trouvent dans `github.com/wailsapp/wails/v3/pkg/updater`. Abonnez-vous au moyen de `app.Event.On(name, fn)` ; la fonction de rappel reçoit un `*application.CustomEvent` dont le champ `Data` contient la charge utile typée indiquée dans la [référence des événements](#rfrence-des-vnements). Effectuez une assertion de type, pas un décodage JSON :

```go {title="main.go"}
import (
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
)

// …

app.Event.On(updater.EventUpdateAvailable, func(e *application.CustomEvent) {
    rel, ok := e.Data.(*updater.Release)
    if !ok { return }
    log.Printf("update found: %s", rel.Version)
})

app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
    p, ok := e.Data.(updater.Progress)
    if !ok { return }
    log.Printf("%d / %d bytes (%.0f KB/s)", p.Written, p.Total, p.Rate/1024)
})

app.Event.On(updater.EventError, func(e *application.CustomEvent) {
    info, ok := e.Data.(updater.ErrorInfo)
    if !ok { return }
    log.Printf("update failed during %s: %s", info.Stage, info.Message)
})
```

Toutes les constantes Go disponibles :

| Constante | Chaîne de transport |
| --- | --- |
| `updater.EventCheckStarted` | `wails:updater:check-started` |
| `updater.EventUpdateAvailable` | `wails:updater:update-available` |
| `updater.EventNoUpdate` | `wails:updater:no-update` |
| `updater.EventDownloadStarted` | `wails:updater:download-started` |
| `updater.EventDownloadProgress` | `wails:updater:download-progress` |
| `updater.EventDownloadComplete` | `wails:updater:download-complete` |
| `updater.EventVerifying` | `wails:updater:verifying` |
| `updater.EventInstalling` | `wails:updater:installing` |
| `updater.EventUpdateReady` | `wails:updater:update-ready` |
| `updater.EventError` | `wails:updater:error` |
| `updater.EventMeta` | `wails:updater:meta` |
| `updater.EventWindowReady` | `wails:updater:window:ready` |
| `updater.EventUserInstall` | `wails:updater:user:install` |
| `updater.EventUserSkip` | `wails:updater:user:skip` |
| `updater.EventUserRemind` | `wails:updater:user:remind` |
| `updater.EventUserCancel` | `wails:updater:user:cancel` |
| `updater.EventUserRestart` | `wails:updater:user:restart` |

### Depuis JavaScript

Les constantes se trouvent sous `Updater.Events` dans `@wailsio/runtime`. Elles portent les mêmes noms qu’en Go et sont organisées par sous-espace de noms (`User.*`, `Window.*`) afin de faciliter leur découverte par autocomplétion :

```js
import { Events, Updater } from "@wailsio/runtime";

Events.On(Updater.Events.UpdateAvailable, (e) => {
    console.log("update found:", e.data.version);
});

Events.On(Updater.Events.DownloadProgress, (e) => {
    const p = e.data;
    console.log(`${p.written} / ${p.total} bytes (${(p.rate/1024).toFixed(1)} KB/s)`);
});

Events.On(Updater.Events.Error, (e) => {
    const info = e.data;
    console.error(`update failed during ${info.stage}: ${info.message}`);
});
```

Les événements d’action utilisateur que votre code HTML personnalisé renvoie *à l’hôte* se trouvent sous `Updater.Events.User` :

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### Référence des événements

Côté abonnement (hôte → page) :

| Constante (Go) | Constante (JS) | Charge utile | Quand |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | aucune | Avant chaque aller-retour de `Check` |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` a trouvé une version plus récente |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | aucune | `Check` a confirmé que l’application est à jour |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | Le transfert des octets commence |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | Environ 10 Hz pendant le téléchargement |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | Tous les octets ont été écrits, avant la vérification |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | La vérification de la signature ou du condensat commence |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | La décompression et la préparation commencent |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | Redémarrage en attente |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | Échec à une étape quelconque |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | Une fois par session, avant la relecture de l’instantané |

Côté page (page → hôte) — votre code s’y abonne si vous écrivez un modèle personnalisé :

| Constante (Go) | Constante (JS) | Quand |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | La fenêtre a fini de se charger ; l’hôte retransmet l’état actuel |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | Action principale dans l’état `available` |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | Action principale dans l’état `ready` |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | « Ignorer cette version » |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | « Me le rappeler plus tard » |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | Bouton de fermeture |

## Référence de l’API

### `updater.Config`

| Champ | Type | Remarques |
| --- | --- | --- |
| `CurrentVersion` | `string` | **Obligatoire.** Même chaîne que celle utilisée pour étiqueter les versions (sans préfixe `v`) |
| `Providers` | `[]updater.Provider` | **Obligatoire.** Chaîne de repli ordonnée |
| `PublicKey` | `[]byte` | PEM ou octets bruts. Facultatif, mais sans cette valeur, les versions signées sont rejetées par sécurité |
| `CheckInterval` | `time.Duration` | Une valeur non nulle démarre une boucle d’interrogation en arrière-plan qui appelle `CheckAndInstall` |
| `Platform` | `string` | Remplace `runtime.GOOS` pour la sélection de l’artefact |
| `Arch` | `string` | Remplace `runtime.GOARCH` pour la sélection de l’artefact |
| `Channel` | `string` | Actuellement fourni à titre informatif ; filtrage par canal propre au fournisseur |
| `Window` | `updater.WindowOption` | `nil` (valeurs par défaut intégrées), `&BuiltinWindow{…}`, `BYOWindow(handle)` ou `WindowNone` |

### Méthodes de `*updater.Updater`

| Signature | Fonction |
| --- | --- |
| `Init(cfg Config) error` | Configure l’outil de mise à jour. Renvoie `ErrAlreadyConfigured` au deuxième appel |
| `State() State` | Phase actuelle du cycle de vie |
| `CurrentVersion() string` | Version transmise à `Init` |
| `Check(ctx) (*Release, error)` | Parcourt la chaîne de fournisseurs. `(rel, nil)` = trouvée, `(nil, nil)` = à jour, `(nil, err)` = tous ont échoué |
| `DownloadAndInstall(ctx) error` | Télécharge en continu, vérifie, extrait (s’il s’agit d’une archive), puis prépare. Nécessite un appel préalable à `Check` |
| `CheckAndInstall(ctx) error` | Raccourci : ouvre la fenêtre, appelle `Check`, puis `DownloadAndInstall` si une mise à jour est trouvée |
| `Restart(ctx) error` | Lance le processus auxiliaire, appelle `Host.Quit`, puis quitte ; le processus auxiliaire effectue le remplacement et relance l’application |
| `DownloadedPath() string` | Emplacement sur le disque de la mise à jour préparée, ou `""` s’il n’y en a aucune |
| `SkipVersion(v string)` | Enregistre `v` comme version ignorée ; les appels ultérieurs à `Check` la considèrent comme étant à jour |
| `SkippedVersion() string` | Lit la version actuellement ignorée |
| `StopPeriodicCheck()` | Annule le minuteur démarré par `Config.CheckInterval` et attend la fin de la boucle |

### Erreurs

| Valeur sentinelle | Renvoyée par |
| --- | --- |
| `ErrAlreadyConfigured` | `Init` après la première réussite |
| `ErrNotConfigured` | Toute opération avant `Init` |
| `ErrNoPendingRelease` | `DownloadAndInstall` sans appel préalable à `Check` |
| `ErrDownloadInProgress` | `DownloadAndInstall` appelé alors qu’un autre est en cours |
| `ErrNotReady` | `Restart` sans mise à jour préparée |

## Fonctionnement du remplacement

`Restart` réexécute le binaire actuel après avoir défini des variables d’environnement sentinelles. Au démarrage, `application.New` les détecte et bascule en mode auxiliaire :

1. Le processus auxiliaire attend jusqu’à 30 s que le PID du processus parent se termine (`platformIsAlive` effectue des interrogations avec `syscall.OpenProcess` + `GetExitCodeProcess` sous Windows, et `os.FindProcess` + `proc.Signal(syscall.Signal(0))` sous Unix).
2. Le processus auxiliaire sauvegarde la cible (copie simple pour les fichiers, copie récursive pour les répertoires de paquet macOS `.app`).
3. Le processus auxiliaire remplace la cible par l’artefact préparé et réessaie jusqu’à 20 fois, avec un délai de 500 ms entre les tentatives :
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. Les descripteurs de fichiers ouverts sur l’ancien inode restent valides.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. Windows autorise le renommage des fichiers dont l’image est encore mappée, mais pas leur suppression ; lors de la mise à jour suivante, le processus auxiliaire supprime tous les fichiers frères `.old.*` restants dont le mappage noyau qui les maintenait référencés a été libéré.

4. Le processus auxiliaire rétablit le mode d’exécution d’origine sur le nouveau binaire (le fichier téléchargé a été créé avec l’umask par défaut, qui supprime `+x` sous Unix ; sous Windows, cette opération est sans effet).
5. Le processus auxiliaire supprime les variables d’environnement du mode auxiliaire et relance le binaire désormais remplacé.
6. Le processus auxiliaire se termine.

Si le lancement échoue, le processus auxiliaire restaure la sauvegarde. Si le processus parent ne se termine pas dans les 30 s, le processus auxiliaire abandonne avant de modifier la cible, afin que l’utilisateur conserve une application fonctionnelle même si une boîte de dialogue de fermeture bloque `Quit`.

Pour les paquets macOS `.app` distribués sous forme de `.zip` — le format recommandé —, l’archive est décompressée entre la vérification et l’état prêt, afin que le processus auxiliaire dispose d’un véritable répertoire à mettre en place par remplacement.

## Vérification périodique

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

Lorsque `CheckInterval > 0`, une goroutine en arrière-plan appelle `CheckAndInstall` selon l’intervalle configuré. Les déclenchements qui surviennent alors qu’un autre flux est déjà en cours — vérification, téléchargement, validation ou installation — sont ignorés : les machines à états concurrentes ne sont pas prises en charge.

Pour effectuer une vérification silencieuse en arrière-plan et n’afficher quelque chose que lorsqu’une mise à jour est trouvée, définissez `Window: updater.WindowNone` et réagissez à `EventUpdateAvailable` depuis votre propre interface utilisateur.

## Ignorer et rappeler

Le bouton « Ignorer cette version » de la fenêtre par défaut enregistre la version disponible via `SkipVersion(rel.Version)`. Les appels ultérieurs à `Check` trouvent la même version et considèrent l’application comme à jour, jusqu’à ce que l’utilisateur mette à jour `CurrentVersion`, ce qui se produit automatiquement après la réussite de `Restart`. « Me le rappeler plus tard » ferme simplement la fenêtre sans rien enregistrer.

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## Liste de contrôle pour la distribution

Avant de publier une version que le programme de mise à jour installera :

1. **Choisissez le bon format d’archive.** macOS : `.zip` du paquet `.app`. Linux : un seul binaire ou `.tar.gz`. Windows : un seul `.exe` ou `.zip`. `.dmg` / `.msi` / `.pkg` ne sont pas pris en charge.
2. **Signez l’artefact** avec la clé privée correspondant à `Config.PublicKey`. Pour les flux publiés par un fournisseur (keygen.sh, AppCast), suivez la procédure de signature propre à chaque fournisseur. Pour GitHub Releases avec `ChecksumAsset`, générez un fichier `SHA256SUMS` avec `sha256sum` / `shasum -a 256`.
3. **Faites correspondre la chaîne de version.** `Config.CurrentVersion` et l’étiquette de version de la publication doivent correspondre exactement (par exemple, `1.0.0` ↔ étiquette `v1.0.0` ; le préfixe `v` est supprimé côté fournisseur).
4. **Testez le remplacement sur la plateforme cible** au moins une fois avant la distribution : la signature du code, la notarisation et la gestion par Gatekeeper sont propres à chaque plateforme et ne sont pas prises en charge par le programme de mise à jour lui-même.

## Dépannage

**« signature requires a public key but none configured »** — la publication contient un champ `Signature`, mais `Config.PublicKey` est vide. Définissez la clé publique ou modifiez votre pipeline de publication afin qu’il n’inclue pas de signature.

**« digest mismatch »** — les octets téléchargés ne correspondent pas à ceux annoncés par le fournisseur. Il s’agit généralement d’un téléchargement partiel dû à un incident réseau ou d’un artefact corrompu. Relancer l’opération résout souvent le problème.

**La fenêtre s’ouvre, mais disparaît immédiatement, sans Markdown ni progression** — votre HTML personnalisé n’a pas appelé `wails:runtime:ready`. Consultez le shim [Remplacer le modèle](#remplacer-le-modle).

**La mise à jour sous Windows ne se termine jamais ; le journal du processus auxiliaire indique « remove old (attempt N): Access is denied »** — cela ne se produit que dans les versions de cette PR antérieures à `de764fb` ; l’implémentation actuelle renomme l’ancien fichier pour le mettre de côté et ne rencontre pas ce problème. Effectuez la mise à niveau.

**Gatekeeper sous macOS bloque le binaire remplacé** — la signature du code doit être préservée de bout en bout. Signez le `.app` d’origine *et* signez de nouveau le binaire relancé si votre pipeline de build modifie les autorisations lors de la mise à jour.

## Voir aussi

- Exemple exécutable : [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)
- Dépôt de démonstration pour les tests : [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- Tutoriel : [Permettre à une application Wails de se mettre à jour elle-même](/tutorials/04-self-update-a-wails-app/)
