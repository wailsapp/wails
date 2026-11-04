package codesignature

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"unicode/utf16"
)

var (
	oidData           = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	oidEncryptedData  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 6}
	oidKeyBag         = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 10, 1, 1}
	oidShroudedKeyBag = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 10, 1, 2}
	oidCertBag        = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 10, 1, 3}
	oidX509Cert       = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 22, 1}
	oidPBES2          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 13}
	oidPBKDF2         = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 12}
	oidTripleDES      = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 12, 1, 3}
	oidSHA1           = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
	oidSHA256         = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
)

// ErrUnsupportedPFXEncryption identifies an authenticated PFX encryption profile
// that a platform signer may support. Authentication failures never wrap it.
var ErrUnsupportedPFXEncryption = errors.New("unsupported PFX encryption profile")

type contentInfo struct {
	Type    asn1.ObjectIdentifier
	Content asn1.RawValue `asn1:"optional,explicit,tag:0"`
}
type digestInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	Digest    []byte
}
type pfxMAC struct {
	Digest     digestInfo
	Salt       []byte
	Iterations int `asn1:"optional,default:1"`
}
type pfxDocument struct {
	Version int
	Safe    contentInfo
	MAC     pfxMAC
}
type encryptedValue struct {
	Algorithm pkix.AlgorithmIdentifier
	Data      []byte
}
type safeBag struct {
	ID         asn1.ObjectIdentifier
	Value      asn1.RawValue   `asn1:"explicit,tag:0"`
	Attributes []asn1.RawValue `asn1:"optional,set"`
}
type pfxDecoder struct {
	password string
	budget   int
	bags     int
	key      crypto.Signer
	certs    []*x509.Certificate
}

// DecodePFX authenticates a password-protected PFX before decrypting its bags.
// Supported encryption is PBES2/PBKDF2/AES-CBC and legacy SHA1/3-key-3DES. RC2
// containers must be re-exported; integrity is never skipped or downgraded.
func DecodePFX(data []byte, password string) (Identity, error) {
	if len(data) == 0 || len(data) > 16<<20 || len(password) > 4096 {
		return Identity{}, errors.New("PFX or password exceeds supported size")
	}
	var document pfxDocument
	if err := decodeDER(data, &document); err != nil {
		return Identity{}, fmt.Errorf("decode PFX: %w", err)
	}
	if document.Version != 3 || !document.Safe.Type.Equal(oidData) {
		return Identity{}, errors.New("PFX requires version 3 password-MAC integrity")
	}
	var contents []byte
	if err := decodeDER(document.Safe.Content.Bytes, &contents); err != nil {
		return Identity{}, err
	}
	decoder := pfxDecoder{password: password, budget: 2_000_000}
	if err := decoder.verifyMAC(contents, document.MAC); err != nil {
		return Identity{}, err
	}
	var safes []contentInfo
	if err := decodeDER(contents, &safes); err != nil {
		return Identity{}, err
	}
	if len(safes) > 64 {
		return Identity{}, errors.New("too many PFX safe contents")
	}
	for _, safe := range safes {
		plain, err := decoder.safeContents(safe)
		if err != nil {
			return Identity{}, err
		}
		err = decoder.readBags(plain)
		clear(plain)
		if err != nil {
			return Identity{}, err
		}
	}
	return identityFromParts(decoder.key, decoder.certs)
}

func decodeDER(data []byte, value any) error {
	rest, err := asn1.Unmarshal(data, value)
	if err != nil {
		return err
	}
	if len(rest) != 0 {
		return errors.New("trailing ASN.1 data")
	}
	return nil
}

func (d *pfxDecoder) useIterations(iterations int) error {
	if iterations < 1 || iterations > 1_000_000 || iterations > d.budget {
		return errors.New("PFX key derivation exceeds supported iteration budget")
	}
	d.budget -= iterations
	return nil
}

func digestHash(oid asn1.ObjectIdentifier) (func() hash.Hash, error) {
	switch oid.String() {
	case "1.3.14.3.2.26":
		return sha1.New, nil
	case "2.16.840.1.101.3.4.2.1":
		return sha256.New, nil
	case "2.16.840.1.101.3.4.2.2":
		return sha512.New384, nil
	case "2.16.840.1.101.3.4.2.3":
		return sha512.New, nil
	default:
		return nil, fmt.Errorf("unsupported PFX digest algorithm %s", oid)
	}
}

func pfxPassword(password string) []byte {
	units := utf16.Encode([]rune(password))
	result := make([]byte, 2*(len(units)+1))
	for i, value := range units {
		binary.BigEndian.PutUint16(result[2*i:], value)
	}
	return result
}

