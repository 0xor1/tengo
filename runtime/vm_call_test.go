package runtime_test

import "testing"

func TestCall(t *testing.T) {
	expect(t, `a := { b: func(x) { return x + 2 } }; out = a.b(5)`, 7)
	expect(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a.b.c(5)`, 7)
	expect(t, `a := { b: { c: func(x) { return x + 2 } } }; out = a["b"].c(5)`, 7)
}
