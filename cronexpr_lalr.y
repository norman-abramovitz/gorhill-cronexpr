// The grammar of a field, read bottom up by an LALR(1) parser that goyacc
// builds from this file:
//
//	go generate -tags parser_lalr .
//
// writes cronexpr_lalr.go. The parser shifts tokens onto a stack until the
// top of it matches the right side of a rule, then reduces them to the
// rule's left side and runs the rule's action.

%{
//go:build parser_lalr

package cronexpr
%}

%union {
	tok   fieldToken
	entry fieldEntry
}

%token <tok> NUMBER NAME LAST WEEKDAY LASTWEEKDAY HASH BAD
%token <tok> '*' '?' '-' '/' ',' '(' ')' '#'

%type <entry> entry
%type <tok> value wildcard

%%

// A list is left-recursive: each entry is reduced as soon as its comma
// arrives, so the stack never holds more than one entry.
field:
	item
|	field ',' item

item:
	/* empty, as in `1,,2` */
|	entry
	{
		lalrlex.(*lalrFieldLexer).add($1)
	}

entry:
	wildcard
	{ $$ = entryOf(formWildcard, $1, $1, "", "", "") }
|	wildcard '/' NUMBER
	{ $$ = entryOf(formWildcardStep, $1, $3, "", "", $3.text) }
|	value
	{ $$ = entryOf(formValue, $1, $1, $1.text, "", "") }
|	value '-' value
	{ $$ = entryOf(formRange, $1, $3, $1.text, $3.text, "") }
|	value '/' NUMBER
	{ $$ = entryOf(formValueStep, $1, $3, $1.text, "", $3.text) }
|	value '-' value '/' NUMBER
	{ $$ = entryOf(formRangeStep, $1, $5, $1.text, $3.text, $5.text) }
|	HASH
	{ $$ = entryOf(formHash, $1, $1, "", "", "") }
|	HASH '/' NUMBER
	{ $$ = entryOf(formHashStep, $1, $3, "", "", $3.text) }
|	HASH '(' value '-' value ')'
	{ $$ = entryOf(formHashRange, $1, $6, $3.text, $5.text, "") }
|	HASH '(' value '-' value ')' '/' NUMBER
	{ $$ = entryOf(formHashRangeStep, $1, $8, $3.text, $5.text, $8.text) }
|	LAST
	{ $$ = entryOf(formLastDom, $1, $1, "", "", "") }
|	LAST '-' NUMBER
	{ $$ = entryOf(formLastDomOffset, $1, $3, $3.text, "", "") }
|	LASTWEEKDAY
	{ $$ = entryOf(formLastWorkdom, $1, $1, "", "", "") }
|	value WEEKDAY
	{ $$ = entryOf(formWorkdom, $1, $2, $1.text, "", "") }
|	value LAST
	{ $$ = entryOf(formDowOfLastWeek, $1, $2, $1.text, "", "") }
|	value '#' LAST
	{ $$ = entryOf(formDowOfLastWeek, $1, $3, $1.text, "", "") }
|	value '#' NUMBER
	{ $$ = entryOf(formDowOfSpecificWeek, $1, $3, $1.text, $3.text, "") }

wildcard:
	'*'
|	'?'

value:
	NUMBER
|	NAME

%%
