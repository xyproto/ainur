package ainur

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestStatic(t *testing.T) {
	assert.Equal(t, MustExamineStatic("testdata/go_upx_amd64_static"), true)
}

func TestDynamic(t *testing.T) {
	assert.Equal(t, MustExamineStatic("testdata/nano_voidlinux"), false)
}
