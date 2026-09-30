---
title: "Snippet Expander"
description: "Une application de bureau créée avec Wails"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Capture d’écran de Snippet Expander](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Capture d’écran de la fenêtre Select Snippet de Snippet Expander

![Capture d’écran de Snippet Expander](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Capture d’écran de l’écran Add Snippet de Snippet Expander

![Capture d’écran de Snippet Expander](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Capture d’écran de la fenêtre Search & Paste de Snippet Expander

[Snippet Expander](https://snippetexpander.org) est « votre petit assistant d’extraits de texte extensibles », pour Linux.

Snippet Expander comprend une application avec interface graphique créée avec Wails pour gérer les extraits et les paramètres, ainsi qu’un mode de fenêtre Search & Paste permettant de sélectionner et de coller rapidement un extrait.

L’interface graphique basée sur Wails, l’interface en ligne de commande en go-lang et le démon d’expansion automatique en vala-lang communiquent tous avec un démon en go-lang via D-Bus. Ce démon effectue l’essentiel du travail : il gère la base de données des extraits et les paramètres communs, et fournit notamment des services pour développer et coller les extraits.

Consultez le [code source](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38) pour découvrir comment l’application Wails envoie de l’interface utilisateur au backend des messages qui sont ensuite transmis au démon. Vous verrez également comment elle s’abonne à un événement D-Bus afin de détecter les modifications apportées aux extraits par une autre instance de l’application ou de l’interface en ligne de commande, puis de les afficher instantanément dans l’interface utilisateur au moyen d’un événement Wails.