func (d *pfxDecoder) verifyMAC(data []byte, mac pfxMAC) error {
	if parameters := mac.Digest.Algorithm.Parameters.FullBytes; len(parameters) != 0 && !bytes.Equal(parameters, []byte{5, 0}) {
		return errors.New("invalid PFX digest parameters")
	}
	if err := d.useIterations(mac.Iterations); err != nil {
		return err
	}
	if len(mac.Salt) < 1 || len(mac.Salt) > 1024 {
		return errors.New("invalid PFX MAC salt")
	}
	newHash, err := digestHash(mac.Digest.Algorithm.Algorithm)
	if err != nil {
		return err
	}
	password := pfxPassword(d.password)
	defer clear(password)
	key := derivePFX(newHash, 3, password, mac.Salt, mac.Iterations, newHash().Size())
	defer clear(key)
	hash := hmac.New(newHash, key)
	_, _ = hash.Write(data)
	if !hmac.Equal(hash.Sum(nil), mac.Digest.Digest) {
		return errors.New("PFX integrity check failed (wrong password or modified data)")
	}
	return nil
}

// derivePFX follows RFC 7292 Appendix B with a bounded output length. Hashes,
// HMAC and block ciphers themselves are provided by Go's crypto packages.
func derivePFX(newHash func() hash.Hash, purpose byte, password, salt []byte, iterations, length int) []byte {
	h := newHash()
	blockSize := h.BlockSize()
	repeat := func(value []byte) []byte {
		if len(value) == 0 {
			return nil
		}
		output := make([]byte, blockSize*((len(value)+blockSize-1)/blockSize))
		for i := range output {
			output[i] = value[i%len(value)]
		}
		return output
	}
	diversifier := bytes.Repeat([]byte{purpose}, blockSize)
	input := append(repeat(salt), repeat(password)...)
	defer clear(input)
	output := make([]byte, 0, length+h.Size())
	for len(output) < length {
		h.Reset()
		_, _ = h.Write(diversifier)
		_, _ = h.Write(input)
		a := h.Sum(nil)
		for i := 1; i < iterations; i++ {
			h.Reset()
			_, _ = h.Write(a)
			a = h.Sum(a[:0])
		}
		output = append(output, a...)
		for start := 0; start < len(input); start += blockSize {
			carry := uint16(1)
			for i := blockSize - 1; i >= 0; i-- {
				sum := uint16(input[start+i]) + uint16(a[i%len(a)]) + carry
				input[start+i] = byte(sum)
				carry = sum >> 8
			}
		}
		clear(a)
	}
	return output[:length]
}

func (d *pfxDecoder) safeContents(info contentInfo) ([]byte, error) {
	switch {
	case info.Type.Equal(oidData):
		var data []byte
		err := decodeDER(info.Content.Bytes, &data)
		return bytes.Clone(data), err
	case info.Type.Equal(oidEncryptedData):
		var value struct {
			Version int
			Info    struct {
				Type      asn1.ObjectIdentifier
				Algorithm pkix.AlgorithmIdentifier
				Data      []byte `asn1:"tag:0"`
			}
		}
		if err := decodeDER(info.Content.Bytes, &value); err != nil {
			return nil, err
		}
		if value.Version != 0 || !value.Info.Type.Equal(oidData) {
			return nil, errors.New("unsupported encrypted PFX safe contents")
		}
		return d.decrypt(encryptedValue{value.Info.Algorithm, value.Info.Data})
	default:
		return nil, fmt.Errorf("unsupported PFX content type %s", info.Type)
	}
}

func (d *pfxDecoder) readBags(data []byte) error {
	var bags []safeBag
	if err := decodeDER(data, &bags); err != nil {
		return err
	}
	d.bags += len(bags)
	if d.bags > 64 {
		return errors.New("too many PFX safe bags")
	}
	for _, bag := range bags {
		switch {
		case bag.ID.Equal(oidCertBag):
			var cert struct {
				ID   asn1.ObjectIdentifier
				Data []byte `asn1:"explicit,tag:0"`
			}
			if err := decodeDER(bag.Value.Bytes, &cert); err != nil {
				return err
			}
			if !cert.ID.Equal(oidX509Cert) {
				return errors.New("PFX contains a non-X509 certificate")
			}
			parsed, err := x509.ParseCertificate(bytes.Clone(cert.Data))
			if err != nil {
				return err
			}
			d.certs = append(d.certs, parsed)
		case bag.ID.Equal(oidKeyBag):
			if err := d.readKey(bag.Value.Bytes); err != nil {
				return err
			}
		case bag.ID.Equal(oidShroudedKeyBag):
			var encrypted encryptedValue
			if err := decodeDER(bag.Value.Bytes, &encrypted); err != nil {
				return err
			}
			key, err := d.decrypt(encrypted)
			if err != nil {
				return err
			}
			err = d.readKey(key)
			clear(key)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported PFX safe bag %s", bag.ID)
		}
	}
	return nil
}

func (d *pfxDecoder) readKey(data []byte) error {
	if d.key != nil {
		return errors.New("PFX contains multiple private keys")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return err
	}
	key, ok := parsed.(crypto.Signer)
	if !ok {
		return errors.New("PFX private key cannot sign")
	}
	d.key = key
	return nil
}

