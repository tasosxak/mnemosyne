package internals

import (
	"strconv"
)
import (
	"bufio"
	"io"
	"strings"
)

type frame struct {
	i            int
	s            string
	line, column int
}
type Lexer struct {
	// The lexer runs in its own goroutine, and communicates via channel 'ch'.
	ch      chan frame
	ch_stop chan bool
	// We record the level of nesting because the action could return, and a
	// subsequent call expects to pick up where it left off. In other words,
	// we're simulating a coroutine.
	// TODO: Support a channel-based variant that compatible with Go's yacc.
	stack []frame
	stale bool

	// The 'l' and 'c' fields were added for
	// https://github.com/wagerlabs/docker/blob/65694e801a7b80930961d70c69cba9f2465459be/buildfile.nex
	// Since then, I introduced the built-in Line() and Column() functions.
	l, c int

	parseResult interface{}

	// The following line makes it easy for scripts to insert fields in the
	// generated code.
	// [NEX_END_OF_LEXER_STRUCT]
}

// NewLexerWithInit creates a new Lexer object, runs the given callback on it,
// then returns it.
func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer {
	yylex := new(Lexer)
	if initFun != nil {
		initFun(yylex)
	}
	yylex.ch = make(chan frame)
	yylex.ch_stop = make(chan bool, 1)
	var scan func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int)
	scan = func(in *bufio.Reader, ch chan frame, ch_stop chan bool, family []dfa, line, column int) {
		// Index of DFA and length of highest-precedence match so far.
		matchi, matchn := 0, -1
		var buf []rune
		n := 0
		checkAccept := func(i int, st int) bool {
			// Higher precedence match? DFAs are run in parallel, so matchn is at most len(buf), hence we may omit the length equality check.
			if family[i].acc[st] && (matchn < n || matchi > i) {
				matchi, matchn = i, n
				return true
			}
			return false
		}
		var state [][2]int
		for i := 0; i < len(family); i++ {
			mark := make([]bool, len(family[i].startf))
			// Every DFA starts at state 0.
			st := 0
			for {
				state = append(state, [2]int{i, st})
				mark[st] = true
				// As we're at the start of input, follow all ^ transitions and append to our list of start states.
				st = family[i].startf[st]
				if -1 == st || mark[st] {
					break
				}
				// We only check for a match after at least one transition.
				checkAccept(i, st)
			}
		}
		atEOF := false
		stopped := false
		for {
			if n == len(buf) && !atEOF {
				r, _, err := in.ReadRune()
				switch err {
				case io.EOF:
					atEOF = true
				case nil:
					buf = append(buf, r)
				default:
					panic(err)
				}
			}
			if !atEOF {
				r := buf[n]
				n++
				var nextState [][2]int
				for _, x := range state {
					x[1] = family[x[0]].f[x[1]](r)
					if -1 == x[1] {
						continue
					}
					nextState = append(nextState, x)
					checkAccept(x[0], x[1])
				}
				state = nextState
			} else {
			dollar: // Handle $.
				for _, x := range state {
					mark := make([]bool, len(family[x[0]].endf))
					for {
						mark[x[1]] = true
						x[1] = family[x[0]].endf[x[1]]
						if -1 == x[1] || mark[x[1]] {
							break
						}
						if checkAccept(x[0], x[1]) {
							// Unlike before, we can break off the search. Now that we're at the end, there's no need to maintain the state of each DFA.
							break dollar
						}
					}
				}
				state = nil
			}

			if state == nil {
				lcUpdate := func(r rune) {
					if r == '\n' {
						line++
						column = 0
					} else {
						column++
					}
				}
				// All DFAs stuck. Return last match if it exists, otherwise advance by one rune and restart all DFAs.
				if matchn == -1 {
					if len(buf) == 0 { // This can only happen at the end of input.
						break
					}
					lcUpdate(buf[0])
					buf = buf[1:]
				} else {
					text := string(buf[:matchn])
					buf = buf[matchn:]
					matchn = -1
					select {
					case ch <- frame{matchi, text, line, column}:
						{
						}
					case stopped = <-ch_stop:
						{
						}
					}
					if stopped {
						break
					}
					if len(family[matchi].nest) > 0 {
						scan(bufio.NewReader(strings.NewReader(text)), ch, ch_stop, family[matchi].nest, line, column)
					}
					if atEOF {
						break
					}
					for _, r := range text {
						lcUpdate(r)
					}
				}
				n = 0
				for i := 0; i < len(family); i++ {
					state = append(state, [2]int{i, 0})
				}
			}
		}
		ch <- frame{-1, "", line, column}
	}
	go scan(bufio.NewReader(in), yylex.ch, yylex.ch_stop, dfas, 0, 0)
	return yylex
}

