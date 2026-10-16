---
title: "高度なバインディング"
description: "ディレクティブ、コード注入、カスタム ID などの高度なバインディング手法"
slug: "features/bindings/advanced"
sourcePath: "features/bindings/advanced.md"
---

このガイドでは、Wails v3 のバインディング生成プロセスをカスタマイズし、最適化するための高度な手法について説明します。

## ディレクティブによる生成コードのカスタマイズ

### カスタムコードの注入

`//wails:inject` ディレクティブを使用すると、生成されるバインディングにカスタム JavaScript/TypeScript コードを注入できます。

```go
//wails:inject console.log("Hello from Wails!");
type MyService struct {}

func (s *MyService) Greet(name string) string {
    return "Hello, " + name
}
```

これにより、`MyService` サービス用に生成される JavaScript/TypeScript ファイルに、指定したコードが注入されます。

条件付き注入を使用して、特定の出力形式を対象にすることもできます。

```go
//wails:inject j*:console.log("Hello JS!");  // JavaScript only
//wails:inject t*:console.log("Hello TS!");  // TypeScript only
```

### 追加ファイルの組み込み

`//wails:include` ディレクティブを使用すると、生成されるバインディングに追加ファイルを組み込めます。

```go
//wails:include js/*.js
package mypackage
```

このディレクティブは通常、パッケージのドキュメントコメントで使用し、生成されるバインディングに追加の JavaScript/TypeScript ファイルを組み込みます。

### 型とメソッドを内部用としてマークする

`//wails:internal` ディレクティブは型またはメソッドを内部用としてマークし、フロントエンドへのエクスポートを防ぎます。

```go
//wails:internal
type InternalModel struct {
    Field string
}

//wails:internal
func (s *MyService) InternalMethod() {}
```

これは、Go コード内でのみ使用し、フロントエンドに公開すべきではない型やメソッドに便利です。

### メソッドの無視

`//wails:ignore` ディレクティブを使用すると、バインディング生成時にメソッドが完全に無視されます。

```go
//wails:ignore
func (s *MyService) IgnoredMethod() {}
```

これは `//wails:internal` と似ていますが、メソッドを内部用としてマークするのではなく、完全に無視します。

### カスタムメソッド ID

`//wails:id` ディレクティブはメソッドのカスタム ID を指定し、デフォルトのハッシュベース ID を上書きします。

```go
//wails:id 42
func (s *MyService) CustomIDMethod() {}
```

これは、コードをリファクタリングする際に互換性を維持するのに役立ちます。

## 複雑な型の扱い

### ネストされた構造体

バインディングジェネレーターは、ネストされた構造体を自動的に処理します。

```go
type Address struct {
    Street string
    City   string
    State  string
    Zip    string
}

type Person struct {
    Name    string
    Address Address
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Address: Address{
            Street: "123 Main St",
            City:   "Anytown",
            State:  "CA",
            Zip:    "12345",
        },
    }
}
```

生成される JavaScript/TypeScript コードには、`Person` と `Address` の両方のクラスが含まれます。

### マップとスライス

マップとスライスも自動的に処理されます。

```go
type Person struct {
    Name       string
    Attributes map[string]string
    Friends    []string
}

func (s *MyService) GetPerson() Person {
    return Person{
        Name: "John Doe",
        Attributes: map[string]string{
            "hair": "brown",
            "eyes": "blue",
        },
        Friends: []string{"Jane", "Bob", "Alice"},
    }
}
```

JavaScript では、マップはオブジェクト、スライスは配列として表現されます。TypeScript では、マップは `Record<K, V>`、スライスは `T[]` として表現されます。

### ジェネリック型

バインディングジェネレーターはジェネリック型をサポートしています。

```go
type Result[T any] struct {
    Data  T
    Error string
}

func (s *MyService) GetResult() Result[string] {
    return Result[string]{
        Data:  "Hello, World!",
        Error: "",
    }
}
```

生成される TypeScript コードには、`Result` のジェネリッククラスが含まれます。

