---
title: "서비스"
description: "서비스를 사용하여 모듈식의 재사용 가능한 애플리케이션 컴포넌트 구축"
slug: "features/bindings/services"
sourcePath: "features/bindings/services.md"
---

## 서비스 아키텍처

Wails <strong>서비스</strong>는 모듈식의 독립적인 컴포넌트로 애플리케이션 로직을 체계적으로 구성할 수 있게 해 줍니다. 서비스는 시작 및 종료 훅을 통해 수명 주기를 인식하고, 프런트엔드에 자동으로 바인딩되며, 종속성을 주입할 수 있고, 완전히 격리된 상태로 테스트할 수 있습니다.

## 빠른 시작

```go
type GreetService struct {
    prefix string
}

func NewGreetService(prefix string) *GreetService {
    return &GreetService{prefix: prefix}
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

// Register
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewGreetService("Hello, ")),
    },
})
```

**이것으로 끝입니다!** 이제 JavaScript에서 `Greet`을(를) 호출할 수 있습니다.

```js
import { Greet } from './bindings/changeme/GreetService';

const message = await Greet("World");
console.log(message);  // "Hello, World!"
```

## 서비스 만들기

### 기본 서비스

```go
type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}
```

**등록:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**핵심 사항:**

- **내보낸 메서드**(PascalCase)만 바인딩됩니다.
- 서비스는 <strong>싱글턴</strong>입니다(인스턴스 하나).
- 메서드는 `(value, error)`을(를) 반환할 수 있습니다.

### 상태가 있는 서비스

```go
type CounterService struct {
    count int
    mu    sync.RWMutex
}

func (c *CounterService) Increment() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
    return c.count
}

func (c *CounterService) GetCount() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}

func (c *CounterService) Reset() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count = 0
}
```

**중요:** 서비스는 모든 창에서 공유됩니다. 스레드 안전성을 위해 항상 뮤텍스를 사용하세요.

### 종속성이 있는 서비스

```go
type UserService struct {
    db     *sql.DB
    logger *slog.Logger
}

func NewUserService(db *sql.DB, logger *slog.Logger) *UserService {
    return &UserService{
        db:     db,
        logger: logger,
    }
}

func (u *UserService) GetUser(id int) (*User, error) {
    u.logger.Info("Getting user", "id", id)
    
    var user User
    err := u.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}
```

**종속성과 함께 등록:**

```go
db, _ := sql.Open("sqlite3", "app.db")
logger := slog.Default()

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewUserService(db, logger)),
    },
})
```

## 서비스 수명 주기

### ServiceStartup

애플리케이션이 시작될 때 호출됩니다.

```go
func (u *UserService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    u.logger.Info("UserService starting up")
    
    // Initialise resources
    if err := u.db.Ping(); err != nil {
        return fmt.Errorf("database not available: %w", err)
    }
    
    // Run migrations
    if err := u.runMigrations(); err != nil {
        return fmt.Errorf("migrations failed: %w", err)
    }
    
    // Start background tasks
    go u.backgroundSync(ctx)
    
    return nil
}
```

**사용 사례:**

- 리소스 초기화
- 구성 검증
- 마이그레이션 실행
- 백그라운드 작업 시작
- 외부 서비스에 연결

**중요:**

- 서비스는 <strong>등록 순서</strong>대로 시작됩니다.
- <strong>애플리케이션 시작을 중단</strong>하려면 오류를 반환하세요.
- 취소에는 `ctx`을(를) 사용하세요.

### ServiceShutdown

애플리케이션이 종료될 때 호출됩니다.

```go
func (u *UserService) ServiceShutdown() error {
    u.logger.Info("UserService shutting down")
    
    // Close database
    if err := u.db.Close(); err != nil {
        return fmt.Errorf("failed to close database: %w", err)
    }
    
    // Cleanup resources
    u.cleanup()
    
    return nil
}
```

**사용 사례:**

- 연결 닫기
- 상태 저장
- 리소스 정리
- 버퍼 비우기
- 백그라운드 작업 취소

**중요:**

- 서비스는 <strong>역순</strong>으로 종료됩니다.
- 애플리케이션 컨텍스트는 이미 취소된 상태입니다.
- <strong>경고를 기록</strong>하려면 오류를 반환하세요(종료를 막지는 않습니다).

### 전체 수명 주기 예제

```go
type DatabaseService struct {
    db     *sql.DB
    logger *slog.Logger
    cancel context.CancelFunc
}

func NewDatabaseService(logger *slog.Logger) *DatabaseService {
    return &DatabaseService{logger: logger}
}

func (d *DatabaseService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    d.logger.Info("Starting database service")
    
    // Open database
    db, err := sql.Open("sqlite3", "app.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    d.db = db
    
    // Test connection
    if err := db.Ping(); err != nil {
        return fmt.Errorf("database not available: %w", err)
    }
    
    // Start background cleanup
    ctx, cancel := context.WithCancel(ctx)
    d.cancel = cancel
    go d.periodicCleanup(ctx)
    
    return nil
}

func (d *DatabaseService) ServiceShutdown() error {
    d.logger.Info("Shutting down database service")
    
    // Cancel background tasks
    if d.cancel != nil {
        d.cancel()
    }
    
    // Close database
    if d.db != nil {
        if err := d.db.Close(); err != nil {
            return fmt.Errorf("failed to close database: %w", err)
        }
    }
    
    return nil
}

func (d *DatabaseService) periodicCleanup(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            d.cleanup()
        }
    }
}
```

