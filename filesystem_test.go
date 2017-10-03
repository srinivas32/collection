package collection

import (
	"fmt"
	"math"
	"path/filepath"
	"testing"
)

func TestEnumerateDirectoryOptions_UniqueBits(t *testing.T) {
	isPowerOfTwo := func(subject float64) bool {
		a := math.Abs(math.Log2(subject))
		b := math.Floor(a)

		return a-b < .0000001
	}

	if !isPowerOfTwo(64) {
		t.Log("isPowerOfTwo decided 64 is not a power of two.")
		t.FailNow()
	}

	if isPowerOfTwo(91) {
		t.Log("isPowerOfTwo decided 91 is a power of two.")
		t.FailNow()
	}

	seen := make(map[DirectoryOptions]struct{})

	declared := []DirectoryOptions{
		DirectoryOptionsExcludeFiles,
		DirectoryOptionsExcludeDirectories,
		DirectoryOptionsRecursive,
	}

	for _, option := range declared {
		if _, ok := seen[option]; ok {
			t.Logf("Option: %d has already been declared.", option)
			t.Fail()
		}
		seen[option] = struct{}{}

		if !isPowerOfTwo(float64(option)) {
			t.Logf("Option should have been a power of two, got %g instead.", float64(option))
			t.Fail()
		}
	}
}

func TestDirectory_Enumerate(t *testing.T) {
	subject := Directory{
		Location: filepath.Join(".", "testdata", "foo"),
	}

	testCases := []struct {
		options  DirectoryOptions
		expected map[string]struct{}
	}{
		{
			options: 0,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "a.txt"): struct{}{},
				filepath.Join("testdata", "foo", "c.txt"): struct{}{},
				filepath.Join("testdata", "foo", "bar"):   struct{}{},
			},
		},
		{
			options: DirectoryOptionsExcludeFiles,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "bar"): struct{}{},
			},
		},
		{
			options: DirectoryOptionsExcludeDirectories,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "a.txt"): struct{}{},
				filepath.Join("testdata", "foo", "c.txt"): struct{}{},
			},
		},
		{
			options: DirectoryOptionsRecursive,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "bar"):          struct{}{},
				filepath.Join("testdata", "foo", "bar", "b.txt"): struct{}{},
				filepath.Join("testdata", "foo", "a.txt"):        struct{}{},
				filepath.Join("testdata", "foo", "c.txt"):        struct{}{},
			},
		},
		{
			options: DirectoryOptionsExcludeFiles | DirectoryOptionsRecursive,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "bar"): struct{}{},
			},
		},
		{
			options: DirectoryOptionsRecursive | DirectoryOptionsExcludeDirectories,
			expected: map[string]struct{}{
				filepath.Join("testdata", "foo", "a.txt"):        struct{}{},
				filepath.Join("testdata", "foo", "bar", "b.txt"): struct{}{},
				filepath.Join("testdata", "foo", "c.txt"):        struct{}{},
			},
		},
		{
			options:  DirectoryOptionsExcludeDirectories | DirectoryOptionsExcludeFiles,
			expected: map[string]struct{}{},
		},
		{
			options:  DirectoryOptionsExcludeFiles | DirectoryOptionsRecursive | DirectoryOptionsExcludeDirectories,
			expected: map[string]struct{}{},
		},
	}

	for _, tc := range testCases {
		subject.Options = tc.options
		t.Run(fmt.Sprintf("%d", uint(tc.options)), func(t *testing.T) {
			for entry := range subject.Enumerate(nil) {
				cast := entry.(string)
				if _, ok := tc.expected[cast]; !ok {
					t.Logf("unexpected result: %q", cast)
					t.Fail()
				}
				delete(tc.expected, cast)
			}

			if len(tc.expected) != 0 {
				for unseenFile := range tc.expected {
					t.Logf("missing file: %q", unseenFile)
				}
				t.Fail()
			}
		})
	}
}
