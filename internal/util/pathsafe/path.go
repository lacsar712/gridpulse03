package pathsafe

import (
	"errors"
	"path/filepath"
	"strings"
)

func JoinUnder(root, rel string) (string, error) {
	if strings.TrimSpace(rel) == "" {
		return "", errors.New("empty path")
	}
	if filepath.IsAbs(rel) {
		return "", errors.New("absolute path")
	}
	clean := filepath.Clean(rel)
	full := filepath.Join(root, clean)
	relOut, err := filepath.Rel(filepath.Clean(root), full)
	if err != nil {
		return "", err
	}
	if relOut == ".." || strings.HasPrefix(relOut, ".."+string(filepath.Separator)) {
		return "", errors.New("path escapes root")
	}
	return full, nil
}