## 서비스 옵션

### 사용자 지정 이름

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&MyService{}, application.ServiceOptions{
            Name: "CustomServiceName",
        }),
    },
})
```

**사용 사례:**

- 동일한 서비스 유형의 여러 인스턴스
- 더 명확한 로깅
- 더 나은 디버깅

### HTTP 경로

서비스에서 HTTP 요청을 처리할 수 있습니다.

```go
type FileService struct {
    root string
}

func (f *FileService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Serve files from root directory
    http.FileServer(http.Dir(f.root)).ServeHTTP(w, r)
}

// Register with route
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&FileService{root: "./files"}, application.ServiceOptions{
            Route: "/files",
        }),
    },
})
```

**접속:** `http://wails.localhost/files/...`

**사용 사례:**

- 파일 제공
- 사용자 지정 API
- WebSocket 엔드포인트
- 미디어 스트리밍

## 서비스 패턴

### 리포지토리 패턴

```go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) GetByID(id int) (*User, error) {
    // Database query
}

func (r *UserRepository) Create(user *User) error {
    // Insert user
}

func (r *UserRepository) Update(user *User) error {
    // Update user
}

func (r *UserRepository) Delete(id int) error {
    // Delete user
}
```

### 서비스 계층 패턴

```go
type UserService struct {
    repo   *UserRepository
    logger *slog.Logger
}

func (s *UserService) RegisterUser(email, password string) (*User, error) {
    // Validate
    if !isValidEmail(email) {
        return nil, errors.New("invalid email")
    }
    
    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &User{
        Email:        email,
        PasswordHash: string(hash),
        CreatedAt:    time.Now(),
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, err
    }
    
    s.logger.Info("User registered", "email", email)
    return user, nil
}
```

### 팩토리 패턴

```go
type ServiceFactory struct {
    db     *sql.DB
    logger *slog.Logger
}

func (f *ServiceFactory) CreateUserService() *UserService {
    return &UserService{
        repo:   &UserRepository{db: f.db},
        logger: f.logger,
    }
}

func (f *ServiceFactory) CreateOrderService() *OrderService {
    return &OrderService{
        repo:   &OrderRepository{db: f.db},
        logger: f.logger,
    }
}
```

### 이벤트 기반 패턴

```go
type OrderService struct {
    app *application.App
}

func (o *OrderService) CreateOrder(items []Item) (*Order, error) {
    order := &Order{
        Items:     items,
        CreatedAt: time.Now(),
    }
    
    // Save order
    if err := o.saveOrder(order); err != nil {
        return nil, err
    }
    
    // Emit event
    o.app.Event.Emit("order-created", order)
    
    return order, nil
}
```

## 종속성 주입

### 생성자 주입

```go
type EmailService struct {
    smtp   *smtp.Client
    logger *slog.Logger
}

func NewEmailService(smtp *smtp.Client, logger *slog.Logger) *EmailService {
    return &EmailService{
        smtp:   smtp,
        logger: logger,
    }
}

// Register
smtpClient := createSMTPClient()
logger := slog.Default()

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewEmailService(smtpClient, logger)),
    },
})
```

### 애플리케이션 주입

```go
type NotificationService struct {
    app *application.App
}

func NewNotificationService(app *application.App) *NotificationService {
    return &NotificationService{app: app}
}

func (n *NotificationService) Notify(message string) {
    // Use application to emit events
    n.app.Event.Emit("notification", message)

    // For native OS notifications, register the notifications service from
    // pkg/services/notifications and call notifier.SendNotification(...).
}

// Register after app creation
app := application.New(application.Options{})
app.RegisterService(application.NewService(NewNotificationService(app)))
```

### 서비스 간 종속성

```go
type OrderService struct {
    userService  *UserService
    emailService *EmailService
}

func NewOrderService(userService *UserService, emailService *EmailService) *OrderService {
    return &OrderService{
        userService:  userService,
        emailService: emailService,
    }
}

// Register in order
userService := &UserService{}
emailService := &EmailService{}
orderService := NewOrderService(userService, emailService)

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(userService),
        application.NewService(emailService),
        application.NewService(orderService),
    },
})
```

## 서비스 테스트

### 단위 테스트

