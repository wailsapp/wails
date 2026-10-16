---
title: "보안 모범 사례"
description: "Wails 애플리케이션을 안전하게 보호하기"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## 개요

데스크톱 애플리케이션에서 보안은 매우 중요합니다. 다음 사례를 따라 애플리케이션을 안전하게 보호하세요.

## 입력값 검증

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
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### 인증서 검증하기

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

## 파일 작업

### 경로 검증하기

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

## 요청 빈도 제한

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

- 모든 입력값을 검증하세요
- 네트워크 호출에 HTTPS를 사용하세요
- 민감한 데이터를 암호화하세요
- 안전한 비밀번호 해싱을 사용하세요
- 요청 빈도 제한을 구현하세요
- 의존성을 최신 상태로 유지하세요
- 보안 이벤트를 기록하세요
- 운영체제 키체인을 사용하세요

### ❌ 하지 말아야 할 일

- 사용자 입력을 신뢰하지 마세요
- 비밀번호를 평문으로 저장하지 마세요
- 비밀 정보를 하드코딩하지 마세요
- 인증서 검증을 생략하지 마세요
- 로그에 민감한 데이터를 노출하지 마세요
- 취약한 암호화를 사용하지 마세요
- 보안 업데이트를 무시하지 마세요

## 보안 체크리스트

- [ ] 모든 사용자 입력값 검증 완료
- [ ] bcrypt로 비밀번호 해싱 완료
- [ ] 민감한 데이터 암호화 완료
- [ ] 모든 네트워크 호출에 HTTPS 사용
- [ ] 요청 빈도 제한 구현 완료
- [ ] 파일 경로 검증 완료
- [ ] 의존성을 최신 상태로 유지
- [ ] 보안 로깅 활성화
- [ ] 오류 메시지에서 정보가 유출되지 않음
- [ ] 코드의 취약점 검토 완료

## 다음 단계

- [아키텍처](/guides/architecture/) - 애플리케이션 아키텍처 패턴
- [모범 사례](/features/bindings/best-practices/) - 바인딩 모범 사례
