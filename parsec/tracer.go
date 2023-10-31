package parsec

import "fmt"

var (
	num       = 0
	traceFlag = false
)

func EnableTrace(f func()) {
	defer func() {
		traceFlag = false
		num = 0
	}()
	traceFlag = true
	num = 0
	f()
}

func Trace[K TK, R any](name string, p Parser[K, R]) Parser[K, R] {
	return parser[K, R](func(toks []Token[K]) Output[K, R] {
		if traceFlag {
			// fmt.Println(toks)
			fmt.Printf("[%-3d] %s\n", num, name)
		}
		num++
		out := p.Parse(toks)
		num--
		if traceFlag {
			fmt.Printf("[%-3d] %s\n", num, out)
		}
		return out
	})
}
