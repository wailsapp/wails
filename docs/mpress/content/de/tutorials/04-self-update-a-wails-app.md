---
title: "Wails-App mit automatischer Aktualisierung"
description: "Erstellen Sie eine Wails-v3-Anwendung, die sich über GitHub Releases selbst aktualisiert – von `wails3 init` über die Verifizierung signierter Releases bis zum Austausch im Hilfsmodus."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

In diesem Tutorial fügen Sie einer neuen Wails-v3-Anwendung einen integrierten Updater hinzu. Am Ende kann die App:

- GitHub Releases bei Bedarf und optional regelmäßig prüfen.
- Das passende Asset für das verwendete Betriebssystem und die verwendete Architektur herunterladen.
- Einen SHA-256-Hash und optional eine Ed25519-Signatur anhand der heruntergeladenen Bytes verifizieren.
- Versionshinweise im standardmäßigen Aktualisierungsfenster des Frameworks anzeigen.
- Die ausgeführte Binärdatei austauschen und die App neu starten – ohne eine separate Hilfsprogrammdatei auszuliefern.

Wir verwenden **GitHub Releases** als Aktualisierungsquelle, da der Dienst kostenlos ist und keine eigene Infrastruktur erfordert. Dieselben Verfahren funktionieren mit [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) und [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast). Weitere Informationen finden Sie nach Abschluss im [Updater-Leitfaden](/guides/updater/).

@note{type="tip" title="Voraussetzungen"}
- Go 1.25 oder neuer
- Installierte `wails3`-CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- Ein GitHub-Repository, in das Sie Releases übertragen können
- Kenntnisse aus dem [Tutorial zum QR-Code-Dienst](/tutorials/01-creating-a-service/) sind hilfreich, aber nicht erforderlich

@end

<br/>

@steps
### Mit einer neuen Wails-App beginnen
Erstellen Sie mit der Vanilla-Vorlage das Grundgerüst eines neuen Projekts:

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

Sie sollten nun ein Verzeichnis mit `main.go`, `frontend/` und einer `Taskfile.yml` haben. Prüfen Sie, ob sich das Projekt erstellen und starten lässt:

```bash
wails3 task dev
```

Ein leeres Wails-Fenster sollte geöffnet werden. Beenden Sie die App und fahren Sie fort.

### Updater-Importe hinzufügen
Öffnen Sie `main.go` und fügen Sie Ihren Importen die beiden Updater-Pakete hinzu:

```go {title="main.go" ins="6-7"}
package main

import (
    _ "embed"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)
```

Damit werden der Updater selbst und der Anbieter für GitHub Releases eingebunden.

### Updater konfigurieren
`app.Updater` ist bereits in jede `*application.App` eingebunden. Sie müssen lediglich `Init` aufrufen:

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

Platzieren Sie dies nach `application.New` und vor `app.Run()`.

@note{type="note" title="Format der Versionszeichenfolge"}
Übergeben Sie dieselbe Version, mit der Sie Releases taggen, jedoch **ohne** das vorangestellte `v`. Der Anbieter entfernt `v` seinerseits aus den Tag-Namen. `1.0.0` hier ↔ `v1.0.0` auf GitHub.

@end

### Menüeintrag zum Auslösen der Aktualisierung hinzufügen
Fügen Sie in derselben `main.go` einen Menüeintrag „Nach Updates suchen…“ hinzu:

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

`CheckAndInstall` öffnet das Aktualisierungsfenster des Frameworks, führt `Check` aus und führt bei einem gefundenen Release automatisch `DownloadAndInstall` aus. Wenn nichts Neues verfügbar ist, bleibt das Fenster im Zustand „Auf dem neuesten Stand“ geöffnet. Der Benutzer schließt es mit der Schaltfläche **Schließen**.

@note{type="caution" title="In einer Goroutine ausführen"}
`CheckAndInstall` blockiert, bis Verifizierung und Installation abgeschlossen sind. Ein direkter Aufruf beim Anklicken des Menüeintrags würde den UI-Thread blockieren. Kapseln Sie den Aufruf in `go func()`.

