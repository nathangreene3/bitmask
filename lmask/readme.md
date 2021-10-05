# LMask

An `LMask` is a bitmask of arbitrary precision. The bit capacity may be specified as any non-negative value.

## Examples

### Fibonacci numbers

```go
var (
    fibs     = Zero(4).SetBits(0, 1, 2)
    maxCount = 50
)

for n := fibs.Count(); n < maxCount; n++ {
    b0 := fibs.PrevBit(fibs.BitCap())
    b1 := fibs.PrevBit(b0) + b0
    if fibs.BitCap() <= b1 {
        fibs.SetBitCap(b1 << 1)
    }

    fibs.SetBit(b1)
}
```

### Primes (Sieve of Eratosthenes)

```go
var (
    bitCap     = WordBitCap << 2
    sqrtBitCap = int(math.Sqrt(float64(bitCap)))
    primes     = Max(bitCap).ClrBits(0, 1)
)

for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
    if primes.MasksBit(p) {
        for k := p * p; k < bitCap; k += p {
            primes.ClrBit(k)
        }
    }
}
```

### Squares

```go
var (
    bitCap  = WordBitCap + 1
    squares = One(bitCap)
)

for i := 1; ; i += 2 {
    s := squares.PrevBit(bitCap) + i
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
