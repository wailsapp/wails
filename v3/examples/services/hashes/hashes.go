package hashes

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Hashes = struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

type Service struct{}

func (h *Service) Generate(s string) Hashes {
	md5Hash := md5.Sum([]byte(s))
	sha1Hash := sha1.Sum([]byte(s))
	sha256Hash := sha256.Sum256([]byte(s))

	return Hashes{
		MD5:    hex.EncodeToString(md5Hash[:]),
		SHA1:   hex.EncodeToString(sha1Hash[:]),
		SHA256: hex.EncodeToString(sha256Hash[:]),
	}
}

func New() *Service {
	return &Service{}
}

func (h *Service) ServiceName() string {
	return "Hashes Service"
}

func (h *Service) ServiceStartup(context.Context, application.ServiceOptions) error {
	return nil
}

func (h *Service) ServiceShutdown() error { return nil }
