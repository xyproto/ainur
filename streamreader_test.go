package ainur

import (
	"bytes"
	"io"
	"testing"

	"github.com/bmizerany/assert"
)

func TestStreamReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte("aaaabbbbcccc"))
	r, err := NewStreamReader(buf, 6)
	assert.Equal(t, nil, err)

	found := false
	for {
		b, err := r.Next()
		if err == io.EOF {
			break
		}

		if bytes.Index(b, []byte("bbbb")) != -1 {
			found = true
		}
	}

	assert.Equal(t, true, found)
}
