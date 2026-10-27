---
title: "Signature de code"
description: "Guide de signature de vos applications Wails sur toutes les plateformes"
slug: "guides/build/signing"
sourcePath: "guides/build/signing.md"
---

## Signer votre application

Ce guide explique comment signer vos applications Wails pour macOS, Windows et Linux. Wails v3 fournit des outils CLI intégrés pour la signature de code, la notarisation et la gestion des clés PGP.

- **macOS** – Signez et notarisez vos applications macOS
- **Windows** – Signez vos exécutables et paquets Windows
- **Linux** – Signez les paquets DEB et RPM avec des clés PGP

## Matrice de signature multiplateforme

Cette matrice indique les formats que vous pouvez signer depuis chaque plateforme source :

| Format cible | Depuis Windows | Depuis macOS | Depuis Linux |
| --- | :---: | :---: | :---: |
| EXE/MSI Windows | ✅ | ✅ | ✅ |
| Bundle .app macOS | ❌ | ✅ | ❌ |
| Notarisation macOS | ❌ | ✅ | ❌ |
| DEB Linux | ✅ | ✅ | ✅ |
| RPM Linux | ✅ | ✅ | ✅ |

@note{type="tip"}
Les paquets Windows et Linux peuvent être signés depuis **n’importe quelle plateforme**. La signature macOS nécessite un Mac en raison des outils imposés par Apple.

@end

### Moteurs de signature

Wails sélectionne automatiquement le meilleur moteur de signature disponible :

| Plateforme | Moteur natif | Moteur multiplateforme |
| --- | --- | --- |
| Windows | `signtool.exe` (SDK Windows) | Intégré |
| macOS | `codesign` (Xcode) | Non disponible |
| Linux | S/O | Intégré |

Sur la plateforme native, Wails utilise les outils natifs pour une compatibilité maximale. Lors d’une compilation croisée, il utilise la prise en charge intégrée de la signature.

## Démarrage rapide

Le moyen le plus rapide de configurer la signature consiste à utiliser l’assistant de configuration :

```bash
wails3 setup
```

