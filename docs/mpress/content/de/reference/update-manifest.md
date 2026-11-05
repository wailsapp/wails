---
title: "Update-Manifest-Protokoll"
description: "Das offene JSON-Protokoll, mit dem Wails-Apps selbstständige Updates erkennen und verifizieren. Es kann über jeden statischen Dateihost oder dynamischen Update-Server bereitgestellt werden."
slug: "reference/update-manifest"
sourcePath: "reference/update-manifest.md"
---

Das Wails Update-Manifest-Protokoll ist ein kleiner, offener JSON-Vertrag zwischen einer Wails-Anwendung und einer Update-Quelle. Jeder Dienst, der eine JSON-Datei über HTTPS bereitstellen kann, kann Wails-Updates bereitstellen: ein S3-Bucket, GitHub Pages, ein CDN oder ein dynamischer Update-Server, der Releases abhängig von einer Lizenz freigibt.

Die Clientseite ist als `endpoint`-Provider (`github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint`) Bestandteil des Frameworks. Diese Seite ist die Referenz für das Übertragungsformat und richtet sich an alle, die die Serverseite implementieren.

## Entwurfsziele

1. **Für statisches Hosting geeignet.** Eine einzige Manifestdatei pro Kanal, die die Artefakte aller Plattformen auflistet, ist eine vollständige Implementierung. Kein Servercode erforderlich.
2. **Für dynamische Server geeignet.** Der Client sendet bei jeder Prüfung `platform`, `arch`, `version` und `channel`. Dadurch kann ein Server mit genau einem Artefakt antworten, Lizenzierungsregeln anwenden oder `204 No Content` zurückgeben, wenn der Aufrufer auf dem neuesten Stand ist.
3. **Verifizierung zuerst.** Das Manifest enthält für jedes Artefakt Prüfsummen und Signaturen. Der Wails-Updater verifiziert sie anhand eines öffentlichen Schlüssels, der zur Build-Zeit im Anwendungs-Binary verankert wurde. Die Update-Quelle legt niemals ihre eigene Vertrauenswurzel fest.

## Die Anfrage

Der Client sendet einen `GET` an die konfigurierte Manifest-URL, mit `Accept: application/json` sowie allen von der Anwendung konfigurierten Headern (zum Beispiel `Authorization: License <key>`).

Die URL kann Platzhalter enthalten, die der Client bei jeder Prüfung ersetzt:

| Platzhalter | Ersetzt durch |
| --- | --- |
| `{{platform}}` | Das aktuell verwendete Betriebssystem als Go-Wert `GOOS` (`darwin`, `windows`, `linux`) |
| `{{arch}}` | Die aktuell verwendete Architektur als Go-Wert `GOARCH` (`amd64`, `arm64`, ...) |
| `{{version}}` | Die derzeit installierte Version |
| `{{channel}}` | Der konfigurierte Release-Kanal, sofern festgelegt |

Jeder der vier Werte, der nicht durch einen Platzhalter verwendet wird, wird als gleichnamiger Abfrageparameter angehängt (`channel` nur, wenn konfiguriert). Daher sind die beiden folgenden Konfigurationen gültig und gleichwertig:

```text
# Dynamic server: reads query parameters
https://updates.example.com/check
  -> GET /check?platform=darwin&arch=arm64&version=1.0.0&channel=stable

# Static host: one manifest per platform/arch/channel path
https://cdn.example.com/updates/{{platform}}/{{arch}}/{{channel}}.json
  -> GET /updates/darwin/arm64/stable.json?version=1.0.0
```

Statische Hosts ignorieren einfach die empfangenen Abfrageparameter.

## Die Antwort

| Status | Bedeutung |
| --- | --- |
| `200 OK` | Es folgt ein Manifest. Der Client entscheidet, ob es sich um ein Upgrade handelt. |
| `204 No Content` | Der Server hat die Versionen verglichen und der Aufrufer ist auf dem neuesten Stand. |
| `404 Not Found` | Nichts veröffentlicht (wird wie „auf dem neuesten Stand“ behandelt). |
| Alle anderen Werte | Ein Fehler. Der Updater wechselt zum nächsten konfigurierten Provider. |

