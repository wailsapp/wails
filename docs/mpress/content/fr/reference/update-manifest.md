---
title: "Protocole de manifeste de mise à jour"
description: "Le protocole JSON ouvert utilisé par les applications Wails pour détecter et vérifier leurs propres mises à jour, qui peut être servi depuis n’importe quel hébergeur de fichiers statiques ou serveur de mise à jour dynamique."
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Le protocole Wails Update Manifest est un petit contrat JSON ouvert entre une application Wails et une source de mises à jour. Tout service capable de servir un fichier JSON via HTTPS peut fournir des mises à jour Wails : un compartiment S3, GitHub Pages, un CDN ou un serveur de mise à jour dynamique qui conditionne les versions à la possession d’une licence.

Le côté client est fourni avec le framework sous la forme du fournisseur `endpoint` (`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`). Cette page constitue la référence du format d’échange pour toute personne qui implémente le côté serveur.

## Objectifs de conception

1. **Adapté à l’hébergement statique.** Un seul fichier manifeste par canal, répertoriant l’artefact de chaque plateforme, constitue une implémentation complète. Aucun code serveur n’est nécessaire.
2. **Adapté aux serveurs dynamiques.** Le client envoie `platform`, `arch`, `version` et `channel` lors de chaque vérification. Un serveur peut ainsi répondre avec exactement un artefact, appliquer des règles de licence ou renvoyer `204 No Content` lorsque l’appelant est à jour.
3. **Priorité à la vérification.** Le manifeste contient les sommes de contrôle et les signatures propres à chaque artefact, et le programme de mise à jour de Wails les vérifie à l’aide d’une clé publique épinglée dans le binaire de l’application lors de la compilation. La source de mise à jour ne choisit jamais sa propre racine de confiance.

## La requête

Le client envoie une requête `GET` à l’URL de manifeste configurée avec `Accept: application/json`, ainsi que tous les en-têtes configurés par l’application (par exemple `Authorization: License <key>`).

L’URL peut contenir des espaces réservés, que le client remplace à chaque vérification :

| Espace réservé | Remplacé par |
| --- | --- |
| `{{platform}}` | Le système d’exploitation en cours d’exécution sous la forme d’une valeur Go `GOOS` (`darwin`, `windows`, `linux`) |
| `{{arch}}` | L’architecture en cours d’exécution sous la forme d’une valeur Go `GOARCH` (`amd64`, `arm64`, ...) |
| `{{version}}` | La version actuellement installée |
| `{{channel}}` | Le canal de publication configuré, le cas échéant |

Toute valeur parmi les quatre qui n’est pas utilisée par un espace réservé est ajoutée comme paramètre de requête portant le même nom (`channel` uniquement s’il est configuré). Les deux configurations suivantes sont donc valides et équivalentes :

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

Les hébergeurs statiques ignorent simplement les paramètres de requête qu’ils reçoivent.

## La réponse

| État | Signification |
| --- | --- |
| `200 OK` | Un manifeste suit. Le client détermine s’il s’agit d’une mise à niveau. |
| `204 No Content` | Le serveur a comparé les versions et l’appelant est à jour. |
| `404 Not Found` | Aucune version publiée (traité de la même manière qu’un système à jour). |
| Toute autre valeur | Une erreur. Le programme de mise à jour passe au fournisseur configuré suivant. |

Un corps `200` est un document manifeste :

```json
{
  "schemaVersion": 1,
  "version": "2.1.0",
  "channel": "stable",
  "name": "Summer Release",
  "notes": "## What's new\n\n- Faster startup\n- New themes",
  "publishedAt": "2026-07-03T10:00:00Z",
  "artifacts": [
    {
      "url": "MyApp-2.1.0-darwin-arm64.zip",
      "platform": "darwin",
      "arch": "arm64",
      "filetype": "zip",
      "size": 8388608,
      "digestAlgo": "sha512",
      "digest": "base64-encoded digest bytes",
      "signatureAlgo": "ed25519ph",
      "signature": "base64-encoded signature bytes"
    },
    {
      "url": "MyApp-2.1.0-windows-amd64.zip",
      "platform": "windows",
      "arch": "amd64",
      "filetype": "zip",
      "size": 9437184,
      "digestAlgo": "sha512",
      "digest": "...",
      "signatureAlgo": "ed25519ph",
      "signature": "..."
    }
  ]
}
```

