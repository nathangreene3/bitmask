package bitmask

import (
	"math"
	"math/bits"
	"testing"
)

// ------------------------------------------------------------------------------------
// Demos
// ------------------------------------------------------------------------------------

func TestPrimes(t *testing.T) {
	primes := ClrBits(Max, 0, 1)
	for p, sqrtBitCap := 2, int(math.Sqrt(float64(BitCap))); p <= sqrtBitCap; p = NextBit(primes, p) {
		if MasksBit(primes, p) {
			for k := p * p; k < BitCap; k += p {
				primes = ClrBit(primes, k)
			}
		}
	}

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
			if !MasksBit(primes, i) {
				t.Errorf("expected %d to be masked as prime\n", i)
			}
		} else if MasksBit(primes, i) {
			t.Errorf("expected %d to not be masked as prime\n", i)
		}
	}
}

func TestSquares(t *testing.T) {
	var squares uint
	for n := 0; n < BitCap; n++ {
		if n2 := n * n; n2 < BitCap {
			squares = SetBit(squares, n2)
		}
	}

	isSquare := func(n int) bool {
		r := int(math.Sqrt(float64(n)))
		return n == r*r
	}

	for n := 0; n < BitCap; n++ {
		if MasksBit(squares, n) {
			if !isSquare(n) {
				t.Errorf("\nexpected %d to be square\n", n)
			}
		} else if isSquare(n) {
			t.Errorf("\nexpected %d to not be square\n", n)
		}
	}
}

// ------------------------------------------------------------------------------------
// Focused tests
// ------------------------------------------------------------------------------------

func TestNextBit(t *testing.T) {
	nextBit0 := func(a uint, bit int) int {
		if bit = clamp(bit+1, 0, BitCap); bit < BitCap {
			if a = a >> bit << bit; 0 < a {
				return bits.TrailingZeros(a)
			}
		}

		return BitCap
	}

	nextBit1 := func(a uint, bit int) int {
		bit = clamp(bit+1, 0, BitCap)
		if a = a >> bit << bit; 0 < a {
			return bits.TrailingZeros(a)
		}

		return BitCap
	}

	nextBit2 := func(a uint, bit int) int {
		bit = clamp(bit+1, 0, BitCap)
		return bits.TrailingZeros(a >> bit << bit)
	}

	tests := []struct {
		a        uint
		bit, exp int
	}{
		{
			a:   0,
			bit: -1,
			exp: BitCap,
		},
		{
			a:   0,
			bit: 0,
			exp: BitCap,
		},
		{
			a:   1,
			bit: -1,
			exp: 0,
		},
		{
			a:   1,
			bit: 0,
			exp: BitCap,
		},
		{
			a:   2,
			bit: -1,
			exp: 1,
		},
		{
			a:   2,
			bit: 0,
			exp: 1,
		},
		{
			a:   2,
			bit: 1,
			exp: BitCap,
		},
		{
			a:   3,
			bit: -1,
			exp: 0,
		},
		{
			a:   3,
			bit: 0,
			exp: 1,
		},
		{
			a:   3,
			bit: 1,
			exp: BitCap,
		},
		{
			a:   1 << (BitCap - 2),
			bit: -1,
			exp: BitCap - 2,
		},
		{
			a:   1 << (BitCap - 2),
			bit: BitCap - 3,
			exp: BitCap - 2,
		},
		{
			a:   1 << (BitCap - 2),
			bit: BitCap - 1,
			exp: BitCap,
		},
		{
			a:   1 | 1<<(BitCap-2),
			bit: -1,
			exp: 0,
		},
		{
			a:   1 | 1<<(BitCap-2),
			bit: 0,
			exp: BitCap - 2,
		},
		{
			a:   1 | 1<<(BitCap-2),
			bit: BitCap - 3,
			exp: BitCap - 2,
		},
		{
			a:   1 | 1<<(BitCap-2),
			bit: BitCap - 2,
			exp: BitCap,
		},
		{
			a:   1 << (BitCap - 1),
			bit: -1,
			exp: BitCap - 1,
		},
		{
			a:   1 << (BitCap - 1),
			bit: BitCap - 2,
			exp: BitCap - 1,
		},
		{
			a:   1 << (BitCap - 1),
			bit: BitCap - 1,
			exp: BitCap,
		},
		{
			a:   1 | 1<<(BitCap-1),
			bit: -1,
			exp: 0,
		},
		{
			a:   1 | 1<<(BitCap-1),
			bit: 0,
			exp: BitCap - 1,
		},
		{
			a:   1 | 1<<(BitCap-1),
			bit: BitCap - 2,
			exp: BitCap - 1,
		},
		{
			a:   1 | 1<<(BitCap-1),
			bit: BitCap - 1,
			exp: BitCap,
		},
	}

	for _, test := range tests {
		if rec := nextBit0(test.a, test.bit); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}

		if rec := nextBit1(test.a, test.bit); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}

		if rec := nextBit2(test.a, test.bit); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}
	}
}

