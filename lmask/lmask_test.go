package lmask

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

// -------------------------------------------------------------------------
// TODO: Finish testing
// -------------------------------------------------------------------------

func TestFromToBigInt(t *testing.T) {
	type testCase struct {
		exp *big.Int
	}

	tcs := []testCase{
		{exp: big.NewInt(0)},
		{exp: big.NewInt(1)},
		{exp: big.NewInt(WordMax >> 1)},
		{exp: big.NewInt(0).Lsh(big.NewInt(1), 2*WordBitCap)},
	}

	for _, tc := range tcs {
		if rec := FromBigInt(tc.exp).BigInt(); tc.exp.Cmp(rec) != 0 {
			t.Errorf("\nexpected %v\nreceived %v\n", tc, rec)
		}
	}
}

func TestFromBits(t *testing.T) {
	type testCase struct {
		bitCap int
		bits   []int
		exp    *LMask
	}

	tcs := []testCase{
		{
			bitCap: 4 * WordBitCap,
			bits:   []int{0, WordBitCap - 1, WordBitCap, 2*WordBitCap - 1, 2 * WordBitCap, 3*WordBitCap - 1, 3 * WordBitCap, 4*WordBitCap - 1},
			exp:    &LMask{bitCap: 4 * WordBitCap, words: []uint{1 | 1<<(WordBitCap-1), 1 | 1<<(WordBitCap-1), 1 | 1<<(WordBitCap-1), 1 | 1<<(WordBitCap-1)}},
		},
	}

	for _, tc := range tcs {
		if rec := FromBits(tc.bitCap, tc.bits...); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %q\nreceived %q\n", tc.exp, rec)
		}
	}
}

func TestFromWords(t *testing.T) {
	type testCase struct {
		words []uint
		exp   *LMask
	}

	tcs := []testCase{
		{
			words: nil,
			exp:   &LMask{bitCap: 0, words: []uint{}},
		},
		{
			words: []uint{uint(1<<WordBitCap - 1), uint(1<<WordBitCap - 1), uint(1<<WordBitCap - 1), uint(1<<WordBitCap - 1)},
			exp:   &LMask{bitCap: 4 * WordBitCap, words: []uint{uint(1<<WordBitCap - 1), uint(1<<WordBitCap - 1), uint(1<<WordBitCap - 1), uint(1<<WordBitCap - 1)}},
		},
	}

	for _, tc := range tcs {
		if rec := FromWords(tc.words...); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestMax(t *testing.T) {
	type testCase struct {
		bitCap int
		exp    *LMask
	}

	tcs := []testCase{
		{
			bitCap: 0,
			exp:    &LMask{bitCap: 0, words: []uint{}},
		},
		{
			bitCap: 1,
			exp:    &LMask{bitCap: 1, words: []uint{1}},
		},
		{
			bitCap: WordBitCap - 1,
			exp:    &LMask{bitCap: WordBitCap - 1, words: []uint{WordMax >> 1}},
		},
		{
			bitCap: WordBitCap,
			exp:    &LMask{bitCap: WordBitCap, words: []uint{WordMax}},
		},
		{
			bitCap: WordBitCap + 1,
			exp:    &LMask{bitCap: WordBitCap + 1, words: []uint{WordMax, 1}},
		},
		{
			bitCap: 2*WordBitCap - 1,
			exp:    &LMask{bitCap: 2*WordBitCap - 1, words: []uint{WordMax, WordMax >> 1}},
		},
		{
			bitCap: 2 * WordBitCap,
			exp:    &LMask{bitCap: 2 * WordBitCap, words: []uint{WordMax, WordMax}},
		},
	}

	for _, tc := range tcs {
		if a := Max(tc.bitCap); !a.Masks(tc.exp) || !tc.exp.Masks(a) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, a)
		}
	}
}

