package codesignature

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

var oidTimestampInfo = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 1, 4}

type messageImprint struct {
	Algorithm pkix.AlgorithmIdentifier
	Digest    []byte
}
type timestampInfo struct {
	Version  int
	Policy   asn1.ObjectIdentifier
	Imprint  messageImprint
	Serial   *big.Int
	Time     time.Time `asn1:"generalized"`
	Accuracy struct {
		Seconds int `asn1:"optional"`
		Millis  int `asn1:"optional,tag:0"`
		Micros  int `asn1:"optional,tag:1"`
	} `asn1:"optional"`
	Ordering   bool          `asn1:"optional"`
	Nonce      *big.Int      `asn1:"optional"`
	TSA        asn1.RawValue `asn1:"optional,tag:0"`
	Extensions asn1.RawValue `asn1:"optional,tag:1"`
}

func requestTimestamp(ctx context.Context, signature []byte, options Options) ([]byte, error) {
	address, err := url.Parse(options.TimestampURL)
	if err != nil || address.Host == "" || address.User != nil || address.Fragment != "" || (address.Scheme != "http" && address.Scheme != "https") {
		return nil, errors.New("invalid timestamp server URL")
	}
	nonce, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256(signature)
	query := struct {
		Version      int
		Imprint      messageImprint
		Nonce        *big.Int
		Certificates bool
	}{1, messageImprint{algorithm(oidSHA256), digest[:]}, nonce, true}
	body, err := asn1.Marshal(query)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, address.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/timestamp-query")
	request.Header.Set("Accept", "application/timestamp-reply")
	client := options.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: time.Minute}
	}
	copyClient := *client
	copyClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	copyClient.Jar = nil
	response, err := copyClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("timestamp request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("timestamp server returned HTTP %d", response.StatusCode)
	}
	body, err = io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(body) > 4<<20 {
		return nil, errors.New("timestamp response exceeds 4 MiB")
	}
	var reply struct {
		Status struct {
			Code    int
			Text    []string       `asn1:"optional"`
			Failure asn1.BitString `asn1:"optional"`
		}
		Token asn1.RawValue `asn1:"optional"`
	}
	if err := decodeDER(body, &reply); err != nil {
		return nil, fmt.Errorf("invalid timestamp response: %w", err)
	}
	if (reply.Status.Code != 0 && reply.Status.Code != 1) || len(reply.Token.FullBytes) == 0 {
		return nil, fmt.Errorf("timestamp authority rejected request (status %d)", reply.Status.Code)
	}
	if err := verifyTimestamp(reply.Token.FullBytes, digest[:], nonce, options.TimestampRoots, time.Now()); err != nil {
		return nil, err
	}
	return reply.Token.FullBytes, nil
}

