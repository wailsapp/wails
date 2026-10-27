---
title: "Настройка"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Экспериментальная функция"}
Мастер настройки появился недавно и в основном протестирован в Linux. Если у вас возникнут проблемы, [сообщите о них](https://github.com/wailsapp/wails/issues/4904) и вместо мастера выполните [инструкции по установке вручную](/getting-started/installation/#platform-specific-dependencies).

@end

## Быстрый старт

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# 2. Run the setup wizard
wails3 setup
```

Мастер открывается в браузере и помогает проверить зависимости, задать параметры проекта по умолчанию и при необходимости настроить кроссплатформенную сборку.

После этого можно создать первый проект:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Возможности мастера

- **Проверка зависимостей** — проверяет Go, npm и инструменты для целевой платформы
- **Настройка параметров по умолчанию** — сведения об авторе, префикс идентификатора пакета и предпочтительные шаблоны
- **Кроссплатформенная сборка** — необязательная настройка Docker для сборки на любой хост-системе
- **Подписание кода** — необязательная настройка для macOS, Windows и Linux

Конфигурация сохраняется в `$XDG_CONFIG_HOME/wails/defaults.yaml`, а при невозможности — в `~/.config/wails/defaults.yaml`, и используется командой `wails3 init`. Она не связана с файлом `build/config.yml` проекта.

## Подкоманды

```bash
wails3 signing      # Configure code signing
wails3 entitlements # Configure macOS entitlements
```

## Возникли проблемы?

1. Выполните `wails3 doctor` для диагностики проблем
2. Выполните [инструкции по установке вручную](/getting-started/installation/#platform-specific-dependencies)
3. [Сообщите о проблеме](https://github.com/wailsapp/wails/issues/4904) и приложите вывод команды `wails3 doctor`