```typescript
export class Result<T> {
    "Data": T;
    "Error": string;

    constructor(source: Partial<Result<T>> = {}) {
        if (!("Data" in source)) {
            this["Data"] = null as any;
        }
        if (!("Error" in source)) {
            this["Error"] = "";
        }

        Object.assign(this, source);
    }

    static createFrom<T>(source: string | object = {}): Result<T> {
        let parsedSource = typeof source === "string" ? JSON.parse(source) : source;
        return new Result<T>(parsedSource as Partial<Result<T>>);
    }
}
```

### インターフェース

`-i` フラグを使用すると、バインディングジェネレーターはクラスの代わりに TypeScript インターフェースを生成できます。

```bash
wails3 generate bindings -ts -i
```

これにより、すべてのモデルに対して TypeScript インターフェースが生成されます。

```typescript
export interface Person {
    Name: string;
    Attributes: Record<string, string>;
    Friends: string[];
}
```

## バインディング生成の最適化

### ID の代わりに名前を使用する

デフォルトでは、バインディングジェネレーターはメソッド呼び出しにハッシュベースの ID を使用します。代わりに名前を使用するには、`-names` フラグを使用できます。

```bash
wails3 generate bindings -names
```

これにより、バインドされたメソッドを<strong>完全修飾</strong>名（`<package>.<Service>.<Method>`）で呼び出し、位置に基づく `$0`、`$1`、… というパラメーター名を使用するコードが生成されます。

```javascript
export function Greet($0) {
    let $resultPromise = $Call.ByName("main.GreetService.Greet", $0);
    return $resultPromise;
}
```

これにより生成コードの可読性が高まり、デバッグしやすくなりますが、効率がわずかに低下する可能性があります。

### ランタイムのバンドル

デフォルトでは、生成コードは `@wailsio/runtime` npm パッケージから Wails ランタイムをインポートします。`-b` フラグを使用すると、生成コードにランタイムをバンドルできます。

```bash
wails3 generate bindings -b
```

これにより、生成されるファイルにランタイムコードが直接組み込まれ、npm パッケージが不要になります。

### インデックスファイルの無効化

インデックスファイルが不要な場合は、`-noindex` フラグを使用して生成を無効にできます。

```bash
wails3 generate bindings -noindex
```

これは、サービスとモデルをそれぞれのファイルから直接インポートしたい場合に便利です。

## 実践的な例

### 認証サービス

カスタムディレクティブを使用した認証サービスの例を次に示します。

```go
package auth

//wails:inject console.log("Auth service initialized");
type AuthService struct {
    // Private fields
    users map[string]User
}

type User struct {
    Username string
    Email    string
    Role     string
}

type LoginRequest struct {
    Username string
    Password string
}

type LoginResponse struct {
    Success bool
    User    User
    Token   string
    Error   string
}

// Login authenticates a user
func (s *AuthService) Login(req LoginRequest) LoginResponse {
    // Implementation...
}

// GetCurrentUser returns the current user
func (s *AuthService) GetCurrentUser() User {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *AuthService) validateCredentials(username, password string) bool {
    // Implementation...
}
```

### データ処理サービス

ジェネリック型を使用したデータ処理サービスの例を次に示します。

```go
package data

type ProcessingResult[T any] struct {
    Data  T
    Error string
}

type DataService struct {}

// Process processes data and returns a result
func (s *DataService) Process(data string) ProcessingResult[map[string]int] {
    // Implementation...
}

// ProcessBatch processes multiple data items
func (s *DataService) ProcessBatch(data []string) ProcessingResult[[]map[string]int] {
    // Implementation...
}

// Internal helper method
//wails:internal
func (s *DataService) parseData(data string) (map[string]int, error) {
    // Implementation...
}
```

### 条件付きコード注入

出力形式ごとに条件付きでコードを注入する例を次に示します。

```go
//wails:inject j*:/**
//wails:inject j*: * @param {string} arg
//wails:inject j*: * @returns {Promise<void>}
//wails:inject j*: */
//wails:inject j*:export async function CustomMethod(arg) {
//wails:inject t*:export async function CustomMethod(arg: string): Promise<void> {
//wails:inject     await InternalMethod("Hello " + arg + "!");
//wails:inject }
type Service struct{}
```

これにより、JavaScript と TypeScript の出力にそれぞれ異なるコードが挿入され、各言語に適した型注釈が付与されます。
