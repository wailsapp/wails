---
title: "セキュリティのベストプラクティス"
description: "Wails アプリケーションを安全に保つ"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 概要

デスクトップアプリケーションでは、セキュリティが極めて重要です。以下のプラクティスに従って、アプリケーションを安全に保ってください。

## 入力の検証

### 必ず検証する

```go
func (s *UserService) CreateUser(email, password string) (*User, error) {
    // Validate email
    if !isValidEmail(email) {
        return nil, errors.New("invalid email")
    }
    
    // Validate password strength
    if len(password) < 8 {
        return nil, errors.New("password too short")
    }
    
    // Sanitise input
    email = strings.TrimSpace(email)
    email = html.EscapeString(email)
    
    // Continue...
}
```

### HTML をサニタイズする

```go
import "html"

func (s *Service) SaveComment(text string) error {
    // Escape HTML
    text = html.EscapeString(text)
    
    // Validate length
    if len(text) > 1000 {
        return errors.New("comment too long")
    }
    
    return s.db.Save(text)
}
```

## 認証

### パスワードを安全に保存する

```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), err
}

func verifyPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### セッション管理

```go
type Session struct {
    UserID    int
    Token     string
    ExpiresAt time.Time
}

func (a *AuthService) CreateSession(userID int) (*Session, error) {
    token := generateSecureToken()
    
    session := &Session{
        UserID:    userID,
        Token:     token,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    
    return session, a.saveSession(session)
}
```

## データ保護

### 機密データを暗号化する

```go
import "crypto/aes"
import "crypto/cipher"

func encrypt(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    return gcm.Seal(nonce, nonce, data, nil), nil
}
```

### 安全なストレージ

```go
// Use OS keychain for sensitive data
import "github.com/zalando/go-keyring"

func saveAPIKey(key string) error {
    return keyring.Set("myapp", "api_key", key)
}

func getAPIKey() (string, error) {
    return keyring.Get("myapp", "api_key")
}
```

## ネットワークセキュリティ

### HTTPS を使用する

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### 証明書を検証する

```go
import "crypto/tls"

func secureClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
            },
        },
    }
}
```

## ファイル操作

### パスを検証する

```go
func readFile(path string) ([]byte, error) {
    // Prevent path traversal
    if strings.Contains(path, "..") {
        return nil, errors.New("invalid path")
    }
    
    // Check file exists in allowed directory
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, err
    }
    
    if !strings.HasPrefix(absPath, allowedDir) {
        return nil, errors.New("access denied")
    }
    
    return os.ReadFile(absPath)
}
```

## レート制限

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
    var recent []time.Time
    for _, t := range r.requests[key] {
        if now.Sub(t) < r.window {
            recent = append(recent, t)
        }
    }
    
    if len(recent) >= r.limit {
        return false
    }
    
    r.requests[key] = append(recent, now)
    return true
}
```

## ログのサニタイズ

Wails は IPC ログ内の機密データを自動的にマスクし、秘密情報の意図しない露出を防ぎます。この機能は、適切な既定設定で**デフォルトで有効**になっています。

### 既定の保護

以下のフィールド名は自動的にマスクされます（大文字と小文字を区別しない部分文字列一致）。

- **認証**: `password`, `passwd`, `pwd`, `token`, `bearer`, `jwt`, `access_token`, `refresh_token`, `secret`, `apikey`, `api_key`, `auth`, `authorization`, `credential`
- **暗号**: `private`, `privatekey`, `private_key`, `signing`, `encryption_key`
- **セッション**: `session`, `sessionid`, `session_id`, `cookie`, `csrf`, `xsrf`

さらに、値に含まれる以下のパターンも検出されます。

- JWT トークン (`eyJhbG...`)
- Bearer トークン (`Bearer xxx`)
- 一般的な API キー形式 (`sk_live_xxx`, `pk_test_xxx`)

### カスタム設定

利用可能なすべてのオプションを使ってサニタイズを設定します。

```go
app := application.New(application.Options{
    Name: "MyApp",
    SanitizeOptions: &application.SanitizeOptions{
        // RedactFields: additional field names to redact (merged with defaults)
        RedactFields: []string{"cardNumber", "cvv", "ssn"},

        // RedactPatterns: additional regex patterns to match values
        RedactPatterns: []*regexp.Regexp{
            regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`), // card numbers
            regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`), // SSN format
        },

        // CustomSanitizeFunc: full control - return (value, true) to override
        CustomSanitizeFunc: func(key string, value any, path string) (any, bool) {
            // Custom handling for specific paths
            if strings.HasPrefix(path, "payment.") && key != "amount" {
                return "[PAYMENT_REDACTED]", true
            }
            return nil, false // fall through to default logic
        },

        // Replacement: custom replacement string (default: "***")
        Replacement: "[REDACTED]",

        // DisableDefaults: if true, only use explicitly specified fields/patterns
        // DisableDefaults: false,

        // Disabled: completely disable sanitization
        // Disabled: false,
    },
})
```

### 設定オプション

| オプション | 説明 |
| --- | --- |
| `RedactFields` | マスクする追加のフィールド名（既定値と統合） |
| `RedactPatterns` | 値との照合に使う追加の正規表現 |
| `CustomSanitizeFunc` | 完全に制御するための関数。値を置き換えるには `(value, true)` を返す |
| `DisableDefaults` | 明示的に指定されたフィールドとパターンだけを使用 |
| `Replacement` | カスタムの置換文字列（デフォルト：`***`） |
| `Disabled` | サニタイズを完全に無効化する（注意して使用） |

### 公開サニタイザー API

独自のデータにもサニタイザーを使用できます。

```go
// Get the application's sanitizer
sanitizer := app.Sanitizer()

// Sanitize a map
cleanData := sanitizer.SanitizeMap(sensitiveData)

// Sanitize JSON
cleanJSON := sanitizer.SanitizeJSON(jsonBytes)
```

## ベストプラクティス

### ✅ 推奨事項

- すべての入力を検証する
- ネットワーク通信には HTTPS を使用する
- 機密データを暗号化する
- 安全なパスワードハッシュを使用する
- レート制限を実装する
- 依存関係を最新の状態に保つ
- セキュリティイベントをログに記録する
- OS のキーチェーンを使用する

### ❌ 禁止事項

- ユーザー入力を信用しない
- パスワードを平文で保存しない
- シークレットをハードコードしない
- 証明書の検証を省略しない
- 機密データをログに出力しない (Wails はデフォルトで IPC ログをサニタイズします)
- 脆弱な暗号化を使用しない
- セキュリティ更新を無視しない
- 本番環境でログのサニタイズを無効にしない

## セキュリティチェックリスト

- [ ] すべてのユーザー入力を検証済み
- [ ] bcrypt でパスワードをハッシュ化済み
- [ ] 機密データを暗号化済み
- [ ] すべてのネットワーク通信で HTTPS を使用
- [ ] レート制限を実装済み
- [ ] ファイルパスを検証済み
- [ ] 依存関係が最新
- [ ] セキュリティログを有効化済み
- [ ] エラーメッセージから情報が漏えいしない
- [ ] 脆弱性についてコードをレビュー済み
- [ ] アプリケーション固有のフィールドに対してログのサニタイズを設定している
- [ ] 独自のログ出力で機密フィールドを露出していない

## 次のステップ

- [アーキテクチャ](/guides/architecture/) - アプリケーションのアーキテクチャパターン
- [ベストプラクティス](/features/bindings/best-practices/) - バインディングのベストプラクティス
