package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
)

// SelfupdateKeygenOptions holds options for the selfupdate keygen command.
type SelfupdateKeygenOptions struct {
	Output     string `name:"o" description:"Output directory for generated keys" default:"."`
	KeyType    string `name:"type" description:"Key type: ecdsa or pgp" default:"ecdsa"`
	CommonName string `name:"cn" description:"Common name for the certificate" default:"Wails Update Signing Key"`
	ValidYears int    `name:"years" description:"Certificate validity in years" default:"10"`
	Force      bool   `name:"f" description:"Overwrite existing keys"`
}

// SelfupdateKeygen generates cryptographic keys for signing application updates.
func SelfupdateKeygen(options *SelfupdateKeygenOptions) error {
	pterm.DefaultHeader.Println("Wails Selfupdate Key Generator")

	switch options.KeyType {
	case "ecdsa":
		return generateECDSAKeys(options)
	case "pgp":
		return fmt.Errorf("PGP key generation not yet implemented - use gpg to generate keys")
	default:
		return fmt.Errorf("unsupported key type: %s (use 'ecdsa' or 'pgp')", options.KeyType)
	}
}

func generateECDSAKeys(options *SelfupdateKeygenOptions) error {
	// Create output directory if it doesn't exist.
	if err := os.MkdirAll(options.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	privateKeyPath := filepath.Join(options.Output, "update_private.pem")
	publicKeyPath := filepath.Join(options.Output, "update_public.pem")

	// Check if keys already exist.
	if !options.Force {
		if _, err := os.Stat(privateKeyPath); err == nil {
			return fmt.Errorf("private key already exists at %s (use -f to overwrite)", privateKeyPath)
		}
		if _, err := os.Stat(publicKeyPath); err == nil {
			return fmt.Errorf("public key already exists at %s (use -f to overwrite)", publicKeyPath)
		}
	}

	pterm.Info.Println("Generating ECDSA P-256 key pair...")

	// Generate ECDSA private key.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create a self-signed certificate.
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: options.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(options.ValidYears, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode private key to PEM.
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode certificate (public key) to PEM.
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Write private key (with restrictive permissions).
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}
	pterm.Success.Printf("Private key written to: %s\n", privateKeyPath)

	// Write public key/certificate.
	if err := os.WriteFile(publicKeyPath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}
	pterm.Success.Printf("Public key written to: %s\n", publicKeyPath)

	// Print usage instructions.
	pterm.Println()
	pterm.DefaultBox.WithTitle("Usage Instructions").Println(`
1. Keep the private key (update_private.pem) SECURE!
   - Store it in a secure location (e.g., CI/CD secrets)
   - NEVER commit it to version control

2. Embed the public key in your application:
   selfupdate.NewWithConfig(&selfupdate.Config{
       CurrentVersion: "1.0.0",
       Repository:     "owner/repo",
       Signature: &selfupdate.SignatureConfig{
           Type:      selfupdate.SignatureECDSA,
           PublicKey: string(publicKeyPEM),
       },
   })

3. Sign your releases with the private key:
   wails3 tool selfupdate sign -key update_private.pem -file myapp

4. Upload both the binary and .sig file to your release
`)

	return nil
}

// SelfupdateSignOptions holds options for the selfupdate sign command.
type SelfupdateSignOptions struct {
	KeyFile string `name:"key" description:"Path to the private key file" required:"true"`
	File    string `name:"file" description:"Path to the file to sign" required:"true"`
	Output  string `name:"o" description:"Output path for signature file (default: <file>.sig)"`
}

// SelfupdateSign signs a file with the given private key.
func SelfupdateSign(options *SelfupdateSignOptions) error {
	pterm.DefaultHeader.Println("Wails Selfupdate Signer")

	// Read the private key.
	privateKeyPEM, err := os.ReadFile(options.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block from private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Read the file to sign.
	data, err := os.ReadFile(options.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	pterm.Info.Printf("Signing %s (%d bytes)...\n", options.File, len(data))

	// Sign the file using ECDSA.
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hashData(data))
	if err != nil {
		return fmt.Errorf("failed to sign file: %w", err)
	}

	// Determine output path.
	outputPath := options.Output
	if outputPath == "" {
		outputPath = options.File + ".sig"
	}

	// Write the signature.
	if err := os.WriteFile(outputPath, signature, 0644); err != nil {
		return fmt.Errorf("failed to write signature: %w", err)
	}

	pterm.Success.Printf("Signature written to: %s\n", outputPath)

	return nil
}

// hashData computes SHA256 hash of data for ECDSA signing.
func hashData(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// SelfupdateVerifyOptions holds options for the selfupdate verify command.
type SelfupdateVerifyOptions struct {
	KeyFile string `name:"key" description:"Path to the public key file" required:"true"`
	File    string `name:"file" description:"Path to the file to verify" required:"true"`
	SigFile string `name:"sig" description:"Path to the signature file (default: <file>.sig)"`
}

// SelfupdateVerify verifies a file signature.
func SelfupdateVerify(options *SelfupdateVerifyOptions) error {
	pterm.DefaultHeader.Println("Wails Selfupdate Verifier")

	// Read the public key (certificate).
	publicKeyPEM, err := os.ReadFile(options.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block from public key")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	publicKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not ECDSA")
	}

	// Read the file to verify.
	data, err := os.ReadFile(options.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Read the signature.
	sigPath := options.SigFile
	if sigPath == "" {
		sigPath = options.File + ".sig"
	}

	signature, err := os.ReadFile(sigPath)
	if err != nil {
		return fmt.Errorf("failed to read signature: %w", err)
	}

	pterm.Info.Printf("Verifying %s...\n", options.File)

	// Verify the signature.
	if !ecdsa.VerifyASN1(publicKey, hashData(data), signature) {
		pterm.Error.Println("Signature verification FAILED!")
		return fmt.Errorf("signature verification failed")
	}

	pterm.Success.Println("Signature verification PASSED!")

	return nil
}
