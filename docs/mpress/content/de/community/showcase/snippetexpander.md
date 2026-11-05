---
title: "Snippet Expander"
description: "Eine mit Wails entwickelte Desktopanwendung"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Screenshot von Snippet Expander](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Screenshot des Fensters „Snippet auswählen“ von Snippet Expander

![Screenshot von Snippet Expander](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Screenshot der Ansicht „Snippet hinzufügen“ von Snippet Expander

![Screenshot von Snippet Expander](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Screenshot des Fensters „Suchen und Einfügen“ von Snippet Expander

[Snippet Expander](https://snippetexpander.org) ist „dein kleiner Helfer für erweiterbare Text-Snippets“ unter Linux.

Snippet Expander umfasst eine mit Wails entwickelte GUI-Anwendung zum Verwalten von Snippets und Einstellungen sowie einen Fenstermodus „Suchen und Einfügen“, mit dem sich ein Snippet schnell auswählen und einfügen lässt.

Die auf Wails basierende GUI, die go-lang-CLI und der automatische Erweiterungs-Daemon in vala-lang kommunizieren alle über D-Bus mit einem go-lang-Daemon. Der Daemon erledigt den Großteil der Arbeit: Er verwaltet die Datenbank mit Snippets und gemeinsamen Einstellungen und stellt unter anderem Dienste zum Erweitern und Einfügen von Snippets bereit.

Sieh dir den [Quellcode](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38) an, um zu erfahren, wie die Wails-App Nachrichten von der UI an das Backend sendet, die anschließend an den Daemon weitergeleitet werden. Außerdem abonniert sie ein D-Bus-Ereignis, um Änderungen an Snippets durch eine andere Instanz der App oder CLI zu überwachen und sie über ein Wails-Ereignis sofort in der UI anzuzeigen.
