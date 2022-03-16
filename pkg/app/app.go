package app

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	config2 "dubbo.apache.org/dubbo-go/v3/config"
	"github.com/imkuqin-zw/courier/pkg/defers"
)

func Run(ops ...Option) {
	Init(ops...)
	Start(ops...)
	WaitSignals()
}

func Init(ops ...Option) {
	initOpts(ops...)
	loadEnvAndFlag()
	loadAppCfgFile()
	loadDubboV3()
}

func Start(ops ...Option) {
	if rc == nil || o == nil {
		panic("please call Init before start")
	}
	applyOpts(o, ops...)

	if err := initConsumer(); err != nil {
		panic(err)
	}
	if err := initProvider(); err != nil {
		panic(err)
	}
	if err := initShutdown(); err != nil {
		panic(err)
	}

	rc.Start()
}

func WaitSignals() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, shutdownSignals...)
	s := <-sig
	logger.Info("got os signal " + s.String())
	timeout := rc.Shutdown.GetTimeout()
	if shutdownTimeout > timeout {
		timeout = shutdownTimeout
	}
	time.AfterFunc(timeout, func() {
		logger.Warn("application graceful shutdown timeout")
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	})

	go func() {
		defers.DoBeforeHook()
		config2.BeforeShutdown()
		dumpSignal := map[os.Signal]struct{}{
			syscall.SIGQUIT: {},
			syscall.SIGILL:  {},
			syscall.SIGTRAP: {},
			syscall.SIGABRT: {},
		}
		if _, ok := dumpSignal[s]; ok {
			debug.WriteHeapDump(os.Stdout.Fd())
		}
		defers.DoAfterHook()
		logger.Info("application graceful shutdown")
		os.Exit(0)
	}()
	t := time.Tick(time.Second * 3)
	for {
		out := false
		select {
		case <-sig:
			logger.Warn("application force shutdown, you need to wait at least 3 seconds")
		case <-t:
			out = true
		}
		if out {
			break
		}
	}
	<-sig
	logger.Warn("application force shutdown")
	os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
}
