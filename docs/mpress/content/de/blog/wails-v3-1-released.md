---
title: "Wails v3.1.0: Mehr Go, weniger Build-Hürden"
description: "Experimentelle Desktop-Builds ohne CGO und Wake: eine einfachere Möglichkeit, Wails-Anwendungen zu bauen, zu paketieren und zu signieren."
authors: ["leaanthony"]
tags: ["wails", "v3", "release"]
date: "2026-11-01"
slug: "blog/wails-v3-1-released"
sourcePath: "blog/wails-v3-1-released.md"
---

Bei Wails war diese Woche viel los: Builds ohne CGO, eine einfachere Projekteinrichtung und weniger Hürden beim Paketieren und Signieren.

Mit Wails v3.1.0 beginnen unsere wöchentlichen Feature-Releases, durch die du regelmäßig neue Funktionen erhältst. Fehlerkorrekturen veröffentlichen wir weiterhin als nächtliche Patch-Updates, sobald sie bereit sind.

## Korrekturen und Optimierungen dieser Woche

Diese Woche gab es im Framework, in der CLI, in den Build-Werkzeugen und in der CI **100 Korrekturen, 21 Leistungsverbesserungen und vier Änderungen zur Sicherheitshärtung**. Zu den wichtigsten Änderungen gehören eine zuverlässigere Fensterpositionierung bei gemischten DPI-Werten, Korrekturen für den Linux-Infobereich und den Dunkelmodus, eine schlankere Bereitstellung von Assets und ein geringerer Speicherverbrauch beim Paketieren von MSIX-Dateien. [Einige der Änderungen findest du unten](#einige-highlights-dieser-woche).

## Experimentelle Builds ohne CGO

Mit purego kann Go native Systembibliotheken ohne CGo aufrufen, die übliche Brücke von Go zu C-Code. Wir begannen kurz nach dem Erscheinen mit ersten Untersuchungen, doch die Pflege zweier Backends war eine Herausforderung. Da purego inzwischen ausgereifter ist und generative KI diese Pflege besser handhabbar macht, können wir nun experimentelle Backends für macOS und Linux bereitstellen.

Du kannst Desktop-Anwendungen jetzt ohne C-Compiler bauen und kreuzkompilieren:

```sh
WAILS_PUREGO=1 wails3 build
```

Nachdem Sie Ihre Frontend-Assets erstellt haben, können Sie auch direkt mit Go bauen:

```sh
CGO_ENABLED=0 go build -tags purego,production .
```

Das macOS-Backend verwendet AppKit und WebKit; unter Linux werden GTK4/WebKitGTK 6.0 und GTK3 unterstützt. Benutzer benötigen weiterhin native Laufzeitbibliotheken, und für das Paketieren und Signieren gelten separate Anforderungen. CGo bleibt der Standard, während wir Feedback sammeln.

[Details zur Implementierung und zu den Plattformen](https://github.com/wailsapp/wails/pull/6521)

## Warum Wake?

Wails v2 verbarg zu viel vom Build-Prozess. In v3 machten Taskfiles alles zugänglich, doch die meisten Benutzer möchten lediglich einige Anpassungen vornehmen, ohne ein Build-System pflegen zu müssen.

Wake bündelt gängige Optionen in `wails.yaml`: Anwendungen, zugehörige CLI-Werkzeuge, Server, Zielplattformen und Paketierung. Task bleibt unser zentrales Build-System und steht dir zur Verfügung, wenn du seine Flexibilität benötigst.

Wake unterstützt jetzt Live-Neuladen, inkrementelle Builds, zwischengespeicherte Ausgaben und parallele Paketierung. Du hast bereits Binärdateien? `--skip-build` paketiert sie, ohne dass Go oder Frontend-Werkzeuge installiert sein müssen.

Aktiviere Wake mit `WAILS_EXP_USE_WAKE=1`.

### Geführter Einstieg

Wenn Wake aktiviert ist, führt dich `wails3 init -ui` durch die Auswahl von Backend, Zielplattform und Paket, prüft deine Werkzeuge und speichert die Projektkonfiguration.

@image{src="/assets/screenshots/wake/application.png" alt="Der Wake-Projektassistent mit Optionen für Desktop-Anwendung und Backend, einschließlich des experimentellen Pure Go." expand}

Wähle für jede Zielplattform Paketformate und Signierungsoptionen aus:

@image{src="/assets/screenshots/wake/packages.png" alt="Paketauswahl in Wake für Linux mit den Formaten AppImage, DEB, RPM, Void und Arch Linux sowie der Option, Release-Ausgaben zu signieren." expand}

Die folgende Werkzeugprüfung zeigt einen standardmäßigen Linux-Build:

@image{src="/assets/screenshots/wake/tools.png" alt="Wake prüft den C-Compiler, Go, npm und pkg-config für einen standardmäßigen Linux-Build." expand}

### Was kann ich auf meinem Rechner bauen?

Für Kompilierung, Paketierung und Signierung gelten unterschiedliche Anforderungen. Folgendes kannst du auf den jeweiligen Hostsystemen ausführen:

@table{header search sort column-separators align=["left", "centre", "centre", "centre", "left"]}
| Vorgang | Windows-Host | macOS-Host | Linux-Host | Voraussetzungen |
| Desktop-Anwendungen mit purego kompilieren | ✓ | ✓ | ✓ | Go und das experimentelle purego-Backend für macOS-/Linux-Zielplattformen |
| MSIX-Pakete erstellen | ✓ | ✓ | ✓ | Gebaute Anwendung; kein Windows SDK erforderlich |
| NSIS-Installationsprogramme erstellen | ✓ | ✓ | ✓ | Gebaute Anwendung und `makensis` |
| macOS-Anwendungen und DMGs erstellen | ✓ | ✓ | ✓ | Gebaute Anwendung; Go übernimmt die Paketerstellung |
| DEB-, RPM-, Arch- und Void-Pakete erstellen | ✓ | ✓ | ✓ | Gebaute Anwendung |
| AppImages erstellen | ✓ | ✓ | ✓ | AppImage-Laufzeit der Zielplattform; Linux-GUI-Laufzeitverzeichnis für plattformübergreifende Paketierung |
| Windows-Pakete mit einer exportierten Identität signieren | ✓ | ✓ | ✓ | Exportiertes RSA-PFX-Zertifikat und privater Schlüssel |
| Apple-Anwendungen und DMGs mit einer exportierten Identität signieren | ✓ | ✓ | ✓ | Exportierte RSA-PKCS#12-Identität; siehe Validierungshinweis unten |
| DEB-, RPM- und Arch-Pakete signieren | ✓ | ✓ | ✓ | Exportierter OpenPGP-Schlüssel |
| Zur Apple-Notarisierung einreichen und Tickets anhängen | ✓ | ✓ | ✓ | Apple-Team-API-Schlüssel und eine entsprechend signierte Anwendung |
| Apples `codesign` oder Schlüsselbundprofile verwenden | – | ✓ | – | macOS und die entsprechende Signierungsidentität |
@end

CGo-Builds benötigen weiterhin eine C-Toolchain. Benutzerdefinierte Hooks, die DMGs verändern, erfordern macOS und `hdiutil`.

### Signierung und Notarisierung

Wake unterstützt die Signierung von Windows-, Apple- und Linux-Paketen, einschließlich Ad-hoc-Signaturen für Apple. Die plattformübergreifende Signierung verwendet exportierte Identitäten; Identitäten aus Zertifikatspeichern und Schlüsselbünden benötigen weiterhin native Werkzeuge. Anmeldeinformationen auf Benutzerebene halten private Schlüssel aus der Projektkonfiguration heraus. Die Signierung von AppImage und Void ist noch nicht enthalten.

Die Apple-Notarisierung erfolgt mit Go und einem Team-API-Schlüssel. Wir haben die Einreichung und das Anhängen von Tickets unter Linux getestet; Gatekeeper hat das Ergebnis akzeptiert. Diese Anwendung wurde mit Apples [`codesign`](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/Procedures/Procedures.html) signiert. Ob Apple die von Wake mit Go erzeugten Developer-ID-Signaturen akzeptiert, ist noch nicht verifiziert.

### Umstieg auf Wake

Die Migration unterstützt alle neun integrierten Vorlagen. Sie ist unumkehrbar: Das Aktivieren von `wails.yaml` ersetzt das Taskfile-Routing. Nicht unterstützte benutzerdefinierte Befehle und Überschreibungen werden abgelehnt und müssen überprüft werden.

[Wake-Funktionen, Anforderungen und Migrationsdetails](https://github.com/wailsapp/wails/pull/6531)

## Einige Highlights dieser Woche

- **Fenster erscheinen auf Displays mit gemischten DPI-Werten an der richtigen Stelle.** Bildschirmgrenzen und Positionen werden jetzt auf macOS- und Linux-Systemen mit unterschiedlichen Anzeigeskalierungen korrekt verarbeitet. [#6346](https://github.com/wailsapp/wails/pull/6346)
- **Linux-Infobereich und Theme-Erkennung funktionieren zuverlässiger.** Nach dem Anwendungsstart erstellte Symbole für den Infobereich gehen nicht mehr verloren, Tooltips werden befüllt und die Erkennung des Dunkelmodus liest die Darstellungseinstellung des Desktops bereits beim Start aus.
- **Go-Event-Callbacks behalten für jeden Listener ihre Reihenfolge bei.** Jeder Listener wird seriell in der Reihenfolge der Auslösung ausgeführt; verschiedene Listener können gleichzeitig ausgeführt werden, und für gleichzeitige Auslösungen ist keine Reihenfolge festgelegt. Event-Hooks werden unter einer Sperre als Snapshot erfasst, um ein Datenrennen zu vermeiden. [#6275](https://github.com/wailsapp/wails/pull/6275), [#6276](https://github.com/wailsapp/wails/pull/6276)
- **Der Build-Cache unterstützt mehr reale Projekte.** Leerzeichen in Anführungszeichen oder mit Escape-Zeichen führen nicht mehr zu Fehltreffern im Go-Ausgabe-Cache, und gleichzeitige Cache-Speichervorgänge werden serialisiert, um beschädigte Dateien zu verhindern. [#6495](https://github.com/wailsapp/wails/pull/6495), [#6488](https://github.com/wailsapp/wails/pull/6488)
- **Die Asset-Bereitstellung führt weniger unnötige Arbeit aus.** Bei bekannten Dateiendungen wird die MIME-Erkennung übersprungen, und bei deaktivierter Debug-Protokollierung werden keine Nachrichten mehr vorbereitet, die anschließend verworfen würden. Beide Änderungen reduzieren den Arbeitsaufwand im Anfragepfad. [#6414](https://github.com/wailsapp/wails/pull/6414), [#6420](https://github.com/wailsapp/wails/pull/6420)
- **Die MSIX-Paketierung vermeidet Puffer in der Größe der Binärdatei.** Ausführbare Dateien werden in das Paket gestreamt, statt vollständig in den Arbeitsspeicher eingelesen zu werden. Beim Kopieren einer repräsentativen Datei mit 16 MiB wurden zuvor etwa 16.8 MB pro Vorgang reserviert. [#6472](https://github.com/wailsapp/wails/pull/6472)

## Probiere es aus

Beide Funktionen sind experimentell. Probiere einen plattformübergreifenden Build aus, paketiere deine Anwendung und teile uns mit, wo du auf Hindernisse stößt. Vielen Dank an alle, die testen, Probleme melden und Beiträge leisten.

Ein besonderer Dank gilt beesophie aus unserer Discord-Community. Sie hat [Cook](https://github.com/SimTXNG/cook) entwickelt, ein beeindruckendes Build-System, das uns einige Lücken in unserem eigenen System aufgezeigt hat. Die Möglichkeiten von Cook halfen uns zu erkennen, wo wir das Bauen mit Wails vereinfachen konnten, und waren eine echte Inspiration für diese Arbeit.

Danke, Beesophie!
