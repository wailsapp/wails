---
title: "Frontend-Runtime"
description: "Das Wails-JavaScript-Runtime-Paket für die Frontend-Integration"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Die Wails-Frontend-Runtime ist die Standardbibliothek für Wails-Anwendungen. Sie stellt zahlreiche Funktionen bereit, die Sie in Ihren Anwendungen verwenden können, darunter:

- Fensterverwaltung
- Dialoge
- Browser-Integration
- Zwischenablage
- Menüs
- Systeminformationen
- Ereignisse
- Kontextmenüs
- Bildschirme
- WML (Wails Markup Language)

Die Runtime ist für die Integration zwischen Go und dem Frontend erforderlich. Es gibt 2 Möglichkeiten, die Runtime zu integrieren:

- Das Paket `@wailsio/runtime` verwenden
- Ein vorgefertigtes Bundle verwenden

## Das npm-Paket verwenden

Das Paket `@wailsio/runtime` ist ein JavaScript-Paket, das vom Frontend aus Zugriff auf die Wails-Runtime ermöglicht. Es wird von allen Standardvorlagen verwendet und ist die empfohlene Methode, die Runtime in Ihre Anwendung zu integrieren. Wenn Sie das Paket `@wailsio/runtime` verwenden, werden nur die von Ihnen genutzten Teile der Runtime eingebunden.

Das Paket ist auf npm verfügbar und kann wie folgt installiert werden:

```shell
npm install --save @wailsio/runtime
```

## Ein vorgefertigtes Bundle verwenden

Einige Projekte verwenden keinen JavaScript-Bundler und bevorzugen möglicherweise eine vorgefertigte, gebündelte Version der Runtime. Diese Version können Sie lokal mit folgendem Befehl erzeugen:

```shell
wails3 generate runtime
```

Der Befehl erzeugt im aktuellen Verzeichnis eine Datei `runtime.js` (und `runtime-debug.js`). Diese Datei ist ein ES-Modul, das Ihre Anwendungsskripte genau wie das npm-Paket importieren können. Die API wird jedoch zusätzlich über das globale Window-Objekt exportiert, sodass Sie sie in einfacheren Anwendungen wie folgt verwenden können:

```html
<html>
    <head>
        <script type="module" src="./runtime.js"></script>
        <script>
            window.onload = function () {
                wails.Window.SetTitle("A new window title");
            }
        </script>
    </head>
    <!--- ... -->
</html>
```

@note{type="caution"}
Das Tag `<script>`, das die Runtime lädt, muss unbedingt das Attribut `type="module"` enthalten. Warten Sie außerdem, bis die Seite vollständig geladen ist, bevor Sie die API aufrufen, da Skripte mit dem Attribut `type="module"` asynchron ausgeführt werden.

@end

## Initialisierung

Neben den API-Funktionen unterstützt die Runtime Kontextmenüs und das Ziehen von Fenstern. Diese Funktionen arbeiten erst nach der Initialisierung der Runtime wie erwartet. Auch wenn Sie die API nicht verwenden, müssen Sie an einer Stelle in Ihrem Frontend-Code eine Importanweisung einfügen, die ausschließlich wegen ihrer Seiteneffekte ausgeführt wird:

```javascript
import "@wailsio/runtime";
```

Ihr Bundler sollte die vorhandenen Seiteneffekte erkennen und den gesamten erforderlichen Initialisierungscode in den Build aufnehmen.

@note{type="info"}
Wenn Sie das vorgefertigte Bundle bevorzugen, genügt es, wie oben gezeigt ein Script-Tag hinzuzufügen.

@end

## Vite-Plugin für typisierte Ereignisse

Die Runtime enthält ein Vite-Plugin, das während der Entwicklung HMR (Hot Module Replacement) für typisierte Ereignisse ermöglicht.

### Einrichtung

