package buildmeta

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Revision reads SHA-1 file references without launching Git or loading its object database.
func Revision(ctx context.Context, root string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	directory, err := repositoryDirectory(root)
	if err != nil {
		return "", err
	}
	gitDirectory := filepath.Join(directory, ".git")
	marker, err := os.Stat(gitDirectory)
	if err != nil {
		return "", err
	}
	if !marker.IsDir() {
		data, err := os.ReadFile(gitDirectory)
		if err != nil {
			return "", err
		}
		value, found := strings.CutPrefix(strings.TrimSpace(string(data)), "gitdir: ")
		if !found || value == "" {
			return "", fmt.Errorf("invalid application Git directory marker")
		}
		gitDirectory = resolveRepositoryPath(directory, value)
	}
	commonDirectory := gitDirectory
	data, err := os.ReadFile(filepath.Join(gitDirectory, "commondir"))
	if err == nil {
		commonDirectory = resolveRepositoryPath(gitDirectory, strings.TrimSpace(string(data)))
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := checkRepositoryFormat(commonDirectory); err != nil {
		return "", err
	}
	revision, err := readReference(ctx, gitDirectory, commonDirectory, "HEAD", map[string]bool{})
	if err != nil {
		return "", fmt.Errorf("resolve application Git HEAD: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return revision, nil
}

func resolveRepositoryPath(base, value string) string {
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Join(base, value)
}

func checkRepositoryFormat(directory string) error {
	file, err := os.Open(filepath.Join(directory, "config"))
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	extensions := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") {
			section, _, closed := strings.Cut(line[1:], "]")
			if !closed {
				return fmt.Errorf("invalid application Git config section")
			}
			extensions = strings.EqualFold(strings.TrimSpace(section), "extensions")
			continue
		}
		if !extensions {
			continue
		}
		key, value, _ := strings.Cut(line, "=")
		key = strings.ToLower(strings.TrimSpace(key))
		if key != "objectformat" && key != "refstorage" {
			continue
		}
		value, err = repositoryConfigValue(value)
		if err != nil {
			return err
		}
		supported := key == "objectformat" && value == "sha1" || key == "refstorage" && value == "files"
		if !supported {
			return fmt.Errorf("application Git %s %q is not supported by revision lookup", key, value)
		}
	}
	return scanner.Err()
}

// Delimiters introduce comments only outside quotes; quoted whitespace is data.
func repositoryConfigValue(raw string) (string, error) {
	var value strings.Builder
	quoted, escaped := false, false
	end := 0
	for _, c := range strings.TrimLeft(raw, " \t") {
		if escaped {
			switch c {
			case 'n':
				c = '\n'
			case 't':
				c = '\t'
			case 'b':
				c = '\b'
			case '\\', '"':
			default:
				return "", fmt.Errorf("invalid application Git config escape")
			}
			value.WriteRune(c)
			end = value.Len()
			escaped = false
			continue
		}
		switch c {
		case '\\':
			escaped = true
		case '"':
			quoted = !quoted
		case '#', ';':
			if !quoted {
				return value.String()[:end], nil
			}
			value.WriteRune(c)
			end = value.Len()
		default:
			value.WriteRune(c)
			if quoted || c != ' ' && c != '\t' {
				end = value.Len()
			}
		}
	}
	if escaped || quoted {
		return "", fmt.Errorf("unterminated application Git config value")
	}
	return value.String()[:end], nil
}

func readReference(ctx context.Context, directory, common, reference string, seen map[string]bool) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if seen[reference] || len(seen) >= 32 {
		return "", fmt.Errorf("cyclic or overly deep Git symbolic reference %q", reference)
	}
	seen[reference] = true
	base := common
	if reference == "HEAD" || strings.HasPrefix(reference, "refs/worktree/") || strings.HasPrefix(reference, "refs/bisect/") || strings.HasPrefix(reference, "refs/rewritten/") {
		base = directory
	}
	data, err := os.ReadFile(filepath.Join(base, filepath.FromSlash(reference)))
	if errors.Is(err, os.ErrNotExist) && reference != "HEAD" {
		return packedRevision(ctx, common, reference)
	}
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(string(data))
	if target, found := strings.CutPrefix(value, "ref: "); found {
		if !validReference(target) {
			return "", fmt.Errorf("invalid application Git symbolic reference %q", target)
		}
		return readReference(ctx, directory, common, target, seen)
	}
	return validateRevision(value)
}

func validReference(reference string) bool {
	if !strings.HasPrefix(reference, "refs/") || strings.ContainsAny(reference, `\ :?*[~^`) || strings.Contains(reference, "..") || strings.Contains(reference, "@{") {
		return false
	}
	for _, component := range strings.Split(reference, "/") {
		if component == "" || strings.HasPrefix(component, ".") || strings.HasSuffix(component, ".") || strings.HasSuffix(component, ".lock") {
			return false
		}
		for _, c := range component {
			if c < 32 || c == 127 {
				return false
			}
		}
	}
	return true
}

func validateRevision(value string) (string, error) {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 20 || strings.Trim(value, "0") == "" {
		return "", fmt.Errorf("application Git HEAD has no valid SHA-1 revision")
	}
	return strings.ToLower(value), nil
}

func packedRevision(ctx context.Context, directory, reference string) (string, error) {
	file, err := os.Open(filepath.Join(directory, "packed-refs"))
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[1] == reference {
			return validateRevision(fields[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("application Git reference %q does not exist", reference)
}

// Find the nearest marker first: an unsupported inner repository must not make
// revision lookup search upward and report an unrelated ancestor's revision.
func repositoryDirectory(root string) (string, error) {
	directory, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Lstat(filepath.Join(directory, ".git")); err == nil {
			return directory, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("no application Git repository above %s", root)
		}
		directory = parent
	}
}
