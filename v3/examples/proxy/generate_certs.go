// +build ignore

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to get source file location")
	}
	scriptDir := filepath.Dir(filename)

	// Create certs directory relative to script location
	certsDir := filepath.Join(scriptDir, "certs")
	fmt.Printf("Creating certificates in: %s\n", certsDir)

	if err := os.MkdirAll(certsDir, 0755); err != nil {
		log.Fatalf("Failed to create certs directory: %v", err)
	}

	// Generate CA certificate
	fmt.Println("Generating CA certificate...")
	caCert, caKey, err := generateCA()
	if err != nil {
		log.Fatalf("Failed to generate CA: %v", err)
	}

	// Save CA certificate and key
	if err := saveCertAndKey(filepath.Join(certsDir, "ca.crt"), filepath.Join(certsDir, "ca.key"), caCert, caKey); err != nil {
		log.Fatalf("Failed to save CA: %v", err)
	}

	// Generate server certificate for app-local.wails-awesome.io
	fmt.Println("Generating server certificate for app-local.wails-awesome.io...")
	serverCert, serverKey, err := generateServerCert(caCert, caKey, []string{
		"app-local.wails-awesome.io",
		"*.app-local.wails-awesome.io",
		"localhost",
		"127.0.0.1",
	})
	if err != nil {
		log.Fatalf("Failed to generate server certificate: %v", err)
	}

	// Save server certificate and key
	if err := saveCertAndKey(filepath.Join(certsDir, "server.crt"), filepath.Join(certsDir, "server.key"), serverCert, serverKey); err != nil {
		log.Fatalf("Failed to save server certificate: %v", err)
	}

	fmt.Println("\n✅ Certificates generated successfully!")
	fmt.Println("\n📁 Certificate files created:")
	fmt.Printf("   - CA Certificate: %s\n", filepath.Join(certsDir, "ca.crt"))
	fmt.Printf("   - CA Private Key: %s\n", filepath.Join(certsDir, "ca.key"))
	fmt.Printf("   - Server Certificate: %s\n", filepath.Join(certsDir, "server.crt"))
	fmt.Printf("   - Server Private Key: %s\n", filepath.Join(certsDir, "server.key"))

	fmt.Println("\n🔒 To trust the CA certificate:")
	fmt.Println("   Windows:")
	fmt.Println("     1. Double-click certs/ca.crt")
	fmt.Println("     2. Click 'Install Certificate'")
	fmt.Println("     3. Select 'Local Machine'")
	fmt.Println("     4. Select 'Place all certificates in the following store'")
	fmt.Println("     5. Browse and select 'Trusted Root Certification Authorities'")
	fmt.Println("     6. Finish the wizard")
	fmt.Println("")
	fmt.Println("   macOS:")
	fmt.Printf("     sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s\n", filepath.Join(certsDir, "ca.crt"))
	fmt.Println("")
	fmt.Println("   Linux:")
	fmt.Printf("     sudo cp %s /usr/local/share/ca-certificates/wails-proxy-example.crt\n", filepath.Join(certsDir, "ca.crt"))
	fmt.Println("     sudo update-ca-certificates")
	fmt.Println("")
	fmt.Println("📌 Add to hosts file (as administrator/root):")
	fmt.Println("   127.0.0.1    app-local.wails-awesome.io")
}

func generateCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	// Generate RSA key
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate CA key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Wails Proxy Example"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Generate certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	return cert, caKey, nil
}

func generateServerCert(caCert *x509.Certificate, caKey *rsa.PrivateKey, hosts []string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Generate RSA key for server
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate server key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Wails Proxy Example Server"},
			Country:      []string{"US"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}

	// Add hosts and IPs
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	// Generate certificate signed by CA
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create server certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse server certificate: %w", err)
	}

	return cert, serverKey, nil
}

func saveCertAndKey(certPath, keyPath string, cert *x509.Certificate, key *rsa.PrivateKey) error {
	// Save certificate
	certOut, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %w", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyOut.Close()

	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	}); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Verify we can load the certificate
	_, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("failed to verify certificate and key pair: %w", err)
	}

	return nil
}