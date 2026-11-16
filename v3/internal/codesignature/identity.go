package codesignature

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Identity struct {
	Signer      crypto.Signer
	Certificate *x509.Certificate
	Chain       []*x509.Certificate
}

type Options struct {
	// TimestampAttribute defaults to the standard RFC3161 unsigned attribute.
	TimestampAttribute asn1.ObjectIdentifier
	TimestampURL       string
	HTTPClient         *http.Client
	// TimestampRoots overrides system roots for a private timestamp authority.
	TimestampRoots *x509.CertPool
}

func identityFromParts(key crypto.Signer, certificates []*x509.Certificate) (Identity, error) {
	if key == nil {
		return Identity{}, errors.New("PFX does not contain a private signing key")
	}
	public, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return Identity{}, err
	}
	result := Identity{Signer: key}
	for _, cert := range certificates {
		if bytes.Equal(public, cert.RawSubjectPublicKeyInfo) {
			if result.Certificate != nil {
				return Identity{}, errors.New("PFX contains multiple matching signing certificates")
			}
			result.Certificate = cert
		} else {
			result.Chain = append(result.Chain, cert)
		}
	}
	if result.Certificate == nil {
		return Identity{}, errors.New("PFX certificate does not match the private key")
	}
	return result, validateIdentity(result)
}

func validateIdentity(identity Identity) error {
	if identity.Signer == nil || identity.Certificate == nil {
		return errors.New("signing requires a private key and certificate")
	}
	key, ok := identity.Signer.Public().(*rsa.PublicKey)
	if !ok || key.N.BitLen() < 2048 || key.N.BitLen() > 8192 {
		return errors.New("portable signing requires a 2048–8192 bit RSA key")
	}
	public, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	if !bytes.Equal(public, identity.Certificate.RawSubjectPublicKeyInfo) {
		return errors.New("signing certificate does not match the private key")
	}
	now := time.Now()
	if now.Before(identity.Certificate.NotBefore) || now.After(identity.Certificate.NotAfter) {
		return errors.New("signing certificate is outside its validity period")
	}
	if identity.Certificate.KeyUsage != 0 && identity.Certificate.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return errors.New("certificate does not permit digital signatures")
	}
	if len(identity.Certificate.ExtKeyUsage) > 0 || len(identity.Certificate.UnknownExtKeyUsage) > 0 {
		allowed := false
		for _, usage := range identity.Certificate.ExtKeyUsage {
			allowed = allowed || usage == x509.ExtKeyUsageCodeSigning || usage == x509.ExtKeyUsageAny
		}
		if !allowed {
			return errors.New("certificate does not permit code signing")
		}
	}
	if identity.Certificate.SerialNumber == nil || identity.Certificate.SerialNumber.Sign() < 0 {
		return fmt.Errorf("invalid signing certificate serial number")
	}
	return nil
}