func TestPrevBit(t *testing.T) {
	tests := []struct {
		a        uint
		bit, exp int
	}{
		{
			a:   0,
			bit: -1,
			exp: -1,
		},
		{
			a:   0,
			bit: 0,
			exp: -1,
		},
		{
			a:   0,
			bit: BitCap,
			exp: -1,
		},
		{
			a:   1,
			bit: -1,
			exp: -1,
		},
		{
			a:   1,
			bit: 0,
			exp: -1,
		},
		{
			a:   1,
			bit: 1,
			exp: 0,
		},
		{
			a:   1,
			bit: BitCap,
			exp: 0,
		},
		{
			a:   2,
			bit: -1,
			exp: -1,
		},
		{
			a:   2,
			bit: 0,
			exp: -1,
		},
		{
			a:   2,
			bit: 1,
			exp: -1,
		},
		{
			a:   2,
			bit: BitCap,
			exp: 1,
		},
		{
			a:   3,
			bit: -1,
			exp: -1,
		},
		{
			a:   3,
			bit: 0,
			exp: -1,
		},
		{
			a:   3,
			bit: 1,
			exp: 0,
		},
		{
			a:   3,
			bit: BitCap,
			exp: 1,
		},
		{
			a:   1 << (BitCap - 1),
			bit: -1,
			exp: -1,
		},
		{
			a:   1 << (BitCap - 1),
			bit: BitCap - 2,
			exp: -1,
		},
		{
			a:   1 << (BitCap - 1),
			bit: BitCap - 1,
			exp: -1,
		},
		{
			a:   1 << (BitCap - 1),
			bit: BitCap,
			exp: BitCap - 1,
		},
	}

	for _, test := range tests {
		if rec := PrevBit(test.a, test.bit); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}
	}
}

func TestMaxInt(t *testing.T) {
	intMax := 1<<(BitCap-1) - 1
	intMin := -intMax

	max0 := func(a, b int) int {
		if a < b {
			return b
		}

		return a
	}

	max1 := func(a, b int) int {
		// Source: https://web.archive.org/web/20130821015554/http://bob.allegronetwork.com/prog/tricks.html
		a -= b
		a &= ^a >> (BitCap - 1)
		a += b
		return a
	}

	tests := []struct {
		a, b, exp int
	}{
		{
			a:   0,
			b:   0,
			exp: 0,
		},
		{
			a:   0,
			b:   1,
			exp: 1,
		},
		{
			a:   1,
			b:   0,
			exp: 1,
		},
		{
			a:   0,
			b:   1 << (BitCap - 2),
			exp: 1 << (BitCap - 2),
		},
		{
			a:   1 << (BitCap - 2),
			b:   0,
			exp: 1 << (BitCap - 2),
		},
		{
			a:   0,
			b:   intMin,
			exp: 0,
		},
		{
			a:   intMin,
			b:   0,
			exp: 0,
		},
		{
			a:   intMin,
			b:   intMax,
			exp: intMax,
		},
		{
			a:   intMax,
			b:   intMin,
			exp: intMax,
		},
	}

	for _, test := range tests {
		if rec := max0(test.a, test.b); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}

		if rec := max1(test.a, test.b); test.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", test.exp, rec)
		}
	}
}
