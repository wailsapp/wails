---
title: "Runtime frontend"
description: "Le package d’exécution JavaScript de Wails pour l’intégration au frontend"
slug: "reference/frontend-runtime"
sourcePath: "reference/frontend-runtime.md"
---

Le runtime frontend de Wails constitue la bibliothèque standard des applications Wails. Il fournit de nombreuses fonctionnalités utilisables dans vos applications, notamment :

- Gestion des fenêtres
- Boîtes de dialogue
- Intégration au navigateur
- Presse-papiers
- Menus
- Informations système
- Événements
- Menus contextuels
- Écrans
- WML (langage de balisage Wails)

Le runtime est nécessaire à l’intégration entre Go et le frontend. Il existe 2 façons d’intégrer le runtime :

- Utiliser le package `@wailsio/runtime`
- Utiliser un bundle précompilé

## Utiliser le package npm

Le package `@wailsio/runtime` est un package JavaScript qui permet d’accéder au runtime Wails depuis le frontend. Tous les modèles standard l’utilisent et il s’agit de la méthode recommandée pour intégrer le runtime à votre application. Avec le package `@wailsio/runtime`, seules les parties du runtime que vous utilisez sont incluses.

Le package est disponible sur npm et peut être installé avec :

```shell
npm install --save @wailsio/runtime
```

## Utiliser un bundle précompilé

Certains projets n’utilisent pas d’outil de regroupement JavaScript et peuvent préférer une version précompilée du runtime sous forme de bundle. Vous pouvez générer cette version localement à l’aide de la commande suivante :

```shell
wails3 generate runtime
```

La commande génère un fichier `runtime.js` (ainsi qu’un fichier `runtime.debug.js`) dans le répertoire courant. Ce fichier est un module ES que les scripts de votre application peuvent importer comme le package npm, mais l’API est également exportée dans l’objet global window ; les applications plus simples peuvent donc l’utiliser comme suit :

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
Il est important d’inclure l’attribut `type="module"` dans la balise `<script>` qui charge le runtime et d’attendre le chargement complet de la page avant d’appeler l’API, car les scripts dotés de l’attribut `type="module"` s’exécutent de manière asynchrone.

@end

## Initialisation

Outre les fonctions de l’API, le runtime prend en charge les menus contextuels et le déplacement des fenêtres. Ces fonctionnalités ne fonctionneront comme prévu qu’après l’initialisation du runtime. Même si vous n’utilisez pas l’API, veillez à inclure quelque part dans votre code frontend une instruction d’importation avec effet de bord :

```javascript
import "@wailsio/runtime";
```

Votre outil de regroupement devrait détecter la présence d’effets de bord et inclure dans le build tout le code d’initialisation requis.

@note{type="info"}
Si vous préférez le bundle précompilé, il suffit d’ajouter une balise script comme indiqué ci-dessus.

@end

## Plugin Vite pour les événements typés

Le runtime comprend un plugin Vite qui active la prise en charge du HMR (remplacement de module à chaud) pour les événements typés pendant le développement.

### Configuration

Ajoutez le plugin à votre fichier `vite.config.ts` :

```typescript
import { defineConfig } from 'vite'
import wails from '@wailsio/runtime/plugins/vite'

export default defineConfig({
  plugins: [wails()],
})
```

### Avantages

- **Rechargement automatique** : les liaisons d’événements sont automatiquement régénérées et rechargées lorsque vous exécutez `wails3 generate bindings`
- **Mode développement** : fonctionne de manière transparente avec `wails3 dev` pour appliquer les mises à jour instantanément
- **Sécurité des types** : prise en charge complète de TypeScript avec autocomplétion et vérification des types

### Utilisation avec l’enregistrement des événements

Enregistrez vos événements dans Go :

```go
type UserData struct {
    ID   string
    Name string
}

func init() {
    application.RegisterEvent[UserData]("user-updated")
}
```

Générez les liaisons :

```bash
wails3 generate bindings

# Or, to include TypeScript definitions
wails3 generate bindings -ts

# For more options, see:
wails3 generate bindings -help
```

Utilisez les événements typés dans votre frontend :

```typescript
import { Events } from '@wailsio/runtime'
import { UserUpdated } from './bindings/events'

// Type-safe event with autocomplete
Events.Emit(UserUpdated({
    ID: "123",
    Name: "John Doe"
}))
```

## Référence de l’API

Le runtime est organisé en modules, chacun fournissant des fonctionnalités spécifiques. Importez uniquement ce dont vous avez besoin :

