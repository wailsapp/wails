---
title: "Lire des fichiers audio et vidéo locaux"
description: "Lire des médias intégrés sous Linux avec des URL blob de taille limitée et les libérer après utilisation."
slug: "guides/linux-media"
sourcePath: "guides/linux-media.md"
---

Utilisez `Media.SetSource` pour lire de courts clips locaux dans une application Wails v3. Sous Linux, WebKitGTK délègue la lecture multimédia à GStreamer, qui ne peut pas charger directement les URL `wails://`. Cette fonction reçoit le clip via un flux Wails et attribue une URL blob au lecteur.

Sur ordinateur, les transferts utilisent le transport de ressources Wails existant sans ouvrir de socket d’écoute. L’API fonctionne avec les éléments audio et vidéo. Les médias HTTP/HTTPS ordinaires peuvent utiliser directement le `src` du lecteur.

## Enregistrer vos fichiers multimédias

Exposez un système de fichiers contenant les clips que votre frontend peut lire, puis enregistrez le gestionnaire multimédia sur un flux nommé :

```go
import (
    "embed"
    "io/fs"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
    "github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed clips
var clips embed.FS

func registerMedia(app *application.App) {
    mediaFiles, err := fs.Sub(clips, "clips")
    if err != nil {
        log.Fatal(err)
    }
    handler, err := media.NewHandler(mediaFiles, 32 << 20) // 32 MiB per file
    if err != nil {
        log.Fatal(err)
    }
    app.HandleStream("media", handler)
}
```

Appelez `registerMedia(app)` après avoir créé l’application et avant d’appeler `app.Run()`.

Pour les fichiers sur disque, utilisez `os.OpenRoot(directory)` et transmettez `root.FS()` à `media.NewHandler`. Gardez la racine ouverte jusqu’au retour de `app.Run()`, puis fermez-la. Choisissez un répertoire ne contenant que les fichiers accessibles au frontend ; `os.Root` limite l’accès même si un lien symbolique pointe hors de ce répertoire.

## Charger un clip

Créez un lecteur :

```html
<video id="player" controls></video>
```

Pour un frontend utilisant le runtime npm :

```javascript
import { Media } from '@wailsio/runtime';

const player = document.getElementById('player');

try {
    await Media.SetSource(player, 'media', 'welcome.mp4');
} catch (error) {
    if (error.name !== 'AbortError') {
        console.error('Could not load the clip:', error);
    }
}
```

Pour une application utilisant le runtime intégré, remplacez l’import par :

```javascript
import { Media } from '/wails/runtime.js';
```

Le deuxième argument est le nom du flux enregistré. Le troisième est un chemin de fichier séparé par des barres obliques, relatif à la racine de son système de fichiers, comme `welcome.mp4` ou `tutorials/intro.mp4`. Ce n’est ni une URL ni un chemin du système d’exploitation.

La promesse est résolue lorsque la source est attribuée. Le lecteur la décode ensuite. Gérez l’événement `error` du lecteur pour détecter les codecs non pris en charge et appelez `player.play()` depuis une interaction utilisateur si vous devez démarrer la lecture vous-même.

## Remplacer ou libérer un clip

Appelez à nouveau `Media.SetSource` pour changer de clip. Cela annule tout chargement précédent en attente pour ce lecteur, afin qu’une réponse lente ne remplace pas la dernière sélection. Le clip précédent reste disponible jusqu’au chargement réussi de son remplaçant. Son URL blob est alors révoquée.

Appelez `Media.ClearSource(player)` lorsque vous fermez un lecteur ou démontez son composant :

```javascript
Media.ClearSource(player);
```

Cela annule tout chargement en attente, réinitialise le lecteur et libère l’URL blob. Appelez cette fonction avant de retirer l’élément du document. Un composant de framework doit l’appeler dans son hook de démontage ou de destruction. Retirer uniquement l’élément ne libère pas l’URL blob. Utilisez systématiquement ces fonctions pour gérer la source du lecteur.

