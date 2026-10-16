---
title: "Praktik Terbaik Keamanan"
description: "Amankan aplikasi Wails Anda"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Ikhtisar

Keamanan sangat penting untuk aplikasi desktop. Ikuti praktik berikut untuk menjaga keamanan aplikasi Anda.

## Validasi Input

### Selalu Lakukan Validasi

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

### Sanitasi HTML

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

## Autentikasi

### Penyimpanan Kata Sandi yang Aman

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

### Manajemen Sesi

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

## Perlindungan Data

### Enkripsi Data Sensitif

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

### Penyimpanan Aman

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

## Keamanan Jaringan

### Gunakan HTTPS

```go
func makeAPICall(url string) (*Response, error) {
    // Always use HTTPS
    if !strings.HasPrefix(url, "https://") {
        return nil, errors.New("only HTTPS allowed")
    }
    
    return http.Get(url)
}
```

### Verifikasi Sertifikat

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

## Operasi File

### Validasi Jalur File

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

## Pembatasan Laju

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

## Praktik Terbaik

### ✅ Lakukan

- Validasi semua input
- Gunakan HTTPS untuk panggilan jaringan
- Enkripsi data sensitif
- Gunakan hashing kata sandi yang aman
- Terapkan pembatasan laju
- Selalu perbarui dependensi
- Catat peristiwa keamanan dalam log
- Gunakan keychain sistem operasi

### ❌ Jangan Lakukan

- Jangan percayai input pengguna
- Jangan simpan kata sandi sebagai teks biasa
- Jangan tanamkan rahasia langsung dalam kode
- Jangan lewatkan verifikasi sertifikat
- Jangan tampilkan data sensitif dalam log
- Jangan gunakan enkripsi yang lemah
- Jangan abaikan pembaruan keamanan

## Daftar Periksa Keamanan

- [ ] Semua input pengguna telah divalidasi
- [ ] Kata sandi di-hash dengan bcrypt
- [ ] Data sensitif telah dienkripsi
- [ ] HTTPS digunakan untuk semua panggilan jaringan
- [ ] Pembatasan laju telah diterapkan
- [ ] Jalur file telah divalidasi
- [ ] Dependensi telah diperbarui
- [ ] Pencatatan log keamanan telah diaktifkan
- [ ] Pesan kesalahan tidak membocorkan informasi
- [ ] Kode telah ditinjau untuk menemukan kerentanan

## Langkah Berikutnya

- [Arsitektur](/guides/architecture/) - Pola arsitektur aplikasi
- [Praktik Terbaik](/features/bindings/best-practices/) - Praktik terbaik binding