```javascript
import { Events, Window, Clipboard } from '@wailsio/runtime'
```

### Événements

Système d’événements permettant la communication entre Go et JavaScript.

#### On()

Enregistre une fonction de rappel pour un événement.

```typescript
function On(eventName: string, callback: (event: WailsEvent) => void): () => void
function On<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**Renvoie :** fonction de désabonnement

**Exemple :**

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

Enregistre une fonction de rappel qui ne s’exécute qu’une seule fois.

```typescript
function Once(eventName: string, callback: (event: WailsEvent) => void): () => void
function Once<T>(eventType: EventType<T>, callback: (event: WailsEvent<T>) => void): () => void
```

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

Events.Once('app-ready', () => {
    console.log('App initialized')
})
```

#### Emit()

Émet un événement vers le backend Go ou d’autres fenêtres.

```typescript
function Emit(name: string, data?: any): Promise<boolean>
function Emit<T>(event: Event<T>): Promise<boolean>
```

**Renvoie :** une promesse résolue avec `true` si l’événement a été annulé, et avec `false` dans le cas contraire

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

// Basic event emission
const wasCancelled = await Events.Emit('button-clicked', { buttonId: 'submit' })

// With typed events (TypeScript)
import { UserLogin } from './bindings/events'

const cancelled = await Events.Emit(UserLogin({
    UserID: "123",
    Username: "john_doe",
    LoginTime: new Date().toISOString()
}))

if (cancelled) {
    console.log('Login was cancelled by a hook')
}
```

@note{type="info"}
La valeur renvoyée indique si l’événement a été annulé par un hook. La plupart des événements ne peuvent pas être annulés et renvoient toujours `false`.

@end

#### Off()

Supprime des écouteurs d’événements.

```typescript
function Off(...eventNames: string[]): void
```

**Exemple :**

```javascript
import { Events } from '@wailsio/runtime'

Events.Off('user-logged-in', 'user-logged-out')
```

#### OffAll()

Supprime tous les écouteurs d’événements.

```typescript
function OffAll(): void
```

### Fenêtre

Méthodes de gestion des fenêtres. L’exportation par défaut correspond à la fenêtre actuelle.

```javascript
import { Window } from '@wailsio/runtime'

// Current window
await Window.SetTitle('New Title')
await Window.Center()

// Get another window
const otherWindow = Window.Get('secondary')
await otherWindow.Show()
```

#### Visibilité

**Show()** — Affiche la fenêtre

```typescript
function Show(): Promise<void>
```

**Hide()** — Masque la fenêtre

```typescript
function Hide(): Promise<void>
```

**Close()** - Ferme la fenêtre

```typescript
function Close(): Promise<void>
```

#### Taille et position

**SetSize(width, height)** - Définit la taille de la fenêtre

```typescript
function SetSize(width: number, height: number): Promise<void>
```

**Size()** - Récupère la taille de la fenêtre

```typescript
function Size(): Promise<{ width: number, height: number }>
```

**SetPosition(x, y)** - Définit la position absolue

```typescript
function SetPosition(x: number, y: number): Promise<void>
```

**Position()** - Récupère la position absolue

```typescript
function Position(): Promise<{ x: number, y: number }>
```

**Center()** - Centre la fenêtre

```typescript
function Center(): Promise<void>
```

**Exemple :**

```javascript
import { Window } from '@wailsio/runtime'

// Resize and center
await Window.SetSize(800, 600)
await Window.Center()

