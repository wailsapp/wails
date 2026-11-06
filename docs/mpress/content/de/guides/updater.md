---
title: "Updater"
description: "In-App-Selbstaktualisierungen für Wails v3 — austauschbare Provider, kryptografische Verifizierung, atomarer Austausch und eine standardmäßige Benutzeroberfläche, die Sie anpassen oder ersetzen können."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

Der Updater stellt Softwareupdates innerhalb der App bereit, ohne dass Sie eine eigene Pipeline zum Herunterladen, Verifizieren und Austauschen erstellen müssen. Er baut auf `app.Updater` auf, akzeptiert einen oder mehrere austauschbare `Provider`s (GitHub Releases, keygen.sh, Sparkle AppCast, das offene Wails-Update-Manifest-Protokoll oder eigene Implementierungen), authentifiziert Downloads anhand eines konfigurierten öffentlichen Schlüssels, tauscht die laufende Binärdatei sicher aus und meldet jeden Übergang über den standardmäßigen Wails-Event-Bus.

![Das standardmäßige Updater-Fenster im Zustand „Update bereit“ — zustandsabhängiges Symbol, Versionsanzeige (v1.0.0 → v2.0.1 · 8.8 MB), als Markdown gerenderte Versionshinweise einschließlich einer GFM-Tabelle und eine einzige primäre Aktion.](/assets/updater/default-window-ready.png)

## Schnellstart

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

Dadurch wird das Updatefenster des Frameworks geöffnet, GitHub geprüft, das plattformspezifische Artefakt heruntergeladen und verifiziert, die Binärdatei ausgetauscht und anschließend auf den Neustart durch den Benutzer gewartet.

## Lebenszyklus

`app.Updater` ist ein Zustandsautomat mit den folgenden Zuständen (`updater.State`):

| Zustand | Zeitpunkt |
| --- | --- |
| `unconfigured` | Bevor `Init` aufgerufen wird |
| `idle` | Nach `Init`, vor jeder Prüfung |
| `checking` | Eine `Check` wird gerade ausgeführt |
| `up-to-date` | Laut der letzten Provider-Antwort verwendet der Aufrufer die aktuelle Version |
| `available` | Eine neue Version wurde gefunden, wird aber noch nicht heruntergeladen |
| `downloading` | Bytes werden vom Provider gestreamt |
| `verifying` | Der Download ist abgeschlossen; Signatur bzw. Prüfsumme wird geprüft |
| `installing` | Verifizierte Bytes werden entpackt und in das Staging-Verzeichnis umbenannt |
| `ready` | Das Update wurde bereitgestellt; rufen Sie zum Anwenden `Restart` auf |
| `error` | Ein vorheriger Schritt ist fehlgeschlagen |

