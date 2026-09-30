---
title: "Service de génération de codes QR"
description: "Créez un service de génération de codes QR pour découvrir les services Wails"
slug: "tutorials/01-creating-a-service"
sourcePath: "tutorials/01-creating-a-service.md"
---

Dans Wails, un **service** est une structure Go qui contient la logique métier que vous souhaitez mettre à la disposition de votre frontend. Les services organisent votre code en regroupant les fonctionnalités associées.

Considérez un service comme un ensemble de méthodes que votre code JavaScript peut appeler. Après la génération des bindings, chaque méthode publique du service peut être appelée depuis le frontend.

Dans ce tutoriel, nous allons créer un service de génération de codes QR pour illustrer ces concepts. À la fin, vous saurez créer des services, gérer les dépendances et connecter votre code Go au frontend.

<br/>

@steps
### Créer le fichier du service de génération de codes QR
Dans le répertoire de votre application, créez un fichier nommé `qrservice.go` :

```go {title="qrservice.go"}
package main

import (
    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}
```

**Voici ce qui se passe :**

- `QRService` est une structure vide qui contiendra nos méthodes de génération de codes QR
- `NewQRService()` est une fonction constructeur qui crée une instance de notre service
- `Generate()` est une méthode qui reçoit un texte et une taille, puis renvoie le code QR sous forme de tableau d’octets PNG
- Conformément à la convention de Go qui consiste à renvoyer l’erreur en dernière valeur, la méthode renvoie `([]byte, error)`
- Nous utilisons le paquet `github.com/skip2/go-qrcode` pour effectuer la génération du code QR

 <br/>

### Enregistrer le service
Créer un service ne suffit pas : nous devons l’**enregistrer** auprès de l’application Wails afin qu’elle connaisse son existence et puisse générer ses bindings.

L’enregistrement s’effectue dans `main.go` lors de la création de votre application. Transmettez les instances de votre service à l’option `Services` :

```go {title="main.go" ins="7-9"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewService(NewQRService()),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**Voici ce qui se passe :**

- `application.NewService()` encapsule votre service afin que Wails puisse le gérer
- Nous appelons `NewQRService()` pour créer une instance de notre service
- Le service est ajouté au slice `Services` dans les options de l’application
- Wails va maintenant analyser les méthodes publiques de ce service afin de les mettre à la disposition du frontend

 <br/>

### Installer les dépendances
Nous avons référencé le paquet `github.com/skip2/go-qrcode` dans notre code, mais nous ne l’avons pas encore téléchargé. Go doit connaître cette dépendance et la télécharger dans votre projet.

Depuis le répertoire de votre projet, exécutez cette commande dans votre terminal :

```bash
go mod tidy
```

**Voici ce qui se passe :**

- `go mod tidy` recherche les instructions d’importation dans vos fichiers Go
- La commande télécharge tous les paquets manquants, tels que `go-qrcode`, et les ajoute à `go.mod`
- Elle supprime également toutes les dépendances qui ne sont plus utilisées
- Votre projet dispose ainsi de tout le code nécessaire pour être compilé correctement

Une sortie devrait indiquer que le paquet de génération de codes QR a été téléchargé et ajouté à votre projet.

 <br/>

### Générer les bindings
Pour appeler ces méthodes depuis votre frontend, nous devons générer les bindings. Pour cela, exécutez `wails generate bindings` dans le répertoire racine de votre projet.

@note{type="info"}
Lors de la toute première exécution de cette commande dans un projet, le générateur de bindings analyse minutieusement votre code et ses dépendances. Cette opération peut parfois prendre un peu plus de temps que prévu ; les exécutions suivantes seront toutefois beaucoup plus rapides.

@end

Après avoir exécuté cette commande, vous devriez voir dans votre terminal une sortie semblable à la suivante :

```bash
 % wails3 generate bindings
 INFO  Processed: 337 Packages, 1 Service, 1 Method, 0 Enums, 0 Models in 740.196125ms.
 INFO  Output directory: /Users/leaanthony/myproject/frontend/bindings
