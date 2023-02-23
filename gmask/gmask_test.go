package gmask

import (
	"fmt"
	"math"
	"math/bits"
	"testing"

	"github.com/onsi/gomega"
)

// ---------------------------------------------------------------------
// 	Reporting functionality
// ---------------------------------------------------------------------

func TestBitCap(t *testing.T) {
	testBitCap[uint](t)
	testBitCap[uint64](t)
	testBitCap[uint32](t)
	testBitCap[uint16](t)
	testBitCap[uint8](t)
}

func testBitCap[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)

	switch any(Zero[T]()).(type) {
	case uint:
		g.Expect(BitCap[uint]()).To(gomega.Equal(bits.UintSize))
	case uint64:
		g.Expect(BitCap[uint64]()).To(gomega.Equal(64))
	case uint32:
		g.Expect(BitCap[uint32]()).To(gomega.Equal(32))
	case uint16:
		g.Expect(BitCap[uint16]()).To(gomega.Equal(16))
	case uint8:
		g.Expect(BitCap[uint8]()).To(gomega.Equal(8))
	}
}

// ---------------------------------------------------------------------
// 	Constructors
// ---------------------------------------------------------------------

func TestMax(t *testing.T) {
	testMax[uint](t)
	testMax[uint64](t)
	testMax[uint32](t)
	testMax[uint16](t)
	testMax[uint8](t)
}

func testMax[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)

	g.Expect(Max[T]()).To(gomega.Equal(Zero[T]() - 1))
	g.Expect(Max[T]() + 1).To(gomega.Equal(Zero[T]()))
}

func TestFromBits(t *testing.T) {
	testFromBits[uint](t)
	testFromBits[uint64](t)
	testFromBits[uint32](t)
	testFromBits[uint16](t)
	testFromBits[uint8](t)
}

func testFromBits[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected FromBits[%T](%d) = %d\n" +
		"received %d\n"

	var exp T = 1 | 1<<(BitCap[T]()-1)

	bits := []int{0, BitCap[T]() - 1}
	rec := FromBits[T](bits...)

	if exp != rec {
		t.Errorf(errFmt, exp, bits, exp, rec)
	}
}

// ---------------------------------------------------------------------
// 	Bitwise functionality
// ---------------------------------------------------------------------

func TestAnd(t *testing.T) {
	testAnd[uint](t)
	testAnd[uint64](t)
	testAnd[uint32](t)
	testAnd[uint16](t)
	testAnd[uint8](t)
}

func testAnd[T UInteger](t *testing.T) {
	var (
		g          = gomega.NewWithT(t)
		lastBitSet = FromBits[T](BitCap[T]() - 1)
		max        = Max[T]()
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, 0},

		{1, 0, 0},
		{0, 1, 0},
		{1, 1, 1},

		{lastBitSet, 0, 0},
		{0, lastBitSet, 0},
		{1, lastBitSet, 0},
		{lastBitSet, 1, 0},
		{lastBitSet, lastBitSet, lastBitSet},
		{lastBitSet, max, lastBitSet},
		{max, lastBitSet, lastBitSet},

		{max, 0, 0},
		{0, max, 0},
		{max, max, max},
	}

	for _, test := range tests {
		g.Expect(And(test.a, test.b)).To(gomega.Equal(test.exp))
	}
}

func TestAndNot(t *testing.T) {
	testAndNot[uint](t)
	testAndNot[uint64](t)
	testAndNot[uint32](t)
	testAndNot[uint16](t)
	testAndNot[uint8](t)
}

func testAndNot[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected AndNot[%T](%d, %d) to be %d\n" +
		"received %d\n"

	var (
		a  = T(1)
		b  = T(a) << (BitCap[T]() - 1)
		ab = a | b
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, 0},

		{a, 0, a},
		{0, a, 0},
		{a, a, 0},

		{b, 0, b},
		{0, b, 0},
		{b, b, 0},

		{ab, a, b},
		{a, ab, 0},
		{ab, b, a},
		{b, ab, 0},
		{ab, ab, 0},
	}

	for _, test := range tests {
		if rec := AndNot(test.a, test.b); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.b, test.exp, rec)
		}
	}
}

func TestNAnd(t *testing.T) {
	testNAnd[uint](t)
	testNAnd[uint64](t)
	testNAnd[uint32](t)
	testNAnd[uint16](t)
	testNAnd[uint8](t)
}