Ein `200`-Body ist ein Manifestdokument:

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

### Felder der obersten Ebene

| Feld | Typ | Erforderlich | Hinweise |
| --- | --- | --- | --- |
| `schemaVersion` | int | nein | Protokollversion. Wenn sie fehlt, wird `1` angenommen. Clients lehnen neuere Werte ab, die sie nicht verstehen. |
| `version` | string | **ja** | SemVer 2.0.0, mit oder ohne vorangestelltes `v`. |
| `channel` | string | nein | Informativ. Ein für einen anderen Kanal konfigurierter Client behandelt das Manifest so, als wäre kein Update verfügbar. |
| `name` | string | nein | Für Menschen lesbarer Release-Titel, der im Update-Fenster angezeigt wird. |
| `notes` | string | nein | Versionshinweise in Markdown, die im Update-Fenster gerendert werden. |
| `publishedAt` | string | nein | Zeitstempel gemäß RFC 3339. |
| `artifacts` | array | **ja** | Ein Eintrag pro herunterladbarem Artefakt. Die Reihenfolge gibt die Präferenz des Herausgebers an. |
| `metadata` | object | nein | Frei definierbare Schlüssel/Wert-Daten, die an die Anwendung durchgereicht werden. |

Clients ignorieren unbekannte Felder, sodass Server eigene hinzufügen können, ohne bestehende Clients zu beeinträchtigen. Serverspezifische Ergänzungen gehören in `metadata`.

### Artefaktfelder

