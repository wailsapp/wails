package packagesign

import (
	"bytes"
	"context"
	"crypto/md5"  // Required by the dpkg-sig v4 archive format.
	"crypto/sha1" // Required by the dpkg-sig v4 archive format.
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
)

type debMember struct {
	name string
	size int64
	md5  []byte
	sha1 []byte
}

func debMembers(ctx context.Context, src io.Reader) ([]debMember, error) {
	var magic [8]byte
	if _, err := io.ReadFull(src, magic[:]); err != nil || string(magic[:]) != "!<arch>\n" {
		return nil, fmt.Errorf("invalid DEB ar archive")
	}
	var members []debMember
	seen := make(map[string]bool)
	for {
		var header [60]byte
		_, err := io.ReadFull(src, header[:])
		if err == io.EOF {
			break
		}
		if err != nil || string(header[58:]) != "`\n" {
			return nil, fmt.Errorf("invalid DEB member header")
		}
		name := strings.TrimSuffix(strings.TrimSpace(string(header[:16])), "/")
		size, err := strconv.ParseInt(strings.TrimSpace(string(header[48:58])), 10, 64)
		if err != nil || size < 0 || name == "" || strings.ContainsAny(name, " \t\r\n/\\") || seen[name] {
			return nil, fmt.Errorf("invalid or duplicate DEB member %q", name)
		}
		seen[name] = true
		md, sha := md5.New(), sha1.New()
		if _, err := io.CopyN(io.MultiWriter(md, sha), contextReader{ctx, src}, size); err != nil {
			return nil, fmt.Errorf("read DEB member %s: %w", name, err)
		}
		if size%2 != 0 {
			var padding [1]byte
			if _, err := io.ReadFull(src, padding[:]); err != nil || padding[0] != '\n' {
				return nil, fmt.Errorf("invalid DEB member padding")
			}
		}
		members = append(members, debMember{name, size, md.Sum(nil), sha.Sum(nil)})
		if len(members) > 4096 {
			return nil, fmt.Errorf("too many DEB archive members")
		}
	}
	if len(members) < 3 || members[0].name != "debian-binary" || members[0].size != 4 ||
		!strings.HasPrefix(members[1].name, "control.tar") || !strings.HasPrefix(members[2].name, "data.tar") {
		return nil, fmt.Errorf("DEB archive is missing its package members")
	}
	return members, nil
}

func signDEB(ctx context.Context, src *os.File, output string, signer *signingKey, role string) error {
	if role == "" {
		role = "builder"
	}
	if len(role) > 10 || strings.Trim(role, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-") != "" {
		return fmt.Errorf("DEB signing role must contain 1 to 10 letters, digits, underscores or hyphens")
	}
	original := sha256.New()
	members, err := debMembers(ctx, io.TeeReader(src, original))
	if err != nil {
		return err
	}
	name, err := debSignatureName(members, role)
	if err != nil {
		return err
	}
	identity := fmt.Sprintf("%X", signer.key.Fingerprint)
	if _, id := signer.entity.PrimarySelfSignature(); id != nil {
		identity = strings.NewReplacer("\r", " ", "\n", " ").Replace(id.Name)
	}
	var signature bytes.Buffer
	cleartext, err := clearsign.Encode(&signature, signer.key, signer.config)
	if err != nil {
		return err
	}
	fmt.Fprintf(cleartext, "Version: 4\nSigner: %s\nDate: %s\nRole: %s\nFiles:\n", identity, signer.config.Now().Format(time.RFC1123Z), role)
	for _, member := range members {
		fmt.Fprintf(cleartext, " %x %x %d %s\n", member.md5, member.sha1, member.size, member.name)
	}
	if err := cleartext.Close(); err != nil {
		return fmt.Errorf("sign DEB manifest: %w", err)
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return writeOutput(output, func(dst io.Writer) error {
		copied := sha256.New()
		if _, err := io.Copy(io.MultiWriter(dst, copied), contextReader{ctx, src}); err != nil {
			return err
		}
		if !bytes.Equal(original.Sum(nil), copied.Sum(nil)) {
			return fmt.Errorf("DEB input changed while signing")
		}
		header := fmt.Sprintf("%-16s%-12d%-6d%-6d%-8o%-10d`\n", name, signer.config.Now().Unix(), 0, 0, 0644, signature.Len())
		if _, err := io.WriteString(dst, header); err != nil {
			return err
		}
		if _, err := dst.Write(signature.Bytes()); err != nil {
			return err
		}
		if signature.Len()%2 != 0 {
			_, err = io.WriteString(dst, "\n")
		}
		return err
	})
}

func debSignatureName(members []debMember, role string) (string, error) {
	used := make(map[string]bool, len(members))
	for _, member := range members {
		used[member.name] = true
	}
	name := "_gpg" + role
	for n := 0; n < 100; n++ {
		candidate := name
		if n != 0 {
			candidate += strconv.Itoa(n)
		}
		if !used[candidate] {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("too many DEB signatures for role %q", role)
}
