package thinknum

import (
	"os"
	"testing"
)

func TestCanWrite(t *testing.T) {

	homePath, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	var scenarios = []struct {
		path     string
		expected bool
	}{
		{"/root", false},
		{"/tmp", true},
		{homePath, true},
		{"out", true},
	}

	for _, s := range scenarios {
		t.Run(s.path, func(t *testing.T) {
			got := isDirWritable(s.path)
			if got != s.expected {
				t.Errorf("Wrong result for %s. Expected: %v, got: %v", s.path, s.expected, got)
			}
		})
	}
}
