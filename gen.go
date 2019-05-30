//go:generate go run gen.go -out atomic_amd64.s -stubs atomic_amd64.go
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

type R operand.Op
type Fn func(R, []R) []R

func main() {
	cases := []struct {
		Name  string
		NArgs int
		NRets int
		B     int
		Fn    Fn
	}{
		{"Add", 1, 0, 8, X(func(p R, r R) { LOCK(); ADDB(r, p) })},
		{"Add", 1, 0, 16, X(func(p R, r R) { LOCK(); ADDW(r, p) })},
		{"Add", 1, 0, 32, X(func(p R, r R) { LOCK(); ADDL(r, p) })},
		{"Add", 1, 0, 64, X(func(p R, r R) { LOCK(); ADDQ(r, p) })},

		{"AddWithCarry", 1, 0, 8, X(func(p R, r R) { LOCK(); ADCB(r, p) })},
		{"AddWithCarry", 1, 0, 16, X(func(p R, r R) { LOCK(); ADCW(r, p) })},
		{"AddWithCarry", 1, 0, 32, X(func(p R, r R) { LOCK(); ADCL(r, p) })},
		{"AddWithCarry", 1, 0, 64, X(func(p R, r R) { LOCK(); ADCQ(r, p) })},

		{"And", 1, 0, 8, X(func(p R, r R) { LOCK(); ANDB(r, p) })},
		{"And", 1, 0, 16, X(func(p R, r R) { LOCK(); ANDW(r, p) })},
		{"And", 1, 0, 32, X(func(p R, r R) { LOCK(); ANDL(r, p) })},
		{"And", 1, 0, 64, X(func(p R, r R) { LOCK(); ANDQ(r, p) })},

		{"BitTestAndComplement", 1, 1, 16, X(func(p R, r R) R { LOCK(); BTCW(r, p); return FLAG(SETCS, 16) })},
		{"BitTestAndComplement", 1, 1, 32, X(func(p R, r R) R { LOCK(); BTCL(r, p); return FLAG(SETCS, 32) })},
		{"BitTestAndComplement", 1, 1, 64, X(func(p R, r R) R { LOCK(); BTCQ(r, p); return FLAG(SETCS, 64) })},

		{"BitTestAndReset", 1, 1, 16, X(func(p R, r R) R { LOCK(); BTRW(r, p); return FLAG(SETCS, 16) })},
		{"BitTestAndReset", 1, 1, 32, X(func(p R, r R) R { LOCK(); BTRL(r, p); return FLAG(SETCS, 32) })},
		{"BitTestAndReset", 1, 1, 64, X(func(p R, r R) R { LOCK(); BTRQ(r, p); return FLAG(SETCS, 64) })},

		{"BitTestAndSet", 1, 1, 16, X(func(p R, r R) R { LOCK(); BTSW(r, p); return FLAG(SETCS, 16) })},
		{"BitTestAndSet", 1, 1, 32, X(func(p R, r R) R { LOCK(); BTSL(r, p); return FLAG(SETCS, 32) })},
		{"BitTestAndSet", 1, 1, 64, X(func(p R, r R) R { LOCK(); BTSQ(r, p); return FLAG(SETCS, 64) })},

		{"CompareAndSwap", 2, 1, 8, X(func(p R, r []R) R {
			MOVB(r[0], reg.AL)
			LOCK()
			CMPXCHGB(r[1], p)
			return reg.AL
		})},
		{"CompareAndSwap", 2, 1, 16, X(func(p R, r []R) R {
			MOVW(r[0], reg.AX)
			LOCK()
			CMPXCHGW(r[1], p)
			return reg.AX
		})},
		{"CompareAndSwap", 2, 1, 32, X(func(p R, r []R) R {
			MOVL(r[0], reg.EAX)
			LOCK()
			CMPXCHGL(r[1], p)
			return reg.EAX
		})},
		{"CompareAndSwap", 2, 1, 64, X(func(p R, r []R) R {
			MOVQ(r[0], reg.RAX)
			LOCK()
			CMPXCHGQ(r[1], p)
			return reg.RAX
		})},
		{"CompareAndSwap2x", 4, 2, 64, func(p R, r []R) []R {
			MOVQ(r[0], reg.RAX)
			MOVQ(r[1], reg.RDX)
			MOVQ(r[2], reg.RBX)
			MOVQ(r[3], reg.RCX)
			LOCK()
			CMPXCHG16B(p)
			return []R{reg.RAX, reg.RDX}
		}},

		{"Decrement", 0, 0, 8, X(func(p R) { LOCK(); DECB(p) })},
		{"Decrement", 0, 0, 16, X(func(p R) { LOCK(); DECW(p) })},
		{"Decrement", 0, 0, 32, X(func(p R) { LOCK(); DECL(p) })},
		{"Decrement", 0, 0, 64, X(func(p R) { LOCK(); DECQ(p) })},

		{"DecrementAndCheckZero", 0, 1, 8, X(func(p R) R { LOCK(); DECB(p); return FLAG(SETEQ, 8) })},
		{"DecrementAndCheckZero", 0, 1, 16, X(func(p R) R { LOCK(); DECW(p); return FLAG(SETEQ, 16) })},
		{"DecrementAndCheckZero", 0, 1, 32, X(func(p R) R { LOCK(); DECL(p); return FLAG(SETEQ, 32) })},
		{"DecrementAndCheckZero", 0, 1, 64, X(func(p R) R { LOCK(); DECQ(p); return FLAG(SETEQ, 64) })},

		{"Increment", 0, 0, 8, X(func(p R) { LOCK(); INCB(p) })},
		{"Increment", 0, 0, 16, X(func(p R) { LOCK(); INCW(p) })},
		{"Increment", 0, 0, 32, X(func(p R) { LOCK(); INCL(p) })},
		{"Increment", 0, 0, 64, X(func(p R) { LOCK(); INCQ(p) })},

		{"IncrementAndCheckZero", 0, 1, 8, X(func(p R) R { LOCK(); INCB(p); return FLAG(SETEQ, 8) })},
		{"IncrementAndCheckZero", 0, 1, 16, X(func(p R) R { LOCK(); INCW(p); return FLAG(SETEQ, 16) })},
		{"IncrementAndCheckZero", 0, 1, 32, X(func(p R) R { LOCK(); INCL(p); return FLAG(SETEQ, 32) })},
		{"IncrementAndCheckZero", 0, 1, 64, X(func(p R) R { LOCK(); INCQ(p); return FLAG(SETEQ, 64) })},

		{"Negate", 0, 0, 8, X(func(p R) { LOCK(); NEGB(p) })},
		{"Negate", 0, 0, 16, X(func(p R) { LOCK(); NEGW(p) })},
		{"Negate", 0, 0, 32, X(func(p R) { LOCK(); NEGL(p) })},
		{"Negate", 0, 0, 64, X(func(p R) { LOCK(); NEGQ(p) })},

		{"Not", 0, 0, 8, X(func(p R) { LOCK(); NOTB(p) })},
		{"Not", 0, 0, 16, X(func(p R) { LOCK(); NOTW(p) })},
		{"Not", 0, 0, 32, X(func(p R) { LOCK(); NOTL(p) })},
		{"Not", 0, 0, 64, X(func(p R) { LOCK(); NOTQ(p) })},

		{"Or", 1, 0, 8, X(func(p R, r R) { LOCK(); ORB(r, p) })},
		{"Or", 1, 0, 16, X(func(p R, r R) { LOCK(); ORW(r, p) })},
		{"Or", 1, 0, 32, X(func(p R, r R) { LOCK(); ORL(r, p) })},
		{"Or", 1, 0, 64, X(func(p R, r R) { LOCK(); ORQ(r, p) })},

		{"SubtractWithBorrow", 1, 0, 8, X(func(p R, r R) { LOCK(); SBBB(r, p) })},
		{"SubtractWithBorrow", 1, 0, 16, X(func(p R, r R) { LOCK(); SBBW(r, p) })},
		{"SubtractWithBorrow", 1, 0, 32, X(func(p R, r R) { LOCK(); SBBL(r, p) })},
		{"SubtractWithBorrow", 1, 0, 64, X(func(p R, r R) { LOCK(); SBBQ(r, p) })},

		{"Subtract", 1, 0, 8, X(func(p R, r R) { LOCK(); SUBB(r, p) })},
		{"Subtract", 1, 0, 16, X(func(p R, r R) { LOCK(); SUBW(r, p) })},
		{"Subtract", 1, 0, 32, X(func(p R, r R) { LOCK(); SUBL(r, p) })},
		{"Subtract", 1, 0, 64, X(func(p R, r R) { LOCK(); SUBQ(r, p) })},

		{"Xor", 1, 0, 8, X(func(p R, r R) { LOCK(); XORB(r, p) })},
		{"Xor", 1, 0, 16, X(func(p R, r R) { LOCK(); XORW(r, p) })},
		{"Xor", 1, 0, 32, X(func(p R, r R) { LOCK(); XORL(r, p) })},
		{"Xor", 1, 0, 64, X(func(p R, r R) { LOCK(); XORQ(r, p) })},

		{"XorAndCheckZero", 1, 1, 8, X(func(p R, r R) R { LOCK(); XORB(r, p); return FLAG(SETEQ, 8) })},
		{"XorAndCheckZero", 1, 1, 16, X(func(p R, r R) R { LOCK(); XORW(r, p); return FLAG(SETEQ, 16) })},
		{"XorAndCheckZero", 1, 1, 32, X(func(p R, r R) R { LOCK(); XORL(r, p); return FLAG(SETEQ, 32) })},
		{"XorAndCheckZero", 1, 1, 64, X(func(p R, r R) R { LOCK(); XORQ(r, p); return FLAG(SETEQ, 64) })},

		{"AddAndSwap", 1, 1, 8, X(func(p R, r R) R { LOCK(); XADDB(r, p); return r })},
		{"AddAndSwap", 1, 1, 16, X(func(p R, r R) R { LOCK(); XADDW(r, p); return r })},
		{"AddAndSwap", 1, 1, 32, X(func(p R, r R) R { LOCK(); XADDL(r, p); return r })},
		{"AddAndSwap", 1, 1, 64, X(func(p R, r R) R { LOCK(); XADDQ(r, p); return r })},

		{"Swap", 1, 1, 8, X(func(p R, r R) R { XCHGB(r, p); return r })},
		{"Swap", 1, 1, 16, X(func(p R, r R) R { XCHGW(r, p); return r })},
		{"Swap", 1, 1, 32, X(func(p R, r R) R { XCHGL(r, p); return r })},
		{"Swap", 1, 1, 64, X(func(p R, r R) R { XCHGQ(r, p); return r })},
	}

	// ADD, ADC, AND, BTC, BTR, BTS, CMPXCHG, CMPXCH8B, CMPXCHG16B, DEC, INC, NEG, NOT, OR, SBB, SUB, XOR, XADD, and XCHG

	const natB = 64
	for _, c := range cases {
		for _, s := range []string{"Int", "Uint"} {
			t := fmt.Sprintf("%s%d", strings.ToLower(s), c.B)
			T := fmt.Sprintf("%s%d", s, c.B)
			gen(c.Name, c.NArgs, c.NRets, t, T, c.B, c.Fn)
		}
		if c.B == natB {
			gen(c.Name, c.NArgs, c.NRets, "int", "Int", c.B, c.Fn)
			gen(c.Name, c.NArgs, c.NRets, "uint", "Uint", c.B, c.Fn)
			gen(c.Name, c.NArgs, c.NRets, "uintptr", "Uintptr", c.B, c.Fn)
		}
	}

	Generate()

	asm, err := ioutil.ReadFile("atomic_amd64.s")
	if err != nil {
		panic(err)
	}
	asm = bytes.ReplaceAll(asm, []byte("// LOCK //"), []byte("LOCK"))
	//asm = []byte(dedupFuncs(string(asm)))
	err = ioutil.WriteFile("atomic_amd64.s", []byte(asm), 0)
	if err != nil {
		panic(err)
	}
}