@end

### Einmal ohne Releases ausführen
```bash
wails3 task dev
```

Klicken Sie auf **App → Nach Updates suchen…**. Das Aktualisierungsfenster sollte kurz geöffnet werden, die GitHub-API abfragen, keine neueren Releases als `1.0.0` finden und anschließend den Zustand **Auf dem neuesten Stand** mit einem grünen ✓ anzeigen.

Wenn hier ein Fehler auftritt, handelt es sich normalerweise um eines der folgenden Probleme:

| Symptom | Lösung |
| --- | --- |
| `404 Not Found` | Das Feld `Repository` ist falsch – es muss `owner/repo` lauten |
| `403 rate-limited` | Fügen Sie der github.Config `Token: "ghp_…"` hinzu. Verwenden Sie ein PAT mit dem Berechtigungsumfang `public_repo`. |
| Netzwerkfehler | Prüfen Sie, ob die ausgeführte App `api.github.com` erreichen kann |

### Test-Release veröffentlichen
Erhöhen Sie `currentVersion` in `main.go` auf `1.0.0` oder belassen Sie den Wert. Erstellen Sie die App für eine Plattform, um eine Binärdatei zu erhalten, die Sie an ein Release anhängen können:

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

Erzeugen Sie neben der Binärdatei eine `SHA256SUMS`-Datei:

```bash
cd bin
shasum -a 256 updater-tutorial-* > SHA256SUMS
cat SHA256SUMS
```

Sie sollten eine oder mehrere Zeilen wie diese sehen:

```
abc123…  updater-tutorial-darwin-arm64.zip
```

Veröffentlichen Sie dies nun in Ihrem GitHub-Repository als **v2.0.0**:

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