func testNAnd[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected NAnd[%T](%d, %d) to be %d\n" +
		"received %d\n"

	var (
		a   = T(1)
		na  = ^a
		b   = T(1) << (BitCap[T]() - 1)
		nb  = ^b
		ab  = a | b
		nab = ^ab
		max = Max[T]()
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, max},

		{a, 0, max},
		{0, a, max},
		{a, a, na},

		{b, 0, max},
		{0, b, max},
		{b, b, nb},

		{ab, a, na},
		{a, ab, na},
		{ab, b, nb},
		{b, ab, nb},
		{ab, ab, nab},
	}

	for _, test := range tests {
		if rec := NAnd(test.a, test.b); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.b, test.exp, rec)
		}
	}
}

func TestNOr(t *testing.T) {
	testNOr[uint](t)
	testNOr[uint64](t)
	testNOr[uint32](t)
	testNOr[uint16](t)
	testNOr[uint8](t)
}

func testNOr[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected NOr[%T](%d, %d) to be %d\n" +
		"received %d\n"

	var (
		a   = T(1)
		na  = ^a
		b   = T(1) << (BitCap[T]() - 1)
		nb  = ^b
		ab  = a | b
		nab = ^ab
		max = Max[T]()
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, max},

		{a, 0, na},
		{0, a, na},
		{a, a, na},

		{b, 0, nb},
		{0, b, nb},
		{b, b, nb},

		{ab, a, nab},
		{a, ab, nab},
		{ab, b, nab},
		{b, ab, nab},
		{ab, ab, nab},
	}

	for _, test := range tests {
		if rec := NOr(test.a, test.b); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.b, test.exp, rec)
		}
	}
}

func TestNot(t *testing.T) {
	testNot[uint](t)
	testNot[uint64](t)
	testNot[uint32](t)
	testNot[uint16](t)
	testNot[uint8](t)
}

func testNot[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected Not[%T](%d) to be %d\n" +
		"received %d\n"

	var (
		a   = T(1)
		na  = ^a
		b   = T(a) << (BitCap[T]() - 1)
		nb  = ^b
		ab  = a | b
		nab = ^ab
		max = Max[T]()
	)

	tests := []struct {
		a, exp T
	}{
		{0, max},
		{max, 0},

		{a, na},
		{na, a},

		{b, nb},
		{nb, b},

		{ab, nab},
		{nab, ab},
	}

	for _, test := range tests {
		if rec := Not(test.a); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.exp, rec)
		}
	}
}

func TestOr(t *testing.T) {
	testOr[uint](t)
	testOr[uint64](t)
	testOr[uint32](t)
	testOr[uint16](t)
	testOr[uint8](t)
}

func testOr[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected Or[%T](%d, %d) to be %d\n" +
		"received %d\n"

	var (
		a  = T(1)
		b  = T(a) << (BitCap[T]() - 1)
		ab = a | b
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, 0},

		{a, 0, a},
		{0, a, a},
		{a, a, a},

		{b, 0, b},
		{0, b, b},
		{b, b, b},

		{ab, a, ab},
		{a, ab, ab},
		{ab, b, ab},
		{b, ab, ab},
		{ab, ab, ab},
	}

	for _, test := range tests {
		if rec := Or(test.a, test.b); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.b, test.exp, rec)
		}
	}
}

func TestXNOr(t *testing.T) {
	testXNOr[uint](t)
	testXNOr[uint64](t)
	testXNOr[uint32](t)
	testXNOr[uint16](t)
	testXNOr[uint8](t)
}

func testXNOr[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected XNOr[%T](%d, %d) to be %d\n" +
		"received %d\n"

	var (
		a  = T(1)
		b  = T(a) << (BitCap[T]() - 1)
		ab = a | b
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, Max[T]()},

		{a, 0, ^a},
		{0, a, ^a},
		{a, a, Max[T]()},

		{b, 0, ^b},
		{0, b, ^b},
		{b, b, Max[T]()},

		{ab, a, ^b},
		{a, ab, ^b},
		{ab, b, ^a},
		{b, ab, ^a},
		{ab, ab, Max[T]()},
	}

	for _, test := range tests {
		if rec := XNOr(test.a, test.b); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.b, test.exp, rec)
		}
	}
}

func TestXOr(t *testing.T) {
	testXOr[uint](t)
	testXOr[uint64](t)
	testXOr[uint32](t)
	testXOr[uint16](t)
	testXOr[uint8](t)
}

