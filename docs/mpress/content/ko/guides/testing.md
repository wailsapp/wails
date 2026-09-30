---
title: "테스트"
description: "Wails 애플리케이션 테스트하기"
slug: "guides/testing"
sourcePath: "guides/testing.md"
---

## 개요

테스트를 통해 애플리케이션이 올바르게 작동하는지 확인하고 회귀를 방지할 수 있습니다.

## 단위 테스트

### 서비스 테스트

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

### 모의 객체를 사용한 테스트

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

## 통합 테스트

### 실제 의존성을 사용한 테스트

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

## 프런트엔드 테스트

### JavaScript 단위 테스트

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

### 바인딩 테스트

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

## 모범 사례

### ✅ 권장 사항

- 버그를 수정하기 전에 테스트를 작성하세요
- 경계 사례를 테스트하세요
- 테이블 기반 테스트를 사용하세요
- 외부 의존성을 모의 객체로 대체하세요
- 오류 처리를 테스트하세요
- 테스트가 빠르게 실행되도록 유지하세요

### ❌ 금지 사항

- 오류 사례를 건너뛰지 마세요
- 구현 세부 사항을 테스트하지 마세요
- 불안정한 테스트를 작성하지 마세요
- 테스트 실패를 무시하지 마세요
- 통합 테스트를 건너뛰지 마세요

## 테스트 실행

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

## 다음 단계

- [엔드 투 엔드 테스트](/guides/e2e-testing/) - 전체 사용자 흐름을 테스트하세요
- [모범 사례](/features/bindings/best-practices/) - 모범 사례를 알아보세요
