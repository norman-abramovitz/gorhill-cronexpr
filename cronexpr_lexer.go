//go:build parser_rd || parser_lalr

package cronexpr

import (
	"fmt"
	"strings"
)

// The tokens of a field. Punctuation is its own byte: '*', '?', '-', '/',
// ',', '(', ')', '#'.
const (
	tokEnd    = 0
	tokNumber = 256 + iota // `15`
	tokName                // `jan`, `mon`
	tokL                   // `L`
	tokW                   // `W`
	tokLW                  // `LW`
	tokH                   // `H`
	tokBad                 // anything else
)

type fieldToken struct {
	kind     int
	text     string
	beg, end int // place in the field
}

// fieldLexer splits a lower-case field into tokens. A run of letters is a
// name of the field, L, W, LW or H; a name may carry an L or W on its end
// (`fril`, `julw`), which comes back as a second token.
type fieldLexer struct {
	s      string
	pos    int
	names  map[string]int
	queued fieldToken
}

func newFieldLexer(s string, desc *fieldDescriptor) fieldLexer {
	lex := fieldLexer{s: s}
	switch desc.name {
	case "month":
		lex.names = monthTokens
	case "day-of-week":
		lex.names = dowTokens
	}
	return lex
}

func (lex *fieldLexer) next() fieldToken {
	if lex.queued.kind != tokEnd {
		t := lex.queued
		lex.queued = fieldToken{}
		return t
	}
	beg := lex.pos
	if beg == len(lex.s) {
		return fieldToken{kind: tokEnd, beg: beg, end: beg}
	}
	end := beg + 1
	token := func(kind int) fieldToken {
		lex.pos = end
		return fieldToken{kind: kind, text: lex.s[beg:end], beg: beg, end: end}
	}
	switch c := lex.s[beg]; {
	case '0' <= c && c <= '9':
		for end < len(lex.s) && '0' <= lex.s[end] && lex.s[end] <= '9' {
			end++
		}
		return token(tokNumber)
	case 'a' <= c && c <= 'z':
		for end < len(lex.s) && 'a' <= lex.s[end] && lex.s[end] <= 'z' {
			end++
		}
		switch word := lex.s[beg:end]; {
		case word == "l":
			return token(tokL)
		case word == "w":
			return token(tokW)
		case word == "lw":
			return token(tokLW)
		case word == "h":
			return token(tokH)
		case lex.isName(word):
			return token(tokName)
		case lex.isName(word[:len(word)-1]) && (word[len(word)-1] == 'l' || word[len(word)-1] == 'w'):
			kind := tokL
			if word[len(word)-1] == 'w' {
				kind = tokW
			}
			lex.queued = fieldToken{kind: kind, text: word[len(word)-1:], beg: end - 1, end: end}
			lex.pos = end
			return fieldToken{kind: tokName, text: word[:len(word)-1], beg: beg, end: end - 1}
		}
		return token(tokBad)
	case strings.IndexByte("*?-/,()#", c) >= 0:
		return token(int(c))
	}
	return token(tokBad)
}

func (lex *fieldLexer) isName(word string) bool {
	_, found := lex.names[word]
	return found
}

/******************************************************************************/

// validEntry checks the values of an entry against its field, as the
// field's value pattern does in the layouts: a name of the field, or a
// number from min to max in at most two digits (a year: any digits).
func validEntry(entry *fieldEntry, desc *fieldDescriptor) bool {
	value := func(text string) bool {
		if text[0] >= 'a' {
			return true // the lexer only lets the field's names through
		}
		if desc.name == "year" {
			return true
		}
		return len(text) <= 2 && desc.min <= number(text) && number(text) <= desc.max
	}
	switch entry.form {
	case formValue, formValueStep, formWorkdom, formDowOfLastWeek:
		return value(entry.first)
	case formRange, formRangeStep, formHashRange, formHashRangeStep:
		return value(entry.first) && value(entry.last)
	case formDowOfSpecificWeek:
		return value(entry.first) && len(entry.last) == 1 && '1' <= entry.last[0] && entry.last[0] <= '5'
	case formLastDomOffset:
		return len(entry.first) <= 2 && 1 <= number(entry.first) && number(entry.first) <= 30
	}
	return true
}

// number reads at most two digits. atoi will not do: it looks numbers up,
// and one it lacks (60) reads as 0.
func number(digits string) int {
	v := 0
	for i := 0; i < len(digits) && i < 2; i++ {
		v = v*10 + int(digits[i]-'0')
	}
	return v
}

// entrySyntaxError reports the entry of field s that holds position pos.
func entrySyntaxError(s string, pos int, desc *fieldDescriptor) error {
	beg := strings.LastIndexByte(s[:pos], ',') + 1
	end := strings.IndexByte(s[pos:], ',')
	if end < 0 {
		end = len(s)
	} else {
		end += pos
	}
	return fmt.Errorf("syntax error in %s field: '%s'", desc.name, s[beg:end])
}
