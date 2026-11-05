---
title: "Andere Frontend-Frameworks verwenden"
description: "So verwenden Sie ein Framework ohne integrierte Vorlage, indem Sie Ihr eigenes Vite-Projekt im Verzeichnis frontend ablegen"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails enthält bewusst nur für eine kleine Auswahl von Frameworks integrierte Startvorlagen:

| Vorlage | Sprache |
| --- | --- |
| `vanilla` | TypeScript (Standard) |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

Das bedeutet nicht, dass dies Ihre einzigen Optionen sind. Das Frontend einer Wails-App ist **lediglich ein Webprojekt** – alles, was sich als statisches HTML/CSS/JS erstellen lässt, funktioniert. Wenn für das Framework Ihrer Wahl (Solid, Preact, Lit, Qwik, SvelteKit, Angular, …) keine Vorlage vorhanden ist, können Sie das Projekt in wenigen Minuten selbst erzeugen.

## So wird das Verzeichnis `frontend/` verwendet

Wails ist es gleichgültig, welches Framework sich in `frontend/` befindet. Es setzt lediglich einen kleinen, frameworkunabhängigen Vertrag voraus:

- **`frontend/dist/` wird ausgeliefert.** `main.go` bettet das erstellte Frontend mit `//go:embed all:frontend/dist` ein und stellt es über den Asset-Server bereit. Ihr Build muss ein statisches Bundle nach `frontend/dist/` ausgeben (Vites Standard-Ausgabeverzeichnis).
- **Der Build wird von `frontend/package.json` gesteuert.** Während `wails3 build` führt Wails das Skript `build` des Frontends aus; während `wails3 dev` führt es `dev` aus und leitet den Vite-Entwicklungsserver für Hot Reload weiter.
- **Bindings werden in `frontend/bindings/` erzeugt.** Wails untersucht Ihre registrierten Go-Dienste und schreibt dort ein typsicheres SDK. Sie importieren daraus wie aus jedem anderen Modul:
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **Der Entwicklungsserver wird an einem festen Port ausgeführt.** `wails3 dev` leitet Vite über den in `WAILS_VITE_PORT` angegebenen Port weiter (standardmäßig `9245`). Legen Sie daher `server.port` mit `strictPort: true` auf diesen Port fest. Dies ist eine gewöhnliche Vite-Konfiguration – daran ist kein Wails-Plug-in beteiligt.
- **Optional – typisierte benutzerdefinierte Ereignisse.** Die integrierten Vorlagen registrieren außerdem das Plug-in `@wailsio/runtime/plugins/vite`. Sie benötigen es nur, wenn Sie *typisierte* benutzerdefinierte Ereignisse verwenden: Es fügt Ihre erzeugten Ereignistypdefinitionen in die Runtime ein und lässt den Build fehlschlagen, solange keine Bindings erzeugt wurden. Wenn Sie nur die zeichenkettenbasierte API `Events.On("time", …)` verwenden, können Sie es weglassen.

Alles Weitere – Komponenten, Routing, Zustand und Styling – bleibt vollständig Ihrem Framework überlassen.

## Beliebige Frameworks mit Vite erzeugen

Am schnellsten beginnen Sie mit einer integrierten Vorlage (damit Sie `main.go`, `Taskfile`, Build-Assets und einen funktionierenden Go-Dienst erhalten) und ersetzen anschließend `frontend/` durch ein neues Vite-Projekt für Ihr Framework.

@steps
### Ein Projekt aus der Standardvorlage erstellen
```bash
wails3 init -n myapp
cd myapp
```

### `frontend/` durch eine Vite-App für Ihr Framework ersetzen
Vite kann für die meisten Frameworks mit einem einzigen Befehl ein Projekt erzeugen. Wählen Sie eine Vorlage:

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

Ersetzen Sie `solid` durch eine beliebige Vite-Vorlage: `preact`, `lit`, `svelte`, `vue`, `react`, `vanilla` oder deren `-ts`-Varianten (`solid-ts`, `preact-ts`, …).

### Runtime installieren und Vite auf den Wails-Entwicklungsserver verweisen
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime` stellt die JS-APIs bereit (`Events`, `Browser`, Dialoge, …). Die einzige *erforderliche* Änderung an `vite.config` ist der Port des Entwicklungsservers, damit `wails3 dev` ihn finden kann:

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

Fügen Sie das Plug-in nur dann ebenfalls hinzu, wenn Sie **typisierte benutzerdefinierte Ereignisse** verwenden möchten. Es fügt Ihre erzeugten Ereignistypen ein und setzt voraus, dass die Bindings bereits vor dem Build vorhanden sind:

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Ihre Go-Dienste aufrufen
Erzeugen Sie die Bindings einmal und importieren Sie sie anschließend an beliebiger Stelle in Ihren Komponenten:

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### Ausführen
```bash
wails3 dev
```

@end

@note{type="tip" title="Ein reines JavaScript-Projekt erzeugen"}
Derselbe Befehl erzeugt ein Projekt ohne TypeScript – verwenden Sie einfach eine Vite-Vorlage ohne `-ts`:

```bash
npm create vite@latest frontend -- --template solid
```

@end

## Frameworks mit eigenem Scaffolding-Werkzeug

Einige Frameworks werden nicht über die `create`-Vorlagen von Vite erstellt und verfügen über eigene Werkzeuge. Sie funktionieren dennoch: Erzeugen Sie das Projekt einfach mit dem jeweiligen nativen Befehl und fügen Sie anschließend das Wails-Plug-in hinzu:

- **SvelteKit:** `npx sv create frontend`. Verwenden Sie den statischen Adapter (`@sveltejs/adapter-static`), damit ein statisches Bundle erstellt wird, und deaktivieren Sie SSR.
- **Qwik:** `npm create qwik@latest`. Verwenden Sie den statischen Adapter (SSG).
- **Angular:** Erzeugen Sie das Projekt mit `ng new`, setzen Sie `outputPath` auf `dist` und lassen Sie das Build-Skript auf `ng build` verweisen.

Die Regel ist immer dieselbe: Erzeugen Sie einen statischen Build in `frontend/dist/`, behalten Sie das Vite-Plug-in `@wailsio/runtime` bei (oder importieren Sie die Runtime direkt) und importieren Sie Ihre Go-Bindings aus `frontend/bindings/`.

@note{type="info"}
Wenn Sie eine ausgereifte Einrichtung für ein Framework erstellen, sollten Sie erwägen, sie als [benutzerdefinierte Vorlage](/guides/advanced/custom-templates/) zu veröffentlichen, damit andere sie direkt mit `wails3 init -t` verwenden können.

@end
