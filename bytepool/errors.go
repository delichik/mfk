package bytepool

import "errors"

var (
	ErrOutOfRange = errors.New("out of range")
	ErrClosed     = errors.New("closed")
)
