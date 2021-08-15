package lmask

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestFromBits(t *testing.T) {
	tests := []struct {
		bitCap int
		bits   []int
		exp    *LMask
	}{
		{
			bitCap: 256,
			bits:   []int{0, 63, 64, 127, 128, 191, 192, 255},
			exp: &LMask{
				bitCap: 256,
				words:  []uint{1<<0 | 1<<63, 1<<0 | 1<<63, 1<<0 | 1<<63, 1<<0 | 1<<63},
			},
		},
	}

	for _, test := range tests {
		if rec := FromBits(test.bitCap, test.bits...); !test.exp.Equals(rec) {
			t.Errorf("\nexpected %q\nreceived %q\n", test.exp, rec)
		}
	}
}

func TestFromWords(t *testing.T) {
	tests := []struct {
		words []uint
		exp   *LMask
	}{
		{
			words: nil,
			exp:   &LMask{bitCap: 0, words: []uint{}},
		},
		{
			words: []uint{uint(1<<64 - 1), uint(1<<64 - 1), uint(1<<64 - 1), uint(1<<64 - 1)},
			exp:   &LMask{bitCap: 256, words: []uint{uint(1<<64 - 1), uint(1<<64 - 1), uint(1<<64 - 1), uint(1<<64 - 1)}},
		},
	}

	for _, test := range tests {
		if rec := FromWords(test.words...); !test.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.exp, rec)
		}
	}
}

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

