package lmask

import (
	"fmt"
	"math"
	"testing"
	// yb "github.com/yourbasic/bit"
)

func TestMax(t *testing.T) {
	tests := []struct {
		bitCap int
		exp    *LMask
	}{
		{
			bitCap: 0,
			exp: &LMask{
				bitCap: 0,
				words:  []uint{},
			},
		},
		{
			bitCap: 1,
			exp: &LMask{
				bitCap: 1,
				words:  []uint{1},
			},
		},
		{
			bitCap: WordBitCap - 1,
			exp: &LMask{
				bitCap: WordBitCap - 1,
				words:  []uint{WordMax >> 1},
			},
		},
		{
			bitCap: WordBitCap,
			exp: &LMask{
				bitCap: WordBitCap,
				words:  []uint{WordMax},
			},
		},
		{
			bitCap: 65,
			exp: &LMask{
				bitCap: 65,
				words:  []uint{WordMax, 1},
			},
		},
		{
			bitCap: 127,
			exp: &LMask{
				bitCap: 127,
				words:  []uint{WordMax, WordMax >> 1},
			},
		},
		{
			bitCap: 128,
			exp: &LMask{
				bitCap: 128,
				words:  []uint{WordMax, WordMax},
			},
		},
	}

	for _, test := range tests {
		a := Max(test.bitCap)
		if !a.Masks(test.exp) || !test.exp.Masks(a) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, a)
		}
	}
}

