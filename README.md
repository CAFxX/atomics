# atomics
All the AMD64 atomic instructions, exposed as Go functions: [![GoDoc](https://godoc.org/github.com/CAFxX/atomics?status.svg)](https://godoc.org/github.com/CAFxX/atomics)

⚠️ **PRE-ALPHA** ⚠️

## Features

- `LOCK {ADD,ADC,AND,BTC,BTR,BTS,CMPXCHG,CMPXCH8B,CMPXCHG16B,DEC,INC,NEG,NOT,OR,SBB,SUB,XOR,XADD}` and `XCHG`, for all possible operand sizes
- `XBEGIN`, `XABORT`, `XEND` and `XTEST` (Intel TSX RTM transactional memory instructions)

## Notes

- This is currently an AMD64-only package; none of the functions are available if not building for `amd64`.
- Some of the functions are hard to use reliably due to the nature of the Go compiler (e.g. those that depend on the contents of the FLAGS register, like `AddWithCarry` or `SubtractWithBorrow`); they are included for completeness and because they can be used, if proper care is used.
- Before using the TSX RTM functions, you must ensure that the processors supports them by using e.g. [github.com/intel-go/cpuid](https://github.com/intel-go/cpuid):
    ```
    if cpuid.HasExtendedFeature(cpuid.RTM) {
        // OK to use XBEGIN/XTEST/XABORT/XEND
    }
    ```

## Development

The atomics instructions are defined in `gen.go`. `gen.go` will produce the files `atomic_amd64.{go,s}` by running:

```
go generate gen.go
```

