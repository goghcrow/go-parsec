package parsec

// ----------------------------------------------------------------
// Location
// ----------------------------------------------------------------

type Pos interface {
	Loc() (idx /*include*/, idxEnd /*exclude*/, col, ln int)
}

type VirtualPos string

func (VirtualPos) Loc() (int, int, int, int) { return -1, -1, -1, -1 }

var (
	UnknownPos = VirtualPos("unknown")
	EOFPos     = VirtualPos("end of input")
)

// ----------------------------------------------------------------
// tokSeq
// ----------------------------------------------------------------

type Token[K Ord] interface {
	Pos
	Kind() K
	Lexeme() string
	String() string
}

type virtualToken[K Ord] struct {
	VirtualPos
	name string
}

func (virtualToken[K]) Kind() K          { return *new(K) }
func (v virtualToken[K]) Lexeme() string { return v.name }
func (v virtualToken[K]) String() string { return v.name }

func EOFToken[K Ord]() Token[K] {
	return VirtualToken[K]("<EOF>", EOFPos)
}

func VirtualToken[K Ord](name string, pos VirtualPos) Token[K] {
	return virtualToken[K]{pos, name}
}