func TestOne(t *testing.T) {
	type testCase struct {
		bitCap int
		exp    *LMask
	}

	tcs := []testCase{
		{
			bitCap: 1,
			exp:    &LMask{bitCap: 1, words: []uint{1}},
		},
		{
			bitCap: WordBitCap,
			exp:    &LMask{bitCap: WordBitCap, words: []uint{1}},
		},
		{
			bitCap: WordBitCap + 1,
			exp:    &LMask{bitCap: WordBitCap + 1, words: []uint{1, 0}},
		},
		{
			bitCap: 4 * WordBitCap,
			exp:    &LMask{bitCap: 4 * WordBitCap, words: []uint{1, 0, 0, 0}},
		},
	}

	for _, tc := range tcs {
		if rec := One(tc.bitCap); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestZero(t *testing.T) {
	type testCase struct {
		bitCap int
		exp    *LMask
	}

	tcs := []testCase{
		{
			bitCap: 1,
			exp:    &LMask{bitCap: 1, words: []uint{0}},
		},
		{
			bitCap: WordBitCap,
			exp:    &LMask{bitCap: WordBitCap, words: []uint{0}},
		},
		{
			bitCap: WordBitCap + 1,
			exp:    &LMask{bitCap: WordBitCap + 1, words: []uint{0, 0}},
		},
		{
			bitCap: 4 * WordBitCap,
			exp:    &LMask{bitCap: 4 * WordBitCap, words: []uint{0, 0, 0, 0}},
		},
	}

	for _, tc := range tcs {
		if rec := Zero(tc.bitCap); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestAnd(t *testing.T) {
	type testCase struct {
		a, b, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4*WordBitCap).SetBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			b:   Max(4*WordBitCap).ClrBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Max(4*WordBitCap).ClrBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			b:   Max(4 * WordBitCap),
			exp: Max(4*WordBitCap).ClrBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Copy().And(tc.b); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestBits(t *testing.T) {
	type testCase struct {
		bits []int
	}

	tcs := []testCase{
		{},
		{bits: []int{0}},
		{bits: []int{0, WordBitCap - 1}},
	}

	for _, tc := range tcs {
		rec := Zero(4 * WordBitCap).SetBits(tc.bits...).Bits()
		equal := true
		if len(tc.bits) != len(rec) {
			equal = false
		}

		for i := 0; i < len(tc.bits) && equal; i++ {
			equal = tc.bits[i] == rec[i]
		}

		if !equal {
			t.Errorf("\nexpected %d\nreceived %d\n", tc, rec)
		}
	}
}

func TestNot(t *testing.T) {
	type testCase struct {
		a, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4*WordBitCap).SetBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			exp: Max(4*WordBitCap).ClrBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Copy().Not(); !tc.exp.Equals(tc.exp) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestOr(t *testing.T) {
	type testCase struct {
		a, b, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4*WordBitCap).SetBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			b:   Max(4*WordBitCap).ClrBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			exp: Max(4 * WordBitCap),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Copy().Or(tc.b); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

func TestBitCap(t *testing.T) {
	type testCase struct {
		bitCap int
	}

	tcs := []testCase{
		{bitCap: 0},
		{bitCap: 1},
		{bitCap: WordBitCap},
		{bitCap: 2 * WordBitCap},
	}

	for _, exp := range tcs {
		if rec := FromBits(exp.bitCap).BitCap(); exp.bitCap != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
		}
	}
}

func TestClr(t *testing.T) {
	type testCase struct {
		a, b, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4 * WordBitCap),
			b:   Zero(4 * WordBitCap),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Zero(4 * WordBitCap),
			b:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			b:   Zero(4 * WordBitCap),
			exp: Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			b:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			b:   Max(4 * WordBitCap),
			exp: Zero(4 * WordBitCap),
		},
	}

	for _, tc := range tcs {
		if tc.a.Clr(tc.b); !tc.exp.Equals(tc.a) {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, tc.a)
		}
	}
}

func TestClrAll(t *testing.T) {
	type testCase struct {
		a, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4 * WordBitCap),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   FromBits(4*WordBitCap, 0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Max(4 * WordBitCap),
			exp: Zero(4 * WordBitCap),
		},
	}

	for _, tc := range tcs {
		if tc.a.ClrAll(); !tc.exp.Equals(tc.a) {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, tc.a)
		}
	}
}

func TestClrBit(t *testing.T) {
	type testCase struct {
		a, exp *LMask
		bit    int
	}

	tcs := []testCase{
		{
			a:   Zero(4 * WordBitCap).SetBits(),
			exp: Zero(4 * WordBitCap).SetBits(),
			bit: 0,
		},
		{
			a:   Zero(4 * WordBitCap).SetBits(0),
			exp: Zero(4 * WordBitCap).SetBits(),
			bit: 0,
		},
		{
			a:   Zero(4 * WordBitCap).SetBits(1),
			exp: Zero(4 * WordBitCap).SetBits(),
			bit: 1,
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 1),
			exp: Zero(4 * WordBitCap).SetBits(0),
			bit: 1,
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 1),
			exp: Zero(4 * WordBitCap).SetBits(1),
			bit: 0,
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap).SetBits(0),
			bit: 4*WordBitCap - 1,
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap).SetBits(4*WordBitCap - 1),
			bit: 0,
		},
	}

	for _, tc := range tcs {
		if tc.a.ClrBit(tc.bit); !tc.exp.Equals(tc.a) {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, tc.a)
		}
	}
}

