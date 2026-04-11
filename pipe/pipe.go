// Package pipe provides F#-style pipe operators for Go.
// Since Go has no |> operator, Pipe functions thread a value left-to-right
// through a chain of single-argument functions.
//
// Example:
//
//	result := pipe.Pipe3(
//	    rawInput,
//	    strings.TrimSpace,
//	    strings.ToLower,
//	    validate,
//	)
package pipe

// Pipe2 threads v through f1 then f2.
func Pipe2[A, B, C any](v A, f1 func(A) B, f2 func(B) C) C {
	return f2(f1(v))
}

// Pipe3 threads v through f1 → f2 → f3.
func Pipe3[A, B, C, D any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D) D {
	return f3(f2(f1(v)))
}

// Pipe4 threads v through f1 → f2 → f3 → f4.
func Pipe4[A, B, C, D, E any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E) E {
	return f4(f3(f2(f1(v))))
}

// Pipe5 threads v through f1 → … → f5.
func Pipe5[A, B, C, D, E, F any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F) F {
	return f5(f4(f3(f2(f1(v)))))
}

// Pipe6 threads v through f1 → … → f6.
func Pipe6[A, B, C, D, E, F, G any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F, f6 func(F) G) G {
	return f6(f5(f4(f3(f2(f1(v))))))
}

// Pipe7 threads v through f1 → … → f7.
func Pipe7[A, B, C, D, E, F, G, H any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F, f6 func(F) G, f7 func(G) H) H {
	return f7(f6(f5(f4(f3(f2(f1(v)))))))
}

// Pipe8 threads v through f1 → … → f8.
func Pipe8[A, B, C, D, E, F, G, H, I any](v A, f1 func(A) B, f2 func(B) C, f3 func(C) D, f4 func(D) E, f5 func(E) F, f6 func(F) G, f7 func(G) H, f8 func(H) I) I {
	return f8(f7(f6(f5(f4(f3(f2(f1(v))))))))
}