```

Vous devriez constater qu’un nouveau répertoire nommé `bindings` a été créé dans le répertoire du frontend :

```bash
frontend/
└── bindings
    └── changeme
        ├── index.js
        └── qrservice.js
```

@note{type="tip" title="Conseil de pro"}
Lorsque vous compilez votre application avec `wails3 build`, les bindings sont automatiquement générés et tenus à jour.

@end

 <br/>

### Comprendre les bindings
Examinons les bindings générés dans `bindings/changeme/qrservice.js` :

```js {title="bindings/changeme/qrservice.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 /**
  * QRService handles QR code generation
  * @module
  */

 // eslint-disable-next-line @typescript-eslint/ban-ts-comment
 // @ts-ignore: Unused imports
 import {Call as $Call, Create as $Create} from "@wailsio/runtime";

 /**
  * Generate creates a QR code from the given text
  * @param {string} text
  * @param {number} size
  * @returns {Promise<string> & { cancel(): void }}
  */
 export function Generate(text, size) {
     let $resultPromise = /** @type {any} */($Call.ByID(3576998831, text, size));
     let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
         return $Create.ByteSlice($result);
     }));
     $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
     return $typingPromise;
 }
```

Nous pouvons constater que des bindings sont générés pour la méthode `Generate`. Les noms des paramètres sont conservés, tout comme les commentaires. Une documentation JSDoc est également générée pour la méthode afin de fournir des informations de type à votre IDE.

@note{type="info"}
Il n’est pas nécessaire de comprendre les bindings générés dans les moindres détails, mais il est important de savoir comment ils fonctionnent.

@end

Les bindings fournissent :

- Des fonctions équivalentes à vos méthodes Go
- Une conversion automatique entre les types Go et JavaScript
- Des opérations asynchrones fondées sur les promesses
- Des informations de type sous forme de commentaires JSDoc

@note{type="tip" title="TypeScript"}
Le générateur de bindings permet également de générer des bindings TypeScript. Pour cela, exécutez `wails3 generate bindings -ts`.

@end

Le service généré est réexporté par un fichier `index.js` :

```js {title="bindings/changeme/index.js"}
 // @ts-check
 // Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
 // This file is automatically generated. DO NOT EDIT

 import * as QRService from "./qrservice.js";
 export {
     QRService
 };
```

Vous pouvez ensuite y accéder au moyen du chemin d’importation simplifié `./bindings/changeme`, constitué uniquement du chemin de votre paquet Go, sans indiquer de nom de fichier.

@note{type="info"}
Les chemins d’importation simplifiés ne sont disponibles qu’avec les outils de regroupement frontend. Si vous préférez un frontend standard qui n’utilise aucun outil de regroupement, vous devrez importer manuellement `index.js` ou `qrservice.js`.

@end

 <br/>

### Utiliser les bindings dans le frontend
Nous pouvons maintenant appeler notre service Go depuis JavaScript ! Les bindings générés rendent cette opération simple et sûre du point de vue des types.

Mettez à jour `frontend/src/main.js` pour utiliser les nouveaux bindings :

```js {title="frontend/src/main.js"}
 import { QRService } from './bindings/changeme';

 async function generateQR() {
     const text = document.getElementById('text').value;
     if (!text) {
         alert('Please enter some text');
         return;
     }

     try {
         // Generate QR code as base64
         const qrCodeBase64 = await QRService.Generate(text, 256);

         // Display the QR code
         const qrDiv = document.getElementById('qrcode');
         qrDiv.src = `data:image/png;base64,${qrCodeBase64}`;

     } catch (err) {
         console.error('Failed to generate QR code:', err);
         alert('Failed to generate QR code: ' + err);
     }
 }

 export function initializeQRGenerator() {
     const button = document.getElementById('generateButton');
     button.addEventListener('click', generateQR);
 }
