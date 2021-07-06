package umask

import (
	"math"
	"testing"
)

func TestBitmask(t *testing.T) {
	{
		// Sieve of Eratosthenes
		primes := Max.ClrBits(0, 1) // Zero and one are not prime, but cannot be ruled out by this method without being artificially removed at initialization.
		for m := 2; m < WordBitCap; m++ {
			if primes.MasksBit(m) {
				for n := m << 1; n < WordBitCap; n += m {
					primes = primes.ClrBits(n)
				}
			}
		}

		isPrime := func(n int) bool {
			for d := 2; d < n; d++ {
				if n%d == 0 {
					return false
				}
			}

			return 1 < n
		}

		for i := 0; i < WordBitCap; i++ {
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
		for i := 0; i < WordBitCap; i++ {
			if i2 := i * i; i2 < WordBitCap {
				squares = squares.SetBits(i2)
			}
		}

		isSquare := func(n int) bool {
			r := int(math.Sqrt(float64(n)))
			return n == r*r
		}

		for n := 0; n < WordBitCap; n++ {
			if squares.MasksBit(n) {
				if !isSquare(n) {
					t.Errorf("\nexpected %d to be square\n", n)
				}
			} else if isSquare(n) {
				t.Errorf("\nexpected %d to not be square\n", n)
			}
		}
	}

	// tests := []struct {
	// 	a                              bm.Bitmask
	// 	expBin, expOct, expDec, expHex string
	// }{
	// 	{
	// 		a:      Zero(),
	// 		expBin: "0",
	// 		expOct: "0",
	// 		expDec: "0",
	// 		expHex: "0",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1)),
	// 		expBin: "1",
	// 		expOct: "1",
	// 		expDec: "1",
	// 		expHex: "1",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(2)),
	// 		expBin: "10",
	// 		expOct: "2",
	// 		expDec: "2",
	// 		expHex: "2",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(2)),
	// 		expBin: "11",
	// 		expOct: "3",
	// 		expDec: "3",
	// 		expHex: "3",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(4)),
	// 		expBin: "100",
	// 		expOct: "4",
	// 		expDec: "4",
	// 		expHex: "4",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(4)),
	// 		expBin: "101",
	// 		expOct: "5",
	// 		expDec: "5",
	// 		expHex: "5",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(2), UMask(4)),
	// 		expBin: "110",
	// 		expOct: "6",
	// 		expDec: "6",
	// 		expHex: "6",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(2), UMask(4)),
	// 		expBin: "111",
	// 		expOct: "7",
	// 		expDec: "7",
	// 		expHex: "7",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(8)),
	// 		expBin: "1000",
	// 		expOct: "10",
	// 		expDec: "8",
	// 		expHex: "8",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(8)),
	// 		expBin: "1001",
	// 		expOct: "11",
	// 		expDec: "9",
	// 		expHex: "9",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(2), UMask(8)),
	// 		expBin: "1010",
	// 		expOct: "12",
	// 		expDec: "10",
	// 		expHex: "a",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(2), UMask(8)),
	// 		expBin: "1011",
	// 		expOct: "13",
	// 		expDec: "11",
	// 		expHex: "b",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(4), UMask(8)),
	// 		expBin: "1100",
	// 		expOct: "14",
	// 		expDec: "12",
	// 		expHex: "c",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(4), UMask(8)),
	// 		expBin: "1101",
	// 		expOct: "15",
	// 		expDec: "13",
	// 		expHex: "d",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(2), UMask(4), UMask(8)),
	// 		expBin: "1110",
	// 		expOct: "16",
	// 		expDec: "14",
	// 		expHex: "e",
	// 	},
	// 	{
	// 		a:      Zero().Set(UMask(1), UMask(2), UMask(4), UMask(8)),
	// 		expBin: "1111",
	// 		expOct: "17",
	// 		expDec: "15",
	// 		expHex: "f",
	// 	},
	// }

	// for _, test := range tests {
	// 	if rec := test.a.Bin(); test.expBin != rec {
	// 		t.Errorf("\n   given %v\nexpected %s\nreceived %s\n", test.a, test.expBin, rec)
	// 	}

	// 	if rec := test.a.Oct(); test.expOct != rec {
	// 		t.Errorf("\n   given %v\nexpected %s\nreceived %s\n", test.a, test.expOct, rec)
	// 	}

	// 	if rec := test.a.Dec(); test.expDec != rec {
	// 		t.Errorf("\n   given %v\nexpected %s\nreceived %s\n", test.a, test.expDec, rec)
	// 	} else if exp := test.a.String(); exp != rec {
	// 		t.Errorf("\n   given %v\nexpected %s (UMask.String()) to equal %s (UMask.Dec())\n", test.a, exp, rec)
	// 	}

	// 	if rec := test.a.Hex(); test.expHex != rec {
	// 		t.Errorf("\n   given %v\nexpected %s\nreceived %s\n", test.a, test.expHex, rec)
	// 	}
	// }
}
