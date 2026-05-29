package main

import (
	"strings"
	"testing"
	"time"
)

func TestHelpAliases(t *testing.T) {
	for _, arg := range []string{"/?", "/H", "/HELP", "-h", "-help", "--help"} {
		opt, files, err := parseArgs([]string{arg})
		if err != nil {
			t.Fatalf("%s: %v", arg, err)
		}
		if !opt.ShowHelp {
			t.Fatalf("%s did not set ShowHelp", arg)
		}
		if len(files) != 0 {
			t.Fatalf("%s produced files: %v", arg, files)
		}
	}
}

func TestVersionAliases(t *testing.T) {
	for _, arg := range []string{"-v", "--version"} {
		opt, _, err := parseArgs([]string{arg})
		if err != nil {
			t.Fatalf("%s: %v", arg, err)
		}
		if !opt.ShowVersion {
			t.Fatalf("%s did not set ShowVersion", arg)
		}
	}
}

func TestDefaultSelectors(t *testing.T) {
	opt, files, err := parseArgs([]string{"file.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if !opt.SelectCreated || !opt.SelectModified || !opt.SelectAccessed || opt.SelectorSeen {
		t.Fatalf("bad default selectors: %+v", opt)
	}
	if len(files) != 1 || files[0] != "file.txt" {
		t.Fatalf("bad files: %v", files)
	}
}

func TestExplicitSelectors(t *testing.T) {
	tests := []struct {
		args    []string
		c, m, a bool
	}{
		{[]string{"-a", "file"}, false, false, true},
		{[]string{"-m", "file"}, false, true, false},
		{[]string{"-c", "file"}, true, false, false},
		{[]string{"-am", "file"}, false, true, true},
		{[]string{"-ac", "file"}, true, false, true},
		{[]string{"-mc", "file"}, true, true, false},
		{[]string{"-amc", "file"}, true, true, true},
	}
	for _, tt := range tests {
		opt, _, err := parseArgs(tt.args)
		if err != nil {
			t.Fatalf("%v: %v", tt.args, err)
		}
		if opt.SelectCreated != tt.c || opt.SelectModified != tt.m || opt.SelectAccessed != tt.a || !opt.SelectorSeen {
			t.Fatalf("%v: selectors got c=%v m=%v a=%v seen=%v", tt.args, opt.SelectCreated, opt.SelectModified, opt.SelectAccessed, opt.SelectorSeen)
		}
	}
}

func TestBooleanClusters(t *testing.T) {
	opt, _, err := parseArgs([]string{"-amcnpr", "file"})
	if err != nil {
		t.Fatal(err)
	}
	if !opt.SelectAccessed || !opt.SelectModified || !opt.SelectCreated || !opt.NoCreate || !opt.Parents || !opt.Runas {
		t.Fatalf("cluster not fully parsed: %+v", opt)
	}
}

func TestValueOptionsSeparated(t *testing.T) {
	opt, files, err := parseArgs([]string{"-amc", "-d", "2026-05-28 21:30:00", "-s", "ref.txt", "file"})
	if err == nil {
		t.Fatalf("expected base source conflict, got opt=%+v files=%v", opt, files)
	}

	opt, files, err = parseArgs([]string{"-d", "2026-05-28 21:30:00", "-amc", "file"})
	if err != nil {
		t.Fatal(err)
	}
	if opt.DateString != "2026-05-28 21:30:00" || !opt.SelectCreated || !opt.SelectModified || !opt.SelectAccessed {
		t.Fatalf("bad date parse: %+v", opt)
	}
	if len(files) != 1 || files[0] != "file" {
		t.Fatalf("bad files: %v", files)
	}

	opt, _, err = parseArgs([]string{"-s", "reference.txt", "target.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if opt.SyncFile != "reference.txt" {
		t.Fatalf("bad sync file: %+v", opt)
	}

	opt, _, err = parseArgs([]string{"-t", "202605282130.00", "file"})
	if err != nil {
		t.Fatal(err)
	}
	if opt.TouchTime != "202605282130.00" {
		t.Fatalf("bad touch time: %+v", opt)
	}

	opt, _, err = parseArgs([]string{"-e", "-0100", "file"})
	if err != nil {
		t.Fatal(err)
	}
	if !opt.UseEdit || opt.EditOffset != -60 {
		t.Fatalf("bad negative edit: %+v", opt)
	}
}

func TestRejectAttachedValuesAndValueInCluster(t *testing.T) {
	bad := [][]string{
		{"-d2026-05-28T21:30:00", "file"},
		{"-t202605282130.00", "file"},
		{"-sreference.txt", "file"},
		{"-e0130", "file"},
		{"-amct", "202605282130.00", "file"},
		{"-amce", "0130", "file"},
	}
	for _, args := range bad {
		_, _, err := parseArgs(args)
		if err == nil {
			t.Fatalf("expected error for %v", args)
		}
	}
}

func TestConflicts(t *testing.T) {
	bad := [][]string{
		{"-i", "-f", "file"},
		{"-s", "ref", "-d", "2026-05-28 21:30:00", "file"},
		{"-s", "ref", "-t", "202605282130.00", "file"},
		{"-d", "2026-05-28 21:30:00", "-t", "202605282130.00", "file"},
	}
	for _, args := range bad {
		_, _, err := parseArgs(args)
		if err == nil {
			t.Fatalf("expected conflict for %v", args)
		}
	}
}

func TestOffsetParser(t *testing.T) {
	tests := map[string]int64{
		"30":      30,
		"+30":     30,
		"-30":     -30,
		"0130":    90,
		"+0130":   90,
		"-0130":   -90,
		"010000":  3600,
		"-020000": -7200,
	}
	for s, want := range tests {
		got, err := parseOffset(s)
		if err != nil {
			t.Fatalf("%s: %v", s, err)
		}
		if got != want {
			t.Fatalf("%s: got %d want %d", s, got, want)
		}
	}
	for _, s := range []string{"", "1", "001", "0060", "0160", "abcdef", "+"} {
		if _, err := parseOffset(s); err == nil {
			t.Fatalf("expected bad offset %q", s)
		}
	}
}

func TestDateParser(t *testing.T) {
	valid := []string{
		"2026-05-28 21:30:00",
		"2026-05-28T21:30:00",
		"2026-05-28 21:30:00.123456789",
		"2026-05-28T21:30:00.123456789",
		"2026-05-28T21:30:00Z",
	}
	for _, s := range valid {
		if _, err := parseDate(s); err != nil {
			t.Fatalf("valid date %q rejected: %v", s, err)
		}
	}
	for _, s := range []string{"", "now", "@1700000000", "2026-05-28", "2026/05/28 21:30:00"} {
		if _, err := parseDate(s); err == nil {
			t.Fatalf("invalid date %q accepted", s)
		}
	}
}

func TestTouchTimeParser(t *testing.T) {
	for _, s := range []string{"05282130", "2605282130.00", "202605282130.00"} {
		if _, err := parseTouchTime(s); err != nil {
			t.Fatalf("valid touch time %q rejected: %v", s, err)
		}
	}
	for _, s := range []string{"202605282130.0", "202605282130.99", "202613282130.00", "202602302130.00"} {
		if _, err := parseTouchTime(s); err == nil {
			t.Fatalf("invalid touch time %q accepted", s)
		}
	}
}

func TestHelpTextContainsFinalContract(t *testing.T) {
	needles := []string{
		"-s  --sync <FILE>",
		"-e  --edit <OFFSET>",
		"-r  --runas",
		"-e  [+|-][[hh]mm]SS",
		"Default behavior updates all three Windows timestamps",
	}
	for _, needle := range needles {
		if !strings.Contains(HelpText, needle) {
			t.Fatalf("help missing %q", needle)
		}
	}
}

func TestParseDateUTC(t *testing.T) {
	tm, err := parseDate("2026-05-28T21:30:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if tm.Location() != time.UTC {
		// time.ParseInLocation returns UTC location for UTC input.
		t.Fatalf("expected UTC location, got %v", tm.Location())
	}
}
