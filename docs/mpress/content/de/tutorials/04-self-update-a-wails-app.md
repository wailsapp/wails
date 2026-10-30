---
title: "Selbstaktualisierende Wails-App"
description: "Erstellen Sie eine Wails-v3-Anwendung, die sich über GitHub Releases selbst aktualisiert – von `wails3 init` über die Verifizierung signierter Releases bis zum Austausch im Hilfsmodus."
slug: "tutorials/04-self-update-a-wails-app"
sourcePath: "tutorials/04-self-update-a-wails-app.md"
---

In diesem Tutorial fügen Sie einer neuen Wails-v3-Anwendung einen integrierten Updater hinzu. Anschließend kann die App:

- GitHub Releases bei Bedarf und optional zeitgesteuert auf Updates prüfen.
- Das passende Asset für das verwendete Betriebssystem und die Architektur herunterladen.
- Einen SHA-256-Digest und optional eine Ed25519-Signatur anhand der heruntergeladenen Bytes verifizieren.
- Versionshinweise im standardmäßigen Updatefenster des Frameworks anzeigen.
- Die laufende Binärdatei austauschen und die App neu starten – ohne eine separate Hilfsprogrammdatei auszuliefern.

Wir verwenden **GitHub Releases** als Updatequelle, da der Dienst kostenlos ist und keine Infrastruktur erfordert. Dieselben Muster funktionieren mit [keygen.sh](/guides/updater/#keygensh--updaterproviderskeygen) und [Sparkle AppCast](/guides/updater/#sparkle-appcast--updaterprovidersappcast). Lesen Sie nach Abschluss den [Updater-Leitfaden](/guides/updater/).

@note{type="tip" title="Voraussetzungen"}
- Go 1.26 oder neuer
- Installierte `wails3`-CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`)
- Ein GitHub-Repository, in das Sie Releases übertragen können
- Kenntnisse aus dem Tutorial zum [QR-Code-Dienst](/tutorials/01-creating-a-service/) sind hilfreich, aber nicht erforderlich

@end

<br/>

@steps
### Mit einer neuen Wails-App beginnen
Erstellen Sie mit der Vanilla-Vorlage ein neues Projektgerüst:

```bash
wails3 init -n updater-tutorial -t vanilla
cd updater-tutorial
```

Sie sollten nun über ein Verzeichnis mit `main.go`, `frontend/` und einer `Taskfile.yml` verfügen. Vergewissern Sie sich, dass sich das Projekt erstellen und starten lässt:

```bash
wails3 task dev
```

Ein leeres Wails-Fenster sollte geöffnet werden. Beenden Sie die App und fahren Sie fort.

### Updater-Importe hinzufügen
Öffnen Sie `main.go` und fügen Sie den Importen die beiden Updater-Pakete hinzu:

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
`app.Updater` ist bereits mit jeder `*application.App` verbunden. Sie müssen lediglich `Init` aufrufen:

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

Fügen Sie dies nach `application.New` und vor `app.Run()` ein.

@note{type="note" title="Format der Versionszeichenfolge"}
Übergeben Sie dieselbe Version, die Sie für Release-Tags verwenden, jedoch **ohne** das vorangestellte `v`. Der Anbieter entfernt `v` selbst aus den Tag-Namen. `1.0.0` hier ↔ `v1.0.0` auf GitHub.

@end

### Menüeintrag zum Auslösen des Updates hinzufügen
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

`CheckAndInstall` öffnet das Updatefenster des Frameworks, führt `Check` aus und führt automatisch `DownloadAndInstall` aus, falls ein Release gefunden wird. Wenn nichts Neues verfügbar ist, bleibt das Fenster im Zustand „Auf dem neuesten Stand“ geöffnet. Der Benutzer schließt es über die Schaltfläche **Schließen**.

@note{type="caution" title="In einer Goroutine ausführen"}
`CheckAndInstall` blockiert, bis Verifizierung und Installation abgeschlossen sind. Ein direkter Aufruf über den Menüklick würde den UI-Thread blockieren. Umschließen Sie den Aufruf mit `go func()`.

@end

### Einmal ohne vorhandene Releases ausführen
```bash
wails3 task dev
```

Klicken Sie auf **App → Nach Updates suchen…**. Das Updatefenster sollte kurz geöffnet werden, die GitHub-API aufrufen, keine Releases finden, die neuer als `1.0.0` sind, und anschließend mit einem grünen ✓ den Zustand **Auf dem neuesten Stand** anzeigen.

Wenn hier ein Fehler auftritt, liegt meist eine der folgenden Ursachen vor:

| Symptom | Lösung |
| --- | --- |
| `404 Not Found` | Das Feld `Repository` ist falsch – es muss `owner/repo` lauten |
| `403 rate-limited` | Fügen Sie der github.Config `Token: "ghp_…"` hinzu. Verwenden Sie ein PAT mit dem Geltungsbereich `public_repo`. |
| Netzwerkfehler | Vergewissern Sie sich, dass die laufende App `api.github.com` erreichen kann |

### Test-Release veröffentlichen
Erhöhen Sie `currentVersion` in `main.go` auf `1.0.0` oder lassen Sie den Wert unverändert. Erstellen Sie die App für eine Plattform, um eine Binärdatei zu erhalten, die Sie an ein Release anhängen können:

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

Sie sollten mindestens eine Zeile wie die folgende sehen:

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
Der standardmäßige Asset-Abgleich wählt anhand der Teilzeichenfolgen `GOOS` und `GOARCH` im Dateinamen aus. Solange der Asset-Name `darwin` (oder `linux`/`windows`) und `arm64` (oder `amd64`/`386`) enthält, wird das Asset gefunden. Informationen zu benutzerdefinierten Abgleichsfunktionen finden Sie im [Updater-Leitfaden](/guides/updater/#github-releases--updaterprovidersgithub).

@end

### App ausführen und Update überprüfen
Führen Sie die App erneut aus, während `currentVersion` weiterhin `1.0.0` ist:

```bash
wails3 task dev
```

Klicken Sie auf **App → Nach Updates suchen…**. Dieses Mal sollte ungefähr Folgendes angezeigt werden:

![Standardmäßiges Updater-Fenster im Zustand „Update bereit“ mit Versionskennzeichnung, als Markdown gerenderten Versionshinweisen und der primären Schaltfläche „Neu starten und anwenden“.](/assets/updater/default-window-ready.png)

- Das Hero-Symbol wechselt von einem blauen ↓ („Update verfügbar“) zu einem grünen ✓ („Update bereit“).
- Der Untertitel zeigt `v1.0.0 → v2.0.0 · <size>` an.
- Der Bereich für Versionshinweise rendert Ihr Markdown mit Fettschrift, Code-Spannen und der Tabelle.
- Während des Downloads füllt sich die Fortschrittsanzeige. Das geht schnell, da die Binärdatei klein ist.

Der Updater stellt die neue Binärdatei in einem temporären Verzeichnis bereit. So schließen Sie das Update ab:

- Klicken Sie auf **Neu starten und anwenden**.
- Ihre App wird beendet, das Hilfsprogramm tauscht die Binärdatei aus und die neue Binärdatei wird gestartet.
- Die neu gestartete App meldet `currentVersion = "1.0.0"`, weil wir diesen Wert fest codiert haben. Die Bytes auf dem Datenträger entsprechen jedoch dem Build v2.0.0.

In einer echten App würde `currentVersion` zur Build-Zeit über `-ldflags` festgelegt. Dadurch weiß die neue Binärdatei, dass sie nun v2.0.0 ist, und eine nachfolgende Prüfung findet kein Update.

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

### Kryptografische Signierung hinzufügen (für die Produktion empfohlen)
Der SHA256SUMS-Pfad überprüft die *Integrität* – die Bytes entsprechen den auf GitHub gespeicherten Daten –, nicht jedoch die *Authentizität*: Die Bytes müssen aus Ihrer Release-Pipeline stammen und dürfen nicht von einem kompromittierten Maintainer-Konto erzeugt worden sein. Signieren Sie zum Schutz vor Manipulation jedes Release mit einem Ed25519-Schlüssel:

```bash
# One-time: generate the keypair
ssh-keygen -t ed25519 -f updater-key -N "" -C "wails-updater"
#   updater-key      — keep secret (build server, HSM, password manager)
#   updater-key.pub  — bundle in your app
```

Signieren Sie für jedes Release den SHA-256-Hashwert jedes Assets mit Ihrem privaten Schlüssel. Ein kleines Go-Hilfsprogramm:

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

Der standardmäßige GitHub-Provider ruft derzeit keine separate Signaturdatei ab. Sie können einen [benutzerdefinierten Provider schreiben](/guides/updater/#writing-your-own-provider), der dies übernimmt, oder zu **keygen.sh** wechseln. Dieser Dienst signiert jedes Artefakt serverseitig und stellt sowohl den Hashwert als auch die Signatur über seine API bereit.

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

Wenn `PublicKey` festgelegt ist, muss jedes Release, das eine `Signature` enthält, erfolgreich anhand dieses Schlüssels verifiziert werden. Die Release-Quelle kann keinen eigenen Schlüssel einsetzen – genau das ist der Zweck der Out-of-Band-Fixierung zur Build-Zeit.

### Fenster anpassen
Das Standardfenster deckt den üblichen Anwendungsfall ab. Wenn Sie mehr Kontrolle benötigen, stehen Ihnen drei Erweiterungsmöglichkeiten zur Verfügung. Wählen Sie je nach gewünschtem Anpassungsumfang eine davon aus:

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

Die vollständige Variablenliste finden Sie im Abschnitt [Design über CSS-Variablen](/guides/updater/#theme-via-css-variables).

[Benutzerdefiniertes HTML]
```go
//go:embed updater-window.html
var updaterHTML string

app.Updater.Init(updater.Config{
    // …
    Window: &updater.BuiltinWindow{HTML: updaterHTML},
})
```

Ihr HTML muss `updater:*`-Ereignisse abonnieren und über den Wails-Ereigniskanal `updater:user:*`-Aktionen auslösen. Den JS-Shim finden Sie unter [Vorlage ersetzen](/guides/updater/#replace-the-template).

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

Dies ist hilfreich, wenn Sie bereits über eine eigene Fensterinfrastruktur verfügen und der Updater diese steuern soll, statt ein weiteres Fenster zu öffnen. Eine vollständig benutzerdefinierte HTML-Vorlage, die von denselben Updater-Ereignissen wie die Standardvorlage gesteuert wird, sieht so aus:

![Ein selbst bereitgestelltes Updater-Fenster mit einem pink-orangefarbenen Verlaufshintergrund und einem benutzerdefinierten Kartenlayout mit abgerundeten Ecken, das zeigt, dass die Standardoberfläche vollständig ersetzt werden kann.](/assets/updater/byo-custom-window.png)

@note{type="caution" title="`AllowSimpleEventEmit` ist erforderlich"}
Der benutzerdefinierte HTML-Shim des Updaters steuert Installieren / Überspringen / Erinnern / Neustarten über die `wails:event:emit:`-postMessage-Abkürzung. Aus Sicherheitsgründen ist diese Abkürzung nur verfügbar, wenn dieses Feld aktiviert ist. Wenn Sie es vergessen, bleiben die Schaltflächen ohne sichtbare Wirkung. Aktivieren Sie es nicht für Fenster, die HTML laden, das Sie nicht vollständig kontrollieren. Das Bedrohungsmodell finden Sie im Abschnitt [Eigenes Fenster verwenden](/guides/updater/#bring-your-own-window) des Leitfadens.

@end

@end

### Automatische Prüfungen im Hintergrund ausführen
So prüfen Sie zeitgesteuert statt per Menüklick oder zusätzlich dazu:

```go {ins="5"}
app.Updater.Init(updater.Config{
    CurrentVersion: currentVersion,
    Providers:      []updater.Provider{gh},
    PublicKey:      updaterPublicKey,
    CheckInterval:  6 * time.Hour,
})
```

Bei jedem Intervall wird derselbe `CheckAndInstall`-Ablauf wie bei einem manuellen Klick ausgeführt. Legen Sie `Window: updater.WindowNone` fest, wenn die regelmäßige Prüfung keine Meldung anzeigen soll, bis tatsächlich etwas gefunden wurde. Abonnieren Sie dann selbst `EventUpdateAvailable`, um zu entscheiden, welche Benutzeroberfläche angezeigt werden soll.

@end

## Fertig

Sie haben jetzt eine Wails-App, die:

- GitHub Releases bei Bedarf und zeitgesteuert auf Updates prüft.
- Release Notes als Markdown in einem professionell gestalteten Standardfenster darstellt.
- Downloads anhand eines von Ihnen veröffentlichten SHA-256-Hashwerts verifiziert.
- Optional eine Ed25519-Signatur anhand eines öffentlichen Schlüssels verifiziert, den Sie zur Build-Zeit einbetten.
- Die laufende Binärdatei direkt ersetzt und die App automatisch neu startet.

## Nächste Schritte

- Der [Updater-Leitfaden](/guides/updater/) enthält die vollständige API-Referenz, sämtliche Ereignisse und Konfigurationsoptionen sowie die Austauschmechanik des Hilfsmodus.
- Ein vollständiges, funktionsfähiges Beispiel zum Klonen finden Sie unter [`v3/examples/updater`](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/updater).
- Das Test-Zielrepository [`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo) zeigt die empfohlene Struktur der Release-Assets.

