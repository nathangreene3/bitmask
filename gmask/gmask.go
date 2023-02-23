package gmask

import (
	"errors"
	"math/bits"
)

var (
	ErrInvalidRange = errors.New("one or more values out of range")
)

// UInteger defines supported types that may be used as a bitmask.
type UInteger interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

// ---------------------------------------------------------------------
// 	Reporting functionality
// ---------------------------------------------------------------------

// BitCap returns the number of bits available in a bitmask.
func BitCap[T UInteger]() int {
	return bits.Len64(uint64(^T(0))) // Equivalent to Len64(Max[uint64]())
}

// ---------------------------------------------------------------------
// 	Constructors
// ---------------------------------------------------------------------

// FromBits ...
func FromBits[T UInteger](bits ...int) T {
	var a T
	for i := 0; i < len(bits); i++ {
		a |= 1 << bits[i]
	}

	return a
}

// Max returns the maximum-valued bitmask.
func Max[T UInteger]() T {
	return ^T(0)
}

// One returns T(1).
func One[T UInteger]() T {
	return 1
}

// Zero returns T(0).
func Zero[T UInteger]() T {
	return 0
}

// ---------------------------------------------------------------------
// 	Bitwise functionality
// ---------------------------------------------------------------------

// And ...
func And[T UInteger](a, b T) T {
	return a & b
}

// AndNot ...
func AndNot[T UInteger](a, b T) T {
	return a &^ b
}

// NAnd ...
func NAnd[T UInteger](a, b T) T {
	return ^(a & b)
}

// NOr ...
func NOr[T UInteger](a, b T) T {
	return ^(a | b)
}

// Not ...
func Not[T UInteger](a T) T {
	return ^a
}

// Or ...
func Or[T UInteger](a, b T) T {
	return a | b
}

// XNOr ...
func XNOr[T UInteger](a, b T) T {
	return ^(a ^ b)
}

// XOr ...
func XOr[T UInteger](a, b T) T {
	return a ^ b
}

// ---------------------------------------------------------------------
// 	Set functionality
// ---------------------------------------------------------------------

// Bits returns the bits that are set in a bitmask.
func Bits[T UInteger](a T) []int {
	// Source: https://lemire.me/blog/2023/02/07/bit-hacking-with-go-code/

	b := uint64(a)
	setBits := make([]int, 0, bits.OnesCount64(b))

	for b > 0 {
		setBits = append(setBits, bits.TrailingZeros64(b))
		b &= b - 1
	}

	return setBits
}

// Clr ...
func Clr[T UInteger](a, b T) T {
	return a &^ b
}

// ClrBit ...
func ClrBit[T UInteger](a T, bit int) T {
	return a &^ (1 << bit)
}

// ClrBits ...
func ClrBits[T UInteger](a T, bits ...int) T {
	for i := 0; i < len(bits); i++ {
		a &^= 1 << bits[i]
	}

	return a
}

// Count ...
func Count[T UInteger](a T) int {
	return bits.OnesCount(uint(a))
}

// Masks ...
func Masks[T UInteger](a, b T) bool {
	return a&b == b
}

// MasksBit ...
func MasksBit[T UInteger](a T, bit int) bool {
	b := T(1) << bit
	return a&b == b
}

// NextBit ...
func NextBit[T UInteger](a T, bit int) int {
	bit = clamp(bit+1, 0, BitCap[T]())
	return bits.TrailingZeros64(uint64(a) &^ (1<<bit - 1))
}

// PrevBit ...
func PrevBit[T UInteger](a T, bit int) int {
	bit = clamp(bit, 0, BitCap[T]())
	return 63 - bits.LeadingZeros64(uint64(a)&(1<<bit-1))
}

// Set ...
func Set[T UInteger](a, b T) T {
	return a | b
}

// SetBit ...
func SetBit[T UInteger](a T, bit int) T {
	return a | (1 << bit)
}

// SetBits ...
func SetBits[T UInteger](a T, bits ...int) T {
	for i := 0; i < len(bits); i++ {
		a |= 1 << bits[i]
	}

	return a
}

// Toggle ...
func Toggle[T UInteger](a, b T) T {
	return a ^ b
}

// ToggleBit ...
func ToggleBit[T UInteger](a T, bit int) T {
	return a ^ (1 << bit)
}

// ToggleBits ...
func ToggleBits[T UInteger](a T, bits ...int) T {
	for i := 0; i < len(bits); i++ {
		a ^= 1 << bits[i]
	}

	return a
}

// ---------------------------------------------------------------------
//	Helper functionality
// ---------------------------------------------------------------------

// clamp returns a value restricted to the range [a, b]. The value
// corresponds to the first observed condition in the following:
//
//   - b < a: panic
//   - n < a: a
//   - b < n: b
//   - else: n
func clamp(n, a, b int) int {
	switch {
	case b < a:
		panic(ErrInvalidRange)
	case n < a:
		return a
	case b < n:
		return b
	default:
		return n
	}
}
