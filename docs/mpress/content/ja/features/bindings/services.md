---
title: "サービス"
description: "サービスを使用して、モジュール化された再利用可能なアプリケーションコンポーネントを構築します"
slug: "features/bindings/services"
sourcePath: "features/bindings/services.md"
---

## サービスアーキテクチャ

Wailsの<strong>サービス</strong>を使用すると、モジュール化された自己完結型のコンポーネントでアプリケーションロジックを体系的に整理できます。サービスはライフサイクルに対応し、起動時と終了時のフックを備えています。また、フロントエンドに自動的にバインドされ、依存性を注入でき、単独で完全にテストできます。

## クイックスタート

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

**これだけです！** これで、JavaScriptから`Greet`を呼び出せます。

```js
import { Greet } from './bindings/changeme/GreetService';

const message = await Greet("World");
console.log(message);  // "Hello, World!"
```

## サービスの作成

### 基本的なサービス

```go
type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}
```

**登録：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**要点：**

- バインドされるのは、**エクスポートされたメソッド**（PascalCase）のみです
- サービスは<strong>シングルトン</strong>です（インスタンスは1つ）
- メソッドは`(value, error)`を返せます

### 状態を持つサービス

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

**重要：** サービスはすべてのウィンドウで共有されます。スレッドセーフにするため、必ずミューテックスを使用してください。

### 依存関係を持つサービス

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

**依存関係とともに登録：**

```go
db, _ := sql.Open("sqlite3", "app.db")
logger := slog.Default()

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewUserService(db, logger)),
    },
})
```

## サービスのライフサイクル

### ServiceStartup

アプリケーションの起動時に呼び出されます：

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

**ユースケース：**

- リソースを初期化する
- 設定を検証する
- マイグレーションを実行する
- バックグラウンドタスクを開始する
- 外部サービスに接続する

**重要：**

- サービスは<strong>登録順</strong>に起動します
- エラーを返すと、<strong>アプリケーションの起動を中止</strong>できます
- キャンセルには`ctx`を使用します

### ServiceShutdown

アプリケーションの終了時に呼び出されます：

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

**ユースケース：**

- 接続を閉じる
- 状態を保存する
- リソースをクリーンアップする
- バッファをフラッシュする
- バックグラウンドタスクをキャンセルする

**重要：**

- サービスは<strong>逆順</strong>に終了します
- アプリケーションコンテキストはすでにキャンセルされています
- エラーを返すと<strong>警告がログに記録</strong>されます（終了処理は中止されません）

### ライフサイクルの完全な例

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

## サービスオプション

### カスタム名

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&MyService{}, application.ServiceOptions{
            Name: "CustomServiceName",
        }),
    },
})
```

**ユースケース：**

- 同じサービスタイプの複数のインスタンス
- より明確なログ記録
- デバッグの改善

### HTTPルート

サービスでHTTPリクエストを処理できます：

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

**アクセス：** `http://wails.localhost/files/...`

**ユースケース：**

- ファイル配信
- カスタムAPI
- WebSocketエンドポイント
- メディアストリーミング

## サービスパターン

### リポジトリパターン

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

### サービスレイヤーパターン

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

### ファクトリーパターン

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

### イベント駆動パターン

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

## 依存性注入

### コンストラクタインジェクション

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

### アプリケーションインジェクション

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

### サービス間の依存関係

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

## サービスのテスト

### 単体テスト

```go
func TestCalculatorService_Add(t *testing.T) {
    calc := &CalculatorService{}
    
    result := calc.Add(2, 3)
    
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### 依存関係を含むテスト

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

### ライフサイクルのテスト

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

## ベストプラクティス

### ✅ 推奨事項

- **単一責任** - 1つのサービスに1つの目的を持たせる
- **コンストラクターインジェクション** - 依存関係を明示的に渡す
- **スレッドセーフな状態** - ミューテックスを使用する
- **エラーを返す** - パニックを発生させない
- **重要なイベントをログに記録する** - 構造化ログを使用する
- **単独でテストする** - 依存関係をモック化する

### ❌ 非推奨事項

- **グローバル状態を使用しない** - 依存関係を渡す
- **起動処理をブロックしない** - ServiceStartupを短時間で完了させる
- **シャットダウン処理を無視しない** - 必ずクリーンアップする
- **循環依存を作らない** - 慎重に設計する
- **内部メソッドを公開しない** - 非公開のままにする
- **スレッドセーフ性を忘れない** - サービスは共有される

## 完全な例

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

## 次のステップ

- [メソッドバインディング](/features/bindings/methods/) - GoのメソッドをJavaScriptにバインドする方法を学ぶ
- [モデル](/features/bindings/models/) - 複雑なデータ構造をバインドする
- [イベント](/features/events/system/) - イベントを使用してPub/Sub通信を行う
- [ベストプラクティス](/features/bindings/best-practices/) - サービスの設計パターンとベストプラクティス

---

**ご質問がありますか？** [Discord](https://discord.gg/JDdSxwjhGf)で質問するか、[サービスの例](https://github.com/wailsapp/wails/tree/master/v3/examples)を確認してください。
