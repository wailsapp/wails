---
title: "服務"
description: "使用服務建構模組化且可重複使用的應用程式元件"
slug: "features/bindings/services"
sourcePath: "features/bindings/services.md"
---

## 服務架構

Wails <strong>服務</strong>提供一種結構化方式，讓您能以模組化且自成一體的元件組織應用程式邏輯。服務能感知生命週期，並提供啟動與關閉掛鉤；服務會自動繫結至前端、支援相依性注入，也可完全獨立測試。

## 快速開始

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

<strong>就是這麼簡單！</strong>現在可以從 JavaScript 呼叫`Greet`：

```js
import { Greet } from './bindings/changeme/GreetService';

const message = await Greet("World");
console.log(message);  // "Hello, World!"
```

## 建立服務

### 基本服務

```go
type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}
```

**註冊：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**重點：**

- 只會繫結<strong>匯出的方法</strong>（PascalCase）
- 服務是<strong>單例</strong>（只有一個執行個體）
- 方法可以傳回`(value, error)`

### 具備狀態的服務

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

<strong>重要：</strong>所有視窗會共用服務。務必使用互斥鎖以確保執行緒安全。

### 具備相依性的服務

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

**連同相依項目一起註冊：**

```go
db, _ := sql.Open("sqlite3", "app.db")
logger := slog.Default()

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewUserService(db, logger)),
    },
})
```

## 服務生命週期

### ServiceStartup

應用程式啟動時呼叫：

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

**使用情境：**

- 初始化資源
- 驗證組態
- 執行移轉
- 啟動背景工作
- 連線至外部服務

**重要：**

- 服務會依照<strong>註冊順序</strong>啟動
- 傳回錯誤可<strong>阻止應用程式啟動</strong>
- 使用`ctx`進行取消

### ServiceShutdown

應用程式關閉時呼叫：

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

**使用情境：**

- 關閉連線
- 儲存狀態
- 清理資源
- 將緩衝區中的資料寫出至目的地
- 取消背景工作

**重要：**

- 服務會依照<strong>相反順序</strong>關閉
- 應用程式的執行上下文已取消
- 傳回錯誤以<strong>記錄警告</strong>（不會阻止關閉）

### 完整的生命週期範例

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

## 服務選項

### 自訂名稱

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&MyService{}, application.ServiceOptions{
            Name: "CustomServiceName",
        }),
    },
})
```

**使用情境：**

- 同一服務型別的多個執行個體
- 更清楚的記錄
- 更容易偵錯

### HTTP 路由

服務可以處理 HTTP 請求：

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

**存取：**`http://wails.localhost/files/...`

**使用情境：**

- 提供檔案
- 自訂 API
- WebSocket 端點
- 媒體串流

## 服務模式

### 儲存庫模式

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

### 服務層模式

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

### 工廠模式

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

### 事件驅動模式

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

## 相依性注入

### 建構函式注入

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

### 應用程式注入

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

### 服務間相依性

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

## 測試服務

### 單元測試

```go
func TestCalculatorService_Add(t *testing.T) {
    calc := &CalculatorService{}
    
    result := calc.Add(2, 3)
    
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### 使用相依項目進行測試

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

### 測試生命週期

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

## 最佳實務

### ✅ 建議做法

- **單一職責**——一個服務只負責一項用途
- **建構函式注入**——明確傳入相依項目
- **執行緒安全的狀態**——使用互斥鎖
- **傳回錯誤**——不要引發 panic
- **記錄重要事件**——使用結構化日誌記錄
- **獨立測試**——模擬相依項目

### ❌ 避免的做法

- **不要使用全域狀態**——傳入相依項目
- **不要阻塞啟動流程**——讓 ServiceStartup 快速完成
- **不要忽略關閉流程**——務必執行清理
- **不要建立循環相依性**——審慎設計
- **不要公開內部方法**——將其保留為私有方法
- **不要忘記執行緒安全**——服務會由多處共用

## 完整範例

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

## 後續步驟

- [方法繫結](/features/bindings/methods/)——瞭解如何將 Go 方法繫結至 JavaScript
- [模型](/features/bindings/models/)——繫結複雜的資料結構
- [事件](/features/events/system/)——使用事件進行發布／訂閱通訊
- [最佳實務](/features/bindings/best-practices/)——服務設計模式與最佳實務

---

<strong>有問題嗎？</strong>請在[Discord](https://discord.gg/JDdSxwjhGf)中提問，或查看[服務範例](https://github.com/wailsapp/wails/tree/master/v3/examples)。