```

**Voici ce qui se passe :**

- Nous importons `QRService` depuis les bindings générés
- `QRService.Generate()` appelle notre méthode Go ; cet appel renvoie une promesse, nous utilisons donc `await`
- La méthode Go renvoie `[]byte`, que Wails convertit automatiquement en chaîne base64 pour JavaScript
- Nous créons une URL de données avec la chaîne base64 afin d’afficher l’image PNG
- Le bloc `try/catch` traite toutes les erreurs provenant du code Go, telles qu’une entrée non valide
- Si notre code Go renvoie une erreur, la promesse est rejetée et nous interceptons l’erreur ici

Mettez maintenant à jour `index.html` afin d’utiliser les nouveaux bindings dans la fonction `initializeQRGenerator` :

```html {title="frontend/src/index.html"}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>QR Code Generator</title>
            <style>
                body {
                font-family: Arial, sans-serif;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                height: 100vh;
                margin: 0;
            }
                #qrcode {
                margin-bottom: 20px;
                width: 256px;
                height: 256px;
                display: flex;
                align-items: center;
                justify-content: center;
            }
                #controls {
                display: flex;
                gap: 10px;
            }
                #text {
                padding: 5px;
            }
                #generateButton {
                padding: 5px 10px;
                cursor: pointer;
            }
            </style>
</head>
<body>
<img id="qrcode"/>
<div id="controls">
    <input type="text" id="text" placeholder="Enter text">
        <button id="generateButton">Generate QR Code</button>
</div>

<script type="module">
    import { initializeQRGenerator } from './main.js';
    document.addEventListener('DOMContentLoaded', initializeQRGenerator);
</script>
</body>
</html>
```

Exécutez `wails3 dev` pour démarrer le serveur de développement. L’application devrait s’ouvrir après quelques secondes.

Saisissez du texte et cliquez sur le bouton « Générer le code QR ». Un code QR devrait apparaître au centre de la page :

![Code QR](/assets/qr1.png)

 <br/>

 <br/>

### Autre approche : gestionnaire HTTP
Jusqu’ici, nous avons abordé les sujets suivants :

- Création d’un service
- Génération des bindings
- Utilisation des bindings dans notre code frontend

**Pourquoi utiliser un gestionnaire HTTP ?**

Les liaisons de méthodes conviennent parfaitement aux opérations sur les données, mais il existe une autre approche pour servir des fichiers, des images ou d’autres médias. Au lieu de tout convertir en base64 et de l’envoyer par l’intermédiaire des liaisons, vous pouvez faire fonctionner votre service comme un mini-serveur web.

Cette approche est utile dans les cas suivants :

- Vous servez des images, des vidéos ou des fichiers volumineux
- Vous souhaitez utiliser des balises HTML `<img>` ou `<video>` standard avec des attributs `src`
- Vous avez besoin d’accéder directement aux ressources par leur URL

Si votre service implémente la méthode standard `ServeHTTP(w http.ResponseWriter, r *http.Request)` de Go, Wails peut le rendre accessible en tant que point de terminaison HTTP. Étendons notre service de codes QR pour prendre en charge cette possibilité :

```go {title="qrservice.go" ins="4-5,37-65"}
package main