func gp(b int) reg.GPVirtual {
	switch b {
	case 8:
		return GP8()
	case 16:
		return GP16()
	case 32:
		return GP32()
	case 64:
		return GP64()
	default:
		panic("unknown bit size")
	}
}

func arg(n int) string {
	return fmt.Sprintf("a%d", n)
}

func gen(name string, nargs, nrets int, t, T string, b int, fn Fn) {
	_args := []string{}
	for i := 0; i < nargs; i++ {
		_args = append(_args, arg(i))
	}
	args := ", " + strings.Join(_args, ", ") + " " + t
	if nargs == 0 {
		args = ""
	}

	_rets := []string{}
	for i := 0; i < nrets; i++ {
		_rets = append(_rets, t)
	}
	rets := "(" + strings.Join(_rets, ", ") + ")"

	signature := fmt.Sprintf("func(addr *%s %s) %s", t, args, rets)

	TEXT(name+T, NOSPLIT, signature)

	p := operand.Mem{Base: Load(Param("addr"), GP64())}

	var r []R
	for i := 0; i < nargs; i++ {
		r = append(r, Load(ParamIndex(i+1), gp(b)))
	}

	r = fn(p, r)

	for i, o := range r {
		Store(o.(reg.Register), ReturnIndex(i))
	}

	RET()
}

