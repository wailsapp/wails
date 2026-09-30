---
title: "Application Wails à mise à jour automatique"
description: "Créez une application Wails v3 qui se met à jour elle-même depuis les versions GitHub, de `wails3 init` à la vérification des versions signées et au remplacement en mode assistant."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

Dans ce tutoriel, vous ajouterez un programme de mise à jour intégré à une nouvelle application Wails v3. À la fin, l’application pourra :

- Rechercher à la demande les nouvelles versions publiées sur GitHub (et, facultativement, à intervalles réguliers).
- Télécharger la ressource adaptée au système d’exploitation et à l’architecture en cours d’utilisation.
- Vérifier une empreinte SHA-256 (et, facultativement, une signature Ed25519) à partir des octets téléchargés.
- Afficher les notes de version dans la fenêtre de mise à jour par défaut du framework.
- Remplacer le binaire en cours d’exécution et relancer l’application, sans distribuer d’exécutable assistant distinct.

Nous utiliserons les **versions GitHub** comme source des mises à jour, car elles sont gratuites et ne nécessitent aucune infrastructure. Les mêmes méthodes fonctionnent avec [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) et [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast). Une fois ce tutoriel terminé, consultez le [guide du programme de mise à jour](/guides/updater/).

@note{type="tip" title="Prérequis"}
- Go 1.25 ou version ultérieure
- CLI `wails3` installée (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- Un dépôt GitHub dans lequel vous pouvez publier des versions
- Il est utile, mais pas indispensable, de connaître le [tutoriel du service de codes QR](/tutorials/01-creating-a-service/)

@end

<br/>

@steps
### Commencer avec une nouvelle application Wails
Générez la structure d’un nouveau projet avec le modèle vanilla :

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

Vous devez maintenant disposer d’un répertoire contenant `main.go`, `frontend/` et un `Taskfile.yml`. Vérifiez que le projet se compile et démarre :

```bash
wails3 task dev
```

Une fenêtre Wails vide doit s’ouvrir. Fermez l’application, puis continuez.

### Ajouter l’importation du programme de mise à jour
Ouvrez `main.go` et ajoutez les deux packages du programme de mise à jour aux importations :

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

Ces importations ajoutent le programme de mise à jour lui-même et le fournisseur de versions GitHub.

### Configurer le programme de mise à jour
`app.Updater` est déjà intégré à chaque `*application.App` : il vous suffit d’appeler `Init` :

```go {title="main.go"}
const currentVersion = "1.0.0"

gh, err := github.New(github.Config{
    Repository:    "yourorg/your-repo",   // ← change this
    ChecksumAsset: "SHA256SUMS",          // sibling file with sha256 digests
})
if err != nil {
    log.Fatalf("github.New: %v", err)
}

if err := app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
}); err != nil {
    log.Fatalf("Updater.Init: %v", err)
}
```

Placez ce code après `application.New` et avant `app.Run()`.

@note{type="note" title="Format de la chaîne de version"}
Transmettez la même version que celle utilisée pour étiqueter les versions publiées, **sans** le `v` initial. De son côté, le fournisseur retire `v` des noms d’étiquettes. `1.0.0` ici ↔ `v1.0.0` sur GitHub.

@end

### Ajouter un élément de menu qui déclenche la mise à jour
Dans le même `main.go`, ajoutez une entrée de menu « Rechercher des mises à jour… » :

```go {title="main.go"}
menu := app.Menu.New()
app.Menu.SetApplicationMenu(menu)
appMenu := menu.AddSubmenu("App")
appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
    go func() {
        if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
            app.Logger.Error("update", "error", err)
        }
    }()
})
```

`CheckAndInstall` ouvre la fenêtre de mise à jour du framework, exécute `Check` puis, si une version est trouvée, exécute automatiquement `DownloadAndInstall`. Lorsqu’aucune nouveauté n’est disponible, la fenêtre reste ouverte dans l’état « À jour » ; l’utilisateur la ferme avec le bouton **Fermer**.

@note{type="caution" title="Exécuter l’opération dans une goroutine"}
`CheckAndInstall` reste bloqué jusqu’à la fin de la vérification et de l’installation. L’appeler directement lors du clic sur le menu bloquerait le thread de l’interface utilisateur. Encapsulez l’appel dans `go func()`.

@end

### Effectuer un premier essai sans version publiée
```bash
wails3 task dev
```

Cliquez sur **Application → Rechercher des mises à jour…**. La fenêtre de mise à jour devrait s’ouvrir brièvement, interroger l’API GitHub, ne trouver aucune version plus récente que `1.0.0`, puis afficher l’état **À jour** avec une coche verte ✓.

Si une erreur se produit à ce stade, elle correspond généralement à l’un des cas suivants :

| Symptôme | Solution |
| --- | --- |
| `404 Not Found` | Le champ `Repository` est incorrect : il doit être défini sur `owner/repo` |
| `403 rate-limited` | Ajoutez `Token: "ghp_…"` à github.Config (utilisez un jeton d’accès personnel doté de la portée `public_repo`) |
| Erreurs réseau | Vérifiez que l’application en cours d’exécution peut accéder à `api.github.com` |

### Publier une version de test
Faites passer `currentVersion` à `1.0.0` dans `main.go` (ou laissez la valeur actuelle). Compilez l’application pour une plateforme afin d’obtenir un binaire à joindre à une version publiée :

@tabs
[macOS]
```bash
wails3 task build:darwin
# produces bin/updater-tutorial.app
# zip it for the release asset:
cd bin && zip -r updater-tutorial-darwin-arm64.zip updater-tutorial.app && cd ..
```

[Linux]
```bash
wails3 task build:linux
# produces bin/updater-tutorial
mv bin/updater-tutorial bin/updater-tutorial-linux-amd64
```

[Windows]
```bash
wails3 task build:windows
# produces bin/updater-tutorial.exe
mv bin/updater-tutorial.exe bin/updater-tutorial-windows-amd64.exe
```

@end

Générez un fichier `SHA256SUMS` à côté du binaire :

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

Une ou plusieurs lignes semblables à celles-ci doivent s’afficher :

```
abc123…  updater-tutorial-darwin-arm64.zip
```

Publiez maintenant ces fichiers en tant que **v2.0.0** dans votre dépôt GitHub :

```bash
gh release create v2.0.0 \
    --title "v2.0.0" \
    --notes "First update for the self-update tutorial.

- **Bold** Markdown renders in the update window
- \`Code spans\` too
- Lists work
- GFM tables work" \
    bin/SHA256SUMS bin/updater-tutorial-*
```

@note{type="note" title="Nommage des ressources"}
Le mécanisme de correspondance des ressources par défaut effectue la sélection à partir des sous-chaînes `GOOS` et `GOARCH` présentes dans le nom du fichier. Il trouvera la ressource si son nom contient `darwin` (ou `linux` / `windows`) ainsi que `arm64` (ou `amd64` / `386`). Pour définir des mécanismes de correspondance personnalisés, consultez le [guide du programme de mise à jour](/guides/updater/#github-releases--updaterprovidersgithub).

@end

### Exécuter l’application et vérifier la mise à jour
Avec `currentVersion` toujours défini sur `1.0.0`, exécutez de nouveau l’application :

```bash
wails3 task dev
```

Cliquez sur **Application → Rechercher des mises à jour…**. Cette fois, vous devriez voir quelque chose de semblable à ceci :

![Fenêtre par défaut du programme de mise à jour dans l’état Mise à jour prête, affichant la pastille de version, les notes de version rendues au format Markdown et le bouton principal Redémarrer et appliquer.](/assets/updater/default-window-ready.png)

- L’icône principale passe d’une flèche bleue ↓ (« Mise à jour disponible ») à une coche verte ✓ (« Mise à jour prête »).
- Le sous-titre affiche `v1.0.0 → v2.0.0 · <size>`.
- Le panneau des notes de version affiche votre Markdown avec le texte en gras, les fragments de code et le tableau.
- La barre de progression se remplit pendant le téléchargement (ce sera rapide, car le binaire est petit).

Le programme de mise à jour place le nouveau binaire dans un répertoire temporaire. Pour terminer la mise à jour :

- Cliquez sur **Redémarrer et appliquer**.
- Votre application se ferme, l’assistant remplace le binaire, puis le nouveau binaire redémarre.
- L’application relancée indique `currentVersion = "1.0.0"` (car nous avons codé cette valeur en dur), mais les octets sur le disque correspondent à la compilation v2.0.0.

Dans une véritable application, `currentVersion` serait défini au moment de la compilation au moyen de `-ldflags`, afin que le nouveau binaire sache qu’il s’agit désormais de la version v2.0.0 et qu’une vérification ultérieure ne trouve aucune mise à jour.

### Relier `currentVersion` à la compilation
Remplacez la constante par une variable définie au moment de la compilation :

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

Puis, dans votre commande de compilation :

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

Vous pouvez aussi ajouter `-ldflags` à votre `Taskfile.yml` afin que la valeur soit récupérée depuis `git describe --tags`.

### Ajouter une signature cryptographique (recommandé en production)
La méthode de vérification utilisant SHA256SUMS vérifie l’*intégrité* (les octets correspondent à ceux stockés par GitHub), mais pas l’*authenticité* (ces octets ont bien été produits par votre pipeline de publication et ne proviennent pas d’un compte de responsable compromis). Pour résister à la falsification, signez chaque version publiée avec une clé Ed25519 :

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

Pour chaque version publiée, signez avec votre clé privée le condensat SHA-256 de chaque artefact. Voici un petit utilitaire Go :

```go {title="cmd/sign-release/main.go"}
package main

import (
    "crypto/ed25519"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

func main() {
    priv, _ := os.ReadFile("updater-key")
    key := ed25519.PrivateKey(priv) // raw 64-byte private key

    f, _ := os.Open(os.Args[1])
    defer f.Close()
    h := sha256.New()
    _, _ = io.Copy(h, f)
    sig := ed25519.Sign(key, h.Sum(nil))
    fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
```

Le fournisseur GitHub par défaut ne récupère actuellement aucun fichier de signature distinct. Vous pouvez [écrire un fournisseur personnalisé](/guides/updater/#writing-your-own-provider) qui le fait, ou passer à **keygen.sh**, qui signe chaque artefact côté serveur et expose le condensat et la signature via son API.

Intégrez la clé publique à votre application :

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

Lorsque `PublicKey` est défini, la vérification de la signature avec cette clé doit réussir pour toute version publiée qui fournit un `Signature`. La source de la version publiée ne peut pas lui substituer sa propre clé : c’est précisément l’intérêt de l’épinglage hors bande au moment de la compilation.

### Personnaliser la fenêtre
La fenêtre par défaut répond aux besoins courants. Si vous avez besoin de davantage de contrôle, trois possibilités s’offrent à vous ; choisissez-en une selon le niveau de personnalisation souhaité :

@tabs
[CSS uniquement]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

Consultez la section [Thème au moyen de variables CSS](/guides/updater/#theme-via-css-variables) pour obtenir la liste complète des variables.

[HTML personnalisé]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

Votre code HTML doit s’abonner aux événements `updater:*` et émettre des actions `updater:user:*` via le canal d’événements Wails. Consultez [Remplacer le modèle](/guides/updater/#replace-the-template) pour obtenir la couche d’adaptation JS.

[Utiliser votre propre fenêtre]
```go
myWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:                "My App Updater",
    Width:                520, Height: 460,
    HTML:                 updaterHTML,
    AllowSimpleEventEmit: true,  // required — see security note
})
app.Updater.Init(updater.Config{
    // …
    Window: updater.BYOWindow(myWin.AsUpdaterWindow()),
})
```

Cette option est utile si vous disposez déjà de votre propre infrastructure de fenêtres et souhaitez que le programme de mise à jour la pilote au lieu d’ouvrir une autre fenêtre. Voici à quoi ressemble un modèle HTML entièrement personnalisé, piloté par les mêmes événements du programme de mise à jour que le modèle par défaut :

![Une fenêtre de mise à jour personnalisée, avec un arrière-plan en dégradé rose-orange et une disposition personnalisée en carte aux angles arrondis, montrant que l’interface utilisateur par défaut est entièrement remplaçable.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` est requis"}
La couche d’adaptation pour HTML personnalisé du programme de mise à jour pilote les actions Install / Skip / Remind / Restart au moyen du raccourci `wails:event:emit:` postMessage. Pour des raisons de sécurité, ce raccourci n’est disponible que si ce champ l’autorise. Si vous oubliez de le définir, les boutons restent inopérants sans afficher d’erreur. Ne l’activez pas dans les fenêtres qui chargent du code HTML que vous ne contrôlez pas entièrement ; consultez la section [Utiliser votre propre fenêtre](/guides/updater/#bring-your-own-window) du guide pour connaître le modèle de menace.

@end

@end

### Exécuter les vérifications automatiques en arrière-plan
Pour effectuer la vérification périodiquement plutôt que lors d’un clic dans le menu, ou en complément de celui-ci :

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

À chaque intervalle, le même flux `CheckAndInstall` que lors d’un clic manuel est exécuté. Définissez `Window: updater.WindowNone` si vous souhaitez que la vérification périodique reste silencieuse jusqu’à ce qu’un élément soit effectivement trouvé ; abonnez-vous alors vous-même à `EventUpdateAvailable` pour choisir l’expérience utilisateur à présenter.

@end

## Vous avez terminé

Vous disposez maintenant d’une application Wails qui :

- Recherche les mises à jour dans les versions publiées sur GitHub à la demande et périodiquement.
- Affiche les notes de version au format Markdown dans une fenêtre par défaut soignée.
- Vérifie les téléchargements à l’aide d’un condensat SHA-256 que vous publiez.
- Peut vérifier une signature Ed25519 à l’aide d’une clé publique que vous intégrez au moment de la compilation.
- Remplace sur place le binaire en cours d’exécution et relance automatiquement l’application.

## Étapes suivantes

- Le [guide du programme de mise à jour](/guides/updater/) contient la documentation de référence complète de l’API, tous les événements, toutes les options de configuration et le mécanisme de remplacement en mode auxiliaire.
- Consultez [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater) pour obtenir un exemple fonctionnel complet que vous pouvez cloner.
- Le dépôt cible de test [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo) illustre l’organisation recommandée des artefacts de version.

## Points à surveiller en production

- **Signature du code sous macOS** — Gatekeeper exige que le binaire de remplacement soit signé et notarié. Signez votre paquet `.app` *avant* de le compresser au format ZIP pour la version publiée. Le programme de mise à jour conserve les octets à l’identique ; il ne signe rien à nouveau.
- **Antivirus sous Windows** — les fichiers `.exe` non signés téléchargés depuis Internet peuvent déclencher des avertissements SmartScreen. Signez votre binaire avec un certificat Authenticode, ou acceptez que les utilisateurs de machines soumises à des restrictions strictes puissent devoir ajouter votre application à leur liste d’autorisation.
- **Versions atomiques** — publiez `SHA256SUMS` et vos binaires ensemble, et non dans des commits distincts. Le programme de mise à jour récupère le fichier annexe séparément du binaire ; s’ils divergent, la vérification du condensat échoue de manière sûre.
- **Versions ignorées** — le bouton « Skip This Version » de la fenêtre par défaut enregistre localement la version ignorée. Si vous publiez une mise à jour de sécurité critique, attribuez-lui un nouveau numéro de version afin qu’elle ne soit pas automatiquement ignorée par les utilisateurs ayant refusé une version antérieure.
