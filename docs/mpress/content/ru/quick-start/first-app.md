---
title: "Ваше первое приложение"
description: "Создайте работающее приложение Wails за 10 минут"
slug: "quick-start/first-app"
sourcePath: "quick-start/first-app.md"
---

Мы создадим простое приложение с приветствием, демонстрирующее основные концепции Wails:

- Бэкенд на Go, управляющий логикой
- Фронтенд, вызывающий функции Go
- Типобезопасные привязки
- Горячая перезагрузка во время разработки

**Время выполнения:** 10 минут

@note{type="tip" title="Совет по повышению производительности для пользователей Windows 11"}
Рассмотрите возможность хранения проектов на [Dev Drive](https://learn.microsoft.com/en-us/windows/dev-drive/). Диски Dev Drive оптимизированы для рабочих нагрузок разработчиков и по сравнению с обычными дисками NTFS могут значительно сократить время сборки и повысить скорость доступа к диску — на величину до 30%.

@end

## Создание проекта

@steps
### Создайте проект
```bash
wails3 init -n myapp
cd myapp
```

Будет создан новый проект с шаблоном Vanilla + Vite по умолчанию (HTML/CSS/TypeScript и сборщик Vite).

@note{type="tip" title="Другие шаблоны"}
Выберите `-t react`, `-t vue` или `-t svelte` в зависимости от предпочитаемого фреймворка. По умолчанию эти шаблоны используют TypeScript; для обычного JavaScript используйте `-t vanilla-js` или `-t react-js`. Выполните `wails3 init -l`, чтобы просмотреть все доступные шаблоны, или [подключите собственный фронтенд-фреймворк](/guides/dev/frontend-frameworks/).

@end

### Разберитесь в структуре проекта
```
myapp/
├── main.go              # Application entry point
├── greetservice.go      # Greet service
├── frontend/            # Your UI code
│   ├── index.html       # HTML entry point
│   ├── src/
│   │   └── main.ts      # Frontend TypeScript
│   ├── public/
│   │   └── style.css    # Styles
│   ├── package.json     # Frontend dependencies
│   ├── tsconfig.json    # TypeScript configuration
│   └── vite.config.ts   # Vite bundler config
├── build/               # Build configuration
└── Taskfile.yml         # Build tasks
```

### Запустите приложение
```bash
wails3 dev
```

@note{type="info" title="Первый запуск"}
Первый запуск может занять больше времени, чем ожидалось, поскольку устанавливаются зависимости фронтенда, создаются привязки и выполняются другие операции. Последующие запуски проходят гораздо быстрее.

@end

Откроется приложение с интерфейсом приветствия. Введите своё имя и нажмите «Приветствовать» — бэкенд на Go обработает введённые данные и вернёт приветствие.

@end

## Как это работает

Разберём код, благодаря которому всё это работает.

### Бэкенд на Go

Откройте `greetservice.go`:

```go {title="greetservice.go"}
package main

import (
	"fmt"
)

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
```

**Основные концепции:**

1. **Сервис** — структура Go с экспортируемыми методами
2. **Экспортируемый метод** — имя `Greet` начинается с прописной буквы, поэтому метод доступен фронтенду
3. **Простая логика** — принимает имя и возвращает приветствие
4. **Типобезопасность** — типы входных и выходных данных определены

@note{type="tip" title="Сервисы и привязки"}
**Сервисы** — это автономные модули Go, предоставляющие функциональность фронтенду. Они представляют собой обычные структуры Go с экспортируемыми методами, которые вы регистрируете в поле `Services` конфигурации приложения.

**Привязки** — это автоматически создаваемый SDK для TypeScript/JavaScript, позволяющий фронтенду вызывать эти сервисы. При выполнении `wails3 dev` или `wails3 build` Wails анализирует зарегистрированные сервисы и создаёт типобезопасные привязки в `frontend/bindings/`.

Сервисы можно представить как API бэкенда, а привязки — как клиентскую библиотеку, которая взаимодействует с ним.

@end

### Регистрация сервиса

Откройте `main.go` и найдите регистрацию сервиса:

```go {title="main.go" highlight="4-6"}
err := application.New(application.Options{
    Name: "myapp",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    // ... other options
})
```

Так вы зарегистрируете `GreetService` в Wails, сделав все его экспортируемые методы доступными фронтенду.

### Фронтенд

Откройте `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
import {GreetService} from "../bindings/changeme";

window.greet = async () => {
    const nameElement = document.getElementById('name');
    const resultElement = document.getElementById('result');

    const name = nameElement.value;
    if (!name) {
        return;
    }

    try {
        const result = await GreetService.Greet(name);
        resultElement.innerText = result;
    } catch (err) {
        console.error(err);
    }
};
```

**Основные концепции:**

1. **Автоматически созданные привязки** — `GreetService` импортируется из созданного кода
2. **Типобезопасные вызовы** — имена и сигнатуры методов соответствуют коду Go
3. **Асинхронность по умолчанию** — все вызовы Go возвращают Promise
4. **Обработка ошибок** — ошибки из Go перехватываются конструкцией try/catch

@note{type="info" title="Где находятся привязки?"}
Созданные привязки находятся в `frontend/bindings/`. Они создаются автоматически при выполнении `wails3 dev` или `wails3 build`.

**Никогда не редактируйте эти файлы вручную** — они создаются заново при каждой сборке.

@end

## Настройка приложения

Добавим новую функцию, чтобы разобраться в рабочем процессе.

### Добавление функции «Приветствовать нескольких»

@steps
### Добавьте метод в GreetService
Добавьте следующий код в `greetservice.go`:

```go {title="greetservice.go"}
func (g *GreetService) GreetMany(names []string) []string {
    greetings := make([]string, len(names))
    for i, name := range names {
        greetings[i] = fmt.Sprintf("Hello %s!", name)
    }
    return greetings
}
```

### Приложение будет автоматически пересобрано
Сохраните файл, и `wails3 dev` автоматически пересоберёт код Go и перезапустит приложение.

@note{type="info" title="Автоматическая пересборка"}
Изменения в коде Go запускают автоматическую пересборку и перезапуск. Изменения во фронтенде применяются с помощью горячей перезагрузки без перезапуска.

@end

### Используйте метод во фронтенде
Добавьте следующий код в `frontend/src/main.js`:

```javascript {title="frontend/src/main.js"}
window.greetMany = async () => {
    const names = ['Alice', 'Bob', 'Charlie'];
    const greetings = await GreetService.GreetMany(names);
    console.log(greetings);
};
```

Откройте консоль браузера и вызовите `greetMany()` — вы увидите массив приветствий.

@end

## Сборка для рабочей среды

Когда приложение будет готово к распространению:

```bash
wails3 build
```

**Что при этом происходит:**

- Код Go компилируется с оптимизациями
- Фронтенд собирается для рабочей среды (с минификацией)
- В `bin/` создаётся нативный исполняемый файл

@tabs{sync-key="os"}
[Windows]
**Результат:** `bin/myapp.exe`

Запустите двойным щелчком. Зависимости не требуются (WebView2 входит в состав Windows).

[macOS]
**Результат:** `bin/myapp.app`

Перетащите в папку Applications или запустите двойным щелчком.

[Linux]
**Результат:** `bin/myapp`

Запустите с помощью `./bin/myapp` или создайте файл `.desktop` для средства запуска.

@end

@note{type="tip" title="Кроссплатформенная сборка"}
Хотите выполнить сборку для других платформ? См. [Кроссплатформенная сборка →](/guides/build/cross-platform/)

@end

## Что мы узнали

**Структура проекта**

- `main.go` для серверной части на Go
- `frontend/` для кода пользовательского интерфейса
- `Taskfile.yml` для задач сборки

**Сервисы**

- Создавайте структуры Go с экспортируемыми методами
- Регистрируйте с помощью `application.NewService()`
- Методы автоматически становятся доступны во фронтенде

**Привязки**

- Автоматически создаваемые определения TypeScript
- Типобезопасные вызовы функций
- Асинхронность по умолчанию (Promise)

**Процесс разработки**

- `wails3 dev` для горячей перезагрузки
- При изменениях в коде Go приложение автоматически пересобирается и перезапускается
- Изменения во фронтенде мгновенно применяются с помощью горячей перезагрузки

---

**Есть вопросы?** Присоединяйтесь к [Discord](https://discord.gg/JDdSxwjhGf) и задайте их сообществу.
