---
title: "Архитектурные шаблоны"
description: "Шаблоны проектирования приложений Wails"
slug: "guides/architecture"
sourcePath: "guides/architecture.md"
---

## Обзор

Проверенные шаблоны организации приложения Wails.

## Шаблон сервисного слоя

### Структура

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

### Реализация

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

## Шаблон «Репозиторий»

### Структура

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

## Событийно-ориентированная архитектура

### Шина событий

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

### Использование

```go
// Subscribe
eventBus.Subscribe("user.created", func(data interface{}) {
    user := data.(*User)
    sendWelcomeEmail(user)
})

// Publish
eventBus.Publish("user.created", user)
```

## Внедрение зависимостей

### Ручное внедрение зависимостей

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

### Использование Wire

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

## Управление состоянием

### Централизованное состояние

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

## Рекомендации

### ✅ Рекомендуется

- Разделяйте обязанности
- Используйте интерфейсы
- Внедряйте зависимости
- Корректно обрабатывайте ошибки
- Не перегружайте сервисы лишними обязанностями
- Документируйте архитектуру

### ❌ Не рекомендуется

- Не создавайте объекты, отвечающие за всё
- Не связывайте компоненты слишком тесно
- Не пренебрегайте обработкой ошибок
- Не игнорируйте конкурентное выполнение
- Не усложняйте архитектуру без необходимости

## Дальнейшие шаги

- [Безопасность](/guides/security/) — рекомендации по безопасности
- [Рекомендации](/features/bindings/best-practices/) — рекомендации по привязкам
