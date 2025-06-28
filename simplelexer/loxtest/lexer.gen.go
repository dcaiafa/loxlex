package loxtest

var _lexerMode0 = []uint32{
	4, 22, 36, 43, 17, 0, 5, 9, 10, 1, 13, 13, 1, 32,
	32, 1, 34, 34, 2, 48, 57, 3, 13, 0, 3, 9, 10, 1,
	13, 13, 1, 32, 32, 1, 4, 0, 6, 0, 0, 1, 1, 5,
	0, 7, 0, 1, 48, 57, 3, 3, 2,
}

var _lexerMode1 = []uint32{
	22, 40, 45, 40, 69, 40, 81, 40, 93, 105, 40, 117, 129, 141,
	153, 160, 172, 184, 196, 208, 220, 232, 17, 0, 5, 32, 33, 10,
	34, 34, 14, 35, 91, 10, 92, 92, 2, 93, 1114111, 10, 4, 0,
	0, 5, 0, 23, 0, 7, 34, 34, 7, 85, 85, 4, 110, 110,
	7, 114, 114, 7, 116, 116, 7, 117, 117, 19, 120, 120, 17, 11,
	0, 3, 48, 57, 21, 65, 70, 21, 97, 102, 21, 11, 0, 3,
	48, 57, 1, 65, 70, 1, 97, 102, 1, 11, 0, 3, 48, 57,
	6, 65, 70, 6, 97, 102, 6, 11, 0, 3, 48, 57, 3, 65,
	70, 3, 97, 102, 3, 11, 0, 3, 48, 57, 8, 65, 70, 8,
	97, 102, 8, 11, 0, 3, 48, 57, 9, 65, 70, 9, 97, 102,
	9, 11, 0, 3, 48, 57, 5, 65, 70, 5, 97, 102, 5, 6,
	0, 0, 2, 0, 3, 3, 11, 0, 3, 48, 57, 11, 65, 70,
	11, 97, 102, 11, 11, 0, 3, 48, 57, 12, 65, 70, 12, 97,
	102, 12, 11, 0, 3, 48, 57, 13, 65, 70, 13, 97, 102, 13,
	11, 0, 3, 48, 57, 15, 65, 70, 15, 97, 102, 15, 11, 0,
	3, 48, 57, 16, 65, 70, 16, 97, 102, 16, 11, 0, 3, 48,
	57, 18, 65, 70, 18, 97, 102, 18, 11, 0, 3, 48, 57, 20,
	65, 70, 20, 97, 102, 20,
}

var _lexerModes = [][]uint32{

	_lexerMode0,

	_lexerMode1,
}

const _stateNonGreedyAccepting = 1

const (
	_lexerConsume  = 0
	_lexerAccept   = 1
	_lexerDiscard  = 2
	_lexerTryAgain = 3
	_lexerEOF      = 4
	_lexerError    = -1
)

type _LexerStateMachine struct {
	token     int
	state     int
	mode      []uint32
	modeStack _Stack[[]uint32]
}

func (l *_LexerStateMachine) PushRune(r rune) int {
	if l.mode == nil {
		l.mode = _lexerMode0
	}

	mode := l.mode

	// Find the table row corresponding to state.
	i := int(mode[int(l.state)])
	count := int(mode[i])
	i++
	end := i + count

	// The format of each row is as follows:
	//
	//   stateFlags uint32
	//   gotoCount uint32
	//   [gotoCount]struct{
	//     rangeBegin uint32
	//     rangeEnd   uint32
	//     gotoState  uint32
	//   }
	//   [actionCount]struct {
	//     actionType  uint32
	//     actionParam uint32
	//   }
	//
	// Where 'actionCount' is determined by the amount of uint32 left in the row.

	flags := mode[i]
	gotoN := int(mode[i+1])
	i += 2

	if flags&_stateNonGreedyAccepting == 0 {
		// Use binary-search to find the next state.
		b := 0
		e := gotoN
		for b < e {
			j := b + (e-b)/2
			k := i + j*3
			switch {
			case r >= rune(mode[k]) && r <= rune(mode[k+1]):
				l.state = int(mode[k+2])
				return _lexerConsume
			case r < rune(mode[k]):
				e = j
			case r > rune(mode[k+1]):
				b = j + 1
			default:
				panic("not reached")
			}
		}
	}

	// Move 'i' to the beginning of the actions section.
	i += gotoN * 3

	for ; i < end; i += 2 {
		switch mode[i] {
		case 1: // PushMode
			modeIndex := int(mode[i+1])
			l.modeStack.Push(mode)
			l.mode = _lexerModes[modeIndex]
		case 2: // PopMode
			if len(l.modeStack) == 0 {
				return _lexerError
			}
			l.mode = l.modeStack.Peek(0)
			l.modeStack.Pop(1)
		case 3: // Accept
			l.token = int(mode[i+1])
			l.state = 0
			return _lexerAccept
		case 4: // Discard
			l.state = 0
			return _lexerDiscard
		case 5: // Accum
			l.state = 0
			return _lexerTryAgain
		}
	}

	if l.state == 0 && r == -1 {
		return _lexerEOF
	}

	return _lexerError
}

func (l *_LexerStateMachine) Reset() {
	l.mode = nil
	l.state = 0
}

func (l *_LexerStateMachine) Token() int {
	return l.token
}
