package bitmask

import (
	"math"
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

}
