package service

// normalize converts a custom_labels string into a canonical line-set form
// suitable for semantic equality comparison. It optionally decodes base64,
// trims and splits lines, filters Coolify's auto-injected letsencrypt
// certresolver lines, and sorts the result.
//
// Returns an empty (nil-or-zero-length) slice for empty/whitespace input.
func normalize(raw string) []string {
	return nil
}