func verifyTimestamp(token, digest []byte, nonce *big.Int, roots *x509.CertPool, now time.Time) error {
	var outer contentInfo
	if err := decodeDER(token, &outer); err != nil {
		return err
	}
	if !outer.Type.Equal(oidSignedData) {
		return errors.New("timestamp token is not CMS SignedData")
	}
	var signed signedData
	if err := decodeDER(outer.Content.Bytes, &signed); err != nil {
		return err
	}
	if len(signed.Signers) != 1 || len(signed.Digests) != 1 || !signed.Content.Type.Equal(oidTimestampInfo) {
		return errors.New("invalid timestamp CMS content or signer count")
	}
	var content []byte
	if err := decodeDER(signed.Content.Content.Bytes, &content); err != nil {
		return err
	}
	var info timestampInfo
	if err := decodeDER(content, &info); err != nil {
		return err
	}
	if nonce == nil || info.Version != 1 || !info.Imprint.Algorithm.Algorithm.Equal(oidSHA256) || !bytes.Equal(info.Imprint.Digest, digest) || info.Nonce == nil || nonce.Cmp(info.Nonce) != 0 {
		return errors.New("timestamp nonce or message imprint does not match the request")
	}
	if info.Time.Before(now.Add(-5*time.Minute)) || info.Time.After(now.Add(5*time.Minute)) {
		return errors.New("timestamp generation time is not current")
	}
	signer := signed.Signers[0]
	if signer.Signer.Serial == nil || signer.Version != 1 || !signer.Digest.Algorithm.Equal(signed.Digests[0].Algorithm) {
		return errors.New("invalid timestamp signer algorithm or version")
	}
	certificates, err := cmsCertificates(signed.Certificates.Bytes)
	if err != nil {
		return err
	}
	var certificate *x509.Certificate
	for _, candidate := range certificates {
		if bytes.Equal(candidate.RawIssuer, signer.Signer.Issuer.FullBytes) && candidate.SerialNumber.Cmp(signer.Signer.Serial) == 0 {
			certificate = candidate
		}
	}
	if certificate == nil {
		return errors.New("timestamp signer certificate is missing")
	}
	if len(certificate.ExtKeyUsage) != 1 || certificate.ExtKeyUsage[0] != x509.ExtKeyUsageTimeStamping || len(certificate.UnknownExtKeyUsage) != 0 {
		return errors.New("timestamp certificate must have only the timestamping extended key usage")
	}
	critical := false
	for _, extension := range certificate.Extensions {
		if extension.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 15}) && certificate.KeyUsage&(x509.KeyUsageDigitalSignature|x509.KeyUsageContentCommitment) == 0 {
			return errors.New("timestamp certificate key usage does not permit signatures")
		}
		if extension.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 37}) {
			critical = extension.Critical
		}
	}
	if !critical {
		return errors.New("timestamp certificate extended key usage is not critical")
	}
	if err := verifyTimestampChain(certificate, certificates, roots, info.Time); err != nil {
		return fmt.Errorf("untrusted timestamp certificate: %w", err)
	}
	if err := verifyTimestampAttributes(signer, certificate, content); err != nil {
		return err
	}
	attrs, err := asn1.MarshalWithParams(signer.Signed, "set")
	if err != nil {
		return err
	}
	signatureAlgorithm, err := cmsSignatureAlgorithm(signer.Algorithm.Algorithm, signer.Digest.Algorithm)
	if err != nil {
		return err
	}
	if err := certificate.CheckSignature(signatureAlgorithm, attrs, signer.Signature); err != nil {
		return fmt.Errorf("invalid timestamp CMS signature: %w", err)
	}
	return nil
}

func cmsCertificates(data []byte) ([]*x509.Certificate, error) {
	var result []*x509.Certificate
	for len(data) > 0 {
		if len(result) >= 64 {
			return nil, errors.New("too many CMS certificates")
		}
		var value asn1.RawValue
		rest, err := asn1.Unmarshal(data, &value)
		if err != nil {
			return nil, err
		}
		certificate, err := x509.ParseCertificate(value.FullBytes)
		if err != nil {
			return nil, err
		}
		result = append(result, certificate)
		data = rest
	}
	return result, nil
}

func cmsSignatureAlgorithm(signature, digest asn1.ObjectIdentifier) (x509.SignatureAlgorithm, error) {
	switch digest.String() {
	case "2.16.840.1.101.3.4.2.1":
		if signature.Equal(oidRSA) || signature.String() == "1.2.840.113549.1.1.11" {
			return x509.SHA256WithRSA, nil
		}
		if signature.String() == "1.2.840.10045.4.3.2" {
			return x509.ECDSAWithSHA256, nil
		}
	case "2.16.840.1.101.3.4.2.2":
		if signature.Equal(oidRSA) || signature.String() == "1.2.840.113549.1.1.12" {
			return x509.SHA384WithRSA, nil
		}
		if signature.String() == "1.2.840.10045.4.3.3" {
			return x509.ECDSAWithSHA384, nil
		}
	case "2.16.840.1.101.3.4.2.3":
		if signature.Equal(oidRSA) || signature.String() == "1.2.840.113549.1.1.13" {
			return x509.SHA512WithRSA, nil
		}
		if signature.String() == "1.2.840.10045.4.3.4" {
			return x509.ECDSAWithSHA512, nil
		}
	}
	return x509.UnknownSignatureAlgorithm, errors.New("unsupported timestamp signature algorithm")
}

