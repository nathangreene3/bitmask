package umask

import (
	"math/bits"
	"strconv"
)

const (
	// BitCap is the maximum number of bits in a UMask.
	BitCap = 32 << (^uint(0) >> 32 & 1) // Source: bits.UintSize

	// Max is the maximum UMask.
	Max UMask = 1<<BitCap - 1
)

// UMask is either a 32 or 64-bit bitmask.
type UMask uint

func One() UMask {
	return 1
}

func Zero() UMask {
	return 0
}

// ------------------------------------------------------------------------------------
// Logic functionality TODO: Add N-variants
// ------------------------------------------------------------------------------------

// And returns a bitmask with only the bits set that are common to both bitmasks.
func (a UMask) And(b UMask) UMask {
	return a & b
}

// AndNot ...
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

// Not inverts a bitmask. This is equivalent to calling UMax.Xor(a).
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

// ------------------------------------------------------------------------------------
// Additional functionality
// ------------------------------------------------------------------------------------

// Base returns a string representing a bitmask in a given base n where 2 <= n <= 36.
func (a UMask) Base(n int) string {
	return strconv.FormatUint(uint64(a), n)
}

// Bin returns a string representing a bitmask in binary.
func (a UMask) Bin() string {
	return strconv.FormatUint(uint64(a), 2)
}

// Clr returns a bitmask with the bits of each given bitmask b cleared from a.
func (a UMask) Clr(b UMask) UMask {
	return a &^ b
}

// ClrBit returns a bitmask with the given bits cleared from a.
func (a UMask) ClrBit(bit int) UMask {
	return a &^ (1 << bit)
}

// ClrBits ...
func (a UMask) ClrBits(bits ...int) UMask {
	if 0 < a {
		for i := 0; i < len(bits); i++ {
			a &^= 1 << bits[i]
		}
	}

	return a
}

// Count ...
func (a UMask) Count() int {
	return bits.OnesCount(uint(a))
}

// Dec returns a string representing a bitmask in decimal.
func (a UMask) Dec() string {
	return strconv.FormatUint(uint64(a), 10)
}

// Hex returns a string representing a bitmask in hexidecimal.
func (a UMask) Hex() string {
	return strconv.FormatUint(uint64(a), 16)
}

// Len ...
func (a UMask) Len() int {
	return bits.Len(uint(a))
}

// Masks determines if the bits set in b are set in a.
func (a UMask) Masks(b UMask) bool {
	return a&b == b
}

// MasksBit determines if a bit is set.
func (a UMask) MasksBit(bit int) bool {
	b := UMask(1) << bit
	return a&b == b
}

// NextBit ...
func (a UMask) NextBit(bit int) int {
	if bit++; bit < BitCap {
		if a = a >> bit << bit; 0 < a {
			return bits.TrailingZeros(uint(a))
		}
	}

	return BitCap
}

// Oct returns a string representing a bitmask in decimal.
func (a UMask) Oct() string {
	return strconv.FormatUint(uint64(a), 8)
}

// Set returns a Bitmask with bits set in each b. This is equivalent to repeatedly calling a.Or(b) for each b.
func (a UMask) Set(b UMask) UMask {
	return a.Or(b)
}

// SetBits ...
func (a UMask) SetBits(bits ...int) UMask {
	if a < Max {
		for i := 0; i < len(bits); i++ {
			a |= 1 << bits[i]
		}
	}

	return a
}

// Lsh returns a Bitmask shifted to the left n times.
func (a UMask) ShiftLeft(bits int) UMask {
	return a << bits
}

// ShiftRight returns a Bitmask shifted to the right n times.
func (a UMask) ShiftRight(bits int) UMask {
	return a >> bits
}

// String returns a string representing a bitmask in decimal.
func (a UMask) String() string {
	return strconv.FormatUint(uint64(a), 10)
}
