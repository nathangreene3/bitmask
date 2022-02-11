package umask

import (
	"math"
	"testing"
)

// -------------------------------------------------------------------------
// TODO: Finish testing
// -------------------------------------------------------------------------

func TestAnd(t *testing.T) {
	type testStruct struct {
		a, b, exp UMask
	}

	tcs := []testStruct{
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Max.ClrBits(0, BitCap-1),
			exp: Zero,
		},
		{
			a:   Max.ClrBits(0, BitCap-1),
			b:   Max,
			exp: Max.ClrBits(0, BitCap-1),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.And(tc.b); tc.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestBitLen(t *testing.T) {
	var exp int // Iterates over range [0,BitCap]
	if rec := Zero.BitLen(); exp != rec {
		t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
	}

	for exp++; exp < BitCap; exp++ {
		if rec := Zero.SetBits(exp - 1).BitLen(); exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
		}
	}
}

func TestBits(t *testing.T) {
	type testCase struct {
		bits []int
	}

	tcs := []testCase{
		{bits: []int{}},
		{bits: []int{0}},
		{bits: []int{0, BitCap - 1}},
	}

	for _, tc := range tcs {
		var (
			rec   = Zero.SetBits(tc.bits...).Bits()
			equal = len(tc.bits) == len(rec)
		)

		for i := 0; i < len(tc.bits) && equal; i++ {
			equal = tc.bits[i] == rec[i]
		}

		if !equal {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.bits, rec)
		}
	}
}

func TestClr(t *testing.T) {
	type testCase struct {
		a, b, exp UMask
	}

	tcs := []testCase{
		{
			a:   Zero,
			b:   Zero,
			exp: Zero,
		},
		{
			a:   Zero,
			b:   Zero.SetBits(0, BitCap-1),
			exp: Zero,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Zero,
			exp: Zero.SetBits(0, BitCap-1),
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Zero.SetBits(0, BitCap-1),
			exp: Zero,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Max,
			exp: Zero,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Clr(tc.b); tc.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, rec)
		}
	}
}

func TestClrBit(t *testing.T) {
	type testCase struct {
		a, exp UMask
		bit    int
	}

	tcs := []testCase{
		{
			a:   Zero.SetBits(),
			exp: Zero.SetBits(),
			bit: 0,
		},
		{
			a:   Zero.SetBits(0),
			exp: Zero.SetBits(),
			bit: 0,
		},
		{
			a:   Zero.SetBits(1),
			exp: Zero.SetBits(),
			bit: 1,
		},
		{
			a:   Zero.SetBits(0, 1),
			exp: Zero.SetBits(0),
			bit: 1,
		},
		{
			a:   Zero.SetBits(0, 1),
			exp: Zero.SetBits(1),
			bit: 0,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			exp: Zero.SetBits(0),
			bit: BitCap - 1,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			exp: Zero.SetBits(BitCap - 1),
			bit: 0,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.ClrBit(tc.bit); tc.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, rec)
		}
	}
}

func TestClrBits(t *testing.T) {
	type testCase struct {
		a, b, exp UMask
	}

	tcs := []testCase{
		{
			a:   Zero,
			b:   Zero,
			exp: Zero,
		},
		{
			a:   Zero,
			b:   Zero.SetBits(0, BitCap-1),
			exp: Zero,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Zero,
			exp: Zero.SetBits(0, BitCap-1),
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Zero.SetBits(0, BitCap-1),
			exp: Zero,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Max,
			exp: Zero,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.ClrBits(tc.b.Bits()...); tc.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, rec)
		}
	}
}

