---
title: "テスト"
description: "Wails アプリケーションをテストする"
slug: "guides/testing"
sourcePath: "guides/testing.md"
---

## 概要

テストにより、アプリケーションが正しく動作することを確認し、リグレッションを防止できます。

## 単体テスト

### サービスのテスト

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

### モックを使用したテスト

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

## 結合テスト

### 実際の依存関係を使用したテスト

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

## フロントエンドのテスト

### JavaScript の単体テスト

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

### バインディングのテスト

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

## ベストプラクティス

### ✅ 推奨事項

- バグを修正する前にテストを作成する
- エッジケースをテストする
- テーブル駆動テストを使用する
- 外部依存関係をモック化する
- エラー処理をテストする
- テストを高速に保つ

### ❌ 禁止事項

- エラーケースを省略しない
- 実装の詳細をテストしない
- 不安定なテストを作成しない
- テストの失敗を無視しない
- 結合テストを省略しない

## テストの実行

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

## 次のステップ

- [エンドツーエンドテスト](/guides/e2e-testing/) - ユーザーフロー全体をテストする
- [ベストプラクティス](/features/bindings/best-practices/) - ベストプラクティスを学ぶ
