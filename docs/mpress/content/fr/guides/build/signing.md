---
title: "Signature du code"
description: "Guide pour signer vos applications Wails sur toutes les plateformes"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Signer le code de votre application

Ce guide explique comment signer vos applications Wails pour macOS, Windows et Linux. Wails v3 fournit des outils CLI intégrés pour la signature du code, la notarisation et la gestion des clés PGP.

- **macOS** — Signez et notarisez vos applications macOS
- **Windows** — Signez vos exécutables et paquets Windows
- **Linux** — Signez les paquets DEB et RPM avec des clés PGP

## Matrice de signature multiplateforme

Cette matrice indique ce que vous pouvez signer depuis chaque plateforme source :

| Format cible | Depuis Windows | Depuis macOS | Depuis Linux |
| --- | :---: | :---: | :---: |
| Windows EXE/DLL/MSIX/APPX, PFX exporté | ✅ | ✅ | ✅ |
| MSI Windows ou identité du magasin de certificats/fournisseur | ✅ | ❌ | ❌ |
| .app / Mach-O / DMG macOS, identité exportée | ✅ | ✅ | ✅ |
| Notarisation macOS, identifiants d’API d’équipe | ✅ | ✅ | ✅ |
| Identité du trousseau macOS ou profil de notarisation | ❌ | ✅ | ❌ |
| DEB Linux | ✅ | ✅ | ✅ |
| RPM Linux | ✅ | ✅ | ✅ |

@note{type="tip"}
Les identités Apple exportées et les identifiants d’API d’équipe permettent la signature et la notarisation portables. Les identités du trousseau macOS et les profils de notarisation nécessitent un Mac. Les restrictions propres à chaque format sont décrites ci-dessous.

@end

### Moteurs de signature

Wails sélectionne automatiquement le meilleur moteur de signature disponible :

| Plateforme | Moteur natif | Moteur multiplateforme |
| --- | --- | --- |
| Windows | `signtool.exe` pour les identités du magasin/fournisseur, les MSI ou les PFX RC2 hérités | Intégré pour les PFX RSA exportés et les formats PE/MSIX |
| macOS | `codesign` pour les identités du trousseau | Intégré pour les identités PKCS#12 exportées ou la signature ad hoc |
| Linux | N/D | Intégré |

Le type d’identifiant sélectionne le moteur de signature Apple sur chaque hôte. Les identités exportées utilisent Go ; les identités du trousseau macOS utilisent `codesign`. La récupération, la vérification et l’agrafage des tickets utilisent Go.

## Démarrage rapide

Le moyen le plus rapide de configurer la signature consiste à utiliser l’assistant de configuration :

```bash
wails3 setup
```