import (
    "net/http"
    "strconv"

    "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation
type QRService struct {
    // We can add state here if needed
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**Voici ce qui se passe :**

- `ServeHTTP` est l’interface standard de Go pour traiter les requêtes HTTP
- Nous analysons les paramètres de requête de l’URL (`?text=hello&size=256`)
- Nous appelons notre méthode `Generate()` existante pour créer le code QR
- Nous définissons le type de contenu sur `image/png` afin que les navigateurs sachent qu’il s’agit d’une image
- Nous écrivons directement les octets PNG bruts dans la réponse : aucune conversion en base64 n’est nécessaire !

Mettez maintenant à jour `main.go` pour indiquer la route sur laquelle le service de codes QR doit être accessible :

```go {title="main.go" ins="8-10"}
 func main() {

     app := application.New(application.Options{
         Name:        "myproject",
         Description: "A demo of using raw HTML & CSS",
         LogLevel:    slog.LevelDebug,
         Services: []application.Service{
             application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                 Route: "/qrservice",
             }),
         },
         Assets: application.AssetOptions{
             Handler: application.AssetFileServerFS(assets),
         },
         Mac: application.MacOptions{
             ApplicationShouldTerminateAfterLastWindowClosed: true,
         },
     })

     app.Window.NewWithOptions(application.WebviewWindowOptions{
         Title:  "myproject",
         Width:  600,
         Height: 400,
     })

     // Run the application. This blocks until the application has been exited.
     err := app.Run()

     // If an error occurred while running the application, log it and exit.
     if err != nil {
         log.Fatal(err)
     }
 }
```

**Voici ce qui se passe :**

- Nous ajoutons `application.ServiceOptions` pour configurer la manière dont le service est exposé
- `Route: "/qrservice"` rend le gestionnaire HTTP accessible à l’adresse `/qrservice`
- Désormais, toute requête envoyée à `/qrservice?text=hello` appellera notre méthode `ServeHTTP`
- Si `Route` n’est pas défini, la fonctionnalité de gestionnaire HTTP est désactivée

@note{type="info"}
Si vous ne définissez pas explicitement l’option `Route`, le gestionnaire HTTP ne sera pas accessible depuis le frontend.

@end

Enfin, mettez à jour `main.js` pour utiliser simplement l’attribut `src` d’une image au lieu d’un encodage en base64 :

```js {title="frontend/src/main.js"}
async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Make the image source the path to the QR code service, passing the text
    img.src = `/qrservice?text=${encodeURIComponent(text)}`
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**Voici ce qui se passe :**

- Nous avons supprimé l’importation et l’appel à `await QRService.Generate()`
- À la place, nous définissons simplement `img.src` pour qu’il pointe vers notre point de terminaison HTTP
- `encodeURIComponent()` échappe de manière sûre les caractères spéciaux dans l’URL
- Le navigateur envoie automatiquement une requête HTTP GET lorsque nous définissons `src`
- Cette approche est plus simple et plus efficace pour les images : aucune conversion en base64 n’est nécessaire !

En relançant l’application, vous devriez obtenir le même code QR :

![Code QR](/assets/qr1.png)

 <br/>

 <br/>

### Prise en charge des configurations dynamiques
**Le problème des routes codées en dur :**

Dans l’exemple précédent, nous avons utilisé une route `/qrservice` codée en dur dans notre code JavaScript. Cela crée un couplage étroit entre votre configuration Go et votre code frontend.

Si vous modifiez l’option `Route` dans `main.go` sans mettre à jour `main.js`, l’application ne fonctionnera plus :

```go {title="main.go" ins="3"}
        // ...
            application.NewServiceWithOptions(NewQRService(), application.ServiceOptions{
                Route: "/services/qr",
            }),
        // ...
```

Les routes codées en dur conviennent aux applications simples, mais elles rendent votre code fragile et plus difficile à maintenir.

**La solution : une configuration dynamique**

Les liaisons de méthodes et les gestionnaires HTTP peuvent fonctionner ensemble ! Nous pouvons utiliser les liaisons pour indiquer au frontend quelle route utiliser, ce qui rend la configuration dynamique et élimine le chemin codé en dur.

Voici comment cela fonctionne :

1. La méthode de cycle de vie `ServiceStartup` s’exécute au démarrage de votre application
2. Nous enregistrons la route configurée dans les options
3. Nous ajoutons une méthode `URL()` que le frontend peut appeler pour obtenir la route appropriée
4. Désormais, le frontend demande sa route au service Go au lieu de la deviner

Commencez par implémenter l’interface `ServiceStartup` et ajouter une nouvelle méthode `URL` :

```go {title="qrservice.go" ins="4,6,10,15,23-27,46-55"}
package main

import (
    "context"
    "net/http"
    "net/url"
    "strconv"

    "github.com/skip2/go-qrcode"
    "github.com/wailsapp/wails/v3/pkg/application"
)

// QRService handles QR code generation
type QRService struct {
    route string
}

// NewQRService creates a new QR service
func NewQRService() *QRService {
    return &QRService{}
}

// ServiceStartup runs at application startup.
func (s *QRService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    s.route = options.Route
    return nil
}

// Generate creates a QR code from the given text
func (s *QRService) Generate(text string, size int) ([]byte, error) {
    // Generate the QR code
    qr, err := qrcode.New(text, qrcode.Medium)
    if err != nil {
        return nil, err
    }

    // Convert to PNG
    png, err := qr.PNG(size)
    if err != nil {
        return nil, err
    }

    return png, nil
}

// URL returns an URL that may be used to fetch
// a QR code with the given text and size.
// It returns an error if the HTTP handler is not available.
func (s *QRService) URL(text string, size int) (string, error) {
    if s.route == "" {
        return "", errors.New("http handler unavailable")
    }

    return fmt.Sprintf("%s?text=%s&size=%d", s.route, url.QueryEscape(text), size), nil
}

func (s *QRService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract the text parameter from the request
    text := r.URL.Query().Get("text")
    if text == "" {
        http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
        return
    }
    // Extract Size parameter from the request
    sizeText := r.URL.Query().Get("size")
    if sizeText == "" {
        sizeText = "256"
    }
    size, err := strconv.Atoi(sizeText)
    if err != nil {
        http.Error(w, "Invalid 'size' parameter", http.StatusBadRequest)
        return
    }

    // Generate the QR code
    qrCodeData, err := s.Generate(text, size)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Write the QR code data to the response
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrCodeData)
}
```

**Voici ce qui se passe :**

- Nous avons ajouté un champ `route` pour stocker la route configurée provenant de `ServiceStartup`
- `ServiceStartup(ctx, options)` est appelée au démarrage de l’application : nous y enregistrons la route
- La méthode `URL()` construit l’URL complète avec les paramètres de requête
- Si aucune route n’est configurée (la route est vide), nous renvoyons une erreur
- `url.QueryEscape()` encode le texte de manière sûre pour son utilisation dans une URL
- Cette méthode sera accessible au frontend par l’intermédiaire des liaisons

Mettez maintenant à jour `main.js` pour utiliser la méthode `URL` à la place d’un chemin codé en dur :

```js {title="frontend/src/main.js" ins="1,11-12"}
import { QRService } from "./bindings/changeme";

async function generateQR() {
    const text = document.getElementById('text').value;
    if (!text) {
        alert('Please enter some text');
        return;
    }

    const img = document.getElementById('qrcode');
    // Invoke the URL method to obtain an URL for the given text.
    img.src = await QRService.URL(text, 256);
}

export function initializeQRGenerator() {
    const button = document.getElementById('generateButton');
    if (button) {
        button.addEventListener('click', generateQR);
    } else {
        console.error('Generate button not found');
    }
}
```

**Voici ce qui se passe :**

- Nous importons `QRService` afin d’utiliser de nouveau les liaisons
- Au lieu de coder `/qrservice` en dur, nous appelons `await QRService.URL(text, 256)`
- Le service Go construit l’URL avec la route et les paramètres appropriés
- Désormais, si vous modifiez la route dans `main.go`, le frontend utilise automatiquement la nouvelle route
- Plus besoin de synchroniser manuellement la configuration Go et le code frontend !

Le fonctionnement devrait être identique à celui de l’exemple précédent, mais modifier la route du service dans `main.go` ne provoquera plus de dysfonctionnement du frontend.

@note{type="info"}
Si une méthode Go renvoie une erreur non nil, la promesse côté JS sera rejetée et les instructions await lèveront une exception.

@end

 <br/>

@end
