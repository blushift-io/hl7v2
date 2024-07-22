package query

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	eof      = rune(0)
	leftArr  = '['
	rightArr = ']'
	itemSep  = '.'

	capLetters  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits      = "0123456789"
	beginDigits = "123456789"
	segChars    = capLetters + digits
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=itemType -linecomment
type itemType int

const (
	itemErr    itemType = iota //error
	itemSeg                    //segment
	itemSegIdx                 //segment_index
	itemField                  //field
	itemRep                    //repetition
	itemComp                   //component
	itemSub                    //subcomponent
	itemEOF                    //eof
)

func (t itemType) isIdx() bool {
	return t == itemSegIdx || t == itemRep
}

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	if i.typ == itemEOF {
		return "EOF"
	}

	return i.val
}

type stateFn func(*lexer) stateFn
type lexer struct {
	name  string
	input string
	state stateFn
	on    itemType
	start int
	pos   int
	width int
	items chan item
}

func lex(name, input string) *lexer {
	return &lexer{
		name:  name,
		input: input,
		state: lexMain,
		on:    itemSeg,
		items: make(chan item, 2),
	}
}

func (l *lexer) nextItem() item {
	for {
		select {
		case i := <-l.items:
			return i
		default:
			l.state = l.state(l)
		}
	}
}

func (l *lexer) run() {
	for state := lexMain; state != nil; {
		state = state(l)
	}

	close(l.items)
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
	l.on++
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}

	l.backup()

	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}

	l.backup()
}

func (l *lexer) errorf(f string, args ...interface{}) stateFn {
	l.items <- item{itemErr, fmt.Sprintf(f, args...)}
	return nil
}

func lexMain(l *lexer) stateFn {
	for r := l.next(); r != eof; {
		if unicode.IsSpace(r) {
			l.ignore()
			continue
		}

		if l.on >= itemEOF {
			return l.errorf("invalid location query: %s", l.input)
		}

		switch {
		case strings.ContainsRune(capLetters, r):
			return lexSeg
		case r == leftArr:
			l.ignore()
			return lexIdx
		case r == itemSep:
			l.ignore()
			return lexPos
		}

		return l.errorf("invalid location query, unexpected character at %d: %s", l.pos, l.input)
	}

	return lexEnd
}

func lexSeg(l *lexer) stateFn {
	l.acceptRun(segChars)
	if l.pos-l.start > 3 {
		return l.errorf("invalid segment name")
	}

	l.emit(itemSeg)
	return lexMain
}

func lexIdx(l *lexer) stateFn {
	if !l.on.isIdx() {
		return l.errorf("item type %v is not indexable", l.on)
	}

	if !l.accept(digits) {
		return l.errorf("invalid index at %d: %s", l.pos, l.input)
	}

	l.acceptRun(digits)

	if l.peek() != rightArr {
		return l.errorf("expected index close at %d: %s", l.pos, l.input)
	}

	l.emit(l.on)
	l.next()
	l.ignore()

	return lexMain
}

func lexPos(l *lexer) stateFn {
	if l.on.isIdx() {
		l.on++
	}

	if l.accept("0") {
		l.emit(l.on)
		return lexMain
	}

	if !l.accept(beginDigits) {
		return l.errorf("invalid position at %d: %s", l.pos, l.input)
	}

	l.acceptRun(digits)
	l.emit(l.on)

	return lexMain
}

func lexEnd(l *lexer) stateFn {
	l.emit(itemEOF)
	return nil
}
