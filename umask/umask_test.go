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
			p       = i * j
		)

		for i <= j {
			if p < n {
				i++
				p += j
				continue
			}

			if n < p {
				j--
				p -= i
				continue
			}

			if p == n {
				factors = factors.SetBits(i, j)
			}

			i++
			j--
			p += j - i + 1
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
	var (
		fibs = Zero.SetBits(0, 1, 2)
	)

	for n := fibs.Count(); n < BitCap; n++ {
		var (
			b0 = fibs.PrevBit(BitCap)
			b1 = fibs.PrevBit(b0) + b0
		)

		fibs = fibs.SetBit(b1)
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
		square := squares.PrevBit(BitCap) + i
		if BitCap <= square {
			break
		}

		squares = squares.SetBit(square)
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
