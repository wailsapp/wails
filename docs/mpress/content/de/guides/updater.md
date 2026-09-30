---
title: "Updater"
description: "In-App-Selbstaktualisierungen für Wails v3 – austauschbare Provider, kryptografische Verifizierung, atomarer Austausch und eine standardmäßige Benutzeroberfläche, die Sie gestalten oder ersetzen können."
slug: "guides/updater"
sourcePath: "guides/updater.md"
---

Der Updater stellt Softwareupdates innerhalb der Anwendung bereit, ohne dass Sie eine eigene Pipeline zum Herunterladen, Verifizieren und Austauschen entwickeln müssen. Er basiert auf `app.Updater`, akzeptiert einen oder mehrere austauschbare `Provider` (GitHub Releases, keygen.sh, Sparkle AppCast, das offene Wails-Update-Manifest-Protokoll oder eigene Provider), authentifiziert Downloads anhand eines konfigurierten öffentlichen Schlüssels, tauscht die laufende Binärdatei sicher aus und meldet jeden Zustandsübergang über den standardmäßigen Wails-Event-Bus.

![Das standardmäßige Updater-Fenster im Zustand „Update bereit“ – zustandsabhängiges Symbol, Versionsanzeige (v1.0.0 → v2.0.1 · 8.8 MB), als Markdown gerenderte Versionshinweise einschließlich einer GFM-Tabelle und eine einzige primäre Aktion.](/assets/updater/default-window-ready.png)

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
| `idle` | Nach `Init`, bevor überhaupt eine Prüfung stattgefunden hat |
| `checking` | Ein Aufruf von `Check` läuft |
| `up-to-date` | Laut der letzten Provider-Antwort ist der Aufrufer auf dem aktuellen Stand |
| `available` | Eine neue Version wurde gefunden; der Download hat noch nicht begonnen |
| `downloading` | Bytes werden vom Provider gestreamt |
| `verifying` | Der Download ist abgeschlossen; Signatur/Prüfsumme wird geprüft |
| `installing` | Verifizierte Bytes werden entpackt und in das Staging-Verzeichnis verschoben |
| `ready` | Das Update wurde bereitgestellt; rufen Sie zum Anwenden `Restart` auf |
| `error` | Ein vorheriger Schritt ist fehlgeschlagen |

