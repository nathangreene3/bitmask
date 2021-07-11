package umask

import (
	"math"
	"testing"
)

func TestBitmask(t *testing.T) {
	{
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

	{
		// Squares
		var squares UMask
		for i := 0; i < BitCap; i++ {
			if i2 := i * i; i2 < BitCap {
				squares = squares.SetBits(i2)
			}
		}

		isSquare := func(n int) bool {
			r := int(math.Sqrt(float64(n)))
			return n == r*r
		}

		for n := 0; n < BitCap; n++ {
			if squares.MasksBit(n) {
				if !isSquare(n) {
					t.Errorf("\nexpected %d to be square\n", n)
				}
			} else if isSquare(n) {
				t.Errorf("\nexpected %d to not be square\n", n)
			}
		}
	}

	{
		// Fibs
		// fibs:=Zero().SetBits(0,1)
		// for i:=
	}
}
