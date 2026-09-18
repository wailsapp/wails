---
title: "바인딩 모범 사례"
description: "Go-JavaScript 바인딩을 위한 디자인 패턴과 모범 사례"
slug: "features/bindings/best-practices"
sourcePath: "features/bindings/best-practices.md"
---

## 바인딩 모범 사례

깔끔하고 성능이 뛰어나며 안전한 바인딩을 만들려면 바인딩 설계에 <strong>검증된 패턴</strong>을 따르세요. 이 가이드에서는 유지보수하기 쉬운 애플리케이션을 위한 API 설계 원칙, 성능 최적화, 보안 패턴, 오류 처리 및 테스트 전략을 다룹니다.

## API 설계 원칙

### 1. 단일 책임

각 서비스에는 하나의 명확한 목적만 있어야 합니다:

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

### 2. 명확한 메서드 이름

의미가 명확하고 동작을 나타내는 이름을 사용하세요:

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

### 3. 일관된 반환 타입

오류는 항상 명시적으로 반환하세요:

```go
// ❌ Bad: Inconsistent error handling
func (s *Service) GetData() interface{}  // How to handle errors?
func (s *Service) SaveData(data string)  // Silent failures?

// ✅ Good: Explicit errors
func (s *Service) GetData() (Data, error)
func (s *Service) SaveData(data string) error
```

### 4. 입력 유효성 검사

모든 입력의 유효성을 Go 측에서 검사하세요:

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

## 성능 패턴

### 1. 일괄 작업

작업을 일괄 처리하여 브리지 호출을 줄이세요:

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

### 2. 페이지네이션

방대한 데이터 세트를 반환하지 마세요:

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

### 3. 캐싱

비용이 많이 드는 작업의 결과를 캐시하세요:

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

### 4. 이벤트를 사용한 스트리밍

데이터 스트리밍에는 이벤트를 사용하세요:

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

## 보안 패턴

### 1. 입력 정제

사용자 입력은 항상 정제하세요:

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

### 2. 인증

민감한 작업을 보호하세요:

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

### 3. 요청 속도 제한

악용을 방지하세요:

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

## 오류 처리 패턴

### 1. 설명이 명확한 오류

오류에 관련 맥락을 포함하세요:

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

### 2. 오류 타입

상황별로 처리할 수 있도록 타입이 지정된 오류를 사용하세요:

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

### 3. 오류 복구

오류를 적절하게 처리하세요:

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

## 테스트 패턴

### 1. 단위 테스트

서비스를 격리하여 테스트하세요:

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

### 2. 통합 테스트

실제 종속성을 사용하여 테스트하세요:

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

### 3. 모의 서비스

테스트 가능한 인터페이스를 만드세요:

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

## 모범 사례 요약

### ✅ 권장 사항

- **단일 책임** - 서비스 하나에 목적 하나
- **명확한 이름 지정** - 의미가 명확한 메서드 이름
- **입력 유효성 검사** - 항상 Go 측에서 수행
- **오류 반환** - 명시적인 오류 처리
- **일괄 작업** - 브리지 호출 감소
- **이벤트 사용** - 데이터 스트리밍에 사용
- **입력 정제** - 인젝션 방지
- **철저한 테스트** - 단위 테스트 및 통합 테스트
- **메서드 문서화** - 주석이 JSDoc으로 변환됨
- **API 버전 관리** - 변경에 대비

### ❌ 금지 사항

- **갓 오브젝트를 만들지 마세요** - 서비스의 책임 범위를 명확하게 유지하세요
- **프런트엔드를 신뢰하지 마세요** - 모든 항목의 유효성을 검사하세요
- **방대한 데이터 세트를 반환하지 마세요** - 페이지네이션을 사용하세요
- **실행을 블로킹하지 마세요** - 오래 걸리는 작업에는 goroutine을 사용하세요
- **오류를 무시하지 마세요** - 모든 오류 상황을 처리하세요
- **테스트를 생략하지 마세요** - 일찍부터 자주 테스트하세요
- **하드코딩하지 마세요** - 설정을 사용하세요
- **내부 구현을 노출하지 마세요** - 구현은 비공개로 유지하세요

## 다음 단계

- [메서드](/features/bindings/methods/) - 메서드 바인딩의 기본 사항을 알아보세요
- [서비스](/features/bindings/services/) - 서비스 아키텍처를 이해하세요
- [모델](/features/bindings/models/) - 복잡한 데이터 구조를 바인딩하세요
- [이벤트](/features/events/system/) - 게시/구독에 이벤트를 사용하세요

---

**질문이 있으신가요?** [Discord](https://discord.gg/JDdSxwjhGf)에서 질문하거나 [바인딩 예제](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)를 확인하세요.
