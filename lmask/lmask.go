package lmask

import (
	"errors"
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

const (
	// WordBitCap is the maximum number of bits in a word.
	WordBitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// WordMax is the maximum word.
	WordMax = 1<<WordBitCap - 1
)

// errBitCap indicates bitmask operations can only be applied when the bit capacities are equal.
var errBitCap = errors.New("undefined on unequal bit capacities")

// LMask is an arbitrarily sized bitmask.
type LMask struct {
	bitCap int
	words  []uint
}

// ------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------

// FromBits ...
func FromBits(bitCap int, bits ...int) *LMask {
	return Zero(bitCap).SetBits(bits...)
}

// FromWords returns a bitmask set with a given list of uints. The bit capacity will be the number of words times the WordBitCap (a multiple of 32 or 64).
func FromWords(words ...uint) *LMask {
	return &LMask{bitCap: WordBitCap * len(words), words: append(make([]uint, 0, len(words)), words...)}
}

// Max returns a bitmask with all bits set.
func Max(bitCap int) *LMask {
	return Zero(bitCap).Not()
}

// One returns a bitmask with bit zero set. This is equivalent to the number one in binary.
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

// ------------------------------------------------------------------------------------
// Logic functionality
// ------------------------------------------------------------------------------------

// And sets each bit in a if the bit in b is also set. Otherwise, the bit in a is unset.
func (a *LMask) And(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] &= b.words[i]
	}

	return a
}

// AndNot sets each bit in a if the bit in a is set and the bit in b is not set. Otherwise, the bit in a is unset.
func (a *LMask) AndNot(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] &^= b.words[i]
	}

	return a
}

// NAnd sets each bit in a if the bit is not set in both a and b. Otherwise, the bit in a is unset.
func (a *LMask) NAnd(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] = ^(a.words[i] & b.words[i])
	}

	return a
}

// NOr ...
func (a *LMask) NOr(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
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

// Or sets each bit in a if either bit in a or b is set. Otherwise, the bit in a is unset.
func (a *LMask) Or(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] |= b.words[i]
	}

	return a
}

// XNOr sets each bit in a if either both bits in a and b are set or unset. Otherwise, the bit in a is unset.
func (a *LMask) XNOr(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] = ^(a.words[i] ^ b.words[i])
	}

	return a
}

// XOr sets each bit in a if exactly one bit in a and b is set. Otherwise, the bit in a is unset.
func (a *LMask) XOr(b *LMask) *LMask {
	if a.bitCap != b.bitCap {
		panic(errBitCap)
	}

	for i := 0; i < len(a.words); i++ {
		a.words[i] ^= b.words[i]
	}

	return a
}

// ------------------------------------------------------------------------------------
// Additional functionality
// ------------------------------------------------------------------------------------

// Bin returns the bits of the bitmask. The left-most bit is the leading bit.
func (a *LMask) Bin() string {
	var sb strings.Builder
	sb.Grow(a.bitCap)
	for i := len(a.words) - 1; 0 <= i; i-- {
		sb.WriteString(strconv.FormatUint(uint64(a.words[i]), 2))
	}

	return sb.String()
}

// BitCap returns the bit capacity.
func (a *LMask) BitCap() int {
	return a.bitCap
}

// Clr unsets each bit in a that is set in b.
func (a *LMask) Clr(b *LMask) *LMask {
	return a.AndNot(b)
}

// ClrBit unsets a bit.
func (a *LMask) ClrBit(bit int) *LMask {
	if a.bitCap != 0 {
		a.words[bit/WordBitCap] &^= 1 << (bit % WordBitCap)
	}

	return a
}