Den aktuellen Zustand können Sie jederzeit mit `app.Updater.State()` abrufen. Bei jedem Übergang wird außerdem ein Wails-Event ausgelöst (siehe [Events](#ereignisse)).

`Restart` wartet, bis der Helfer `application.New` erreicht hat, bevor die laufende App zum Beenden aufgefordert wird. Das standardmäßige Startzeitlimit beträgt 30 Sekunden. Wenn Ihre App vor `application.New` eine längere Initialisierung durchführt, legen Sie für `Config.HelperReadyTimeout` eine längere Dauer fest, beispielsweise `time.Minute`. Null wählt den Standardwert; negative Zeitspannen werden abgelehnt. Bei einer Zeitüberschreitung während des Starts gibt `Restart` den Fehler `updater.ErrHelperNotReady` zurück und lässt die laufende App geöffnet.

Das Update wird in einem `wails-update-*`-Verzeichnis neben der installierten Binärdatei bereitgestellt (unter macOS neben dem `.app`-Bundle). Dieses Verzeichnis muss daher beschreibbar sein.

Innerhalb eines AppImage ersetzt und startet der Updater die durch `$APPIMAGE` benannte Datei `.AppImage` neu, nicht den schreibgeschützten Einhängepunkt.

Das standardmäßige Fenster zeigt automatisch den aktuellen Zustand an. Wenn `Check` beispielsweise kein Upgrade zurückgibt, wird dies dem Benutzer angezeigt, und er schließt das Fenster mit **Schließen**:

![Das standardmäßige Updater-Fenster im Zustand „Aktuell“ — grünes Häkchen, Überschrift „Sie sind auf dem neuesten Stand“ und eine einzige Schaltfläche „Schließen“.](/assets/updater/default-window-up-to-date.png)

## Provider

Ein `Provider` ist jede Implementierung, die diese Schnittstelle erfüllt:

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

Vier Implementierungen sind im Quellbaum enthalten.

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

Der standardmäßige Artefaktabgleich sucht im Dateinamen nach den Teilzeichenfolgen `GOOS` und `GOARCH` und erkennt dabei gängige Aliasse (`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`). Für benutzerdefinierte Benennungsschemata:

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

`ChecksumAsset` ist der Name eines weiteren Release-Artefakts, dessen Inhalt aus `<sha256>  <filename>`-Zeilen besteht (das Format, das `sha256sum` und `shasum -a 256` erzeugen). Der Provider ruft es während `Check` ab, sucht die zum ausgewählten Artefakt passende Zeile und befüllt `Release.Verification.Digest`, damit das Framework den Download verifiziert.

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

Der Provider überträgt die artefaktspezifische SHA-512-Prüfsumme und Ed25519ph-Signatur von keygen.sh automatisch in den `Release.Verification`-Block des Frameworks — eine zusätzliche Konfiguration ist nicht erforderlich.

**Tokenformate:** keygen.sh-Token enthalten ein Rollenpräfix (`admi-` / `prod-` / `envi-` / `user-`). Die im Dashboard angezeigte reine UUID ist der *Bezeichner* des Tokens, nicht sein geheimer Wert. Das Geheimnis ist nur beim Erstellen des Tokens sichtbar. Einzelheiten finden Sie in der [Authentifizierungsdokumentation](https://keygen.sh/docs/api/authentication/) von keygen.sh.

### Sparkle AppCast — `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

Lässt sich unverändert in eine vorhandene Sparkle-/WinSparkle-Infrastruktur integrieren. Liest `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>` und `sparkle:channel` aus dem Feed.

Die DSA-Signaturen (`sparkle:dsaSignature`) von Sparkle 1 werden nicht unterstützt. Projekte, die dieses Signaturschema verwenden, sollten auf EdDSA (Sparkle 2) umstellen.

### Wails Update Manifest — `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

Verwendet das offene [Wails-Update-Manifest-Protokoll](/reference/update-manifest/): ein JSON-Dokument, das die neueste Version und deren plattformspezifische Artefakte beschreibt und Prüfsummen sowie Signaturen direkt enthält. Dasselbe Dokument funktioniert sowohl auf einem statischen Datei-Host (S3, GitHub Pages oder einem beliebigen CDN — veröffentlichen Sie pro Kanal ein Manifest, das alle Plattformen aufführt) als auch auf einem dynamischen Updateserver (der Provider sendet bei jeder Prüfung `platform`, `arch`, `version` und `channel`, sodass der Server genau ein Artefakt zurückgeben oder den Zugriff von einer Lizenz abhängig machen kann).

Mit URL-Platzhaltern lassen sich statische Strukturen in einer einzigen Konfigurationszeile angeben:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

Konfigurierte Header werden bei jeder Manifestanfrage gesendet. Für Artefaktdownloads werden sie nur auf dem Host des Manifests und nur ohne Herabstufung von `https` auf `http` wiederverwendet. Bei jeder ursprungsübergreifenden oder herabstufenden Weiterleitung wird der Header `Authorization` entfernt.

Die CLI übernimmt die Veröffentlichung: `wails3 updater manifest` berechnet Prüfsummen, signiert und beschreibt Ihre Release-Dateien mit einem einzigen Befehl; `wails3 updater verify` prüft das Ergebnis vor dem Hochladen erneut. Siehe [Veröffentlichen mit der wails3-CLI](/reference/update-manifest/#publishing-with-the-wails3-cli).

### Fallback-Kette

`Config.Providers` ist geordnet. Der Updater durchläuft die Einträge nacheinander: Der erste Provider, der eine Version zurückgibt, gewinnt. Der erste Provider, der „aktuell“ meldet, beendet die Kette vorzeitig (der Fallback ist für den Fall „primärer Provider nicht erreichbar“ vorgesehen, nicht für „Provider sind unterschiedlicher Meinung“). Bei einem Fehler wird mit dem nächsten Provider fortgefahren.

```go
app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers: []updater.Provider{
        kg, // primary: licensed customers
        gh, // fallback: public mirror
    },
})
```

### Eigenen Provider erstellen

Drei Methoden und bei einer typischen Implementierung etwa 150 Zeilen. Der Updater übernimmt Verifizierung, atomares Staging, Austausch und Fenster; der Providercode ermittelt die nächste Version und streamt die Bytes:

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

Verwenden Sie die im Quellbaum enthaltenen Provider als Referenz — jeder besteht aus einer einzelnen Go-Datei.

## Kryptografische Verifizierung

Der Verifizierer des Frameworks authentifiziert Releases mit `Config.PublicKey` als Vertrauensanker:

```go
//go:embed publickey.pem
var publicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: "1.0.0",
    Providers:      []updater.Provider{...},
    PublicKey:      publicKey,
})
```

Unterstützte Algorithmen (`Release.Verification.SignatureAlgo`):

| Algorithmus | Signierter Inhalt | Hinweise |
| --- | --- | --- |
| `ed25519` | Der SHA-256-Digest des Artefakts | Von Sparkle EdDSA verwendet |
| `ed25519ph` | Das vollständige Artefakt über den Pre-Hash von Ed25519ph (intern SHA-512) | Von keygen.sh verwendet |
| `ecdsa-p256` | Der SHA-256-Digest des Artefakts | Sowohl unformatierte `r∥s`- als auch DER-Signaturen werden akzeptiert |

Außerdem nur Digest (`DigestAlgo`: `sha256` / `sha512`), wenn ein Release einen Hash, aber keine Signatur enthält.

`Config.PublicKey` ist der EINZIGE Vertrauensanker für die Signaturprüfung – die Release-Quelle kann keinen eigenen Schlüssel einsetzen. Releases, die eine `Signature` enthalten, obwohl keine `Config.PublicKey` konfiguriert ist, werden sicher abgelehnt. Der Verifizierer berechnet den Digest während des Downloads in einem Streaming-Durchlauf. Daher erfordert die Prüfung selbst bei Updates von mehreren GB keinen zusätzlichen Datenträgerdurchlauf.

@note{type="caution" title="Nur Digest ≠ kryptografische Verifizierung"}
Ein Release, das nur `Digest` enthält, wird anhand der TLS-Verbindung der Registry und der von der Registry selbst bereitgestellten Integritätsgarantie authentifiziert – nicht anhand eines von Ihnen kontrollierten kryptografischen Vertrauensankers. Verwenden Sie nur einen Digest, um Bitfehler zu erkennen, und Signaturen, um Manipulationen durch eine kompromittierte Release-Pipeline abzuwehren.

@end

### Signierschlüssel erzeugen

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

Der private Schlüssel ist PKCS#8-PEM, der öffentliche Schlüssel PKIX-PEM. `Config.PublicKey` akzeptiert die Datei `.pub` unverändert. Ebenfalls akzeptiert werden der unformatierte, 32 Byte lange Schlüssel oder dessen Base64-Darstellung, die `genkey` zur direkten Einbettung ausgibt. Signieren Sie Releases mit `wails3 updater manifest -key updater.key ...` oder `wails3 updater sign`; siehe [Veröffentlichen mit der wails3-CLI](/reference/update-manifest/#publishing-with-the-wails3-cli).

Oder in Go:

```go
import "crypto/ed25519"
import "crypto/rand"

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
// Persist `priv` securely (HSM, signing CI, etc.); embed `pub` in your binary.
```

## Artefaktformate

Provider streamen die Bytes der jeweils veröffentlichten Datei; das Framework entpackt sie anschließend vor dem Austausch:

- **Einzelne Binärdatei** (z. B. `myapp-linux-amd64`) – wird unverändert verwendet. Unter Linux üblich.
- **`.zip`** – wird direkt am Zielort entpackt. Das Archiv muss genau einen Eintrag auf oberster Ebene enthalten (üblicherweise ein macOS-`.app`-Bundle oder eine einzelne Binärdatei). Empfohlenes Paketformat für macOS.
- **`.tar.gz`** / **`.tgz`** – wird unter derselben Regel eines einzigen Eintrags auf oberster Ebene direkt am Zielort entpackt. Nützlich für Linux-Distributionen, die neben der Binärdatei einen Runtime-Verzeichnisbaum ausliefern.

Archive mit mehr als einem Eintrag auf oberster Ebene werden abgelehnt: Das Framework tauscht ein einzelnes Ziel auf dem Datenträger aus. Daher ist „dieses Archiv am Zielort einsetzen“ mehrdeutig, wenn das Archiv mehrere Elemente enthält. `.dmg` und `.pkg` (macOS) sowie `.msi` (Windows) werden in v1 nicht unterstützt – verteilen Sie stattdessen ein `.zip` des Bundles. Beim Entpacken wird Zip-Slip-Schutz erzwungen, Symlinks, die aus dem Archivstamm herausführen, werden abgelehnt und die unkomprimierte Gesamtgröße (2 GiB) sowie die Anzahl der Einträge (50 000) werden begrenzt.

## Das Standardfenster

`app.Updater.CheckAndInstall(ctx)` öffnet ein vom Framework verwaltetes Fenster der Größe 520 × 540 mit:

- Einem statusabhängigen Hauptsymbol (blaues ↓ für verfügbar/wird heruntergeladen, grünes ✓ für bereit/aktuell, rotes ! für Fehler)
- Einer Versionsplakette: `v1.0.0 → v2.0.1 · 8.8 MB`
- Einem scrollbaren Bereich für Release-Hinweise mit **gerendertem Markdown** (Absätze, Fett-/Kursivschrift, Listen, GFM-Tabellen, Inline-Code, abgegrenzte Codeblöcke, h1–h3, Links)
- Einer einzigen primären Aktion pro Status (Installieren / Neu starten und anwenden / Erneut versuchen)
- Sekundären Aktionen im Ghost-Stil (Diese Version überspringen / Später erinnern)
- Dunkel-/Hellmodus über `prefers-color-scheme`
- Einem unbestimmten Fortschrittsschimmer, wenn die Gesamtgröße unbekannt ist

Es lauscht im Wails-Event-Bus auf `updater:*`-Events und sendet `updater:user:*`-Aktionen an Go zurück.

### Design über CSS-Variablen

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

Das Standard-Stylesheet stellt diese Variablen bereit – Sie können jede davon überschreiben:

| Variable | Standard (hell) | Standard (dunkel) |
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
| `--font` | Systemschriftarten | — |

### Vorlage ersetzen

Stellen Sie eigenes HTML bereit; es muss lediglich auf `wails:updater:*`-Events lauschen und `wails:updater:user:*`-Aktionen senden:

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML-Fenster werden ohne Ursprung des Asset-Servers geladen und können `/wails/runtime.js` daher nicht dynamisch abrufen. Es gibt zwei Möglichkeiten, von einem solchen Fenster aus mit dem Host zu kommunizieren:

1. **Schreiben Sie einfach HTML.** Das Framework injiziert automatisch einen minimalen `window.wails.Events`-Shim in jedes Fenster, das mit gesetztem `WebviewWindowOptions.AllowSimpleEventEmit = true` und `HTML` geöffnet wird – genau das geschieht bei den integrierten und den BYO-Pfaden des Updaters. Kein Build-Schritt erforderlich. Das folgende Beispiel verwendet diesen Pfad.
2. **Bündeln Sie `@wailsio/runtime` mit Ihrem bevorzugten Bundler** (Vite, esbuild, Rollup) und importieren Sie es zur Build-Zeit in Ihr benutzerdefiniertes HTML. `Events.On` funktioniert sofort, da es ausschließlich clientseitig arbeitet; `Events.Emit` verwendet den Fetch-Transport der Runtime, der aufgrund des Null-Ursprungs nicht funktioniert. Installieren Sie daher über den Hook [`setTransport`](/guides/custom-transport/) der Runtime einen kleinen postMessage-Transport, der über `window._wails.invoke("wails:event:emit:<name>")` weiterleitet. Die Framework-Injektion führt keine Aktion aus, wenn `window.wails.Events` bereits im Gültigkeitsbereich vorhanden ist, sodass sich die beiden Ansätze nicht gegenseitig beeinträchtigen.

In beiden Fällen ruft das JavaScript in Ihrem benutzerdefinierten HTML dieselbe API `Events.On` / `Events.Emit` auf:

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

Der Shim stellt die Teilmenge der modernen Runtime bereit, die Ereignisse mit einfachen Namen benötigen: `Events.On(name, cb)` gibt eine Funktion zum Abbestellen zurück, und `Events.Emit(nameOrEventObject)` leitet über den geschützten postMessage-Pfad `wails:event:emit:` an den Host weiter. Er wird einmal beim Laden der Seite installiert, bevor eines Ihrer Inline-Skripte ausgeführt wird.

Wenn Sie den Shim *überschreiben möchten* oder die vollständige Runtime auf andere Weise laden, setzen Sie `window.wails.Events`, bevor das erste `<script>`-Tag der Seite ausgeführt wird; die Injektion überspringt sich dann selbst.

### Fensterrahmen

Überschreiben Sie die Fensteroptionen (Größe, rahmenlos, immer im Vordergrund), ohne das HTML zu ändern:

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

### Eigenes Fenster verwenden

Steuern Sie den Aktualisierungsablauf über ein selbst erstelltes `*application.WebviewWindow`. Der Updater ruft `Show()` / `Close()` / `EmitEvent()` für Ihr Fenster auf; Ihr HTML bestimmt, was dargestellt wird:

![Ein eigenes Updater-Fenster mit einem pink-orangefarbenen Verlaufshintergrund, einer einzelnen weißen abgerundeten Karte, benutzerdefinierter Typografie und denselben Updater-Ereignissen, die den sichtbaren Zustand steuern. Zeigt, wie vollständig sich die Standardbenutzeroberfläche ersetzen lässt.](/assets/updater/byo-custom-window.png)

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

Ihr HTML verwendet `window.wails.Events.On` / `Events.Emit` genau wie die integrierte Vorlage. Der automatisch vom Framework injizierte Shim wird in jedes Fenster mit `AllowSimpleEventEmit: true` eingebunden, unabhängig davon, ob das Fenster dem Framework oder Ihnen gehört. Informationen zur API finden Sie unter [Vorlage ersetzen](#vorlage-ersetzen).

@note{type="caution" title="`AllowSimpleEventEmit` ist für eigene Updater-Fenster erforderlich"}
Aus Sicherheitsgründen macht das Framework die postMessage-Abkürzung `wails:event:emit:` von diesem Feld abhängig: Ein Fenster, bei dem es nicht gesetzt ist, kann keine benutzerdefinierten hostseitigen Ereignisse erzeugen. Der Shim für benutzerdefiniertes HTML des Updaters sendet `updater:user:*`-Ereignisse über diese Abkürzung. Vergessen Sie das Feld bei einem eigenen Fenster, werden daher sämtliche Schaltflächenklicks stillschweigend verworfen: Der Benutzer klickt auf Installieren, und nichts geschieht.

Lassen Sie `AllowSimpleEventEmit` für jedes Fenster **deaktiviert**, das HTML lädt, über das Sie nicht die vollständige Kontrolle haben (Remote-URLs, von Benutzern bereitgestellte Inhalte). Ist die Option aktiviert, kann jedes JavaScript auf der Seite – einschließlich XSS-Senken – jeden `app.Event.On(name, …)`-Handler auslösen. Die Abkürzung übermittelt nur einfache Namen ohne Nutzlast und kann den Binding-/Call-Pfad nicht erreichen. Sie kann jedoch weiterhin privilegierte Handler für benutzerdefinierte Ereignisse in Ihrem Go-Code auslösen, wenn diese Handler allein anhand des Ereignisnamens handeln.

Beim *integrierten* Updater-Fenster setzt das Framework diese Option intern; nur Aufrufer mit eigenem Fenster müssen daran denken.

@end

### Ohne Benutzeroberfläche

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

Es wird zu keinem Zeitpunkt ein Fenster geöffnet. Abonnieren Sie `updater:*`-Ereignisse über Ihre eigene Benutzeroberfläche oder Ihr vorhandenes Hauptfenster und rufen Sie `app.Updater.CheckAndInstall(ctx)` in einem Schaltflächen-Handler auf. Dies eignet sich für regelmäßige Hintergrundprüfungen, die nur sichtbar werden sollen, wenn etwas gefunden wurde, oder für Anwendungen, die den Aktualisierungsablauf in einen benutzerdefinierten Einstellungsbereich integrieren.

## Ereignisse

Sowohl Go als auch JavaScript abonnieren Ereignisse über den standardmäßigen Wails-Ereignisbus. **Geben Sie die Übertragungszeichenfolgen nicht manuell ein**, sondern verwenden Sie die Konstanten, die vom Updater-Paket (Go) beziehungsweise Runtime-Paket (JS) exportiert werden. Beide Ebenen verwenden dieselben Namen, deren Synchronität ein Regressionstest sicherstellt.

### Aus Go

Die Konstanten befinden sich in `github.com/wailsapp/wails/v3/pkg/updater`. Abonnieren Sie sie über `app.Event.On(name, fn)`. Der Callback erhält ein `*application.CustomEvent`, dessen Feld `Data` die in der [Ereignisreferenz](#ereignisreferenz) aufgeführte typisierte Nutzlast enthält. Verwenden Sie eine Typassertion und keine JSON-Dekodierung:

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

Alle verfügbaren Go-Konstanten:

| Konstante | Übertragungszeichenfolge |
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

### Aus JavaScript

Die Konstanten befinden sich unter `Updater.Events` in `@wailsio/runtime`. Sie tragen dieselben Namen wie in Go und sind für eine bessere Auffindbarkeit bei der automatischen Vervollständigung nach Unter-Namensräumen (`User.*`, `Window.*`) gegliedert:

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

Die Benutzeraktionsereignisse, die Ihr benutzerdefiniertes HTML an den Host *zurücksendet*, befinden sich unter `Updater.Events.User`:

```js
import { Updater } from "@wailsio/runtime";

document.getElementById("install-btn").addEventListener("click", () => {
    // The framework window does this internally via the postMessage shim;
    // shown here for BYO templates that need to drive the flow themselves.
    window._wails.invoke("wails:event:emit:" + Updater.Events.User.Install);
});
```

### Ereignisreferenz

Abonnementseite (Host → Seite):

| Konstante (Go) | Konstante (JS) | Nutzlast | Zeitpunkt |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | keine | Vor jedem `Check`-Roundtrip |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` hat eine neuere Version gefunden |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | keine | `Check` hat bestätigt, dass die aktuelle Version installiert ist |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | Byte-Streaming beginnt |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | ~10 Hz während des Downloads |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | Alle Bytes geschrieben, vor der Verifizierung |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | Signatur-/Digest-Prüfung beginnt |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | Entpacken und Staging beginnen |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | Neustart steht aus |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | Eine beliebige Phase ist fehlgeschlagen |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | Einmal pro Sitzung vor der Wiedergabe des Snapshots |

Seitenseitig (Seite → Host) – Ihr Code abonniert diese Ereignisse, wenn Sie eine benutzerdefinierte Vorlage erstellen:

| Konstante (Go) | Konstante (JS) | Zeitpunkt |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | Fenster vollständig geladen; der Host gibt den aktuellen Zustand wieder |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | Primäre Aktion im Zustand `available` |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | Primäre Aktion im Zustand `ready` |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | „Diese Version überspringen“ |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | „Später erinnern“ |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | Schließen-Schaltfläche |

## API-Referenz

### `updater.Config`

| Feld | Typ | Hinweise |
| --- | --- | --- |
| `CurrentVersion` | `string` | **Erforderlich.** Dieselbe Zeichenfolge, mit der Sie Releases taggen (ohne Präfix `v`) |
| `Providers` | `[]updater.Provider` | **Erforderlich.** Geordnete Fallback-Kette |
| `PublicKey` | `[]byte` | PEM oder Rohbytes. Optional, aber signierte Releases schlagen ohne diesen Wert sicher fehl |
| `CheckInterval` | `time.Duration` | Ein Wert ungleich null startet im Hintergrund eine Abfrageschleife, die `CheckAndInstall` aufruft |
| `Platform` | `string` | Überschreibt `runtime.GOOS` für die Asset-Auswahl |
| `Arch` | `string` | Überschreibt `runtime.GOARCH` für die Asset-Auswahl |
| `Channel` | `string` | Derzeit nur informativ; providerspezifische Kanalfilterung |
| `Window` | `updater.WindowOption` | `nil` (integrierte Standardwerte), `&BuiltinWindow{…}`, `BYOWindow(handle)` oder `WindowNone` |

### Methoden von `*updater.Updater`

| Signatur | Zweck |
| --- | --- |
| `Init(cfg Config) error` | Konfiguriert den Updater. Gibt beim zweiten Aufruf `ErrAlreadyConfigured` zurück |
| `State() State` | Aktuelle Lebenszyklusphase |
| `CurrentVersion() string` | Die an `Init` übergebene Version |
| `Check(ctx) (*Release, error)` | Durchläuft die Provider-Kette. `(rel, nil)` = gefunden, `(nil, nil)` = aktuell, `(nil, err)` = alle fehlgeschlagen |
| `DownloadAndInstall(ctx) error` | Streamt, verifiziert, entpackt (bei Archiven) und führt das Staging durch. Erfordert einen vorherigen Aufruf von `Check` |
| `CheckAndInstall(ctx) error` | Komfortmethode: Öffnet das Fenster, führt `Check` aus und anschließend `DownloadAndInstall`, falls ein Update gefunden wurde |
| `Restart(ctx) error` | Startet den Helper, ruft `Host.Quit` auf und beendet die Anwendung; der Helper tauscht die Anwendung aus und startet sie neu |
| `DownloadedPath() string` | Speicherort des bereitgestellten Updates auf dem Datenträger oder `""`, falls keines vorhanden ist |
| `SkipVersion(v string)` | Vermerkt `v` als übersprungen; nachfolgende `Check`-Aufrufe behandeln diese Version als aktuell |
| `SkippedVersion() string` | Liest die derzeit übersprungene Version |
| `StopPeriodicCheck()` | Beendet den von `Config.CheckInterval` gestarteten Timer und wartet, bis die Schleife zurückkehrt |

### Fehler

| Sentinel-Fehler | Zurückgegeben von |
| --- | --- |
| `ErrAlreadyConfigured` | `Init` nach dem ersten Erfolg |
| `ErrNotConfigured` | Jeder Vorgang vor `Init` |
| `ErrNoPendingRelease` | `DownloadAndInstall` ohne vorheriges `Check` |
| `ErrDownloadInProgress` | `DownloadAndInstall` wird aufgerufen, während bereits ein anderer Vorgang läuft |
| `ErrNotReady` | `Restart` ohne bereitgestelltes Update |

## So funktioniert der Austausch

`Restart` führt das aktuelle Binärprogramm mit gesetzten Sentinel-Umgebungsvariablen erneut aus. `application.New` erkennt sie beim Start und wechselt in den Hilfsmodus:

1. Das Hilfsprogramm wartet bis zu 30 s darauf, dass die übergeordnete PID beendet wird (`platformIsAlive` fragt unter Windows über `syscall.OpenProcess` + `GetExitCodeProcess` und unter Unix über `os.FindProcess` + `proc.Signal(syscall.Signal(0))` ab).
2. Das Hilfsprogramm sichert das Ziel (Kopie bei Dateien, rekursive Kopie bei macOS-Verzeichnissen für `.app`-Bundles).
3. Das Hilfsprogramm ersetzt das Ziel durch das bereitgestellte Artefakt. Bei Fehlern versucht es dies bis zu 20-mal erneut und wartet zwischen den Versuchen jeweils 500 ms:
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. Offene Dateideskriptoren, die auf den alten Inode verweisen, bleiben gültig.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. Windows erlaubt das Umbenennen von Dateien, deren Image noch gemappt ist, nicht jedoch deren Löschen. Beim nächsten Update entfernt das Hilfsprogramm alle verbliebenen `.old.*`-Geschwisterdateien, deren zugehöriges Kernel-Mapping freigegeben wurde.

4. Das Hilfsprogramm stellt beim neuen Binärprogramm den ursprünglichen Modus der ausführbaren Datei wieder her (die heruntergeladene Datei wurde mit der standardmäßigen umask erstellt, wodurch `+x` unter Unix entfällt; unter Windows ist dies eine wirkungslose Operation).
5. Das Hilfsprogramm entfernt die Umgebungsvariablen des Hilfsmodus und startet das nun ersetzte Binärprogramm erneut.
6. Das Hilfsprogramm wird beendet.

Schlägt der Start fehl, stellt das Hilfsprogramm die Sicherung wieder her. Wird der übergeordnete Prozess nicht innerhalb von 30 s beendet, bricht das Hilfsprogramm ab, bevor es das Ziel verändert. So behält der Benutzer eine funktionsfähige App, selbst wenn ein Dialog beim Beenden `Quit` blockiert.

Bei macOS-`.app`-Bundles, die als `.zip` verteilt werden (die empfohlene Paketierung), wird das Archiv zwischen der Verifizierung und der Bereitschaftsmeldung entpackt, damit das Hilfsprogramm ein echtes Verzeichnis austauschen kann.

## Regelmäßige Prüfung

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

Wenn `CheckInterval > 0` gilt, ruft eine Hintergrund-Goroutine `CheckAndInstall` im konfigurierten Intervall auf. Ticks, die eintreffen, während bereits ein anderer Ablauf ausgeführt wird (Prüfen, Herunterladen, Verifizieren oder Installieren), werden verworfen — parallele Zustandsautomaten werden nicht unterstützt.

Legen Sie für eine stille Hintergrundabfrage, die nur bei einem Fund sichtbar wird, `Window: updater.WindowNone` fest und reagieren Sie in Ihrer eigenen Benutzeroberfläche auf `EventUpdateAvailable`.

## Überspringen und später erinnern

Die Schaltfläche „Diese Version überspringen“ im Standardfenster speichert die verfügbare Version über `SkipVersion(rel.Version)`. Nachfolgende `Check`s finden dieselbe Version und behandeln die App als aktuell, bis der Benutzer `CurrentVersion` aktualisiert; dies geschieht nach einem erfolgreichen `Restart` automatisch. „Später erinnern“ schließt lediglich das Fenster, ohne etwas zu speichern.

```go
// Reading what the user skipped (e.g. to surface in app settings)
if v := app.Updater.SkippedVersion(); v != "" {
    log.Printf("user skipped %s", v)
}

// Programmatically clearing the skip:
app.Updater.SkipVersion("")
```

## Checkliste für die Distribution

Bevor Sie ein Release veröffentlichen, das der Updater installieren soll:

1. **Wählen Sie das richtige Archivformat.** macOS: `.zip` des `.app`-Bundles. Linux: einzelnes Binärprogramm oder `.tar.gz`. Windows: einzelne `.exe`- oder `.zip`-Datei. `.dmg` / `.msi` / `.pkg` werden nicht unterstützt.
2. **Signieren Sie das Artefakt** mit dem privaten Schlüssel, der zu `Config.PublicKey` gehört. Befolgen Sie bei vom Anbieter veröffentlichten Feeds (keygen.sh, AppCast) den Signierablauf des jeweiligen Anbieters. Erzeugen Sie für GitHub Releases mit `ChecksumAsset` eine `SHA256SUMS`-Datei mit `sha256sum` / `shasum -a 256`.
3. **Gleichen Sie die Versionszeichenfolge ab.** `Config.CurrentVersion` und das Versions-Tag des Releases müssen exakt übereinstimmen (z. B. `1.0.0` ↔ Tag `v1.0.0`; das führende `v` wird auf Anbieterseite entfernt).
4. **Testen Sie den Austausch auf der Zielplattform** vor der Auslieferung mindestens einmal — Codesignierung, Notarisierung und die Handhabung durch Gatekeeper sind plattformspezifisch und werden nicht vom Updater selbst abgedeckt.

## Fehlerbehebung

**„Für die Signatur ist ein öffentlicher Schlüssel erforderlich, es wurde jedoch keiner konfiguriert“** — das Release enthält ein `Signature`-Feld, aber `Config.PublicKey` ist leer. Legen Sie den öffentlichen Schlüssel fest oder ändern Sie Ihre Release-Pipeline so, dass sie keine Signatur einfügt.

**„Prüfsumme stimmt nicht überein“** — die heruntergeladenen Bytes entsprechen nicht den Angaben des Anbieters. Ursache ist meist ein unvollständiger Download aufgrund einer kurzen Netzwerkstörung oder ein beschädigtes Artefakt. Ein erneuter Versuch behebt das Problem häufig.

**Das Fenster wird geöffnet, verschwindet aber sofort; weder Markdown noch Fortschritt werden angezeigt** — Ihr benutzerdefiniertes HTML hat `wails:runtime:ready` nicht aufgerufen. Siehe den Shim unter [Vorlage ersetzen](#vorlage-ersetzen).

**Das Windows-Update wird nie abgeschlossen; im Protokoll des Hilfsprogramms steht „remove old (attempt N): Access is denied“** — dies tritt nur bei Versionen dieses PRs vor `de764fb` auf. Die aktuelle Implementierung benennt die alte Datei stattdessen um und ist daher nicht betroffen. Führen Sie ein Upgrade durch.

**macOS Gatekeeper blockiert das ausgetauschte Binärprogramm** — die Codesignierung muss über den gesamten Ablauf hinweg erhalten bleiben. Signieren Sie das ursprüngliche `.app` *und* signieren Sie das neu gestartete Binärprogramm erneut, wenn Ihre Build-Pipeline die Berechtigungen während des Updates ändert.

## Siehe auch

- Ausführbares Beispiel: [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/v3.0.0/v3/examples/updater)
- Demo-Repository für Tests: [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- Tutorial: [Selbstaktualisierungen zu einer Wails-App hinzufügen](/tutorials/04-self-update-a-wails-app/)
