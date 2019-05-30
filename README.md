# atomics
All the AMD64 atomic instructions, exposed as Go functions: [![GoDoc](https://godoc.org/github.com/CAFxX/atomics?status.svg)](https://godoc.org/github.com/CAFxX/atomics)

⚠️ **PRE-ALPHA** ⚠️

## Features

- `LOCK {ADD,ADC,AND,BTC,BTR,BTS,CMPXCHG,CMPXCH8B,CMPXCHG16B,DEC,INC,NEG,NOT,OR,SBB,SUB,XOR,XADD} and XCHG`
- `XBEGIN, XABORT, XEND, XTEST` (Intel TSX RTM transactional memory instructions)

## Notes

- Before using the TSX RTM functions, you must ensure that the processors supports them by using e.g. github.com/intel-go/cpuid:
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

