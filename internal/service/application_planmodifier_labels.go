package service

import (
	"encoding/base64"
	"regexp"
	"sort"
	"strings"
)

// leCertResolverPattern matches Coolify's auto-injected LetsEncrypt cert
// resolver line on Traefik routers. Key part is case-insensitive (matches
// what Coolify produces); value is enforced exact "letsencrypt" by
// hasExactLetsencryptValue — custom resolver names (custom-ca, mycorp,
// uppercase LETSENCRYPT, etc.) are preserved as genuine user intent.
var leCertResolverPattern = regexp.MustCompile(
	`(?i)^traefik\.http\.routers\.[a-z0-9-]+\.tls\.certresolver` +
		`\s*=\s*letsencrypt\s*$`,
)

// filterCertResolver drops any lines matching Coolify's auto-injected
// letsencrypt cert resolver pattern. Other lines (including custom
// certresolver values) are preserved.
func filterCertResolver(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		if leCertResolverPattern.MatchString(ln) && hasExactLetsencryptValue(ln) {
			continue
		}
		out = append(out, ln)
	}
	return out
}

// hasExactLetsencryptValue confirms the value (right of `=`) is exactly
// "letsencrypt" (lowercase, possibly with surrounding whitespace).
func hasExactLetsencryptValue(line string) bool {
	idx := strings.Index(line, "=")
	if idx < 0 {
		return false
	}
	val := strings.TrimSpace(line[idx+1:])
	return val == "letsencrypt"
}

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
	trimmed := make([]string, 0, len(lines))
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" {
			continue
		}
		trimmed = append(trimmed, t)
	}

	filtered := filterCertResolver(trimmed)

	sort.Strings(filtered)
	return filtered
}