func TestOne(t *testing.T) {
	tests := []struct {
		bitCap int
		bits   []int
		exp    *LMask
	}{
		{
			bitCap: 64,
			bits:   []int{2, 3, 5, 7},
			exp: &LMask{
				bitCap: 64,
				words:  []uint{1 | 1<<2 | 1<<3 | 1<<5 | 1<<7},
			},
		},
		{
			bitCap: 128,
			bits:   []int{2, 3, 5, 7, 127},
			exp: &LMask{
				bitCap: 128,
				words:  []uint{1 | 1<<2 | 1<<3 | 1<<5 | 1<<7, 1 << 63},
			},
		},
	}

	for _, test := range tests {
		a := One(test.bitCap).SetBits(test.bits...)
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
		{
			a:   Max(256).ClrBits(0, 63, 64, 127, 128, 185, 186, 255),
			b:   Max(256),
			exp: Max(256).ClrBits(0, 63, 64, 127, 128, 185, 186, 255),
		},
	}

	for _, test := range tests {
		rec := test.a.Copy().And(test.b)
		if !test.exp.Masks(rec) || !rec.Masks(test.exp) || !test.exp.Equals(rec) {
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

func TestJSON(t *testing.T) {
	tests := []struct {
		expLMask *LMask
		expJSON  string
	}{
		{
			expLMask: nil,
			expJSON:  "null",
		},
		{
			expLMask: Zero(256),
			expJSON:  "0",
		},
		{
			expLMask: One(256),
			expJSON:  "1",
		},
		{
			expLMask: FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			expJSON:  big.NewInt(0).SetBits([]big.Word{1<<0 | 1<<63, 1<<0 | 1<<63, 1<<0 | 1<<63, 1<<0 | 1<<63}).Text(10),
		},
		{
			expLMask: Max(256),
			expJSON:  big.NewInt(0).SetBits([]big.Word{1<<64 - 1, 1<<64 - 1, 1<<64 - 1, 1<<64 - 1}).Text(10),
		},
		{
			// Test max with bit cap that isn't a multiple of a word bit cap
			expLMask: Max(255),
			expJSON:  big.NewInt(0).SetBits([]big.Word{1<<64 - 1, 1<<64 - 1, 1<<64 - 1, 1<<63 - 1}).Text(10),
		},
	}

	for _, test := range tests {
		recJSON := test.expLMask.JSON()
		if test.expJSON != recJSON {
			t.Errorf("\nexpected %q\nreceived %q\n", test.expJSON, recJSON)
			continue
		}

		recLMask, err := FromJSON(recJSON)
		if err != nil {
			t.Error(err)
			continue
		}

		if test.expLMask == nil {
			if recLMask.bitCap != 0 || len(recLMask.words) != 0 {
				t.Errorf("\nexpected %v\nreceived %v\n", Zero(0), recLMask)
			}

			continue
		}

		recLMask.SetBitCap(test.expLMask.BitCap())
		if !test.expLMask.Equals(recLMask) {
			t.Errorf("\nexpected %v\nreceived %v\n", test.expLMask, recLMask)
		}
	}
}

func TestLeftRight(t *testing.T) {
	tests := []struct {
		a, expLeft, expRight *LMask
		left, right          int
	}{
		// Shifts less than 64
		{
			a:        FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			left:     0,
			right:    0,
			expLeft:  FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			expRight: FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
		},
		{
			a:        FromBits(64, 0, 63),
			left:     1,
			right:    1,
			expLeft:  FromBits(64, 1),
			expRight: FromBits(64, 62),
		},
		{
			a:        FromBits(128, 0, 63, 64, 127),
			left:     1,
			right:    1,
			expLeft:  FromBits(128, 1, 64, 65),
			expRight: FromBits(128, 62, 63, 126),
		},
		{
			a:        FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			left:     1,
			right:    1,
			expLeft:  FromBits(256, 1, 64, 65, 128, 129, 192, 193),
			expRight: FromBits(256, 62, 63, 126, 127, 190, 191, 254),
		},
		{
			a:        FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			left:     63,
			right:    63,
			expLeft:  FromBits(256, 63, 126, 127, 190, 191, 254, 255),
			expRight: FromBits(256, 0, 1, 64, 65, 128, 129, 192),
		},

		// Shifts larger than or equal to 64
		{
			a:        FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			left:     64,
			right:    64,
			expLeft:  FromBits(256, 64, 127, 128, 191, 192, 255),
			expRight: FromBits(256, 0, 63, 64, 127, 128, 191),
		},
		{
			a:        FromBits(256, 0, 63, 64, 127, 128, 191, 192, 255),
			left:     128,
			right:    128,
			expLeft:  FromBits(256, 128, 191, 192, 255),
			expRight: FromBits(256, 0, 63, 64, 127),
		},
	}

	for _, test := range tests {
		if rec := test.a.Copy().Left(test.left); !test.expLeft.Equals(rec) {
			t.Errorf("\nexpected %s\nreceived %s\n", test.expLeft.Fmt(2), rec.Fmt(2))
		}

		if rec := test.a.Copy().Right(test.right); !test.expRight.Equals(rec) {
			t.Errorf("\nexpected %s\nreceived %s\n", test.expRight.Fmt(2), rec.Fmt(2))
		}
	}
}

func TestNextPrevBit(t *testing.T) {
	for bitCap := 0; bitCap <= 256; bitCap += 8 {
		tests := append(make([]LMask, 0, 8), *Zero(bitCap))
		if 0 < bitCap {
			tests = append(
				tests,
				*One(bitCap),
				*Zero(bitCap).SetBits(0, bitCap-1), // End bits
				*Max(bitCap).ClrBits(0, bitCap-1),  // Middle bits
				*Max(bitCap),
			)
		}

		if 8 <= bitCap {
			tests = append(
				tests,
				*Zero(bitCap).SetBits(0, 2, 4, 6, bitCap-8, bitCap-6, bitCap-4, bitCap-2), // Some even bits
				*Zero(bitCap).SetBits(1, 3, 5, 7, bitCap-7, bitCap-5, bitCap-3, bitCap-1), // Some odd bits
			)
		}

		for _, test := range tests {
			for i := 0; i < test.bitCap; i++ {
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
				exp := test.Copy()
				rec := Zero(test.bitCap)
				for i := test.NextBit(-1); i < test.bitCap; i = test.NextBit(i) {
					rec.SetBit(i)
				}

				if !exp.Equals(rec) {
					t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
				}
			}

			{
				exp := Zero(test.bitCap)
				rec := test.Copy()
				for i := test.PrevBit(test.bitCap); -1 < i; i = test.PrevBit(i) {
					rec.ClrBit(i)
				}

				if !exp.Equals(rec) {
					t.Errorf("\nexpected %v\nreceived %v\n", exp, rec)
				}
			}
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		a   *LMask
		exp string
	}{
		{
			a:   nil,
			exp: "<nil>",
		},
		{
			a:   Zero(256),
			exp: "0",
		},
		{
			a:   Max(256),
			exp: big.NewInt(0).SetBits([]big.Word{1<<64 - 1, 1<<64 - 1, 1<<64 - 1, 1<<64 - 1}).Text(10),
		},
		{
			//
			a:   Max(255),
			exp: big.NewInt(0).SetBits([]big.Word{1<<64 - 1, 1<<64 - 1, 1<<64 - 1, 1<<63 - 1}).Text(10),
		},
	}

	for _, test := range tests {
		if rec := test.a.String(); test.exp != rec {
			t.Errorf("\nexpected %q\nrecieved %q\n", test.exp, rec)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		a, exp *LMask
	}{
		{
			a:   &LMask{bitCap: 0, words: []uint{}},
			exp: &LMask{bitCap: 0, words: []uint{}},
		},
		{
			a:   &LMask{bitCap: 1, words: []uint{1<<WordBitCap - 1}},
			exp: &LMask{bitCap: 1, words: []uint{1}},
		},
		{
			a:   &LMask{bitCap: 63, words: []uint{1<<WordBitCap - 1}},
			exp: &LMask{bitCap: 63, words: []uint{1<<(WordBitCap-1) - 1}},
		},
		{
			a:   &LMask{bitCap: 193, words: []uint{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1}},
			exp: &LMask{bitCap: 193, words: []uint{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1}},
		},
	}

	for _, test := range tests {
		if rec := test.a.Copy().trim(); !test.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nrecieved %v\n", test.exp, rec)
		}
	}
}

// --------------------------------------------------------------------
// Applications
// --------------------------------------------------------------------

// TestFibonacciNumbers generates n Fibonacci numbers. This test
// demonstrates updating the bit capacity when the ith bit is not known
// in advance. That is, we don't know what the bit capacity should be
// when computing n Fibonacci numbers.
func TestFibonacciNumbers(t *testing.T) {
	fibs := Zero(4).SetBits(0, 1, 2)
	maxCount := 50
	for n := fibs.Count(); n < maxCount; n++ {
		b0 := fibs.PrevBit(fibs.BitCap())
		b1 := fibs.PrevBit(b0) + b0
		if bc := fibs.BitCap(); bc <= b1 {
			fibs.SetBitCap(b1 << 1) // Double the current bit capacity similar to append
		}

		fibs.SetBit(b1)
	}

	for a0, a1 := 0, 1; a1 < fibs.BitCap(); a0, a1 = a1, a0+a1 {
		if !fibs.MasksBit(a1) {
			t.Errorf("\nexpected %d to be masked\n", a1)
		}
	}

	if recCount := fibs.Count(); maxCount != recCount {
		t.Errorf("\nexpected count of %d\nreceived %d\n", maxCount, recCount)
	}
}

// TestPrimes generates all primes less than the given bit capacity.
// This test demonstrates setting, querying, and clearing bits within
// the constraint of a fixed bit capacity. Iteration through the
// bitmask is also determined entirely by querying the next bit set and
// acting accordingly.
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

	for i := 0; i < bitCap; i++ {
		if isPrime(i) {
			if !primes.MasksBit(i) {
				t.Errorf("\nexpected %d to be masked as prime\n", i)
			}
		} else if primes.MasksBit(i) {
			t.Errorf("\nexpected %d to not be masked as prime\n", i)
		}
	}
}

// TestSquares computes all squares within a given bit capacity. This
// test demonstrates setting a bit using information retreived from the
// previous bit.
func TestSquares(t *testing.T) {
	// 1+3+5+...+(2n-1) is odd for n in N

	var (
		bitCap  = WordBitCap + 1
		squares = One(bitCap)
	)

	for i := 1; ; i += 2 {
		square := squares.PrevBit(bitCap) + i
		if bitCap <= square {
			break
		}

		squares.SetBit(square)
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