## Wichtige Punkte für den Produktivbetrieb

- **Codesignierung unter macOS** – Gatekeeper verlangt, dass die ausgetauschte Binärdatei signiert und notarisiert ist. Signieren Sie Ihr `.app`-Bundle, *bevor* Sie es für das Release als ZIP-Datei verpacken. Der Updater übernimmt die Bytes unverändert und signiert nichts erneut.
- **Virenschutz unter Windows** – unsignierte `.exe`-Dateien, die aus dem Internet heruntergeladen wurden, können SmartScreen-Warnungen auslösen. Signieren Sie Ihre Binärdatei mit einem Authenticode-Zertifikat. Andernfalls müssen Sie damit rechnen, dass Benutzer auf stark eingeschränkten Systemen Ihre App möglicherweise auf die Positivliste setzen müssen.
- **Atomare Releases** – veröffentlichen Sie `SHA256SUMS` und Ihre Binärdateien gemeinsam, nicht in separaten Commits. Der Updater ruft die Sidecar-Datei getrennt von der Binärdatei ab. Wenn beide voneinander abweichen, schlägt die Hashwertprüfung sicher fehl.
- **Übersprungene Versionen** – die Schaltfläche „Diese Version überspringen“ im Standardfenster speichert die übersprungene Version lokal. Vergeben Sie für ein kritisches Sicherheitsupdate eine neue Versionsnummer, damit es bei Benutzern, die ein früheres Release verworfen haben, nicht automatisch übersprungen wird.
