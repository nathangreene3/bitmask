# Bitmask

```go
go get github.com/nathangreene3/bitmask
```

A `uint` is treated as an implementation of a bitmask.
 
 ## Examples

### Fibonacci numbers

```go
var fibs uint = SetBits(0, 1, 2)
for n := Count(fibs); n < BitCap; n++ {
	var b int = PrevBit(fibs, BitCap)
	fibs = SetBit(fibs, PrevBit(fibs, b)+b)
}

for a0, a1 := 0, 1; a1 < BitCap; a0, a1 = a1, a0+a1 {
	if !MasksBit(fibs, a1) {
		t.Errorf("\nexpected %d to be masked\n", a1)
	}
}
```

### Prime numbers (sieve of Eratosthenes)

```go
var (
	sqrtBitCap int  = int(math.Sqrt(float64(BitCap)))
	primes     uint = ClrBits(Max, 0, 1)
)

for p := 2; p <= sqrtBitCap; p = NextBit(primes, p) {
	if MasksBit(primes, p) {
		for k := p * p; k < BitCap; k += p {
			primes = ClrBit(primes, k)
		}
	}
}
```

### Square numbers

```go
var squares uint
for n := 0; n < BitCap; n++ {
	if n2 := n * n; n2 < BitCap {
		squares = SetBit(squares, n2)
	}
}
```
