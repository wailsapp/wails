---
title: "Praktik Terbaik Keamanan"
description: "Amankan aplikasi Wails Anda"
slug: "guides/security"
sourcePath: "guides/security.md"
---

## Ikhtisar

Keamanan sangat penting bagi aplikasi desktop. Ikuti praktik berikut untuk menjaga keamanan aplikasi Anda.

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

### Pengelolaan Sesi

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

Gunakan nonce baru untuk setiap pesan yang dienkripsi dengan kunci tertentu. Pada toolchain Go yang diwajibkan oleh Wails v3, `NewGCMWithRandomNonce` membuat nonce dan menambahkannya sebagai prefiks secara otomatis. Simpan byte yang dikembalikan secara utuh; `decrypt` membaca nonce dari prefiks tersebut. Simpan kunci secara terpisah dari teks sandi.

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

Gunakan kunci AES acak yang sesuai dan patuhi batas jumlah pesan per kunci pada AEAD. Lihat [referensi cipher Go](https://pkg.go.dev/crypto/cipher#NewGCMWithRandomNonce).

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

### Verifikasi Sertifikat

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

## Operasi File

### Batasi Akses File

Terima path yang relatif terhadap direktori yang dikendalikan aplikasi. `os.OpenRoot` membatasi resolusi sistem file pada direktori tersebut, termasuk untuk symlink. Menolak `..` atau hanya memeriksa prefiks string tidak memberikan batasan ini.

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

Misalnya, panggil `readFile(appDataDir, "notes/today.json")`. Pilih `appDataDir` dalam kode backend tepercaya; jangan biarkan frontend menentukan root yang diizinkan. Root membatasi penelusuran sistem file, bukan menentukan file mana di dalamnya yang boleh dibaca oleh pengguna.

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
- Jangan lewati verifikasi sertifikat
- Jangan ungkapkan data sensitif dalam log
- Jangan gunakan enkripsi yang lemah
- Jangan abaikan pembaruan keamanan

## Daftar Periksa Keamanan

- [ ] Semua input pengguna telah divalidasi
- [ ] Kata sandi di-hash dengan bcrypt
- [ ] Data sensitif telah dienkripsi
- [ ] HTTPS digunakan untuk semua panggilan jaringan
- [ ] Pembatasan laju telah diterapkan
- [ ] Path file telah divalidasi
- [ ] Dependensi telah diperbarui
- [ ] Pencatatan log keamanan telah diaktifkan
- [ ] Pesan kesalahan tidak membocorkan informasi
- [ ] Kode telah ditinjau untuk menemukan kerentanan

## Langkah Berikutnya

- [Arsitektur](/guides/architecture/) - Pola arsitektur aplikasi
- [Praktik Terbaik](/features/bindings/best-practices/) - Praktik terbaik binding
