package ainur

import (
	"bytes"
	"io"
	"testing"
)

func TestStreamReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte("aaaabbbbcccc"))
	r, err := NewStreamReader(buf, 6)

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	found := false
	for {
		b, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Unexpected error during reading: %v", err)
			break
		}
		if bytes.Contains(b, []byte("bbbb")) {
			found = true
		}
	}

	// Check if "bbbb" was found
	if !found {
		t.Errorf("Expected to find 'bbbb' in the stream, but did not")
	}
}
