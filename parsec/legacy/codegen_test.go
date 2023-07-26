package legacy

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"
)

func TestGenSeq(*testing.T) {
	return
	funcs := template.FuncMap{
		"Range": func(from, to int) []int {
			var i int
			var xs []int
			for i = from; i <= to; i++ {
				xs = append(xs, i)
			}
			return xs
		},
		"Sub1": func(n int) int { return n - 1 },
	}

	tplSeq := template.New("Seq").Funcs(funcs)
	tplSeq = template.Must(tplSeq.Parse(`
{{define "subTx"}}{{range $i := Range 1 (Sub1 (Sub1 .N)) }}T{{$i}}, {{end}}T{{(Sub1 .N)}}{{end}}
{{define "subTupleX"}}Tuple{{(Sub1 .N)}}[{{template "subTx" .}}]{{end}}
{{define "tx"}}{{range $i := Range 1 (Sub1 .N) }}T{{$i}}, {{end}}T{{.N}}{{end}}
{{define "tupleX"}}Tuple{{.N}}[{{template "tx" .}}]{{end}}
{{define "px"}}{{range $i := Range 1 (Sub1 (Sub1 .N)) }}p{{$i}}, {{end}}p{{(Sub1 .N)}}{{end}}
{{define "subVx"}}{{range $i := Range 1 (Sub1 .N) }}V{{$i}}: step.Val.V{{$i}}, {{end}}{{end}}
type Tuple{{.N}}[{{template "tx" .}} any] struct {
	{{range $i := Range 1 (Sub1 .N) }}V{{$i}} T{{$i}}
	{{end}}V{{.N}} T{{.N}}
}
func Seq{{.N}}[K Ord, {{template "tx" .}} any](
	{{range $i := Range 1 (Sub1 .N) }}p{{$i}} Parser[K, T{{$i}}],
	{{end}}p{{.N}} Parser[K, T{{.N}}],
) Parser[K, {{template "tupleX" .}}] {
	s{{Sub1 .N}} := Seq{{Sub1 .N}}({{template "px" .}})
	return parser[K, {{template "tupleX" .}}](func(toks []Token[K]) Output[K, {{template "tupleX" .}}] {
		out{{Sub1 .N}} := s{{Sub1 .N}}.Parse(toks)
		if !out{{Sub1 .N}}.Success {
			return failOf[K, {{template "subTupleX" .}}, {{template "tupleX" .}}](out{{Sub1 .N}})
		}

		var err *Error
		steps := out{{Sub1 .N}}.Candidates
		var xs []Result[K, {{template "tupleX" .}}]
		for _, step := range steps {
			out{{.N}} := p{{.N}}.Parse(step.next)
			err = betterError(out{{Sub1 .N}}.Error, out{{.N}}.Error)
			if out{{.N}}.Success {
				for _, candidate := range out{{.N}}.Candidates {
					xs = append(xs, Result[K, {{template "tupleX" .}}]{
						Val: {{template "tupleX" .}}{{"{"}}{{template "subVx" .}}V{{.N}}: candidate.Val},
						next: candidate.next,
					})
				}
			}
		}
		return newOutput(xs, err, len(xs) != 0)
	})
}
`))

	tplTuple := template.New("tuple").Funcs(funcs)
	tplTuple = template.Must(tplTuple.Parse(`{{define "tx"}}{{range $i := Range 1 (Sub1 .N) }}T{{$i}}, {{end}}T{{.N}}{{end}}
`))

	output, err := os.OpenFile("./generated.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	buf := bytes.NewBufferString(`package parsec

type Tuple1[T1 any] struct {
	V1 T1
}

func Seq1[K Ord, T1 any](p1 Parser[K, T1]) Parser[K, Tuple1[T1]] {
	return parser[K, Tuple1[T1]](func(toks []Token[K]) Output[K, Tuple1[T1]] {
		var xs []Result[K, Tuple1[T1]]
		out := p1.Parse(toks)
		if out.Success {
			for _, candidate := range out.Candidates {
				xs = append(xs, Result[K, Tuple1[T1]]{
					Val:  Tuple1[T1]{candidate.Val},
					next: candidate.next,
				})
			}
		}
		return newOutput(xs, out.Error, len(xs) != 0)
	})
}
`)
	const n = 11
	for i := 2; i < n; i++ {
		buf1 := bytes.NewBufferString("")
		err = tplSeq.Execute(buf1, struct{ N int }{N: i})
		if err != nil {
			panic(err)
		}
		buf.WriteString("\n\n")
		buf.WriteString(strings.Trim(buf1.String(), "\n"))
	}

	_, _ = output.WriteString(strings.Trim(buf.String(), "\n"))
}
