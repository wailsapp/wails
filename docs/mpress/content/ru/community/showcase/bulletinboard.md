---
title: "BulletinBoard"
description: "Настольное приложение, созданное с помощью Wails"
slug: "community/showcase/bulletinboard"
sourcePath: "community/showcase/bulletinboard.md"
---

![BulletinBoard](/assets/showcase-images/bboard.webp)

Приложение [BulletinBoard](https://github.com/raguay/BulletinBoard) — это универсальная доска для статических сообщений и диалогов, позволяющих сценарию получать данные от пользователя. В приложении есть текстовый интерфейс (TUI) для создания новых диалогов, которые затем можно использовать для получения данных от пользователя. Приложение рассчитано на постоянную работу в вашей системе: оно показывает информацию по мере необходимости, а затем скрывается. У меня настроен процесс, который отслеживает файл в системе и при его изменении отправляет содержимое в BulletinBoard. Он отлично вписывается в мои рабочие процессы. Также предусмотрен [рабочий процесс Alfred](https://github.com/raguay/MyAlfred/blob/master/Alfred%205/EmailIt.alfredworkflow) для отправки информации в программу. Этот рабочий процесс также предназначен для работы с [EmailIt](https://github.com/raguay/EmailIt).
