---
title: "Tests"
description: "Testez votre application Wails"
slug: "guides/testing"
sourcePath: "guides/testing.md"
---

## Vue d’ensemble

Les tests garantissent le bon fonctionnement de votre application et évitent les régressions.

## Tests unitaires

### Test des services

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

### Tests avec des objets simulés

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

## Tests d’intégration

### Tests avec des dépendances réelles

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

## Tests du frontend

### Tests unitaires JavaScript

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

### Test des liaisons

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

## Bonnes pratiques

### ✅ À faire

- Écrivez des tests avant de corriger les bogues
- Testez les cas limites
- Utilisez des tests pilotés par des tables
- Simulez les dépendances externes
- Testez la gestion des erreurs
- Veillez à ce que les tests restent rapides

### ❌ À ne pas faire

- Ne négligez pas les cas d’erreur
- Ne testez pas les détails d’implémentation
- N’écrivez pas de tests instables
- N’ignorez pas les échecs des tests
- Ne négligez pas les tests d’intégration

## Exécution des tests

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

## Étapes suivantes

- [Tests de bout en bout](/guides/e2e-testing/) – Testez des parcours utilisateur complets
- [Bonnes pratiques](/features/bindings/best-practices/) – Découvrez les bonnes pratiques
