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

## Sanitasi Log

Wails secara otomatis menyamarkan data sensitif dalam log IPC untuk mencegah terbukanya rahasia secara tidak sengaja. Fitur ini **diaktifkan secara default** dengan pengaturan bawaan yang sesuai.

### Perlindungan Default

Nama bidang berikut disamarkan secara otomatis (pencocokan substring tanpa membedakan huruf besar dan kecil):

- **Autentikasi**: `password`, `passwd`, `pwd`, `token`, `bearer`, `jwt`, `access_token`, `refresh_token`, `secret`, `apikey`, `api_key`, `auth`, `authorization`, `credential`
- **Kriptografi**: `private`, `privatekey`, `private_key`, `signing`, `encryption_key`
- **Sesi**: `session`, `sessionid`, `session_id`, `cookie`, `csrf`, `xsrf`

Selain itu, pola berikut dideteksi dalam nilai:

- Token JWT (`eyJhbG...`)
- Token Bearer (`Bearer xxx`)
- Format kunci API yang umum (`sk_live_xxx`, `pk_test_xxx`)

### Konfigurasi Khusus

Konfigurasikan sanitasi dengan semua opsi yang tersedia:

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

### Opsi Konfigurasi

| Opsi | Deskripsi |
| --- | --- |
| `RedactFields` | Nama bidang tambahan untuk disamarkan (digabungkan dengan bawaan) |
| `RedactPatterns` | Ekspresi reguler tambahan untuk mencocokkan nilai |
| `CustomSanitizeFunc` | Fungsi kendali penuh; kembalikan `(value, true)` untuk mengganti nilai |
| `DisableDefaults` | Hanya gunakan bidang dan pola yang ditentukan secara eksplisit |
| `Replacement` | Teks pengganti khusus (default: `***`) |
| `Disabled` | Nonaktifkan sanitasi sepenuhnya (gunakan dengan hati-hati) |

### API Sanitizer Publik

Gunakan sanitizer untuk data Anda sendiri:

```go
// Get the application's sanitizer
sanitizer := app.Sanitizer()

// Sanitize a map
cleanData := sanitizer.SanitizeMap(sensitiveData)

// Sanitize JSON
cleanJSON := sanitizer.SanitizeJSON(jsonBytes)
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
- Jangan tampilkan data sensitif dalam log (Wails melakukan sanitasi log IPC secara default)
- Jangan gunakan enkripsi yang lemah
- Jangan abaikan pembaruan keamanan
- Jangan nonaktifkan sanitasi log di produksi

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
- [ ] Sanitasi log dikonfigurasi untuk bidang khusus aplikasi
- [ ] Bidang sensitif tidak terekspos dalam log khusus

## Langkah Berikutnya

- [Arsitektur](/guides/architecture/) - Pola arsitektur aplikasi
- [Praktik Terbaik](/features/bindings/best-practices/) - Praktik terbaik binding