func TestClrBits(t *testing.T) {
	type testCase struct {
		a, b, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4 * WordBitCap),
			b:   Zero(4 * WordBitCap),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Zero(4 * WordBitCap),
			b:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			b:   Zero(4 * WordBitCap),
			exp: Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			b:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, 4*WordBitCap-1),
			b:   Max(4 * WordBitCap),
			exp: Zero(4 * WordBitCap),
		},
	}

	for _, tc := range tcs {
		if tc.a.ClrBits(tc.b.Bits()...); !tc.exp.Equals(tc.a) {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, tc.a)
		}
	}
}

func TestJSON(t *testing.T) {
	type testCase struct {
		expLMask *LMask
		expJSON  string
	}

	tcs := []testCase{
		{
			expLMask: nil,
			expJSON:  "null",
		},
		{
			expLMask: Zero(4 * WordBitCap),
			expJSON:  "0",
		},
		{
			expLMask: One(4 * WordBitCap),
			expJSON:  "1",
		},
		{
			expLMask: FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			expJSON:  big.NewInt(0).SetBits([]big.Word{1 | 1<<(WordBitCap-1), 1 | 1<<(WordBitCap-1), 1 | 1<<(WordBitCap-1), 1 | 1<<(WordBitCap-1)}).Text(10),
		},
		{
			expLMask: Max(4 * WordBitCap),
			expJSON:  big.NewInt(0).SetBits([]big.Word{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1}).Text(10),
		},
		{
			expLMask: Max(4*WordBitCap - 1),
			expJSON:  big.NewInt(0).SetBits([]big.Word{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<(WordBitCap-1) - 1}).Text(10),
		},
	}

	for _, tc := range tcs {
		recJSON := tc.expLMask.JSON()
		if tc.expJSON != recJSON {
			t.Errorf("\nexpected %q\nreceived %q\n", tc.expJSON, recJSON)
			continue
		}

		recLMask, err := FromJSON(recJSON)
		if err != nil {
			t.Error(err)
			continue
		}

		if tc.expLMask == nil {
			if recLMask.bitCap != 0 || len(recLMask.words) != 0 {
				t.Errorf("\nexpected %v\nreceived %v\n", Zero(0), recLMask)
			}

			continue
		}

		recLMask.SetBitCap(tc.expLMask.BitCap())
		if !tc.expLMask.Equals(recLMask) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.expLMask, recLMask)
		}
	}
}

func TestBitLen(t *testing.T) {
	type testCase struct {
		a   *LMask
		exp int
	}

	tcs := []testCase{
		{
			a:   FromBits(0),
			exp: 0,
		},
		{
			a:   FromBits(1),
			exp: 0,
		},
		{
			a:   FromBits(1, 0),
			exp: 1,
		},
		{
			a:   FromBits(2),
			exp: 0,
		},
		{
			a:   FromBits(2, 0),
			exp: 1,
		},
		{
			a:   FromBits(2, 1),
			exp: 2,
		},
		{
			a:   FromBits(2, 0, 1),
			exp: 2,
		},
		{
			a:   FromBits(WordBitCap),
			exp: 0,
		},
		{
			a:   FromBits(WordBitCap, 0),
			exp: 1,
		},
		{
			a:   FromBits(WordBitCap, WordBitCap-1),
			exp: WordBitCap,
		},
		{
			a:   FromBits(WordBitCap, 0, WordBitCap-1),
			exp: WordBitCap,
		},
		{
			a:   FromBits(2 * WordBitCap),
			exp: 0,
		},
		{
			a:   FromBits(2*WordBitCap, 0),
			exp: 1,
		},
		{
			a:   FromBits(2*WordBitCap, 2*WordBitCap-1),
			exp: 2 * WordBitCap,
		},
		{
			a:   FromBits(2*WordBitCap, 0, 2*WordBitCap-1),
			exp: 2 * WordBitCap,
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.BitLen(); tc.exp != rec {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, rec)
		}
	}
}

