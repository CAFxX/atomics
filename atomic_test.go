package atomics_test

import (
	"testing"
	"."
	"github.com/intel-go/cpuid"
)

func TestCompareAndSwapUint8(t *testing.T) {
    a := byte(0xFF)
    atomics.CompareAndSwapUint8(&a, 0xFF, 42) 
    if a != 42 {
        t.Fatal("no swap")
    }
}

func TestCompareAndSwap2xUint64(t *testing.T) {
    a := [2]uint64{1, 2}
    atomics.CompareAndSwap2xUint64(&a[0], 1, 2, 3, 4)
    if a[0] != 3 || a[1] != 4 {
        t.Fatal("no swap")
    }
}

func TestAndUint64(t *testing.T) {
    a := uint64(15)
    atomics.AndUint64(&a, 170)
    if a != 10 {
        t.Fatal("no and")
    }
}

func TestAddWithCarryUint8(t *testing.T) {
    a, b := uint8(100), uint8(50)
	atomics.AddUint8(&a, 156)
	atomics.AddWithCarryUint8(&b, 1)
    if a != 0 || b != 52 {
        t.Fatal("no add with carry")
    }
}

func TestRTM(t *testing.T) {
	if !cpuid.HasExtendedFeature(cpuid.RTM) {
		t.Skip("RTM not supported")
	}
	s := atomics.XBEGIN()
	if s == 0xFFFFFFFF {
		atomics.XABORT(42)
	} else if s>>24 != 42 {
		t.Fatal(s)
	}
}

func TestFlags(t *testing.T) {
	a := uint8(40)

	atomics.AddUint8(&a, 2)
	if f := atomics.Flags(); f & 1 != 0 {
		t.Fatalf("unexpected carry %x", f)
	}

	atomics.AddUint8(&a, 255)
	if f := atomics.Flags(); f & 1 != 1 {
		t.Fatalf("missing carry %x", f)
	}
}