func LOCK() {
	Comment("LOCK //")
}

func X(fn interface{}) func(R, []R) []R {
	switch f := fn.(type) {
	case func(R, []R) []R:
		return f
	case func(R, []R) R:
		return func(p R, r []R) []R { return []R{f(p, r)} }
	case func(R, []R):
		return func(p R, r []R) []R { f(p, r); return nil }
	case func(R, R) R:
		return func(p R, r []R) []R { return []R{f(p, r[0])} }
	case func(R, R):
		return func(p R, r []R) []R { f(p, r[0]); return nil }
	case func(R) R:
		return func(p R, _ []R) []R { return []R{f(p)} }
	case func(R):
		return func(p R, _ []R) []R { f(p); return nil }
	default:
		panic(fmt.Sprintf("unknown fn type: %+v", fn))
	}
}

func FLAG(flagFn func(operand.Op), bits int) R {
	r := gp(bits)
	if bits == 8 {
		flagFn(r)
	} else {
		flagFn(r.As8())
	}
	return r
}

func dedupFuncs(asm string) string {
	asm += "\n"
	// drop empty and comment lines
	asm = regexp.MustCompile("(?m)^(//.*)?\n+^").ReplaceAllString(asm, "")
	// find function boundaries
	fi := regexp.MustCompile("(?m)^(TEXT.*\n|$)").FindAllStringIndex(asm, -1)
	m := map[string][]string{}
	o := []string{}
	for i, f := range fi[:len(fi)-1] {
		// start is where the func declaration starts, body is where the body starts
		start, body, end := f[0], f[1], fi[i+1][0]
		// find duplicates
		if _, found := m[asm[body:end]]; !found {
			// store the order
			o = append(o, asm[body:end])
		}
		m[asm[body:end]] = append(m[asm[body:end]], asm[start:body])
	}
	// emit deduped funcs
	asm = "// Code generated by command: go run gen.go -out atomic_amd64.s -stubs atomic_amd64.go. DO NOT EDIT.\n" +
		"#include \"textflag.h\"\n"
	re := regexp.MustCompile("^TEXT (.*?\\(SB\\))")
	for _, body := range o {
		decls := m[body]
		fname := re.FindStringSubmatch(decls[len(decls)-1])
		for _, decl := range decls[:len(decls)-1] {
			asm += decl
			asm += "	JMP " + fname[1] + "\n"
		}
		asm += decls[len(decls)-1]
		asm += body
	}
	return asm
}
