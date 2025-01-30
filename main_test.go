package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCompareVersions tests the version comparison logic
func TestCompareVersions(t *testing.T) {
	testCases := []struct {
		name     string
		version1 string
		version2 string
		expected bool
	}{
		{
			name:     "Major version comparison",
			version1: "v1.0.0",
			version2: "v2.0.0",
			expected: true,
		},
		{
			name:     "Same version should return false",
			version1: "v1.0.0",
			version2: "v1.0.0",
			expected: false,
		}, {
			name:     "Minor version comparison",
			version1: "v1.1.0",
			version2: "v1.2.0",
			expected: true,
		}, {
			name:     "Same version should return false",
			version1: "v1.1.0",
			version2: "v1.1.0",
			expected: false,
		}, {
			name:     "Patch version comparison",
			version1: "v1.1.1",
			version2: "v1.1.2",
			expected: true,
		}, {
			name:     "Same version should return false",
			version1: "v1.1.1",
			version2: "v1.1.1",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compareVersions(tc.version1, tc.version2)
			if result != tc.expected {
				t.Errorf("compareVersions(%s, %s) = %v; want %v",
					tc.version1, tc.version2, result, tc.expected)
			}
		})
	}
}