@note{type="note" title="Benennung der Assets"}
Der standardmäßige Asset-Abgleich wählt anhand der Teilzeichenfolgen `GOOS` und `GOARCH` im Dateinamen aus. Solange der Asset-Name `darwin` oder `linux` / `windows` sowie `arm64` oder `amd64` / `386` enthält, wird das Asset gefunden. Informationen zu benutzerdefinierten Abgleichsfunktionen finden Sie im [Updater-Leitfaden](/guides/updater/#github-releases--updaterprovidersgithub).

@end

### App ausführen und Aktualisierung überprüfen
Führen Sie die App erneut aus, während `currentVersion` weiterhin auf `1.0.0` gesetzt ist:

```bash
wails3 task dev
```

Klicken Sie auf **App → Nach Updates suchen…**. Dieses Mal sollten Sie ungefähr Folgendes sehen:

![Standardmäßiges Updater-Fenster im Zustand „Update bereit“ mit Versionsplakette, als Markdown gerenderten Versionshinweisen und der primären Schaltfläche „Neu starten und anwenden“.](/assets/updater/default-window-ready.png)

- Das Hauptsymbol wechselt vom blauen ↓ („Update verfügbar“) zum grünen ✓ („Update bereit“).
- Der Untertitel zeigt `v1.0.0 → v2.0.0 · <size>` an.
- Der Bereich für Versionshinweise rendert Ihr Markdown einschließlich Fettschrift, Code-Spannen und Tabelle.
- Während des Downloads füllt sich der Fortschrittsbalken. Das geht schnell, da die Binärdatei klein ist.

Der Updater stellt die neue Binärdatei in einem temporären Verzeichnis bereit. So schließen Sie die Aktualisierung ab:

- Klicken Sie auf **Neu starten und anwenden**.
- Ihre App wird beendet, das Hilfsprogramm tauscht die Binärdatei aus und die neue Binärdatei wird gestartet.
- Die neu gestartete App meldet `currentVersion = "1.0.0"`, da wir diesen Wert fest codiert haben. Die Bytes auf dem Datenträger entsprechen jedoch dem Build v2.0.0.

In einer echten App würde `currentVersion` zur Build-Zeit über `-ldflags` festgelegt. Dadurch weiß die neue Binärdatei, dass sie nun Version v2.0.0 hat, und bei einer anschließenden Prüfung wird keine Aktualisierung gefunden.

### `currentVersion` mit dem Build verknüpfen
Ersetzen Sie die Konstante durch eine zur Build-Zeit gesetzte Variable:

```go {title="main.go" ins="2,4"}
var (
    currentVersion = "dev" // overridden by -ldflags at release time
)
```

Verwenden Sie anschließend in Ihrem Build-Befehl:

```bash
wails3 task build:darwin -- -ldflags "-X main.currentVersion=2.0.0"
```

Alternativ können Sie `-ldflags` zu Ihrer `Taskfile.yml` hinzufügen, damit der Wert aus `git describe --tags` übernommen wird.

### Kryptografische Signierung hinzufügen (für den Produktionseinsatz empfohlen)
Der SHA256SUMS-Pfad verifiziert die *Integrität* – die Bytes entsprechen den auf GitHub gespeicherten Bytes –, aber nicht die *Authentizität*: dass diese Bytes von Ihrer Release-Pipeline und nicht über ein kompromittiertes Maintainer-Konto erzeugt wurden. Um Manipulationen zu erschweren, signieren Sie jedes Release mit einem Ed25519-Schlüssel:

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

Signieren Sie für jedes Release den SHA-256-Digest jedes Assets mit Ihrem privaten Schlüssel. Ein kleines Go-Hilfsprogramm:

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

Der standardmäßige GitHub-Provider ruft derzeit keine separate Signaturdatei ab. Sie können einen [benutzerdefinierten Provider schreiben](/guides/updater/#writing-your-own-provider), der dies übernimmt, oder zu **keygen.sh** wechseln. Dieser Dienst signiert jedes Artefakt serverseitig und stellt sowohl den Digest als auch die Signatur über seine API bereit.

Betten Sie den öffentlichen Schlüssel in Ihre App ein:

```go {title="main.go" ins="1,6"}
//go:embed updater-key.pub
var updaterPublicKey []byte

app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    PublicKey:      updaterPublicKey,
    Providers:      []updater.Provider{gh},
})
```

Wenn `PublicKey` festgelegt ist, muss jedes Release, das eine `Signature` enthält, erfolgreich mit diesem Schlüssel verifiziert werden. Die Release-Quelle kann keinen eigenen Schlüssel einschleusen – genau das ist der Zweck der externen Schlüsselbindung zur Build-Zeit.

### Fenster anpassen
Das Standardfenster deckt den üblichen Anwendungsfall ab. Wenn Sie mehr Kontrolle benötigen, stehen Ihnen drei Ausweichmöglichkeiten zur Verfügung. Wählen Sie je nach gewünschtem Anpassungsumfang eine davon:

@tabs
[Nur CSS]
```go
app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{
        CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
    },
})
```

Die vollständige Variablenliste finden Sie im Abschnitt [Theme über CSS-Variablen](/guides/updater/#theme-via-css-variables).

[Benutzerdefiniertes HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

Ihr HTML muss `updater:*`-Ereignisse abonnieren und `updater:user:*`-Aktionen über den Wails-Ereigniskanal auslösen. Den JS-Shim finden Sie unter [Vorlage ersetzen](/guides/updater/#replace-the-template).

[Eigenes Fenster verwenden]
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

Dies ist nützlich, wenn Sie bereits über eine eigene Fensterinfrastruktur verfügen und möchten, dass der Updater diese steuert, anstatt ein weiteres Fenster zu öffnen. Eine vollständig benutzerdefinierte HTML-Vorlage, die von denselben Updater-Ereignissen wie die Standardvorlage gesteuert wird, sieht so aus:

![Ein eigenes Updater-Fenster mit einem pink-orangefarbenen Verlaufshintergrund und einem benutzerdefinierten Kartenlayout mit abgerundeten Ecken, das zeigt, dass die Standardoberfläche vollständig ersetzt werden kann.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` ist erforderlich"}
Der benutzerdefinierte HTML-Shim des Updaters steuert Installieren / Überspringen / Erinnern / Neustarten über die `wails:event:emit:`-postMessage-Abkürzung. Aus Sicherheitsgründen ist diese Abkürzung nur verfügbar, wenn dieses Feld aktiviert ist. Wenn Sie dies vergessen, bleiben die Schaltflächen ohne sichtbare Reaktion. Aktivieren Sie das Feld nicht für Fenster, die HTML laden, über das Sie nicht die vollständige Kontrolle haben. Das Bedrohungsmodell finden Sie im Abschnitt [Eigenes Fenster verwenden](/guides/updater/#bring-your-own-window) des Leitfadens.

@end

@end

### Automatische Prüfungen im Hintergrund ausführen
So führen Sie die Prüfung zeitgesteuert statt durch einen Menüklick oder zusätzlich dazu aus:

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

Jeder Zeitgeberdurchlauf führt denselben `CheckAndInstall`-Ablauf wie ein manueller Klick aus. Legen Sie `Window: updater.WindowNone` fest, wenn die regelmäßige Prüfung im Hintergrund bleiben soll, bis tatsächlich etwas gefunden wird. Abonnieren Sie dann selbst `EventUpdateAvailable`, um zu entscheiden, welche Benutzeroberfläche angezeigt werden soll.

@end

## Fertig

Sie haben jetzt eine Wails-App, die:

- GitHub Releases bei Bedarf und zeitgesteuert auf Updates prüft.
- Release Notes als Markdown in einem ansprechenden Standardfenster darstellt.
- Downloads anhand eines von Ihnen veröffentlichten SHA-256-Digests verifiziert.
- Optional eine Ed25519-Signatur anhand eines öffentlichen Schlüssels verifiziert, den Sie zur Build-Zeit einbetten.
- Die laufende Binärdatei direkt ersetzt und die App automatisch neu startet.

## Nächste Schritte

- Der [Updater-Leitfaden](/guides/updater/) enthält die vollständige API-Referenz, sämtliche Ereignisse und Konfigurationsoptionen sowie die Austauschmechanik des Hilfsmodus.
- Ein vollständiges, funktionsfähiges Beispiel zum Klonen finden Sie unter [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/master/v3/examples/updater).
- Das Testziel-Repository [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo) zeigt die empfohlene Anordnung der Release-Assets.

## Wichtige Hinweise für den Produktivbetrieb

- **Codesignierung unter macOS** – Gatekeeper setzt voraus, dass die ausgetauschte Binärdatei signiert und notarisiert ist. Signieren Sie Ihr `.app`-Bundle, *bevor* Sie es für das Release in eine ZIP-Datei packen. Der Updater behält die Bytes unverändert bei und signiert nichts erneut.
- **Antivirenprogramme unter Windows** – unsignierte, aus dem Internet heruntergeladene `.exe`-Dateien können SmartScreen-Warnungen auslösen. Signieren Sie Ihre Binärdatei mit einem Authenticode-Zertifikat, oder nehmen Sie in Kauf, dass Benutzer auf stark eingeschränkten Rechnern Ihre App möglicherweise auf die Positivliste setzen müssen.
- **Atomare Releases** – veröffentlichen Sie `SHA256SUMS` (und Ihre Binärdateien) gemeinsam, nicht in separaten Commits. Der Updater ruft die Sidecar-Datei getrennt von der Binärdatei ab. Weichen sie voneinander ab, schlägt die Digest-Prüfung fehl und blockiert die Installation.
- **Übersprungene Versionen** – die Schaltfläche „Diese Version überspringen“ im Standardfenster speichert die übersprungene Version lokal. Wenn Sie ein kritisches Sicherheitsupdate veröffentlichen, vergeben Sie eine neue Versionsnummer, damit es bei Benutzern, die ein früheres Release verworfen haben, nicht automatisch übersprungen wird.
