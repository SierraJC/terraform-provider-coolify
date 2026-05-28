package service

import "encoding/base64"

// base64StdEncodingHelper is a test-only helper to construct base64-encoded
// label strings without importing encoding/base64 directly in the table-driven
// test file (keeps the test file focused on assertions, not encoding).
func base64StdEncodingHelper(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
