package ainur

import (
	"testing"
)

func TestStatic(t *testing.T) {
	if result := MustExamineStatic("testdata/go_upx_amd64_static"); !result {
		t.Errorf("Expected true for static file, got %v", result)
	}
}

func TestDynamic(t *testing.T) {
	if result := MustExamineStatic("testdata/nano_voidlinux"); result {
		t.Errorf("Expected false for dynamic file, got %v", result)
	}
}