func TestLRSh(t *testing.T) {
	type testCase struct {
		a, expLeft, expRight UMask
		left, right          int
	}

	tcs := []testCase{
		// Shifts less than BitCap
		{
			a:        Zero.SetBits(0, BitCap-1),
			left:     0,
			right:    0,
			expLeft:  Zero.SetBits(0, BitCap-1),
			expRight: Zero.SetBits(0, BitCap-1),
		},
		{
			a:        Zero.SetBits(0, BitCap-1),
			left:     1,
			right:    1,
			expLeft:  Zero.SetBits(1),
			expRight: Zero.SetBits(BitCap - 2),
		},
		{
			a:        Zero.SetBits(0, BitCap-1),
			left:     BitCap - 1,
			right:    BitCap - 1,
			expLeft:  Zero.SetBits(BitCap - 1),
			expRight: Zero.SetBits(0),
		},

		// Shifts larger than or equal to BitCap
		{
			a:        Zero.SetBits(0, BitCap-1),
			left:     BitCap,
			right:    BitCap,
			expLeft:  Zero,
			expRight: Zero,
		},
		{
			a:        Zero.SetBits(0, BitCap-1),
			left:     BitCap + 1,
			right:    BitCap + 1,
			expLeft:  Zero,
			expRight: Zero,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.LSh(tc.left); tc.expLeft != rec {
			t.Errorf("\nexpected %s\nreceived %s\n", tc.expLeft.Fmt(2), rec.Fmt(2))
		}

		if rec := tc.a.RSh(tc.right); tc.expRight != rec {
			t.Errorf("\nexpected %s\nreceived %s\n", tc.expRight.Fmt(2), rec.Fmt(2))
		}
	}
}

func TestNextPrevBit(t *testing.T) {
	type testCase struct {
		a UMask
	}

	tcs := []testCase{
		{a: Zero},
		{a: One},
		{a: Zero.SetBits(0, BitCap-1)}, // End bits
		{a: Max.ClrBits(0, BitCap-1)},  // Middle bits
		{a: Max},
		{a: Zero.SetBits(0, 2, 4, 6, BitCap-8, BitCap-6, BitCap-4, BitCap-2)}, // Some even bits
		{a: Zero.SetBits(1, 3, 5, 7, BitCap-7, BitCap-5, BitCap-3, BitCap-1)}, // Some odd bits
	}

	for _, tc := range tcs {
		for i := 0; i < BitCap; i++ {
			if tc.a.MasksBit(i) {
				if rec := tc.a.NextBit(i - 1); i != rec {
					t.Errorf("\nexpected next bit to be %d\nreceived next bit %d\n", i, rec)
				}

				if rec := tc.a.PrevBit(i + 1); i != rec {
					t.Errorf("\nexpected previous bit to be %d\nreceived next bit %d\n", i, rec)
				}
			} else {
				if rec := tc.a.NextBit(i - 1); i == rec {
					t.Errorf("\nexpected next bit to NOT be %d\nreceived next bit %d\n", i, rec)
				}

				if rec := tc.a.PrevBit(i + 1); i == rec {
					t.Errorf("\nexpected previous bit to NOT be %d\nreceived next bit %d\n", i, rec)
				}
			}
		}

		{
			exp, rec := tc.a, Zero
			for i := tc.a.NextBit(-1); i < BitCap; i = tc.a.NextBit(i) {
				rec = rec.SetBit(i)
			}

			if exp != rec {
				t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
			}
		}

		{
			exp, rec := Zero, tc.a
			for i := tc.a.PrevBit(BitCap); -1 < i; i = tc.a.PrevBit(i) {
				rec = rec.ClrBit(i)
			}

			if exp != rec {
				t.Errorf("\nexpected %v\nreceived %v\n", exp, rec)
			}
		}
	}
}

func TestNot(t *testing.T) {
	type testCase struct {
		a, exp UMask
	}

	tcs := []testCase{
		{
			a:   Zero.SetBits(0, BitCap-1),
			exp: Max.ClrBits(0, BitCap-1),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Not(); tc.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestOr(t *testing.T) {
	type testCase struct {
		a, b, exp UMask
	}

	tcs := []testCase{
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Max.ClrBits(0, BitCap-1),
			exp: Max,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Or(tc.b); tc.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestXOr(t *testing.T) {
	type testCase struct {
		a, b, exp UMask
	}

	tcs := []testCase{
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Max.ClrBits(0, BitCap-1),
			exp: Max,
		},
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Zero.SetBits(0, BitCap-1),
			exp: Zero,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.XOr(tc.b); tc.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

// -------------------------------------------------------------------------
// Applications
// -------------------------------------------------------------------------

func TestDivisors(t *testing.T) {
	// gcd returns the largest divisor of two given numbers.
	var gcd = func(a, b int) int {
		switch {
		case a == 0:
			if b == 0 {
				panic("gcd(0,0) is undefined")
			}

			return b
		case b == 0:
			return a
		}

		if a < 0 {
			a = -a
		}

		if b < 0 {
			b = -b
		}

		for 0 < b {
			a, b = b, a%b
		}

		return a
	}

	var n0, n1, dn int = 1, BitCap, 1
	for n := n0; n <= n1; n += dn {
		// -----------------------------------------------------------------
		// For any a on range [1,n], we search for b such that a*b = n.
		// While a <= b, we iterate computing d = a*b until d = n. When
		// d = n, we set bits a and b in the bitmask as divisors of n. We
		// then increment a, decrement b, and increment d by b-a+1. If
		// d < n, then we increment a and increment d by b instead of
		// recomputing d as a*b. If n < d, then we decrement b and decrement
		// d by a for the same reason. If b < a, then we have run through
		// all possible divisors and hault.
		// -----------------------------------------------------------------

		var (
			divisors UMask = Zero
			a, b     int   = 1, n
			d        int   = n // a*b
		)

		for a <= b {
			switch {
			case d < n:
				a++
				d += b
			case n < d:
				b--
				d -= a
			default:
				divisors = divisors.SetBits(a, b)
				a++
				b--
				d += b - a + 1
			}
		}

		for m := 0; m <= n; m++ {
			if gcd(m, n) == m {
				if !divisors.MasksBit(m) {
					t.Errorf("\nexpected %d to be set\n", m)
				}
			} else if divisors.MasksBit(m) {
				t.Errorf("\nexpected %d to not be set\n", m)
			}
		}
	}
}

func TestFibonacciNumbers(t *testing.T) {
	// ---------------------------------------------------------------------
	// Given a(0) = a(1) = 1, the nth Fibonacci number is defined as
	// a(n) := a(n-2) + a(n-1). Initializing the bitmask with bits 1 and 2
	// set, we can then get the largest set bit a(n-1) add it to the
	// next largest set bit a(n-2) and set the result as a(n).
	// ---------------------------------------------------------------------

	var fibs UMask = Zero.SetBits(1, 2)
	for n := fibs.Count(); n < BitCap; n++ {
		var b int = fibs.PrevBit(BitCap)
		fibs = fibs.SetBit(fibs.PrevBit(b) + b)
	}

	for a0, a1 := 0, 1; a1 < BitCap; a0, a1 = a1, a0+a1 {
		if !fibs.MasksBit(a1) {
			t.Errorf("\nexpected %d to be masked\n", a1)
		}
	}
}

func TestPrimes(t *testing.T) {
	// ---------------------------------------------------------------------
	// Sieve of Eratosthenes
	// ---------------------------------------------------------------------
	// We begin by setting a bitmask to the largest value possible and
	// imediately clear zero and one as trivial cases. We then iterate from
	// two up to and including the square-root of the bit capacity. With
	// each iteration, we check if the current value p is prime. If it is
	// prime, it will be in the bitmask and we can then iterate again for
	// each k from p^2 up to the bit capacity clearing each bit k. We then
	// update p as the next set bit in the bitmask.
	// ---------------------------------------------------------------------

	var (
		sqrtBitCap int   = int(math.Sqrt(float64(BitCap)))
		primes     UMask = Max.ClrBits(0, 1)
	)

	for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
		if primes.MasksBit(p) {
			for k := p * p; k < BitCap; k += p {
				primes = primes.ClrBit(k)
			}
		}
	}

	// isPrime is a simple method determining if a number is prime or not.
	var isPrime = func(n int) bool {
		if n < 2 {
			return false
		}

		for d, r := 2, int(math.Sqrt(float64(n))); d <= r; d++ {
			if n%d == 0 {
				return false
			}
		}

		return true
	}

	for i := 0; i < BitCap; i++ {
		if isPrime(i) {
			if !primes.MasksBit(i) {
				t.Errorf("\nexpected %d to be masked as prime\n", i)
			}
		} else if primes.MasksBit(i) {
			t.Errorf("\nexpected %d to not be masked as prime\n", i)
		}
	}
}

func TestSquares(t *testing.T) {
	// ---------------------------------------------------------------------
	// We generate squares exploiting the following theorem.
	//
	// 	An integer s is a square number if and only if for some odd integer
	//  n, we have s = 1+3+5+...+n.
	//
	// An odd integer n is always of the form 2m+1 for each integer
	// m = 0,1,2,..., so define S(m) = 1+3+5+...+(2m+1). We get the
	// recursive relation S(m-1) + (2m+1) = S(m). Beginning with n = 2*0+1,
	// we set each square bit S(m) in the bitmask by finding the largest
	// set bit, a square S(m-1), then adding n = 2m+1 to it. We never have
	// to track the index m due to the recurrence relation, but the next
	// index m to be set is always equal to the current bitmask count.
	// Similarly, the current largest index m set in the bitmask is the
	// bitmask count minus one.
	// ---------------------------------------------------------------------

	var squares UMask = One                                            // S(0) := 2*0 + 1
	t.Logf("S(%d) = %d\n", squares.Count()-1, squares.PrevBit(BitCap)) // Prints S(0)

	for n := 1; ; n += 2 {
		var s int = squares.PrevBit(BitCap) + n // S(m) = S(m-1) + (2m+1)
		if BitCap <= s {
			break
		}

		t.Logf("S(%d) = %d\n", squares.Count(), s) // Prints S(m), which has not been set in the bitmask yet
		squares = squares.SetBit(s)
	}

	var isSquare = func(n int) bool {
		var r int = int(math.Sqrt(float64(n)))
		return n == r*r
	}

	for n := 0; n < BitCap; n++ {
		if isSquare(n) {
			if !squares.MasksBit(n) {
				t.Errorf("\nexpected %d to be masked as square\n", n)
			}
		} else if squares.MasksBit(n) {
			t.Errorf("\nexpected %d to not be masked as square\n", n)
		}
	}
}

// -------------------------------------------------------------------------
// Benchmarks
// -------------------------------------------------------------------------

func BenchmarkPrimes(b *testing.B) {
	benchmarkPrimes(b)
	benchmarkPrimesNextBit(b)
}

func benchmarkPrimes(b *testing.B) bool {
	f := func(b *testing.B) {
		var (
			primes     UMask
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = Max.ClrBits(0, 1)
			sqrtBitCap = int(math.Sqrt(float64(BitCap)))
			for p := 2; p <= sqrtBitCap; p++ {
				if primes.MasksBit(p) {
					for m := p << 1; m < BitCap; m += p {
						primes = primes.ClrBit(m)
					}
				}
			}
		}

		_, _ = primes, sqrtBitCap
	}

	return b.Run("", f)
}

func benchmarkPrimesNextBit(b *testing.B) bool {
	f := func(b *testing.B) {
		var (
			primes     UMask
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = Max.ClrBits(0, 1)
			sqrtBitCap = int(math.Sqrt(float64(BitCap)))
			for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
				if primes.MasksBit(p) {
					for k := p * p; k < BitCap; k += p {
						primes = primes.ClrBit(k)
					}
				}
			}
		}

		_, _ = primes, sqrtBitCap
	}

	return b.Run("", f)
}
