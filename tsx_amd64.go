package atomics

// XBEGIN starts a TSX RTM transaction.
// The processor must support TSX or a panic will occur.
func XBEGIN() uint32

// XEND commits a TSX RTM transaction. A transaction must be open.
// The processor must support TSX or a panic will occur.
func XEND()

// XTEST returns true if we are currently inside a TSX RTM transaction.
// The processor must support TSX or a panic will occur.
func XTEST() bool

// XABORT aborts the current TSX RTM transaction, if any. The provided
// reason will be available in the return value of XBEGIN. If no
// transaction is executing, it is a no-op.
// The processor must support TSX or a panic will occur.
func XABORT(reason uint8)
