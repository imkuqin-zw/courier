// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
