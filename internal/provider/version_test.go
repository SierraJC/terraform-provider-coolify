package provider

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version       string
		expectedMajor int
		expectedMinor int
		expectedPatch int
		expectedBeta  int
	}{
		{"1.2.3", 1, 2, 3, 0},
		{"1.2.3-beta.1", 1, 2, 3, 1},
		{"2.0.0", 2, 0, 0, 0},
		{"2.0.0-beta.2", 2, 0, 0, 2},
		{"0.1.0", 0, 1, 0, 0},
		{"0.1.0-beta.3", 0, 1, 0, 3},
		{"10.20.30", 10, 20, 30, 0},
		{"10.20.30-beta.4", 10, 20, 30, 4},
		{"0.0.1", 0, 0, 1, 0},
		{"0.0.1-beta.5", 0, 0, 1, 5},
		{"99.99.99", 99, 99, 99, 0},
		{"99.99.99-beta.99", 99, 99, 99, 99},
		{"v4.0.0-beta.420.6", 4, 0, 0, 420},
		{"1.2", 1, 2, 0, 0},
		{"5", 5, 0, 0, 0},
		{"v1.2.3", 1, 2, 3, 0},
	}

	for _, test := range tests {
		major, minor, patch, beta := parseVersion(test.version)
		if major != test.expectedMajor || minor != test.expectedMinor || patch != test.expectedPatch || beta != test.expectedBeta {
			t.Errorf("parseVersion(%q) = (%d, %d, %d, %d); want (%d, %d, %d, %d)",
				test.version, major, minor, patch, beta, test.expectedMajor, test.expectedMinor, test.expectedPatch, test.expectedBeta)
		}
	}
}

func TestIsVersionCompatible(t *testing.T) {
	tests := []struct {
		currentVersion string
		minVersion     string
		expected       bool
	}{
		{"1.2.3", "1.2.3", true},
		{"1.2.4", "1.2.3", true},
		{"1.2.2", "1.2.3", false},
		{"2.0.0", "1.2.3", true},
		{"1.2.3-beta.1", "1.2.3-beta.1", true},
		{"1.2.3-beta.2", "1.2.3-beta.1", true},
		{"1.2.3-beta.1", "1.2.3-beta.2", false},
		{"1.2.3", "1.2.3-beta.1", true},
		{"1.2.3-beta.1", "1.2.3", false},
		{"0.1.0", "0.1.0", true},
		{"0.1.1", "0.1.0", true},
		{"0.0.9", "0.1.0", false},
		{"10.20.30", "10.20.30", true},
		{"10.20.31", "10.20.30", true},
		{"10.20.29", "10.20.30", false},
		{"99.99.99", "99.99.99", true},
		{"99.99.99-beta.99", "99.99.99-beta.98", true},
		{"99.99.99-beta.98", "99.99.99-beta.99", false},
	}

	for _, test := range tests {
		result := isVersionCompatible(test.currentVersion, test.minVersion)
		if result != test.expected {
			t.Errorf("isVersionCompatible(%q, %q) = %v; want %v",
				test.currentVersion, test.minVersion, result, test.expected)
		}
	}
}