```go
func TestCalculatorService_Add(t *testing.T) {
    calc := &CalculatorService{}
    
    result := calc.Add(2, 3)
    
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### 종속성을 사용한 테스트

```go
func TestUserService_GetUser(t *testing.T) {
    // Create mock database
    db, mock, _ := sqlmock.New()
    defer db.Close()
    
    // Set expectations
    rows := sqlmock.NewRows([]string{"id", "name"}).
        AddRow(1, "Alice")
    mock.ExpectQuery("SELECT").WillReturnRows(rows)
    
    // Create service
    service := NewUserService(db, slog.Default())
    
    // Test
    user, err := service.GetUser(1)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
}
```

### 수명 주기 테스트

```go
func TestDatabaseService_Lifecycle(t *testing.T) {
    service := NewDatabaseService(slog.Default())
    
    // Test startup
    ctx := context.Background()
    err := service.ServiceStartup(ctx, application.ServiceOptions{})
    if err != nil {
        t.Fatalf("startup failed: %v", err)
    }
    
    // Test functionality
    // ...
    
    // Test shutdown
    err = service.ServiceShutdown()
    if err != nil {
        t.Fatalf("shutdown failed: %v", err)
    }
}
```

## 모범 사례

### ✅ 권장 사항

- **단일 책임** - 하나의 서비스는 하나의 목적만 담당하도록 하세요
- **생성자 주입** - 의존성을 명시적으로 전달하세요
- **스레드 안전 상태** - 뮤텍스를 사용하세요
- **오류 반환** - 패닉을 발생시키지 마세요
- **중요한 이벤트 로깅** - 구조화된 로깅을 사용하세요
- **격리된 환경에서 테스트** - 의존성을 모의 객체로 대체하세요

### ❌ 피해야 할 사항

- **전역 상태를 사용하지 마세요** - 의존성을 전달하세요
- **시작을 지연시키지 마세요** - ServiceStartup을 빠르게 완료하세요
- **종료 처리를 무시하지 마세요** - 항상 정리 작업을 수행하세요
- **순환 의존성을 만들지 마세요** - 신중하게 설계하세요
- **내부 메서드를 노출하지 마세요** - 비공개로 유지하세요
- **스레드 안전성을 잊지 마세요** - 서비스는 공유됩니다

## 전체 예제

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log/slog"
    "sync"
    "time"
    
    "github.com/wailsapp/wails/v3/pkg/application"
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}

type UserService struct {
    db     *sql.DB
    logger *slog.Logger
    cache  map[int]*User
    mu     sync.RWMutex
}

func NewUserService(logger *slog.Logger) *UserService {
    return &UserService{
        logger: logger,
        cache:  make(map[int]*User),
    }
}

func (u *UserService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
    u.logger.Info("Starting UserService")
    
    // Open database
    db, err := sql.Open("sqlite3", "users.db")
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    u.db = db
    
    // Create table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return fmt.Errorf("failed to create table: %w", err)
    }
    
    // Preload cache
    if err := u.loadCache(); err != nil {
        return fmt.Errorf("failed to load cache: %w", err)
    }
    
    return nil
}

func (u *UserService) ServiceShutdown() error {
    u.logger.Info("Shutting down UserService")
    
    if u.db != nil {
        return u.db.Close()
    }
    
    return nil
}

func (u *UserService) GetUser(id int) (*User, error) {
    // Check cache first
    u.mu.RLock()
    if user, ok := u.cache[id]; ok {
        u.mu.RUnlock()
        return user, nil
    }
    u.mu.RUnlock()
    
    // Query database
    var user User
    err := u.db.QueryRow(
        "SELECT id, name, email, created_at FROM users WHERE id = ?",
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user %d not found", id)
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    // Update cache
    u.mu.Lock()
    u.cache[id] = &user
    u.mu.Unlock()
    
    return &user, nil
}

func (u *UserService) CreateUser(name, email string) (*User, error) {
    result, err := u.db.Exec(
        "INSERT INTO users (name, email) VALUES (?, ?)",
        name, email,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    id, _ := result.LastInsertId()
    
    user := &User{
        ID:        int(id),
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
    
    // Update cache
    u.mu.Lock()
    u.cache[int(id)] = user
    u.mu.Unlock()
    
    u.logger.Info("User created", "id", id, "email", email)
    
    return user, nil
}

func (u *UserService) loadCache() error {
    rows, err := u.db.Query("SELECT id, name, email, created_at FROM users")
    if err != nil {
        return err
    }
    defer rows.Close()
    
    u.mu.Lock()
    defer u.mu.Unlock()
    
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
            return err
        }
        u.cache[user.ID] = &user
    }
    
    return rows.Err()
}

func main() {
    app := application.New(application.Options{
        Name: "User Management",
        Services: []application.Service{
            application.NewService(NewUserService(slog.Default())),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

## 다음 단계

- [메서드 바인딩](/features/bindings/methods/) - Go 메서드를 JavaScript에 바인딩하는 방법을 알아보세요
- [모델](/features/bindings/models/) - 복잡한 데이터 구조를 바인딩하세요
- [이벤트](/features/events/system/) - 게시/구독 통신에 이벤트를 사용하세요
- [모범 사례](/features/bindings/best-practices/) - 서비스 설계 패턴 및 모범 사례

---

**질문이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [서비스 예제](https://github.com/wailsapp/wails/tree/master/v3/examples)를 확인하세요.