type dfa struct {
	acc          []bool           // Accepting states.
	f            []func(rune) int // Transitions.
	startf, endf []int            // Transitions at start and end of input.
	nest         []dfa
}

var dfas = []dfa{
	// [ \t]
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 9:
				return 1
			case 32:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 9:
				return -1
			case 32:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// out
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 111:
				return 1
			case 116:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 116:
				return -1
			case 117:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 116:
				return 3
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// ,
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 44:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 44:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// ;
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 59:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 59:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// do
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 100:
				return 1
			case 111:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 111:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 111:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// channel
	{[]bool{false, false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return 1
			case 101:
				return -1
			case 104:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 104:
				return 2
			case 108:
				return -1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return 3
			case 99:
				return -1
			case 101:
				return -1
			case 104:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 104:
				return -1
			case 108:
				return -1
			case 110:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 104:
				return -1
			case 108:
				return -1
			case 110:
				return 5
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return -1
			case 101:
				return 6
			case 104:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 104:
				return -1
			case 108:
				return 7
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 99:
				return -1
			case 101:
				return -1
			case 104:
				return -1
			case 108:
				return -1
			case 110:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1, -1}, nil},

	// when
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return -1
			case 110:
				return -1
			case 119:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return 2
			case 110:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return 3
			case 104:
				return -1
			case 110:
				return -1
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return -1
			case 110:
				return 4
			case 119:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 104:
				return -1
			case 110:
				return -1
			case 119:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// from
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 102:
				return 1
			case 109:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 102:
				return -1
			case 109:
				return -1
			case 111:
				return -1
			case 114:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 102:
				return -1
			case 109:
				return -1
			case 111:
				return 3
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 102:
				return -1
			case 109:
				return 4
			case 111:
				return -1
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 102:
				return -1
			case 109:
				return -1
			case 111:
				return -1
			case 114:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// to
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 116:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return 2
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 116:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// \(
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 40:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 40:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \)
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 41:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 41:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \[
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 91:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 91:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \~
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 126:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 126:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \|
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 124:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 124:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \]
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 93:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 93:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// >=
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 62:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return 2
			case 62:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			case 62:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// >
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 62:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 62:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// send
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return 2
			case 110:
				return -1
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return -1
			case 110:
				return 3
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return 4
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			case 115:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// <>
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 60:
				return 1
			case 62:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 60:
				return -1
			case 62:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 60:
				return -1
			case 62:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// =
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 61:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 61:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// #
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 35:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 35:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// <=
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 60:
				return 1
			case 61:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 60:
				return -1
			case 61:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 60:
				return -1
			case 61:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// \/
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 47:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 47:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// <
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 60:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 60:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \+
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 43:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 43:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \-
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// @
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 64:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 64:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \^
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 94:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 94:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// \*
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 42:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 42:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// false
	{[]bool{false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 102:
				return 1
			case 108:
				return -1
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return 2
			case 101:
				return -1
			case 102:
				return -1
			case 108:
				return -1
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 108:
				return 3
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 108:
				return -1
			case 115:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return 5
			case 102:
				return -1
			case 108:
				return -1
			case 115:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 102:
				return -1
			case 108:
				return -1
			case 115:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1}, nil},

	// true
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 114:
				return -1
			case 116:
				return 1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 114:
				return 2
			case 116:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 114:
				return -1
			case 116:
				return -1
			case 117:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return 4
			case 114:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 114:
				return -1
			case 116:
				return -1
			case 117:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// on
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 110:
				return -1
			case 111:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 110:
				return 2
			case 111:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 110:
				return -1
			case 111:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// end
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return 1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return -1
			case 110:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return 3
			case 101:
				return -1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 100:
				return -1
			case 101:
				return -1
			case 110:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// and
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 97:
				return 1
			case 100:
				return -1
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 110:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return 3
			case 110:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 100:
				return -1
			case 110:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// or
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 111:
				return 1
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 114:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 111:
				return -1
			case 114:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// <-
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 60:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return 2
			case 60:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 60:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// not
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 110:
				return 1
			case 111:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 110:
				return -1
			case 111:
				return 2
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 110:
				return -1
			case 111:
				return -1
			case 116:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 110:
				return -1
			case 111:
				return -1
			case 116:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// ite
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 105:
				return 1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 105:
				return -1
			case 116:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return 3
			case 105:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 101:
				return -1
			case 105:
				return -1
			case 116:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// bool
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 98:
				return 1
			case 108:
				return -1
			case 111:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 108:
				return -1
			case 111:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 108:
				return -1
			case 111:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 108:
				return 4
			case 111:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 98:
				return -1
			case 108:
				return -1
			case 111:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// int
	{[]bool{false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 105:
				return 1
			case 110:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 105:
				return -1
			case 110:
				return 2
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 105:
				return -1
			case 110:
				return -1
			case 116:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 105:
				return -1
			case 110:
				return -1
			case 116:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// string
	{[]bool{false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 105:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return 1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 105:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 105:
				return -1
			case 110:
				return -1
			case 114:
				return 3
			case 115:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 105:
				return 4
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 105:
				return -1
			case 110:
				return 5
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return 6
			case 105:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 103:
				return -1
			case 105:
				return -1
			case 110:
				return -1
			case 114:
				return -1
			case 115:
				return -1
			case 116:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, nil},

	// real
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 108:
				return -1
			case 114:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return 2
			case 108:
				return -1
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return 3
			case 101:
				return -1
			case 108:
				return -1
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 108:
				return 4
			case 114:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 97:
				return -1
			case 101:
				return -1
			case 108:
				return -1
			case 114:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// [-]?([0-9]*[.])[0-9]+
	{[]bool{false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return 1
			case 46:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 4
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			case 46:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 4
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1}, nil},

	// -[0-9]+
	{[]bool{false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 45:
				return 1
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 45:
				return -1
			}
			switch {
			case 48 <= r && r <= 57:
				return 2
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// [0-9]+
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch {
			case 48 <= r && r <= 57:
				return 1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// [a-zA-Z_][a-zA-Z0-9_]*
	{[]bool{false, true, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 95:
				return 1
			}
			switch {
			case 48 <= r && r <= 57:
				return -1
			case 65 <= r && r <= 90:
				return 1
			case 97 <= r && r <= 122:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 95:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 2
			case 65 <= r && r <= 90:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 95:
				return 2
			}
			switch {
			case 48 <= r && r <= 57:
				return 2
			case 65 <= r && r <= 90:
				return 2
			case 97 <= r && r <= 122:
				return 2
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1}, nil},

	// \"[a-zA-Z0-9_\*\+\$\(\)\.]*\"
	{[]bool{false, false, true, false}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 34:
				return 1
			case 36:
				return -1
			case 40:
				return -1
			case 41:
				return -1
			case 42:
				return -1
			case 43:
				return -1
			case 46:
				return -1
			case 95:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return -1
			case 97 <= r && r <= 122:
				return -1
			case 48 <= r && r <= 57:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 34:
				return 2
			case 36:
				return 3
			case 40:
				return 3
			case 41:
				return 3
			case 42:
				return 3
			case 43:
				return 3
			case 46:
				return 3
			case 95:
				return 3
			}
			switch {
			case 65 <= r && r <= 90:
				return 3
			case 97 <= r && r <= 122:
				return 3
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 34:
				return -1
			case 36:
				return -1
			case 40:
				return -1
			case 41:
				return -1
			case 42:
				return -1
			case 43:
				return -1
			case 46:
				return -1
			case 95:
				return -1
			}
			switch {
			case 65 <= r && r <= 90:
				return -1
			case 97 <= r && r <= 122:
				return -1
			case 48 <= r && r <= 57:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 34:
				return 2
			case 36:
				return 3
			case 40:
				return 3
			case 41:
				return 3
			case 42:
				return 3
			case 43:
				return 3
			case 46:
				return 3
			case 95:
				return 3
			}
			switch {
			case 65 <= r && r <= 90:
				return 3
			case 97 <= r && r <= 122:
				return 3
			case 48 <= r && r <= 57:
				return 3
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1}, nil},

	// (\/\*([\r\n]|([^*\]|[\r\n]))*\*\/)
	{[]bool{false, false, false, false, false, false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 42:
				return -1
			case 47:
				return 1
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 42:
				return 2
			case 47:
				return -1
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return 3
			case 13:
				return 3
			case 42:
				return 4
			case 47:
				return 5
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return 5
		},
		func(r rune) int {
			switch r {
			case 10:
				return 3
			case 13:
				return 3
			case 42:
				return 4
			case 47:
				return 5
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return 5
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 42:
				return -1
			case 47:
				return 6
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return 3
			case 13:
				return 3
			case 42:
				return 4
			case 47:
				return 5
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return 5
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			case 13:
				return -1
			case 42:
				return -1
			case 47:
				return -1
			case 91:
				return -1
			case 93:
				return -1
			case 124:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, []int{ /* End-of-input transitions */ -1, -1, -1, -1, -1, -1, -1}, nil},

	// \n
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			switch r {
			case 10:
				return 1
			}
			return -1
		},
		func(r rune) int {
			switch r {
			case 10:
				return -1
			}
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},

	// .
	{[]bool{false, true}, []func(rune) int{ // Transitions
		func(r rune) int {
			return 1
		},
		func(r rune) int {
			return -1
		},
	}, []int{ /* Start-of-input transitions */ -1, -1}, []int{ /* End-of-input transitions */ -1, -1}, nil},
}

func NewLexer(in io.Reader) *Lexer {
	return NewLexerWithInit(in, nil)
}

func (yyLex *Lexer) Stop() {
	yyLex.ch_stop <- true
}

// Text returns the matched text.
func (yylex *Lexer) Text() string {
	return yylex.stack[len(yylex.stack)-1].s
}

// Line returns the current line number.
// The first line is 0.
func (yylex *Lexer) Line() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].line
}

// Column returns the current column number.
// The first column is 0.
func (yylex *Lexer) Column() int {
	if len(yylex.stack) == 0 {
		return 0
	}
	return yylex.stack[len(yylex.stack)-1].column
}

func (yylex *Lexer) next(lvl int) int {
	if lvl == len(yylex.stack) {
		l, c := 0, 0
		if lvl > 0 {
			l, c = yylex.stack[lvl-1].line, yylex.stack[lvl-1].column
		}
		yylex.stack = append(yylex.stack, frame{0, "", l, c})
	}
	if lvl == len(yylex.stack)-1 {
		p := &yylex.stack[lvl]
		*p = <-yylex.ch
		yylex.stale = false
	} else {
		yylex.stale = true
	}
	return yylex.stack[lvl].i
}
func (yylex *Lexer) pop() {
	yylex.stack = yylex.stack[:len(yylex.stack)-1]
}
func (yylex Lexer) Error(e string) {
	panic(e)
}

// Lex runs the lexer. Always returns 0.
// When the -s option is given, this function is not generated;
// instead, the NN_FUN macro runs the lexer.
func (yylex *Lexer) Lex(lval *yySymType) int {
OUTER0:
	for {
		switch yylex.next(0) {
		case 0:
			{ /* Skip blanks and tabs. */
			}
		case 1:
			{
				return OUT
			}
		case 2:
			{
				return COMMA
			}
		case 3:
			{
				return SC
			}
		case 4:
			{
				return DO
			}
		case 5:
			{
				return CHANNEL
			}
		case 6:
			{
				return WHEN
			}
		case 7:
			{
				return FROM
			}
		case 8:
			{
				return TO
			}
		case 9:
			{
				return LPAR
			}
		case 10:
			{
				return RPAR
			}
		case 11:
			{
				return LBR
			}
		case 12:
			{
				return TAB
			}
		case 13:
			{
				return PIPE
			}
		case 14:
			{
				return RBR
			}
		case 15:
			{
				return GE
			}
		case 16:
			{
				return G
			}
		case 17:
			{
				return SEND
			}
		case 18:
			{
				return NEQ
			}
		case 19:
			{
				return EQ
			}
		case 20:
			{
				return DIESI
			}
		case 21:
			{
				return LE
			}
		case 22:
			{
				return DIV
			}
		case 23:
			{
				return L
			}
		case 24:
			{
				return ADD
			}
		case 25:
			{
				return MINUS
			}
		case 26:
			{
				return AT
			}
		case 27:
			{
				return EXP
			}
		case 28:
			{
				return TIMES
			}
		case 29:
			{
				return FALSE
			}
		case 30:
			{
				return TRUE
			}
		case 31:
			{
				return ON
			}
		case 32:
			{
				return END
			}
		case 33:
			{
				return AND
			}
		case 34:
			{
				return OR
			}
		case 35:
			{
				return ASSIGN
			}
		case 36:
			{
				return NOT
			}
		case 37:
			{
				return ITE
			}
		case 38:
			{
				return BOOL
			}
		case 39:
			{
				return INT
			}
		case 40:
			{
				return STR
			}
		case 41:
			{
				return REAL
			}
		case 42:
			{
				lval.realn, _ = strconv.ParseFloat(yylex.Text(), 64)
				return REALNUM
			}
		case 43:
			{
				lval.n, _ = strconv.Atoi(yylex.Text())
				return NUM
			}
		case 44:
			{
				lval.n, _ = strconv.Atoi(yylex.Text())
				return NUM
			}
		case 45:
			{
				lval.name = yylex.Text()
				return ID
			}
		case 46:
			{
				lval.name = yylex.Text()
				return LITERAL
			}
		case 47:
			{ /* comments */
			}
		case 48:
			{
			}
		case 49:
			{
			}
		default:
			break OUTER0
		}
		continue
	}
	yylex.pop()

	return 0
}
