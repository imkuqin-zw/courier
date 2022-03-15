package defers

type Stage uint32

const (
	stageMin Stage = iota
	//StageBeforeStop before app stop
	StageBeforeStop
	//StageAfterStop after app stop
	StageAfterStop
	stageMax
)

var hooks = make(map[Stage]*Defer)

func RegisterBeforeHook(fns ...func() error) {
	hook, ok := hooks[StageBeforeStop]
	if !ok {
		hook = NewDefer()
		hooks[StageBeforeStop] = hook

	}
	hook.Register(fns...)
}

func RegisterAfterHook(fns ...func() error) {
	hook, ok := hooks[StageAfterStop]
	if !ok {
		hook = NewDefer()
		hooks[StageAfterStop] = hook

	}
	hook.Register(fns...)
}

func DoBeforeHook() {
	hook, ok := hooks[StageBeforeStop]
	if !ok {
		return

	}
	hook.Done()
}

func DoAfterHook() {
	hook, ok := hooks[StageAfterStop]
	if !ok {
		return

	}
	hook.Done()
}