func (d *pfxDecoder) decrypt(value encryptedValue) ([]byte, error) {
	var block cipher.Block
	var iv []byte
	var err error
	switch {
	case value.Algorithm.Algorithm.Equal(oidPBES2):
		block, iv, err = d.pbes2(value.Algorithm.Parameters.FullBytes)
	case value.Algorithm.Algorithm.Equal(oidTripleDES):
		block, iv, err = d.tripleDES(value.Algorithm.Parameters.FullBytes)
	default:
		return nil, fmt.Errorf("%w: %s; re-export using AES or 3DES (RC2 is not supported)", ErrUnsupportedPFXEncryption, value.Algorithm.Algorithm)
	}
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() || len(value.Data) == 0 || len(value.Data)%block.BlockSize() != 0 {
		return nil, errors.New("invalid PFX ciphertext or IV length")
	}
	output := make([]byte, len(value.Data))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(output, value.Data)
	padding := int(output[len(output)-1])
	invalid := padding < 1 || padding > block.BlockSize() || padding > len(output)
	if !invalid {
		for _, b := range output[len(output)-padding:] {
			invalid = invalid || int(b) != padding
		}
	}
	if invalid {
		clear(output)
		return nil, errors.New("invalid PFX encryption padding")
	}
	return output[:len(output)-padding], nil
}

func (d *pfxDecoder) tripleDES(parameters []byte) (cipher.Block, []byte, error) {
	var params struct {
		Salt       []byte
		Iterations int
	}
	if err := decodeDER(parameters, &params); err != nil {
		return nil, nil, err
	}
	if params.Iterations < 1 || params.Iterations > 333_333 || len(params.Salt) > 1024 {
		return nil, nil, errors.New("invalid PFX 3DES parameters")
	}
	if err := d.useIterations(3 * params.Iterations); err != nil {
		return nil, nil, err
	}
	password := pfxPassword(d.password)
	defer clear(password)
	key := derivePFX(sha1.New, 1, password, params.Salt, params.Iterations, 24)
	defer clear(key)
	iv := derivePFX(sha1.New, 2, password, params.Salt, params.Iterations, 8)
	block, err := des.NewTripleDESCipher(key)
	return block, iv, err
}

func (d *pfxDecoder) pbes2(parameters []byte) (cipher.Block, []byte, error) {
	var params struct {
		KDF        pkix.AlgorithmIdentifier
		Encryption pkix.AlgorithmIdentifier
	}
	if err := decodeDER(parameters, &params); err != nil {
		return nil, nil, err
	}
	if !params.KDF.Algorithm.Equal(oidPBKDF2) {
		return nil, nil, errors.New("PFX PBES2 requires PBKDF2")
	}
	var kdf struct {
		Salt       []byte
		Iterations int
		KeyLength  int                      `asn1:"optional"`
		PRF        pkix.AlgorithmIdentifier `asn1:"optional"`
	}
	if err := decodeDER(params.KDF.Parameters.FullBytes, &kdf); err != nil {
		return nil, nil, err
	}
	if kdf.Iterations < 1 || kdf.Iterations > 1_000_000 {
		return nil, nil, errors.New("invalid PFX PBKDF2 iteration count")
	}
	if len(kdf.Salt) < 1 || len(kdf.Salt) > 1024 {
		return nil, nil, errors.New("invalid PFX PBKDF2 salt")
	}
	keySize := 0
	switch params.Encryption.Algorithm.String() {
	case "2.16.840.1.101.3.4.1.2":
		keySize = 16
	case "2.16.840.1.101.3.4.1.22":
		keySize = 24
	case "2.16.840.1.101.3.4.1.42":
		keySize = 32
	default:
		return nil, nil, errors.New("PFX PBES2 requires AES-CBC")
	}
	if kdf.KeyLength != 0 && kdf.KeyLength != keySize {
		return nil, nil, errors.New("PFX PBKDF2 key length differs from AES key size")
	}
	var newHash func() hash.Hash
	switch kdf.PRF.Algorithm.String() {
	case "", "1.2.840.113549.2.7":
		newHash = sha1.New
	case "1.2.840.113549.2.9":
		newHash = sha256.New
	case "1.2.840.113549.2.10":
		newHash = sha512.New384
	case "1.2.840.113549.2.11":
		newHash = sha512.New
	default:
		return nil, nil, errors.New("unsupported PFX PBKDF2 PRF")
	}
	if parameters := kdf.PRF.Parameters.FullBytes; len(parameters) != 0 && !bytes.Equal(parameters, []byte{5, 0}) {
		return nil, nil, errors.New("invalid PFX PBKDF2 PRF parameters")
	}
	blocks := (keySize + newHash().Size() - 1) / newHash().Size()
	if err := d.useIterations(kdf.Iterations * blocks); err != nil {
		return nil, nil, err
	}
	key, err := pbkdf2.Key(newHash, d.password, kdf.Salt, kdf.Iterations, keySize)
	if err != nil {
		return nil, nil, err
	}
	defer clear(key)
	var iv []byte
	if err := decodeDER(params.Encryption.Parameters.FullBytes, &iv); err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	return block, iv, err
}
