---
title: "Best Practices für Bindings"
description: "Entwurfsmuster und Best Practices für Go-JavaScript-Bindings"
slug: "features/bindings/best-practices"
sourcePath: "features/bindings/best-practices.md"
---

## Best Practices für Bindings

Nutzen Sie für den Entwurf von Bindings **bewährte Muster**, um übersichtliche, leistungsfähige und sichere Bindings zu erstellen. Dieser Leitfaden behandelt Grundsätze des API-Designs, Leistungsoptimierung, Sicherheitsmuster, Fehlerbehandlung und Teststrategien für wartbare Anwendungen.

## Grundsätze des API-Designs

### 1. Einzelverantwortung

Jeder Service sollte genau einen klar definierten Zweck erfüllen:

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

### 2. Eindeutige Methodennamen

Verwenden Sie aussagekräftige, handlungsorientierte Namen:

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

### 3. Einheitliche Rückgabetypen

Geben Sie Fehler immer explizit zurück:

```go
// ❌ Bad: Inconsistent error handling
func (s *Service) GetData() interface{}  // How to handle errors?
func (s *Service) SaveData(data string)  // Silent failures?

// ✅ Good: Explicit errors
func (s *Service) GetData() (Data, error)
func (s *Service) SaveData(data string) error
```

### 4. Eingabevalidierung

Validieren Sie alle Eingaben auf der Go-Seite:

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

## Leistungsmuster

### 1. Batchoperationen

Reduzieren Sie Aufrufe über die Bridge, indem Sie sie bündeln:

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

### 2. Paginierung

Geben Sie keine riesigen Datensätze zurück:

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

### 3. Caching

Cachen Sie aufwendige Operationen:

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

### 4. Streaming mit Events

Verwenden Sie Events für das Streaming von Daten:

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

## Sicherheitsmuster

### 1. Bereinigung von Eingaben

Bereinigen Sie Benutzereingaben immer:

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

### 2. Authentifizierung

Schützen Sie sensible Operationen:

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

### 3. Ratenbegrenzung

Verhindern Sie Missbrauch:

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

## Muster für die Fehlerbehandlung

### 1. Aussagekräftige Fehler

Geben Sie bei Fehlern den Kontext an:

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

### 2. Fehlertypen

Verwenden Sie typisierte Fehler für eine gezielte Behandlung:

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

### 3. Fehlerbehebung

Behandeln Sie Fehler angemessen:

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

## Testmuster

### 1. Unit-Tests

Testen Sie Services isoliert:

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

### 2. Integrationstests

Testen Sie mit echten Abhängigkeiten:

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

### 3. Mock-Services

Erstellen Sie testbare Schnittstellen:

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

## Zusammenfassung der Best Practices

### ✅ Empfohlen

- **Einzelverantwortung** – ein Service, ein Zweck
- **Eindeutige Benennung** – aussagekräftige Methodennamen
- **Eingaben validieren** – immer auf der Go-Seite
- **Fehler zurückgeben** – explizite Fehlerbehandlung
- **Operationen bündeln** – Aufrufe über die Bridge reduzieren
- **Events verwenden** – für das Streaming von Daten
- **Eingaben bereinigen** – Injection-Angriffe verhindern
- **Gründlich testen** – Unit- und Integrationstests
- **Methoden dokumentieren** – Kommentare werden zu JSDoc
- **API versionieren** – Änderungen einplanen

### ❌ Nicht tun

- **Keine God Objects erstellen** – Services auf ihren Zweck beschränken
- **Dem Frontend nicht vertrauen** – alles validieren
- **Keine riesigen Datensätze zurückgeben** – Paginierung verwenden
- **Nicht blockieren** – für lang laufende Operationen Goroutinen verwenden
- **Fehler nicht ignorieren** – alle Fehlerfälle behandeln
- **Tests nicht auslassen** – frühzeitig und häufig testen
- **Werte nicht fest codieren** – Konfiguration verwenden
- **Keine Interna offenlegen** – Implementierung nicht öffentlich zugänglich machen

## Nächste Schritte

- [Methoden](/features/bindings/methods/) – Grundlagen der Methodenbindung kennenlernen
- [Services](/features/bindings/services/) – Servicearchitektur verstehen
- [Modelle](/features/bindings/models/) – Komplexe Datenstrukturen binden
- [Ereignisse](/features/events/system/) – Ereignisse für Publish/Subscribe verwenden

---

**Fragen?** Stelle sie auf [Discord](https://discord.gg/JDdSxwjhGf) oder sieh dir die [Binding-Beispiele](https://github.com/wailsapp/wails/tree/master/v3/examples/binding) an.
