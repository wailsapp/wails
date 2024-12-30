package application

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var alreadyRunningError = errors.New("application is already running")
var secondInstanceBuffer = make(chan string, 1)
var once sync.Once

// SecondInstanceData contains information about the second instance launch
type SecondInstanceData struct {
	Args           []string          `json:"args"`
	WorkingDir     string            `json:"workingDir"`
	AdditionalData map[string]string `json:"additionalData,omitempty"`
}

// SingleInstanceOptions defines options for single instance functionality
type SingleInstanceOptions struct {
	// UniqueID is used to identify the application instance
	// This should be unique per application, e.g. "com.myapp.myapplication"
	UniqueID string

	// OnSecondInstanceLaunch is called when a second instance of the application is launched
	// The callback receives data about the second instance launch
	OnSecondInstanceLaunch func(data SecondInstanceData)

	// AdditionalData allows passing custom data from second instance to first
	AdditionalData map[string]string

	// ExitCode is the exit code to use when the second instance exits
	ExitCode int

	// EncryptionKey is a 32-byte key used for encrypting instance communication
	// If not provided (zero array), data will be sent unencrypted
	EncryptionKey [32]byte
}

// platformLock is the interface that platform-specific lock implementations must implement
type platformLock interface {
	// acquire attempts to acquire the lock
	acquire(uniqueID string) error
	// release releases the lock and cleans up resources
	release()
	// notify sends data to the first instance
	notify(data string) error
}

// singleInstanceManager handles the single instance functionality
type singleInstanceManager struct {
	options *SingleInstanceOptions
	lock    platformLock
	app     *App
}

func newSingleInstanceManager(app *App, options *SingleInstanceOptions) (*singleInstanceManager, error) {
	if options == nil {
		return nil, nil
	}

	manager := &singleInstanceManager{
		options: options,
		app:     app,
	}

	// Launch second instance data listener
	once.Do(func() {
		go func() {
			for encryptedData := range secondInstanceBuffer {
				var secondInstanceData SecondInstanceData
				var jsonData []byte
				var err error

				// Check if encryption key is non-zero
				var zeroKey [32]byte
				if options.EncryptionKey != zeroKey {
					// Try to decrypt the data
					jsonData, err = decrypt(options.EncryptionKey, encryptedData)
					if err != nil {
						continue // Skip invalid data
					}
				} else {
					jsonData = []byte(encryptedData)
				}

				if err := json.Unmarshal(jsonData, &secondInstanceData); err == nil && manager.options.OnSecondInstanceLaunch != nil {
					manager.options.OnSecondInstanceLaunch(secondInstanceData)
				}
			}
		}()
	})

	// Create platform-specific lock
	lock, err := newPlatformLock(manager)
	if err != nil {
		return nil, err
	}

	manager.lock = lock

	// Try to acquire the lock
	err = lock.acquire(options.UniqueID)
	if err != nil {
		return manager, err
	}

	return manager, nil
}

func (m *singleInstanceManager) cleanup() {
	if m == nil || m.lock == nil {
		return
	}
	m.lock.release()
}

// encrypt encrypts data using AES-256-GCM
func encrypt(key [32]byte, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	encrypted := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// decrypt decrypts data using AES-256-GCM
func decrypt(key [32]byte, encrypted string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	if len(data) < 12 {
		return nil, errors.New("invalid encrypted data")
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := data[:12]
	ciphertext := data[12:]

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

// notifyFirstInstance sends data to the first instance of the application
func (m *singleInstanceManager) notifyFirstInstance() error {
	data := SecondInstanceData{
		Args:           os.Args,
		WorkingDir:     getCurrentWorkingDir(),
		AdditionalData: m.options.AdditionalData,
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Check if encryption key is non-zero
	var zeroKey [32]byte
	if m.options.EncryptionKey != zeroKey {
		encrypted, err := encrypt(m.options.EncryptionKey, serialized)
		if err != nil {
			return err
		}
		return m.lock.notify(encrypted)
	}

	return m.lock.notify(string(serialized))
}

func getCurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// getLockPath returns the path to the lock file for Unix systems
func getLockPath(uniqueID string) string {
	// Use system temp directory
	tmpDir := os.TempDir()
	lockFileName := uniqueID + ".lock"
	return filepath.Join(tmpDir, lockFileName)
}