// ClrBits unsets several bits.
func (a *LMask) ClrBits(bits ...int) *LMask {
	if a.bitCap != 0 {
		for i := 0; i < len(bits); i++ {
			a.words[bits[i]/WordBitCap] &^= 1 << (bits[i] % WordBitCap)
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

// Equal determines if two bitmasks are equal. Equality is defined as having the same bit capacity and the same bits set.
func (a *LMask) Equals(b *LMask) bool {
	if a != b {
		if a.bitCap != b.bitCap {
			return false
		}

		for i := 0; i < len(a.words); i++ {
			if a.words[i] != b.words[i] {
				return false
			}
		}
	}

	return true
}

// Masks determines if a masks b.
func (a *LMask) Masks(b *LMask) bool {
	if a != b {
		if a.bitCap != b.bitCap {
			panic(errBitCap)
		}

		for i := 0; i < len(a.words); i++ {
			if a.words[i]&b.words[i] != b.words[i] {
				return false
			}
		}
	}

	return true
}

// MasksBit determines if a bit is set in a.
func (a *LMask) MasksBit(bit int) bool {
	c := uint(1 << (bit % WordBitCap))
	return a.words[bit/WordBitCap]&c == c
}

// NextBit returns the next set bit in a. If no set bit is next, then the bit capacity is returned.
func (a *LMask) NextBit(bit int) int {
	if bit++; bit < a.bitCap {
		i, r := bit/WordBitCap, bit%WordBitCap
		if w := a.words[i] >> r << r; w != 0 {
			return bits.TrailingZeros(w) + i*WordBitCap
		}

		for i++; i < len(a.words); i++ {
			if a.words[i] != 0 {
				return min(bits.TrailingZeros(a.words[i])+i*WordBitCap, a.bitCap)
			}
		}
	}

	return a.bitCap
}

// Set sets the bits of b in a.
func (a *LMask) Set(b *LMask) *LMask {
	return a.Or(b)
}

// SetBit sets a bit in a.
func (a *LMask) SetBit(bit int) *LMask {
	a.words[bit/WordBitCap] |= 1 << (bit % WordBitCap)
	return a.trim()
}

// SetBitCap sets the bit capacity. If the bit capacity is decreasing, any higher-order bits will be lost. If the bit capacity is increasing, new bits will be unset.
func (a *LMask) SetBitCap(bitCap int) *LMask {
	if a.bitCap == bitCap {
		return a
	}

	a.bitCap = bitCap
	n := bitCap / WordBitCap
	if n*WordBitCap < bitCap {
		n++
	}

	switch {
	case n < len(a.words):
		a.words = a.words[:n]
		return a.trim()
	case len(a.words) < n:
		words := make([]uint, n)
		copy(words[:len(a.words)], a.words)
		a.words = words
		return a
	default:
		return a.trim()
	}
}

// SetBits sets several bits.
func (a *LMask) SetBits(bits ...int) *LMask {
	for i := 0; i < len(bits); i++ {
		a.words[bits[i]/WordBitCap] |= 1 << (bits[i] % WordBitCap)
	}

	return a.trim()
}

// ShiftLeft ...
func (a *LMask) ShiftLeft(bits int) *LMask {
	if a.bitCap != 0 {
		k := bits % WordBitCap / a.bitCap
		if 0 < k {
			copy(a.words[:len(a.words)-k], a.words[k:])
			for i := 0; i < k; i++ {
				a.words[i] = 0
			}
		}

		if r := bits % WordBitCap; 0 < r {
			d := WordBitCap - r
			h0 := a.words[k] >> d
			a.words[k] <<= r
			for i := k + 1; i < len(a.words); i++ {
				h1 := a.words[i] >> d
				a.words[i] = (a.words[i] << r) | h0
				h0 = h1
			}
		}
	}

	return a.trim()
}

// ShiftRight ...
func (a *LMask) ShiftRight(bits int) *LMask {
	if a.bitCap != 0 {
		k := bits / a.bitCap
		if 0 < k {
			copy(a.words[k:], a.words[:len(a.words)-k])
			for i := len(a.words) - k; i < len(a.words); i++ {
				a.words[i] = 0
			}
		}

		if r := bits % WordBitCap; 0 < r {
			d := WordBitCap - r
			h0 := a.words[k] << d
			a.words[k] >>= r
			for i := k - 1; 0 <= i; i-- {
				h1 := a.words[i] << d
				a.words[i] = (a.words[i] >> r) | h0
				h0 = h1
			}
		}
	}

	return a.trim()
}

// String returns a string-representation of a bitmask formatted as [0, 0, ..., 0]
func (a *LMask) String() string {
	return fmt.Sprint(a.words)
}

// Words returns a copy of the words in the bitmask.
func (a *LMask) Words() []uint {
	return append(make([]uint, 0, len(a.words)), a.words...)
}

// ------------------------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------------------------

// min returns the minimum value.
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// trim unsets any leading bits outside the range [0, bitCap).
func (a *LMask) trim() *LMask {
	if len(a.words) != 0 {
		if r := a.bitCap % WordBitCap; r != 0 {
			a.words[len(a.words)-1] &^= WordMax << r
		}
	}

	return a
}