Fügen Sie das Plugin zu Ihrer `vite.config.ts` hinzu:

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails("./bindings")],
})
```

Das Argument ist das Verzeichnis der generierten Bindings relativ zum Vite-Projekt. Generieren Sie die Bindings, bevor Sie das Frontend bauen.

### Vorteile

- **Automatisches Neuladen**: Ereignis-Bindings werden automatisch neu generiert und geladen, wenn Sie `wails3 generate bindings` ausführen
- **Entwicklungsmodus**: Funktioniert nahtlos mit `wails3 dev` und ermöglicht sofortige Aktualisierungen
- **Typsicherheit**: Vollständige TypeScript-Unterstützung mit automatischer Vervollständigung und Typprüfung

### Verwendung mit der Ereignisregistrierung

Registrieren Sie Ihre Ereignisse in Go:

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

Generieren Sie die Bindings:

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

Verwenden Sie typisierte Ereignisse in Ihrem Frontend:

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated, {
    ID: "123",
    Name: "John Doe"
})
```

## API-Referenz

Die Runtime ist in Module gegliedert, die jeweils bestimmte Funktionen bereitstellen. Importieren Sie nur, was Sie benötigen:

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### Ereignisse

Ereignissystem für die Kommunikation zwischen Go und JavaScript. Generierte Ereignisexporte wie `UserUpdated` sind Zeichenkettenkonstanten und keine Factory-Funktionen für Nutzdaten. Der Generator erweitert `CustomEvents`, um registrierte Ereignisnamen den Nutzdatentypen zuzuordnen. Die folgenden Typen (`WailsEventName`, `WailsEventCallback`, `WailsEventData`) werden vom Events-Modul exportiert.

#### On()

Registriert einen Callback für ein Ereignis.

```typescript
function On<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void
```

**Rückgabewert:** Funktion zum Abmelden

**Beispiel:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event listening
const unsubscribe = Events.On('user-logged-in', (event) => {
    console.log('User:', event.data.username)
})

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

Events.On(UserLogin, (event) => {
    // event.data is typed as UserLoginData
    console.log('User:', event.data.username)
})

// Later: unsubscribe()
```

#### Once()

Registriert einen Callback, der nur einmal ausgeführt wird.

```typescript
function Once<E extends WailsEventName = WailsEventName>(eventName: E, callback: WailsEventCallback<E>): () => void
```

**Beispiel:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Sendet ein Ereignis an das Go-Backend oder andere Fenster.

```typescript
function Emit<E extends WailsEventName = WailsEventName>(name: E, data: WailsEventData<E>): Promise<boolean>
function Emit<E extends WailsEventName = WailsEventName>(name: WailsEventData<E> extends null | void ? E : never): Promise<boolean>
```

**Rückgabewert:** Promise, die zu `true` aufgelöst wird, wenn das Ereignis abgebrochen wurde, andernfalls zu `false`

**Beispiel:**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event emission
const wasCancelled = await Events.Emit('button-clicked', { buttonId: 'submit' })

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

const cancelled = await Events.Emit(UserLogin, {
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
})

if (cancelled) {
    console.log('Login was cancelled by a hook')
}
```

@note{type="info"}
Der Rückgabewert gibt an, ob das Ereignis durch einen Hook abgebrochen wurde. Die meisten Ereignisse können nicht abgebrochen werden und geben immer `false` zurück.

@end

#### Off()

Entfernt Ereignis-Listener.

```typescript
function Off(...eventNames: string[]): void
```

**Beispiel:**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

Entfernt alle Ereignis-Listener.

```typescript
function OffAll(): void
```

### Fenster

Methoden zur Fensterverwaltung. Der Standardexport ist das aktuelle Fenster.

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### Sichtbarkeit

**Show()** – Zeigt das Fenster an

```typescript
function Show(): Promise<void>
```

**Hide()** – Blendet das Fenster aus

```typescript
function Hide(): Promise<void>
```

**Close()** – Schließt das Fenster

```typescript
function Close(): Promise<void>
```

#### Größe und Position

**SetSize(width, height)** – Legt die Fenstergröße fest

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** – Ruft die Fenstergröße ab

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** – Legt die absolute Position fest

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** – Ruft die absolute Position ab

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** – Zentriert das Fenster

```typescript
function Center(): Promise<void>
```

**Beispiel:**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### Fensterstatus

**Minimise()** – Minimiert das Fenster

```typescript
function Minimise(): Promise<void>
```

**Maximise()** – Maximiert das Fenster

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** – Wechselt in den Vollbildmodus

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** – Stellt das Fenster aus dem minimierten, maximierten oder Vollbildzustand wieder her

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** – Prüft, ob das Fenster minimiert ist

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** – Prüft, ob das Fenster maximiert ist

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** – Prüft, ob das Fenster im Vollbildmodus ist

