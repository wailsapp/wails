---
title: "列挙型"
description: "Go 定数からの列挙型の自動生成"
slug: "features/bindings/enums"
sourcePath: "features/bindings/enums.md"
---

## 列挙型のバインディング

Wails v3 のバインディングジェネレーターは、**Go の定数型を自動的に検出し、TypeScript の列挙型または JavaScript の const オブジェクトを生成します**。登録も設定も不要です。Go で型と定数を定義するだけで、残りはジェネレーターが処理します。

@note{type="info"}
Wails v2 とは異なり、`EnumBind`<strong>を呼び出したり、列挙型を手動で登録したりする</strong>必要はありません。ジェネレーターがソースコードから自動的に検出します。

@end

## クイックスタート

**Go で、定数を持つ名前付き型を定義します。**

```go
type Status string

const (
    StatusActive  Status = "active"
    StatusPending Status = "pending"
    StatusClosed  Status = "closed"
)
```

**構造体またはサービスメソッドでその型を使用します。**

```go
type Ticket struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status Status `json:"status"`
}
```

**バインディングを生成します。**

```bash
wails3 generate bindings
```

ジェネレーターの出力には、モデル数とともに列挙型の数が表示されます。

```
3 Enums, 5 Models
```

**フロントエンドで使用します。**

```javascript
import { Ticket, Status } from './bindings/changeme/models'

const ticket = new Ticket({
    id: 1,
    title: "Bug report",
    status: Status.StatusActive
})
```

<strong>以上です。</strong>列挙型は Go と JavaScript/TypeScript の両方で適用されます。

## 列挙型の定義

Wails の列挙型は、基底となる基本型を持つ<strong>名前付き型</strong>と、その型の<strong>const 宣言</strong>を組み合わせたものです。

### 文字列列挙型

```go
// Title is a title
type Title string

const (
    // Mister is a title
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
    Mrs    Title = "Mrs"
    Dr     Title = "Dr"
)
```

**生成される TypeScript：**

```typescript
/**
 * Title is a title
 */
export enum Title {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = "",

    /**
     * Mister is a title
     */
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
    Mrs = "Mrs",
    Dr = "Dr",
}
```

**生成される JavaScript：**

```javascript
/**
 * Title is a title
 * @readonly
 * @enum {string}
 */
export const Title = {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero: "",

    /**
     * Mister is a title
     */
    Mister: "Mr",
    Miss: "Miss",
    Ms: "Ms",
    Mrs: "Mrs",
    Dr: "Dr",
};
```

### 整数列挙型

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityMedium Priority = 1
    PriorityHigh   Priority = 2
)
```

**生成される TypeScript：**

```typescript
export enum Priority {
    /**
     * The Go zero value for the underlying type of the enum.
     */
    $zero = 0,

    PriorityLow = 0,
    PriorityMedium = 1,
    PriorityHigh = 2,
}
```

### 型エイリアスの列挙型

Go の型エイリアス（`=`）も使用できますが、生成される出力は少し異なります。TypeScript ネイティブの`enum`ではなく、型定義と const オブジェクトが生成されます。

```go
// Age is an integer with some predefined values
type Age = int

const (
    NewBorn    Age = 0
    Teenager   Age = 12
    YoungAdult Age = 18

    // Oh no, some grey hair!
    MiddleAged Age = 50
    Mathusalem Age = 1000 // Unbelievable!
)
```

**生成される TypeScript：**

```typescript
/**
 * Age is an integer with some predefined values
 */
export type Age = number;

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

**生成される JavaScript：**

```javascript
/**
 * Age is an integer with some predefined values
 * @typedef {number} Age
 */

/**
 * Predefined constants for type Age.
 * @namespace
 */
export const Age = {
    NewBorn: 0,
    Teenager: 12,
    YoungAdult: 18,

    /**
     * Oh no, some grey hair!
     */
    MiddleAged: 50,

    /**
     * Unbelievable!
     */
    Mathusalem: 1000,
};
```

@note{type="tip"}
**名前付き型**（`type Title string`）からは、`$zero`メンバーを持つ TypeScript ネイティブの`enum`宣言が生成されます。 **型エイリアス**（`type Age = int`）からは、`$zero`を持たない`type`と`const`の名前空間ペアが生成されます。

@end

## `$zero`値

名前付き型の各列挙型には、基底型の<strong>Go ゼロ値</strong>を表す特別な`$zero`メンバーが含まれます。

| 基底型 | `$zero`値 |
| --- | --- |
| `string` | `""` |
| `int`、`int8`、`int16`、`int32`、`int64` | `0` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `0` |
| `float32`、`float64` | `0` |
| `bool` | `false` |

構造体のフィールドで列挙型を使用し、値が指定されていない場合、コンストラクターはデフォルトで`$zero`を使用します。

```typescript
export class Person {
    "Title": Title;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;  // defaults to ""
        }
        Object.assign(this, $$source);
    }
}
```

これにより、クラスを生成するときに型安全な初期化が保証され、列挙型フィールドが`undefined`になることはありません。TypeScript インターフェースを生成する場合（`-i`を使用）はコンストラクターがないため、通常どおりフィールドが存在しないことがあります。

## 構造体での列挙型の使用

構造体のフィールドが列挙型である場合、生成されるコードはプリミティブ型にフォールバックせず、その型を<strong>保持します</strong>。

```go
type Person struct {
    Title Title
    Name  string
    Age   Age
}
```

**生成される TypeScript：**

