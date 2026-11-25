---
title: "Feedback"
description: "So geben Sie Feedback und melden Probleme für Wails v3"
slug: "feedback"
sourcePath: "feedback.md"
---

Wir freuen uns über Ihr Feedback und ermutigen Sie ausdrücklich dazu! Suchen Sie bitte nach bestehenden Issues oder Diskussionen, bevor Sie neue erstellen. Sie können auf folgende Arten beitragen:

@tabs
[Fehler]
Wenn Sie einen Fehler finden, [erstellen Sie bitte ein Issue](https://github.com/wailsapp/wails/issues/new/choose) auf GitHub und verwenden Sie dafür die Vorlage für Fehlerberichte.

- Beschreiben Sie den Fehler eindeutig anhand eines einfachen, reproduzierbaren Beispiels. Falls aus der Dokumentation nicht klar hervorgeht, was *geschehen sollte*, erwähnen Sie dies im Bericht.
- Fügen Sie die Ausgabe von `wails3 doctor` in Ihren Bericht ein.
- Falls der Fehler in einem Verhalten besteht, das nicht der aktuellen Dokumentation entspricht, gehen Sie bitte außerdem wie folgt vor:
  - Aktualisieren Sie ein vorhandenes Beispiel im Verzeichnis `v3/examples` oder erstellen Sie ein neues Beispiel, das das Problem eindeutig veranschaulicht.
  - Erstellen Sie einen [PR](https://github.com/wailsapp/wails/pulls), der auf das Issue verweist.


@note{type="caution"}
*Denken Sie daran*: Unerwartetes Verhalten ist nicht zwangsläufig ein Fehler – möglicherweise entspricht es lediglich nicht Ihren Erwartungen. Verwenden Sie dafür `Suggestions`.

@end

Sie können Fehler auch im Discord-Kanal [#v3](https://discord.gg/bdj28QNHmT) diskutieren.

[Fehlerbehebungen]
Wenn Sie eine Fehlerbehebung oder eine Verbesserung der Dokumentation haben, gehen Sie bitte wie folgt vor:

- Erstellen Sie einen Pull Request im [Wails-Repository](https://github.com/wailsapp/wails) und befolgen Sie dabei die [Richtlinien für Beiträge](https://github.com/wailsapp/wails/blob/master/CONTRIBUTING.md).
- Verweisen Sie in der PR-Beschreibung auf alle zugehörigen Issues.

[Erweiterungen]
Neue Funktionen und Änderungen am öffentlichen Verhalten werden über einen Pull-Request-Entwurf für einen **WEP (Wails Enhancement Proposal)** vorgeschlagen, nicht über ein Feature-Request-Issue.

- Lesen Sie die Beschreibung des [WEP-Prozesses](https://github.com/wailsapp/wails/blob/v3.0.0/v3/wep/README.md).
- Kopieren Sie die Vorlage und erstellen Sie anschließend einen PR-Entwurf mit dem Titel `[WEP] <title>`, der ausschließlich den WEP und unterstützendes Material enthält.
- Sie können eine Idee zunächst informell in den [GitHub Discussions](https://github.com/wailsapp/wails/discussions) oder im Discord-Kanal [#v3](https://discord.gg/bdj28QNHmT) diskutieren. Für eine Entscheidung der Maintainer ist jedoch ein WEP-PR erforderlich.

[Unterstützen]
- Zeigen Sie Ihre Unterstützung für Fehler, WEPs und Diskussionen mit der :thumbsup:-Reaktion auf GitHub.
- Fügen Sie bitte *nicht* einfach Kommentare wie „+1“ oder „Ich auch“ hinzu.
- Fügen Sie einen Kommentar hinzu, wenn Sie etwas Wesentliches beitragen können, beispielsweise „Dieser Fehler betrifft auch ARM-Builds“ oder „Ein anderer Ansatz wäre …“.

@end

Bekannte Probleme und laufende Arbeiten finden Sie [hier](https://github.com/orgs/wailsapp/projects/6).
