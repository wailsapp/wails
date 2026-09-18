---
title: "Tests"
description: "Testen Sie Ihre Wails-Anwendung"
slug: "guides/testing"
sourcePath: "guides/testing.md"
---

## Übersicht

Tests stellen sicher, dass Ihre Anwendung ordnungsgemäß funktioniert, und verhindern Regressionen.

## Unit-Tests

### Services testen

```go
func TestUserService_Create(t *testing.T) {
    service := &UserService{
        users: make(map[string]*User),
    }
    
    user, err := service.Create("john@example.com", "password123")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if user.Email != "john@example.com" {
        t.Errorf("expected email john@example.com, got %s", user.Email)
    }
}
```

### Tests mit Mocks

```go
type MockDB struct {
    users map[string]*User
}

func (m *MockDB) Create(user *User) error {
    m.users[user.ID] = user
    return nil
}

func TestUserService_WithMock(t *testing.T) {
    mockDB := &MockDB{users: make(map[string]*User)}
    service := &UserService{db: mockDB}
    
    user, err := service.Create("test@example.com", "pass")
    if err != nil {
        t.Fatal(err)
    }
    
    if len(mockDB.users) != 1 {
        t.Error("expected 1 user in mock")
    }
}
```

## Integrationstests

### Tests mit realen Abhängigkeiten

```go
func TestIntegration(t *testing.T) {
    // Setup test database
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
    user, err := service.Create("test@example.com", "password")
    if err != nil {
        t.Fatal(err)
    }
    
    // Verify in database
    var count int
    db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    if count != 1 {
        t.Errorf("expected 1 user, got %d", count)
    }
}
```

## Frontend-Tests

### JavaScript-Unit-Tests

```javascript
// Using Vitest
import { describe, it, expect } from 'vitest'
import { formatDate } from './utils'

describe('formatDate', () => {
    it('formats date correctly', () => {
        const date = new Date('2024-01-01')
        expect(formatDate(date)).toBe('2024-01-01')
    })
})
```

### Bindings testen

```javascript
import { vi } from 'vitest'
import { GetUser } from './bindings/changeme/userservice'

// Mock the binding
vi.mock('./bindings/changeme/userservice', () => ({
    GetUser: vi.fn()
}))

describe('User Component', () => {
    it('loads user data', async () => {
        GetUser.mockResolvedValue({ name: 'John', email: 'john@example.com' })
        
        // Test your component
        const user = await GetUser(1)
        expect(user.name).toBe('John')
    })
})
```

## Bewährte Methoden

### ✅ Empfohlen

- Schreiben Sie Tests, bevor Sie Fehler beheben
- Testen Sie Grenzfälle
- Verwenden Sie tabellengesteuerte Tests
- Simulieren Sie externe Abhängigkeiten mit Mocks
- Testen Sie die Fehlerbehandlung
- Halten Sie Tests schnell

### ❌ Nicht empfohlen

- Überspringen Sie keine Fehlerfälle
- Testen Sie keine Implementierungsdetails
- Schreiben Sie keine unzuverlässigen Tests
- Ignorieren Sie keine fehlgeschlagenen Tests
- Überspringen Sie keine Integrationstests

## Tests ausführen

```bash
# Run Go tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestUserService

# Run frontend tests
cd frontend && npm test

# Run with watch mode
cd frontend && npm test -- --watch
```

## Nächste Schritte

- [End-to-End-Tests](/guides/e2e-testing/) – Testen Sie vollständige Benutzerabläufe
- [Bewährte Methoden](/features/bindings/best-practices/) – Lernen Sie bewährte Methoden kennen