```typescript
export class Person {
    "Title": Title;
    "Name": string;
    "Age": Age;

    constructor($$source: Partial<Person> = {}) {
        if (!("Title" in $$source)) {
            this["Title"] = Title.$zero;
        }
        if (!("Name" in $$source)) {
            this["Name"] = "";
        }
        if (!("Age" in $$source)) {
            this["Age"] = 0;
        }

        Object.assign(this, $$source);
    }
}
```

`Title`フィールドの型は`string`ではなく`Title`です。これにより、IDE で列挙値に対する完全な自動補完と型チェックを利用できます。

## インポートしたパッケージの列挙型

別のパッケージで定義された列挙型も完全にサポートされています。対応するパッケージディレクトリに生成されます。

```go
// services/types.go
package services

type Title string

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
    Ms     Title = "Ms"
)
```

```go
// main.go
package main

import "myapp/services"

func (*GreetService) Greet(name string, title services.Title) string {
    return "Hello " + string(title) + " " + name
}
```

`Title`列挙型は`services`モデルファイルに生成され、インポートパスは自動的に解決されます。

```typescript
// bindings/changeme/services/models.ts
export enum Title {
    $zero = "",
    Mister = "Mr",
    Miss = "Miss",
    Ms = "Ms",
}
```

## 列挙型のメソッド

Go では列挙型にメソッドを追加できます。これらはバインディングの生成には影響しませんが、サーバー側で便利な機能を提供します。

```go
type Title string

func (t Title) String() string {
    return string(t)
}

const (
    Mister Title = "Mr"
    Miss   Title = "Miss"
)
```

その型に Go メソッドが存在するかどうかにかかわらず、生成される列挙型は同一です。

## コメントとドキュメント

ジェネレーターは、Go のコメントを生成出力内の JSDoc として保持します。

- <strong>型のコメント</strong>は、列挙型のドキュメントコメントになります
- <strong>const グループのコメント</strong>は、セクションの区切りになります
- <strong>個々の const のコメント</strong>は、メンバーのドキュメントコメントになります
- <strong>インラインコメント</strong>は、可能な限り保持されます

これにより、列挙値にカーソルを合わせると、IDE にドキュメントが表示されます。

## サポートされる基底型

バインディングジェネレーターは、次の Go 基底型を持つ列挙型をサポートします。

| Go の型 | 列挙型として使用可能 |
| --- | :---: |
| `string` | はい |
| `int`、`int8`、`int16`、`int32`、`int64` | はい |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | はい |
| `float32`、`float64` | はい |
| `byte`（`uint8`） | はい |
| `rune`（`int32`） | はい |
| `bool` | はい |
| `complex64`、`complex128` | いいえ |

## 制限事項

列挙型の生成では、以下は<strong>サポートされていません</strong>：

- **ジェネリック型** — 型パラメーターがあると定数を検出できません
- **カスタムの `json.Marshaler` または `encoding.TextMarshaler`** を持つ型 — カスタムシリアライズを使用すると、生成された値が実行時の動作と一致しない可能性があるため、ジェネレーターはこれらをスキップします
- **値を静的に評価または表現できない定数** — 定数には、基底型で表現できる既知の値が必要です。標準的な `iota` パターンはコンパイラーによって具体的な値に解決されるため、問題なく機能します
- **複素数型** — `complex64` と `complex128` は列挙型の基底型にはできません

## 完全な例

**Go：**

```go
package main

import (
    "github.com/wailsapp/wails/v3/pkg/application"
)

// BackgroundType defines the type of background
type BackgroundType string

const (
    BackgroundSolid    BackgroundType = "solid"
    BackgroundGradient BackgroundType = "gradient"
    BackgroundImage    BackgroundType = "image"
)

type BackgroundConfig struct {
    Type  BackgroundType `json:"type"`
    Value string         `json:"value"`
}

type ThemeService struct{}

func (*ThemeService) GetBackground() BackgroundConfig {
    return BackgroundConfig{
        Type:  BackgroundSolid,
        Value: "#ffffff",
    }
}

func (*ThemeService) SetBackground(config BackgroundConfig) error {
    // Apply background
    return nil
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&ThemeService{}),
        },
    })
    app.Window.New()
    app.Run()
}
```

**フロントエンド（TypeScript）：**

```typescript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Get current background
const bg = await GetBackground()

// Check the type using enum values
if (bg.type === BackgroundType.BackgroundSolid) {
    console.log("Solid background:", bg.value)
}

// Set a new background
await SetBackground(new BackgroundConfig({
    type: BackgroundType.BackgroundGradient,
    value: "linear-gradient(to right, #000, #fff)"
}))
```

**フロントエンド（JavaScript）：**

```javascript
import { GetBackground, SetBackground } from './bindings/changeme/themeservice'
import { BackgroundConfig, BackgroundType } from './bindings/changeme/models'

// Use enum values for type-safe comparisons
const bg = await GetBackground()

switch (bg.type) {
    case BackgroundType.BackgroundSolid:
        applySolid(bg.value)
        break
    case BackgroundType.BackgroundGradient:
        applyGradient(bg.value)
        break
    case BackgroundType.BackgroundImage:
        applyImage(bg.value)
        break
}
```

## 次のステップ

@cards{cols="2"}
📖 データモデル
構造体、型マッピング、モデル生成について説明します。

[詳しく見る →](/features/bindings/models/)

---
🚀 メソッドバインディング
Go のメソッドをフロントエンドにバインドします。

[詳しく見る →](/features/bindings/methods/)

---
⚙ 高度なバインディング
ディレクティブ、コードインジェクション、カスタム ID について説明します。

[詳しく見る →](/features/bindings/advanced/)

---
✓ ベストプラクティス
バインディングの設計パターンについて説明します。

[詳しく見る →](/features/bindings/best-practices/)

@end

---

**質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf) で質問するか、[バインディングの例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)を確認してください。
