---
title: "Asset-Server"
description: "Wie Wails v3 Ihre Web-Assets in der Entwicklung und Produktion bereitstellt und einbettet"
slug: "contributing/asset-server"
sourcePath: "contributing/asset-server.md"
---

## Übersicht

Jede Wails-Anwendung wird als **einzelne native ausführbare Datei** ausgeliefert, die Folgendes kombiniert:

1. Ihr *Go*-Backend
2. Ein *Web*-Frontend (HTML + JS + CSS)

Der **Asset-Server** verbindet diese Komponenten. Er verfügt über **zwei Betriebsmodi**, die zur Kompilierzeit über Go-Build-Tags ausgewählt werden:

| Modus | Tag | Zweck |
| --- | --- | --- |
| **Entwicklung** | `//go:build !production` | Schnelle Iteration mit Hot Reload |
| **Produktion** | `//go:build production` | Eingebettete Assets ohne externe Abhängigkeiten |

Die Implementierung befindet sich in `v3/internal/assetserver/` und ist übersichtlich auf mehrere Dateien aufgeteilt:

```
build_dev.go              # ⬅️ dev-only entrypoint (!production build tag)
build_production.go       # ⬅️ production-only entrypoint (production build tag)
assetserver.go            # Shared core
assetserver_dev.go        # Dev proxy/disk handler
assetserver_webview.go    # WebView-side adapter
assetserver_darwin.go     # OS-specific helpers (also linux/windows variants)
asset_fileserver.go       # Shared static file logic
content_type_sniffer.go   # MIME type detection
mimecache.go              # Cached MIME lookups
ringqueue.go              # Tiny in-memory LRU
options.go                # Configuration struct
middleware.go             # http.Handler middleware type
bundled_assetserver.go    # Hand-written wrapper around embedded bundles
bundledassets/            # Embedded runtime JS assets
```

---

## Entwicklungsmodus

### Lebenszyklus

1. `wails3 dev` startet und **erzeugt den Prozess für Ihren Frontend-Entwicklungsserver** (Vite, SvelteKit, React-SWC …), indem die in `build/Taskfile.yml` definierte Task ausgeführt wird (üblicherweise `npm run dev`).
2. Die CLI setzt `WAILS_VITE_PORT` auf den Wails-Entwicklungsport und `FRONTEND_DEVSERVER_URL` auf die **vollständige** URL (`http://host:port` / `https://host:port`), die auf den laufenden Entwicklungsserver des Frameworks verweist. Siehe `internal/commands/dev.go`.
3. Der Entwicklungs-Asset-Server, der über `//go:build !production` in `build_dev.go` einkompiliert wird, liest `FRONTEND_DEVSERVER_URL` über `GetDevServerURL()` und leitet Datenverkehr, der nicht für die Runtime bestimmt ist, per Reverse-Proxy dorthin weiter.
4. Statische Dateien (`/assets/logo.svg`) können über `asset_fileserver.go` **direkt vom Datenträger bereitgestellt werden**, um die Geschwindigkeit zu erhöhen. Alle unbekannten Anfragen werden an den Entwicklungsserver des Frameworks **weitergeleitet**, sodass Hot-Module-Replacement *unmittelbar* erfolgt.

```
┌─────────┐  /wails/runtime.js     ┌─────────────┐
│ Browser │ ── embedded runtime ──▶│   Runtime   │
├─────────┤                        └─────────────┘
│   JS    │  / (index.html)        proxy / -> Vite via FRONTEND_DEVSERVER_URL
└─────────┘ ◀─────────────┐
              AssetServer │
                          ▼
                   ┌────────────┐
                   │  Vite Dev  │
                   │   Server   │
                   └────────────┘
```

### Funktionen

- **Live Reload** – Vite, SvelteKit usw. stellen HMR über WebSocket bereit; der Entwicklungs-Asset-Server leitet die Verbindung transparent weiter.
- **Source-Map-Unterstützung** – da die Assets nicht gebündelt sind, ordnen die Entwicklerwerkzeuge Ihres Browsers Fehler dem ursprünglichen Quellcode zu.
- **Keine erneute Go-Kompilierung** – nur das Frontend wird neu gebaut; der Go-Code läuft weiter, bis Sie `.go`-Dateien ändern.

### Framework wechseln

Der Entwicklungs-Proxy ist **frameworkunabhängig**. Beim Starten Ihrer Entwicklungs-Task stellt die Wails-CLI zwei Umgebungsvariablen bereit:

| Umgebungsvariable | Quelle | Bedeutung |
| --- | --- | --- |
| `WAILS_VITE_PORT` | `internal/commands/dev.go` (Konstante `wailsVitePort`) | Standardmäßiger Entwicklungsport (9245, sofern nicht `--port` übergeben wird) – Ihre Vite-Konfiguration sollte diesen Wert berücksichtigen |
| `FRONTEND_DEVSERVER_URL` | `internal/commands/dev.go` | Vollständige URL, an die Wails Anfragen weiterleitet; wird in Go über `assetserver.GetDevServerURL()` (`build_dev.go`) gelesen |

Im v3-Quellbaum gibt es keine Umgebungsvariable `VITE_PORT`, `FRONTEND_DEV_PORT` oder `WAILSDEV_VERBOSE`.

Fügen Sie eine neue Vorlage hinzu → definieren Sie deren Entwicklungs-Task → der Asset-Server funktioniert ohne weitere Anpassungen.

---

## Produktionsmodus

Wenn Sie `wails3 build` ausführen, führt die Pipeline folgende Schritte aus:

