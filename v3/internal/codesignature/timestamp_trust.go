package codesignature

import (
	"bytes"
	"crypto/x509"
	"time"
)

func verifyTimestampGoChain(certificate *x509.Certificate, certificates []*x509.Certificate, roots *x509.CertPool, at time.Time) error {
	intermediates := x509.NewCertPool()
	for _, candidate := range certificates {
		if !bytes.Equal(candidate.Raw, certificate.Raw) {
			intermediates.AddCert(candidate)
		}
	}
	_, err := certificate.Verify(x509.VerifyOptions{Roots: roots, Intermediates: intermediates, CurrentTime: at, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageTimeStamping}})
	return err
}
