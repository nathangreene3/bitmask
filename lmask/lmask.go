package lmask

import (
	"math/big"
	"math/bits"
)

const (
	// WordBitCap is the maximum number of bits in a word.
	WordBitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// WordMax is the maximum word.
	WordMax = 1<<WordBitCap - 1

	// errUneqBitCaps ...
	errUneqBitCaps = "unequal bit capacities"
)

// LMask is an arbitrarily sized bitmask.
type LMask struct {
	bitCap int
	words  []uint
}

// --------------------------------------------------------------------
// Constructors
// --------------------------------------------------------------------

// FromBigInt returns a bitmask from a given big integer. The bit
// capacity will be a multiple of the word bit capacity.
func FromBigInt(n *big.Int) *LMask {
	bigWords := n.Bits()
	words := make([]uint, 0, len(bigWords))
	for i := 0; i < len(bigWords); i++ {
		words = append(words, uint(bigWords[i]))
	}

	return &LMask{bitCap: len(words) * WordBitCap, words: words}
}

// FromBits returns a bitmask of a given bit capacity with and
// specified bits set.
func FromBits(bitCap int, bits ...int) *LMask {
	return Zero(bitCap).SetBits(bits...)
}

// FromJSON returns a bitmask decoded from a json-encoded string. The
// bit capacity will be a multiple of the word bit capacity.
func FromJSON(s string) (*LMask, error) {
	var a LMask
	if err := a.UnmarshalText([]byte(s)); err != nil {
		return nil, err
	}

	return &a, nil
}

// FromWords returns a bitmask set with a given list of uints. The bit
// capacity will be the number of words times the word bit capacity.
func FromWords(words ...uint) *LMask {
	return &LMask{bitCap: len(words) * WordBitCap, words: append(make([]uint, 0, len(words)), words...)}
}

// Max returns a bitmask with all bits set.
func Max(bitCap int) *LMask {
	return Zero(bitCap).Not()
}

// One returns a bitmask with bit zero set.
func One(bitCap int) *LMask {
	return Zero(bitCap).SetBit(0)
}

// Zero returns a bitmask with no bits set.
func Zero(bitCap int) *LMask {
	n := bitCap / WordBitCap
	if n*WordBitCap != bitCap {
		n++
	}

	return &LMask{bitCap: bitCap, words: make([]uint, n)}
}

// --------------------------------------------------------------------
// Logic functionality
// --------------------------------------------------------------------

// And sets each bit in a if the bit in b is also set. Otherwise, the
// bit in a is unset.
func (a *LMask) And(b *LMask) *LMask {
	if a != b {
		if a.bitCap != b.bitCap {
			panic(errUneqBitCaps)
		}

		for i := 0; i < len(a.words); i++ {
			a.words[i] &= b.words[i]
		}
	}

	return a
}

// AndNot sets each bit in a if the bit in a is set and the bit in b
// is not set. Otherwise, the bit in a is unset.
func (a *LMask) AndNot(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errUneqBitCaps)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] &^= b.words[i]
	}

	return a
}

// NAnd sets each bit in a if the bit is not set in both a and b.
// Otherwise, the bit in a is unset.
func (a *LMask) NAnd(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errUneqBitCaps)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] = ^(a.words[i] & b.words[i])
	}

	return a
}

// NOr sets each bit in a if the bit in a and b is unset. Otherwise,
// the bit is unset.
func (a *LMask) NOr(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errUneqBitCaps)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] = ^(a.words[i] | b.words[i])
	}

	return a
}

// Not inverts each bit in a.
func (a *LMask) Not() *LMask {
	for i := 0; i < len(a.words); i++ {
		a.words[i] = ^a.words[i]
	}

	return a.trim()
}

// Or sets each bit in a if either bit in a or b is set. Otherwise,
// the bit in a is unset.
func (a *LMask) Or(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errUneqBitCaps)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] |= b.words[i]
	}

	return a
}

// XNOr sets each bit in a if either both bits in a and b are set or
// unset. Otherwise, the bit in a is unset.
func (a *LMask) XNOr(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errUneqBitCaps)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] = ^(a.words[i] ^ b.words[i])
	}

	return a
}

// XOr sets each bit in a if exactly one bit in a and b is set.
// Otherwise, the bit in a is unset.
func (a *LMask) XOr(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errUneqBitCaps)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] ^= b.words[i]
	}

	return a
}

// --------------------------------------------------------------------
// Additional functionality
// --------------------------------------------------------------------

// BigInt returns an equivalent big integer.
func (a *LMask) BigInt() *big.Int {
	if a == nil {
		return nil
	}

	bigWords := make([]big.Word, 0, len(a.words))
	for i := 0; i < len(a.words); i++ {
		bigWords = append(bigWords, big.Word(a.words[i]))
	}

	return big.NewInt(0).SetBits(bigWords)
}

