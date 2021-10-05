package umask

import (
	"math"
	"testing"
)

// --------------------------------------------------------------------
// TODO: FINISH TESTING
// --------------------------------------------------------------------

func TestAnd(t *testing.T) {
	tests := []struct {
		a, b, exp UMask
	}{
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

	for _, test := range tests {
		if rec := test.a.And(test.b); test.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestBitLen(t *testing.T) {
	var exp int
	a := Zero
	if rec := a.BitLen(); exp != rec {
		t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
	}

	for exp++; exp < BitCap; exp++ {
		a := Zero.SetBits(exp - 1)
		if rec := a.BitLen(); exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
		}
	}
}

func TestBits(t *testing.T) {
	tests := [][]int{
		{},
		{0},
		{0, BitCap - 1},
	}

	for _, test := range tests {
		var (
			rec   = Zero.SetBits(test...).Bits()
			equal = true
		)

		if len(test) != len(rec) {
			equal = false
		}

		for i := 0; i < len(test) && equal; i++ {
			equal = test[i] == rec[i]
		}

		if !equal {
			t.Errorf("\nexpected %d\nreceived %d\n", test, rec)
		}
	}
}

func TestClr(t *testing.T) {
	tests := []struct {
		a, b, exp UMask
	}{
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

	for _, test := range tests {
		if rec := test.a.Clr(test.b); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}
	}
}

func TestClrBit(t *testing.T) {
	tests := []struct {
		a, exp UMask
		bit    int
	}{
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

	for _, test := range tests {
		if rec := test.a.ClrBit(test.bit); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}
	}
}

func TestClrBits(t *testing.T) {
	tests := []struct {
		a, b, exp UMask
	}{
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

	for _, test := range tests {
		if rec := test.a.ClrBits(test.b.Bits()...); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}
	}
}

func TestLRSh(t *testing.T) {
	tests := []struct {
		a, expLeft, expRight UMask
		left, right          int
	}{
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

	for _, test := range tests {
		if rec := test.a.LSh(test.left); test.expLeft != rec {
			t.Errorf("\nexpected %s\nreceived %s\n", test.expLeft.Fmt(2), rec.Fmt(2))
		}

		if rec := test.a.RSh(test.right); test.expRight != rec {
			t.Errorf("\nexpected %s\nreceived %s\n", test.expRight.Fmt(2), rec.Fmt(2))
		}
	}
}

func TestNextPrevBit(t *testing.T) {
	tests := []UMask{
		Zero,
		One,
		Zero.SetBits(0, BitCap-1), // End bits
		Max.ClrBits(0, BitCap-1),  // Middle bits
		Max,
		Zero.SetBits(0, 2, 4, 6, BitCap-8, BitCap-6, BitCap-4, BitCap-2), // Some even bits
		Zero.SetBits(1, 3, 5, 7, BitCap-7, BitCap-5, BitCap-3, BitCap-1), // Some odd bits
	}

	for _, test := range tests {
		for i := 0; i < BitCap; i++ {
			if test.MasksBit(i) {
				if rec := test.NextBit(i - 1); i != rec {
					t.Errorf("\nexpected next bit to be %d\nreceived next bit %d\n", i, rec)
				}

				if rec := test.PrevBit(i + 1); i != rec {
					t.Errorf("\nexpected previous bit to be %d\nreceived next bit %d\n", i, rec)
				}
			} else {
				if rec := test.NextBit(i - 1); i == rec {
					t.Errorf("\nexpected next bit to NOT be %d\nreceived next bit %d\n", i, rec)
				}

				if rec := test.PrevBit(i + 1); i == rec {
					t.Errorf("\nexpected previous bit to NOT be %d\nreceived next bit %d\n", i, rec)
				}
			}
		}

		{
			exp := test
			rec := Zero
			for i := test.NextBit(-1); i < BitCap; i = test.NextBit(i) {
				rec = rec.SetBit(i)
			}

			if exp != rec {
				t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
			}
		}

		{
			exp := Zero
			rec := test
			for i := test.PrevBit(BitCap); -1 < i; i = test.PrevBit(i) {
				rec = rec.ClrBit(i)
			}

			if exp != rec {
				t.Errorf("\nexpected %v\nreceived %v\n", exp, rec)
			}
		}
	}
}

func TestNot(t *testing.T) {
	tests := []struct {
		a, exp UMask
	}{
		{
			a:   Zero.SetBits(0, BitCap-1),
			exp: Max.ClrBits(0, BitCap-1),
		},
	}

	for _, test := range tests {
		if rec := test.a.Not(); test.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		a, b, exp UMask
	}{
		{
			a:   Zero.SetBits(0, BitCap-1),
			b:   Max.ClrBits(0, BitCap-1),
			exp: Max,
		},
	}

	for _, test := range tests {
		if rec := test.a.Or(test.b); test.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestXOr(t *testing.T) {
	tests := []struct {
		a, b, exp UMask
	}{
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

	for _, test := range tests {
		if rec := test.a.XOr(test.b); test.exp != rec {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

// --------------------------------------------------------------------
// Applications
// --------------------------------------------------------------------

func TestFactor(t *testing.T) {
	gcd := func(a, b int) int {
		for 0 < b {
			a, b = b, a%b
		}

		return a
	}

	n0, n1, dn := 0, BitCap, 1
	for n := n0; n <= n1; n += dn {
		var (
			factors = Zero.SetBits(1, n)
			i, j    = 2, n >> 1
			d       = i * j
		)

		for i <= j {
			switch {
			case d < n:
				i++
				d += j
			case n < d:
				j--
				d -= i
			default:
				factors = factors.SetBits(i, j)
				i++
				j--
				d += j - i + 1
			}
		}

		for m := 0; m <= n; m++ {
			if gcd(m, n) == m {
				if !factors.MasksBit(m) {
					t.Errorf("\nexpected %d to be set\n", m)
				}
			} else if factors.MasksBit(m) {
				t.Errorf("\nexpected %d to not be set\n", m)
			}
		}
	}
}

func TestFibonacciNumbers(t *testing.T) {
	fibs := Zero.SetBits(0, 1, 2)
	for n := fibs.Count(); n < BitCap; n++ {
		b := fibs.PrevBit(BitCap)
		fibs = fibs.SetBit(fibs.PrevBit(b) + b)
	}

	for a0, a1 := 0, 1; a1 < BitCap; a0, a1 = a1, a0+a1 {
		if !fibs.MasksBit(a1) {
			t.Errorf("\nexpected %d to be masked\n", a1)
		}
	}
}

func TestPrimes(t *testing.T) {
	// Sieve of Eratosthenes

	var (
		sqrtBitCap = int(math.Sqrt(float64(BitCap)))
		primes     = Max.ClrBits(0, 1)
	)

	for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
		if primes.MasksBit(p) {
			for k := p * p; k < BitCap; k += p {
				primes = primes.ClrBit(k)
			}
		}
	}

	// isPrime is a simple method determining if a number is prime or not.
	isPrime := func(n int) bool {
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
	squares := One
	for i := 1; ; i += 2 {
		s := squares.PrevBit(BitCap) + i
		if BitCap <= s {
			break
		}

		squares = squares.SetBit(s)
	}

	isSquare := func(n int) bool {
		r := int(math.Sqrt(float64(n)))
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
			primes = Max.ClrBit(0).ClrBit(1)
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
			primes = Max.ClrBit(0).ClrBit(1)
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