Cette opération écrit une configuration de signature partagée dans `~/.config/wails/defaults.yaml`, qui est **prise en compte lors de la signature sur toutes les plateformes** (voir [Ordre de priorité de la configuration](#ordre-de-priorit-de-la-configuration)). L’étape de signature :

- Affiche les outils de configuration permettant de générer ou d’exporter des certificats et des clés, comme OpenSSL et GnuPG. Les identités Apple et Windows exportées utilisent le signataire intégré ; les identités du magasin de certificats/fournisseur Windows et les paquets MSI nécessitent Windows SignTool.
- Sous macOS, répertorie les certificats Developer ID de votre trousseau.
- Sous Linux, répertorie vos clés GPG et peut en **générer et exporter** une nouvelle.
- Sous Windows, peut **générer un certificat autosigné** (à des fins de test) avec OpenSSL.
- **Stocke les mots de passe de manière sécurisée dans le gestionnaire d’identifiants de votre système** (et non dans les Taskfiles).

@note{type="tip"}
Les mots de passe sont stockés dans le gestionnaire d’identifiants natif de votre système (Trousseau macOS, Gestionnaire d’informations d’identification Windows ou Secret Service sous Linux). Votre configuration de signature est ainsi sécurisée et fonctionne dans tous vos projets Wails.

@end

## Ordre de priorité de la configuration

Lorsque vous exécutez une tâche de signature, chaque option est résolue dans l’ordre suivant (la première correspondance l’emporte) :

1. Une option explicite transmise à `wails3 tool sign` (par exemple `--pgp-key`, `--certificate` ou `--identity`).
2. La **variable correspondante du Taskfile du projet** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, etc.).
3. La **configuration globale** dans `~/.config/wails/defaults.yaml` (écrite par `wails3 setup`).

Les variables du Taskfile sont donc des **remplacements facultatifs** : si une variable n’est pas définie, la clé, le certificat ou l’identité configuré globalement est utilisé. Si aucune de ces trois sources ne fournit de valeur, la commande de signature affiche une erreur claire qui vous indique comment la configurer.

## Configuration par projet

Pour configurer la signature d’un seul projet plutôt que globalement, exécutez l’assistant propre au projet depuis le répertoire de celui-ci. Il écrit la valeur `vars` dans les fichiers `build/<platform>/Taskfile.yml` du projet :

```bash
wails3 signing                                   # all detected platforms
wails3 signing --platform windows --platform linux
```

### Configuration manuelle

Vous pouvez également modifier manuellement les Taskfiles propres à chaque plateforme. Modifiez la section `vars` située en haut de chaque fichier :

@tabs
[macOS]
Modifiez `build/darwin/Taskfile.yml` :

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  # ENTITLEMENTS: "build/darwin/entitlements.plist"
```

Exécutez ensuite :

```bash
wails3 task darwin:sign           # Sign only
wails3 task darwin:sign:notarize  # Sign and notarize
```

[Windows]
Modifiez `build/windows/Taskfile.yml` :

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

Le mot de passe est récupéré dans le gestionnaire d’identifiants du système (exécutez `wails3 signing` pour le configurer).

Exécutez ensuite :

```bash
wails3 task windows:sign           # Sign executable
wails3 task windows:sign:installer # Sign NSIS installer
```

[Linux]
Modifiez `build/linux/Taskfile.yml` :

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

Le mot de passe est récupéré dans le gestionnaire d’identifiants du système (exécutez `wails3 signing` pour le configurer).

Exécutez ensuite :

```bash
wails3 task linux:sign:deb       # Sign DEB package
wails3 task linux:sign:rpm       # Sign RPM package
wails3 task linux:sign:packages  # Sign all packages
```

@end

Vous pouvez également inspecter directement l’état de la signature depuis le système :

```bash
# List available macOS code-signing identities
security find-identity -v -p codesigning

# List PGP keys (Linux package signing)
gpg --list-keys
```

Pour configurer interactivement la signature de chaque plateforme, utilisez l’assistant :

```bash
wails3 signing
```

## Signature du code sous macOS

### Prérequis

- Compte Apple Developer (99 $/an)
- Certificat Developer ID Application
- Pour une identité du trousseau macOS : outils de ligne de commande Xcode installés
- Pour la signature portable : une identité RSA PKCS#12 exportée ; aucun outil de signature Apple n’est requis

### Signature portable avec Wake

Wake lit directement `wails.yaml`. Il n’exécute pas les Taskfiles. Configurez une
identité Developer ID exportée et une clé d’API d’équipe App Store Connect :

```yaml
darwin:
  signing:
    certificate: signing/developer-id.p12
    credential: APPLE_CERT_PASSWORD
    timestamp_server: https://timestamp.digicert.com
  notarization:
    api_key: signing/AuthKey.p8
    key_id: ABCDEFGHIJ
    issuer_id: 00000000-0000-0000-0000-000000000000
```

`credential` désigne la variable d’environnement qui contient le mot de passe du certificat.
Conservez les deux fichiers de clé privée hors du contrôle de version. Activez la signature et
la notarisation dans le profil de publication décrit dans le
[guide de compilation YAML](/experimental/yaml-builds/).

Le signataire portable prend en charge les fichiers Mach-O Intel 64 bits et Apple Silicon, les binaires universels, les bundles d’app canoniques contenant du code imbriqué, les frameworks et les DMG UDIF de moins de 4 Gio. Il rejette les structures non prises en charge avant de publier une sortie. La signature ad hoc utilise `identity: "-"` et ne peut pas être notariée. Les identités PKCS#12 peuvent utiliser le chiffrement AES/PBKDF2 ou l’ancien chiffrement 3DES ; les exportations RC2 doivent être réexportées avec un algorithme de chiffrement pris en charge.

Les signatures notariées portables nécessitent explicitement un service d’horodatage SHA-2 conforme à la RFC 3161. Les anciennes signatures d’horodatage SHA-1 sont rejetées. L’envoi par API d’équipe prend en charge les téléversements uniques jusqu’à 5 Gio ; les clés d’API individuelles et les téléversements multiparties ne sont pas pris en charge. La vérification et l’agrafage du ticket utilisent le protocole de service Apple observé avec une racine Apple Root G3 épinglée, sans `xcrun stapler`. Gatekeeper reste l’autorité en matière de confiance de distribution et de révocation.

Les sections suivantes décrivent le processus existant fondé sur le trousseau macOS, utilisé par les anciens projets et par les configurations Wake qui sélectionnent une identité du trousseau.

### Identités de signature

Vérifiez les identités de signature disponibles :

```bash
security find-identity -v -p codesigning
```

Sortie :

```
Found 2 signing identities:

  Developer ID Application: Your Company (ABCD1234) [valid]
    Hash: ABC123DEF456...

  Apple Development: your@email.com (XYZ789) [valid]
    Hash: DEF789ABC123...
```

@note{type="tip"}
Pour distribuer votre application en dehors de l’App Store, vous avez besoin d’un certificat **Developer ID Application**.

@end

### Configuration

Modifiez `build/darwin/Taskfile.yml` et définissez les variables de signature :

```yaml
vars:
  SIGN_IDENTITY: "Developer ID Application: Your Company (TEAMID)"
  KEYCHAIN_PROFILE: "my-notarize-profile"
  ENTITLEMENTS: "build/darwin/entitlements.plist"
```

| Variable | Obligatoire | Description |
| --- | --- | --- |
| `SIGN_IDENTITY` | Oui | Votre identifiant Developer ID (par exemple, « Developer ID Application: Your Company (TEAMID) ») |
| `KEYCHAIN_PROFILE` | Pour la notarisation | Nom du profil du trousseau contenant les identifiants enregistrés |
| `ENTITLEMENTS` | Non | Chemin du fichier de droits |

Exécutez ensuite :

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Droits

Les droits déterminent les fonctionnalités auxquelles votre application peut accéder. Les applications Wails nécessitent généralement des droits différents en développement et en production :

- **Développement** : nécessite des droits pour le JIT, la mémoire non signée et le débogage
- **Production** : droits minimaux (uniquement l’accès réseau)

Utilisez l’assistant de configuration interactif pour générer les deux fichiers :

```bash
wails3 entitlements
```

Cette opération crée :

- `build/darwin/entitlements.dev.plist` — pour les builds de développement
- `build/darwin/entitlements.plist` — pour les builds de production/signés

**Préréglages disponibles :**\

| Préréglage | Description |
| --- | --- |
| Développement | JIT, mémoire non signée, débogage, réseau |
| Production | Réseau uniquement (minimal et le plus sécurisé) |
| Les deux | Crée les fichiers de développement et de production (recommandé) |
| App Store | Bac à sable activé avec accès au réseau et aux fichiers |
| Personnalisé | Choisissez les droits individuellement |

@note{type="note"}
La tâche `run` du Taskfile darwin utilise automatiquement `entitlements.dev.plist`. Les tâches `sign` utilisent `entitlements.plist` pour les builds de production.

@end

Dans les variables de votre Taskfile, définissez ensuite `ENTITLEMENTS` afin qu’il pointe vers le fichier approprié.

### Notarisation

Apple exige que toutes les applications distribuées soient notariées.

@steps
### **Enregistrez vos identifiants dans le trousseau** (configuration unique). Exécutez `wails3 signing`, qui vous demande ces valeurs et appelle `notarytool` en arrière-plan, ou appelez directement `notarytool` :
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Définissez KEYCHAIN_PROFILE dans votre Taskfile** afin qu’il corresponde au nom du profil ci-dessus.
### **Signez et notariez votre application** :
```bash
wails3 task darwin:sign:notarize
```

### **Vérifiez la notarisation** :
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
La notarisation prend généralement 1-2 minutes. Le ticket est automatiquement agrafé à votre application.

@end

## Signature de code Windows

### Prérequis

- Certificat de signature de code (fourni par DigiCert, Sectigo, etc.)
- Pour la signature portable : une identité RSA PFX exportée et chiffrée avec AES/PBKDF2 ou l’ancien algorithme 3DES ; aucun outil de signature externe n’est nécessaire
- Pour les identités du magasin de certificats ou d’un fournisseur, les paquets MSI ou les anciens PFX RC2 : Windows avec les outils de signature du SDK Windows sur `PATH`

### Génération d’un certificat autosigné (tests)

Si vous devez seulement tester le pipeline de signature, exécutez `wails3 setup`, ouvrez l’onglet **Windows** de l’étape de signature, puis choisissez **Générer un certificat autosigné**. OpenSSL crée alors un `.pfx` de signature de code et enregistre son chemin dans la configuration globale.

@note{type="caution"}
Les certificats autosignés sont réservés **aux tests et à la distribution interne** : ils déclenchent des avertissements SmartScreen pour les utilisateurs finaux. Les versions publiques nécessitent un certificat délivré par une autorité de certification de confiance.

@end

### Configuration

Modifiez `build/windows/Taskfile.yml` et définissez les variables de signature :

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Variable | Obligatoire | Description |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Remplacement | Chemin du fichier de certificat .pfx/.p12 (utilise la configuration globale par défaut si cette variable n’est pas définie) |
| `SIGN_THUMBPRINT` | Remplacement | Empreinte du certificat dans le magasin de certificats Windows (solution de remplacement à `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Non | URL du serveur d’horodatage (valeur par défaut : http://timestamp.digicert.com) |

@note{type="note"}
Ces variables sont des **valeurs de remplacement** facultatives. Si elles ne sont pas définies, le certificat configuré globalement avec `wails3 setup` est utilisé. Consultez [Priorité de la configuration](#ordre-de-priorit-de-la-configuration).

@end

@note{type="note"}
Le mot de passe du certificat est enregistré dans le trousseau de votre système, et non dans le Taskfile. Exécutez `wails3 signing` pour le configurer, ou définissez la variable d’environnement `WAILS_WINDOWS_CERT_PASSWORD` dans l’environnement d’intégration continue.

@end

Exécutez ensuite :

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Signature multiplateforme

Les identités RSA PFX exportées permettent de signer avec Go les fichiers EXE/DLL, les exécutables d’installation NSIS et les packages MSIX/APPX individuels sous Windows, macOS et Linux. Les signatures existantes peuvent être remplacées. La signature prépare la sortie et ne la publie qu’après la réussite de la signature et de la vérification de l’horodatage, tout en préservant l’original en cas d’échec. Utilisez `--output` pour écrire un fichier signé distinct.

Wake lit ces paramètres directement depuis `wails.yaml` :

```yaml
windows:
  signing:
    certificate: signing/windows.pfx
    credential: WINDOWS_SIGNING_PASSWORD
    timestamp_server: http://timestamp.digicert.com
```

`credential` indique le nom de la variable d’environnement contenant le mot de passe, jamais le mot de passe lui-même. Activez `sign: true` pour la cible de publication ; consultez le [guide de compilation YAML](/experimental/yaml-builds/). Wake n’appelle pas les Taskfiles. Les projets existants peuvent continuer à utiliser la configuration Taskfile ci-dessus.

Le profil portable utilise RSA/SHA-256 et valide les horodatages RFC 3161, y compris la chaîne de certificats de confiance du signataire de l’horodatage. Le service par défaut est DigiCert. L’éditeur MSIX doit correspondre au sujet du certificat de signature ; les bundles de packages et les archives d’au moins 4 Gio ne sont pas pris en charge. La taille des fichiers PFX est limitée à 16 Mio. Les autres algorithmes de certificat et le chiffrement RC2 sont refusés. Pour la signature portable, réexportez un PFX RC2 avec AES ou 3DES ; sous Windows, un PFX RC2 authentifié utilise SignTool à la place. Un mot de passe incorrect ou un PFX endommagé ne déclenche jamais ce mécanisme de repli.

Les clés provenant du magasin de certificats, d’un jeton matériel ou d’un fournisseur conservent le workflow SignTool natif de Windows avec `thumbprint`. La signature MSI requiert également SignTool sous Windows. La configuration d’un certificat exporté ne rend pas ces workflows portables.

### Formats Windows pris en charge

| Format | Extension | Remarques |
| --- | --- | --- |
| Exécutables et bibliothèques | .exe, .dll | Signature PE portable, y compris pour les programmes d’installation NSIS |
| Programmes d’installation | .msi | SignTool pour Windows requis |
| Packages d’application | .msix, .appx | Signature portable de packages individuels ; aucun bundle |

## Signature des packages Linux

Les packages Linux (DEB et RPM) sont signés à l’aide de clés PGP/GPG. Contrairement à la signature de code sous Windows et macOS, la signature des packages Linux prouve que le package provient d’une source de confiance, et non que le système d’exploitation considère le code comme fiable.

Un fichier de clé privée exporté, au format ASCII armor ou binaire, utilise le signataire Go natif sous Windows, macOS et Linux. La signature d’un fichier DEB ou RPM existant préserve le contenu de son package sans extraction ni recompression et ne nécessite aucune installation de GnuPG, `dpkg-sig` ou `rpmsign`. Les packages Arch utilisent le même chargeur de clés pour écrire une signature détachée. GnuPG reste utile pour générer ou exporter des clés et vérifier les signatures de manière indépendante. Les configurations DEB/RPM existantes qui indiquent une identité de trousseau plutôt qu’un fichier de clé conservent leurs outils de signature externes.

### Prérequis

- Paire de clés PGP (peut être générée avec Wails)

### Génération d’une clé PGP

La méthode la plus simple consiste à utiliser l’assistant de configuration : exécutez `wails3 setup`, ouvrez l’onglet **Linux** de l’étape de signature, puis choisissez **Créer une nouvelle clé GPG**. L’assistant :

- génère une clé RSA 4096 dans votre trousseau GPG (laissez la phrase secrète vide pour obtenir une clé utilisable sans intervention et adaptée à la CI),
- **l’exporte vers `~/.wails/signing/<keyid>.asc`** (le fichier utilisé pour signer une compilation), et
- enregistre à la fois l’identifiant de la clé et le chemin d’exportation dans `~/.config/wails/defaults.yaml`, afin qu’ils soient utilisés automatiquement lors de la signature.

@note{type="note"}
Si une clé a été configurée **uniquement par son identifiant** (par exemple, si elle a été créée avant que l’assistant n’exporte automatiquement les clés), l’onglet Linux l’exporte dans un fichier lors de votre prochaine visite et renseigne le chemin : aucune étape manuelle n’est nécessaire. Les compilations utilisent un *fichier* de clé pour signer ; c’est donc le chemin qui importe.

@end

Vous pouvez également effectuer cette opération manuellement avec `gpg` :

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Paramètres recommandés : RSA 4096 bits, expiration après 1 an, protection par un mot de passe robuste.

@note{type="caution"}
Protégez votre clé privée ! Stockez-la sous forme chiffrée et sauvegardez-la en lieu sûr.

@end

### Configuration

Modifiez `build/linux/Taskfile.yml` et définissez les variables de signature :

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variable | Obligatoire | Description |
| --- | --- | --- |
| `PGP_KEY` | Remplacement | Chemin du fichier de clé privée PGP exporté (si cette valeur n’est pas définie, utilise la clé configurée globalement au moyen de `wails3 setup`) |
| `SIGN_ROLE` | Non | Rôle de signature DEB (valeur par défaut : builder) |

@note{type="note"}
Le mot de passe de la clé PGP est stocké dans le trousseau de votre système, et non dans le Taskfile. Exécutez `wails3 signing` pour le configurer, ou définissez la variable d’environnement `WAILS_PGP_PASSWORD` dans la CI.

@end

Exécutez ensuite :

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### Rôles de signature DEB

Pour les packages DEB, vous pouvez indiquer le rôle de signature au moyen de `SIGN_ROLE` :

- `origin` : signature provenant de l’origine du package
- `maint` : signature du responsable du package
- `archive` : signature du responsable de l’archive
- `builder` : signature du générateur du package (valeur par défaut)

### Signature multiplateforme

Les packages Linux peuvent être signés depuis n’importe quelle plateforme. La même configuration Taskfile et les mêmes commandes fonctionnent sous Windows et macOS.

### Affichage des informations sur la clé

```bash
gpg --show-keys signing-key.asc
```

Sortie :

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Vérification des packages Linux

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### Distribution de votre clé publique

Les utilisateurs ont besoin de votre clé publique pour vérifier les packages :

```bash
# Export public key for distribution
gpg --armor --export "your@email.com" > myapp-signing.pub.asc

# Users can import it:
# For DEB (apt):
sudo apt-key add myapp-signing.pub.asc
# Or for modern apt:
sudo cp myapp-signing.pub.asc /etc/apt/trusted.gpg.d/

# For RPM:
sudo rpm --import myapp-signing.pub.asc
```

## Intégration à GitHub Actions

Dans les environnements de CI, les mots de passe sont fournis au moyen de variables d’environnement plutôt que par le trousseau du système :

| Variable d’environnement | Description |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Mot de passe du certificat Windows |
| `WAILS_PGP_PASSWORD` | Mot de passe de la clé PGP pour les packages Linux |

Vous pouvez également transmettre directement les variables Taskfile :

```bash
wails3 task darwin:sign SIGN_IDENTITY="$SIGN_IDENTITY" KEYCHAIN_PROFILE="$KEYCHAIN_PROFILE"
```

### Workflow macOS

```yaml
name: Build and Sign macOS

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.MACOS_CERTIFICATE }}
          CERTIFICATE_PASSWORD: ${{ secrets.MACOS_CERTIFICATE_PASSWORD }}
        run: |
          echo $CERTIFICATE_BASE64 | base64 --decode > certificate.p12
          security create-keychain -p "" build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p "" build.keychain
          security import certificate.p12 -k build.keychain -P "$CERTIFICATE_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "" build.keychain

      - name: Store Notarization Credentials
        env:
          APPLE_ID: ${{ secrets.APPLE_ID }}
          APPLE_TEAM_ID: ${{ secrets.APPLE_TEAM_ID }}
          APPLE_APP_PASSWORD: ${{ secrets.APPLE_APP_PASSWORD }}
        run: |
          xcrun notarytool store-credentials "notarize-profile" \
            --apple-id "$APPLE_ID" \
            --team-id "$APPLE_TEAM_ID" \
            --password "$APPLE_APP_PASSWORD"

      - name: Build, Sign, and Notarize
        env:
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
        run: |
          wails3 task darwin:sign:notarize \
            SIGN_IDENTITY="$SIGN_IDENTITY" \
            KEYCHAIN_PROFILE="notarize-profile"

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-macOS
          path: bin/*.app
```

### Workflow Windows

```yaml
name: Build and Sign Windows

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Import Certificate
        env:
          CERTIFICATE_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
        run: |
          $certBytes = [Convert]::FromBase64String($env:CERTIFICATE_BASE64)
          [IO.File]::WriteAllBytes("certificate.pfx", $certBytes)

      - name: Build and Sign
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: MyApp-Windows
          path: bin/*.exe
```

### Workflow multiplateforme (exécuteur Linux)

Signez les packages Windows et Linux depuis un seul exécuteur Linux :

```yaml
name: Build and Sign (Cross-Platform)

on:
  push:
    tags: ['v*']

jobs:
  build-and-sign:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

      - name: Install Build Dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y nsis rpm

      # Import certificates
      - name: Import Certificates
        env:
          WINDOWS_CERT_BASE64: ${{ secrets.WINDOWS_CERTIFICATE }}
          PGP_KEY_BASE64: ${{ secrets.PGP_PRIVATE_KEY }}
        run: |
          echo "$WINDOWS_CERT_BASE64" | base64 -d > certificate.pfx
          echo "$PGP_KEY_BASE64" | base64 -d > signing-key.asc

      # Build and sign Windows
      - name: Build and Sign Windows
        env:
          WAILS_WINDOWS_CERT_PASSWORD: ${{ secrets.WINDOWS_CERTIFICATE_PASSWORD }}
        run: |
          wails3 task windows:sign SIGN_CERTIFICATE=certificate.pfx

      # Build and sign Linux packages
      - name: Build and Sign Linux Packages
        env:
          WAILS_PGP_PASSWORD: ${{ secrets.PGP_PASSWORD }}
        run: |
          wails3 task linux:sign:packages PGP_KEY=signing-key.asc

      # Cleanup secrets
      - name: Cleanup
        if: always()
        run: rm -f certificate.pfx signing-key.asc

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: signed-binaries
          path: |
            bin/*.exe
            bin/*.deb
            bin/*.rpm
```

@note{type="note"}
L’utilisation d’un exécuteur Linux pour la signature multiplateforme simplifie la CI/CD en supprimant le besoin d’exécuteurs Windows distincts. Les identités du trousseau Apple et les profils de notarisation nécessitent un exécuteur macOS. Les identités exportées et les identifiants d’API d’équipe prennent en charge le workflow Wake portable décrit ci-dessus.

@end

## Référence de la CLI

### wails3 signing

Assistant interactif permettant de configurer la signature pour votre projet.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

L’assistant vous guide dans les opérations suivantes :

- **macOS** : sélection d’un certificat Developer ID et configuration des identifiants de notarisation (appelle `xcrun notarytool store-credentials`).
- **Windows** : choix entre un fichier de certificat et une empreinte numérique, puis définition du mot de passe et du serveur d’horodatage.
- **Linux** : utilisation d’une clé PGP existante ou génération d’une nouvelle clé (appelle `gpg`), puis configuration du rôle de signature.

### wails3 entitlements

Assistant interactif permettant de configurer les autorisations macOS.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Préréglages :**

- **Développement** : crée `entitlements.dev.plist` avec la compilation JIT, le débogage et l’accès réseau
- **Production** : crée `entitlements.plist` avec un minimum d’autorisations
- **Les deux** : crée les deux fichiers (recommandé)
- **App Store** : crée des autorisations avec bac à sable pour le Mac App Store
- **Personnalisé** : choisissez chaque autorisation et le fichier cible

### wails3 sign

Signe les fichiers binaires et les paquets pour la plateforme actuelle ou indiquée. Il s’agit d’une commande enveloppe qui appelle la tâche de signature propre à la plateforme concernée.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Cette commande exécute la tâche `<platform>:sign` correspondante, qui utilise la configuration de signature de votre Taskfile.

### wails3 tool sign

Commande de bas niveau permettant de signer directement un fichier précis. Utilisée en interne par les Taskfiles.

```bash
wails3 tool sign [flags]
```

**Options communes :**\

| Option | Description |
| --- | --- |
| `--input` | Chemin du fichier à signer |
| `--output` | Chemin de sortie (facultatif ; par défaut, le fichier est modifié sur place) |
| `--verbose` | Active la sortie détaillée |

**Options Windows/macOS :**\

| Option | Description |
| --- | --- |
| `--certificate` | Chemin du certificat PKCS#12 (.pfx/.p12) |
| `--password` | Mot de passe du certificat |
| `--timestamp` | URL du serveur d’horodatage |

**Options propres à macOS :**\

| Option | Description |
| --- | --- |
| `--identity` | Identité de signature (utilisez « - » pour une signature ad hoc) |
| `--entitlements` | Chemin du fichier plist des autorisations |
| `--hardened-runtime` | Active l’environnement d’exécution renforcé (valeur par défaut : true) |
| `--notarize` | Soumet le fichier pour notarisation |
| `--keychain-profile` | Profil du trousseau pour la notarisation |
| `--notary-key` | Fichier `.p8` de clé privée de l’API d’équipe App Store Connect ; utilise le client de l’API Go Notary |
| `--notary-key-id` | Identifiant de la clé d’API de l’équipe ; requis avec `--notary-key` |
| `--notary-issuer-id` | Identifiant de l’émetteur de l’équipe ; requis avec `--notary-key` |

**Options propres à Windows :**\

| Option | Description |
| --- | --- |
| `--thumbprint` | Empreinte numérique du certificat dans le magasin de certificats Windows |

**Options propres à Linux :**\

| Option | Description |
| --- | --- |
| `--pgp-key` | Chemin de la clé privée PGP |
| `--pgp-password` | Mot de passe de la clé PGP |
| `--role` | Rôle de signature DEB (origin/maint/archive/builder) |

### Inspection de l’état de la signature (outils natifs)

`wails3 signing` configure les identifiants de signature. Pour inspecter l’état de la signature, utilisez directement les outils natifs :

| Tâche | Commande |
| --- | --- |
| Répertorier les identités de signature de code macOS | `security find-identity -v -p codesigning` |
| Enregistrer les identifiants de notarisation | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Inspecter un fichier de clé PGP | `gpg --show-keys <key.asc>` |
| Générer une paire de clés PGP | `gpg --full-generate-key` |
| Exporter une clé publique | `gpg --armor --export <email>` |

## Dépannage

### Problèmes sous macOS

**« Aucun certificat Developer ID trouvé »**

- Vérifiez que votre certificat est installé dans le trousseau
- Vérifiez avec `security find-identity -v -p codesigning` qu’il n’a pas expiré
- Assurez-vous de disposer d’un certificat « Developer ID Application » et pas seulement d’un certificat « Apple Development »

**« Échec de la notarisation »**

- Consultez le journal de notarisation : `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Vérifiez que l’environnement d’exécution renforcé est activé
- Vérifiez que votre application ne contient aucun fichier binaire non signé

**« Échec de Codesign »**

- Vérifiez que le trousseau est déverrouillé : `security unlock-keychain`
- Vérifiez les autorisations des fichiers dans le paquet de l’application

### Problèmes sous Windows

**« Certificat introuvable »**

- Vérifiez que le chemin du certificat est correct
- Vérifiez le mot de passe du certificat
- Assurez-vous que le certificat est valide (ni expiré ni révoqué).

**« Erreur du serveur d’horodatage »**

- Essayez un autre serveur d’horodatage :
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Problèmes sous Linux

**« Clé PGP non valide »**

- Assurez-vous que le fichier de clé est au format ASCII avec armure.
- Vérifiez avec `gpg --show-keys <key.asc>` que la clé n’a pas expiré.
- Vérifiez que le mot de passe est correct.

**« Échec de la vérification de la signature »**

- Assurez-vous que la clé publique est correctement importée.
- Vérifiez que le paquet n’a pas été modifié après sa signature.

## Ressources supplémentaires

### Documentation officielle

- [Guide Apple sur la signature de code](https://developer.apple.com/support/code-signing/)
- [Documentation Apple sur la notarisation](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Signature de code Microsoft](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Signature des paquets Debian](https://wiki.debian.org/SecureApt)
- [Signature des paquets RPM](https://rpm-software-management.github.io/rpm/manual/signatures.html)