// BitCap returns the bit capacity.
func (a *LMask) BitCap() int {
	return a.bitCap
}

// BitLen returns the minimum number of bits required to represent a
// bitmask exactly.
func (a *LMask) BitLen() int {
	for i := len(a.words) - 1; 0 <= i; i-- {
		if 0 < a.words[i] {
			return bits.Len(a.words[i]) + i*WordBitCap
		}
	}

	return 0
}

// Bits returns the bits that are set in a bitmask.
func (a *LMask) Bits() []int {
	bits := make([]int, 0, a.Count())
	for bit := a.NextBit(-1); bit < a.bitCap; bit = a.NextBit(bit) {
		bits = append(bits, bit)
	}

	return bits
}

// Clr unsets each bit in a that is set in b.
func (a *LMask) Clr(b *LMask) *LMask {
	return a.AndNot(b)
}

// ClrBit unsets a bit.
func (a *LMask) ClrBit(bit int) *LMask {
	if 0 < a.bitCap {
		k := bit / WordBitCap
		a.words[k] &^= 1 << (bit - k*WordBitCap)
	}

	return a
}

// ClrBits unsets several bits.
func (a *LMask) ClrBits(bits ...int) *LMask {
	if 0 < a.bitCap {
		for i := 0; i < len(bits); i++ {
			k := bits[i] / WordBitCap
			a.words[k] &^= 1 << (bits[i] - k*WordBitCap)
		}
	}

	return a
}

// Copy returns a copy of a bitmask.
func (a *LMask) Copy() *LMask {
	return &LMask{bitCap: a.bitCap, words: append(make([]uint, 0, len(a.words)), a.words...)}
}

// Count returns the number of bits set.
func (a *LMask) Count() int {
	var c int
	for i := 0; i < len(a.words); i++ {
		c += bits.OnesCount(a.words[i])
	}

	return c
}

// Equal determines if two bitmasks are equal. Equality is defined as
// having the same bit capacity and the same bits set.
func (a *LMask) Equals(b *LMask) bool {
	if a == b {
		return true
	}

	if a.bitCap != b.bitCap {
		return false
	}

	for i := 0; i < len(a.words); i++ {
		if a.words[i] != b.words[i] {
			return false
		}
	}

	return true
}

// Fmt returns a string formatted as the integer represented by a
// bitmask in a given base. This supports bases on range [2, 62].
func (a *LMask) Fmt(base int) string {
	return a.BigInt().Text(base)
}

// JSON returns a json-encoded string representing a bitmask.
func (a *LMask) JSON() string {
	b, _ := a.MarshalText()
	return string(b)
}

// LSh shifts all set bits by a given amount. That is, each set bit i
// will be unset and bit i+bits will be set.
func (a *LMask) LSh(bits int) *LMask {
	if 0 < a.bitCap {
		k := bits / WordBitCap
		if 0 < k {
			copy(a.words[:len(a.words)-k], a.words[k:])
			for i := 0; i < k; i++ {
				a.words[i] = 0
			}
		}

		if r := bits - k*WordBitCap; 0 < r {
			d := WordBitCap - r
			h0 := a.words[k] >> d
			a.words[k] <<= r
			for k++; k < len(a.words); k++ {
				h1 := a.words[k] >> d
				a.words[k] = (a.words[k] << r) | h0
				h0 = h1
			}
		}
	}

	return a.trim()
}

// MarshalText returns text representing a bitmask.
func (a *LMask) MarshalText() ([]byte, error) {
	if a == nil {
		// TODO: big.Int.MarshalText returns <nil>. Determine if this
		// should do the same.
		return []byte("null"), nil
	}

	return []byte(a.Fmt(10)), nil
}

// MarshalJSON returns json-encoded bytes representing a bitmask.
func (a *LMask) MarshalJSON() ([]byte, error) {
	return a.MarshalText()
}

// Masks determines if a masks b.
func (a *LMask) Masks(b *LMask) bool {
	if a == b {
		return true
	}

	if b.bitCap < a.bitCap {
		return false
	}

	for i := 0; i < len(a.words); i++ {
		if a.words[i]&b.words[i] != b.words[i] {
			return false
		}
	}

	return true
}

// MasksBit determines if a bit is set in a.
func (a *LMask) MasksBit(bit int) bool {
	k := bit / WordBitCap
	c := uint(1 << (bit - k*WordBitCap))
	return a.words[k]&c == c
}

// NextBit returns the next set bit in a. If no set bit is next, then
// the bit capacity is returned.
func (a *LMask) NextBit(bit int) int {
	bit = clamp(bit+1, 0, a.bitCap)
	if k := bit / WordBitCap; k < len(a.words) {
		r := bit - k*WordBitCap
		if w := a.words[k] >> r << r; 0 < w {
			return bits.TrailingZeros(w) + k*WordBitCap
		}

		for k++; k < len(a.words); k++ {
			if 0 < a.words[k] {
				return min(bits.TrailingZeros(a.words[k])+k*WordBitCap, a.bitCap)
			}
		}
	}

	return a.bitCap
}

