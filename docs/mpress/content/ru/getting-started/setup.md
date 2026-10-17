---
title: "Настройка"
slug: "getting-started/setup"
sourcePath: "getting-started/setup.md"
---

@note{type="caution" title="Экспериментальная функция"}
Мастер настройки появился недавно и в основном протестирован в Linux. Если вы столкнётесь с проблемами, [сообщите о них](https://github.com/wailsapp/wails/issues/4904) и вместо мастера выполните [инструкции по установке вручную](/getting-started/installation/#platform-specific-dependencies).

@end

## Быстрый старт

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# 2. Run the setup wizard
wails3 setup
```

Мастер открывается в браузере и помогает проверить зависимости, задать настройки проекта по умолчанию и при необходимости настроить сборку для нескольких платформ.

После этого можно приступать к созданию первого проекта:

```bash
wails3 init -n myapp -t vanilla
cd myapp && wails3 dev
```

## Возможности мастера

- **Проверка зависимостей** — проверяет Go, npm и инструменты платформы
- **Настройка значений по умолчанию** — сведения об авторе, префикс идентификатора пакета и предпочитаемые шаблоны
- **Кроссплатформенная сборка** — необязательная настройка Docker для сборки на любой основной системе
- **Подписание кода** — необязательная настройка для macOS, Windows и Linux

Конфигурация сохраняется в `~/.config/wails/config.yaml` и используется командой `wails3 init`.

## Подкоманды

```bash
wails3 setup signing      # Configure code signing
wails3 setup entitlements # Configure macOS entitlements
```

## Возникли проблемы?

1. Для диагностики проблем выполните `wails3 doctor`
2. Выполните [инструкции по установке вручную](/getting-started/installation/#platform-specific-dependencies)
3. [Сообщите о проблеме](https://github.com/wailsapp/wails/issues/4904) и приложите вывод команды `wails3 doctor`