### Champs de premier niveau

| Champ | Type | Obligatoire | Remarques |
| --- | --- | --- | --- |
| `schemaVersion` | int | non | Version du protocole. Une omission signifie `1`. Les clients rejettent les valeurs correspondant à des versions plus récentes que celles qu’ils prennent en charge. |
| `version` | string | **oui** | Conforme à la spécification SemVer 2.0.0, avec ou sans `v` initial. |
| `channel` | string | non | À titre informatif. Un client configuré pour un autre canal considère que le manifeste ne propose aucune mise à jour. |
| `name` | string | non | Titre de la version lisible par l’utilisateur, affiché dans la fenêtre de mise à jour. |
| `notes` | string | non | Notes de version au format Markdown, affichées dans la fenêtre de mise à jour. |
| `publishedAt` | string | non | Horodatage RFC 3339. |
| `artifacts` | array | **oui** | Une entrée par artefact téléchargeable. L’ordre exprime la préférence de l’éditeur. |
| `metadata` | object | non | Données clé/valeur de forme libre, transmises à l’application. |

Les clients ignorent les champs inconnus, ce qui permet aux serveurs d’ajouter leurs propres champs sans rompre la compatibilité. Les ajouts propres au serveur doivent être placés dans `metadata`.

### Champs des artefacts

