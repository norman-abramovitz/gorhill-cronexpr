//go:build parser_lalr

package cronexpr

import (
	"strings"
	"sync"
)

//go:generate go tool goyacc -l -p lalr -v /dev/null -o cronexpr_lalr.go cronexpr_lalr.y

// recognizeField reads a field bottom up with the LALR(1) parser built from
// cronexpr_lalr.y.
func recognizeField(s string, desc *fieldDescriptor, entries []fieldEntry) ([]fieldEntry, error) {
	s = strings.ToLower(s)
	lex := &lalrFieldLexer{fieldLexer: newFieldLexer(s, desc), desc: desc, entries: entries}
	parser := lalrParsers.Get().(lalrParser)
	parser.Parse(lex)
	lalrParsers.Put(parser)
	if lex.failed {
		return nil, entrySyntaxError(s, lex.errPos, desc)
	}
	return lex.entries, nil
}

// lalrParsers keeps parsers to use again. goyacc's parser holds its stack:
// 16 symbols of about 140 bytes each, as a %union is a struct holding every
// member. A new one for each field was two thirds of the bytes Parse
// allocated, and most of its time went to collecting them.
var lalrParsers = sync.Pool{New: func() any { return lalrNewParser() }}

// lalrFieldLexer hands the parser its tokens and keeps the entries the
// grammar's actions find.
type lalrFieldLexer struct {
	fieldLexer
	desc    *fieldDescriptor
	entries []fieldEntry
	last    fieldToken // the token handed over last
	failed  bool
	errPos  int
}

func (lex *lalrFieldLexer) Lex(lval *lalrSymType) int {
	t := lex.next()
	lex.last = t
	lval.tok = t
	switch t.kind {
	case tokNumber:
		return NUMBER
	case tokName:
		return NAME
	case tokL:
		return LAST
	case tokW:
		return WEEKDAY
	case tokLW:
		return LASTWEEKDAY
	case tokH:
		return HASH
	case tokBad:
		return BAD
	}
	return t.kind // punctuation is its own byte; the end is 0
}

func (lex *lalrFieldLexer) Error(string) {
	lex.fail(lex.last.beg)
}

func (lex *lalrFieldLexer) fail(pos int) {
	if !lex.failed {
		lex.failed, lex.errPos = true, pos
	}
}

func (lex *lalrFieldLexer) add(entry fieldEntry) {
	if !validEntry(&entry, lex.desc) {
		lex.fail(entry.sbeg)
	}
	lex.entries = append(lex.entries, entry)
}

// entryOf builds an entry that runs from token beg to token end.
func entryOf(form int, beg, end fieldToken, first, last, step string) fieldEntry {
	return fieldEntry{form: form, first: first, last: last, step: step, sbeg: beg.beg, send: end.end}
}
