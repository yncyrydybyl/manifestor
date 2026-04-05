package grab

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	unsafeChars    = regexp.MustCompile(`[^a-z0-9._-]`)
	multipleHyphen = regexp.MustCompile(`-{2,}`)
)

// Sanitize transforms a file name into a clean, lowercase,
// filesystem-friendly form. The extension is preserved.
//
// Rules:
//   - lowercase everything
//   - replace spaces and unsafe characters with hyphens
//   - collapse consecutive hyphens
//   - trim leading/trailing hyphens from the stem
//   - preserve the original extension (lowercased)
func Sanitize(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	stem := strings.TrimSuffix(name, filepath.Ext(name))

	stem = strings.ToLower(stem)
	stem = unsafeChars.ReplaceAllString(stem, "-")
	stem = multipleHyphen.ReplaceAllString(stem, "-")
	stem = strings.Trim(stem, "-")

	if stem == "" {
		stem = "file"
	}

	return stem + ext
}
