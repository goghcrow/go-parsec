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
	UnknownPos = VirtualPos("<unknown>")
	EOFPos     = VirtualPos("<EOF>")
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

// ----------------------------------------------------------------
// tokSeq
// ----------------------------------------------------------------

type tokSeq[K Ord] []Token[K]

func (t tokSeq[K]) beginTok() Token[K] {
	if len(t) == 0 {
		return nil
	} else {
		return t[0]
	}
}

func (t tokSeq[K]) beginPos() Pos {
	if len(t) == 0 {
		return UnknownPos
	} else {
		return t[0]
	}
}

func (t tokSeq[K]) equals(other tokSeq[K]) bool {
	if len(t) == 0 && len(other) == 0 {
		return true
	}
	if len(t) == 0 || len(other) == 0 || t[0] != other[0] {
		return false
	}
	return true
}