func TestLRSh(t *testing.T) {
	type testCase struct {
		a, expLeft, expRight *LMask
		left, right          int
	}

	tcs := []testCase{
		// Shifts less than WordBitCap
		{
			a:        FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			left:     0,
			right:    0,
			expLeft:  FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			expRight: FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
		},
		{
			a:        FromBits(WordBitCap, 0, WordBitCap-1),
			left:     1,
			right:    1,
			expLeft:  FromBits(WordBitCap, 1),
			expRight: FromBits(WordBitCap, WordBitCap-2),
		},
		{
			a:        FromBits(2*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1),
			left:     1,
			right:    1,
			expLeft:  FromBits(2*WordBitCap, 1, WordBitCap, WordBitCap+1),
			expRight: FromBits(2*WordBitCap, WordBitCap-2, WordBitCap-1, 2*WordBitCap-2),
		},
		{
			a:        FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			left:     1,
			right:    1,
			expLeft:  FromBits(4*WordBitCap, 1, WordBitCap, WordBitCap+1, 2*WordBitCap, 2*WordBitCap+1, 3*WordBitCap, 3*WordBitCap+1),
			expRight: FromBits(4*WordBitCap, WordBitCap-2, WordBitCap-1, 2*WordBitCap-2, 2*WordBitCap-1, 3*WordBitCap-2, 3*WordBitCap-1, 4*WordBitCap-2),
		},
		{
			a:        FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			left:     WordBitCap - 1,
			right:    WordBitCap - 1,
			expLeft:  FromBits(4*WordBitCap, WordBitCap-1, 2*WordBitCap-2, 2*WordBitCap-1, 3*WordBitCap-2, 3*WordBitCap-1, 4*WordBitCap-2, 4*WordBitCap-1),
			expRight: FromBits(4*WordBitCap, 0, 1, WordBitCap, WordBitCap+1, 2*WordBitCap, 2*WordBitCap+1, 3*WordBitCap),
		},

		// Shifts larger than or equal to WordBitCap
		{
			a:        FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			left:     WordBitCap,
			right:    WordBitCap,
			expLeft:  FromBits(4*WordBitCap, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			expRight: FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1),
		},
		{
			a:        FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			left:     2 * WordBitCap,
			right:    2 * WordBitCap,
			expLeft:  FromBits(4*WordBitCap, 2*WordBitCap, 3*WordBitCap-1, 3*WordBitCap, 4*WordBitCap-1),
			expRight: FromBits(4*WordBitCap, 0, WordBitCap-1, WordBitCap, 2*WordBitCap-1),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Copy().LSh(tc.left); !tc.expLeft.Equals(rec) {
			t.Errorf("\nexpected %s\nreceived %s\n", tc.expLeft.Fmt(2), rec.Fmt(2))
		}

		if rec := tc.a.Copy().RSh(tc.right); !tc.expRight.Equals(rec) {
			t.Errorf("\nexpected %s\nreceived %s\n", tc.expRight.Fmt(2), rec.Fmt(2))
		}
	}
}

func TestNextPrevBit(t *testing.T) {
	type testCase struct {
		a *LMask
	}

	for bitCap := 0; bitCap <= 4*WordBitCap; bitCap += 8 {
		tcs := append(make([]testCase, 0, 8), testCase{a: Zero(bitCap)})
		if 0 < bitCap {
			tcs = append(
				tcs,
				testCase{a: One(bitCap)},
				testCase{a: Zero(bitCap).SetBits(0, bitCap-1)}, // End bits
				testCase{a: Max(bitCap).ClrBits(0, bitCap-1)},  // Middle bits
				testCase{a: Max(bitCap)},
			)
		}

		if 8 <= bitCap {
			tcs = append(
				tcs,
				testCase{a: Zero(bitCap).SetBits(0, 2, 4, 6, bitCap-8, bitCap-6, bitCap-4, bitCap-2)}, // Some even bits
				testCase{a: Zero(bitCap).SetBits(1, 3, 5, 7, bitCap-7, bitCap-5, bitCap-3, bitCap-1)}, // Some odd bits
			)
		}

		for _, tc := range tcs {
			bitCap := tc.a.BitCap()
			for i := 0; i < bitCap; i++ {
				if tc.a.MasksBit(i) {
					if rec := tc.a.NextBit(i - 1); i != rec {
						t.Errorf("\nexpected next bit to be %d\nreceived next bit %d\n", i, rec)
					}

					if rec := tc.a.PrevBit(i + 1); i != rec {
						t.Errorf("\nexpected previous bit to be %d\nreceived next bit %d\n", i, rec)
					}
				} else {
					if rec := tc.a.NextBit(i - 1); i == rec {
						t.Errorf("\nexpected next bit to NOT be %d\nreceived next bit %d\n", i, rec)
					}

					if rec := tc.a.PrevBit(i + 1); i == rec {
						t.Errorf("\nexpected previous bit to NOT be %d\nreceived next bit %d\n", i, rec)
					}
				}
			}

			{
				bitCap := tc.a.BitCap()
				exp, rec := tc.a.Copy(), Zero(bitCap)
				for i := tc.a.NextBit(-1); i < bitCap; i = tc.a.NextBit(i) {
					rec.SetBit(i)
				}

				if !exp.Equals(rec) {
					t.Errorf("\nexpected %d\nreceived %d\n", exp, rec)
				}
			}

			{
				bitCap := tc.a.BitCap()
				exp, rec := Zero(bitCap), tc.a.Copy()
				for i := tc.a.PrevBit(bitCap); -1 < i; i = tc.a.PrevBit(i) {
					rec.ClrBit(i)
				}

				if !exp.Equals(rec) {
					t.Errorf("\nexpected %v\nreceived %v\n", exp, rec)
				}
			}
		}
	}
}

func TestSetAll(t *testing.T) {
	type testCase struct {
		a, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4 * WordBitCap),
			exp: Max(4 * WordBitCap),
		},
		{
			a:   FromBits(4*WordBitCap, 0, 4*WordBitCap-1),
			exp: Max(4 * WordBitCap),
		},
		{
			a:   Max(4 * WordBitCap),
			exp: Max(4 * WordBitCap),
		},
	}

	for _, tc := range tcs {
		if tc.a.SetAll(); !tc.exp.Equals(tc.a) {
			t.Errorf("\nexpected %d\nreceived %d\n", tc.exp, tc.a)
		}
	}
}

