---
title: "Windows-UAC-Konfiguration"
description: "Konfigurieren Sie die Benutzerkontensteuerung (UAC) für Ihre Wails-Anwendung unter Windows"
slug: "guides/windows-uac"
sourcePath: "guides/windows-uac.md"
---

Relevante Plattformen: <span class="mpress-badge mpress-badge-note">Windows</span>

<br/>

Die Windows-Benutzerkontensteuerung (UAC) bestimmt die Ausführungsberechtigungen Ihrer Wails-Anwendung. Wails-v3-Anwendungen enthalten standardmäßig eine explizite UAC-Konfiguration in ihrem Windows-Manifest. Dies gewährleistet ein konsistentes Verhalten auf verschiedenen Computern.

## UAC-Ausführungsebenen

Windows-Anwendungen können über ihre Manifestdatei unterschiedliche Ausführungsebenen anfordern. Wails v3 fügt automatisch eine UAC-Konfiguration mit einer standardmäßigen Ausführungsebene ein, die Sie an die Anforderungen Ihrer Anwendung anpassen können.

### Verfügbare Ausführungsebenen

| Ebene | Beschreibung | Anwendungsfall |
| --- | --- | --- |
| `asInvoker` | Wird mit denselben Berechtigungen wie der übergeordnete Prozess ausgeführt | Standard für die meisten Anwendungen |
| `highestAvailable` | Wird mit den höchsten für den Benutzer verfügbaren Berechtigungen ausgeführt | Anwendungen, die möglicherweise erhöhte Zugriffsrechte benötigen |
| `requireAdministrator` | Erfordert immer Administratorrechte | Systemdienstprogramme, Installationsprogramme |

### Standardkonfiguration

Wails-v3-Anwendungen enthalten in ihrem Windows-Manifest eine standardmäßige UAC-Konfiguration:

```xml
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

Diese Konfiguration stellt sicher, dass Ihre Anwendung:

- Mit denselben Berechtigungen wie der startende Prozess ausgeführt wird
- Standardmäßig keine Rechteerhöhung erfordert
- Auf verschiedenen Computern konsistent funktioniert
- Bei normalen Benutzern keine UAC-Eingabeaufforderungen auslöst

## UAC-Konfiguration anpassen

Da Wails v3 die Anpassung der Build-Ressourcen vorsieht, können Sie die UAC-Konfiguration ändern, indem Sie die Vorlage des Windows-Manifests direkt bearbeiten.

### Manifestvorlage finden

Die Vorlage des Windows-Manifests befindet sich unter:

```
build/windows/wails.exe.manifest
```

### Ausführungsebene ändern

Um die Ausführungsebene zu ändern, bearbeiten Sie das Attribut `level` im Element `requestedExecutionLevel`:

```xml {title="build/windows/wails.exe.manifest"}
<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
        <requestedPrivileges>
            <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
        </requestedPrivileges>
    </security>
</trustInfo>
```

### Beispiele

#### Standardanwendung (Standard)

Die meisten Anwendungen sollten die standardmäßige Ebene `asInvoker` verwenden:

```xml
<requestedExecutionLevel level="asInvoker" uiAccess="false"/>
```

#### Systemdienstprogramm

Anwendungen, die erhöhte Zugriffsrechte benötigen, sofern diese verfügbar sind:

```xml
<requestedExecutionLevel level="highestAvailable" uiAccess="false"/>
```

#### Administrationswerkzeug

Anwendungen, die immer Administratorrechte erfordern:

```xml
<requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
```

## UI-Zugriff

Das Attribut `uiAccess` steuert, ob Ihre Anwendung mit UI-Elementen interagieren kann, die höhere Berechtigungen besitzen. In den meisten Fällen sollte es auf `false` gesetzt bleiben.

Setzen Sie den Wert nur dann auf `true`, wenn Ihre Anwendung Folgendes benötigt:

- Eingaben an andere Anwendungen senden
- Die Benutzeroberfläche anderer Anwendungen steuern
- Auf UI-Elemente von Prozessen mit höheren Berechtigungen zugreifen

@note{type="caution" title="Anforderungen für den UI-Zugriff"}
Wenn Sie den Wert auf `uiAccess="true"` setzen, muss Ihre Anwendung folgende Anforderungen erfüllen:

- Sie muss mit einem Zertifikat einer vertrauenswürdigen Zertifizierungsstelle digital signiert sein
- Sie muss an einem sicheren Speicherort installiert sein (Program Files oder Windows\System32)

@end

## Build mit benutzerdefinierten UAC-Einstellungen erstellen

Nachdem Sie Ihre Manifestvorlage geändert haben, erstellen Sie den Build Ihrer Anwendung wie gewohnt:

```bash
wails3 build
```

Der Buildprozess bettet Ihre benutzerdefinierte UAC-Konfiguration automatisch in die ausführbare Datei ein.

## UAC-Konfiguration überprüfen

Mit dem Werkzeug `go-winres` können Sie überprüfen, ob Ihre UAC-Einstellungen ordnungsgemäß eingebettet sind:

```bash
go-winres extract --in your-app.exe --out extracted-resources/
```

Prüfen Sie anschließend die extrahierte Manifestdatei, um sicherzustellen, dass Ihre UAC-Konfiguration enthalten ist.

@note{type="tip" title="Persistenz des Manifests"}
Anders als bei einigen anderen Frameworks wird die UAC-Konfiguration von Wails v3 beim Kompilieren direkt in die ausführbare Datei eingebettet. Dadurch bleibt sie erhalten, wenn die Anwendung auf andere Computer kopiert wird.

@end

## Fehlerbehebung

### UAC-Eingabeaufforderungen werden nicht angezeigt

Wenn Sie `requireAdministrator` festgelegt haben, aber keine UAC-Eingabeaufforderungen angezeigt werden:

- Überprüfen Sie, ob das Manifest ordnungsgemäß in Ihre ausführbare Datei eingebettet ist
- Stellen Sie sicher, dass die Anwendung nicht von einem Prozess gestartet wird, der bereits mit erhöhten Rechten ausgeführt wird
- Stellen Sie sicher, dass die Manifestsyntax gültiges XML ist

### Anwendung startet nicht

Wenn Ihre Anwendung nach Änderungen an der UAC-Konfiguration nicht mehr startet:

- Prüfen Sie die Manifestsyntax auf XML-Fehler
- Überprüfen Sie, ob der Wert der Ausführungsebene gültig ist
- Kehren Sie versuchsweise zu `asInvoker` zurück, um das Problem einzugrenzen

### Abweichendes Verhalten auf verschiedenen Computern

Wenn sich das UAC-Verhalten auf verschiedenen Computern unterscheidet:

- Stellen Sie sicher, dass das Manifest in die ausführbare Datei eingebettet ist und nicht als externe Datei vorliegt
- Stellen Sie sicher, dass die ausführbare Datei nach dem Build nicht geändert wurde
- Überprüfen Sie, ob die Windows-UAC-Einstellungen auf dem Zielcomputer aktiviert sind