```typescript
function IsFullscreen(): Promise<boolean>
```

#### Fenstereigenschaften

**SetTitle(title)** – Legt den Fenstertitel fest

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** – Ruft den Fensternamen ab

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** – Legt die Hintergrundfarbe fest

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** – Hält das Fenster im Vordergrund

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** – Macht die Fenstergröße veränderbar

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### Fokus und Bildschirm

**Focus()** – Fokussiert das Fenster

```typescript
function Focus(): Promise<void>
```

**IsFocused()** – Prüft, ob das Fenster fokussiert ist

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** – Ruft den Bildschirm ab, auf dem sich das Fenster befindet

```typescript
function GetScreen(): Promise<Screen>
```

#### Inhalt

**Reload()** – Lädt die Seite neu

```typescript
function Reload(): Promise<void>
```

**ForceReload()** – Erzwingt das Neuladen der Seite und leert den Cache

```typescript
function ForceReload(): Promise<void>
```

#### Zoom

**SetZoom(level)** – Legt die Zoomstufe fest

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** – Ruft die Zoomstufe ab

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** – Vergrößert die Ansicht

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** – Verkleinert die Ansicht

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** – Setzt den Zoom auf 100 % zurück

```typescript
function ZoomReset(): Promise<void>
```

#### Drucken

**Print()** – Öffnet den nativen Druckdialog

```typescript
function Print(): Promise<void>
```

**Beispiel:**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**Hinweis:** Dies öffnet den nativen Druckdialog des Betriebssystems, in dem Benutzer die Druckereinstellungen auswählen und den aktuellen Fensterinhalt drucken können. Im Gegensatz zu `window.print()`, das in Webviews möglicherweise nicht funktioniert, wird dabei die native Druck-API der Plattform verwendet.

### Zwischenablage

Operationen mit der Zwischenablage.

#### SetText()

Legt den Text der Zwischenablage fest.

```typescript
function SetText(text: string): Promise<void>
```

**Beispiel:**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

Ruft den Text aus der Zwischenablage ab.

```typescript
function Text(): Promise<string>
```

**Beispiel:**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### System

Systemmethoden auf niedriger Ebene für die direkte Kommunikation mit dem Backend.

#### invoke()

Sendet eine Rohdatenmeldung direkt an das Backend. Dies umgeht das standardmäßige Bindungssystem und wird von `RawMessageHandler` in den Anwendungsoptionen verarbeitet.

```typescript
function invoke(message: any): void
```

**Beispiel:**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
Diese Funktion arbeitet nach dem Fire-and-forget-Prinzip und gibt keinen Wert zurück. Verwenden Sie Ereignisse, um Antworten vom Backend zu empfangen.

@end

Weitere Informationen finden Sie im [Leitfaden zu Rohdatenmeldungen](/guides/raw-messages/).

### Anwendung

Methoden auf Anwendungsebene.

#### Show()

Zeigt alle Anwendungsfenster an.

```typescript
function Show(): Promise<void>
```

#### Hide()

Blendet alle Anwendungsfenster aus.

```typescript
function Hide(): Promise<void>
```

#### Quit()

Beendet die Anwendung.

```typescript
function Quit(): Promise<void>
```

**Beispiel:**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### Browser

Öffnet URLs im Standardbrowser.

#### OpenURL()

Öffnet eine URL im Systembrowser.

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**Beispiel:**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### Bildschirme

Bildschirminformationen und -verwaltung.

#### GetAll()

Ruft alle Bildschirme ab.

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

Ruft den primären Bildschirm ab.

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

Ruft den derzeit aktiven Bildschirm ab.

```typescript
function GetCurrent(): Promise<Screen>
```

**Bildschirmschnittstelle:**

```typescript
interface Screen {
    ID: string
    Name: string
    ScaleFactor: number
    X: number
    Y: number
    Size: { Width: number, Height: number }
    Bounds: { X: number, Y: number, Width: number, Height: number }
    WorkArea: { X: number, Y: number, Width: number, Height: number }
    IsPrimary: boolean
    Rotation: number
}
```

**Beispiel:**

