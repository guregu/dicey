package dicey

import (
	"fmt"
	"unicode/utf8"
)

// Cribbed from "Lexical Scanning in Go" - Rob Pike
// https://www.youtube.com/watch?v=HxaD_trXwRE
// and the standard library (text/template/parse)

// itemType of a lexed item.
type itemType int

// Types of lexed items.
const (
	itemError itemType = iota
	itemEOF
	itemDice
	itemBonus
	itemAdd
	itemSub
)

// Item is a lexed item.
type item struct {
	Type itemType
	Pos  int
	Val  string
}

func (i item) String() string {
	// switch i.Type {
	// case itemEOF:
	// 	return "EOF"
	// case ItemNamePlaceholder:
	// 	return "$"
	// case ItemValuePlaceholder:
	// 	return "?"
	// }
	return i.Val
}

const (
	eof = -1
)

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	state stateFn
	start int
	pos   int
	width int
	items chan item
}

func (l *lexer) next() rune {
	if (l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) emit(t itemType) {
	l.items <- item{
		Type: t,
		Pos:  l.pos,
		Val:  l.input[l.start:l.pos],
	}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// accepts

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexNumber; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func lexNumber(l *lexer) stateFn {
	var nextFn stateFn
loop:
	for {
		r := l.next()
		switch r {
		case eof:
			break loop
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		case 'd':
			return lexDie
		case '+':
			l.backup()
			nextFn = lexAdd
			break loop
		case '-':
			l.backup()
			nextFn = lexSub
			break loop
		default:
			l.emit(itemError)
			return nil
		}
	}
	if l.pos > l.start {
		l.emit(itemBonus)
	}
	if nextFn == nil {
		l.emit(itemEOF)
	}
	return nextFn
}

func lexDie(l *lexer) stateFn {
loop:
	for {
		r := l.next()
		switch r {
		case eof:
			break loop
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		case '+':
			l.backup()
			l.emit(itemDice)
			return lexAdd
		case '-':
			l.backup()
			l.emit(itemDice)
			return lexSub
		default:
			l.emit(itemError)
			return nil
		}
	}
	if l.pos > l.start {
		l.emit(itemDice)
	}
	l.emit(itemEOF)
	return nil
}

func lexAdd(l *lexer) stateFn {
	l.next()
	l.emit(itemAdd)
	return lexNumber
}

func lexSub(l *lexer) stateFn {
	l.next()
	l.emit(itemSub)
	return lexNumber
}
