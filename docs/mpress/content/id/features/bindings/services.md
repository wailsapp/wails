---
title: "Layanan"
description: "Bangun komponen aplikasi yang modular dan dapat digunakan kembali dengan layanan"
slug: "features/bindings/services"
sourcePath: "features/bindings/services.md"
---

## Arsitektur Layanan

**Layanan** Wails menyediakan cara terstruktur untuk mengatur logika aplikasi dengan komponen modular dan mandiri. Layanan mendukung siklus hidup melalui hook startup dan shutdown, secara otomatis diikat ke frontend, mendukung injeksi dependensi, dan dapat diuji sepenuhnya secara terpisah.

## Mulai Cepat

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

**Selesai!** `Greet` kini dapat dipanggil dari JavaScript:

```js
import { Greet } from './bindings/changeme/GreetService';

const message = await Greet("World");
console.log(message);  // "Hello, World!"
```

## Membuat Layanan

### Layanan Dasar

```go
type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}
```

**Daftarkan:**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**Poin penting:**

- Hanya **metode yang diekspor** (PascalCase) yang diikat
- Layanan merupakan **singleton** (satu instans)
- Metode dapat mengembalikan `(value, error)`

### Layanan dengan State

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

**Penting:** Layanan digunakan bersama oleh semua jendela. Selalu gunakan mutex untuk keamanan thread.

### Layanan dengan Dependensi

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

**Daftarkan dengan dependensi:**

```go
db, _ := sql.Open("sqlite3", "app.db")
logger := slog.Default()

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewUserService(db, logger)),
    },
})
```

## Siklus Hidup Layanan

### ServiceStartup

Dipanggil saat aplikasi dimulai:

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

**Kasus penggunaan:**

- Inisialisasi sumber daya
- Validasi konfigurasi
- Jalankan migrasi
- Mulai tugas latar belakang
- Hubungkan ke layanan eksternal

**Penting:**

- Layanan dimulai sesuai **urutan pendaftaran**
- Kembalikan error untuk **mencegah aplikasi dimulai**
- Gunakan `ctx` untuk pembatalan

### ServiceShutdown

Dipanggil saat aplikasi dimatikan:

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

**Kasus penggunaan:**

- Tutup koneksi
- Simpan state
- Bersihkan sumber daya
- Tulis data yang masih tersimpan dalam buffer ke tujuannya
- Batalkan tugas latar belakang

**Penting:**

- Layanan dimatikan dalam **urutan terbalik**
- Konteks aplikasi sudah dibatalkan
- Kembalikan error untuk **mencatat peringatan** (tidak mencegah proses shutdown)

### Contoh Siklus Hidup Lengkap

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

## Opsi Layanan

### Nama Khusus

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&MyService{}, application.ServiceOptions{
            Name: "CustomServiceName",
        }),
    },
})
```

**Kasus penggunaan:**

- Beberapa instans dari tipe layanan yang sama
- Pencatatan log yang lebih jelas
- Debugging yang lebih mudah

### Rute HTTP

Layanan dapat menangani permintaan HTTP:

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

**Akses:** `http://wails.localhost/files/...`

**Kasus penggunaan:**

- Penyajian file
- API khusus
- Endpoint WebSocket
- Streaming media

## Pola Layanan

### Pola Repository

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

### Pola Lapisan Layanan

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

### Pola Factory

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

### Pola Berbasis Peristiwa

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

## Injeksi Dependensi

### Injeksi Konstruktor

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

### Injeksi Aplikasi

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

### Dependensi Antarlayanan

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

## Menguji Layanan

### Pengujian Unit

```go
func TestCalculatorService_Add(t *testing.T) {
    calc := &CalculatorService{}
    
    result := calc.Add(2, 3)
    
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### Pengujian dengan Dependensi

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

### Pengujian Siklus Hidup

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

## Praktik Terbaik

### ✅ Lakukan

- **Tanggung jawab tunggal** - Satu layanan, satu tujuan
- **Injeksi melalui konstruktor** - Teruskan dependensi secara eksplisit
- **State yang aman untuk thread** - Gunakan mutex
- **Kembalikan error** - Jangan memicu panic
- **Catat peristiwa penting** - Gunakan pencatatan terstruktur
- **Uji secara terisolasi** - Gunakan mock untuk dependensi

### ❌ Jangan Lakukan

- **Jangan gunakan state global** - Teruskan dependensi
- **Jangan hambat proses startup** - Pastikan ServiceStartup berjalan cepat
- **Jangan abaikan proses shutdown** - Selalu lakukan pembersihan
- **Jangan buat dependensi melingkar** - Rancang dengan cermat
- **Jangan ekspos metode internal** - Pertahankan sebagai metode privat
- **Jangan abaikan keamanan thread** - Layanan digunakan bersama

## Contoh Lengkap

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

## Langkah Berikutnya

- [Binding Metode](/features/bindings/methods/) - Pelajari cara mengikat metode Go ke JavaScript
- [Model](/features/bindings/models/) - Ikat struktur data yang kompleks
- [Peristiwa](/features/events/system/) - Gunakan peristiwa untuk komunikasi publikasi/langganan
- [Praktik Terbaik](/features/bindings/best-practices/) - Pola desain layanan dan praktik terbaik

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh layanan](https://github.com/wailsapp/wails/tree/master/v3/examples).