func verifyTimestampAttributes(signer signerInfo, certificate *x509.Certificate, content []byte) error {
	newHash, err := digestHash(signer.Digest.Algorithm)
	if err != nil {
		return err
	}
	hash := newHash()
	_, _ = hash.Write(content)
	seen := map[string]bool{}
	for _, attr := range signer.Signed {
		key := attr.Type.String()
		if seen[key] || len(attr.Values) != 1 {
			return errors.New("duplicate or malformed timestamp signed attribute")
		}
		seen[key] = true
		data := attr.Values[0].FullBytes
		switch key {
		case "1.2.840.113549.1.9.3":
			var oid asn1.ObjectIdentifier
			if err := decodeDER(data, &oid); err != nil || !oid.Equal(oidTimestampInfo) {
				return errors.New("incorrect timestamp signed content type")
			}
		case "1.2.840.113549.1.9.4":
			var value []byte
			if err := decodeDER(data, &value); err != nil || !bytes.Equal(value, hash.Sum(nil)) {
				return errors.New("timestamp CMS digest does not match content")
			}
		case "1.2.840.113549.1.9.16.2.12", "1.2.840.113549.1.9.16.2.47":
			if err := verifyESSCertificate(data, key, certificate); err != nil {
				return err
			}
		}
	}
	if !seen[oidContentType.String()] || !seen[oidMessageDigest.String()] || (!seen["1.2.840.113549.1.9.16.2.12"] && !seen["1.2.840.113549.1.9.16.2.47"]) {
		return errors.New("timestamp CMS is missing required signed attributes")
	}
	return nil
}

func verifyESSCertificate(data []byte, kind string, certificate *x509.Certificate) error {
	var entries struct {
		Certificates []asn1.RawValue
		Policies     []asn1.RawValue `asn1:"optional"`
	}
	if err := decodeDER(data, &entries); err != nil {
		return err
	}
	if len(entries.Certificates) == 0 {
		return errors.New("timestamp certificate identifier is missing")
	}
	var digest []byte
	hashOID := oidSHA1
	var issuer asn1.RawValue
	if kind == "1.2.840.113549.1.9.16.2.47" {
		var value struct {
			Algorithm pkix.AlgorithmIdentifier `asn1:"optional"`
			Hash      []byte
			Issuer    asn1.RawValue `asn1:"optional"`
		}
		if err := decodeDER(entries.Certificates[0].FullBytes, &value); err != nil {
			return err
		}
		hashOID = value.Algorithm.Algorithm
		if len(hashOID) == 0 {
			hashOID = oidSHA256
		}
		digest = value.Hash
		issuer = value.Issuer
	} else {
		var value struct {
			Hash   []byte
			Issuer asn1.RawValue `asn1:"optional"`
		}
		if err := decodeDER(entries.Certificates[0].FullBytes, &value); err != nil {
			return err
		}
		digest = value.Hash
		issuer = value.Issuer
	}
	newHash, err := digestHash(hashOID)
	if err != nil {
		return err
	}
	hash := newHash()
	_, _ = hash.Write(certificate.Raw)
	if !bytes.Equal(digest, hash.Sum(nil)) {
		return errors.New("timestamp certificate identifier does not match signer")
	}
	if len(issuer.FullBytes) > 0 {
		var value struct {
			Names  []asn1.RawValue
			Serial *big.Int
		}
		if err := decodeDER(issuer.FullBytes, &value); err != nil {
			return err
		}
		if value.Serial == nil || value.Serial.Cmp(certificate.SerialNumber) != 0 {
			return errors.New("timestamp certificate serial number mismatch")
		}
		matches := false
		for _, name := range value.Names {
			matches = matches || (name.Class == 2 && name.Tag == 4 && bytes.Equal(name.Bytes, certificate.RawIssuer))
		}
		if !matches {
			return errors.New("timestamp certificate issuer mismatch")
		}
	}
	return nil
}
