package notaryservice

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// MaxUploadBytes is the single-PUT limit of this client. Larger uploads require
// multipart support and are rejected before creating an Apple submission.
const MaxUploadBytes int64 = 5 << 30

var bucketPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$`)
var regionPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)+$`)

type uploadCredentials struct {
	AccessKey    string `json:"awsAccessKeyId"`
	SecretKey    string `json:"awsSecretAccessKey"`
	SessionToken string `json:"awsSessionToken"`
	Bucket       string `json:"bucket"`
	Object       string `json:"object"`
}

// Submit streams a regular ZIP, DMG, or PKG file to the upload location provided
// by Apple. It does not retry POST or PUT automatically. A nonempty returned ID
// means Apple created the submission even when an upload error is returned.
func (c *Client) Submit(ctx context.Context, path string) (Submission, error) {
	file, err := os.Open(path)
	if err != nil {
		return Submission{}, fmt.Errorf("open notarization archive: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return Submission{}, fmt.Errorf("stat notarization archive: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		return Submission{}, errors.New("notarization archive must be a nonempty regular file")
	}
	if info.Size() > MaxUploadBytes {
		return Submission{}, errors.New("portable notarization uploads are limited to 5 GiB; multipart upload is not supported")
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".zip", ".dmg", ".pkg":
	default:
		return Submission{}, errors.New("notarization archive must have a .zip, .dmg, or .pkg extension")
	}
	digest, err := hashArchive(ctx, file)
	if err != nil {
		return Submission{}, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return Submission{}, err
	}
	var result struct {
		Data struct {
			ID         string            `json:"id"`
			Attributes uploadCredentials `json:"attributes"`
		} `json:"data"`
	}
	body := map[string]string{"submissionName": filepath.Base(path), "sha256": digest}
	if err := c.api(ctx, http.MethodPost, "", body, &result); err != nil {
		return Submission{}, err
	}
	submission := Submission{ID: result.Data.ID, Name: filepath.Base(path)}
	if _, err := submissionPath(submission.ID); err != nil {
		return submission, errors.New("notary API returned an invalid submission ID")
	}
	credentials := result.Data.Attributes
	if !bucketPattern.MatchString(credentials.Bucket) || credentials.Object == "" || credentials.AccessKey == "" || credentials.SecretKey == "" || credentials.SessionToken == "" {
		return submission, errors.New("notary API returned incomplete or invalid upload credentials")
	}
	if err := c.upload(ctx, file, info.Size(), digest, credentials); err != nil {
		return submission, fmt.Errorf("upload notarization submission %s: %w", submission.ID, err)
	}
	return submission, nil
}

func hashArchive(ctx context.Context, file *os.File) (string, error) {
	hash := sha256.New()
	buffer := make([]byte, 128*1024)
	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		n, err := file.Read(buffer)
		if n > 0 {
			_, _ = hash.Write(buffer[:n])
		}
		if errors.Is(err, io.EOF) {
			return hex.EncodeToString(hash.Sum(nil)), nil
		}
		if err != nil {
			return "", fmt.Errorf("hash notarization archive: %w", err)
		}
	}
}

func (c *Client) bucketRegion(ctx context.Context, bucket string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://"+bucket+".s3.amazonaws.com/", nil)
	if err != nil {
		return "", err
	}
	response, err := c.do(request, "S3 bucket region")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	// HeadBucket documents the region header even on an anonymous 403 response.
	region := response.Header.Get("X-Amz-Bucket-Region")
	if len(region) > 32 || !regionPattern.MatchString(region) {
		return "", errors.New("S3 did not return a valid bucket region")
	}
	return region, nil
}

func (c *Client) upload(ctx context.Context, file *os.File, size int64, digest string, credentials uploadCredentials) error {
	region, err := c.bucketRegion(ctx, credentials.Bucket)
	if err != nil {
		return err
	}
	path := "/" + credentials.Object
	address := &url.URL{Scheme: "https", Host: credentials.Bucket + ".s3-accelerate.amazonaws.com", Path: path, RawPath: escapeS3Path(path)}
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, address.String(), file)
	if err != nil {
		return errors.New("invalid S3 upload request")
	}
	request.ContentLength = size
	request.Header.Set("Content-Type", "application/octet-stream")
	signUpload(request, credentials, region, digest, time.Now())
	_, err = c.response(request, "S3 upload", maxResponseBytes)
	return err
}

// signUpload implements AWS's documented Signature Version 4 canonical request
// format. S3 object paths must not be cleaned or have repeated slashes removed.
func signUpload(request *http.Request, credentials uploadCredentials, region, digest string, now time.Time) {
	now = now.UTC()
	timestamp, date := now.Format("20060102T150405Z"), now.Format("20060102")
	request.Header.Set("X-Amz-Date", timestamp)
	request.Header.Set("X-Amz-Content-Sha256", digest)
	if credentials.SessionToken != "" {
		request.Header.Set("X-Amz-Security-Token", credentials.SessionToken)
	}
	headers := map[string]string{"host": request.URL.Host}
	for key, values := range request.Header {
		if strings.EqualFold(key, "Authorization") {
			continue
		}
		headers[strings.ToLower(key)] = strings.Join(strings.Fields(strings.Join(values, ",")), " ")
	}
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	sort.Strings(names)
	var canonicalHeaders strings.Builder
	for _, name := range names {
		canonicalHeaders.WriteString(name + ":" + headers[name] + "\n")
	}
	signedHeaders := strings.Join(names, ";")
	canonical := strings.Join([]string{request.Method, escapeS3Path(request.URL.Path), request.URL.RawQuery, canonicalHeaders.String(), signedHeaders, digest}, "\n")
	scope := date + "/" + region + "/s3/aws4_request"
	canonicalHash := sha256.Sum256([]byte(canonical))
	message := "AWS4-HMAC-SHA256\n" + timestamp + "\n" + scope + "\n" + hex.EncodeToString(canonicalHash[:])
	key := []byte("AWS4" + credentials.SecretKey)
	for _, part := range []string{date, region, "s3", "aws4_request"} {
		key = macSHA256(key, part)
	}
	signature := hex.EncodeToString(macSHA256(key, message))
	request.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+credentials.AccessKey+"/"+scope+",SignedHeaders="+signedHeaders+",Signature="+signature)
}

func macSHA256(key []byte, value string) []byte {
	hash := hmac.New(sha256.New, key)
	_, _ = hash.Write([]byte(value))
	return hash.Sum(nil)
}

func escapeS3Path(path string) string {
	const digits = "0123456789ABCDEF"
	var result strings.Builder
	for i := 0; i < len(path); i++ {
		value := path[i]
		if value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z' || value >= '0' && value <= '9' || strings.ContainsRune("-._~/", rune(value)) {
			result.WriteByte(value)
			continue
		}
		result.WriteByte('%')
		result.WriteByte(digits[value>>4])
		result.WriteByte(digits[value&15])
	}
	return result.String()
}
