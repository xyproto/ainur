// Package ainur provides functions for examining ELF files
package ainur

// Only depends on the standard library
import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"io"
	"regexp"
	"strings"
)

const (
	gccMarker   = "GCC: ("
	gnuEnding   = "GNU) "
	clangMarker = "clang version"
	rustMarker  = "rustc version"
	ghcMarker   = "GHC "
	ocamlMarker = "[ocaml]"
)

var (
	// GHCVersionRegex is a regexp for matching Glasgow Haskell Compiler version strings
	GHCVersionRegex = regexp.MustCompile(`GHC\ (\d{1,4}\.)(\d+\.)?(\d+)`)

	// GoVersionRegex is a regexp for matching Go version strings
	GoVersionRegex = regexp.MustCompile(`go(\d+\.)(\d+\.)?(\*|\d+)`)

	// PasVersionRegex is a regexp for matching Free Pascal Compiler version strings
	PasVersionRegex = regexp.MustCompile(`FPC\ (\d+\.)?(\d+\.)?(\*|\d+)`)

	// OcamlVersionRegex is a regexp for matching OCaml version strings
	OcamlVersionRegex = regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)`)

	// GCCVersionRegex0 is another regexp for matching GCC version strings
	GCCVersionRegex0 = regexp.MustCompile(`(\d+\.)(\d+\.)?(\*|\d+)\ `)

	// GCCVersionRegex1 is another regexp for matching GCC version strings
	GCCVersionRegex1 = regexp.MustCompile(`\) (\d{1,4}\.)(\d+\.)?(\*|\d+)\ `)

	// GCCVersionRegex2 is another regexp for matching GCC version strings
	GCCVersionRegex2 = regexp.MustCompile(` (\d{1,4}\.)(\d+\.)?(\*|\d+)`)

	// GCCVersionRegex3 is another regexp for matching GCC version strings
	GCCVersionRegex3 = regexp.MustCompile(`(\d{1,4}\.)(\d+\.)?(\*|\d+)`)

	// GCCVersionRegex4 is another regexp for matching GCC version strings
	GCCVersionRegex4 = regexp.MustCompile(`\) (\d{1,4}\.)(\d+\.)?(\*|\d+).(\d+)`)
)

// compilerVersionFunctions is a slice of functions that can be used
// for discovering a version string from an ELF file, ordered from
// the more specific to the more ambigous ones.
var compilerVersionFunctions = []func(*elf.File) string{
	GoVer,
	OCamlVer,
	GHCVer,
	RustVerUnstripped,
	RustVerStripped,
	DVer,
	GCCVer,
	PasVer,
	TCCVer,
}

// GHCVer returns the GHC compiler version or an empty string
// example output: "GHC 8.6.2"
func GHCVer(f *elf.File) (ver string) {
	sec := f.Section(".comment")
	if sec == nil {
		return
	}
	versionData, errData := sec.Data()
	if errData != nil {
		return
	}
	if bytes.Contains(versionData, []byte(ghcMarker)) {
		// Try the first regexp for picking out the version
		ghcVersion := bytes.TrimSpace(GHCVersionRegex.Find(versionData))
		if len(ghcVersion) > 0 {
			return "GHC " + string(ghcVersion[4:])
		}
	}
	return
}

// GCCVer returns the GCC compiler version or an empty string
// example output: "GCC 6.3.1"
// Also handles clang.
func GCCVer(f *elf.File) (ver string) {
	debug := false
	sec := f.Section(".comment")
	if sec == nil {
		return
	}
	versionData, errData := sec.Data()
	if errData != nil {
		return
	}
	if bytes.Contains(versionData, []byte(gccMarker)) {
		// Check if this is really clang
		if bytes.Contains(versionData, []byte(clangMarker)) {
			clangVersion := bytes.TrimSpace(GCCVersionRegex0.Find(versionData))
			return "Clang " + string(clangVersion)
		}
		// If the bytes are on this form: "GCC: (GNU) 6.3.0GCC: (GNU) 7.2.0",
		// use the largest version number.
		if bytes.Count(versionData, []byte(gccMarker)) > 1 {
			// Split in to 3 parts, always valid for >=2 instances of gccMarker
			elements := bytes.SplitN(versionData, []byte(gccMarker), 3)
			versionA := elements[1]
			versionB := elements[2]
			if bytes.HasPrefix(versionA, []byte(gnuEnding)) {
				versionA = versionA[5:]
			}
			if bytes.HasPrefix(versionB, []byte(gnuEnding)) {
				versionB = versionB[5:]
			}
			if FirstIsGreater(string(versionA), string(versionB)) {
				versionData = versionA
			} else {
				versionData = versionB
			}
		}
		// Try the first regexp for picking out the version
		gccVersion := bytes.TrimSpace(GCCVersionRegex1.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #1 " + string(gccVersion[2:]))
			}
			return "GCC " + string(gccVersion[2:])
		}
		// Try the second regexp for picking out the version
		gccVersion = bytes.TrimSpace(GCCVersionRegex2.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #2 " + string(gccVersion))
			}
			// Check that it does not start with "1.", that may happen
			if !bytes.HasPrefix(gccVersion, []byte("1.")) {
				return "GCC " + string(gccVersion)
			}
		}
		// Try the third regexp for picking out the version
		gccVersion = bytes.TrimSpace(GCCVersionRegex3.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #3 " + string(gccVersion))
			}
			// Check that it does not start with "1.", that may happen
			if !bytes.HasPrefix(gccVersion, []byte("1.")) {
				return "GCC " + string(gccVersion)
			}
		}
		// Try the fourth regexp for picking out the version
		gccVersion = bytes.TrimSpace(GCCVersionRegex4.Find(versionData))
		if len(gccVersion) > 0 {
			if debug {
				println("GCC #4 " + string(gccVersion))
			}
			return "GCC " + string(gccVersion)[2:]
		}
		// See what we've got
		gccVersionString := strings.TrimSpace(string(gccVersion))
		if len(gccVersionString) > 5 {
			if debug {
				println("GCC #4 " + string(gccVersion[5:]))
			}
			// Check that the version number is not "0"
			retver := string(gccVersion)[5:]
			if retver != "0" {
				return "GCC " + retver
			}
		}
		// Failed to find a GCC version string
		return
	}
	return string(versionData)
}

// RustVerUnstripped returns the Rust compiler version or an empty string
// example output: "Rust 1.27.0"
func RustVerUnstripped(f *elf.File) (ver string) {
	// Check if there is debug data in the executable, that may contain the version number
	sec := f.Section(".debug_str")
	if sec == nil {
		return
	}
	r := sec.Open()
	bufferSize := 8192
	margin := 1024
	sr, err := NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		pos1 := bytes.Index(b, []byte(rustMarker))
		if pos1 == -1 {
			continue
		}

		if bufferSize-pos1 < margin {
			continue
		}

		pos1 += len(rustMarker) + 1
		pos2 := bytes.Index(b[pos1:], []byte("("))
		if pos2 == -1 {
			return
		}
		pos2 += pos1
		versionString := strings.TrimSpace(string(b[pos1:pos2]))
		return "Rust " + versionString
	}

	return
}

// RustVerStripped returns the Rust compiler version or an empty string,
// from a stripped Rust executable. Does not contain the Rust
// version number.
// Example output: "Rust (GCC 8.1.0)"
func RustVerStripped(f *elf.File) (ver string) {
	// Check if the .gcc_except_table ELF section exists
	if f.Section(".gcc_except_table") == nil {
		return ""
	}
	// Check if the .rodata ELF section exists
	sec := f.Section(".rodata")
	if sec == nil {
		return
	}
	r := sec.Open()
	bufferSize := 8192
	sr, err := NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		// Look for the rust marker that may appear in new, stripped executables
		pos := bytes.Index(b, []byte("/rustc-"))
		if pos == -1 {
			continue
		}

		// Rust may use GCC for linking
		if gccVersion := GCCVer(f); gccVersion != "" {
			return "Rust (" + GCCVer(f) + ")"
		}
		return "Rust"
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	sr, err = NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		// Look for the rust marker that may appear in old, stripped executables
		rustIndex1 := bytes.Index(b, []byte("__rust_"))
		if rustIndex1 > 0 && b[rustIndex1-1] == 0 {
			// Rust may use GCC for linking
			if gccVersion := GCCVer(f); gccVersion != "" {
				return "Rust (" + GCCVer(f) + ")"
			}
			return "Rust"
		}
	}

	// No rust markers! Probably not created with the Rust compiler.
	return
}

// DVer returns "DMD" if it is detected
// Example output: "DMD"
func DVer(f *elf.File) (ver string) {
	// Check if the .dynstr ELF section exists
	sec := f.Section(".dynstr")
	if sec == nil {
		return
	}
	r := sec.Open()
	bufferSize := 8192
	sr, err := NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		// Look for the DMD marker
		if bytes.Contains(b, []byte("__dmd_")) {
			return "DMD"
		}
	}

	return
}

func offsetOf(f *elf.File, addr uint64) (*elf.Prog, uint64) {
	for _, prog := range f.Progs {
		if prog.Vaddr <= addr && addr <= prog.Vaddr+prog.Filesz-1 {
			return prog, addr - prog.Vaddr
		}
	}
	return nil, 0
}

// buildinfoGoVer returns the Go compiler version embedded in the
// `.go.buildinfo` section, or in the first writable data segment, or
// an empty string if it could not be read.
// example output: "Go 1.22.1"
func buildinfoGoVer(f *elf.File) string {
	// This function more-or-less closely follows
	// `debug/buildinfo` from the standard library, except
	// that it only works on ELF, and it uses
	// StreamReader to read.
	//
	// Most of the code there is not exposed publicly
	// (and the public API takes io.Reader, not elf.File),
	// so we reproduce it here rather than being able to
	// just call into it.
	decodeString := func(data []byte) []byte {
		u, n := binary.Uvarint(data)
		if n <= 0 || u > uint64(len(data)-n) {
			return nil
		}
		return data[n : uint64(n)+u]
	}
	var addr, size uint64
	for _, s := range f.Sections {
		if s.Name == ".go.buildinfo" {
			addr, size = s.Addr, s.Size
			break
		}
	}
	if size == 0 {
		for _, p := range f.Progs {
			if p.Type == elf.PT_LOAD && p.Flags&(elf.PF_X|elf.PF_W) == elf.PF_W {
				addr, size = p.Vaddr, p.Memsz
			}
		}
	}
	if size == 0 {
		return ""
	}

	prog, offset := offsetOf(f, addr)
	if prog == nil {
		return ""
	}
	n := prog.Filesz - offset
	if n > size {
		n = size
	}
	r := io.NewSectionReader(prog, int64(addr-prog.Vaddr), int64(n))

	sr, err := NewStreamReader(r, 8192)
	if err != nil {
		return ""
	}
	const (
		buildInfoAlign = 16
		buildInfoSize  = 32
	)
	var data []byte
	for {
		if data != nil {
			break
		}
		b, err := sr.Next()
		if err != nil {
			return ""
		}
		for {
			i := bytes.Index(b, []byte("\xff Go buildinf:"))
			if i < 0 || len(b)-i < buildInfoSize {
				break
			}
			if i%buildInfoAlign == 0 {
				data = b[i:]
				break
			}
			// round up to next offset
			b = b[(i+buildInfoAlign-1)&^(buildInfoAlign-1):]
		}
	}
	// Decode the blob.
	// The first 14 bytes are buildInfoMagic.
	// The next two bytes indicate pointer size in bytes (4 or 8) and endianness
	// (0 for little, 1 for big).
	// Two virtual addresses to Go strings follow that: runtime.buildVersion,
	// and runtime.modinfo.
	// On 32-bit platforms, the last 8 bytes are unused.
	// If the endianness has the 2 bit set, then the pointers are zero
	// and the 32-byte header is followed by varint-prefixed string data
	// for the two string values we care about.
	//
	// NOTE: The above comment is from Go's debug/buildinfo.
	// In this version, we only extract runtime.buildVersion, not
	// runtime.modinfo.
	ptrSize := int(data[14])
	var result []byte
	if data[15]&2 != 0 {
		result = decodeString(data[32:])
	} else {
		bigEndian := data[15] != 0
		var bo binary.ByteOrder
		if bigEndian {
			bo = binary.BigEndian
		} else {
			bo = binary.LittleEndian
		}
		var addr uint64
		if ptrSize == 4 {
			addr = uint64(bo.Uint32(data[16:]))
		} else if ptrSize == 8 {
			addr = bo.Uint64(data[16:])
		}
		prog, offset := offsetOf(f, addr)
		if prog == nil {
			return ""
		}
		hdr := make([]byte, 2*ptrSize)
		_, err := prog.ReadAt(hdr, int64(offset))
		if err != nil {
			return ""
		}
		var dataAddr, dataLen uint64
		if ptrSize == 4 {
			dataAddr = uint64(bo.Uint32(hdr))
			dataLen = uint64(bo.Uint32(hdr[4:]))
		} else {
			dataAddr = bo.Uint64(hdr)
			dataLen = bo.Uint64(hdr[8:])
		}
		dataProg, dataOffset := offsetOf(f, dataAddr)
		if dataProg == nil {
			return ""
		}
		// Our copy of this code doesn't have
		// access to the Go-internal `saferio` package for reading possibly
		// gibberish lengths of data without OOMing, so instead
		// just check that the length isn't outrageous.
		if dataLen > 1024 {
			return ""
		}
		data := make([]byte, dataLen)
		_, err = prog.ReadAt(data, int64(dataOffset))
		if err != nil {
			return ""
		}
		result = data
	}
	s := string(result)
	if strings.HasPrefix(s, "go") {
		return "Go " + s[2:]
	}
	return ""
}

// GoVer returns the Go compiler version or an empty string
// example output: "Go 1.8.3"
func GoVer(f *elf.File) string {
	buildinfoVer := buildinfoGoVer(f)
	if buildinfoVer != "" {
		return string(buildinfoVer)
	}
	return rodataGoVer(f)
}

func rodataGoVer(f *elf.File) (ver string) {
	sec := f.Section(".rodata")
	if sec == nil {
		return
	}
	r := sec.Open()
	bufferSize := 8192
	margin := 1024
	sr, err := NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		goVersionIndex := GoVersionRegex.FindIndex(b)
		if goVersionIndex == nil {
			continue
		}
		if bufferSize-goVersionIndex[0] < margin {
			continue
		}

		goVersion := string(GoVersionRegex.Find(b))
		if strings.HasPrefix(goVersion, "go") {
			return "Go " + goVersion[2:]
		}
		if goVersion == "" {
			gosec := f.Section(".gosymtab")
			if gosec != nil {
				return "Go (unknown version)"
			}
			return
		}
		return goVersion
	}

	return
}

// PasVer returns the FPC compiler version or an empty string
// example output: "FPC 3.0.2"
func PasVer(f *elf.File) (ver string) {
	sec := f.Section(".data")
	if sec == nil {
		return
	}

	r := sec.Open()
	bufferSize := 8192
	sr, err := NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		indexes := PasVersionRegex.FindIndex(b)
		if indexes == nil {
			continue
		}

		if indexes[0] == bufferSize-(indexes[1]-indexes[0]) {
			continue
		}

		return string(PasVersionRegex.Find(b))
	}

	return
}

// TCCVer returns "TCC" or an empty string
// TCC has no version number, but it does have some signature sections.
func TCCVer(f *elf.File) (ver string) {
	// .note.ABI-tag must be missing
	if f.Section(".note.ABI-tag") != nil {
		// TCC does not normally have this section, not TCC
		return
	}
	if f.Section(".rodata.cst4") == nil {
		// TCC usually has this section, not TCC
		return
	}
	return "TCC"
}

// OCamlVer returns the OCaml compiler version or an empty string
// example output: "OCaml 4.05.0"
func OCamlVer(f *elf.File) (ver string) {
	sec := f.Section(".rodata")
	if sec == nil {
		return
	}
	r := sec.Open()
	bufferSize := 8192
	margin := 1024
	sr, err := NewStreamReader(r, bufferSize)
	if err != nil {
		return
	}

	for {
		b, err := sr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		pos := bytes.Index(b, []byte(ocamlMarker))
		if pos == -1 {
			continue
		}

		if bufferSize-pos < margin {
			continue
		}

		ocamlVersion := "OCaml " + string(OcamlVersionRegex.Find(b))
		if ocamlVersion == "" {
			return "OCaml (unknown version)"
		}
		return ocamlVersion
	}

	return
}

// Compiler takes an *elf.File and tries to find which compiler and version
// it was compiled with, by probing for known locations, strings and patterns.
func Compiler(f *elf.File) string {
	// Loop over the functions that can be used for extracting a version string
	for _, compilerVersion := range compilerVersionFunctions {
		// Call compilerVersion to check if a compiler version is found
		if ver := compilerVersion(f); ver != "" {
			return ver
		}
	}
	return "unknown"
}

// Stripped returns true if symbols can not be retrieved from the given ELF file
func Stripped(f *elf.File) bool {
	return f.SectionByType(elf.SHT_SYMTAB) == nil
}

// Examine tries to discover which compiler and compiler version the given
// file was compiled with.
func Examine(filename string) (string, error) {
	f, err := elf.Open(filename)
	if err != nil {
		if strings.HasPrefix(err.Error(), "bad magic") {
			return "", errors.New(filename + ": Not an ELF")
		}
		return "", err
	}
	defer f.Close()
	return Compiler(f), nil
}

// MustExamine does the same as examine, but panics instead of returning an error
func MustExamine(filename string) string {
	compilerVersion, err := Examine(filename)
	if err != nil {
		panic(err)
	}
	return compilerVersion
}
