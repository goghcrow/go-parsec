package legacy

// type Tuple1[T1 any] struct {
// 	V1 T1
// }
//
// func Seq1[K Ord, T1 any](p1 Parser[K, T1]) Parser[K, Tuple1[T1]] {
// 	return parser[K, Tuple1[T1]](func(toks []Token[K]) Output[K, Tuple1[T1]] {
// 		var xs []Result[K, Tuple1[T1]]
// 		out := p1.Parse(toks)
// 		if out.Success {
// 			for _, candidate := range out.Candidates {
// 				xs = append(xs, Result[K, Tuple1[T1]]{
// 					Val:  Tuple1[T1]{candidate.Val},
// 					next: candidate.next,
// 				})
// 			}
// 		}
// 		return newOutput(xs, out.Error, len(xs) != 0)
// 	})
// }
//
//
// type Tuple2[T1, T2 any] struct {
// 	V1 T1
// 	V2 T2
// }
// func Seq2[K Ord, T1, T2 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// ) Parser[K, Tuple2[T1, T2]] {
// 	s1 := Seq1(p1)
// 	return parser[K, Tuple2[T1, T2]](func(toks []Token[K]) Output[K, Tuple2[T1, T2]] {
// 		out1 := s1.Parse(toks)
// 		if !out1.Success {
// 			return failOf[K, Tuple1[T1], Tuple2[T1, T2]](out1)
// 		}
//
// 		var err *Error
// 		steps := out1.Candidates
// 		var xs []Result[K, Tuple2[T1, T2]]
// 		for _, step := range steps {
// 			out2 := p2.Parse(step.next)
// 			err = betterError(out1.Error, out2.Error)
// 			if out2.Success {
// 				for _, candidate := range out2.Candidates {
// 					xs = append(xs, Result[K, Tuple2[T1, T2]]{
// 						Val: Tuple2[T1, T2]{V1: step.Val.V1, V2: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple3[T1, T2, T3 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// }
// func Seq3[K Ord, T1, T2, T3 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// ) Parser[K, Tuple3[T1, T2, T3]] {
// 	s2 := Seq2(p1, p2)
// 	return parser[K, Tuple3[T1, T2, T3]](func(toks []Token[K]) Output[K, Tuple3[T1, T2, T3]] {
// 		out2 := s2.Parse(toks)
// 		if !out2.Success {
// 			return failOf[K, Tuple2[T1, T2], Tuple3[T1, T2, T3]](out2)
// 		}
//
// 		var err *Error
// 		steps := out2.Candidates
// 		var xs []Result[K, Tuple3[T1, T2, T3]]
// 		for _, step := range steps {
// 			out3 := p3.Parse(step.next)
// 			err = betterError(out2.Error, out3.Error)
// 			if out3.Success {
// 				for _, candidate := range out3.Candidates {
// 					xs = append(xs, Result[K, Tuple3[T1, T2, T3]]{
// 						Val: Tuple3[T1, T2, T3]{V1: step.Val.V1, V2: step.Val.V2, V3: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple4[T1, T2, T3, T4 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// }
// func Seq4[K Ord, T1, T2, T3, T4 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// ) Parser[K, Tuple4[T1, T2, T3, T4]] {
// 	s3 := Seq3(p1, p2, p3)
// 	return parser[K, Tuple4[T1, T2, T3, T4]](func(toks []Token[K]) Output[K, Tuple4[T1, T2, T3, T4]] {
// 		out3 := s3.Parse(toks)
// 		if !out3.Success {
// 			return failOf[K, Tuple3[T1, T2, T3], Tuple4[T1, T2, T3, T4]](out3)
// 		}
//
// 		var err *Error
// 		steps := out3.Candidates
// 		var xs []Result[K, Tuple4[T1, T2, T3, T4]]
// 		for _, step := range steps {
// 			out4 := p4.Parse(step.next)
// 			err = betterError(out3.Error, out4.Error)
// 			if out4.Success {
// 				for _, candidate := range out4.Candidates {
// 					xs = append(xs, Result[K, Tuple4[T1, T2, T3, T4]]{
// 						Val: Tuple4[T1, T2, T3, T4]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple5[T1, T2, T3, T4, T5 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// 	V5 T5
// }
// func Seq5[K Ord, T1, T2, T3, T4, T5 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// 	p5 Parser[K, T5],
// ) Parser[K, Tuple5[T1, T2, T3, T4, T5]] {
// 	s4 := Seq4(p1, p2, p3, p4)
// 	return parser[K, Tuple5[T1, T2, T3, T4, T5]](func(toks []Token[K]) Output[K, Tuple5[T1, T2, T3, T4, T5]] {
// 		out4 := s4.Parse(toks)
// 		if !out4.Success {
// 			return failOf[K, Tuple4[T1, T2, T3, T4], Tuple5[T1, T2, T3, T4, T5]](out4)
// 		}
//
// 		var err *Error
// 		steps := out4.Candidates
// 		var xs []Result[K, Tuple5[T1, T2, T3, T4, T5]]
// 		for _, step := range steps {
// 			out5 := p5.Parse(step.next)
// 			err = betterError(out4.Error, out5.Error)
// 			if out5.Success {
// 				for _, candidate := range out5.Candidates {
// 					xs = append(xs, Result[K, Tuple5[T1, T2, T3, T4, T5]]{
// 						Val: Tuple5[T1, T2, T3, T4, T5]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: step.Val.V4, V5: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple6[T1, T2, T3, T4, T5, T6 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// 	V5 T5
// 	V6 T6
// }
// func Seq6[K Ord, T1, T2, T3, T4, T5, T6 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// 	p5 Parser[K, T5],
// 	p6 Parser[K, T6],
// ) Parser[K, Tuple6[T1, T2, T3, T4, T5, T6]] {
// 	s5 := Seq5(p1, p2, p3, p4, p5)
// 	return parser[K, Tuple6[T1, T2, T3, T4, T5, T6]](func(toks []Token[K]) Output[K, Tuple6[T1, T2, T3, T4, T5, T6]] {
// 		out5 := s5.Parse(toks)
// 		if !out5.Success {
// 			return failOf[K, Tuple5[T1, T2, T3, T4, T5], Tuple6[T1, T2, T3, T4, T5, T6]](out5)
// 		}
//
// 		var err *Error
// 		steps := out5.Candidates
// 		var xs []Result[K, Tuple6[T1, T2, T3, T4, T5, T6]]
// 		for _, step := range steps {
// 			out6 := p6.Parse(step.next)
// 			err = betterError(out5.Error, out6.Error)
// 			if out6.Success {
// 				for _, candidate := range out6.Candidates {
// 					xs = append(xs, Result[K, Tuple6[T1, T2, T3, T4, T5, T6]]{
// 						Val: Tuple6[T1, T2, T3, T4, T5, T6]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: step.Val.V4, V5: step.Val.V5, V6: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple7[T1, T2, T3, T4, T5, T6, T7 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// 	V5 T5
// 	V6 T6
// 	V7 T7
// }
// func Seq7[K Ord, T1, T2, T3, T4, T5, T6, T7 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// 	p5 Parser[K, T5],
// 	p6 Parser[K, T6],
// 	p7 Parser[K, T7],
// ) Parser[K, Tuple7[T1, T2, T3, T4, T5, T6, T7]] {
// 	s6 := Seq6(p1, p2, p3, p4, p5, p6)
// 	return parser[K, Tuple7[T1, T2, T3, T4, T5, T6, T7]](func(toks []Token[K]) Output[K, Tuple7[T1, T2, T3, T4, T5, T6, T7]] {
// 		out6 := s6.Parse(toks)
// 		if !out6.Success {
// 			return failOf[K, Tuple6[T1, T2, T3, T4, T5, T6], Tuple7[T1, T2, T3, T4, T5, T6, T7]](out6)
// 		}
//
// 		var err *Error
// 		steps := out6.Candidates
// 		var xs []Result[K, Tuple7[T1, T2, T3, T4, T5, T6, T7]]
// 		for _, step := range steps {
// 			out7 := p7.Parse(step.next)
// 			err = betterError(out6.Error, out7.Error)
// 			if out7.Success {
// 				for _, candidate := range out7.Candidates {
// 					xs = append(xs, Result[K, Tuple7[T1, T2, T3, T4, T5, T6, T7]]{
// 						Val: Tuple7[T1, T2, T3, T4, T5, T6, T7]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: step.Val.V4, V5: step.Val.V5, V6: step.Val.V6, V7: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// 	V5 T5
// 	V6 T6
// 	V7 T7
// 	V8 T8
// }
// func Seq8[K Ord, T1, T2, T3, T4, T5, T6, T7, T8 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// 	p5 Parser[K, T5],
// 	p6 Parser[K, T6],
// 	p7 Parser[K, T7],
// 	p8 Parser[K, T8],
// ) Parser[K, Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]] {
// 	s7 := Seq7(p1, p2, p3, p4, p5, p6, p7)
// 	return parser[K, Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]](func(toks []Token[K]) Output[K, Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]] {
// 		out7 := s7.Parse(toks)
// 		if !out7.Success {
// 			return failOf[K, Tuple7[T1, T2, T3, T4, T5, T6, T7], Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]](out7)
// 		}
//
// 		var err *Error
// 		steps := out7.Candidates
// 		var xs []Result[K, Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
// 		for _, step := range steps {
// 			out8 := p8.Parse(step.next)
// 			err = betterError(out7.Error, out8.Error)
// 			if out8.Success {
// 				for _, candidate := range out8.Candidates {
// 					xs = append(xs, Result[K, Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]{
// 						Val: Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: step.Val.V4, V5: step.Val.V5, V6: step.Val.V6, V7: step.Val.V7, V8: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// 	V5 T5
// 	V6 T6
// 	V7 T7
// 	V8 T8
// 	V9 T9
// }
// func Seq9[K Ord, T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// 	p5 Parser[K, T5],
// 	p6 Parser[K, T6],
// 	p7 Parser[K, T7],
// 	p8 Parser[K, T8],
// 	p9 Parser[K, T9],
// ) Parser[K, Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]] {
// 	s8 := Seq8(p1, p2, p3, p4, p5, p6, p7, p8)
// 	return parser[K, Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]](func(toks []Token[K]) Output[K, Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]] {
// 		out8 := s8.Parse(toks)
// 		if !out8.Success {
// 			return failOf[K, Tuple8[T1, T2, T3, T4, T5, T6, T7, T8], Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]](out8)
// 		}
//
// 		var err *Error
// 		steps := out8.Candidates
// 		var xs []Result[K, Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
// 		for _, step := range steps {
// 			out9 := p9.Parse(step.next)
// 			err = betterError(out8.Error, out9.Error)
// 			if out9.Success {
// 				for _, candidate := range out9.Candidates {
// 					xs = append(xs, Result[K, Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]{
// 						Val: Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: step.Val.V4, V5: step.Val.V5, V6: step.Val.V6, V7: step.Val.V7, V8: step.Val.V8, V9: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
//
// type Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any] struct {
// 	V1 T1
// 	V2 T2
// 	V3 T3
// 	V4 T4
// 	V5 T5
// 	V6 T6
// 	V7 T7
// 	V8 T8
// 	V9 T9
// 	V10 T10
// }
// func Seq10[K Ord, T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 any](
// 	p1 Parser[K, T1],
// 	p2 Parser[K, T2],
// 	p3 Parser[K, T3],
// 	p4 Parser[K, T4],
// 	p5 Parser[K, T5],
// 	p6 Parser[K, T6],
// 	p7 Parser[K, T7],
// 	p8 Parser[K, T8],
// 	p9 Parser[K, T9],
// 	p10 Parser[K, T10],
// ) Parser[K, Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]] {
// 	s9 := Seq9(p1, p2, p3, p4, p5, p6, p7, p8, p9)
// 	return parser[K, Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]](func(toks []Token[K]) Output[K, Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]] {
// 		out9 := s9.Parse(toks)
// 		if !out9.Success {
// 			return failOf[K, Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9], Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]](out9)
// 		}
//
// 		var err *Error
// 		steps := out9.Candidates
// 		var xs []Result[K, Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]
// 		for _, step := range steps {
// 			out10 := p10.Parse(step.next)
// 			err = betterError(out9.Error, out10.Error)
// 			if out10.Success {
// 				for _, candidate := range out10.Candidates {
// 					xs = append(xs, Result[K, Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]]{
// 						Val: Tuple10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10]{V1: step.Val.V1, V2: step.Val.V2, V3: step.Val.V3, V4: step.Val.V4, V5: step.Val.V5, V6: step.Val.V6, V7: step.Val.V7, V8: step.Val.V8, V9: step.Val.V9, V10: candidate.Val},
// 						next: candidate.next,
// 					})
// 				}
// 			}
// 		}
// 		return newOutput(xs, err, len(xs) != 0)
// 	})
// }
