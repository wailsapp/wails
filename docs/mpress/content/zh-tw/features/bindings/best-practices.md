---
title: "繫結最佳實務"
description: "Go-JavaScript 繫結的設計模式與最佳實務"
slug: "features/bindings/best-practices"
sourcePath: "features/bindings/best-practices.md"
---

## 繫結最佳實務

設計繫結時，請遵循<strong>經過驗證的模式</strong>，以建立簡潔、高效能且安全的繫結。本指南涵蓋 API 設計原則、效能最佳化、安全模式、錯誤處理，以及可維護應用程式的測試策略。

## API 設計原則

### 1. 單一職責

每項服務都應有一個明確用途：

```go
// ❌ Bad: God object
type AppService struct {
    // Does everything
}

func (a *AppService) SaveFile(path string, data []byte) error
func (a *AppService) GetUser(id int) (*User, error)
func (a *AppService) SendEmail(to, subject, body string) error
func (a *AppService) ProcessPayment(amount float64) error

// ✅ Good: Focused services
type FileService struct{}
func (f *FileService) Save(path string, data []byte) error

type UserService struct{}
func (u *UserService) GetByID(id int) (*User, error)

type EmailService struct{}
func (e *EmailService) Send(to, subject, body string) error

type PaymentService struct{}
func (p *PaymentService) Process(amount float64) error
```

### 2. 清楚的方法名稱

使用具描述性且以動作為導向的名稱：

```go
// ❌ Bad: Unclear names
func (s *Service) Do(x string) error
func (s *Service) Handle(data interface{}) interface{}
func (s *Service) Process(input map[string]interface{}) bool

// ✅ Good: Clear names
func (s *FileService) SaveDocument(path string, content string) error
func (s *UserService) AuthenticateUser(email, password string) (*User, error)
func (s *OrderService) CreateOrder(items []Item) (*Order, error)
```

### 3. 一致的回傳型別

一律明確回傳錯誤：

```go
// ❌ Bad: Inconsistent error handling
func (s *Service) GetData() interface{}  // How to handle errors?
func (s *Service) SaveData(data string)  // Silent failures?

// ✅ Good: Explicit errors
func (s *Service) GetData() (Data, error)
func (s *Service) SaveData(data string) error
```

### 4. 輸入驗證

在 Go 端驗證所有輸入：

```go
// ❌ Bad: No validation
func (s *UserService) CreateUser(email, password string) (*User, error) {
    user := &User{Email: email, Password: password}
    return s.db.Create(user)
}

// ✅ Good: Validate first
func (s *UserService) CreateUser(email, password string) (*User, error) {
    // Validate email
    if !isValidEmail(email) {
        return nil, errors.New("invalid email address")
    }
    
    // Validate password
    if len(password) < 8 {
        return nil, errors.New("password must be at least 8 characters")
    }
    
    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &User{
        Email:        email,
        PasswordHash: string(hash),
    }
    
    return s.db.Create(user)
}
```

## 效能模式

### 1. 批次作業

透過批次處理減少橋接呼叫：

```go
// ❌ Bad: N calls
// JavaScript
for (const item of items) {
    await ProcessItem(item)  // N bridge calls
}

// ✅ Good: 1 call
// Go
func (s *Service) ProcessItems(items []Item) ([]Result, error) {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = s.processItem(item)
    }
    return results, nil
}

// JavaScript
const results = await ProcessItems(items)  // 1 bridge call
```

### 2. 分頁

不要回傳龐大的資料集：

```go
// ❌ Bad: Returns everything
func (s *Service) GetAllUsers() ([]User, error) {
    return s.db.FindAll()  // Could be millions
}

// ✅ Good: Paginated
type PageRequest struct {
    Page     int `json:"page"`
    PageSize int `json:"pageSize"`
}

type PageResponse struct {
    Items      []User `json:"items"`
    TotalItems int    `json:"totalItems"`
    TotalPages int    `json:"totalPages"`
    Page       int    `json:"page"`
}

func (s *Service) GetUsers(req PageRequest) (*PageResponse, error) {
    // Validate
    if req.Page < 1 {
        req.Page = 1
    }
    if req.PageSize < 1 || req.PageSize > 100 {
        req.PageSize = 20
    }
    
    // Get total
    total, err := s.db.Count()
    if err != nil {
        return nil, err
    }
    
    // Get page
    offset := (req.Page - 1) * req.PageSize
    users, err := s.db.Find(offset, req.PageSize)
    if err != nil {
        return nil, err
    }
    
    return &PageResponse{
        Items:      users,
        TotalItems: total,
        TotalPages: (total + req.PageSize - 1) / req.PageSize,
        Page:       req.Page,
    }, nil
}
```

### 3. 快取

快取成本高昂的作業：

