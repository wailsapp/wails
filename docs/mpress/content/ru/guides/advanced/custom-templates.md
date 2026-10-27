---
title: "Создание пользовательских шаблонов"
description: "Как создавать, настраивать и размещать собственные шаблоны проектов Wails v3"
slug: "guides/advanced/custom-templates"
sourcePath: "guides/advanced/custom-templates.md"
---

Wails поставляется с набором встроенных шаблонов, но вы можете создавать собственные шаблоны и делиться ими с сообществом. Пользовательский шаблон — это обычный репозиторий Git. После его публикации любой пользователь сможет создать на его основе каркас проекта одной командой.

## Создание каркаса шаблона

Команда `wails3 generate template` создаёт каталог шаблона, готовый к настройке:

```bash
wails3 generate template -name MyTemplate
```

Все флаги:

| Флаг | Описание | Значение по умолчанию |
| --- | --- | --- |
| `-name` | Имя шаблона (обязательно) | — |
| `-author` | Имя автора | — |
| `-description` | Краткое описание, отображаемое в CLI | — |
| `-helpurl` | URL-адрес документации этого шаблона | — |
| `-version` | Начальная версия | `v0.0.1` |
| `-frontend` | Скопировать существующий каталог фронтенда в шаблон | — |
| `-dir` | Куда записать каталог шаблона | Текущий каталог |

Пример со всеми флагами:

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

Созданный каталог имеет следующую структуру:

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

@note{type="tip" title="Прочитайте NEXTSTEPS.md"}
Созданный файл `NEXTSTEPS.md` содержит подробные рекомендации по всем частям шаблона. Прочитайте его перед настройкой. Удалите его перед публикацией: он не должен попадать в проекты, созданные из вашего шаблона.

@end

## Настройка метаданных шаблона

Откройте `template.yaml` и задайте метаданные шаблона:

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

