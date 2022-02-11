# LMask

An `LMask` is an implementation of a bitmask having arbitrary precision.

## Examples

### Fibonacci numbers

```go
var (
	maxCount int    = 50
	fibs     *LMask = Zero(3).SetBits(1, 2)
)

for fibs.Count() < maxCount {
	var b int = fibs.PrevBit(fibs.BitLen())
	b += fibs.PrevBit(b)
	if fibs.BitCap() <= b {
		fibs.SetBitCap(b + 1)
	}

	fibs.SetBit(b)
}
```

### Prime numbers (sieve of Eratosthenes)

```go
var (
    bitCap     int    = WordBitCap << 2
    sqrtBitCap int    = int(math.Sqrt(float64(bitCap)))
    primes     *LMask = Max(bitCap).ClrBits(0, 1)
)

for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
    if primes.MasksBit(p) {
        for k := p * p; k < bitCap; k += p {
            primes.ClrBit(k)
        }
    }
}
```

### Square numbers

```go
var (
    bitCap  = WordBitCap + 1
    squares = One(bitCap)
)

for i := 1; ; i += 2 {
    var s int = squares.PrevBit(bitCap) + i
    if bitCap <= s {
        break
    }

    squares.SetBit(s)
}
```

## TODO

* Finish unit testing.
* Finish documentation.
* Compare performance and features to other implementations.
