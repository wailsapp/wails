---
title: "服务"
description: "使用服务构建模块化、可复用的应用程序组件"
slug: "features/bindings/services"
sourcePath: "features/bindings/services.md"
---

## 服务架构

Wails **服务**提供了一种结构化方式，通过模块化且自包含的组件来组织应用程序逻辑。服务可感知生命周期，提供启动和关闭钩子，会自动绑定到前端，支持依赖注入，并且可以完全隔离测试。

## 快速入门

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

<strong>就这么简单！</strong>现在可以从 JavaScript 调用`Greet`：

```js
import { Greet } from './bindings/changeme/GreetService';

const message = await Greet("World");
console.log(message);  // "Hello, World!"
```

## 创建服务

### 基本服务

```go
type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}
```

**注册：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**要点：**

- 仅绑定**导出的方法**（PascalCase）
- 服务是**单例**（只有一个实例）
- 方法可以返回`(value, error)`

### 带状态的服务

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

<strong>重要：</strong>所有窗口共享服务。务必使用互斥锁来保证线程安全。

### 带依赖项的服务

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

**注册带依赖项的服务：**

```go
db, _ := sql.Open("sqlite3", "app.db")
logger := slog.Default()

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewUserService(db, logger)),
    },
})
```

## 服务生命周期

### ServiceStartup

在应用程序启动时调用：

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

**适用场景：**

- 初始化资源
- 验证配置
- 运行迁移
- 启动后台任务
- 连接外部服务

**重要：**

- 服务按**注册顺序**启动
- 返回错误可**阻止应用程序启动**
- 使用`ctx`执行取消操作

### ServiceShutdown

在应用程序关闭时调用：

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

**适用场景：**

- 关闭连接
- 保存状态
- 清理资源
- 刷新缓冲区
- 取消后台任务

**重要：**

- 服务按**相反顺序**关闭
- 应用程序上下文已取消
- 返回错误可**记录警告**（不会阻止关闭）

### 完整生命周期示例

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

## 服务选项

### 自定义名称

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&MyService{}, application.ServiceOptions{
            Name: "CustomServiceName",
        }),
    },
})
```

**适用场景：**

- 日志更清晰
- 更易于调试

`ServiceOptions.Name`是用于日志记录和调试的显示名称；它不会更改生成的绑定命名空间。每种 Go 类型最多注册一个服务。对于独立配置的实例，请使用不同的包装类型，例如：

```go
type PrimaryStore struct { *Store }
type ArchiveStore struct { *Store }
```

注册`application.NewService(&PrimaryStore{primary})`和`application.NewService(&ArchiveStore{archive})`，而不是注册两个`*Store`值。

### HTTP 路由

服务可以处理 HTTP 请求：

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

**访问地址：**`http://wails.localhost/files/...`

**适用场景：**

- 文件服务
- 自定义 API
- WebSocket 端点
- 媒体流传输

## 服务模式

### 仓储模式

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

### 服务层模式

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

### 工厂模式

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

### 事件驱动模式

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

## 依赖注入

### 构造函数注入

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

### 应用程序注入

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

### 服务间依赖

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

## 测试服务

### 单元测试

```go
func TestCalculatorService_Add(t *testing.T) {
    calc := &CalculatorService{}
    
    result := calc.Add(2, 3)
    
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### 带依赖项的测试

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

### 测试生命周期

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

## 最佳实践

### ✅ 应该做

- **单一职责** - 一个服务只负责一项功能
- **构造函数注入** - 显式传入依赖项
- **线程安全的状态** - 使用互斥锁
- **返回错误** - 不要触发 panic
- **记录重要事件** - 使用结构化日志
- **隔离测试** - 模拟依赖项

### ❌ 不应该做

- **不要使用全局状态** - 传入依赖项
- **不要阻塞启动过程** - 确保 ServiceStartup 快速完成
- **不要忽略关闭过程** - 始终执行清理
- **不要创建循环依赖** - 谨慎设计
- **不要暴露内部方法** - 将其保持为私有方法
- **不要忘记线程安全** - 服务是共享的

## 完整示例

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

## 后续步骤

- [方法绑定](/features/bindings/methods/) - 了解如何将 Go 方法绑定到 JavaScript
- [模型](/features/bindings/models/) - 绑定复杂数据结构
- [事件](/features/events/system/) - 使用事件进行发布/订阅通信
- [最佳实践](/features/bindings/best-practices/) - 服务设计模式和最佳实践

---

<strong>有疑问？</strong>请前往[Discord](https://discord.gg/JDdSxwjhGf)提问，或查看[服务示例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples)。