1. Sie führt den **Produktions-Build** des Frontends aus (`npm run build`), der `frontend/dist/**` erzeugt.
2. Sie **bettet** dieses Verzeichnis über `go:embed` in das anwendungseigene Paket ein (üblicherweise `//go:embed all:frontend/dist` neben `main.go`).
3. Sie kompiliert die ausführbare Go-Datei mit `-tags production`, das vom Taskfile-Wrapper über `EXTRA_TAGS` weitergereicht wird.

`internal/assetserver/build_production.go` ist der Build-Tag-Stub, der den Produktionscodepfad aktiviert. `internal/assetserver/bundled_assetserver.go` ist **manuell geschrieben**: Die Datei kapselt das in `bundledassets/` liegende Runtime-JS und wird nicht generiert.

### Anfrageverarbeitung

Der eigentliche Handler ist `internal/assetserver/assetserver.go` / `asset_fileserver.go`. Das Konzept:

1. Versuchen Sie, die eingebetteten statischen Assets unter dem angeforderten Pfad bereitzustellen.
2. Verwenden Sie für das SPA-Routing ersatzweise `index.html`.
3. Ermitteln Sie den Inhaltstyp, wenn die Dateierweiterung unbekannt ist (`content_type_sniffer.go`).
4. Setzen Sie sinnvolle Cache-Header.

- **MIME-Erkennung** – bei Dateien ohne Erweiterung wird der Inhaltstyp anhand der ersten etwa 512 Byte ermittelt (`content_type_sniffer.go`) und das Ergebnis in `mimecache.go` / `ringqueue.go` zwischengespeichert.
- **Sicherheits-Header** – unterbindet die Navigation über `file://` und setzt `nosniff`.

Da alles eingebettet ist, hat die ausgelieferte ausführbare Datei **keine externen Abhängigkeiten** – auch nicht unter Windows.

---

## Entwicklung ↔ Produktion verbinden

Aus Sicht von `pkg/application` stellen beide Modi **dieselbe öffentliche Schnittstelle** bereit: eine `AssetOptions`-Struktur mit einem `Handler http.Handler` sowie Middlewares und Lebenszyklusverdrahtung innerhalb von `internal/assetserver/`. Die Umschaltung zwischen Entwicklung und Produktion erfolgt vollständig über Go-Build-Tags, sodass der Anwendungscode in beiden Modi identisch ist.

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assetsFS),
    },
})
```

---

## Integration von Frontend-Frameworks

### Vorlagen

Jede mitgelieferte Vorlage (React, Vue, Svelte, Solid, Vanilla, …) enthält:

- `build/Taskfile.yml`
- `frontend/vite.config.ts` (oder eine Entsprechung)

Die Vite-Konfiguration bzw. die entsprechende Konfiguration liest `WAILS_VITE_PORT` und bindet den Entwicklungsserver an diesen Port. Anschließend veröffentlicht die CLI den tatsächlichen `FRONTEND_DEVSERVER_URL`, den der anwendungsinterne Proxy verwendet.

Die Frameworks bleiben vollständig von Go entkoppelt:

- Zur Build-Zeit muss kein Wails-JS-SDK importiert werden — `/wails/runtime.js` wird zur Laufzeit vom Asset-Server bereitgestellt.
- Jedes Framework mit einem HTTP-Entwicklungsserver lässt sich einbinden.

---

## Erweitern und anpassen

Benötigen Sie benutzerdefinierte Header, Authentifizierung oder gzip?

1. Definieren Sie einen `middleware.Middleware` (einen Alias für `func(http.Handler) http.Handler`, der in `internal/assetserver/middleware.go` deklariert ist).
2. Binden Sie ihn über die von `internal/assetserver/options.go` bereitgestellte Konfiguration in Ihren `application.AssetOptions` ein.
3. Das Verhalten ist in der Entwicklungs- und Produktionsumgebung identisch — es gibt keine separate Middleware-Liste für die einzelnen Modi.

---

## Wichtige Quelldateien

| Datei | Aufgabe |
| --- | --- |
| `build_dev.go` / `build_production.go` | Build-Tag-Wrapper zur Auswahl zwischen Entwicklung und Produktion |
| `assetserver.go` / `asset_fileserver.go` | Zentraler HTTP-Handler |
| `assetserver_dev.go` | Reverse-Proxy zu `FRONTEND_DEVSERVER_URL` |
| `bundled_assetserver.go` | Manuell erstellter Wrapper um `bundledassets/` |
| `options.go` | Konfiguration für `application.AssetOptions` |
| `mimecache.go` / `ringqueue.go` | MIME-Cache und kleiner LRU-Cache |

---

## Fallstricke und Fehlersuche

- **Weißer Bildschirm in der Produktionsumgebung** — meist liegt es am SPA-Routing: Stellen Sie sicher, dass Ihr Entwicklungsserver für unbekannte Pfade `index.html` bereitstellt und der Fallback des eingebetteten Produktions-Handlers erreicht wird.
- **404 in der Entwicklungsumgebung** — Ihre Vite-Konfiguration bindet nicht an `WAILS_VITE_PORT`, oder die CLI konnte den Entwicklungsserver nicht erreichen, um `FRONTEND_DEVSERVER_URL` zu setzen.
- **Große Assets** — durch das Einbetten wird die Binärdatei größer. Stellen Sie große Mediendateien von einem separaten Origin bereit oder streamen Sie sie über einen benutzerdefinierten `http.Handler`.

---

Sie wissen nun, wie der Wails-**Asset-Server** Ihren Webcode sowohl in der **Entwicklungs-** als auch in der **Produktionsumgebung** an das native Fenster übermittelt. Wenn Sie diese Schicht beherrschen, können Sie Ladeprobleme zuverlässig beheben, Middlewares hinzufügen oder sogar eine völlig andere Frontend-Toolchain einsetzen.
