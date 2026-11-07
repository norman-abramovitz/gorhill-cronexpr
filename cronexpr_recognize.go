//go:build !parser_percall && !parser_rd && !parser_lalr

package cronexpr

// recognizeField finds the form of each entry of a field. This build uses
// the layouts compiled at startup. The tags parser_percall, parser_rd and
// parser_lalr pick another recognizer, to compare them.
func recognizeField(s string, desc *fieldDescriptor, entries []fieldEntry) ([]fieldEntry, error) {
	return regexRecognizeField(s, desc.layouts, entries), nil
}