func testXOr[T UInteger](t *testing.T) {
	const errFmt = "\n" +
		"expected XOr[%T](%d, %d) to be %d\n" +
		"received %d\n"

	var (
		a  = T(1)
		b  = T(a) << (BitCap[T]() - 1)
		ab = a | b
	)

	tests := []struct {
		a, b, exp T
	}{
		{0, 0, 0},

		{a, 0, a},
		{0, a, a},
		{a, a, 0},

		{b, 0, b},
		{0, b, b},
		{b, b, 0},

		{ab, a, b},
		{a, ab, b},
		{ab, b, a},
		{b, ab, a},
		{ab, ab, 0},
	}

	for _, test := range tests {
		if rec := XOr(test.a, test.b); test.exp != rec {
			t.Errorf(errFmt, test.exp, test.a, test.b, test.exp, rec)
		}
	}
}

// ---------------------------------------------------------------------
// 	Set functionality
// ---------------------------------------------------------------------

// TODO: Test the following.
// 	- ClrBits
// 	- Count
// 	- Masks
// 	- MasksBit
// 	- Set
// 	- SetBit
// 	- SetBits

func TestBits(t *testing.T) {
	testBits[uint](t)
	testBits[uint64](t)
	testBits[uint32](t)
	testBits[uint16](t)
	testBits[uint8](t)
}

func testBits[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)
	bits := []int{0, BitCap[T]() - 1}

	g.Expect(Bits(FromBits[T](bits...))).To(gomega.Equal(bits))
}

func TestClr(t *testing.T) {
	testClr[uint](t)
	testClr[uint64](t)
	testClr[uint32](t)
	testClr[uint16](t)
	testClr[uint8](t)
}

func testClr[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		a, b, exp T
	}{
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			b:   FromBits[T](1, BitCap[T]()-2, BitCap[T]()-1),
			exp: FromBits[T](0),
		},
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			b:   FromBits[T](0, BitCap[T]()-2, BitCap[T]()-1),
			exp: FromBits[T](1),
		},
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			b:   FromBits[T](0, 1, BitCap[T]()-1),
			exp: FromBits[T](BitCap[T]() - 2),
		},
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			b:   FromBits[T](0, 1, BitCap[T]()-2),
			exp: FromBits[T](BitCap[T]() - 1),
		},
	}

	for _, test := range tests {
		g.Expect(Clr(test.a, test.b)).To(gomega.Equal(test.exp))
	}
}

func TestClrBit(t *testing.T) {
	testClrBit[uint](t)
	testClrBit[uint64](t)
	testClrBit[uint32](t)
	testClrBit[uint16](t)
	testClrBit[uint8](t)
}

func testClrBit[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		a, exp T
		bit    int
	}{
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			bit: 0,
			exp: FromBits[T](1, BitCap[T]()-2, BitCap[T]()-1),
		},
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			bit: 1,
			exp: FromBits[T](0, BitCap[T]()-2, BitCap[T]()-1),
		},
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			bit: BitCap[T]() - 2,
			exp: FromBits[T](0, 1, BitCap[T]()-1),
		},
		{
			a:   FromBits[T](0, 1, BitCap[T]()-2, BitCap[T]()-1),
			bit: BitCap[T]() - 1,
			exp: FromBits[T](0, 1, BitCap[T]()-2),
		},
	}

	for _, test := range tests {
		g.Expect(ClrBit(test.a, test.bit)).To(gomega.Equal(test.exp))
	}
}

func TestNextPrevBit(t *testing.T) {
	testNextPrevBit[uint](t)
	testNextPrevBit[uint64](t)
	testNextPrevBit[uint32](t)
	testNextPrevBit[uint16](t)
	testNextPrevBit[uint8](t)
}

func testNextPrevBit[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []T{
		0,                                   // 0000 0000 ... 0000 0000
		1,                                   // 0000 0000 ... 0000 0001
		SetBits[T](0, 0, BitCap[T]()-1),     // 1000 0000 ... 0000 0001
		ClrBits(Max[T](), 0, BitCap[T]()-1), // 0111 1111 ... 1111 1110
		Max[T](),                            // 1111 1111 ... 1111 1111
		SetBits[T](0, 2, 4, 6, BitCap[T]()-8, BitCap[T]()-6, BitCap[T]()-4, BitCap[T]()-2), // 0101 0101 ... 0101 0101
		SetBits[T](1, 3, 5, 7, BitCap[T]()-7, BitCap[T]()-5, BitCap[T]()-3, BitCap[T]()-1), // 1010 1010 ... 1010 1010
	}

	for _, test := range tests {
		for bit := 0; bit < BitCap[T](); bit++ {
			if MasksBit(test, bit) {
				g.Expect(NextBit(test, bit-1)).To(gomega.Equal(bit))
				g.Expect(PrevBit(test, bit+1)).To(gomega.Equal(bit))
			} else {
				g.Expect(NextBit(test, bit-1)).ToNot(gomega.Equal(bit))
				g.Expect(PrevBit(test, bit+1)).ToNot(gomega.Equal(bit))
			}
		}

		{
			var rec T
			for bit := NextBit(test, -1); bit < BitCap[T](); bit = NextBit(test, bit) {
				rec = SetBit(rec, bit)
			}

			g.Expect(rec).To(gomega.Equal(test))
		}

		{
			var (
				exp T
				rec = test
			)

			for bit := PrevBit(test, BitCap[T]()); -1 < bit; bit = PrevBit(test, bit) {
				rec = ClrBit(rec, bit)
			}

			g.Expect(rec).To(gomega.Equal(exp))
		}
	}
}