```go
type CachedService struct {
    cache map[string]interface{}
    mu    sync.RWMutex
    ttl   time.Duration
}

func (s *CachedService) GetData(key string) (interface{}, error) {
    // Check cache
    s.mu.RLock()
    if data, ok := s.cache[key]; ok {
        s.mu.RUnlock()
        return data, nil
    }
    s.mu.RUnlock()
    
    // Fetch data
    data, err := s.fetchData(key)
    if err != nil {
        return nil, err
    }
    
    // Cache it
    s.mu.Lock()
    s.cache[key] = data
    s.mu.Unlock()
    
    // Schedule expiry
    go func() {
        time.Sleep(s.ttl)
        s.mu.Lock()
        delete(s.cache, key)
        s.mu.Unlock()
    }()
    
    return data, nil
}
```

### 4. 使用事件串流

使用事件傳輸串流資料：

```go
// ❌ Bad: Polling
func (s *Service) GetProgress() int {
    return s.progress
}

// JavaScript polls
setInterval(async () => {
    const progress = await GetProgress()
    updateUI(progress)
}, 100)

// ✅ Good: Events
// Service must hold a reference to *application.App to emit events / log:
//
//   type Service struct {
//       app *application.App
//   }
//
//   func NewService(app *application.App) *Service {
//       return &Service{app: app}
//   }
//
//   app := application.New(application.Options{})
//   app.RegisterService(application.NewService(NewService(app)))

func (s *Service) ProcessLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    total := 0
    processed := 0
    
    // Count lines
    for scanner.Scan() {
        total++
    }
    
    // Process
    file.Seek(0, 0)
    scanner = bufio.NewScanner(file)
    
    for scanner.Scan() {
        s.processLine(scanner.Text())
        processed++
        
        // Emit progress
        s.app.Event.Emit("progress", map[string]interface{}{
            "processed": processed,
            "total":     total,
            "percent":   int(float64(processed) / float64(total) * 100),
        })
    }
    
    return scanner.Err()
}

// JavaScript listens
import { Events } from '@wailsio/runtime'

Events.On("progress", (event) => {
    updateProgress(event.data.percent)
})
```

## 安全性模式

### 1. 輸入清理

一律清理使用者輸入：

```go
import (
    "html"
    "strings"
)

func (s *Service) SaveComment(text string) error {
    // Sanitise
    text = strings.TrimSpace(text)
    text = html.EscapeString(text)
    
    // Validate length
    if len(text) == 0 {
        return errors.New("comment cannot be empty")
    }
    if len(text) > 1000 {
        return errors.New("comment too long")
    }
    
    return s.db.SaveComment(text)
}
```

### 2. 驗證身分

保護敏感作業：

```go
type AuthService struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}

func (a *AuthService) Login(email, password string) (string, error) {
    user, err := a.db.FindByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    
    if !a.verifyPassword(user.PasswordHash, password) {
        return "", errors.New("invalid credentials")
    }
    
    // Create session
    token := generateToken()
    a.mu.Lock()
    a.sessions[token] = &Session{
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    a.mu.Unlock()
    
    return token, nil
}

func (a *AuthService) requireAuth(token string) (*Session, error) {
    a.mu.RLock()
    session, ok := a.sessions[token]
    a.mu.RUnlock()
    
    if !ok {
        return nil, errors.New("not authenticated")
    }
    
    if time.Now().After(session.ExpiresAt) {
        return nil, errors.New("session expired")
    }
    
    return session, nil
}

// Protected method
func (a *AuthService) DeleteAccount(token string) error {
    session, err := a.requireAuth(token)
    if err != nil {
        return err
    }
    
    return a.db.DeleteUser(session.UserID)
}
```

### 3. 速率限制

防止濫用：

```go
type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int
    window   time.Duration
}

func (r *RateLimiter) Allow(key string) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    
    // Clean old requests
    if requests, ok := r.requests[key]; ok {
        var recent []time.Time
        for _, t := range requests {
            if now.Sub(t) < r.window {
                recent = append(recent, t)
            }
        }
        r.requests[key] = recent
    }
    
    // Check limit
    if len(r.requests[key]) >= r.limit {
        return false
    }
    
    // Add request
    r.requests[key] = append(r.requests[key], now)
    return true
}

// Usage
func (s *Service) SendEmail(to, subject, body string) error {
    if !s.rateLimiter.Allow(to) {
        return errors.New("rate limit exceeded")
    }
    
    return s.emailer.Send(to, subject, body)
}
```

## 錯誤處理模式

### 1. 描述清楚的錯誤

在錯誤中提供情境資訊：

```go
// ❌ Bad: Generic errors
func (s *Service) LoadFile(path string) ([]byte, error) {
    return os.ReadFile(path)  // "no such file or directory"
}

// ✅ Good: Contextual errors
func (s *Service) LoadFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to load file %s: %w", path, err)
    }
    return data, nil
}
```

