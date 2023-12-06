package ainur

import (
	"testing"
)

func TestVersionCompare(t *testing.T) {
	if !FirstIsGreater("2", "1.0.7.abc") {
		t.Errorf("Expected 2 to be greater than 1.0.7.abc")
	}
	if !FirstIsGreater("2.0", "2.0 alpha1") {
		t.Errorf("Expected 2.0 to be greater than 2.0 alpha1")
	}
	if FirstIsGreater("1.0", "2.0.3") {
		t.Errorf("Expected 1.0 not to be greater than 2.0.3")
	}
	if !FirstIsGreater("2.0", "2.0.rc1") {
		t.Errorf("Expected 2.0 to be greater than 2.0.rc1")
	}
}

func TestVersionSum(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected int
	}{
		{"Empty", []string{}, 0},
		{"SingleZero", []string{"0"}, 0},
		{"SingleNumber", []string{"5"}, 5},
		{"MultipleZeros", []string{"0", "0", "0"}, 0},
		{"Version2.0.0.0", []string{"2", "0", "0", "0"}, 2000},
		{"Version1.2.3", []string{"1", "2", "3"}, 123},
		{"NonNumber", []string{"abc"}, -1},
		{"Mixed", []string{"1", "abc", "3"}, 93},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sum := versionSum(test.parts)
			if sum != test.expected {
				t.Errorf("versionSum(%v) = %v; want %v", test.parts, sum, test.expected)
			}
		})
	}
}

func TestFirstIsGreater(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"EqualVersions", "1.0", "1.0", false},
		{"AGreater", "1.1", "1.0", true},
		{"BGreater", "1.0", "1.1", false},
		{"NonNumberInB", "1.0", "abc", true},
		{"NonNumberInA", "abc", "1.0", false},
		{"LongerVersionInA", "1.0.0", "1.0", false},
		{"LongerVersionInB", "1.0", "1.0.0", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FirstIsGreater(test.a, test.b)
			if result != test.expected {
				t.Errorf("FirstIsGreater(%v, %v) = %v; want %v", test.a, test.b, result, test.expected)
			}
		})
	}
}
