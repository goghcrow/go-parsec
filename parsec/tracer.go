package parsec

import "fmt"

var (
	num       = 0
	traceFlag = false
	errFmt    func(error) string
)

func EnableTrace(f func(), errF func(error) string) {
	defer func() {
		traceFlag = false
		num = 0
		errFmt = nil
	}()
	errFmt = errF
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
			if out.Success {
				fmt.Printf("[%-3d] Success(%v)\n", num, out.Candidates)
			} else {
				if errFmt == nil {
					fmt.Printf("[%-3d] Error(%v)\n", num, out.Error)
				} else {
					fmt.Printf("[%-3d] %s\n", num, errFmt(out.Error))
				}
			}
		}
		return out
	})
}
