//go:build parser_rd

package cronexpr

import "strings"

// recognizeField reads a field top down, by recursive descent: each choice
// is made by looking at the next token only.
//
//	field = [ entry ] { "," [ entry ] }
//	entry = ( "*" | "?" ) [ "/" NUMBER ]
//	      | "H" [ "(" value "-" value ")" ] [ "/" NUMBER ]
//	      | "L" [ "-" NUMBER ]
//	      | "LW"
//	      | value [ "-" value [ "/" NUMBER ]
//	              | "/" NUMBER
//	              | "W"
//	              | "L"
//	              | "#" ( "L" | NUMBER ) ]
//	value = NUMBER | NAME
//
// The entries that start with a value share it, so the value is read once
// and the token after it picks the form (the grammar is left-factored).
func recognizeField(s string, desc *fieldDescriptor, entries []fieldEntry) ([]fieldEntry, error) {
	s = strings.ToLower(s)
	p := rdParser{lex: newFieldLexer(s, desc)}
	p.advance()
	for {
		if p.tok.kind != ',' && p.tok.kind != tokEnd {
			entry := p.entry()
			if p.bad || !validEntry(&entry, desc) {
				return nil, entrySyntaxError(s, p.tok.beg, desc)
			}
			entries = append(entries, entry)
		}
		switch p.tok.kind {
		case tokEnd:
			return entries, nil
		case ',':
			p.advance()
		default:
			return nil, entrySyntaxError(s, p.tok.beg, desc)
		}
	}
}

type rdParser struct {
	lex fieldLexer
	tok fieldToken // the next token, not yet read
	end int        // where the last token read ends
	bad bool       // a token was not the one expected
}

// advance reads the next token and returns its text.
func (p *rdParser) advance() string {
	t := p.tok
	p.end = t.end
	p.tok = p.lex.next()
	return t.text
}

// accept reads the next token when it is kind.
func (p *rdParser) accept(kind int) bool {
	if p.tok.kind != kind {
		return false
	}
	p.advance()
	return true
}

// expect reads the next token, which must be kind.
func (p *rdParser) expect(kind int) string {
	if p.tok.kind != kind {
		p.bad = true
		return ""
	}
	return p.advance()
}

func (p *rdParser) value() string {
	if p.tok.kind != tokNumber && p.tok.kind != tokName {
		p.bad = true
		return ""
	}
	return p.advance()
}

func (p *rdParser) entry() fieldEntry {
	entry := fieldEntry{sbeg: p.tok.beg}
	switch p.tok.kind {
	case '*', '?':
		p.advance()
		entry.form = formWildcard
		if p.accept('/') {
			entry.form = formWildcardStep
			entry.step = p.expect(tokNumber)
		}
	case tokH:
		p.advance()
		entry.form = formHash
		if p.accept('(') {
			entry.form = formHashRange
			entry.first = p.value()
			p.expect('-')
			entry.last = p.value()
			p.expect(')')
		}
		if p.accept('/') {
			entry.form++ // formHashStep or formHashRangeStep
			entry.step = p.expect(tokNumber)
		}
	case tokL:
		p.advance()
		entry.form = formLastDom
		if p.accept('-') {
			entry.form = formLastDomOffset
			entry.first = p.expect(tokNumber)
		}
	case tokLW:
		p.advance()
		entry.form = formLastWorkdom
	default:
		entry.first = p.value()
		entry.form = formValue
		switch {
		case p.accept('-'):
			entry.form = formRange
			entry.last = p.value()
			if p.accept('/') {
				entry.form = formRangeStep
				entry.step = p.expect(tokNumber)
			}
		case p.accept('/'):
			entry.form = formValueStep
			entry.step = p.expect(tokNumber)
		case p.accept(tokW):
			entry.form = formWorkdom
		case p.accept(tokL):
			entry.form = formDowOfLastWeek
		case p.accept('#'):
			if p.accept(tokL) {
				entry.form = formDowOfLastWeek
			} else {
				entry.form = formDowOfSpecificWeek
				entry.last = p.expect(tokNumber)
			}
		}
	}
	entry.send = p.end
	return entry
}