| Feld | Typ | Erforderlich | Hinweise |
| --- | --- | --- | --- |
| `url` | string | **ja** | Absolut oder relativ zur Manifest-URL. Nur `http(s)`. |
| `platform` | string | nein | Go-Wert für `GOOS`. Gängige Aliase (`macos`, `win`, ...) werden akzeptiert. Ein leerer Wert entspricht jeder Plattform. |
| `arch` | string | nein | Go-Wert für `GOARCH`. Gängige Aliase (`x86_64`, `aarch64`, ...) werden akzeptiert. Ein leerer Wert entspricht jeder Architektur. |
| `filename` | string | nein | Standardmäßig das letzte Pfadsegment von `url`. |
| `filetype` | string | nein | Standardmäßig die Dateinamenerweiterung. |
| `size` | int | nein | Bytes; wird für die Anzeige des Downloadfortschritts verwendet. |
| `digestAlgo` / `digest` | string / base64 | nein | `sha256` oder `sha512`. |
| `signatureAlgo` / `signature` | string / base64 | nein | `ed25519`, `ed25519ph` oder `ecdsa-p256`. `signatureAlgo` ist immer erforderlich, wenn `signature` vorhanden ist. Welche Daten die einzelnen Algorithmen signieren, erfahren Sie im [Updater-Leitfaden](/guides/updater/#cryptographic-verification). |

Der Client wählt das **erste** Artefakt aus, dessen `platform` und `arch` zum laufenden System passen. Base64-Werte werden mit und ohne Padding akzeptiert.

### Versionsvergleich

Ob das Manifest ein Upgrade darstellt, wird immer clientseitig nach den Vorrangregeln von SemVer 2.0.0 entschieden: Die `version` des Manifests muss strikt neuer als die installierte Version sein. Dadurch funktioniert statisches Hosting ohne Weiteres korrekt: Das Manifest beschreibt stets die neueste Version, und aktuelle Clients führen einfach keine Aktion aus. Gleichzeitig steht `204` dynamischen Servern weiterhin zur Bandbreiteneinsparung zur Verfügung.

## Verifizierung und Vertrauen

Prüfsummen und Signaturen werden im Manifest übertragen, die Vertrauenswurzel jedoch nicht: Signaturen werden anhand des öffentlichen Schlüssels verifiziert, den die Anwendung zur Build-Zeit über `updater.Config.PublicKey` fest hinterlegt hat. Eine kompromittierte oder ausgetauschte Updatequelle kann keinen eigenen Schlüssel bereitstellen. Ein Artefakt mit Signatur wird sicher abgelehnt, wenn in der Anwendung kein Schlüssel fest hinterlegt ist. Dasselbe gilt für eine Signatur ohne deklariertes `signatureAlgo` oder für eine Signatur, die sich nicht dekodieren lässt: Clients greifen niemals stillschweigend auf eine reine Digest-Verifizierung zurück.

Artefakte, die nur einen Digest enthalten, werden nach erfolgreicher Digest-Prüfung installiert. Dies schützt vor Beschädigungen, setzt für den Schutz vor Manipulationen jedoch TLS und einen nicht kompromittierten Host voraus. Verwenden Sie für alle sicherheitskritischen Inhalte Signaturen.

Ein Artefakt mit dem `ed25519ph`-Verfahren des Frameworks zu signieren, erfordert nur wenige Zeilen Go-Code:

```go
digest := sha512.Sum512(artifactBytes)
sig, _ := privateKey.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
manifest.Artifacts[i].DigestAlgo = "sha512"
manifest.Artifacts[i].Digest = base64.StdEncoding.EncodeToString(digest[:])
manifest.Artifacts[i].SignatureAlgo = "ed25519ph"
manifest.Artifacts[i].Signature = base64.StdEncoding.EncodeToString(sig)
```

In der Praxis müssen Sie diesen Code nur selten schreiben: Die CLI übernimmt dies für Sie.

## Veröffentlichen mit der wails3-CLI

Die Befehlsgruppe `wails3 updater` deckt die gesamte Veröffentlichungspipeline ab. Für eine Veröffentlichung sind drei Befehle erforderlich:

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

`manifest` akzeptiert Dateien oder Verzeichnisse. Schlüsselmaterial, `.json` sowie Begleitdateien mit Prüfsummen und Hinweisen werden automatisch übersprungen. Jedes Artefakt wird für SHA-512 gestreamt; wenn `-key` angegeben ist, wird der Digest mit Ed25519ph signiert. `platform` und `arch` werden aus konventionellen Dateinamen wie `MyApp-2.1.0-darwin-arm64.zip` abgeleitet. Gängige Aliase wie `macOS`, `win64`, `x86_64` und `aarch64` werden erkannt. Für nicht ableitbare Werte wird eine Warnung ausgegeben; das betreffende Artefakt entspricht dann jeder Plattform. Lassen Sie `-url-prefix` weg, um relative URLs auszugeben, und laden Sie das Manifest neben den Artefakten hoch.

`verify` wird bei jeder Abweichung mit einem Exitcode ungleich null beendet und eignet sich daher als CI-Prüfschritt zwischen Build und Veröffentlichung. Für Server, die Manifeste selbst zusammenstellen, gibt `wails3 updater sign -key updater.key <files...>` die Felder `digest`/`signature` jeder Datei als JSON aus, das direkt in Ihr eigenes Dokument übernommen werden kann.

## Authentifizierung

Die Authentifizierung ist Aufgabe des Servers; das Protokoll überträgt lediglich Header. Der Client sendet seine konfigurierten Header bei jeder Manifestanfrage erneut. Bei Artefaktdownloads wird der Header `Authorization` nur gesendet, wenn sich die Artefakt-URL auf demselben Host wie das Manifest befindet und kein Downgrade von `https` auf `http` erfolgt. Bei jeder ursprungsübergreifenden Weiterleitung oder Weiterleitung mit Downgrade wird er entfernt. Dadurch gelangen Anmeldedaten weder an ein CDN oder einen Objektspeicher noch werden sie im Klartext übertragen.

Beispiel mit Lizenzprüfung, das sich gut mit gehosteten Lizenzierungsdiensten kombinieren lässt:

```go
ep, _ := endpoint.New(endpoint.Config{
    URL:     "https://updates.example.com/check",
    Headers: map[string]string{"Authorization": "License " + licenseKey},
})
```

## Clientkonfiguration

Die vollständige Referenz zu `endpoint.Config` und Informationen dazu, wie sich der Provider neben den Providern für GitHub, keygen.sh und AppCast in Fallback-Ketten einfügt, finden Sie im [Updater-Leitfaden](/guides/updater/#providers).
