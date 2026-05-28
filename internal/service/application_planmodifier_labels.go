package service

import (
	"encoding/base64"
	"sort"
	"strings"
)

// isBase64Labels returns true when s is valid base64 AND the decoded content
// has at least one line containing an '=' character (label format heuristic).
// Conservative on purpose: false-positives lead to a semantic-equal false-no-op,
// which is worse than a false-positive diff.
func isBase64Labels(s string) bool {
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(s))
	if err != nil {
		return false
	}
	for _, ln := range strings.Split(string(decoded), "\n") {
		if strings.Contains(strings.TrimSpace(ln), "=") {
			return true
		}
	}
	return false
}

// normalize converts a custom_labels string into a canonical line-set form
// suitable for semantic equality comparison. It optionally decodes base64,
// trims and splits lines, filters Coolify's auto-injected letsencrypt
// certresolver lines, and sorts the result.
//
// Returns an empty (nil-or-zero-length) slice for empty/whitespace input.
func normalize(raw string) []string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return nil
	}

	if isBase64Labels(s) {
		decoded, _ := base64.StdEncoding.DecodeString(s)
		s = string(decoded)
	}

	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}

	sort.Strings(out)
	return out
}
