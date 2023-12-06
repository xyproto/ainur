package ainur

import (
	"testing"
)

func TestVoidLinuxNano(t *testing.T) {
	if result := MustExamine("testdata/nano_voidlinux"); result != "GCC 7.2.0" {
		t.Errorf("Expected GCC 7.2.0, got %s", result)
	}
}

func TestArchLinuxLs(t *testing.T) {
	if result := MustExamine("testdata/ls_archlinux"); result != "GCC 7.1.1" {
		t.Errorf("Expected GCC 7.1.1, got %s", result)
	}
}

func TestClang(t *testing.T) {
	if result := MustExamine("testdata/clang_hello"); result != "Clang 8.2.1" {
		t.Errorf("Expected Clang 8.2.1, got %s", result)
	}
}

func TestTCC(t *testing.T) {
	if result := MustExamine("testdata/tcc_hello"); result != "TCC" {
		t.Errorf("Expected TCC, got %s", result)
	}
}

func TestRustStripped(t *testing.T) {
	if result := MustExamine("testdata/rust_hello_stripped"); result != "Rust (GCC 8.1.0)" {
		t.Errorf("Expected Rust (GCC 8.1.0), got %s", result)
	}
}

func TestRust(t *testing.T) {
	if result := MustExamine("testdata/rust_hello"); result != "Rust 1.27.0-nightly" {
		t.Errorf("Expected Rust 1.27.0-nightly, got %s", result)
	}
}

func TestGo(t *testing.T) {
	if result := MustExamine("testdata/go_hello"); result != "Go 1.20.1" {
		t.Errorf("Expected Go 1.20.1, got %s", result)
	}
}

func TestRust2(t *testing.T) {
	if result := MustExamine("testdata/bat"); result != "Rust (GCC 8.2.1)" {
		t.Errorf("Expected Rust (GCC 8.2.1), got %s", result)
	}
}

func TestD(t *testing.T) {
	if result := MustExamine("testdata/dmd"); result != "DMD" {
		t.Errorf("Expected DMD, got %s", result)
	}
}

func TestGCC1(t *testing.T) {
	if result := MustExamine("testdata/afl-analyze"); result != "GCC 7.2.0" {
		t.Errorf("Expected GCC 7.2.0, got %s", result)
	}
}

func TestPowerPC(t *testing.T) {
	if result := MustExamine("testdata/e500v2"); result != "GCC 4.7.2" {
		t.Errorf("Expected GCC 4.7.2, got %s", result)
	}
}

func TestPowerPC2(t *testing.T) {
	if result := MustExamine("testdata/e500v2_gcc8"); result != "GCC 8.3.0" {
		t.Errorf("Expected GCC 8.3.0, got %s", result)
	}
}

func TestGCC2(t *testing.T) {
	if result := MustExamine("testdata/gcc820"); result != "GCC 8.2.0" {
		t.Errorf("Expected GCC 8.2.0, got %s", result)
	}
}

func TestGHC(t *testing.T) {
	if result := MustExamine("testdata/ghc"); result != "GHC 8.6.2" {
		t.Errorf("Expected GHC 8.6.2, got %s", result)
	}
}

func TestPowerPC3(t *testing.T) {
	if result := MustExamine("testdata/nc_e500v2"); result != "unknown" {
		t.Errorf("Expected unknown, got %s", result)
	}
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