```javascript
import { Screens } from '@wailsio/runtime'

// List all screens
const screens = await Screens.GetAll()
screens.forEach(screen => {
    console.log(`${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})

// Get primary screen
const primary = await Screens.GetPrimary()
console.log('Primary screen:', primary.Name)
```

### Dialoge

Native Betriebssystemdialoge aus JavaScript.

#### Info()

Zeigt einen Informationsdialog an.

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**Beispiel:**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

Zeigt einen Fehlerdialog an.

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

Zeigt einen Warnungsdialog an.

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

Zeigt einen Fragedialog mit benutzerdefinierten Schaltflächen an.

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**Beispiel:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const result = await Dialogs.Question({
    Title: 'Confirm Delete',
    Message: 'Are you sure you want to delete this file?',
    Buttons: [
        { Label: 'Delete', IsDefault: false },
        { Label: 'Cancel', IsDefault: true }
    ]
})

if (result === 'Delete') {
    // Delete the file
}
```

#### OpenFile()

Zeigt einen Dialog zum Öffnen einer Datei an.

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**Beispiel:**

```javascript
import { Dialogs } from '@wailsio/runtime'

const file = await Dialogs.OpenFile({
    Title: 'Select Image',
    Filters: [
        { DisplayName: 'Images', Pattern: '*.png;*.jpg;*.jpeg' },
        { DisplayName: 'All Files', Pattern: '*.*' }
    ]
})

if (file) {
    console.log('Selected:', file)
}
```

#### SaveFile()

Zeigt einen Dialog zum Speichern einer Datei an.

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML (Wails Markup Language)

WML bietet deklarative Attribute für häufig verwendete Aktionen. Fügen Sie HTML-Elementen Attribute hinzu:

#### Attribute

**wml-event** – Löst beim Anklicken ein Ereignis aus

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** – Ruft eine Fenstermethode auf

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** – Gibt das Zielfenster für wml-window an

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** – Öffnet eine URL im Browser

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** – Zeigt vor der Aktion einen Bestätigungsdialog an

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**Beispiel:**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## Vollständiges Beispiel

```javascript
import { Events, Window, Clipboard, Dialogs, Screens } from '@wailsio/runtime'

// Listen for events from Go
Events.On('data-updated', (event) => {
    console.log('Data:', event.data)
    updateUI(event.data)
})

// Window management
document.getElementById('center-btn').addEventListener('click', async () => {
    await Window.Center()
})

document.getElementById('fullscreen-btn').addEventListener('click', async () => {
    const isFullscreen = await Window.IsFullscreen()
    if (isFullscreen) {
        await Window.UnFullscreen()
    } else {
        await Window.Fullscreen()
    }
})

// Clipboard operations
document.getElementById('copy-btn').addEventListener('click', async () => {
    await Clipboard.SetText('Copied from Wails!')
})

// Dialog with confirmation
document.getElementById('delete-btn').addEventListener('click', async () => {
    const result = await Dialogs.Question({
        Title: 'Confirm',
        Message: 'Delete this item?',
        Buttons: [
            { Label: 'Delete' },
            { Label: 'Cancel', IsDefault: true }
        ]
    })

    if (result === 'Delete') {
        await Events.Emit('delete-item', { id: currentItemId })
    }
})

// Screen information
const screens = await Screens.GetAll()
console.log(`Detected ${screens.length} screen(s)`)
screens.forEach(screen => {
    console.log(`- ${screen.Name}: ${screen.Size.Width}x${screen.Size.Height}`)
})
```

## Bewährte Vorgehensweisen

### ✅ Empfohlen

- **Selektiv importieren** – Importieren Sie nur, was Sie benötigen
- **Promises verarbeiten** – Alle Methoden geben Promises zurück
- **WML für einfache Aktionen verwenden** – übersichtlicher als JavaScript
- **Rückgabewerte prüfen** – insbesondere bei Dialogen
- **Ereignisse abbestellen** – nach Abschluss bereinigen

### ❌ Nicht empfohlen

- **await nicht vergessen** – die meisten Methoden sind asynchron
- **Benutzeroberfläche nicht blockieren** – async/await korrekt verwenden
- **Fehler nicht ignorieren** – Ablehnungen immer behandeln

## TypeScript-Unterstützung

Die Runtime enthält vollständige TypeScript-Typdefinitionen:

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
