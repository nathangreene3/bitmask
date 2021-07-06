package bitmask

import "math/bits"

const (
	// WordBitCap is the maximum number of bits in a UMask.
	WordBitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// MaxWord is the maximum uint.
	WordMax = 1<<WordBitCap - 1
)

// And returns the intersection of a and b.
func And(a, b uint) uint {
	return a & b
}

// AndNot ...
func AndNot(a, b uint) uint {
	return a &^ b
}

// Masks determines if a masks b.
func Masks(a, b uint) bool {
	return a&b == b
}

// MasksBit ...
func MasksBit(a uint, bit int) bool {
	b := uint(1) << bit
	return a&b == b
}

// NAnd returns the difference of a-b.
func NAnd(a, b uint) uint {
	return a &^ b
}

// NextBit ...
func NextBit(a uint, bit int) int {
	if bit++; bit < WordBitCap {
		if a = a >> bit << bit; 0 < a {
			return bits.TrailingZeros(a)
		}
	}

	return WordBitCap
}

// NOr ...
func NOr(a, b uint) uint {
	return ^(a | b)
}

// Not returns the inversion of a.
func Not(a uint) uint {
	return ^a
}

// Or returns the union of a and b.
func Or(a, b uint) uint {
	return a | b
}

// XNOr ...
func XNOr(a, b uint) uint {
	return ^(a ^ b)
}

// XOr returns the symmetric difference of a and b.
func XOr(a, b uint) uint {
	return a ^ b
}
