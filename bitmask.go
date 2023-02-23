package bitmask

import (
	"github.com/nathangreene3/bitmask/gmask"
)

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
	return gmask.And(a, b)
}

// AndNot ...
func AndNot(a, b uint) uint {
	return gmask.AndNot(a, b)
}

// NAnd ...
func NAnd(a, b uint) uint {
	return gmask.NAnd(a, b)
}

// NOr ...
func NOr(a, b uint) uint {
	return gmask.NOr(a, b)
}

// Not ...
func Not(a uint) uint {
	return gmask.Not(a)
}

// Or ...
func Or(a, b uint) uint {
	return gmask.Or(a, b)
}

// XNOr ...
func XNOr(a, b uint) uint {
	return gmask.XNOr(a, b)
}

// XOr ...
func XOr(a, b uint) uint {
	return gmask.XOr(a, b)
}

// -------------------------------------------------------------------------
// Set functionality
// -------------------------------------------------------------------------

// Bits returns the bits that are set in a bitmask.
func Bits(a uint) []int {
	return gmask.Bits(a)
}

// Clr ...
func Clr(a, b uint) uint {
	return gmask.Clr(a, b)
}

// ClrBit ...
func ClrBit(a uint, bit int) uint {
	return gmask.ClrBit(a, bit)
}

// ClrBits ...
func ClrBits(a uint, bits ...int) uint {
	return gmask.ClrBits(a, bits...)
}

// Count ...
func Count(a uint) int {
	return gmask.Count(a)
}

// Masks ...
func Masks(a, b uint) bool {
	return gmask.Masks(a, b)
}

// MasksBit ...
func MasksBit(a uint, bit int) bool {
	return gmask.MasksBit(a, bit)
}

// NextBit ...
func NextBit(a uint, bit int) int {
	return gmask.NextBit(a, bit)
}

// PrevBit ...
func PrevBit(a uint, bit int) int {
	return gmask.PrevBit(a, bit)
}

// Set ...
func Set(a, b uint) uint {
	return gmask.Set(a, b)
}

// SetBit ...
func SetBit(a uint, bit int) uint {
	return gmask.SetBit(a, bit)
}

// SetBits ...
func SetBits(a uint, bits ...int) uint {
	return gmask.SetBits(a, bits...)
}