func TestString(t *testing.T) {
	type testCase struct {
		a   *LMask
		exp string
	}

	tcs := []testCase{
		{
			a:   nil,
			exp: "<nil>",
		},
		{
			a:   Zero(4 * WordBitCap),
			exp: "0",
		},
		{
			a:   Max(4 * WordBitCap),
			exp: big.NewInt(0).SetBits([]big.Word{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1}).Text(10),
		},
		{
			a:   Max(4*WordBitCap - 1),
			exp: big.NewInt(0).SetBits([]big.Word{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<(WordBitCap-1) - 1}).Text(10),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.String(); tc.exp != rec {
			t.Errorf("\nexpected %q\nrecieved %q\n", tc.exp, rec)
		}
	}
}

func TestTrim(t *testing.T) {
	type testCase struct {
		a, exp *LMask
	}

	tcs := []testCase{
		{
			a:   &LMask{bitCap: 0, words: []uint{}},
			exp: &LMask{bitCap: 0, words: []uint{}},
		},
		{
			a:   &LMask{bitCap: 1, words: []uint{1<<WordBitCap - 1}},
			exp: &LMask{bitCap: 1, words: []uint{1}},
		},
		{
			a:   &LMask{bitCap: WordBitCap - 1, words: []uint{1<<WordBitCap - 1}},
			exp: &LMask{bitCap: WordBitCap - 1, words: []uint{1<<(WordBitCap-1) - 1}},
		},
		{
			a:   &LMask{bitCap: 3*WordBitCap + 1, words: []uint{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1}},
			exp: &LMask{bitCap: 3*WordBitCap + 1, words: []uint{1<<WordBitCap - 1, 1<<WordBitCap - 1, 1<<WordBitCap - 1, 1}},
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Copy().trim(); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nrecieved %v\n", tc.exp, rec)
		}
	}
}

func TestXOr(t *testing.T) {
	type testCase struct {
		a, b, exp *LMask
	}

	tcs := []testCase{
		{
			a:   Zero(4*WordBitCap).SetBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			b:   Max(4*WordBitCap).ClrBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			exp: Max(4 * WordBitCap),
		},
		{
			a:   Zero(4*WordBitCap).SetBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			b:   Zero(4*WordBitCap).SetBits(0, WordBitCap-1, WordBitCap, 2*WordBitCap-1, 2*WordBitCap, 4*WordBitCap-1),
			exp: Zero(4 * WordBitCap),
		},
	}

	for _, tc := range tcs {
		if rec := tc.a.Copy().XOr(tc.b); !tc.exp.Equals(rec) {
			t.Errorf("\nexpected %v\nreceived %v\n", tc.exp, rec)
		}
	}
}

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

	var n0, n1, dn int = 1, 4 * WordBitCap, 8
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
			divisors *LMask = Zero(n + 1)
			a, b     int    = 1, n
			d        int    = n // a*b
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
				divisors.SetBits(a, b)
				a++
				b--
				d += b - a + 1
			}
		}

		for m := 0; m <= n; m++ {
			if gcd(m, n) == m {
				if !divisors.MasksBit(m) {
					t.Errorf("\nexpected %d to be set\n", m)
				}
			} else if divisors.MasksBit(m) {
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
	// next largest set bit a(n-2) and set the result as a(n). The bit
	// capacity is updated as needed.
	// ---------------------------------------------------------------------

	var (
		maxCount int    = 50
		fibs     *LMask = Zero(3).SetBits(1, 2)
	)

	for fibs.Count() < maxCount {
		var b int = fibs.PrevBit(fibs.BitLen())
		b += fibs.PrevBit(b)
		if fibs.BitCap() <= b {
			fibs.SetBitCap(b + 1)
		}

		fibs.SetBit(b)
	}

	for a0, a1 := 0, 1; a1 < maxCount; a0, a1 = a1, a0+a1 {
		if !fibs.MasksBit(a1) {
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
		bitCap     int    = WordBitCap << 2
		sqrtBitCap int    = int(math.Sqrt(float64(bitCap)))
		primes     *LMask = Max(bitCap).ClrBits(0, 1)
	)

	for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
		if primes.MasksBit(p) {
			for k := p * p; k < bitCap; k += p {
				primes.ClrBit(k)
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

	var (
		bitCap  int    = WordBitCap << 2
		squares *LMask = One(bitCap) // S(0) = 2*0 + 1
	)

	t.Logf("S(%d) = %d\n", squares.Count()-1, squares.PrevBit(bitCap)) // Prints S(0)
	for i := 1; ; i += 2 {
		var s int = squares.PrevBit(bitCap) + i // S(m) = S(m-1) + (2*m+1)
		if bitCap <= s {
			break
		}

		t.Logf("S(%d) = %d\n", squares.Count(), s) // Prints S(m), which has not been set in the bitmask yet
		squares.SetBit(s)
	}

	// isSquare determines if a number is a square.
	var isSquare = func(n int) bool {
		var r int = int(math.Sqrt(float64(n)))
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

// -------------------------------------------------------------------------
// Benchmarks
// -------------------------------------------------------------------------

func BenchmarkPrimes(b *testing.B) {
	const n0, n1, dn = 0, 4 * WordBitCap, 8
	for n := n0; n <= n1; n += dn {
		benchmarkPrimes(b, n)
	}

	for n := n0; n <= n1; n += dn {
		benchmarkPrimesNextBit(b, n)
	}
}

func benchmarkPrimes(b *testing.B, bitCap int) bool {
	var f = func(b *testing.B) {
		var (
			primes     *LMask
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = Max(bitCap).ClrBits(0, 1)
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
	var f = func(b *testing.B) {
		var (
			primes     *LMask
			sqrtBitCap int
		)

		for i := 0; i < b.N; i++ {
			primes = Max(bitCap).ClrBits(0, 1)
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
