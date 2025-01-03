package provider

import (
	"strconv"
	"strings"
)

// parseVersion parses a version string and returns the major, minor, patch, and beta versions as integers.
// The version string is expected to be in the format "major.minor.patch-beta.betaVersion".
// For example, given the version string "4.0.0-beta.360", it will return:
// major = 4, minor = 0, patch = 0, beta = 360.
func parseVersion(version string) (major, minor, patch, beta int) {
	// Example version string: "4.0.0-beta.360"
	parts := strings.Split(version, "-")
	versionParts := strings.Split(parts[0], ".")

	major, _ = strconv.Atoi(versionParts[0])
	minor, _ = strconv.Atoi(versionParts[1])
	patch, _ = strconv.Atoi(versionParts[2])

	if len(parts) > 1 && strings.HasPrefix(parts[1], "beta.") {
		beta, _ = strconv.Atoi(strings.TrimPrefix(parts[1], "beta."))
	}

	return
}

// isVersionCompatible checks if the current version is compatible with the minimum required version.
// It compares the major, minor, patch, and beta versions in sequence to determine compatibility.
func isVersionCompatible(currentVersion, minVersion string) bool {
	currentMajor, currentMinor, currentPatch, currentBeta := parseVersion(currentVersion)
	minMajor, minMinor, minPatch, minBeta := parseVersion(minVersion)

	switch {
	case currentMajor != minMajor:
		return currentMajor > minMajor
	case currentMinor != minMinor:
		return currentMinor > minMinor
	case currentPatch != minPatch:
		return currentPatch > minPatch
	case currentBeta == 0 || minBeta == 0:
		return currentBeta == 0
	default:
		return currentBeta >= minBeta
	}
}