### 2. 錯誤型別

使用具型別的錯誤進行特定處理：

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type NotFoundError struct {
    Resource string
    ID       interface{}
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s not found: %v", e.Resource, e.ID)
}

// Usage
func (s *UserService) GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, &ValidationError{
            Field:   "id",
            Message: "must be positive",
        }
    }
    
    user, err := s.db.Find(id)
    if err == sql.ErrNoRows {
        return nil, &NotFoundError{
            Resource: "User",
            ID:       id,
        }
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    return user, nil
}
```

### 3. 錯誤復原

妥善處理錯誤：

```go
func (s *Service) ProcessWithRetry(data string) error {
    maxRetries := 3
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := s.process(data)
        if err == nil {
            return nil
        }
        
        // Log attempt
        s.app.Logger.Warn("Process failed", 
            "attempt", attempt, 
            "error", err)
        
        // Don't retry on validation errors
        if _, ok := err.(*ValidationError); ok {
            return err
        }
        
        // Wait before retry
        if attempt < maxRetries {
            time.Sleep(time.Duration(attempt) * time.Second)
        }
    }
    
    return fmt.Errorf("failed after %d attempts", maxRetries)
}
```

## 測試模式

### 1. 單元測試

獨立測試各項服務：

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    db := &MockDB{}
    service := &UserService{db: db}
    
    // Test valid input
    user, err := service.CreateUser("test@example.com", "password123")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Email != "test@example.com" {
        t.Errorf("expected email test@example.com, got %s", user.Email)
    }
    
    // Test invalid email
    _, err = service.CreateUser("invalid", "password123")
    if err == nil {
        t.Error("expected error for invalid email")
    }
    
    // Test short password
    _, err = service.CreateUser("test@example.com", "short")
    if err == nil {
        t.Error("expected error for short password")
    }
}
```

### 2. 整合測試

使用實際相依項目進行測試：

```go
func TestUserService_Integration(t *testing.T) {
    // Setup real database
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    // Create schema
    _, err = db.Exec(`CREATE TABLE users (...)`)
    if err != nil {
        t.Fatal(err)
    }
    
    // Test service
    service := &UserService{db: db}
    
    user, err := service.CreateUser("test@example.com", "password123")
    if err != nil {
        t.Fatal(err)
    }
    
    // Verify in database
    var count int
    db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", 
        user.Email).Scan(&count)
    
    if count != 1 {
        t.Errorf("expected 1 user, got %d", count)
    }
}
```

### 3. 模擬服務

建立可測試的介面：

```go
type UserRepository interface {
    Create(user *User) error
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id int) error
}

type UserService struct {
    repo UserRepository
}

// Mock for testing
type MockUserRepository struct {
    users map[string]*User
}

func (m *MockUserRepository) Create(user *User) error {
    m.users[user.Email] = user
    return nil
}

// Test with mock
func TestUserService_WithMock(t *testing.T) {
    mock := &MockUserRepository{
        users: make(map[string]*User),
    }
    
    service := &UserService{repo: mock}
    
    // Test
    user, err := service.CreateUser("test@example.com", "password123")
    if err != nil {
        t.Fatal(err)
    }
    
    // Verify mock was called
    if len(mock.users) != 1 {
        t.Error("expected 1 user in mock")
    }
}
```

## 最佳實務摘要

### ✅ 應做事項

- **單一職責** — 每項服務只負責一個用途
- **清楚命名** — 使用具描述性的方法名稱
- **驗證輸入** — 一律在 Go 端進行
- **回傳錯誤** — 明確處理錯誤
- **批次處理作業** — 減少橋接呼叫
- **使用事件** — 用於串流資料
- **清理輸入** — 防止注入攻擊
- **徹底測試** — 執行單元測試與整合測試
- **記錄方法說明** — 註解會轉換為 JSDoc
- **為 API 建立版本** — 預先規劃變更

### ❌ 不應做事項

- **不要建立上帝物件** — 讓服務聚焦於單一職責
- **不要信任前端** — 驗證所有內容
- **不要回傳龐大的資料集** — 使用分頁
- **不要阻塞** — 對耗時作業使用 goroutine
- **不要忽略錯誤** — 處理所有錯誤情況
- **不要略過測試** — 儘早測試並經常測試
- **不要寫死設定值** — 使用組態
- **不要公開內部實作** — 將實作保持為私有

## 後續步驟

- [方法](/features/bindings/methods/) — 瞭解方法綁定的基礎知識
- [服務](/features/bindings/services/) — 瞭解服務架構
- [模型](/features/bindings/models/) — 綁定複雜的資料結構
- [事件](/features/events/system/) — 使用事件實作發布／訂閱模式

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[綁定範例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
