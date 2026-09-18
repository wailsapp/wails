---
title: "아키텍처 패턴"
description: "Wails 애플리케이션을 위한 디자인 패턴"
slug: "guides/architecture"
sourcePath: "guides/architecture.md"
---

## 개요

검증된 패턴을 사용하여 Wails 애플리케이션을 구성합니다.

## 서비스 계층 패턴

### 구조

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

### 구현

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

## 리포지토리 패턴

### 구조

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

## 이벤트 기반 아키텍처

### 이벤트 버스

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

### 사용법

```go
// Subscribe
eventBus.Subscribe("user.created", func(data interface{}) {
    user := data.(*User)
    sendWelcomeEmail(user)
})

// Publish
eventBus.Publish("user.created", user)
```

## 의존성 주입

### 수동 의존성 주입

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

### Wire 사용

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

## 상태 관리

### 중앙 집중식 상태

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

## 모범 사례

### ✅ 권장 사항

- 관심사를 분리하세요
- 인터페이스를 사용하세요
- 의존성을 주입하세요
- 오류를 적절하게 처리하세요
- 서비스가 한 가지 책임에 집중하도록 하세요
- 아키텍처를 문서화하세요

### ❌ 금지 사항

- 갓 오브젝트를 만들지 마세요
- 컴포넌트를 강하게 결합하지 마세요
- 오류 처리를 생략하지 마세요
- 동시성을 무시하지 마세요
- 과도하게 설계하지 마세요

## 다음 단계

- [보안](/guides/security/) - 보안 모범 사례
- [모범 사례](/features/bindings/best-practices/) - 바인딩 모범 사례