Поле `wailsVersion` **обязательно** и должно иметь значение `3`. Комментарий `# yaml-language-server` в начале файла включает автодополнение и встроенную проверку в VS Code (с [расширением YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) и IDE от JetBrains. Его можно оставить или удалить: на выполнение программы он не влияет.

## Настройка шаблона

### Фронтенд

Содержимое каталога `frontend/` без изменений копируется в каждый проект, созданный из вашего шаблона. Замените содержимое-заполнитель своим фронтендом:

@tabs
[Начать с нуля]
```bash
cd MyTemplate/frontend
npm create vite@latest .
```

Следуйте подсказкам, а затем установите зависимости:

```bash
npm install
```

[Использовать существующий проект]
Чтобы при создании шаблона за один шаг скопировать существующий фронтенд, передайте `-frontend`:

```bash
wails3 generate template -name MyTemplate -frontend ./my-app/frontend
```

Либо затем вручную скопируйте его в каталог `frontend/`.

@end

### Задачи сборки

Файл `Taskfile.tmpl.yml` определяет процесс сборки. Измените задачи `install:frontend:deps` и `build:frontend` в соответствии с инструментами сборки вашего фронтенда:

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Приложение Go

Файл `main.go.tmpl` служит точкой входа приложения. При создании проекта его обрабатывает механизм шаблонов Wails: переменные шаблона, такие как `{{.ProductName}}`, заменяются значениями, указанными пользователем.

Чтобы редактировать его как обычный файл Go с поддержкой IDE, временно переименуйте его в `main.go`, внесите изменения, а перед коммитом снова переименуйте в `main.go.tmpl`.

#### Переменные шаблона

В любом файле `.tmpl` доступны следующие переменные:

| Переменная | Описание | Пример |
| --- | --- | --- |
| `{{.ProjectName}}` | Имя проекта, указанное пользователем | `"MyApp"` |
| `{{.BinaryName}}` | Имя исполняемого файла | `"myapp"` |
| `{{.ProductName}}` | Отображаемое имя продукта | `"My Application"` |
| `{{.ProductDescription}}` | Описание продукта | `"An awesome application"` |
| `{{.ProductVersion}}` | Версия продукта | `"1.0.0"` |
| `{{.ProductCompany}}` | Название компании или имя автора | `"My Company Ltd"` |
| `{{.ProductCopyright}}` | Строка уведомления об авторских правах | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Дополнительные комментарии о продукте | `"Built with Wails"` |
| `{{.ProductIdentifier}}` | Идентификатор продукта в формате обратной DNS-записи | `"com.mycompany.myapp"` |
| `{{.ModulePath}}` | Путь модуля Go | `"github.com/you/myapp"` |
| `{{.WailsVersion}}` | Версия Wails, использованная для создания проекта | `"3.0.0"` |
| `{{.Typescript}}` | `true`, если имя шаблона оканчивается на `-ts` | `true` |
| `{{.Opn}}` | Литеральное значение `{{` — экранируйте внутри шаблонов | `{{` |
| `{{.Cls}}` | Литеральное значение `}}` — экранируйте внутри шаблонов | `}}` |

@note{type="tip"}
Любой файл в шаблоне может быть файлом `.tmpl`, включая файлы HTML, JSON и YAML. Файлы без суффикса `.tmpl` копируются без изменений.

@end

## Локальное тестирование шаблона

Перед публикацией протестируйте шаблон, создав проект по локальному пути:

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Затем убедитесь, что проект работает:

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Проверьте следующее:

- Горячая перезагрузка фронтенда работает
- После изменений в коде Go приложение пересобирается и перезапускается
- Рабочий исполняемый файл в `bin/` запускается и работает правильно

## Публикация на GitHub

@steps
### **Создайте общедоступный репозиторий GitHub** для своего шаблона. В корне репозитория должен находиться файл `template.yaml`.
### **Удалите `NEXTSTEPS.md`** — этот файл содержит рекомендации для авторов шаблонов и не должен попадать в проекты, которые пользователи создают на основе вашего шаблона.
### **Зафиксируйте изменения и отправьте их в удалённый репозиторий**, разместив содержимое каталога шаблона в корне репозитория:
```bash
git init
git add .
git commit -m "Initial template"
git remote add origin https://github.com/yourname/my-template.git
git push -u origin main
```

### **Создайте тег выпуска** согласно правилам семантического версионирования:
```bash
git tag v1.0.0
git push origin v1.0.0
```

@end

Теперь пользователи могут создавать проекты на основе вашего шаблона:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

@note{type="caution" title="Предупреждение о стороннем шаблоне"}
Когда пользователь устанавливает удалённый шаблон, Wails выводит предупреждение о том, что шаблон содержит сторонний код и проект Wails не несёт ответственности за его содержимое. Перед созданием проекта пользователь должен явно подтвердить продолжение.

Как автор шаблона, вы несёте ответственность за безопасность и корректность всего кода в нём.

Шаблоны — это код. Задачи `Taskfile.yml` представляют собой произвольные команды оболочки, а `wails3 dev`, `wails3 build` и `wails3 package` запускают их с правами пользователя вместе со всеми вызываемыми ими сценариями `package.json`. Перед выполнением любой из этих команд проверьте `Taskfile.yml`, каталог `build/` и `frontend/package.json` удалённого шаблона.

@end

## Рекомендации

- **Подготовьте понятный файл `README.md`** — он отображается пользователям после создания проекта. Объясните, как запускать, собирать и настраивать проект.
- **Заполните поле `helpurl`** — укажите ссылку на свой репозиторий или отдельную документацию. Пользователи увидят её в списке шаблонов Wails CLI.
- **Зафиксируйте версии зависимостей фронтенда** в `package.json`, чтобы обновления вышестоящих проектов не нарушали установку.
- **Проверьте шаблон перед созданием тега** — прежде чем объявлять о выпуске сообществу, создайте новый проект из выпуска с этим тегом.
- **Сохраните значение `wailsVersion: 3`** — это поле указывает Wails, для какой основной версии предназначен шаблон. Не изменяйте его.
- **Регулярно обновляйте шаблон** — поддерживайте зависимости в актуальном состоянии и тестируйте шаблон с новыми выпусками Wails.