// ---------------------------------------------------------------------
// 	Helper functionality
// ---------------------------------------------------------------------

func TestClamp(t *testing.T) {
	testClamp(t)
	testClampPanic(t)
}

func testClamp(t *testing.T) {
	g := gomega.NewWithT(t)
	tests := []struct {
		n, a, b, exp int
	}{
		{-1, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 0, 0, 0},

		{-1, 0, 1, 0},
		{0, 0, 1, 0},
		{1, 0, 1, 1},
		{2, 0, 1, 1},

		{-2, -1, 0, -1},
		{-1, -1, 0, -1},
		{0, -1, 0, 0},
		{1, -1, 0, 0},

		{-2, -1, 1, -1},
		{-1, -1, 1, -1},
		{0, -1, 1, 0},
		{1, -1, 1, 1},
		{2, -1, 1, 1},
	}

	for _, test := range tests {
		g.Expect(clamp(test.n, test.a, test.b)).To(gomega.Equal(test.exp))
	}
}

func testClampPanic(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		n, a, b int
	}{
		{0, 0, -1},
		{0, 1, 0},
	}

	for _, test := range tests {
		g.Expect(func() {
			_ = clamp(test.n, test.a, test.b)
		}).To(gomega.PanicWith(ErrInvalidRange))
	}
}

func BenchmarkBits(b *testing.B) {
	benchmarkBits(b, Max[uint]())
	benchmarkBits(b, Max[uint64]())
	benchmarkBits(b, Max[uint32]())
	benchmarkBits(b, Max[uint16]())
	benchmarkBits(b, Max[uint8]())
}

func benchmarkBits[T UInteger](b *testing.B, a T) bool {
	f := func(b *testing.B) {
		var bits []int
		for i := 0; i < b.N; i++ {
			bits = Bits(a)
		}

		_ = bits
	}

	return b.Run(fmt.Sprintf("Bits[%T](%d)", a, a), f)
}

func TestPrimes(t *testing.T) {
	testPrimes[uint](t)
	testPrimes[uint64](t)
	testPrimes[uint32](t)
	testPrimes[uint16](t)
	testPrimes[uint8](t)
}

func testPrimes[T UInteger](t *testing.T) {
	var (
		g          = gomega.NewWithT(t)
		bitCap     = BitCap[T]()
		sqrtBitCap = int(math.Sqrt(float64(bitCap)))
		primes     = ClrBits(Max[T](), 0, 1)
	)

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

	for p := NextBit(primes, -1); p <= sqrtBitCap; p = NextBit(primes, p) {
		if MasksBit(primes, p) {
			for k := p * p; k < bitCap; k += p {
				primes = ClrBit(primes, k)
			}
		}
	}

	for bit := 0; bit < bitCap; bit++ {
		g.Expect(MasksBit(primes, bit)).To(gomega.Equal(isPrime(bit)))
	}
}

func TestLights(t *testing.T) {
	testLights[uint](t)
	testLights[uint64](t)
	testLights[uint32](t)
	testLights[uint16](t)
	testLights[uint8](t)
}

func testLights[T UInteger](t *testing.T) {
	g := gomega.NewWithT(t)

	isSqr := func(n int) bool {
		return int(math.Pow(math.Sqrt(float64(n)), 2.0)) == n
	}

	lights := Zero[T]()

	for i := 1; i < BitCap[T](); i++ {
		for light := i; light < BitCap[T](); light += i {
			lights = ToggleBit(lights, light)
		}
	}

	for light := NextBit(lights, -1); light < BitCap[T](); light = NextBit(lights, light) {
		g.Expect(MasksBit(lights, light)).To(gomega.Equal(isSqr(light)))
	}
}