// PrevBit returns the previous set bit in a. If no set bit is next,
// then -1 is returned.
func (a *LMask) PrevBit(bit int) int {
	bit = clamp(bit, 0, a.bitCap)

	i := bit / WordBitCap
	if i < len(a.words) {
		r := (i+1)*WordBitCap - bit
		if w := a.words[i] << r >> r; 0 < w {
			return -bits.LeadingZeros(w) + (i+1)*WordBitCap - 1
		}
	}

	for i--; 0 <= i; i-- {
		if 0 < a.words[i] {
			return -bits.LeadingZeros(a.words[i]) + (i+1)*WordBitCap - 1
		}
	}

	return -1
}

// RSh shifts all set bits by a given amount. That is, each set bit i
// will be unset and bit i-bits will be set.
func (a *LMask) RSh(bits int) *LMask {
	if 0 < a.bitCap {
		k := bits / WordBitCap
		if 0 < k {
			copy(a.words[k:], a.words[:len(a.words)-k])
			for i := len(a.words) - k; i < len(a.words); i++ {
				a.words[i] = 0
			}
		}

		if r := bits - k*WordBitCap; 0 < r {
			d := WordBitCap - r
			h0 := a.words[len(a.words)-1] << d
			a.words[len(a.words)-1] >>= r
			for i := len(a.words) - 2; k <= i; i-- {
				h1 := a.words[i] << d
				a.words[i] = (a.words[i] >> r) | h0
				h0 = h1
			}
		}
	}

	return a.trim()
}

// Set sets the bits of b in a. Any bits already set in a will remain
// set.
func (a *LMask) Set(b *LMask) *LMask {
	return a.Or(b)
}

// SetBit sets a bit in a.
func (a *LMask) SetBit(bit int) *LMask {
	k := bit / WordBitCap
	a.words[k] |= 1 << (bit - k*WordBitCap)
	return a.trim()
}

// SetBitCap sets the bit capacity. If the bit capacity is decreasing,
// any higher-order bits will be lost. If the bit capacity is
// increasing, new bits will be unset.
func (a *LMask) SetBitCap(bitCap int) *LMask {
	if a.bitCap == bitCap {
		return a
	}

	a.bitCap = bitCap
	n := bitCap / WordBitCap
	if n*WordBitCap < bitCap {
		n++
	}

	if n < len(a.words) {
		a.words = a.words[:n]
		return a.trim()
	}

	if len(a.words) < n {
		words := make([]uint, n)
		copy(words[:len(a.words)], a.words)
		a.words = words
		return a
	}

	return a.trim()
}

// SetBits sets several bits.
func (a *LMask) SetBits(bits ...int) *LMask {
	for i := 0; i < len(bits); i++ {
		k := bits[i] / WordBitCap
		a.words[k] |= 1 << (bits[i] - k*WordBitCap)
	}

	return a.trim()
}

// String returns the base-10 integer representation of a bitmask.
func (a *LMask) String() string {
	return a.Fmt(10)
}

// UnmarshalJSON decodes json-encoded text into a bitmask. If the text
// is "null", no action is taken and no error is returned. Otherwise,
// the text is assumed to be the base-10 integer representation of a
// bitmask. The bit capacity will be a multiple of the word bit
// capacity.
func (a *LMask) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}

	return a.UnmarshalText(b)
}

// UnmarshalText decodes text into a bitmask. The text is assumed to be
// the base-10 integer representation of a bitmask. The bit capacity
// will be a multiple of the word bit capacity.
func (a *LMask) UnmarshalText(b []byte) error {
	n := big.NewInt(0)
	if err := n.UnmarshalJSON(b); err != nil {
		return err
	}

	bigWords := n.Bits()
	words := make([]uint, 0, len(bigWords))
	for len(words) < len(bigWords) {
		words = append(words, uint(bigWords[len(words)]))
	}

	a.bitCap = len(words) * WordBitCap
	a.words = words
	return nil
}

// Words returns a copy of the words in the bitmask.
func (a *LMask) Words() []uint {
	return append(make([]uint, 0, len(a.words)), a.words...)
}

// --------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------

// clamp returns a if n < a, b if b < n, or otherwise n.
func clamp(n, a, b int) int {
	switch {
	case n < a:
		return a
	case b < n:
		return b
	default:
		return n
	}
}

// min returns the minimum value.
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// trim unsets any leading bits greater than the bitmask's bit capacity.
func (a *LMask) trim() *LMask {
	if 0 < len(a.words) {
		if r := a.bitCap - a.bitCap/WordBitCap*WordBitCap; 0 < r {
			a.words[len(a.words)-1] &^= WordMax << r
		}
	}

	return a
}
