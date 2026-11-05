---
title: "CFN Tracker"
description: "Настольное приложение, созданное с помощью Wails"
slug: "community/showcase/cfntracker"
sourcePath: "community/showcase/cfntracker.md"
---

![CFN Tracker](/assets/showcase-images/cfntracker.webp)

[CFN Tracker](https://github.com/williamsjokvist/cfn-tracker) — отслеживайте текущие матчи любого профиля CFN в Street Fighter 6 или V. Чтобы начать, посетите [веб-сайт](https://cfn.williamsjokvist.se/).

## Возможности

- Отслеживание матчей в реальном времени
- Хранение журналов матчей и статистики
- Вывод текущей статистики в OBS с помощью источника «Браузер»
- Поддержка SF6 и SFV
- Возможность создавать собственные темы для источника «Браузер» в OBS с помощью CSS

### Основные технологии, используемые вместе с Wails

- [Task](https://github.com/go-task/task) — оболочка для Wails CLI, упрощающая использование распространённых команд
- [React](https://github.com/facebook/react) — выбран благодаря богатой экосистеме (radix, framer-motion)
- [Bun](https://github.com/oven-sh/bun) — используется благодаря быстрому разрешению зависимостей и высокой скорости сборки
- [Rod](https://github.com/go-rod/rod) — автоматизация безголового браузера для аутентификации и периодической проверки изменений
- [SQLite](https://github.com/mattn/go-sqlite3) — используется для хранения матчей, сеансов и профилей
- [События, отправляемые сервером](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) — HTTP-поток для передачи обновлений отслеживания источникам «Браузер» в OBS
- [i18next](https://github.com/i18next/) — с серверным коннектором для предоставления объектов локализации из слоя Go
- [xstate](https://github.com/statelyai/xstate) — конечные автоматы для процессов аутентификации и отслеживания
