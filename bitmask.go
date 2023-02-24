package bitmask

import "math/bits"

const (
	// BitCap is the maximum number of bits in a bitmask.
	BitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// Max is the maximum bitmask.
	Max = 1<<BitCap - 1
)

// -------------------------------------------------------------------------
// Bitwise functionality
// -------------------------------------------------------------------------

// And ...
func And(a, b uint) uint {
	return a & b
}

// AndNot ...
func AndNot(a, b uint) uint {
	return a &^ b
}

// NAnd ...
func NAnd(a, b uint) uint {
	return ^(a & b)
}

// NOr ...
func NOr(a, b uint) uint {
	return ^(a | b)
}

// Not ...
func Not(a uint) uint {
	return ^a
}

// Or ...
func Or(a, b uint) uint {
	return a | b
}

// XNOr ...
func XNOr(a, b uint) uint {
	return ^(a ^ b)
}

// XOr ...
func XOr(a, b uint) uint {
	return a ^ b
}

// -------------------------------------------------------------------------
// Set functionality
// -------------------------------------------------------------------------

// Bits returns the bits that are set in a bitmask.
func Bits(a uint) []int {
	bits := make([]int, 0, BitCap)
	for bit := NextBit(a, -1); bit < BitCap; bit = NextBit(a, bit) {
		bits = append(bits, bit)
	}

	return bits
}

// Clr ...
func Clr(a, b uint) uint {
	return a &^ b
}

// ClrAll ...
func ClrAll(a uint) uint {
	return 0
}

// ClrBit ...
func ClrBit(a uint, bit int) uint {
	return a &^ (1 << bit)
}

// ClrBits ...
func ClrBits(a uint, bits ...int) uint {
	for i := 0; i < len(bits); i++ {
		a &^= 1 << bits[i]
	}

	return a
}

// Count ...
func Count(a uint) int {
	return bits.OnesCount(a)
}

// Masks ...
func Masks(a, b uint) bool {
	return a&b == b
}

// MasksBit ...
func MasksBit(a uint, bit int) bool {
	b := uint(1) << bit
	return a&b == b
}

// NextBit ...
func NextBit(a uint, bit int) int {
	bit = clamp(bit+1, 0, BitCap)
	return bits.TrailingZeros(a >> bit << bit)
}

// PrevBit ...
func PrevBit(a uint, bit int) int {
	bit = BitCap - clamp(bit, 0, BitCap)
	return BitCap - bits.LeadingZeros(a<<bit>>bit) - 1
}

// Set ...
func Set(a, b uint) uint {
	return a | b
}

// SetAll ...
func SetAll(a uint) uint {
	return Max
}

// SetBit ...
func SetBit(a uint, bit int) uint {
	return a | (1 << bit)
}

// SetBits ...
func SetBits(a uint, bits ...int) uint {
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