Pour le balisage `<video><source ...></video>`, transmettez le nom de fichier choisi à `Media.SetSource(video, "media", name)`. Cette fonction définit le `src` du lecteur parent, qui prend le pas sur ses enfants `<source>`. Effacer cet attribut permet au navigateur de considérer à nouveau ces enfants ; utilisez un lecteur vide lorsque vous gérez toutes les sources via cette API.

Les objets audio détachés fonctionnent également :

```javascript
const sound = new Audio();
await Media.SetSource(sound, 'media', 'notification.mp3');
sound.addEventListener('ended', () => Media.ClearSource(sound), { once: true });
// Call sound.play() from an appropriate user interaction.
// Also clear it if playback is cancelled or the owning component is disposed.
```

## Limiter les téléchargements et annuler le chargement

La limite par défaut est de **32 MiB par source**. Vous pouvez choisir une limite entière positive en octets, inférieure ou supérieure, dans la limite configurée côté Go :

```javascript
const controller = new AbortController();
const loading = Media.SetSource(player, 'media', 'welcome.mp4', {
    maxBytes: 8 * 1024 * 1024,
    signal: controller.signal,
});

// Call controller.abort() to cancel this load.
await loading;
```

Un fichier trop volumineux entraîne un rejet avec `RangeError`. Le gestionnaire Go vérifie la taille du fichier avant de lire son contenu et transfère au maximum la plus petite des limites configurées côté Go et côté frontend. Le frontend vérifie également la taille reçue et rejette les transferts incomplets. Une annulation entraîne un rejet avec `AbortError` ou avec le motif transmis à `AbortController.abort(reason)`.

Les fichiers sont transmis par trames de 64 KiB. Le transport de flux Wails pour ordinateur limite les réponses aux requêtes d’interrogation à 1 MiB ; la mise en mémoire tampon des réponses complètes par WebView2 ne stocke donc pas tout le fichier multimédia dans une seule réponse. Les files de flux disposent de leur propre mécanisme borné de contre-pression. Ces limites de transport ne changent pas le comportement de chargement d’un blob complet par la fonction du lecteur.

**Le fichier entier est téléchargé avant la lecture.** La limite en octets s’applique à chaque fichier, pas à la mémoire totale de l’application. Plusieurs lecteurs, le clip précédent pendant son remplacement, la construction du blob et les médias décodés peuvent consommer de la mémoire supplémentaire. Le transfert commence à l’appel, quel que soit le paramètre `preload` du lecteur. Appelez la fonction lorsque l’utilisateur choisit de charger un clip.

Pour les gros fichiers locaux, cette fonction n’est pas une solution de lecture en continu. Augmenter la limite augmente aussi la consommation mémoire. Les médias déjà hébergés en HTTP/HTTPS doivent utiliser le chargement multimédia natif afin que le navigateur puisse lire en continu et rechercher une position à l’aide de requêtes par plages.

## Dépanner la lecture sous Linux

- Si la lecture locale directe indique **No URI handler implemented for "wails"**, chargez le clip avec `Media.SetSource`. Cela concerne aussi bien la pile GTK4 par défaut que la pile historique `-tags gtk3`.
- Si le chargement réussit mais que le décodage échoue, vérifiez les codecs GStreamer installés sur le système cible. Le MP4 nécessite généralement la prise en charge de la vidéo H.264 et de l’audio AAC ; le MP3 nécessite un décodeur MP3. Testez les formats distribués sur les distributions prises en charge.
- Si le transfert échoue, vérifiez le nom du flux enregistré, le nom de fichier relatif et les permissions du système de fichiers. Modifier un fichier pendant son transfert peut entraîner un transfert incomplet ; réessayez après la fin de son écriture.
- Si votre application définit une politique de sécurité du contenu, autorisez l’origine des ressources Wails dans `connect-src` et `blob:` dans `media-src`. Pour une politique exclusivement locale, ces directives peuvent être `connect-src 'self'; media-src 'self' blob:`. Conservez les autres directives.
- Si un fichier dépasse la limite, choisissez un clip plus court ou moins volumineux, ou définissez une limite explicite adaptée au budget mémoire de l’application.

Exécutez l’[exemple audio-video](https://github.com/wailsapp/wails/tree/master/v3/examples/audio-video) avec `go run .` pour tester les exemples MP3 et MP4 intégrés sur votre machine.
