---
title: "Pola Arsitektur"
description: "Pola desain untuk aplikasi Wails"
slug: "guides/architecture"
sourcePath: "guides/architecture.md"
---

## Ringkasan

Pola teruji untuk menata aplikasi Wails Anda.

## Pola Lapisan Layanan

### Struktur

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

### Implementasi

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

## Pola Repositori

### Struktur

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

## Arsitektur Berbasis Peristiwa

### Bus Peristiwa

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

### Penggunaan

```go
// Subscribe
eventBus.Subscribe("user.created", func(data interface{}) {
    user := data.(*User)
    sendWelcomeEmail(user)
})

// Publish
eventBus.Publish("user.created", user)
```

## Injeksi Dependensi

### Injeksi Dependensi Manual

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

### Menggunakan Wire

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

## Manajemen State

### State Terpusat

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

## Praktik Terbaik

### ✅ Lakukan

- Pisahkan tanggung jawab
- Gunakan antarmuka
- Injeksikan dependensi
- Tangani error dengan benar
- Pastikan layanan tetap terfokus
- Dokumentasikan arsitektur

### ❌ Jangan Lakukan

- Jangan buat objek serba bisa
- Jangan hubungkan komponen secara terlalu erat
- Jangan abaikan penanganan error
- Jangan abaikan konkurensi
- Jangan merancang secara berlebihan

## Langkah Berikutnya

- [Keamanan](/guides/security/) - Praktik terbaik keamanan
- [Praktik Terbaik](/features/bindings/best-practices/) - Praktik terbaik binding
