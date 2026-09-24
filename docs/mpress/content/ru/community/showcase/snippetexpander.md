---
title: "Snippet Expander"
description: "Настольное приложение, созданное с помощью Wails"
slug: "community/showcase/snippetexpander"
sourcePath: "community/showcase/snippetexpander.md"
---

![Снимок экрана Snippet Expander](/assets/showcase-images/snippetexpandergui-select-snippet.png)

Снимок окна «Выбор сниппета» приложения Snippet Expander

![Снимок экрана Snippet Expander](/assets/showcase-images/snippetexpandergui-add-snippet.png)

Снимок экрана «Добавление сниппета» приложения Snippet Expander

![Снимок экрана Snippet Expander](/assets/showcase-images/snippetexpandergui-search-and-paste.png)

Снимок окна «Поиск и вставка» приложения Snippet Expander

[Snippet Expander](https://snippetexpander.org) — это «ваш маленький помощник для разворачивания текстовых сниппетов» в Linux.

Snippet Expander состоит из приложения с графическим интерфейсом, созданного с помощью Wails для управления сниппетами и настройками, с режимом окна «Поиск и вставка» для быстрого выбора и вставки сниппета.

Графический интерфейс на основе Wails, интерфейс командной строки на go-lang и служба автоматического разворачивания на vala-lang обмениваются данными со службой на go-lang через D-Bus. Эта служба выполняет основную часть работы: управляет базой данных сниппетов и общими настройками, а также предоставляет сервисы для разворачивания и вставки сниппетов и других операций.

Ознакомьтесь с [исходным кодом](https://git.sr.ht/~ianmjones/snippetexpander/tree/trunk/item/cmd/snippetexpandergui/app.go#L38), чтобы узнать, как приложение Wails отправляет сообщения из пользовательского интерфейса в серверную часть, откуда они затем передаются службе, а также подписывается на событие D-Bus, чтобы отслеживать изменения сниппетов, внесённые другим экземпляром приложения или интерфейсом командной строки, и мгновенно отображать их в пользовательском интерфейсе с помощью события Wails.
