package loxtest

import (
	_i0 "github.com/dcaiafa/lox_lexer/simplelexer"
)

var _rules = []int32{
	0, 1, 2, 2, 2, 3, 3, 4, 4,
}

var _termCounts = []int32{
	1, 1, 1, 1, 1, 1, 0, 2, 1,
}

var _actions = []int32{
	9, 18, 27, 36, 39, 48, 57, 60, 69, 8, 0, -6, 1, 1,
	2, 2, 3, 4, 8, 0, -4, 1, -4, 2, -4, 3, -4, 8,
	0, -2, 1, -2, 2, -2, 3, -2, 2, 0, 2147483647, 8, 0, -3,
	1, -3, 2, -3, 3, -3, 8, 0, -8, 1, -8, 2, -8, 3,
	-8, 2, 0, -1, 8, 0, -5, 1, 1, 2, 2, 3, 4, 8,
	0, -7, 1, -7, 2, -7, 3, -7,
}

var _goto = []int32{
	9, 18, 18, 18, 18, 18, 18, 19, 18, 8, 1, 3, 2, 5,
	3, 6, 4, 7, 0, 2, 2, 8,
}

type _Bounds struct {
	Begin Token
	End   Token
	Empty bool
}

func _cast[T any](v any) T {
	cv, _ := v.(T)
	return cv
}

type Error struct {
	Token    Token
	Expected []int
}

func _Find(table []int32, y, x int32) (int32, bool) {
	i := int(table[int(y)])
	count := int(table[i])
	i++
	end := i + count
	for ; i < end; i += 2 {
		if table[i] == x {
			return table[i+1], true
		}
	}
	return 0, false
}

type _Lexer interface {
	ReadToken() (Token, int)
}

type _item struct {
	State int32
	Sym   any
}

type lox struct {
	_lex   _Lexer
	_stack _Stack[_item]

	_la    int
	_lasym any

	_qla    int
	_qlasym any
}

func (p *parser) parse(lex _Lexer) bool {
	const accept = 2147483647

	p._lex = lex
	p._qla = -1
	p._stack.Push(_item{})

	p._readToken()

	for {
		topState := p._stack.Peek(0).State
		action, ok := _Find(_actions, topState, int32(p._la))
		if !ok {
			if !p._recover() {
				return false
			}
			continue
		}
		if action == accept {
			break
		} else if action >= 0 { // shift
			p._stack.Push(_item{
				State: action,
				Sym:   p._lasym,
			})
			p._readToken()
		} else { // reduce
			prod := -action
			termCount := _termCounts[int(prod)]
			rule := _rules[int(prod)]
			res := p._act(prod)
			p._stack.Pop(int(termCount))
			topState = p._stack.Peek(0).State
			nextState, _ := _Find(_goto, topState, rule)
			p._stack.Push(_item{
				State: nextState,
				Sym:   res,
			})
		}
	}

	return true
}

// recoverLookahead can be called during an error production action (an action
// for a production that has a @error term) to recover the lookahead that was
// possibly lost in the process of reducing the error production.
func (p *parser) recoverLookahead(typ int, tok Token) {
	if p._qla != -1 {
		panic("recovered lookahead already pending")
	}

	p._qla = p._la
	p._qlasym = p._lasym
	p._la = typ
	p._lasym = tok
}

func (p *parser) _readToken() {
	if p._qla != -1 {
		p._la = p._qla
		p._lasym = p._qlasym
		p._qla = -1
		p._qlasym = nil
		return
	}

	p._lasym, p._la = p._lex.ReadToken()
	if p._la == ERROR {
		p._lasym = p._makeError()
	}
}

func (p *parser) _recover() bool {
	errSym, ok := p._lasym.(Error)
	if !ok {
		errSym = p._makeError()
	}

	for p._la == ERROR {
		p._readToken()
	}

	for {
		save := p._stack

		for len(p._stack) >= 1 {
			state := p._stack.Peek(0).State

			for {
				action, ok := _Find(_actions, state, int32(ERROR))
				if !ok {
					break
				}

				if action < 0 {
					prod := -action
					rule := _rules[int(prod)]
					state, _ = _Find(_goto, state, rule)
					continue
				}

				state = action

				_, ok = _Find(_actions, state, int32(p._la))
				if !ok {
					break
				}

				p._qla = p._la
				p._qlasym = p._lasym
				p._la = ERROR
				p._lasym = errSym
				return true
			}

			p._stack.Pop(1)
		}

		if p._la == EOF {
			return false
		}

		p._stack = save
		p._readToken()
	}
}

func (p *parser) _makeError() Error {
	e := Error{
		Token: p._lasym.(Token),
	}

	// Compile list of allowed tokens at this state.
	// See _Find for the format of the _actions table.
	s := p._stack.Peek(0).State
	i := int(_actions[int(s)])
	count := int(_actions[i])
	i++
	end := i + count
	for ; i < end; i += 2 {
		e.Expected = append(e.Expected, int(_actions[i]))
	}

	return e
}

func (p *parser) _act(prod int32) any {
	switch prod {
	case 1:
		return p.on_S(
			_cast[[]_i0.Token](p._stack.Peek(0).Sym),
		)
	case 2:
		return p.on_token(
			_cast[_i0.Token](p._stack.Peek(0).Sym),
		)
	case 3:
		return p.on_token(
			_cast[_i0.Token](p._stack.Peek(0).Sym),
		)
	case 4:
		return p.on_token__err(
			_cast[Error](p._stack.Peek(0).Sym),
		)
	case 5: // ZeroOrMore
		return _cast[[]_i0.Token](p._stack.Peek(0).Sym)
	case 6: // ZeroOrMore
		{
			var zero []_i0.Token
			return zero
		}
	case 7: // OneOrMore
		return append(
			_cast[[]_i0.Token](p._stack.Peek(1).Sym),
			_cast[_i0.Token](p._stack.Peek(0).Sym),
		)
	case 8: // OneOrMore
		return []_i0.Token{
			_cast[_i0.Token](p._stack.Peek(0).Sym),
		}
	default:
		panic("unreachable")
	}
}
