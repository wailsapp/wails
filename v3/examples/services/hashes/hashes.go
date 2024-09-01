package hashes

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Hashes struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

func (h *Hashes) Generate(s string) Hashes {
	md5Hash := md5.Sum([]byte(s))
	sha1Hash := sha1.Sum([]byte(s))
	sha256Hash := sha256.Sum256([]byte(s))

	return Hashes{
		MD5:    hex.EncodeToString(md5Hash[:]),
		SHA1:   hex.EncodeToString(sha1Hash[:]),
		SHA256: hex.EncodeToString(sha256Hash[:]),
	}
}

func New() *Hashes {
	return &Hashes{}
}

func (h *Hashes) OnShutdown() error { return nil }

func (h *Hashes) Name() string {
	return "Hashes Service"
}

func (h *Hashes) OnStartup(_ context.Context, _ application.ServiceOptions) error {
	return nil
}