Cette opération écrit dans `~/.config/wails/defaults.yaml` une configuration de signature partagée qui est **prise en compte lors de la signature sur toutes les plateformes** (voir [Ordre de priorité de la configuration](#ordre-de-priorit-de-la-configuration)). L’étape de signature de l’assistant :

- Détecte les outils de signature installés pour la plateforme que vous souhaitez cibler — **depuis n’importe quel système hôte** — et affiche la commande d’installation correspondant à votre système d’exploitation (par exemple, `brew install gnupg`, `sudo apt install osslsigncode`, `winget install GnuPG.Gpg4win`).
- Sous macOS, répertorie les certificats Developer ID de votre trousseau.
- Pour Linux, répertorie vos clés GPG et peut en **générer et exporter** une nouvelle.
- Pour Windows, peut **générer un certificat autosigné** (à des fins de test) avec OpenSSL.
- **Stocke les mots de passe de manière sécurisée dans le gestionnaire d’identifiants de votre système** (et non dans les Taskfiles).

@note{type="tip"}
Les mots de passe sont stockés dans le gestionnaire d’identifiants natif de votre système (Trousseau d’accès macOS, Gestionnaire d’informations d’identification Windows ou Secret Service sous Linux). Votre configuration de signature reste ainsi sécurisée et fonctionne dans tous vos projets Wails.

@end

## Ordre de priorité de la configuration

Lorsque vous exécutez une tâche de signature, chaque option de signature est déterminée selon l’ordre suivant (la première correspondance prévaut) :

1. Une option explicite transmise à `wails3 tool sign` (par exemple, `--pgp-key`, `--certificate`, `--identity`).
2. La **variable correspondante du Taskfile du projet** (`PGP_KEY`, `SIGN_CERTIFICATE`/`SIGN_THUMBPRINT`, `SIGN_IDENTITY`, etc.).
3. La **configuration globale** dans `~/.config/wails/defaults.yaml` (écrite par `wails3 setup`).

Les variables du Taskfile constituent donc des **valeurs de remplacement facultatives** : si une variable n’est pas définie, la clé, le certificat ou l’identité configuré globalement est utilisé. Si aucune de ces trois sources ne fournit de valeur, la commande de signature affiche une erreur claire indiquant comment la configurer.

## Configuration par projet

Pour configurer la signature pour un seul projet plutôt que globalement, exécutez l’assistant propre au projet depuis l’intérieur de celui-ci : il écrit les `vars` dans les fichiers `build/<platform>/Taskfile.yml` de ce projet :

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

Le mot de passe est récupéré depuis le gestionnaire d’identifiants du système (exécutez `wails3 signing` pour le configurer).

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

Le mot de passe est récupéré depuis le gestionnaire d’identifiants du système (exécutez `wails3 signing` pour le configurer).

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

Pour configurer interactivement la signature de toutes les plateformes, utilisez l’assistant :

```bash
wails3 signing
```

## Signature de code sous macOS

### Prérequis

- Compte Apple Developer (99 $/an)
- Certificat Developer ID Application
- Outils de ligne de commande Xcode installés

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
Pour distribuer une application en dehors de l’App Store, vous avez besoin d’un certificat **Developer ID Application**.

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
| `SIGN_IDENTITY` | Oui | Votre Developer ID (par exemple, « Developer ID Application: Your Company (TEAMID) ») |
| `KEYCHAIN_PROFILE` | Pour la notarisation | Nom du profil du trousseau contenant les identifiants stockés |
| `ENTITLEMENTS` | Non | Chemin d’accès au fichier d’autorisations |

Exécutez ensuite :

```bash
wails3 task darwin:sign           # Build, package, and sign
wails3 task darwin:sign:notarize  # Build, package, sign, and notarize
```

### Droits

Les droits déterminent les fonctionnalités auxquelles votre application peut accéder. Les applications Wails nécessitent généralement des droits différents pour le développement et la production :

- **Développement** : nécessite des droits pour la compilation JIT, la mémoire non signée et le débogage
- **Production** : droits minimaux (accès réseau uniquement)

Utilisez l’assistant de configuration interactif pour générer les deux fichiers :

```bash
wails3 entitlements
```

Il crée :

- `build/darwin/entitlements.dev.plist` — Pour les builds de développement
- `build/darwin/entitlements.plist` — Pour les builds de production/signés

**Préréglages disponibles :**\

| Préréglage | Description |
| --- | --- |
| Développement | Compilation JIT, mémoire non signée, débogage, réseau |
| Production | Réseau uniquement (configuration minimale et la plus sûre) |
| Les deux | Crée les fichiers de développement et de production (recommandé) |
| App Store | Sandbox activée avec accès au réseau et aux fichiers |
| Personnalisé | Choisissez les droits individuellement |

@note{type="note"}
La tâche `run` du Taskfile Darwin utilise automatiquement `entitlements.dev.plist`. Pour les builds de production, les tâches `sign` utilisent `entitlements.plist`.

@end

Définissez ensuite `ENTITLEMENTS` dans les variables de votre Taskfile de façon à référencer le fichier approprié.

### Notarisation

Apple exige que toutes les applications distribuées soient notarisées.

@steps
### **Stockez vos identifiants dans le trousseau** (configuration unique). Exécutez `wails3 signing` (la commande vous demande ces valeurs et appelle `notarytool` en arrière-plan) ou appelez directement `notarytool` :
```bash
xcrun notarytool store-credentials "my-notarize-profile" \
  --apple-id "your@email.com" \
  --team-id "ABCD1234" \
  --password "app-specific-password"
```

### **Définissez KEYCHAIN_PROFILE dans votre Taskfile** de façon à reprendre le nom du profil ci-dessus.
### **Signez et notarisez votre application** :
```bash
wails3 task darwin:sign:notarize
```

### **Vérifiez la notarisation** :
```bash
spctl --assess --verbose=2 bin/MyApp.app
```

@end

@note{type="note"}
La notarisation prend généralement 1-2 minutes. Le ticket est automatiquement agrafé à votre application.

@end

## Signature du code sous Windows

### Prérequis

- Certificat de signature de code (fourni par DigiCert, Sectigo, etc.)
- Pour la signature native sous Windows : SDK Windows installé (pour `signtool.exe`)
- Pour la signature multiplateforme depuis macOS/Linux : [`osslsigncode`](https://github.com/mtrojnar/osslsigncode) (l’étape de signature `wails3 setup` indique la commande d’installation propre au système hôte)

### Génération d’un certificat autosigné (tests)

Si vous devez simplement tester le pipeline de signature, exécutez `wails3 setup`, ouvrez l’onglet **Windows** de l’étape de signature, puis choisissez **Générer un certificat autosigné**. L’assistant utilise OpenSSL pour créer un fichier `.pfx` de signature de code et enregistre son chemin dans la configuration globale.

@note{type="caution"}
Les certificats autosignés sont réservés **aux tests et à la distribution interne** : ils déclenchent des avertissements SmartScreen pour les utilisateurs finaux. Les versions publiques nécessitent un certificat émis par une autorité de certification de confiance.

@end

### Configuration

Modifiez `build/windows/Taskfile.yml` et définissez les variables de signature :

```yaml
vars:
  SIGN_CERTIFICATE: "path/to/certificate.pfx"
  # Or use thumbprint instead:
  # SIGN_THUMBPRINT: "certificate-thumbprint"
  # TIMESTAMP_SERVER: "http://timestamp.digicert.com"
```

| Variable | Obligatoire | Description |
| --- | --- | --- |
| `SIGN_CERTIFICATE` | Remplacement | Chemin du fichier de certificat .pfx/.p12 (utilise la configuration globale si cette variable n’est pas définie) |
| `SIGN_THUMBPRINT` | Remplacement | Empreinte du certificat dans le magasin de certificats Windows (solution de remplacement à `SIGN_CERTIFICATE`) |
| `TIMESTAMP_SERVER` | Non | URL du serveur d’horodatage (valeur par défaut : http://timestamp.digicert.com) |

@note{type="note"}
Ces variables sont des **valeurs de remplacement** facultatives. Si elles ne sont pas définies, le certificat configuré globalement avec `wails3 setup` est utilisé. Consultez [Priorité des configurations](#ordre-de-priorit-de-la-configuration).

@end

@note{type="note"}
Le mot de passe du certificat est stocké dans le trousseau de votre système, et non dans le Taskfile. Exécutez `wails3 signing` pour le configurer ou définissez la variable d’environnement `WAILS_WINDOWS_CERT_PASSWORD` dans l’environnement CI.

@end

Exécutez ensuite :

```bash
wails3 task windows:sign           # Build and sign executable
wails3 task windows:sign:installer # Build and sign NSIS installer
```

### Signature multiplateforme

Vous pouvez signer des exécutables Windows depuis n’importe quelle plateforme. La même configuration de Taskfile et les mêmes commandes fonctionnent sous macOS et Linux.

### Formats Windows pris en charge

| Format | Extension | Remarques |
| --- | --- | --- |
| Exécutables | .exe | Signature PE standard |
| Programmes d’installation | .msi | Paquets Windows Installer |
| Paquets d’applications | .msix, .appx | Applications Windows modernes |

## Signature des paquets Linux

Les paquets Linux (DEB et RPM) sont signés à l’aide de clés PGP/GPG. Contrairement à la signature de code sous Windows et macOS, la signature des paquets Linux prouve que le paquet provient d’une source de confiance, et non que le système d’exploitation considère le code comme fiable.

### Prérequis

- Paire de clés PGP (peut être générée avec Wails)

### Génération d’une clé PGP

La méthode la plus simple consiste à utiliser l’assistant de configuration : exécutez `wails3 setup`, ouvrez l’onglet **Linux** de l’étape de signature, puis choisissez **Créer une nouvelle clé GPG**. L’assistant :

- génère une clé RSA 4096 dans votre trousseau GPG (laissez la phrase secrète vide pour obtenir une clé utilisable sans intervention et adaptée à l’environnement CI),
- **l’exporte vers `~/.wails/signing/<keyid>.asc`** (le fichier utilisé par un build pour signer), et
- enregistre l’identifiant de la clé et le chemin du fichier exporté dans `~/.config/wails/defaults.yaml` afin que la clé soit automatiquement utilisée lors de la signature.

@note{type="note"}
Si une clé a été configurée **uniquement avec son identifiant** (par exemple, si elle a été créée avant que l’assistant n’exporte automatiquement les clés), l’onglet Linux l’exportera dans un fichier lors de votre prochaine visite et renseignera le chemin : aucune étape manuelle n’est nécessaire. Les builds signent avec un *fichier* de clé ; c’est donc le chemin qui importe.

@end

Vous pouvez également procéder manuellement avec `gpg` :

```bash
# Interactive — the wizard will prompt for name, email, key size and expiry.
gpg --full-generate-key

# Export the key pair to ASCII-armoured files for the Taskfile to consume.
gpg --armor --export-secret-keys "your@email.com" > signing-key.asc
gpg --armor --export "your@email.com" > signing-key.pub.asc
```

Paramètres recommandés : RSA 4096 bits, expiration après 1 an(s) et protection par un mot de passe robuste.

@note{type="caution"}
Protégez votre clé privée ! Stockez-la sous forme chiffrée et sauvegardez-la en lieu sûr.

@end

### Configuration

Modifiez `build/linux/Taskfile.yml` et définissez les variables de signature :

```yaml
vars:
  PGP_KEY: "path/to/signing-key.asc"
  # SIGN_ROLE: "builder"  # Options: origin, maint, archive, builder
```

| Variable | Obligatoire | Description |
| --- | --- | --- |
| `PGP_KEY` | Remplacement | Chemin du fichier de clé privée PGP exporté (si cette valeur n’est pas définie, utilise la clé configurée globalement via `wails3 setup`) |
| `SIGN_ROLE` | Non | Rôle de signature DEB (valeur par défaut : builder) |

@note{type="note"}
Le mot de passe de la clé PGP est stocké dans le trousseau de clés de votre système, et non dans le Taskfile. Exécutez `wails3 signing` pour le configurer ou définissez la variable d’environnement `WAILS_PGP_PASSWORD` dans l’environnement d’intégration continue.

@end

Exécutez ensuite :

```bash
wails3 task linux:sign:deb       # Build and sign DEB package
wails3 task linux:sign:rpm       # Build and sign RPM package
wails3 task linux:sign:packages  # Build and sign all packages
```

### Rôles de signature DEB

Pour les paquets DEB, vous pouvez spécifier le rôle de signature via `SIGN_ROLE` :

- `origin` : signature de l’origine du paquet
- `maint` : signature du responsable du paquet
- `archive` : signature du responsable de l’archive
- `builder` : signature du générateur du paquet (valeur par défaut)

### Signature multiplateforme

Vous pouvez signer les paquets Linux depuis n’importe quelle plateforme. La même configuration du Taskfile et les mêmes commandes fonctionnent sous Windows et macOS.

### Affichage des informations sur la clé

```bash
gpg --show-keys signing-key.asc
```

Sortie :

```
pub   rsa4096 2024-01-15 [SC] [expires: 2025-01-15]
      1234 5678 90AB CDEF 1234 5678 90AB CDEF 1234 5678
uid                      Your Name <your@email.com>
```

### Vérification des paquets Linux

```bash
# Verify DEB signature
dpkg-sig --verify myapp_1.0.0_amd64.deb

# Verify RPM signature
rpm --checksig myapp-1.0.0.x86_64.rpm
```

### Distribution de votre clé publique

Les utilisateurs ont besoin de votre clé publique pour vérifier les paquets :

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

Dans les environnements d’intégration continue, les mots de passe sont fournis au moyen de variables d’environnement plutôt que par le trousseau de clés du système :

| Variable d’environnement | Description |
| --- | --- |
| `WAILS_WINDOWS_CERT_PASSWORD` | Mot de passe du certificat Windows |
| `WAILS_PGP_PASSWORD` | Mot de passe de la clé PGP pour les paquets Linux |

Vous pouvez également transmettre directement les variables du Taskfile :

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

Signez les paquets Windows et Linux depuis un même exécuteur Linux :

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
L’utilisation d’un exécuteur Linux pour la signature multiplateforme simplifie l’intégration et le déploiement continus en évitant de recourir à des exécuteurs Windows distincts. La signature sous macOS nécessite néanmoins un exécuteur macOS en raison des outils natifs imposés par Apple.

@end

## Référence de l’interface en ligne de commande

### wails3 signing

Assistant interactif permettant de configurer la signature de votre projet.

```bash
wails3 signing [flags]

Flags:
  --platform    Platform to configure (darwin, windows, linux). Repeatable.
                If omitted, auto-detects which platforms to configure from the build directory.
```

L’assistant vous guide pour effectuer les opérations suivantes :

- **macOS** : sélection d’un certificat Developer ID et configuration des identifiants de notarisation (appelle `xcrun notarytool store-credentials`).
- **Windows** : choix entre un fichier de certificat et une empreinte numérique, puis définition du mot de passe et du serveur d’horodatage.
- **Linux** : utilisation d’une clé PGP existante ou génération d’une nouvelle clé (appelle `gpg`), puis configuration du rôle de signature.

### wails3 entitlements

Assistant interactif permettant de configurer les droits macOS.

```bash
wails3 entitlements [flags]

Flags:
  --output    Output path for entitlements.plist (default: build/darwin/entitlements.plist)
```

**Préréglages :**

- **Développement** : crée `entitlements.dev.plist` avec la compilation JIT, le débogage et l’accès réseau
- **Production** : crée `entitlements.plist` avec un minimum de droits
- **Les deux** : crée les deux fichiers (recommandé)
- **App Store** : crée des droits avec bac à sable pour le Mac App Store
- **Personnalisé** : choisissez individuellement les droits et le fichier cible

### wails3 sign

Signe les fichiers binaires et les paquets pour la plateforme actuelle ou spécifiée. Cette commande enveloppe et appelle la tâche de signature propre à la plateforme appropriée.

```bash
wails3 sign
wails3 sign GOOS=darwin
wails3 sign GOOS=windows
wails3 sign GOOS=linux
```

Cette commande exécute la tâche `<platform>:sign` correspondante, qui utilise la configuration de signature de votre Taskfile.

### wails3 tool sign

Commande de bas niveau permettant de signer directement un fichier précis. Les Taskfiles l’utilisent en interne.

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
| `--entitlements` | Chemin du fichier plist des droits |
| `--hardened-runtime` | Active l’environnement d’exécution renforcé (valeur par défaut : true) |
| `--notarize` | Soumet le fichier à la notarisation |
| `--keychain-profile` | Profil du trousseau de clés pour la notarisation |

**Options propres à Windows :**\

| Option | Description |
| --- | --- |
| `--thumbprint` | Empreinte du certificat dans le magasin Windows |

**Options propres à Linux :**\

| Option | Description |
| --- | --- |
| `--pgp-key` | Chemin vers la clé privée PGP |
| `--pgp-password` | Mot de passe de la clé PGP |
| `--role` | Rôle de signature DEB (origin/maint/archive/builder) |

### Inspection de l’état de la signature (outils natifs)

`wails3 signing` configure les identifiants de signature. Pour inspecter l’état de la signature, utilisez directement les outils natifs :

| Tâche | Commande |
| --- | --- |
| Répertorier les identités de signature de code macOS | `security find-identity -v -p codesigning` |
| Enregistrer les identifiants de notarisation | `xcrun notarytool store-credentials "<profile>" --apple-id … --team-id … --password …` |
| Inspecter un fichier de clé PGP | `gpg --show-keys <key.asc>` |
| Générer une paire de clés PGP | `gpg --full-generate-key` |
| Exporter une clé publique | `gpg --armor --export <email>` |

## Dépannage

### Problèmes sous macOS

**« Aucun certificat Developer ID trouvé »**

- Vérifiez que votre certificat est installé dans le trousseau
- Vérifiez avec `security find-identity -v -p codesigning` qu’il n’a pas expiré
- Vérifiez que vous disposez d’un certificat « Developer ID Application » (et pas seulement « Apple Development »)

**« Échec de la notarisation »**

- Consultez le journal de notarisation : `xcrun notarytool log <submission-id> --keychain-profile <profile>`
- Vérifiez que l’environnement d’exécution renforcé est activé
- Vérifiez que votre application ne contient aucun binaire non signé

**« Échec de codesign »**

- Vérifiez que le trousseau est déverrouillé : `security unlock-keychain`
- Vérifiez les autorisations des fichiers du paquet de l’application

### Problèmes sous Windows

**« Certificat introuvable »**

- Vérifiez que le chemin du certificat est correct
- Vérifiez le mot de passe du certificat
- Vérifiez que le certificat est valide (ni expiré ni révoqué)

**« Erreur du serveur d’horodatage »**

- Essayez un autre serveur d’horodatage :
  - `http://timestamp.digicert.com`
  - `http://timestamp.sectigo.com`
  - `http://timestamp.comodoca.com`


### Problèmes sous Linux

**« Clé PGP non valide »**

- Vérifiez que le fichier de clé est au format ASCII armored
- Vérifiez avec `gpg --show-keys <key.asc>` que la clé n’a pas expiré
- Vérifiez que le mot de passe est correct

**« Échec de la vérification de la signature »**

- Vérifiez que la clé publique a été correctement importée
- Vérifiez que le paquet n’a pas été modifié après sa signature

## Ressources supplémentaires

### Documentation officielle

- [Guide Apple sur la signature de code](https://developer.apple.com/support/code-signing/)
- [Documentation Apple sur la notarisation](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [Signature de code Microsoft](https://docs.microsoft.com/en-us/windows-hardware/drivers/dashboard/get-a-code-signing-certificate)
- [Signature des paquets Debian](https://wiki.debian.org/SecureApt)
- [Signature des paquets RPM](https://rpm-software-management.github.io/rpm/manual/signatures.html)
