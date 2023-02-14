package ainur

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestVoidLinuxNano(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/nano_voidlinux"), "GCC 7.2.0")
}

func TestArchLinuxLs(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/ls_archlinux"), "GCC 7.1.1")
}

func TestClang(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/clang_hello"), "Clang 8.2.1")
}

func TestTCC(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/tcc_hello"), "TCC")
}

func TestRustStripped(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/rust_hello_stripped"), "Rust (GCC 8.1.0)")
}

func TestRust(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/rust_hello"), "Rust 1.27.0-nightly")
}

func TestRust2(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/bat"), "Rust (GCC 8.2.1)")
}

func TestD(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/dmd"), "DMD")
}

func TestGCC1(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/afl-analyze"), "GCC 7.2.0")
}

func TestPowerPC(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/e500v2"), "GCC 4.7.2")
}

func TestPowerPC2(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/e500v2_gcc8"), "GCC 8.3.0")
}

func TestGCC2(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/gcc820"), "GCC 8.2.0")
}

func TestGHC(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/ghc"), "GHC 8.6.2")
}

func TestPowerPC3(t *testing.T) {
	assert.Equal(t, MustExamine("testdata/nc_e500v2"), "unknown")
}

var doNotOptimiseString string

func BenchmarkVoidLinux(b *testing.B) {
	b.ReportAllocs()

	var result string
	for n := 0; n < b.N; n++ {
		result = MustExamine("testdata/nano_voidlinux")
	}
	doNotOptimiseString = result
}

func TestVersionCompare(t *testing.T) {
	assert.Equal(t, FirstIsGreater("2", "1.0.7.abc"), true)
	assert.Equal(t, FirstIsGreater("2.0", "2.0 alpha1"), true)
	assert.Equal(t, FirstIsGreater("1.0", "2.0.3"), false)
	assert.Equal(t, FirstIsGreater("2.0", "2.0.rc1"), true)
}
