---
title: "Использование других фронтенд-фреймворков"
description: "Как использовать фреймворк без встроенного шаблона, разместив собственный проект Vite в каталоге frontend"
slug: "guides/dev/frontend-frameworks"
sourcePath: "guides/dev/frontend-frameworks.md"
---

Wails поставляется со встроенными начальными шаблонами для намеренно ограниченного набора фреймворков:

| Шаблон | Язык |
| --- | --- |
| `vanilla` | TypeScript (по умолчанию) |
| `vanilla-js` | JavaScript |
| `react` | TypeScript |
| `react-js` | JavaScript |
| `vue` | TypeScript |
| `svelte` | TypeScript |

Но это не значит, что выбор ограничен только ими. Фронтенд приложения Wails — **всего лишь веб-проект**: подойдёт всё, что собирается в статические HTML/CSS/JS. Если для выбранного вами фреймворка (Solid, Preact, Lit, Qwik, SvelteKit, Angular, …) нет шаблона, вы можете создать заготовку проекта самостоятельно за пару минут.

## Как используется каталог `frontend/`

Для Wails неважно, какой фреймворк находится в `frontend/`. Требуется лишь соблюдение небольшого контракта, не зависящего от фреймворка:

- **`frontend/dist/` поставляется с приложением.** `main.go` встраивает собранный фронтенд с помощью `//go:embed all:frontend/dist` и раздаёт его через сервер ресурсов. Результатом сборки должен быть статический пакет в `frontend/dist/` (каталоге вывода Vite по умолчанию).
- **Сборкой управляет `frontend/package.json`.** При выполнении `wails3 build` Wails запускает скрипт фронтенда `build`; при выполнении `wails3 dev` запускается `dev`, а сервер разработки Vite проксируется для горячей перезагрузки.
- **Привязки генерируются в `frontend/bindings/`.** Wails анализирует зарегистрированные сервисы Go и записывает туда типобезопасный SDK. Импортируйте его как любой другой модуль:
  ```js
  import { GreetService } from "./bindings/changeme";
  ```


- **Сервер разработки работает на фиксированном порту.** `wails3 dev` проксирует Vite через порт, указанный в `WAILS_VITE_PORT` (по умолчанию `9245`), поэтому установите `server.port` в значение этого порта и одновременно укажите `strictPort: true`. Это обычная конфигурация Vite — плагин Wails здесь не используется.
- **Необязательно — типизированные пользовательские события.** Встроенные шаблоны также регистрируют плагин `@wailsio/runtime/plugins/vite`. Он нужен, только если вы используете *типизированные* пользовательские события: плагин внедряет сгенерированные определения типов событий в среду выполнения и не позволяет выполнить сборку, пока не будут сгенерированы привязки. Если вы используете только строковый API `Events.On("time", …)`, плагин можно не добавлять.

Всё остальное — компоненты, маршрутизация, состояние и стили — полностью определяется вашим фреймворком.

## Создание заготовки проекта для любого фреймворка с помощью Vite

Самый быстрый способ — начать со встроенного шаблона (чтобы получить `main.go`, `Taskfile`, ресурсы сборки и работающий сервис Go), а затем заменить `frontend/` новым проектом Vite для вашего фреймворка.

@steps
### Создайте проект из шаблона по умолчанию
```bash
wails3 init -n myapp
cd myapp
```

### Замените `frontend/` приложением Vite для вашего фреймворка
Vite позволяет создать заготовку проекта для большинства фреймворков одной командой. Выберите шаблон:

```bash
# From the project root — e.g. Solid, Preact, Lit, Svelte, Vue, React, Vanilla
rm -rf frontend
npm create vite@latest frontend -- --template solid
```

Замените `solid` любым шаблоном Vite: `preact`, `lit`, `svelte`, `vue`, `react`, `vanilla` или их вариантами `-ts` (`solid-ts`, `preact-ts`, …).

### Установите среду выполнения и укажите Vite сервер разработки Wails
```bash
cd frontend
npm install @wailsio/runtime
```

`@wailsio/runtime` предоставляет API JavaScript (`Events`, `Browser`, диалоговые окна, …). Единственное *обязательное* изменение в `vite.config` — порт сервера разработки, чтобы `wails3 dev` мог его найти:

```ts {title="frontend/vite.config.ts"}
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
});
```

Только если вы планируете использовать **типизированные пользовательские события**, также добавьте плагин: он внедряет сгенерированные типы событий и требует, чтобы привязки существовали до начала сборки:

```ts {title="frontend/vite.config.ts" highlight="2,5"}
import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [wails("./bindings")],
  server: { host: "127.0.0.1", port: Number(process.env.WAILS_VITE_PORT) || 9245, strictPort: true },
});
```

### Вызывайте сервисы Go
Один раз сгенерируйте привязки, а затем импортируйте их в любые компоненты:

```bash
wails3 generate bindings
```

```js
import { GreetService } from "./bindings/changeme";

const greeting = await GreetService.Greet("World");
```

### Запустите проект
```bash
wails3 dev
```

@end

@note{type="tip" title="Создайте проект на чистом JavaScript"}
Та же команда создаёт заготовку проекта без TypeScript — просто используйте шаблон Vite без `-ts`:

```bash
npm create vite@latest frontend -- --template solid
```

@end

## Фреймворки с собственными средствами создания проектов

Для некоторых фреймворков проекты создаются не с помощью шаблонов Vite `create`, а собственными инструментами. Они по-прежнему совместимы: просто создайте заготовку проекта штатной командой фреймворка, а затем добавьте плагин Wails:

- **SvelteKit:** `npx sv create frontend`. Используйте статический адаптер (`@sveltejs/adapter-static`), чтобы результатом сборки был статический пакет, и отключите SSR.
- **Qwik:** `npm create qwik@latest`. Используйте статический адаптер (SSG).
- **Angular:** создайте заготовку с помощью `ng new`, задайте для `outputPath` значение `dist` и настройте скрипт сборки на использование `ng build`.

Правило всегда одно и то же: создавайте статическую сборку в `frontend/dist/`, сохраняйте плагин Vite `@wailsio/runtime` (или импортируйте среду выполнения напрямую) и импортируйте привязки Go из `frontend/bindings/`.

@note{type="info"}
Если вы создадите хорошо проработанную конфигурацию для какого-либо фреймворка, рассмотрите возможность опубликовать её как [пользовательский шаблон](/guides/advanced/custom-templates/), чтобы другие разработчики могли напрямую выполнить для него `wails3 init -t`.

@end
