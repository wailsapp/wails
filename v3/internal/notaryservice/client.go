// Package notaryservice submits signed macOS archives to Apple's Notary API.
// It does not sign software or staple notarization tickets.
package notaryservice

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiURL = "https://appstoreconnect.apple.com/notary/v2/submissions"
const maxResponseBytes = 1 << 20
const maxLogBytes = 16 << 20

// Credentials identifies an App Store Connect team API key. PrivateKeyPEM is
// the contents of its downloaded .p8 file, never an Apple ID password.
type Credentials struct {
	KeyID         string
	IssuerID      string
	PrivateKeyPEM []byte
}

// Client is safe for concurrent requests. The provided HTTP transport must be
// safe for concurrent use. Redirects are disabled for every request.
type Client struct {
	http            *http.Client
	key             *ecdsa.PrivateKey
	keyID, issuerID string
}

func New(credentials Credentials, client *http.Client) (*Client, error) {
	if strings.TrimSpace(credentials.KeyID) == "" || strings.TrimSpace(credentials.IssuerID) == "" {
		return nil, errors.New("notary API requires key ID and team issuer ID")
	}
	block, rest := pem.Decode(credentials.PrivateKeyPEM)
	if block == nil || block.Type != "PRIVATE KEY" || len(bytes.TrimSpace(rest)) != 0 {
		return nil, errors.New("notary API requires one PKCS#8 private key PEM block")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("cannot parse notary API PKCS#8 private key")
	}
	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok || key.Curve != elliptic.P256() {
		return nil, errors.New("notary API requires an ECDSA P-256 private key")
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Minute}
	}
	copyClient := *client
	copyClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	copyClient.Jar = nil
	return &Client{http: &copyClient, key: key, keyID: credentials.KeyID, issuerID: credentials.IssuerID}, nil
}

// Submission contains no upload credentials. Submit returns its ID even if the
// subsequent upload fails, so callers can retain it for diagnostics.
type Submission struct {
	ID          string
	Status      string
	Name        string
	CreatedDate string
}

// HTTPError identifies an unsuccessful remote response without exposing its
// body, which might contain credentials or a signed download URL.
type HTTPError struct {
	Operation  string
	StatusCode int
}

func (e *HTTPError) Error() string { return fmt.Sprintf("%s: HTTP %d", e.Operation, e.StatusCode) }

func (c *Client) token(now time.Time) (string, error) {
	header, err := json.Marshal(map[string]string{"alg": "ES256", "kid": c.keyID, "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(map[string]any{"iss": c.issuerID, "iat": now.Unix(), "exp": now.Add(5 * time.Minute).Unix(), "aud": "appstoreconnect-v1"})
	if err != nil {
		return "", err
	}
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	digest := sha256.Sum256([]byte(unsigned))
	r, s, err := ecdsa.Sign(rand.Reader, c.key, digest[:])
	if err != nil {
		return "", fmt.Errorf("sign notary API token: %w", err)
	}
	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func (c *Client) api(ctx context.Context, method, suffix string, body any, result any) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	request, err := http.NewRequestWithContext(ctx, method, apiURL+suffix, bytes.NewReader(data))
	if err != nil {
		return err
	}
	token, err := c.token(time.Now())
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	data, err = c.response(request, "notary API", maxResponseBytes)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, result); err != nil {
		return errors.New("invalid JSON from notary API")
	}
	return nil
}

func (c *Client) do(request *http.Request, operation string) (*http.Response, error) {
	response, err := c.http.Do(request)
	if err != nil {
		var urlErr *url.Error
		for errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return nil, fmt.Errorf("%s request failed: %w", operation, err)
	}
	return response, nil
}

func (c *Client) response(request *http.Request, operation string, limit int64) ([]byte, error) {
	response, err := c.do(request, operation)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		httpError := &HTTPError{Operation: operation, StatusCode: response.StatusCode}
		if operation == "notary API" && response.StatusCode == http.StatusForbidden {
			data, readError := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
			var failure struct {
				Errors []struct {
					Code string `json:"code"`
				} `json:"errors"`
			}
			if readError == nil && len(data) <= maxResponseBytes && json.Unmarshal(data, &failure) == nil {
				for _, item := range failure.Errors {
					if item.Code == "FORBIDDEN.REQUIRED_AGREEMENTS_MISSING_OR_EXPIRED" {
						return nil, fmt.Errorf("%w: required Apple agreements are missing or expired; the Account Holder must review the agreements at https://developer.apple.com/account/ before retrying", httpError)
					}
				}
			}
		}
		return nil, httpError
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", operation, err)
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("%s response exceeds %d bytes", operation, limit)
	}
	return data, nil
}

