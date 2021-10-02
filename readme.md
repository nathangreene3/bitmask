# Bitmask

```go
go get github.com/nathangreene3/bitmask
```

## LMask

An `LMask` is a bitmask of arbitrary precision. The bit capacity may be specified as any non-negative value.

### Examples

#### Fibonacci Numbers

```go
var (
    fibs     = Zero(4).SetBits(0, 1, 2)
    maxCount = 50
)

for n := fibs.Count(); n < maxCount; n++ {
    var (
        b0 = fibs.PrevBit(fibs.BitCap())
        b1 = fibs.PrevBit(b0) + b0
    )

    if fibs.BitCap() <= b1 {
        fibs.SetBitCap(b1 << 1)
    }

    fibs.SetBit(b1)
}
```

#### Primes (Sieve of Eratosthenes)

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

#### Squares

```go
var (
    bitCap  = WordBitCap + 1
    squares = One(bitCap)
)

for i := 1; ; i += 2 {
    var square = squares.PrevBit(bitCap) + i
    if bitCap <= square {
        break
    }

    squares.SetBit(square)
}
```

## UMask

A `UMask` is an aliased `uint` that extends the basic bitmasking operations on a `uint`.

### UMask Examples

#### Fibonacci Numbers

```go
    var (
        fibs = Zero.SetBits(0, 1, 2)
    )

    for n := fibs.Count(); n < BitCap; n++ {
        var (
            b0 = fibs.PrevBit(BitCap)
            b1 = fibs.PrevBit(b0) + b0
        )

        fibs = fibs.SetBit(b1)
    }
```

#### Primes (Sieve of Eratosthenes)

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

#### Squares

```go
squares := One
for i := 1; ; i += 2 {
    square := squares.PrevBit(BitCap) + i
    if BitCap <= square {
        break
    }

    squares = squares.SetBit(square)
}
```

## TODO

* Finish unit testing.
* Finish documentation.
* Implement uint32 and uint64 variants of UMask. Consider uint8, byte, and uint16 as well.
* Compare performance and features to other implementations.