func TestZero(t *testing.T) {
	tests := []struct {
		bitCap int
		bits   []int
		exp    *LMask
	}{
		{
			bitCap: 0,
			bits:   []int{},
			exp: &LMask{
				bitCap: 0,
				words:  []uint{},
			},
		},
		{
			bitCap: 64,
			bits:   []int{2, 3, 5, 7},
			exp: &LMask{
				bitCap: 64,
				words:  []uint{1<<2 | 1<<3 | 1<<5 | 1<<7},
			},
		},
		{
			bitCap: 128,
			bits:   []int{2, 3, 5, 7, 127},
			exp: &LMask{
				bitCap: 128,
				words:  []uint{1<<2 | 1<<3 | 1<<5 | 1<<7, 1 << 63},
			},
		},
	}

	for _, test := range tests {
		a := Zero(test.bitCap).SetBits(test.bits...)
		if !a.Masks(test.exp) || !test.exp.Masks(a) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, a)
		}
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		a, b, exp *LMask
	}{
		{
			a:   Zero(256).SetBits(0, 63, 64, 127, 128, 185, 186, 255),
			b:   Max(256).ClrBits(0, 63, 64, 127, 128, 185, 186, 255),
			exp: Zero(256),
		},
	}

	for _, test := range tests {
		rec := test.a.Copy().And(test.b)
		if !test.exp.Masks(rec) || !rec.Masks(test.exp) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestNot(t *testing.T) {
	tests := []struct {
		a, exp *LMask
	}{
		{
			a:   Zero(256).SetBits(0, 63, 64, 127, 128, 185, 186, 255),
			exp: Max(256).ClrBits(0, 63, 64, 127, 128, 185, 186, 255),
		},
	}

	for _, test := range tests {
		rec := test.a.Copy().Not()
		if !test.exp.Masks(rec) || !rec.Masks(test.exp) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		a, b, exp *LMask
	}{
		{
			a:   Zero(256).SetBits(0, 63, 64, 127, 128, 185, 186, 255),
			b:   Max(256).ClrBits(0, 63, 64, 127, 128, 185, 186, 255),
			exp: Max(256),
		},
	}

	for _, test := range tests {
		rec := test.a.Copy().Or(test.b)
		if !test.exp.Masks(rec) || !rec.Masks(test.exp) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestXOr(t *testing.T) {
	tests := []struct {
		a, b, exp *LMask
	}{
		{
			a:   Zero(256).SetBits(0, 63, 64, 127, 128, 185, 186, 255),
			b:   Max(256).ClrBits(0, 63, 64, 127, 128, 185, 186, 255),
			exp: Max(256),
		},
		{
			a:   Zero(256).SetBits(0, 63, 64, 127, 128, 185, 186, 255),
			b:   Zero(256).SetBits(0, 63, 64, 127, 128, 185, 186, 255),
			exp: Zero(256),
		},
	}

	for _, test := range tests {
		rec := test.a.Copy().XOr(test.b)
		if !test.exp.Masks(rec) || !rec.Masks(test.exp) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

func TestPrimes(t *testing.T) {
	// Sieve of Eratosthenes
	var (
		bitCap     = WordBitCap << 2
		sqrtBitCap = int(math.Sqrt(float64(bitCap)))
		primes     = Max(bitCap).ClrBits(0, 1)
	)

	for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
		if primes.MasksBit(p) {
			for k := p * p; k < bitCap; k += p {
				primes.ClrBit(k)
			}
		}
	}

	// Compare against Your Basic's implementation.
	// Source: https://yourbasic.org/golang/bitmask-flag-set-clear/
	// ybPrimes := yb.New().AddRange(2, bitCap)
	// for p := 2; p <= sqrtBitCap; p = ybPrimes.Next(p) {
	// 	if ybPrimes.Contains(p) {
	// 		for k := p * p; k < bitCap; k += p {
	// 			ybPrimes.Delete(k)
	// 		}
	// 	}
	// }

	// isPrime is a simple method determining if a number is prime or not.
	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}

		r := int(math.Sqrt(float64(n)))
		for d := 2; d <= r; d++ {
			if n%d == 0 {
				return false
			}
		}

		return true
	}

	for i := 0; i < bitCap; i++ {
		if isPrime(i) {
			if !primes.MasksBit(i) {
				t.Errorf("\nexpected %d to be masked as prime\n", i)
			}

			// if !ybPrimes.Contains(i) {
			// 	t.Errorf("\nexpected %d to be masked as prime\n", i)
			// }
		} else {
			if primes.MasksBit(i) {
				t.Errorf("\nexpected %d to not be masked as prime\n", i)
			}

			// if ybPrimes.Contains(i) {
			// 	t.Errorf("\nexpected %d to be masked as prime\n", i)
			// }
		}
	}
}

func TestSquares(t *testing.T) {
	bitCap := WordBitCap << 2
	squares := Zero(bitCap)
	for i := 0; i < bitCap; i++ {
		if i2 := i * i; i2 < bitCap {
			squares = squares.SetBits(i2)
		}
	}

	isSquare := func(n int) bool {
		r := int(math.Sqrt(float64(n)))
		return n == r*r
	}

	for n := 0; n < bitCap; n++ {
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
	n0, n1, dn := 0, 256, 8
	for n := n0; n <= n1; n += dn {
		benchmarkPrimes(b, n)
	}

	for n := n0; n <= n1; n += dn {
		benchmarkPrimesNextBit(b, n)
	}

	// for n := n0; n <= n1; n += dn {
	// 	benchmarkYourBasicPrimes(b, n)
	// }

	benchmarkPrimes(b, 50000000)
	benchmarkPrimesNextBit(b, 50000000)
	// benchmarkYourBasicPrimes(b, 50000000)
}

func benchmarkPrimes(b *testing.B, bitCap int) bool {
	f := func(b *testing.B) {
		var (
			primes     *LMask
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = Max(bitCap).ClrBit(0).ClrBit(1)
			sqrtBitCap = int(math.Sqrt(float64(bitCap)))
			for p := 2; p <= sqrtBitCap; p++ {
				if primes.MasksBit(p) {
					for m := p << 1; m < bitCap; m += p {
						primes.ClrBit(m)
					}
				}
			}
		}

		_, _ = primes, sqrtBitCap
	}

	return b.Run(fmt.Sprintf("LMask: bit cap %d", bitCap), f)
}

func benchmarkPrimesNextBit(b *testing.B, bitCap int) bool {
	f := func(b *testing.B) {
		var (
			primes     *LMask
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = Max(bitCap).ClrBit(0).ClrBit(1)
			sqrtBitCap = int(math.Sqrt(float64(bitCap)))
			for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
				if primes.MasksBit(p) {
					for k := p * p; k < bitCap; k += p {
						primes.ClrBit(k)
					}
				}
			}
		}

		_, _ = primes, sqrtBitCap
	}

	return b.Run(fmt.Sprintf("LMask w/ NextBit: bit cap %d", bitCap), f)
}

// func benchmarkYourBasicPrimes(b *testing.B, bitCap int) bool {
// 	// Source: https://yourbasic.org/golang/bitmask-flag-set-clear/
// 	f := func(b *testing.B) {
// 		var (
// 			primes     *yb.Set
// 			sqrtBitCap int
// 		)

// 		for i := 0; i < b.N; i++ {
// 			primes = yb.New().AddRange(2, bitCap)
// 			sqrtBitCap = int(math.Sqrt(float64(bitCap)))
// 			for p := 2; p <= sqrtBitCap; p = primes.Next(p) {
// 				if primes.Contains(p) {
// 					for k := p * p; k < bitCap; k += p {
// 						primes.Delete(k)
// 					}
// 				}
// 			}
// 		}

// 		_, _ = primes, sqrtBitCap
// 	}

// 	return b.Run(fmt.Sprintf("YourBasicSet: bit cap %d", bitCap), f)
// }
