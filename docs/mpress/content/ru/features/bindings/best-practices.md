---
title: "Рекомендации по привязкам"
description: "Шаблоны проектирования и рекомендации по привязкам Go к JavaScript"
slug: "features/bindings/best-practices"
sourcePath: "features/bindings/best-practices.md"
---

## Рекомендации по привязкам

При проектировании привязок следуйте **проверенным шаблонам**, чтобы создавать понятные, производительные и безопасные привязки. В этом руководстве рассматриваются принципы проектирования API, оптимизация производительности, шаблоны обеспечения безопасности, обработка ошибок и стратегии тестирования приложений, удобных в сопровождении.

## Принципы проектирования API

### 1. Единственная ответственность

У каждого сервиса должно быть одно чётко определённое назначение:

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

### 2. Понятные имена методов

Используйте информативные имена, обозначающие действия:

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

### 3. Единообразные типы возвращаемых значений

Всегда возвращайте ошибки явно:

```go
// ❌ Bad: Inconsistent error handling
func (s *Service) GetData() interface{}  // How to handle errors?
func (s *Service) SaveData(data string)  // Silent failures?

// ✅ Good: Explicit errors
func (s *Service) GetData() (Data, error)
func (s *Service) SaveData(data string) error
```

### 4. Проверка входных данных

Проверяйте все входные данные на стороне Go:

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

## Шаблоны повышения производительности

### 1. Пакетные операции

Сокращайте количество вызовов через мост, объединяя операции в пакеты:

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

### 2. Постраничный вывод

Не возвращайте огромные наборы данных:

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

### 3. Кэширование

Кэшируйте ресурсоёмкие операции:

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

### 4. Потоковая передача с помощью событий

Используйте события для потоковой передачи данных:

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

## Шаблоны обеспечения безопасности

### 1. Очистка входных данных

Всегда очищайте пользовательские данные:

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

### 2. Аутентификация

Защищайте операции, критичные с точки зрения безопасности:

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

### 3. Ограничение частоты запросов

Предотвращайте злоупотребления:

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

## Шаблоны обработки ошибок

### 1. Информативные ошибки

Добавляйте к ошибкам контекст:

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

### 2. Типы ошибок

Используйте типизированные ошибки для специализированной обработки:

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

### 3. Восстановление после ошибок

Корректно обрабатывайте ошибки:

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

## Шаблоны тестирования

### 1. Модульное тестирование

Тестируйте сервисы изолированно:

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

### 2. Интеграционное тестирование

Тестируйте с реальными зависимостями:

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

### 3. Имитированные сервисы

Создавайте интерфейсы, удобные для тестирования:

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

## Краткий обзор рекомендаций

### ✅ Следует

- **Соблюдать принцип единственной ответственности** — один сервис, одно назначение
- **Давать понятные имена** — используйте информативные имена методов
- **Проверять входные данные** — всегда на стороне Go
- **Возвращать ошибки** — обрабатывайте ошибки явно
- **Объединять операции в пакеты** — сокращайте количество вызовов через мост
- **Использовать события** — для потоковой передачи данных
- **Очищать входные данные** — предотвращайте внедрение кода
- **Тщательно тестировать** — проводите модульные и интеграционные тесты
- **Документировать методы** — комментарии преобразуются в JSDoc
- **Версионировать API** — планируйте изменения

### ❌ Не следует

- **Создавать объекты-боги** — сохраняйте узкую специализацию сервисов
- **Доверять клиентской части** — проверяйте всё
- **Возвращать огромные наборы данных** — используйте постраничный вывод
- **Блокировать выполнение** — используйте горутины для длительных операций
- **Игнорировать ошибки** — обрабатывайте все случаи ошибок
- **Пропускать тестирование** — тестируйте как можно раньше и чаще
- **Жёстко задавать значения в коде** — используйте конфигурацию
- **Не раскрывайте внутреннюю реализацию** — оставляйте детали реализации закрытыми

## Дальнейшие шаги

- [Методы](/features/bindings/methods/) — изучите основы привязки методов
- [Сервисы](/features/bindings/services/) — разберитесь в архитектуре сервисов
- [Модели](/features/bindings/models/) — привязывайте сложные структуры данных
- [События](/features/events/system/) — используйте события для публикации и подписки

---

**Есть вопросы?** Задайте их в [Discord](https://discord.gg/JDdSxwjhGf) или ознакомьтесь с [примерами привязок](https://github.com/wailsapp/wails/tree/master/v3/examples/binding).