func submissionPath(id string) (string, error) {
	if len(id) != 36 {
		return "", errors.New("invalid notarization submission ID")
	}
	for i, b := range []byte(id) {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if b == '-' {
				continue
			}
		} else if b >= '0' && b <= '9' || b >= 'a' && b <= 'f' || b >= 'A' && b <= 'F' {
			continue
		}
		return "", errors.New("invalid notarization submission ID")
	}
	return "/" + id, nil
}

func (c *Client) Status(ctx context.Context, id string) (Submission, error) {
	path, err := submissionPath(id)
	if err != nil {
		return Submission{}, err
	}
	var result struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Status      string `json:"status"`
				Name        string `json:"name"`
				CreatedDate string `json:"createdDate"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := c.api(ctx, http.MethodGet, path, nil, &result); err != nil {
		return Submission{}, err
	}
	if !strings.EqualFold(result.Data.ID, id) || result.Data.Attributes.Status == "" {
		return Submission{}, errors.New("notary API returned an incomplete or mismatched submission")
	}
	attributes := result.Data.Attributes
	return Submission{ID: result.Data.ID, Status: attributes.Status, Name: attributes.Name, CreatedDate: attributes.CreatedDate}, nil
}

// Wait polls until Accepted, Invalid, or Rejected. Rejections are returned as
// results, so the caller can retrieve logs before reporting failure. The context
// bounds the overall wait. A nonpositive interval selects ten seconds.
func (c *Client) Wait(ctx context.Context, id string, interval time.Duration) (Submission, error) {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	for {
		result, err := c.Status(ctx, id)
		if err != nil {
			return result, err
		}
		switch result.Status {
		case "Accepted", "Invalid", "Rejected":
			return result, nil
		case "In Progress":
		default:
			return result, fmt.Errorf("notary API returned an unknown submission status %q", result.Status)
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return result, ctx.Err()
		case <-timer.C:
		}
	}
}

// Log retrieves the JSON report without forwarding Apple authorization to the
// temporary download URL. The returned report may contain warnings on success.
func (c *Client) Log(ctx context.Context, id string) ([]byte, error) {
	path, err := submissionPath(id)
	if err != nil {
		return nil, err
	}
	var result struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				URL string `json:"developerLogUrl"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := c.api(ctx, http.MethodGet, path+"/logs", nil, &result); err != nil {
		return nil, err
	}
	if !strings.EqualFold(result.Data.ID, id) {
		return nil, errors.New("notary API returned mismatched log submission")
	}
	address, err := url.Parse(result.Data.Attributes.URL)
	if err != nil || address.Scheme != "https" || address.Host == "" || address.User != nil || address.Fragment != "" {
		return nil, errors.New("notary API returned an invalid HTTPS log URL")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address.String(), nil)
	if err != nil {
		return nil, errors.New("invalid notarization log request")
	}
	data, err := c.response(request, "notarization log", maxLogBytes)
	if err != nil {
		return nil, err
	}
	if !json.Valid(data) {
		return nil, errors.New("notarization log is not valid JSON")
	}
	return data, nil
}
