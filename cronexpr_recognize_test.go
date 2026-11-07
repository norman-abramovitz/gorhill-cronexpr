//go:build parser_rd || parser_lalr

package cronexpr

import (
	"slices"
	"strings"
	"testing"
)

// FuzzRecognizeField checks the recognizer of this build against the
// layouts: both accept the same fields and find the same entries.
func FuzzRecognizeField(f *testing.F) {
	for _, seed := range []string{
		"*", "?/5", "5", "05", "60", "5-20", "22-2/3", "*/0", "1,,2", ",", "",
		"h", "h/15", "h(0-29)", "h(0-29)/10", "h(", "l", "l-3", "l-31", "lw",
		"15w", "5l", "5#l", "5#3", "5#6", "fril", "fri#2", "april", "julw",
		"jan-mar/2", "mon-fri", "sun,sat", "2026", "1970-2099/10", "5-", "-5",
		"5//2", "l-", "w", "x", "5lw", "MON", "Jan", "1\v2", "é",
	} {
		for field := range fieldDescriptors {
			f.Add(byte(field), seed)
		}
	}
	f.Fuzz(func(t *testing.T, field byte, s string) {
		desc := fieldDescriptors[int(field)%len(fieldDescriptors)]
		want := regexRecognizeField(s, desc.layouts, nil)
		wantOK := !slices.ContainsFunc(want, func(e fieldEntry) bool { return e.form == formNone })
		got, err := recognizeField(s, desc, nil)
		if wantOK != (err == nil) {
			t.Fatalf("%s %q: layouts accept %v, recognizer error %v", desc.name, s, wantOK, err)
		}
		if wantOK && !slices.Equal(got, want) {
			t.Fatalf("%s %q:\nlayouts    %+v\nrecognizer %+v", desc.name, s, want, got)
		}
	})
}

// TestRecognizeFieldUpperCase checks a field in capitals reads as in lower case.
func TestRecognizeFieldUpperCase(t *testing.T) {
	for _, s := range []string{"MON-FRI", "JAN,JUL", "LW", "5L", "H(0-29)/10"} {
		got, err := recognizeField(s, &dowDescriptor, nil)
		want, _ := recognizeField(strings.ToLower(s), &dowDescriptor, nil)
		if !slices.Equal(got, want) || (err == nil) != (want != nil) {
			t.Errorf("%q: %+v, %v; lower case %+v", s, got, err, want)
		}
	}
}
