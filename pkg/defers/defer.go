package defers

import (
	"sync"
)

type Defer struct {
	sync.Mutex
	fns []func() error
}

func NewDefer() *Defer {
	return &Defer{
		fns: make([]func() error, 0),
	}
}

func (d *Defer) Register(fns ...func() error) {
	d.Lock()
	defer d.Unlock()
	d.fns = append(d.fns, fns...)
}

func (d *Defer) Done() {
	d.Lock()
	defer d.Unlock()
	for i := len(d.fns) - 1; i >= 0; i-- {
		_ = d.fns[i]()
	}
}
