package umask

import (
	"math/bits"
	"strconv"
)

const (
	// BitCap is the maximum number of bits in a bitmask.
	BitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// Max is the maximum bitmask.
	Max UMask = ^Zero

	// One is a bitmask value of 1.
	One UMask = 1

	// Zero is a bitmask value of 0.
	Zero UMask = 0
)

// UMask is an implementation of a bitmask.
type UMask uint

// -------------------------------------------------------------------------
// Bitwise functionality
// -------------------------------------------------------------------------

// And returns a bitmask with only the bits set that are common to both bitmasks.
func (a UMask) And(b UMask) UMask {
	return a & b
}

// AndNot returns a bitmask with the bits set in a and the bits not set in b.
func (a UMask) AndNot(b UMask) UMask {
	return a &^ b
}

// NAnd ...
func (a UMask) NAnd(b UMask) UMask {
	return ^(a & b)
}

// NOr ...
func (a UMask) NOr(b UMask) UMask {
	return ^(a | b)
}

// Not inverts a bitmask.
func (a UMask) Not() UMask {
	return ^a
}

// Or returns a bitmask with the bits set in either a or b.
func (a UMask) Or(b UMask) UMask {
	return a | b
}

// XNOr ...
func (a UMask) XNOr(b UMask) UMask {
	return ^(a ^ b)
}

// XOr returns the bits of a and b that are set, but not simultaneously set in both a and b.
func (a UMask) XOr(b UMask) UMask {
	return a ^ b
}

// -------------------------------------------------------------------------
// Set functionality
// -------------------------------------------------------------------------

// BitLen returns the minimum number of bits representing a.
func (a UMask) BitLen() int {
	return bits.Len(uint(a))
}

// Bits returns the bits that are set in a bitmask.
func (a UMask) Bits() []int {
	bits := make([]int, 0, BitCap)
	for bit := a.NextBit(-1); bit < BitCap; bit = a.NextBit(bit) {
		bits = append(bits, bit)
	}

	return bits
}

// Clr returns a bitmask with the bits of a given bitmask cleared.
func (a UMask) Clr(b UMask) UMask {
	return a &^ b
}

// ClrBit returns a bitmask with the given bit cleared.
func (a UMask) ClrBit(bit int) UMask {
	return a &^ (1 << bit)
}

// ClrBits returns a bitmask with a given number of bits unset.
func (a UMask) ClrBits(bits ...int) UMask {
	for i := 0; i < len(bits); i++ {
		a &^= 1 << bits[i]
	}

	return a
}

// Count returns the number of bits set in a bitmask.
func (a UMask) Count() int {
	return bits.OnesCount(uint(a))
}

// Fmt returns a representation of a bitmask in a given base on range [2, 36].
func (a UMask) Fmt(base int) string {
	return strconv.FormatUint(uint64(a), base)
}

// LSh returns a bitmask with all bits shifted to the left a given number of bits.
func (a UMask) LSh(bits int) UMask {
	return a << bits
}

// Masks determines if the bits set in b are set in a.
func (a UMask) Masks(b UMask) bool {
	return a&b == b
}

// MasksBit determines if a bit is set.
func (a UMask) MasksBit(bit int) bool {
	b := One << bit
	return a&b == b
}

// NextBit returns the next set bit. If there is no next set bit, then the bit capacity is returned.
func (a UMask) NextBit(bit int) int {
	bit = clamp(bit+1, 0, BitCap)
	return bits.TrailingZeros(uint(a) >> bit << bit)
}

// PrevBit returns the previous set bit. If there is no previous set bit, then -1 is returned.
func (a UMask) PrevBit(bit int) int {
	bit = BitCap - clamp(bit, 0, BitCap)
	return BitCap - bits.LeadingZeros(uint(a)<<bit>>bit) - 1
}

// RSh returns a bitmask with all bits shifted to the right a given number of bits.
func (a UMask) RSh(bits int) UMask {
	return a >> bits
}

// Set returns a bitmask with bits set in a or b.
func (a UMask) Set(b UMask) UMask {
	return a | b
}

// SetBit sets a bit in a bitmask.
func (a UMask) SetBit(bit int) UMask {
	return a | (1 << bit)
}

// SetBits sets several bits in a bitmask.
func (a UMask) SetBits(bits ...int) UMask {
	for i := 0; i < len(bits); i++ {
		a |= 1 << bits[i]
	}

	return a
}

// -------------------------------------------------------------------------
// Helper functionality
// -------------------------------------------------------------------------

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
