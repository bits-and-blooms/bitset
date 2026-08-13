package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPextByte(t *testing.T) {
	cases := []struct {
		b, m, want uint8
	}{
		{0x00, 0x00, 0x00},
		{0xFF, 0x00, 0x00},
		{0xAB, 0xFF, 0xAB}, // full mask: identity
		{0xDA, 0x66, 0x09}, // bits at positions 1, 2, 5, 6 are 1, 0, 0, 1
		{0x80, 0x80, 0x01},
	}
	for _, c := range cases {
		if got := pextByte(c.b, c.m); got != c.want {
			t.Errorf("pextByte(%#02x, %#02x) = %#02x, want %#02x", c.b, c.m, got, c.want)
		}
	}
}

func TestPdepByte(t *testing.T) {
	cases := []struct {
		b, m, want uint8
	}{
		{0x00, 0x00, 0x00},
		{0xFF, 0x00, 0x00},
		{0xAB, 0xFF, 0xAB}, // full mask: identity
		{0x09, 0x66, 0x42}, // bits 0 and 3 go to positions 1 and 6
		{0x01, 0x80, 0x80},
	}
	for _, c := range cases {
		if got := pdepByte(c.b, c.m); got != c.want {
			t.Errorf("pdepByte(%#02x, %#02x) = %#02x, want %#02x", c.b, c.m, got, c.want)
		}
	}
}

func TestPextPdepRoundTrip(t *testing.T) {
	// pdep is the inverse of pext for the bits selected by the mask.
	for b := 0; b < 256; b++ {
		for _, m := range []uint8{0x00, 0x0F, 0x55, 0xAA, 0xFF} {
			extracted := pextByte(uint8(b), m)
			if got := pdepByte(extracted, m); got != uint8(b)&m {
				t.Fatalf("pdepByte(pextByte(%#02x, %#02x), %#02x) = %#02x, want %#02x", b, m, m, got, uint8(b)&m)
			}
		}
	}
}

func TestGenerateTable(t *testing.T) {
	var single [256]uint8

	single[0] = 1

	out := generateTable("popLUT", single, "population counts")
	if !strings.Contains(out, "// population counts\n") {
		t.Errorf("expected the comment in the output, got %q", out[:40])
	}

	if !strings.Contains(out, "var popLUT = [256]uint8{") {
		t.Error("expected a [256]uint8 table declaration in the output")
	}

	if !strings.Contains(out, "1,") {
		t.Error("expected the table values in the output")
	}

	var double [256][256]uint8

	double[0][0] = 7

	out = generateTable("pextLUT", double, "")
	if !strings.Contains(out, "var pextLUT = [256][256]uint8{") {
		t.Error("expected a [256][256]uint8 table declaration in the output")
	}

	if !strings.Contains(out, "7,") {
		t.Error("expected the table values in the output")
	}
}

// setArgs replaces the process arguments and resets the flag set so that
// main() can be invoked more than once.
func setArgs(t *testing.T, args ...string) {
	t.Helper()

	oldArgs, oldCommandLine := os.Args, flag.CommandLine

	t.Cleanup(func() {
		os.Args, flag.CommandLine = oldArgs, oldCommandLine
	})

	os.Args = args
	flag.CommandLine = flag.NewFlagSet(args[0], flag.ContinueOnError)
}

func TestMainMissingPackageName(t *testing.T) {
	setArgs(t, "pextgen")
	main() // must return after reporting the missing -pkg flag
}

func TestMainGeneratesFile(t *testing.T) {
	dir := t.TempDir()

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.Chdir(oldwd)
		if err != nil {
			t.Fatal(err)
		}
	})

	setArgs(t, "pextgen", "-pkg", "bitset")
	main()

	data, err := os.ReadFile(filepath.Join(dir, "pext.gen.go"))
	if err != nil {
		t.Fatalf("expected main to generate pext.gen.go: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "package bitset") {
		t.Error("expected the generated file to declare the requested package")
	}

	for _, table := range []string{"pextLUT", "pdepLUT", "popLUT"} {
		if !strings.Contains(content, "var "+table+" = ") {
			t.Errorf("expected the generated file to contain the %s table", table)
		}
	}
}
