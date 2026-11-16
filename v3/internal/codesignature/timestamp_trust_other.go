//go:build !darwin

package codesignature

import (
	"crypto/x509"
	"time"
)

func verifyTimestampChain(certificate *x509.Certificate, certificates []*x509.Certificate, roots *x509.CertPool, at time.Time) error {
	return verifyTimestampGoChain(certificate, certificates, roots, at)
}
