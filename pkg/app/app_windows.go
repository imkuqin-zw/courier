// +build windows

package app

import (
	"os"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGQUIT, os.Interrupt}
