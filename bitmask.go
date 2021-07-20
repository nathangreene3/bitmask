package bitmask

import "math/bits"

const (
	// BitCap is the maximum number of bits in a uint.
	BitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// Max is the maximum uint.
	Max = 1<<BitCap - 1
)

// ------------------------------------------------------------------------------------
// Logic functionality
// ------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------
// Additional functionality
// ------------------------------------------------------------------------------------

// Clr ...
func Clr(a, b uint) uint {
	return a &^ b
}

// ClrBit ...
func ClrBit(a uint, bit int) uint {
	return a &^ (1 << bit)
}

// ClrBits ...
func ClrBits(a uint, bits ...int) uint {
	if 0 < a {
		for i := 0; i < len(bits); i++ {
			a &^= 1 << bits[i]
		}
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

// clamp returns min if a < min, max if max < a, or otherwise a.
func clamp(a, min, max int) int {
	switch {
	case a < min:
		return min
	case max < a:
		return max
	default:
		return a
	}
}
