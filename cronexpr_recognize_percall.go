//go:build parser_percall

package cronexpr

// recognizeField compiles the field's layouts each time it reads a field,
// to show what compiling them once at startup saves.
func recognizeField(s string, desc *fieldDescriptor, entries []fieldEntry) ([]fieldEntry, error) {
	return regexRecognizeField(s, compileLayouts(desc, nil), entries), nil
}
