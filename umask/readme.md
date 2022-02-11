# UMask

```go
go get github.com/nathangreene3/bitmask/umask
```

A `UMask` is an implementation of a bitmask.

## Examples

### Fibonacci numbers

```go
var fibs UMask = Zero.SetBits(1, 2)
for n := fibs.Count(); n < BitCap; n++ {
    var b int = fibs.PrevBit(BitCap)
    fibs = fibs.SetBit(fibs.PrevBit(b) + b)
}
```

### Prime numbers (Sieve of Eratosthenes)

```go
var (
    sqrtBitCap int   = int(math.Sqrt(float64(BitCap)))
    primes     UMask = Max.ClrBits(0, 1)
)

for p := 2; p <= sqrtBitCap; p = primes.NextBit(p) {
    if primes.MasksBit(p) {
        for k := p * p; k < BitCap; k += p {
            primes = primes.ClrBit(k)
        }
    }
}
```

### Square numbers

```go
var squares UMask = One
for i := 1; ; i += 2 {
    var s int = squares.PrevBit(BitCap) + i
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
