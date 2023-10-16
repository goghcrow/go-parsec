package example

import (
	"fmt"
	"testing"
)

func TestCons(t *testing.T) {
	xs := CONS(1, CONS(2, CONS(3, CONS(4, 5))))
	t.Log(xs)
	t.Log(c5car(xs))
	t.Log(c5cadr(xs))
	t.Log(c5caddr(xs))
	t.Log(c5cadddr(xs))
	t.Log(c5caddddr(xs))
}

//goland:noinspection GoExportedFuncWithUnexportedType
func CONS[T1, T2 any](car T1, cdr T2) cons[T1, T2] {
	return cons[T1, T2]{car, cdr}
}

type cons[T1, T2 any] struct {
	car T1
	cdr T2
}

func (c cons[T1, T2]) String() string { return fmt.Sprintf("(%v, %v)", c.car, c.cdr) }

func c2car[T1, T2 any](t2 cons[T1, T2]) T1 { return t2.car }
func c2cdr[T1, T2 any](t2 cons[T1, T2]) T2 { return t2.cdr }

func c3car[T1, T2, T3 any](t3 cons[T1, cons[T2, T3]]) T1  { return t3.car }
func c3cadr[T1, T2, T3 any](t3 cons[T1, cons[T2, T3]]) T2 { return t3.cdr.car }
func c3cddr[T1, T2, T3 any](t3 cons[T1, cons[T2, T3]]) T3 { return t3.cdr.cdr }

func c4car[T1, T2, T3, T4 any](t4 cons[T1, cons[T2, cons[T3, T4]]]) T1   { return t4.car }
func c4cadr[T1, T2, T3, T4 any](t4 cons[T1, cons[T2, cons[T3, T4]]]) T2  { return t4.cdr.car }
func c4caddr[T1, T2, T3, T4 any](t4 cons[T1, cons[T2, cons[T3, T4]]]) T3 { return t4.cdr.cdr.car }
func c4cdddr[T1, T2, T3, T4 any](t4 cons[T1, cons[T2, cons[T3, T4]]]) T4 { return t4.cdr.cdr.cdr }

func c5car[T1, T2, T3, T4, T5 any](t5 cons[T1, cons[T2, cons[T3, cons[T4, T5]]]]) T1 { return t5.car }
func c5cadr[T1, T2, T3, T4, T5 any](t5 cons[T1, cons[T2, cons[T3, cons[T4, T5]]]]) T2 {
	return t5.cdr.car
}
func c5caddr[T1, T2, T3, T4, T5 any](t5 cons[T1, cons[T2, cons[T3, cons[T4, T5]]]]) T3 {
	return t5.cdr.cdr.car
}
func c5cadddr[T1, T2, T3, T4, T5 any](t5 cons[T1, cons[T2, cons[T3, cons[T4, T5]]]]) T4 {
	return t5.cdr.cdr.cdr.car
}
func c5caddddr[T1, T2, T3, T4, T5 any](t5 cons[T1, cons[T2, cons[T3, cons[T4, T5]]]]) T5 {
	return t5.cdr.cdr.cdr.cdr
}

type P[T any] struct{}

func S2[T1, T2 any](
	p1 P[T1], p2 P[T2],
) P[cons[T1, T2]] {
	return P[cons[T1, T2]]{}
}
func S3[T1, T2, T3 any](
	p1 P[T1],
	p2 P[T2],
	p3 P[T3],
) P[cons[T1, cons[T2, T3]]] {
	return S2(p1, S2(p2, p3))
}
func S4[T1, T2, T3, T4 any](
	p1 P[T1],
	p2 P[T2],
	p3 P[T3],
	p4 P[T4],
) P[cons[T1, cons[T2, cons[T3, T4]]]] {
	var p = S2(p1, S2(p2, S2(p3, p4)))
	return p
}
func S5[T1, T2, T3, T4, T5 any](
	p1 P[T1],
	p2 P[T2],
	p3 P[T3],
	p4 P[T4],
	p5 P[T5],
) P[cons[T1, cons[T2, cons[T3, cons[T4, T5]]]]] {
	p2345 := S2(p2, S2(p3, S2(p4, p5)))
	// make golang happy
	p := S2(p1, p2345)
	return p
}
func S6[T1, T2, T3, T4, T5, T6 any](
	p1 P[T1],
	p2 P[T2],
	p3 P[T3],
	p4 P[T4],
	p5 P[T5],
	p6 P[T6],
) P[cons[T1, cons[T2, cons[T3, cons[T4, cons[T5, T6]]]]]] {
	// make golang happy
	p3456 := S2(p3, S2(p4, S2(p5, p6)))
	p23456 := S2(p2, p3456)
	p123456 := S2(p1, p23456)
	return p123456
}
