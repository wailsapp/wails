---
title: "Ресурсы выпусков GitHub для средства обновления"
description: "Как средство обновления Wails выбирает артефакты приложения из выпусков GitHub и исключает пакеты установщиков."
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

Поставщик GitHub Releases выбирает ресурс выпуска с помощью настроенного `AssetMatcher`. Если `AssetMatcher` имеет значение `nil`, используется `github.DefaultAssetMatcher`.

## Сопоставление по умолчанию

Средство сопоставления по умолчанию ищет обозначения текущей платформы и архитектуры в имени файла каждого ресурса. Оно распознаёт распространённые псевдонимы архитектур, включая:

- `amd64`, `x86_64` и `x64`
- `arm64` и `aarch64`
- `386`, `i386`, `x86` и `ia32`

Сопутствующие файлы, например подписи и контрольные суммы, игнорируются.

## Ресурсы установщиков

Выпуск GitHub может содержать как исполняемый файл приложения, используемый средством обновления, так и обычный установщик, предназначенный для первоначальной установки. Средство сопоставления по умолчанию игнорирует ресурсы, имя файла которых в нижнем регистре:

- содержит `-installer.`
- содержит `_installer.`
- в точности совпадает с `installer.exe`

Например, для следующих ресурсов Windows:

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher` выбирает `myapp-windows-amd64.exe` и игнорирует установщик. Это не позволяет средству обновления заменить работающее приложение исполняемым файлом установщика NSIS или установщика в аналогичном формате.

Проверка намеренно имеет узкую область действия. Имена приложений, которые лишь содержат слово `installer`, остаются допустимыми, включая:

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## Пользовательские схемы именования

Настройте `AssetMatcher`, если ресурсы выпуска не следуют соглашению об именовании по платформе и архитектуре либо если требуется другая фильтрация установщиков:

```go
import (
    "strings"

    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

gh, err := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, asset := range assets {
            name := strings.ToLower(asset.Name)
            if strings.Contains(name, req.Platform) &&
                strings.Contains(name, req.Arch) &&
                !strings.Contains(name, "-setup.") {
                return i
            }
        }
        return -1
    },
})
```

Пользовательское средство сопоставления полностью заменяет `DefaultAssetMatcher`, поэтому оно отвечает за исключение подписей, контрольных сумм, установщиков и любых других ресурсов, которые не следует устанавливать как обновление приложения.