| Champ | Type | Obligatoire | Remarques |
| --- | --- | --- | --- |
| `url` | chaîne | **oui** | URL absolue ou relative à celle du manifeste. `http(s)` uniquement. |
| `platform` | chaîne | non | Valeur Go `GOOS`. Les alias courants (`macos`, `win`, etc.) sont acceptés. Une valeur vide correspond à toutes les plateformes. |
| `arch` | chaîne | non | Valeur Go `GOARCH`. Les alias courants (`x86_64`, `aarch64`, etc.) sont acceptés. Une valeur vide correspond à toutes les architectures. |
| `filename` | chaîne | non | Prend par défaut la valeur du dernier segment du chemin de `url`. |
| `filetype` | chaîne | non | Prend par défaut la valeur de l’extension du nom de fichier. |
| `size` | entier | non | Nombre d’octets, utilisé pour indiquer la progression du téléchargement. |
| `digestAlgo` / `digest` | chaîne / base64 | non | `sha256` ou `sha512`. |
| `signatureAlgo` / `signature` | chaîne / base64 | non | `ed25519`, `ed25519ph` ou `ecdsa-p256`. `signatureAlgo` est obligatoire dès que `signature` est présent. Consultez le [guide du programme de mise à jour](/guides/updater/#cryptographic-verification) pour savoir ce que signe chaque algorithme. |

Le client sélectionne le **premier** artefact dont les valeurs `platform` et `arch` correspondent au système en cours d’exécution. Les valeurs en base64 sont acceptées avec ou sans remplissage.

### Comparaison des versions

Le client détermine toujours si le manifeste constitue une mise à niveau selon l’ordre de priorité SemVer 2.0.0 : la valeur `version` du manifeste doit correspondre à une version strictement plus récente que la version installée. L’hébergement statique fonctionne ainsi correctement sans configuration particulière : le manifeste décrit toujours la dernière version publiée et les clients à jour ne font tout simplement rien. `204` reste toutefois disponible pour permettre aux serveurs dynamiques d’économiser de la bande passante.

## Vérification et confiance

Les sommes de contrôle et les signatures sont incluses dans le manifeste, mais pas la racine de confiance : les signatures sont vérifiées à l’aide de la clé publique que l’application a épinglée au moment de la compilation via `updater.Config.PublicKey`. Une source de mise à jour compromise ou substituée ne peut pas fournir sa propre clé. La vérification échoue de manière sûre lorsqu’un artefact comporte une signature alors que l’application n’a épinglé aucune clé, lorsque la signature ne déclare aucun `signatureAlgo` ou lorsqu’elle ne peut pas être décodée : les clients ne se rabattent jamais silencieusement sur une vérification limitée à l’empreinte.

Les artefacts qui ne comportent qu’une empreinte sont installés après vérification de celle-ci. Cette vérification protège contre la corruption, mais la résistance aux altérations repose sur TLS et sur le fait que l’hôte lui-même ne soit pas compromis. Fournissez des signatures pour tout élément sensible du point de vue de la sécurité.

Quelques lignes de Go suffisent pour signer un artefact avec le mécanisme `ed25519ph` du framework :

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

En pratique, vous aurez rarement à écrire ce code : la CLI s’en charge pour vous.

## Publication avec la CLI wails3

Le groupe de commandes `wails3 updater` couvre l’intégralité du pipeline de publication. Trois commandes suffisent pour publier une version :

```bash
# Once per application: create the signing keypair.
wails3 updater genkey
# updater.key      keep secret (CI secret store), signs every release
# updater.key.pub  embed in the app and pass as updater.Config.PublicKey

# Per release: digest, sign and describe every artifact in one manifest.
wails3 updater manifest -version 2.1.0 -channel stable \
    -key updater.key -notes-file notes.md \
    -url-prefix "https://cdn.example.com/myapp/2.1.0" \
    bin/updates/

# Before uploading: re-verify the files exactly as a shipped app would.
wails3 updater verify -manifest manifest.json -publickey updater.key.pub
```

`manifest` accepte des fichiers ou des répertoires (les clés, `.json` ainsi que les fichiers annexes de sommes de contrôle et de notes sont automatiquement ignorés), transmet chaque artefact en flux à SHA-512, signe l’empreinte avec Ed25519ph lorsque `-key` est fourni et déduit `platform` et `arch` à partir de noms de fichiers conventionnels tels que `MyApp-2.1.0-darwin-arm64.zip`. Les alias courants comme `macOS`, `win64`, `x86_64` et `aarch64` sont reconnus. Un avertissement est affiché pour toute valeur que la commande ne parvient pas à déduire ; l’artefact correspond alors à toutes les plateformes. Omettez `-url-prefix` pour générer des URL relatives, puis téléversez le manifeste à côté des artefacts.

`verify` renvoie un code de sortie différent de zéro en cas de divergence, ce qui en fait un point de contrôle naturel dans l’intégration continue entre la compilation et la publication. Pour les serveurs qui assemblent eux-mêmes les manifestes, `wails3 updater sign -key updater.key <files...>` affiche les champs `digest`/`signature` de chaque fichier sous forme de JSON prêt à être fusionné dans votre propre document.

## Authentification

L’authentification relève du serveur ; le protocole ne fait que transporter les en-têtes. Le client renvoie ses en-têtes configurés avec chaque requête de manifeste. Lors du téléchargement d’un artefact, l’en-tête `Authorization` n’est envoyé que si l’URL de l’artefact se trouve sur le même hôte que le manifeste et ne passe pas de `https` à `http`. Il est supprimé lors de toute redirection vers une autre origine ou entraînant une telle rétrogradation. Les identifiants ne sont donc jamais divulgués à un CDN ou à un stockage d’objets et ne transitent jamais en clair.

Exemple avec contrôle par licence, qui s’associe naturellement aux services de gestion de licences hébergés :

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## Configuration du client

Consultez le [guide du programme de mise à jour](/guides/updater/#providers) pour obtenir la référence complète de `endpoint.Config` et découvrir comment le fournisseur s’intègre aux chaînes de repli aux côtés des fournisseurs GitHub, keygen.sh et AppCast.
