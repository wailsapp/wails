package hashes

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ---------------- Plugin Setup ----------------

type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (r *Plugin) Shutdown() {}

func (r *Plugin) Name() string {
	return "Hashes Plugin"
}

func (r *Plugin) Init(_ *application.App) error {
	return nil
}

func (r *Plugin) CallableByJS() []string {
	return []string{
		"Generate",
	}
}

func (r *Plugin) InjectJS() string {
	return ""
}

// ---------------- Plugin Methods ----------------

type Hashes struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

func (r *Plugin) Generate(s string) Hashes {
	md5Hash := md5.Sum([]byte(s))
	sha1Hash := sha1.Sum([]byte(s))
	sha256Hash := sha256.Sum256([]byte(s))

	return Hashes{
		MD5:    hex.EncodeToString(md5Hash[:]),
		SHA1:   hex.EncodeToString(sha1Hash[:]),
		SHA256: hex.EncodeToString(sha256Hash[:]),
	}
}
