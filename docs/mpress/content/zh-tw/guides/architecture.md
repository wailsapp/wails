---
title: "架構模式"
description: "Wails 應用程式的設計模式"
slug: "guides/architecture"
sourcePath: "guides/architecture.md"
---

## 概觀

經過實務驗證的 Wails 應用程式組織模式。

## 服務層模式

### 結構

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

### 實作

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

## 儲存庫模式

### 結構

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

## 事件驅動架構

### 事件匯流排

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

### 用法

```go
// Subscribe
eventBus.Subscribe("user.created", func(data interface{}) {
    user := data.(*User)
    sendWelcomeEmail(user)
})

// Publish
eventBus.Publish("user.created", user)
```

## 相依性注入

### 手動相依性注入

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

### 使用 Wire

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

## 狀態管理

### 集中式狀態

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

## 最佳實務

### ✅ 應做事項

- 分離關注點
- 使用介面
- 注入相依性
- 妥善處理錯誤
- 讓服務各司其職
- 記錄架構設計

### ❌ 應避免事項

- 不要建立上帝物件
- 不要讓元件緊密耦合
- 不要略過錯誤處理
- 不要忽略並行處理
- 不要過度設計

## 後續步驟

- [安全性](/guides/security/) - 安全性最佳實務
- [最佳實務](/features/bindings/best-practices/) - 繫結最佳實務
