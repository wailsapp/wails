---
title: "보안 모범 사례"
description: "Wails 애플리케이션을 안전하게 보호하기"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 개요

데스크톱 애플리케이션에서 보안은 매우 중요합니다. 애플리케이션을 안전하게 보호하려면 다음 사례를 따르세요.

## 입력 검증

### 항상 검증하기

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

### HTML 무해화하기

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

## 인증

### 안전한 비밀번호 저장

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

### 세션 관리

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

## 데이터 보호

### 민감한 데이터 암호화하기

주어진 키로 암호화하는 메시지마다 새로운 nonce를 사용하세요. Wails v3에 필요한 Go 툴체인에서는 `NewGCMWithRandomNonce`가 nonce를 자동으로 생성하여 앞에 추가합니다. 반환된 바이트를 변경하지 말고 그대로 저장하세요. `decrypt`는 해당 접두부에서 nonce를 읽습니다. 키는 암호문과 별도로 보관하세요.

```go
import (
    "crypto/aes"
    "crypto/cipher"
)

func encrypt(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aead, err := cipher.NewGCMWithRandomNonce(block)
    if err != nil {
        return nil, err
    }
    return aead.Seal(nil, nil, data, nil), nil
}

func decrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aead, err := cipher.NewGCMWithRandomNonce(block)
    if err != nil {
        return nil, err
    }
    return aead.Open(nil, nil, ciphertext, nil)
}
```

적절한 무작위 생성 AES 키를 사용하고 AEAD의 키별 메시지 한도를 준수하세요. [Go 암호화 참조 문서](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce)를 확인하세요.

### 안전한 저장소

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

## 네트워크 보안

### HTTPS 사용하기

```go
import (
    "errors"
    "net/http"
    "strings"
)

func makeAPICall(url string) (*http.Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return secureClient().Get(url)
}
```

### 인증서 검증하기

```go
import (
    "crypto/tls"
    "errors"
    "net/http"
    "time"
)

func secureClient() *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if req.URL.Scheme != "https" {
                return errors.New("redirect to non-HTTPS URL denied")
            }
            if len(via) >= 10 {
                return errors.New("too many redirects")
            }
            return nil
        },
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
            },
        },
    }
}
```

## 파일 작업

### 파일 접근 제한하기

애플리케이션이 제어하는 디렉터리를 기준으로 한 상대 경로만 허용하세요. `os.OpenRoot`는 심볼릭 링크를 포함한 파일 시스템 경로 확인을 해당 디렉터리 내부로 제한합니다. `..`을 거부하거나 문자열 접두부만 검사해서는 이러한 경계를 보장할 수 없습니다.

```go
import "os"

func readFile(allowedDir, name string) ([]byte, error) {
    root, err := os.OpenRoot(allowedDir)
    if err != nil {
        return nil, err
    }
    defer root.Close()
    return root.ReadFile(name)
}
```

예를 들어 `readFile(appDataDir, "notes/today.json")`을 호출하세요. 신뢰할 수 있는 백엔드 코드에서 `appDataDir`을 선택하고, 프런트엔드가 허용되는 루트를 제공하도록 해서는 안 됩니다. 루트는 파일 시스템 탐색 범위를 제한할 뿐이며, 그 안에서 사용자가 읽을 권한이 있는 파일까지 결정하지는 않습니다.

## 요청 속도 제한

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

## 모범 사례

### ✅ 해야 할 일

- 모든 입력을 검증하세요
- 네트워크 호출에 HTTPS를 사용하세요
- 민감한 데이터를 암호화하세요
- 안전한 비밀번호 해싱을 사용하세요
- 요청 속도 제한을 구현하세요
- 종속성을 최신 상태로 유지하세요
- 보안 이벤트를 기록하세요
- OS 키체인을 사용하세요

### ❌ 하지 말아야 할 일

- 사용자 입력을 신뢰하지 마세요
- 비밀번호를 평문으로 저장하지 마세요
- 비밀 정보를 하드코딩하지 마세요
- 인증서 검증을 건너뛰지 마세요
- 로그에 민감한 데이터를 노출하지 마세요
- 취약한 암호화를 사용하지 마세요
- 보안 업데이트를 무시하지 마세요

## 보안 체크리스트

- [ ] 모든 사용자 입력 검증 완료
- [ ] bcrypt로 비밀번호 해싱 완료
- [ ] 민감한 데이터 암호화 완료
- [ ] 모든 네트워크 호출에 HTTPS 사용
- [ ] 요청 속도 제한 구현 완료
- [ ] 파일 경로 검증 완료
- [ ] 종속성이 최신 상태임
- [ ] 보안 로깅 활성화
- [ ] 오류 메시지에서 정보가 유출되지 않음
- [ ] 코드의 취약점 검토 완료

## 다음 단계

- [아키텍처](/guides/architecture/) - 애플리케이션 아키텍처 패턴
- [모범 사례](/features/bindings/best-practices/) - 바인딩 모범 사례
