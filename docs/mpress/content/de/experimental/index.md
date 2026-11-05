---
title: "Experimentelle Funktionen"
description: "Ein Ort für laufende Experimente in Wails v3 – was sie sind, warum es sie gibt und wie Sie Feedback geben können."
slug: "experimental"
sourcePath: "experimental/index.md"
---

@note{type="caution" title="Hier wird experimentiert"}
Alles in diesem Abschnitt ist per Definition ein Experiment. Die hier beschriebenen Funktionen sind optional, standardmäßig deaktiviert und können zwischen Releases umgestaltet, umbenannt oder vollständig entfernt werden. Bauen Sie nichts, was für Ihr Projekt unverzichtbar ist, darauf auf, wenn Sie nicht auf häufige Änderungen vorbereitet sind.

@end

## Was „experimentell“ bedeutet

Wails entwickelt Experimente öffentlich. Ein Experiment ist eine Idee, die wir für vielversprechend genug halten, um sie Ihnen frühzeitig zur Verfügung zu stellen, zu der wir uns aber noch nicht vollständig bekannt haben. Wir veröffentlichen sie *deshalb*, weil wir aus dem praktischen Einsatz lernen möchten, bevor wir entscheiden, ob sie ein dauerhafter, unterstützter Bestandteil von Wails wird.

Daraus ergeben sich einige Merkmale, die für alles in diesem Abschnitt gelten:

- **Die Aktivierung erfolgt ausdrücklich.** Experimente ändern niemals das Standardverhalten von `wails3`. Sie aktivieren sie bewusst, üblicherweise mit einer Umgebungsvariablen oder einem Build-Flag. Solange sie deaktiviert sind, ändert sich nichts an Ihrem bestehenden Arbeitsablauf.
- **Das Experiment besteht möglicherweise nicht fort.** Einige Experimente werden zu stabilen Funktionen weiterentwickelt. Andere werden bis zur Unkenntlichkeit überarbeitet oder eingestellt. Wir probieren Dinge lieber öffentlich aus und lernen schnell daraus, als nur das zu veröffentlichen, dessen wir uns bereits sicher sind.
- **Die API ist noch nicht festgeschrieben.** Namen, Flags, Standardwerte und Verhalten können sich zwischen Releases ändern, während ein Experiment seine endgültige Form annimmt. Die Release Notes weisen auf die Änderungen hin, aber erwarten Sie nicht die Stabilitätsgarantien stabiler Funktionen.

## Wir möchten Ihr Feedback

Das ist der entscheidende Punkt. Ob Experimente fortgeführt oder eingestellt werden, hängt von den Rückmeldungen der Personen ab, die sie tatsächlich verwenden. Wenn Sie eines davon ausprobieren, möchten wir wirklich Folgendes wissen:

- Hat es für Ihr Projekt funktioniert? Wo gab es Probleme?
- War es schneller, übersichtlicher oder angenehmer – oder hat sich der Wechsel nicht gelohnt?
- Welche Voraussetzungen müssten erfüllt sein, damit Sie es standardmäßig verwenden würden?

Am hilfreichsten ist konkretes Feedback: Was haben Sie ausgeführt, was haben Sie erwartet und was ist tatsächlich geschehen? Für jedes Experiment gibt es in der Kategorie **Experimente** auf GitHub Discussions einen eigenen Thread:

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Diskussionen zu Experimenten" href="https://github.com/wailsapp/wails/discussions/categories/experiments" description="Suchen Sie den Thread zu dem von Ihnen verwendeten Experiment und teilen Sie uns mit, wie es funktioniert hat, was nicht funktioniert oder was fehlt."}
@end

## Aktuelle Experimente

@container{display="grid" columns="2" gap="1rem"}
@linkcard{title="Wake" href="/experimental/wake/" description="Ein alternativer, auf Wails abgestimmter Build-Runner für Ihre bestehenden Taskfiles. Schnellere inkrementelle Builds, strukturierte Ausgabe und standardmäßig parallele Ausführung."}
@linkcard{title="LLM-Steuerung (MCP)" href="/guides/mcp-service/" description="Ein integrierter Model Context Protocol-Server, mit dem LLM-Agenten eine laufende Wails-App untersuchen, testen und steuern können – ohne erforderlichen Benutzercode und über ein Build-Tag aktiviert."}
@end
