# UMask

A `UMask` is an aliased `uint` that extends the basic bitmasking operations on a `uint`.

## Examples

### Fibonacci numbers

```go
    fibs := Zero.SetBits(0, 1, 2)
    for n := fibs.Count(); n < BitCap; n++ {
        b := fibs.PrevBit(BitCap)
        fibs = fibs.SetBit(fibs.PrevBit(b) + b)
    }
```

### Primes (Sieve of Eratosthenes)

```go
var (
    sqrtBitCap = int(math.Sqrt(float64(BitCap)))
    primes     = Max.ClrBits(0, 1)
)

for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
    if primes.MasksBit(p) {
        for k := p * p; k < BitCap; k += p {
            primes = primes.ClrBit(k)
        }
    }
}
```

### Squares

```go
squares := One
for i := 1; ; i += 2 {
    s := squares.PrevBit(BitCap) + i
    if BitCap <= s {
        break
    }

    squares = squares.SetBit(s)
}
```

## TODO

* Finish unit testing.
* Finish documentation.
* Implement uint32 and uint64 variants of UMask. Consider uint8, byte, and uint16 as well.
* Compare performance and features to other implementations.
