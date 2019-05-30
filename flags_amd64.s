#include "textflag.h"

// func Flags() uint8
TEXT ·Flags(SB), NOSPLIT, $0-1
    LAHF
    MOVB AH, ret+0(FP)
    RET
    