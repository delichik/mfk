package queue

import "errors"

var ErrClosed = errors.New("queue has been closed")
var ErrFulled = errors.New("queue is full")
