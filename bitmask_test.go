package bitmask

import (
	"math"
	"testing"
)

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
			divisors uint = SetBits(0)
			a, b     int  = 1, n
			d        int  = n // a*b
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
				divisors = SetBits(divisors, a, b)
				a++
				b--
				d += b - a + 1
			}
		}

		for m := 0; m <= n; m++ {
			if gcd(m, n) == m {
				if !MasksBit(divisors, m) {
					t.Errorf("\nexpected %d to be set\n", m)
				}
			} else if MasksBit(divisors, m) {
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

	var fibs uint = SetBits(0, 1, 2)
	for n := Count(fibs); n < BitCap; n++ {
		var b int = PrevBit(fibs, BitCap)
		fibs = SetBit(fibs, PrevBit(fibs, b)+b)
	}

	for a0, a1 := 0, 1; a1 < BitCap; a0, a1 = a1, a0+a1 {
		if !MasksBit(fibs, a1) {
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
		sqrtBitCap int  = int(math.Sqrt(float64(BitCap)))
		primes     uint = ClrBits(Max, 0, 1)
	)

	for p := 2; p <= sqrtBitCap; p = NextBit(primes, p) {
		if MasksBit(primes, p) {
			for k := p * p; k < BitCap; k += p {
				primes = ClrBit(primes, k)
			}
		}
	}

	// isPrime determines if a number is prime.
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
			if !MasksBit(primes, i) {
				t.Errorf("expected %d to be masked as prime\n", i)
			}
		} else if MasksBit(primes, i) {
			t.Errorf("expected %d to not be masked as prime\n", i)
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

	var squares uint = 1                                               // S(0) = 2*0 + 1
	t.Logf("S(%d) = %d\n", Count(squares)-1, PrevBit(squares, BitCap)) // Prints S(0)

	for n := 1; ; n += 2 {
		var s int = PrevBit(squares, BitCap) + n // S(m) = S(m-1) + (2*m+1)
		if BitCap <= s {
			break
		}

		t.Logf("S(%d) = %d\n", Count(squares), s) // Prints S(m), which has not been set in the bitmask yet
		squares = SetBit(squares, s)
	}

	// isSquare determines if a number is square.
	var isSquare = func(n int) bool {
		var r int = int(math.Sqrt(float64(n)))
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
			primes     uint = Max
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = ClrBits(primes, 0, 1)
			sqrtBitCap = int(math.Sqrt(float64(BitCap)))
			for p := 2; p <= sqrtBitCap; p++ {
				if MasksBit(primes, p) {
					for m := p << 1; m < BitCap; m += p {
						primes = ClrBit(primes, m)
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
			primes     uint = Max
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = ClrBits(primes, 0, 1)
			sqrtBitCap = int(math.Sqrt(float64(BitCap)))
			for p := 2; p <= sqrtBitCap; p = NextBit(primes, p) {
				if MasksBit(primes, p) {
					for k := p * p; k < BitCap; k += p {
						primes = ClrBit(primes, k)
					}
				}
			}
		}

		_, _ = primes, sqrtBitCap
	}

	return b.Run("", f)
}