// Get current size
const { width, height } = await Window.Size()
```

#### État de la fenêtre

**Minimise()** - Réduit la fenêtre

```typescript
function Minimise(): Promise<void>
```

**Maximise()** - Agrandit la fenêtre

```typescript
function Maximise(): Promise<void>
```

**Fullscreen()** - Passe en plein écran

```typescript
function Fullscreen(): Promise<void>
```

**Restore()** - Restaure la fenêtre depuis l’état réduit, agrandi ou plein écran

```typescript
function Restore(): Promise<void>
```

**IsMinimised()** - Vérifie si la fenêtre est réduite

```typescript
function IsMinimised(): Promise<boolean>
```

**IsMaximised()** - Vérifie si la fenêtre est agrandie

```typescript
function IsMaximised(): Promise<boolean>
```

**IsFullscreen()** - Vérifie si la fenêtre est en plein écran

```typescript
function IsFullscreen(): Promise<boolean>
```

#### Propriétés de la fenêtre

**SetTitle(title)** - Définit le titre de la fenêtre

```typescript
function SetTitle(title: string): Promise<void>
```

**Name()** - Récupère le nom de la fenêtre

```typescript
function Name(): Promise<string>
```

**SetBackgroundColour(r, g, b, a)** - Définit la couleur d’arrière-plan

```typescript
function SetBackgroundColour(r: number, g: number, b: number, a: number): Promise<void>
```

**SetAlwaysOnTop(alwaysOnTop)** - Maintient la fenêtre au premier plan

```typescript
function SetAlwaysOnTop(alwaysOnTop: boolean): Promise<void>
```

**SetResizable(resizable)** - Rend la fenêtre redimensionnable

```typescript
function SetResizable(resizable: boolean): Promise<void>
```

#### Focus et écran

**Focus()** - Donne le focus à la fenêtre

```typescript
function Focus(): Promise<void>
```

**IsFocused()** - Vérifie si la fenêtre a le focus

```typescript
function IsFocused(): Promise<boolean>
```

**GetScreen()** - Récupère l’écran sur lequel se trouve la fenêtre

```typescript
function GetScreen(): Promise<Screen>
```

#### Contenu

**Reload()** - Recharge la page

```typescript
function Reload(): Promise<void>
```

**ForceReload()** - Force le rechargement de la page (vide le cache)

```typescript
function ForceReload(): Promise<void>
```

#### Zoom

**SetZoom(level)** - Définit le niveau de zoom

```typescript
function SetZoom(level: number): Promise<void>
```

**GetZoom()** - Récupère le niveau de zoom

```typescript
function GetZoom(): Promise<number>
```

**ZoomIn()** - Augmente le zoom

```typescript
function ZoomIn(): Promise<void>
```

**ZoomOut()** - Réduit le zoom

```typescript
function ZoomOut(): Promise<void>
```

**ZoomReset()** - Réinitialise le zoom à 100 %

```typescript
function ZoomReset(): Promise<void>
```

#### Impression

**Print()** - Ouvre la boîte de dialogue d’impression native

```typescript
function Print(): Promise<void>
```

**Exemple :**

```javascript
import { Window } from '@wailsio/runtime'

// Open print dialog for current window
await Window.Print()
```

**Remarque :** cette méthode ouvre la boîte de dialogue d’impression native du système d’exploitation, permettant à l’utilisateur de sélectionner les paramètres de l’imprimante et d’imprimer le contenu actuel de la fenêtre. Contrairement à `window.print()`, qui peut ne pas fonctionner dans les vues web, elle utilise l’API d’impression native de la plateforme.

### Presse-papiers

Opérations sur le presse-papiers.

#### SetText()

Définit le texte du presse-papiers.

```typescript
function SetText(text: string): Promise<void>
```

**Exemple :**

```javascript
import { Clipboard } from '@wailsio/runtime'

await Clipboard.SetText('Hello from Wails!')
```

#### Text()

Récupère le texte du presse-papiers.

```typescript
function Text(): Promise<string>
```

**Exemple :**

```javascript
import { Clipboard } from '@wailsio/runtime'

const clipboardText = await Clipboard.Text()
console.log('Clipboard:', clipboardText)
```

### Système

Méthodes système de bas niveau permettant de communiquer directement avec le backend.

#### invoke()

Envoie un message brut directement au backend. Cette méthode contourne le système de liaison standard et le message est traité par le `RawMessageHandler` défini dans les options de votre application.

```typescript
function invoke(message: any): void
```

**Exemple :**

```javascript
import { System } from '@wailsio/runtime'

// Send a raw message to the backend
System.invoke('my-custom-message')

// Send structured data as JSON
System.invoke(JSON.stringify({ action: 'update', value: 42 }))
```

@note{type="caution"}
Cette fonction est exécutée sans attendre de réponse et ne renvoie aucune valeur. Utilisez des événements pour recevoir les réponses du backend.

@end

Pour plus de détails, consultez le [guide des messages bruts](/guides/raw-messages/).

### Application

Méthodes au niveau de l’application.

#### Show()

Affiche toutes les fenêtres de l’application.

```typescript
function Show(): Promise<void>
```

#### Hide()

Masque toutes les fenêtres de l’application.

```typescript
function Hide(): Promise<void>
```

#### Quit()

Quitte l’application.

```typescript
function Quit(): Promise<void>
```

**Exemple :**

```javascript
import { Application } from '@wailsio/runtime'

