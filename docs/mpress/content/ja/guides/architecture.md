---
title: "アーキテクチャパターン"
description: "Wailsアプリケーションの設計パターン"
slug: "guides/architecture"
sourcePath: "guides/architecture.md"
---

## 概要

Wailsアプリケーションを整理するための実証済みパターン。

## サービスレイヤーパターン

### 構造

```
app/
├── main.go
├── services/
│   ├── user_service.go
│   ├── data_service.go
│   └── file_service.go
└── models/
    └── user.go
```

### 実装

```go
// Service interface
type UserService interface {
    Create(email, password string) (*User, error)
    GetByID(id int) (*User, error)
    Update(user *User) error
    Delete(id int) error
}

// Implementation
type userService struct {
    app *application.App
    db  *sql.DB
}

func NewUserService(app *application.App, db *sql.DB) UserService {
    return &userService{app: app, db: db}
}
```

## リポジトリパターン

### 構造

```go
// Repository interface
type UserRepository interface {
    Create(user *User) error
    FindByID(id int) (*User, error)
    Update(user *User) error
    Delete(id int) error
}

// Service uses repository
type UserService struct {
    repo UserRepository
}

func (s *UserService) Create(email, password string) (*User, error) {
    user := &User{Email: email}
    return user, s.repo.Create(user)
}
```

## イベント駆動アーキテクチャ

### イベントバス

```go
type EventBus struct {
    app       *application.App
    listeners map[string][]func(interface{})
    mu        sync.RWMutex
}

func (eb *EventBus) Subscribe(event string, handler func(interface{})) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    eb.listeners[event] = append(eb.listeners[event], handler)
}

func (eb *EventBus) Publish(event string, data interface{}) {
    eb.mu.RLock()
    handlers := eb.listeners[event]
    eb.mu.RUnlock()
    
    for _, handler := range handlers {
        go handler(data)
    }
}
```

### 使用方法

```go
// Subscribe
eventBus.Subscribe("user.created", func(data interface{}) {
    user := data.(*User)
    sendWelcomeEmail(user)
})

// Publish
eventBus.Publish("user.created", user)
```

## 依存性注入

### 手動による依存性注入

```go
type App struct {
    userService *UserService
    fileService *FileService
    db          *sql.DB
}

func NewApp() *App {
    db := openDatabase()
    
    return &App{
        db:          db,
        userService: NewUserService(db),
        fileService: NewFileService(db),
    }
}
```

### Wireの使用

```go
// wire.go
//go:build wireinject

func InitializeApp() (*App, error) {
    wire.Build(
        openDatabase,
        NewUserService,
        NewFileService,
        NewApp,
    )
    return nil, nil
}
```

## 状態管理

### 一元化された状態

```go
type AppState struct {
    currentUser *User
    settings    *Settings
    mu          sync.RWMutex
}

func (s *AppState) SetUser(user *User) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.currentUser = user
}

func (s *AppState) GetUser() *User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.currentUser
}
```

## ベストプラクティス

### ✅ 推奨事項

- 関心事を分離する
- インターフェースを使用する
- 依存関係を注入する
- エラーを適切に処理する
- 各サービスの責務を限定する
- アーキテクチャを文書化する

### ❌ 禁止事項

- 神オブジェクトを作成しない
- コンポーネント同士を密結合にしない
- エラー処理を省略しない
- 並行処理を無視しない
- 過剰に設計しない

## 次のステップ

- [セキュリティ](/guides/security/) - セキュリティのベストプラクティス
- [ベストプラクティス](/features/bindings/best-practices/) - バインディングのベストプラクティス