Sie können den aktuellen Zustand jederzeit mit `app.Updater.State()` abrufen. Jeder Zustandsübergang löst außerdem ein Wails-Ereignis aus (siehe [Ereignisse](#ereignisse)).

`Restart` wartet, bis der Hilfsprozess `application.New` erreicht, bevor die laufende Anwendung zum Beenden aufgefordert wird. Das standardmäßige Startzeitlimit beträgt 30 Sekunden. Wenn Ihre Anwendung vor `application.New` eine längere Initialisierung durchführt, setzen Sie `Config.HelperReadyTimeout` auf eine längere Dauer, beispielsweise `time.Minute`. Null wählt den Standardwert; negative Zeitspannen werden abgelehnt. Bei Überschreitung des Startzeitlimits gibt `Restart` den Fehler `updater.ErrHelperNotReady` zurück und lässt die laufende Anwendung geöffnet.

Das Standardfenster zeigt automatisch den aktuellen Zustand an. Wenn `Check` beispielsweise kein Upgrade zurückgibt, wird dies dem Benutzer angezeigt und er schließt das Fenster mit **Schließen**:

![Das standardmäßige Updater-Fenster im Zustand „Aktuell“ – grünes Häkchen, Überschrift „Sie sind auf dem aktuellen Stand“ und eine einzige Schaltfläche „Schließen“.](/assets/updater/default-window-up-to-date.png)

## Provider

Als `Provider` kann alles dienen, was diese Schnittstelle implementiert:

```go
type Provider interface {
    Name() string
    Check(ctx context.Context, req CheckRequest) (*Release, error)
    Download(ctx context.Context, r *Release, dst io.Writer, onProgress func(written, total int64)) error
}
```

Vier Implementierungen sind im Quellbaum enthalten.

### GitHub Releases – `updater/providers/github`

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

Der standardmäßige Artefaktabgleich wählt anhand der Teilzeichenfolgen `GOOS` und `GOARCH` im Dateinamen aus und erkennt dabei gängige Aliasse (`amd64` / `x86_64` / `x64`, `arm64` / `aarch64`, `386` / `i386` / `x86` / `ia32`). Für eigene Benennungsschemata:

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

`ChecksumAsset` ist der Name eines weiteren Release-Artefakts, dessen Inhalt aus `<sha256>  <filename>`-Zeilen besteht (dem Format, das `sha256sum` und `shasum -a 256` erzeugen). Der Provider ruft es während `Check` ab, sucht die zum ausgewählten Artefakt passende Zeile und befüllt `Release.Verification.Digest`, damit das Framework den Download verifiziert.

### keygen.sh – `updater/providers/keygen`

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

Der Provider überträgt die artefaktspezifische SHA-512-Prüfsumme und Ed25519ph-Signatur von keygen.sh automatisch in den `Release.Verification`-Block des Frameworks – eine zusätzliche Verkabelung ist nicht erforderlich.

**Tokenformate:** keygen.sh-Token enthalten ein Rollenpräfix (`admi-` / `prod-` / `envi-` / `user-`). Die im Dashboard angezeigte reine UUID ist der *Bezeichner* des Tokens, nicht sein geheimer Wert – dieser ist nur beim Erstellen des Tokens sichtbar. Einzelheiten finden Sie in der [Authentifizierungsdokumentation](https://keygen.sh/docs/api/authentication/) von keygen.sh.

### Sparkle AppCast – `updater/providers/appcast`

```go {title="appcast provider"}
ac, err := appcast.New(appcast.Config{
    URL:        "https://your.app/appcast.xml", // required
    Channel:    "stable",                       // optional sparkle:channel filter
    HTTPClient: nil,                            // optional client override
})
```

Lässt sich unverändert in eine vorhandene Sparkle-/WinSparkle-Infrastruktur integrieren. Liest `sparkle:shortVersionString`, `<enclosure url type length sparkle:os sparkle:edSignature>` und `sparkle:channel` aus dem Feed.

Die DSA-Signaturen (`sparkle:dsaSignature`) von Sparkle 1 werden nicht unterstützt. Projekte mit diesem Signaturschema sollten auf EdDSA (Sparkle 2) umstellen.

### Wails Update Manifest – `updater/providers/endpoint`

```go {title="endpoint provider"}
ep, err := endpoint.New(endpoint.Config{
    URL:        "https://updates.example.com/check", // required; supports {{platform}} / {{arch}} / {{version}} / {{channel}} placeholders
    Channel:    "stable",                            // optional channel filter
    Headers:    nil,                                 // optional headers, e.g. {"Authorization": "License <key>"}
    HTTPClient: nil,                                 // optional client override; redirect-strip wrapper still applied
})
```

Verwendet das offene [Wails-Update-Manifest-Protokoll](/reference/update-manifest/): ein einzelnes JSON-Dokument, das die neueste Version und deren plattformspezifische Artefakte mit eingebetteten Prüfsummen und Signaturen beschreibt. Dasselbe Dokument funktioniert sowohl auf einem statischen Dateihost (S3, GitHub Pages oder einem beliebigen CDN – veröffentlichen Sie pro Kanal ein Manifest, das alle Plattformen aufführt) als auch auf einem dynamischen Updateserver. Bei jeder Prüfung sendet der Provider `platform`, `arch`, `version` und `channel`, sodass der Server genau ein Artefakt zurückgeben oder den Zugriff von einer Lizenz abhängig machen kann.

Mit URL-Platzhaltern lassen sich statische Layouts in einer einzigen Konfigurationszeile angeben:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL: "https://cdn.example.com/updates/{{platform}}/{{arch}}/stable.json",
})
```

Konfigurierte Header werden bei jeder Manifestanfrage gesendet. Bei Artefaktdownloads werden sie nur auf dem Host des Manifests wiederverwendet, sofern keine Herabstufung von `https` auf `http` erfolgt. Bei jeder ursprungsübergreifenden oder herabstufenden Weiterleitung wird der Header `Authorization` entfernt.

Die CLI übernimmt die Veröffentlichung: `wails3 updater manifest` berechnet Prüfsummen, signiert und beschreibt Ihre Releasedateien mit einem einzigen Befehl; `wails3 updater verify` prüft das Ergebnis vor dem Hochladen erneut. Siehe [Veröffentlichen mit der wails3 CLI](/reference/update-manifest/#publishing-with-the-wails3-cli).

### Fallback-Kette

`Config.Providers` ist geordnet. Der Updater durchläuft die Einträge nacheinander: Der erste Provider, der eine Version zurückgibt, erhält den Zuschlag. Der erste Provider, der „aktuell“ meldet, beendet die Kette vorzeitig (der Fallback ist für den Fall „primärer Provider nicht erreichbar“ vorgesehen, nicht für „Provider widersprechen sich“). Bei einem Fehler fährt der Updater mit dem nächsten Provider fort.

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

Drei Methoden und bei einer typischen Implementierung etwa 150 Zeilen. Der Updater übernimmt Verifizierung, atomares Staging, Austausch und Fenster; der Provider-Code ermittelt das nächste Release und streamt die Bytes:

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

Verwenden Sie die im Repository enthaltenen Provider als Referenz – jeder besteht aus einer einzigen Go-Datei.

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
| `ecdsa-p256` | Der SHA-256-Digest des Artefakts | Sowohl rohe `r∥s`- als auch DER-Signaturen werden akzeptiert |

Zusätzlich ist eine reine Digest-Prüfung (`DigestAlgo`: `sha256` / `sha512`) möglich, wenn ein Release einen Hash, aber keine Signatur enthält.

`Config.PublicKey` ist der EINZIGE Vertrauensanker für die Signaturprüfung – die Release-Quelle kann ihn nicht durch einen eigenen Schlüssel ersetzen. Releases, die ein `Signature` enthalten, obwohl kein `Config.PublicKey` konfiguriert ist, werden sicher abgelehnt. Der Verifizierer berechnet den Digest während des Downloads in einem Streaming-Durchlauf. Daher ist selbst bei Updates von mehreren GB kein zusätzlicher Festplattendurchlauf für die Verifizierung erforderlich.

@note{type="caution" title="Reine Digest-Prüfung ≠ kryptografische Verifizierung"}
Ein Release, das nur `Digest` enthält, wird anhand des TLS der Registry und der von der Registry selbst gebotenen Integritätsgarantie authentifiziert – nicht anhand eines von Ihnen kontrollierten kryptografischen Vertrauensankers. Verwenden Sie die reine Digest-Prüfung zur Erkennung von Bit Rot und Signaturen zum Schutz vor Manipulationen bei einer kompromittierten Release-Pipeline.

@end

### Signaturschlüssel erzeugen

```bash
wails3 updater genkey
# updater.key       — keep secret, use to sign releases
# updater.key.pub   — bundle in your app via go:embed
```

Der private Schlüssel liegt als PKCS#8-PEM und der öffentliche Schlüssel als PKIX-PEM vor. `Config.PublicKey` akzeptiert die Datei `.pub` unverändert. Ebenfalls akzeptiert werden der rohe, 32 Byte lange Schlüssel oder dessen Base64-Darstellung, die `genkey` zur direkten Einbettung ausgibt. Signieren Sie Releases mit `wails3 updater manifest -key updater.key ...` oder `wails3 updater sign`; siehe [Veröffentlichen mit der wails3-CLI](/reference/update-manifest/#publishing-with-the-wails3-cli).

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
- **`.zip`** – wird direkt am Zielort entpackt. Das Archiv muss genau einen Eintrag auf oberster Ebene enthalten (in der Regel ein macOS-`.app`-Bundle oder eine einzelne Binärdatei). Empfohlenes Paketformat für macOS.
- **`.tar.gz`** / **`.tgz`** – wird nach derselben Regel eines einzigen Eintrags auf oberster Ebene direkt am Zielort entpackt. Nützlich für Linux-Distributionen, die neben der Binärdatei einen Runtime-Verzeichnisbaum ausliefern.

Archive mit mehr als einem Eintrag auf oberster Ebene werden abgelehnt: Das Framework tauscht ein einzelnes Ziel im Dateisystem aus. Daher ist „dieses Archiv am Zielort einsetzen“ mehrdeutig, wenn das Archiv mehrere Elemente enthält. `.dmg` und `.pkg` (macOS) sowie `.msi` (Windows) werden in v1 nicht unterstützt – verteilen Sie stattdessen ein `.zip` des Bundles. Beim Entpacken wird Schutz vor Zip Slip durchgesetzt, Symlinks, die aus dem Archivstamm herausführen, werden abgelehnt und sowohl die unkomprimierte Gesamtgröße (2 GiB) als auch die Anzahl der Einträge (50 000) werden begrenzt.

## Das Standardfenster

`app.Updater.CheckAndInstall(ctx)` öffnet ein vom Framework verwaltetes Fenster mit den Abmessungen 520 × 540 und folgenden Elementen:

- Ein zustandsabhängiges Hauptsymbol (blaues ↓ für verfügbar/wird heruntergeladen, grünes ✓ für bereit/aktuell, rotes ! für Fehler)
- Eine Versionsplakette: `v1.0.0 → v2.0.1 · 8.8 MB`
- Ein scrollbares Panel mit Versionshinweisen und **gerendertem Markdown** (Absätze, Fett-/Kursivschrift, Listen, GFM-Tabellen, Inline-Code, umschlossene Codeblöcke, h1–h3, Links)
- Eine einzelne primäre Aktion je Zustand (Installieren / Neu starten und anwenden / Erneut versuchen)
- Sekundäre Aktionen im Ghost-Stil (Diese Version überspringen / Später erinnern)
- Dunkler/heller Modus über `prefers-color-scheme`
- Unbestimmter schimmernder Fortschrittsindikator, wenn die Gesamtgröße unbekannt ist

Das Fenster lauscht auf dem Wails-Event-Bus auf `updater:*`-Ereignisse und sendet `updater:user:*`-Aktionen an Go zurück.

### Theme über CSS-Variablen

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --bg: #1a1a1a; --fg: #fafafa; }`,
    },
})
```

Das Standard-Stylesheet stellt die folgenden Variablen bereit – Sie können jede davon überschreiben:

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
| `--font` | Systemschriftarten mit Fallback-Reihenfolge | — |

### Vorlage ersetzen

Stellen Sie eigenes HTML bereit. Es muss lediglich auf `wails:updater:*`-Ereignisse lauschen und `wails:updater:user:*`-Aktionen auslösen:

```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        HTML: myCustomTemplate,
    },
})
```

InitialHTML-Fenster werden ohne Ursprung des Asset-Servers geladen und können daher `/wails/runtime.js` nicht dynamisch abrufen. Sie haben zwei Möglichkeiten, von einem solchen Fenster mit dem Host zu kommunizieren:

1. **Schreiben Sie einfach HTML.** Das Framework fügt automatisch einen minimalen `window.wails.Events`-Shim in jedes Fenster ein, das mit gesetzten Optionen `WebviewWindowOptions.AllowSimpleEventEmit = true` und `HTML` geöffnet wird – genau das tun die integrierten und BYO-Pfade des Updaters. Ein Build-Schritt ist nicht erforderlich. Das folgende Beispiel verwendet diesen Weg.
2. **Bündeln Sie `@wailsio/runtime` mit Ihrem bevorzugten Bundler** (Vite, esbuild, Rollup) und importieren Sie es zur Build-Zeit in Ihr benutzerdefiniertes HTML. `Events.On` funktioniert direkt, da es ausschließlich clientseitig arbeitet. `Events.Emit` verwendet hingegen den Fetch-Transport der Runtime, der durch den Null-Ursprung nicht funktioniert. Installieren Sie daher über den [`setTransport`](https://wails.io/wails/runtime.js)-Hook der Runtime einen kleinen postMessage-Transport, der über `window._wails.invoke("wails:event:emit:<name>")` weiterleitet. Die Framework-Injektion führt keine Aktion aus, wenn `window.wails.Events` bereits im Gültigkeitsbereich vorhanden ist, sodass sich die beiden Ansätze nicht gegenseitig beeinträchtigen.

In beiden Fällen ruft das JavaScript in Ihrem benutzerdefinierten HTML dieselbe `Events.On`-/`Events.Emit`-API auf:

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

Der Shim stellt die Teilmenge der modernen Runtime bereit, die Ereignisse mit einfachen Namen benötigen: `Events.On(name, cb)` gibt eine Funktion zum Abbestellen zurück, und `Events.Emit(nameOrEventObject)` leitet über den zugriffsbeschränkten postMessage-Pfad `wails:event:emit:` an den Host weiter. Er wird beim Laden der Seite einmalig installiert, bevor eines Ihrer Inline-Skripte ausgeführt wird.

Wenn Sie den Shim *bewusst* überschreiben möchten oder die vollständige Runtime auf anderem Weg laden, setzen Sie `window.wails.Events`, bevor das erste `<script>`-Tag der Seite ausgeführt wird. Die Injektion wird dann übersprungen.

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

Führen Sie den Aktualisierungsablauf mit einem selbst erstellten `*application.WebviewWindow` aus. Der Updater ruft `Show()` / `Close()` / `EmitEvent()` für Ihr Fenster auf – Ihr HTML bestimmt, was dargestellt wird:

![Ein benutzerdefiniertes Updater-Fenster mit einem pink-orangefarbenen Verlaufshintergrund, einer einzelnen weißen abgerundeten Karte, eigener Typografie und denselben Updater-Ereignissen zur Steuerung des sichtbaren Zustands. Es zeigt, wie vollständig sich die Standardbenutzeroberfläche ersetzen lässt.](/assets/updater/byo-custom-window.png)

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

Ihr HTML verwendet `window.wails.Events.On` / `Events.Emit` genau wie die integrierte Vorlage. Der automatisch injizierte Shim des Frameworks wird in jedes Fenster mit `AllowSimpleEventEmit: true` eingebunden, unabhängig davon, ob das Fenster dem Framework oder Ihnen gehört. Die API finden Sie unter [Vorlage ersetzen](#vorlage-ersetzen).

@note{type="caution" title="`AllowSimpleEventEmit` ist für benutzerdefinierte Updater-Fenster erforderlich"}
Aus Sicherheitsgründen lässt das Framework die `wails:event:emit:`-postMessage-Abkürzung nur zu, wenn dieses Feld gesetzt ist: Ein Fenster ohne diese Einstellung kann keine benutzerdefinierten Ereignisse auf der Host-Seite erzeugen. Der Shim für benutzerdefiniertes HTML des Updaters löst über diese Abkürzung `updater:user:*`-Ereignisse aus. Vergessen Sie das Feld bei einem eigenen Fenster, werden daher alle Schaltflächenklicks ohne Rückmeldung verworfen – der Benutzer klickt auf Installieren, und nichts geschieht.

Lassen Sie `AllowSimpleEventEmit` bei jedem Fenster **deaktiviert**, das HTML lädt, über das Sie nicht die vollständige Kontrolle haben (Remote-URLs, von Benutzern bereitgestellte Inhalte). Ist die Option aktiviert, kann jedes JavaScript auf der Seite – einschließlich XSS-Senken – jeden `app.Event.On(name, …)`-Handler auslösen. Die Abkürzung überträgt nur einfache Namen (keine Nutzdaten) und kann den Binding-/Call-Pfad nicht erreichen. Sie kann jedoch weiterhin privilegierte Handler für benutzerdefinierte Ereignisse in Ihrem Go-Code auslösen, wenn diese Handler ausschließlich anhand des Ereignisnamens handeln.

Beim *integrierten* Updater-Fenster des Frameworks ist diese Option intern gesetzt – nur bei eigenen Fenstern müssen Aufrufer daran denken.

@end

### Ohne Benutzeroberfläche

```go
app.Updater.Init(updater.Config{
    // …
    Window: updater.WindowNone,
})
```

Es wird zu keinem Zeitpunkt ein Fenster geöffnet. Abonnieren Sie `updater:*`-Ereignisse über Ihre eigene Benutzeroberfläche oder Ihr vorhandenes Hauptfenster und rufen Sie `app.Updater.CheckAndInstall(ctx)` in einem Schaltflächen-Handler auf. Dies eignet sich für regelmäßige Hintergrundprüfungen, die nur sichtbar werden sollen, wenn etwas gefunden wurde, oder für Apps, die den Aktualisierungsablauf in einen benutzerdefinierten Einstellungsbereich integrieren.

## Ereignisse

Sowohl Go als auch JavaScript abonnieren Ereignisse über den standardmäßigen Wails-Ereignisbus. **Geben Sie die Übertragungszeichenfolgen nicht von Hand ein** – verwenden Sie die Konstanten, die vom Updater-Paket (Go) oder Runtime-Paket (JS) exportiert werden. Beide Ebenen verwenden denselben Satz von Namen, der durch einen Regressionstest synchron gehalten wird.

### Aus Go

Die Konstanten befinden sich in `github.com/wailsapp/wails/v3/pkg/updater`. Abonnieren Sie über `app.Event.On(name, fn)`. Der Callback erhält ein `*application.CustomEvent`, dessen Feld `Data` die in der [Ereignisreferenz](#ereignisreferenz) aufgeführten typisierten Nutzdaten enthält. Verwenden Sie eine Typzusicherung und keine JSON-Dekodierung:

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

Die Konstanten befinden sich unter `Updater.Events` in `@wailsio/runtime`. Sie tragen dieselben Namen wie in Go und sind zur besseren Auffindbarkeit per Autovervollständigung nach Unter-Namensräumen (`User.*`, `Window.*`) gegliedert:

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

Ereignisse für Benutzeraktionen, die Ihr benutzerdefiniertes HTML *zurück* an den Host sendet, befinden sich unter `Updater.Events.User`:

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

| Konstante (Go) | Konstante (JS) | Nutzdaten | Zeitpunkt |
| --- | --- | --- | --- |
| `updater.EventCheckStarted` | `Updater.Events.CheckStarted` | keine | Vor jedem Roundtrip von `Check` |
| `updater.EventUpdateAvailable` | `Updater.Events.UpdateAvailable` | `*Release` | `Check` hat eine neuere Version gefunden |
| `updater.EventNoUpdate` | `Updater.Events.NoUpdate` | keine | `Check` hat bestätigt, dass die aktuelle Version installiert ist |
| `updater.EventDownloadStarted` | `Updater.Events.DownloadStarted` | `*Release` | Das Streaming der Bytes beginnt |
| `updater.EventDownloadProgress` | `Updater.Events.DownloadProgress` | `Progress` | Während des Downloads mit ~10 Hz |
| `updater.EventDownloadComplete` | `Updater.Events.DownloadComplete` | `*Release` | Alle Bytes wurden geschrieben, vor der Verifizierung |
| `updater.EventVerifying` | `Updater.Events.Verifying` | `*Release` | Die Prüfung der Signatur/des Hashwerts beginnt |
| `updater.EventInstalling` | `Updater.Events.Installing` | `*Release` | Entpacken und Staging beginnen |
| `updater.EventUpdateReady` | `Updater.Events.UpdateReady` | `*Release` | Neustart steht aus |
| `updater.EventError` | `Updater.Events.Error` | `ErrorInfo` | Eine beliebige Phase ist fehlgeschlagen |
| `updater.EventMeta` | `Updater.Events.Meta` | `Meta` | Einmal pro Sitzung vor der Snapshot-Wiedergabe |

Auf der Seite (Seite → Host) — Ihr Code abonniert diese Ereignisse, wenn Sie eine benutzerdefinierte Vorlage erstellen:

| Konstante (Go) | Konstante (JS) | Wann |
| --- | --- | --- |
| `updater.EventWindowReady` | `Updater.Events.Window.Ready` | Das Fenster ist vollständig geladen; der Host gibt den aktuellen Zustand erneut wieder |
| `updater.EventUserInstall` | `Updater.Events.User.Install` | Primäre Aktion im Zustand `available` |
| `updater.EventUserRestart` | `Updater.Events.User.Restart` | Primäre Aktion im Zustand `ready` |
| `updater.EventUserSkip` | `Updater.Events.User.Skip` | „Diese Version überspringen“ |
| `updater.EventUserRemind` | `Updater.Events.User.Remind` | „Später erinnern“ |
| `updater.EventUserCancel` | `Updater.Events.User.Cancel` | Schaltfläche „Schließen“ |

## API-Referenz

### `updater.Config`

| Feld | Typ | Hinweise |
| --- | --- | --- |
| `CurrentVersion` | `string` | **Erforderlich.** Dieselbe Zeichenfolge, mit der Sie Releases taggen (ohne das Präfix `v`) |
| `Providers` | `[]updater.Provider` | **Erforderlich.** Geordnete Fallback-Kette |
| `PublicKey` | `[]byte` | PEM- oder Rohbytes. Optional, aber signierte Releases werden ohne sie sicher abgelehnt |
| `CheckInterval` | `time.Duration` | Ein Wert ungleich null startet im Hintergrund eine Abfrageschleife, die `CheckAndInstall` aufruft |
| `Platform` | `string` | Überschreibt `runtime.GOOS` für die Auswahl des Assets |
| `Arch` | `string` | Überschreibt `runtime.GOARCH` für die Auswahl des Assets |
| `Channel` | `string` | Derzeit nur informativ; anbieterspezifische Kanalfilterung |
| `Window` | `updater.WindowOption` | `nil` (integrierte Standardwerte), `&BuiltinWindow{…}`, `BYOWindow(handle)` oder `WindowNone` |

### Methoden von `*updater.Updater`

| Signatur | Zweck |
| --- | --- |
| `Init(cfg Config) error` | Konfiguriert den Updater. Gibt beim zweiten Aufruf `ErrAlreadyConfigured` zurück |
| `State() State` | Aktuelle Lebenszyklusphase |
| `CurrentVersion() string` | Die an `Init` übergebene Version |
| `Check(ctx) (*Release, error)` | Durchläuft die Anbieterkette. `(rel, nil)` = gefunden, `(nil, nil)` = aktuell, `(nil, err)` = alle fehlgeschlagen |
| `DownloadAndInstall(ctx) error` | Streamt, verifiziert, extrahiert (bei einem Archiv) und stellt bereit. Erfordert zuvor `Check` |
| `CheckAndInstall(ctx) error` | Komfortfunktion: Öffnet das Fenster, führt `Check` aus und anschließend bei einem Fund `DownloadAndInstall` |
| `Restart(ctx) error` | Startet das Hilfsprogramm, ruft `Host.Quit` auf und beendet die Anwendung; das Hilfsprogramm tauscht die Dateien aus und startet die Anwendung neu |
| `DownloadedPath() string` | Speicherort des bereitgestellten Updates auf dem Datenträger oder `""`, wenn keines vorhanden ist |
| `SkipVersion(v string)` | Vermerkt `v` als übersprungen; nachfolgende `Check`s behandeln diese Version als aktuell |
| `SkippedVersion() string` | Liest die derzeit übersprungene Version |
| `StopPeriodicCheck()` | Beendet den von `Config.CheckInterval` gestarteten Timer und wartet, bis die Schleife zurückkehrt |

### Fehler

| Sentinelwert | Zurückgegeben von |
| --- | --- |
| `ErrAlreadyConfigured` | `Init` nach dem ersten Erfolg |
| `ErrNotConfigured` | Jede Operation vor `Init` |
| `ErrNoPendingRelease` | `DownloadAndInstall` ohne vorheriges `Check` |
| `ErrDownloadInProgress` | Aufruf von `DownloadAndInstall`, während bereits ein weiterer Aufruf derselben Methode läuft |
| `ErrNotReady` | `Restart` ohne bereitgestelltes Update |

## So funktioniert der Austausch

`Restart` führt die aktuelle Binärdatei mit gesetzten Sentinel-Umgebungsvariablen erneut aus. `application.New` erkennt sie beim Start und wechselt in den Hilfsmodus:

1. Das Hilfsprogramm wartet bis zu 30 s darauf, dass die übergeordnete PID beendet wird (`platformIsAlive` fragt unter Windows über `syscall.OpenProcess` + `GetExitCodeProcess` und unter Unix über `os.FindProcess` + `proc.Signal(syscall.Signal(0))` ab).
2. Das Hilfsprogramm sichert das Ziel (Kopie bei Dateien, rekursive Kopie bei macOS-`.app`-Bundle-Verzeichnissen).
3. Das Hilfsprogramm ersetzt das Ziel durch das bereitgestellte Artefakt. Dabei wiederholt es den Vorgang bei Fehlschlägen bis zu 20-mal und wartet zwischen den Versuchen jeweils 500 ms:
  - **Unix** — `os.RemoveAll(target)` + `os.Rename(newPath, target)`. Offene Dateideskriptoren, die auf den alten Inode verweisen, bleiben gültig.
  - **Windows** — `os.Rename(target, target.old.<nanos>)` + `os.Rename(newPath, target)`. Windows erlaubt das Umbenennen von Dateien, deren Image noch gemappt ist, aber nicht deren Löschen. Beim nächsten Update entfernt das Hilfsprogramm alle verbliebenen gleichgeordneten `.old.*`-Dateien, deren zugehöriges Kernel-Mapping freigegeben wurde.

4. Das Hilfsprogramm stellt für die neue Binärdatei die ursprünglichen Ausführungsrechte wieder her (die heruntergeladene Datei wurde mit der Standard-umask erstellt, wodurch `+x` unter Unix entfällt; unter Windows hat dies keine Wirkung).
5. Das Hilfsprogramm entfernt die Umgebungsvariablen des Hilfsmodus und startet die nun ersetzte Binärdatei erneut.
6. Das Hilfsprogramm wird beendet.

Schlägt der Start fehl, stellt das Hilfsprogramm die Sicherung wieder her. Wird der übergeordnete Prozess nicht innerhalb von 30 s beendet, bricht das Hilfsprogramm ab, bevor es das Ziel verändert. Dadurch behält der Benutzer selbst dann eine funktionsfähige App, wenn ein Dialog beim Beenden `Quit` blockiert.

Bei macOS-`.app`-Bundles, die als `.zip` verteilt werden (die empfohlene Paketierung), wird das Archiv zwischen der Verifizierung und der Bereitschaft entpackt, damit dem Hilfsprogramm ein echtes Verzeichnis für den Austausch zur Verfügung steht.

## Regelmäßige Prüfung

```go
app.Updater.Init(updater.Config{
    // …
    CheckInterval: 6 * time.Hour,
})
```

Bei `CheckInterval > 0` ruft eine Hintergrund-Goroutine `CheckAndInstall` im konfigurierten Intervall auf. Ticks, die eintreffen, während bereits ein anderer Ablauf läuft (Prüfen/Herunterladen/Verifizieren/Installieren), werden verworfen – parallele Zustandsautomaten werden nicht unterstützt.

Für eine stille Abfrage im Hintergrund, die nur bei einem Fund sichtbar wird, setzen Sie `Window: updater.WindowNone` und reagieren Sie in Ihrer eigenen Benutzeroberfläche auf `EventUpdateAvailable`.

## Überspringen und erinnern

Die Schaltfläche „Diese Version überspringen“ im Standardfenster speichert die verfügbare Version über `SkipVersion(rel.Version)`. Nachfolgende `Check`-Aufrufe finden dieselbe Version und behandeln sie als aktuell, bis der Benutzer `CurrentVersion` aktualisiert. Dies geschieht nach einem erfolgreichen `Restart` automatisch. „Später erinnern“ schließt lediglich das Fenster, ohne etwas zu speichern.

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

1. **Wählen Sie das richtige Archivformat.** macOS: `.zip` des `.app`-Bundles. Linux: einzelne Binärdatei oder `.tar.gz`. Windows: einzelne `.exe`- oder `.zip`-Datei. `.dmg` / `.msi` / `.pkg` werden nicht unterstützt.
2. **Signieren Sie das Artefakt** mit dem privaten Schlüssel, der zu `Config.PublicKey` gehört. Befolgen Sie bei von Anbietern veröffentlichten Feeds (keygen.sh, AppCast) den jeweiligen Signierungsablauf. Erzeugen Sie für GitHub Releases mit `ChecksumAsset` eine `SHA256SUMS`-Datei mit `sha256sum` / `shasum -a 256`.
3. **Gleichen Sie die Versionszeichenfolge ab.** `Config.CurrentVersion` und das Versions-Tag des Releases müssen exakt übereinstimmen (z. B. `1.0.0` ↔ Tag `v1.0.0`; das vorangestellte `v` wird auf Anbieterseite entfernt).
4. **Testen Sie den Austausch auf der Zielplattform** mindestens einmal vor der Auslieferung – Codesignierung, Notarisierung und die Handhabung durch Gatekeeper sind plattformspezifisch und werden vom Updater selbst nicht abgedeckt.

## Fehlerbehebung

**„signature requires a public key but none configured“** — das Release enthält ein `Signature`-Feld, aber `Config.PublicKey` ist leer. Legen Sie den öffentlichen Schlüssel fest oder ändern Sie Ihre Release-Pipeline so, dass sie keine Signatur einfügt.

**„digest mismatch“** — die heruntergeladenen Bytes entsprechen nicht den Angaben des Anbieters. Ursache ist meist ein unvollständiger Download aufgrund einer kurzen Netzwerkstörung oder ein beschädigtes Artefakt. Ein erneuter Versuch behebt das Problem häufig.

**Fenster öffnet sich, verschwindet aber sofort; kein Markdown, kein Fortschritt** — Ihr benutzerdefiniertes HTML hat `wails:runtime:ready` nicht aufgerufen. Siehe den Shim unter [Vorlage ersetzen](#vorlage-ersetzen).

**Windows-Update wird nie abgeschlossen; im Hilfsprogrammprotokoll steht „remove old (attempt N): Access is denied“** — dies tritt nur bei Versionen dieses PRs vor `de764fb` auf. Die aktuelle Implementierung benennt die Datei vorübergehend um und ist von diesem Problem nicht betroffen. Führen Sie ein Upgrade durch.

**macOS Gatekeeper blockiert die ausgetauschte Binärdatei** — die Codesignierung muss durchgängig erhalten bleiben. Signieren Sie das ursprüngliche `.app` *und* signieren Sie die neu gestartete Binärdatei erneut, falls Ihre Build-Pipeline die Berechtigungen während des Updates ändert.

## Siehe auch

- Ausführbares Beispiel: [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater)
- Demo-Repository für Tests: [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo)
- Tutorial: [Selbstaktualisierungen zu einer Wails-App hinzufügen](/tutorials/04-self-update-a-wails-app/)