// Add quit button
document.getElementById('quit-btn').addEventListener('click', async () => {
    await Application.Quit()
})
```

### Navigateur

Ouvrez des URL dans le navigateur par défaut.

#### OpenURL()

Ouvre une URL dans le navigateur du système.

```typescript
function OpenURL(url: string | URL): Promise<void>
```

**Exemple :**

```javascript
import { Browser } from '@wailsio/runtime'

await Browser.OpenURL('https://wails.io')
```

### Écrans

Informations sur les écrans et gestion de ceux-ci.

#### GetAll()

Récupère tous les écrans.

```typescript
function GetAll(): Promise<Screen[]>
```

#### GetPrimary()

Récupère l’écran principal.

```typescript
function GetPrimary(): Promise<Screen>
```

#### GetCurrent()

Récupère l’écran actif actuel.

```typescript
function GetCurrent(): Promise<Screen>
```

**Interface Screen :**

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

**Exemple :**

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

### Boîtes de dialogue

Boîtes de dialogue natives du système d’exploitation accessibles depuis JavaScript.

#### Info()

Affiche une boîte de dialogue d’information.

```typescript
function Info(options: MessageDialogOptions): Promise<string>
```

**Exemple :**

```javascript
import { Dialogs } from '@wailsio/runtime'

await Dialogs.Info({
    Title: 'Success',
    Message: 'Operation completed successfully!'
})
```

#### Error()

Affiche une boîte de dialogue d’erreur.

```typescript
function Error(options: MessageDialogOptions): Promise<string>
```

#### Warning()

Affiche une boîte de dialogue d’avertissement.

```typescript
function Warning(options: MessageDialogOptions): Promise<string>
```

#### Question()

Affiche une boîte de dialogue de question avec des boutons personnalisés.

```typescript
function Question(options: MessageDialogOptions): Promise<string>
```

**Exemple :**

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

Affiche une boîte de dialogue d’ouverture de fichier.

```typescript
function OpenFile(options: OpenFileDialogOptions): Promise<string | string[]>
```

**Exemple :**

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

Affiche une boîte de dialogue d’enregistrement de fichier.

```typescript
function SaveFile(options: SaveFileDialogOptions): Promise<string>
```

### WML (langage de balisage Wails)

WML fournit des attributs déclaratifs pour les actions courantes. Ajoutez des attributs aux éléments HTML :

#### Attributs

**wml-event** – Émet un événement lors d’un clic

```html
<button wml-event="save-clicked">Save</button>
```

**wml-window** – Appelle une méthode de fenêtre

```html
<button wml-window="Close">Close Window</button>
<button wml-window="Minimise">Minimize</button>
```

**wml-target-window** – Spécifie la fenêtre cible pour wml-window

```html
<button wml-window="Show" wml-target-window="settings">
    Show Settings
</button>
```

**wml-openurl** – Ouvre une URL dans le navigateur

```html
<a href="#" wml-openurl="https://wails.io">Visit Wails</a>
```

**wml-confirm** – Affiche une boîte de dialogue de confirmation avant l’action

```html
<button wml-window="Close" wml-confirm="Are you sure you want to close?">
    Close
</button>
```

**Exemple :**

```html
<div>
    <button wml-event="save-clicked">Save</button>
    <button wml-window="Minimise">Minimize</button>
    <button wml-window="Close" wml-confirm="Close window?">Close</button>
    <a href="#" wml-openurl="https://github.com/wailsapp/wails">GitHub</a>
</div>
```

## Exemple complet

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

## Bonnes pratiques

### ✅ À faire

- **Effectuez des importations sélectives** – Importez uniquement ce dont vous avez besoin
- **Gérez les promesses** – Toutes les méthodes renvoient des promesses
- **Utilisez WML pour les actions simples** – Le code est plus clair qu’en JavaScript
- **Vérifiez les valeurs de retour** – En particulier pour les boîtes de dialogue
- **Désabonnez-vous des événements** – Effectuez le nettoyage lorsque vous avez terminé

### ❌ À ne pas faire

- **N’oubliez pas await** – La plupart des méthodes sont asynchrones
- **Ne bloquez pas l’interface utilisateur** – Utilisez correctement async/await
- **N’ignorez pas les erreurs** – Gérez toujours les rejets de promesse

## Prise en charge de TypeScript

Le runtime inclut des définitions TypeScript complètes :

```typescript
import { Events, Window } from '@wailsio/runtime'

Events.On('custom-event', (event) => {
    // TypeScript knows event.data, event.name, event.sender
    console.log(event.data)
})

// All methods are fully typed
const size: { width: number, height: number } = await Window.Size()
```
