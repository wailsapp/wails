package signing

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// PGPKeyConfig holds configuration for generating a PGP key
type PGPKeyConfig struct {
	// Name is the name associated with the key (e.g., "John Doe")
	Name string
	// Email is the email associated with the key
	Email string
	// Comment is an optional comment (e.g., "Package Signing Key")
	Comment string
	// KeyBits is the RSA key size (default: 4096)
	KeyBits int
	// Expiry is how long until the key expires (0 = no expiry)
	Expiry time.Duration
	// Password is the password to encrypt the private key (empty = no encryption)
	Password string
}

// PGPKeyResult contains the generated key pair
type PGPKeyResult struct {
	// PrivateKeyPath is the path to the private key file
	PrivateKeyPath string
	// PublicKeyPath is the path to the public key file
	PublicKeyPath string
	// Fingerprint is the key fingerprint
	Fingerprint string
	// KeyID is the short key ID
	KeyID string
}

// GeneratePGPKey generates a new PGP key pair for package signing
func GeneratePGPKey(config PGPKeyConfig, privateKeyPath, publicKeyPath string) (*PGPKeyResult, error) {
	// Validate config
	if config.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if config.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Set defaults
	keyBits := config.KeyBits
	if keyBits == 0 {
		keyBits = 4096
	}

	// Generate RSA key
	rsaKey, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Create the entity
	now := time.Now()
	uid := packet.NewUserId(config.Name, config.Comment, config.Email)
	if uid == nil {
		return nil, fmt.Errorf("invalid user ID parameters")
	}

	// Calculate expiry
	var lifetimeSecs uint32
	if config.Expiry > 0 {
		lifetimeSecs = uint32(config.Expiry.Seconds())
	}

	// Create packet config
	packetConfig := &packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		RSABits:                keyBits,
	}

	// Create the primary key
	primaryKey := packet.NewRSAPrivateKey(now, rsaKey)

	entity := &openpgp.Entity{
		PrimaryKey: &primaryKey.PublicKey,
		PrivateKey: primaryKey,
		Identities: make(map[string]*openpgp.Identity),
	}

	// Create self-signature
	isPrimaryID := true
	selfSig := &packet.Signature{
		CreationTime:         now,
		SigType:              packet.SigTypePositiveCert,
		PubKeyAlgo:           packet.PubKeyAlgoRSA,
		Hash:                 crypto.SHA256,
		IsPrimaryId:          &isPrimaryID,
		FlagsValid:           true,
		FlagSign:             true,
		FlagCertify:          true,
		IssuerKeyId:          &primaryKey.KeyId,
		PreferredSymmetric:   []uint8{uint8(packet.CipherAES256), uint8(packet.CipherAES192), uint8(packet.CipherAES128)},
		PreferredHash:        []uint8{uint8(crypto.SHA256), uint8(crypto.SHA384), uint8(crypto.SHA512)},
		PreferredCompression: []uint8{uint8(packet.CompressionZLIB), uint8(packet.CompressionZIP)},
	}

	if lifetimeSecs > 0 {
		selfSig.KeyLifetimeSecs = &lifetimeSecs
	}

	// Sign the user ID
	if err := selfSig.SignUserId(uid.Id, &primaryKey.PublicKey, primaryKey, packetConfig); err != nil {
		return nil, fmt.Errorf("failed to sign user ID: %w", err)
	}

	entity.Identities[uid.Id] = &openpgp.Identity{
		Name:          uid.Id,
		UserId:        uid,
		SelfSignature: selfSig,
	}

	// Encrypt the private key if password is provided
	if config.Password != "" {
		if err := entity.PrivateKey.Encrypt([]byte(config.Password)); err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %w", err)
		}
	}

	// Write private key
	if err := writeArmoredKey(privateKeyPath, entity, true); err != nil {
		return nil, fmt.Errorf("failed to write private key: %w", err)
	}

	// Set restrictive permissions on private key
	if err := os.Chmod(privateKeyPath, 0600); err != nil {
		return nil, fmt.Errorf("failed to set private key permissions: %w", err)
	}

	// Write public key
	if err := writeArmoredKey(publicKeyPath, entity, false); err != nil {
		return nil, fmt.Errorf("failed to write public key: %w", err)
	}

	return &PGPKeyResult{
		PrivateKeyPath: privateKeyPath,
		PublicKeyPath:  publicKeyPath,
		Fingerprint:    fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint),
		KeyID:          entity.PrimaryKey.KeyIdString(),
	}, nil
}

// writeArmoredKey writes an entity to a file in ASCII-armored format
func writeArmoredKey(path string, entity *openpgp.Entity, private bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var blockType string
	if private {
		blockType = openpgp.PrivateKeyType
	} else {
		blockType = openpgp.PublicKeyType
	}

	w, err := armor.Encode(f, blockType, nil)
	if err != nil {
		return err
	}

	if private {
		if err := entity.SerializePrivate(w, nil); err != nil {
			w.Close()
			return err
		}
	} else {
		if err := entity.Serialize(w); err != nil {
			w.Close()
			return err
		}
	}

	return w.Close()
}

// GetPGPKeyInfo reads a PGP key file and returns information about it
func GetPGPKeyInfo(keyPath string) (*PGPKeyInfo, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	var reader io.Reader = bytes.NewReader(keyData)

	// Check if it's ASCII armored
	if bytes.HasPrefix(keyData, []byte("-----BEGIN PGP")) {
		block, err := armor.Decode(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decode armored key: %w", err)
		}
		reader = block.Body
	}

	// Read the entity
	entity, err := openpgp.ReadEntity(packet.NewReader(reader))
	if err != nil {
		return nil, fmt.Errorf("failed to read PGP entity: %w", err)
	}

	info := &PGPKeyInfo{
		KeyID:       entity.PrimaryKey.KeyIdString(),
		Fingerprint: fmt.Sprintf("%X", entity.PrimaryKey.Fingerprint),
		CreatedAt:   entity.PrimaryKey.CreationTime,
		HasPrivate:  entity.PrivateKey != nil,
		IsEncrypted: entity.PrivateKey != nil && entity.PrivateKey.Encrypted,
	}

	// Get user IDs
	for name := range entity.Identities {
		info.UserIDs = append(info.UserIDs, name)
	}

	// Check expiry
	for _, id := range entity.Identities {
		if id.SelfSignature != nil && id.SelfSignature.KeyLifetimeSecs != nil {
			expiry := entity.PrimaryKey.CreationTime.Add(time.Duration(*id.SelfSignature.KeyLifetimeSecs) * time.Second)
			info.ExpiresAt = &expiry
			break
		}
	}

	return info, nil
}

// PGPKeyInfo contains information about a PGP key
type PGPKeyInfo struct {
	KeyID       string
	Fingerprint string
	UserIDs     []string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	HasPrivate  bool
	IsEncrypted bool
}

// IsExpired returns true if the key has expired
func (i *PGPKeyInfo) IsExpired() bool {
	if i.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*i.ExpiresAt)
}
