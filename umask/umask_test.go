package umask

import (
	"math"
	"testing"
)

func TestPrimes(t *testing.T) {
	// Sieve of Eratosthenes
	primes := Max.ClrBits(0, 1) // Zero and one are not prime, but cannot be ruled out by this method without being artificially removed at initialization.
	for p, sqrtBitCap := 2, int(math.Sqrt(float64(BitCap))); p <= sqrtBitCap; p = primes.NextBit(p) {
		if primes.MasksBit(p) {
			for k := p * p; k < BitCap; k += p {
				primes = primes.ClrBit(k)
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
			if !primes.MasksBit(i) {
				t.Errorf("expected %d to be masked as prime\n", i)
			}
		} else if primes.MasksBit(i) {
			t.Errorf("expected %d to not be masked as prime\n", i)
		}
	}
}

func TestSquares(t *testing.T) {
	var (
		squares    = One
		sqrtBitCap = int(math.Sqrt(float64(BitCap)))
		i1         = sqrtBitCap + int(math.Log2(float64(sqrtBitCap)))<<1 - 1
	)

	for i := 1; i <= i1; i += 2 {
		// 1+3+5+...+(2n-1) is odd for n in N
		squares = squares.SetBit(squares.PrevBit(BitCap) + i)
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

func TestNextPrevBit(t *testing.T) {
	if exp, rec := 63, One.ShiftLeft(63).PrevBit(BitCap); exp != rec {
		t.Fatalf("\nexpected %d\nreceived %d\n", exp, rec)
	}

	tests := []UMask{
		Zero,
		One,
		Zero.SetBits(0, 63),                      // End bits
		Zero.SetBits(0, 2, 4, 6, 56, 58, 60, 62), // Some even bits
		Zero.SetBits(1, 3, 5, 7, 57, 59, 61, 63), // Some odd bits
		Max.ClrBits(0, 63),                       // Middle bits
		Max,
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
			rec := Zero
			for i := test.NextBit(-1); i < BitCap; i = test.NextBit(i) {
				rec = rec.SetBit(i)
			}

			if exp := test; exp != rec {
				t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
			}
		}

		{
			rec := test
			for i := test.PrevBit(BitCap); -1 < i; i = test.PrevBit(i) {
				rec = rec.ClrBit(i)
			}

			if exp := Zero; exp != rec {
				t.Errorf("\nexpected %d\nreceived %d\n", test, rec)
			}
		}
	}
}
