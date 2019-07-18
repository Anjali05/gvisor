// Copyright 2019 The gVisor Authors.
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

// +build linux

package flipcall

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"syscall"

	"gvisor.dev/gvisor/pkg/abi/linux"
)

func (ep *Endpoint) futexSwitchToPeer() error {
	// Update connection state to indicate that the peer should be active.
	if !atomic.CompareAndSwapUint32(ep.connState(), ep.activeState, ep.inactiveState) {
		return fmt.Errorf("unexpected connection state before FUTEX_WAKE: %v", atomic.LoadUint32(ep.connState()))
	}

	// Wake the peer's Endpoint.futexSwitchFromPeer().
	if err := ep.futexWakeConnState(1); err != nil {
		return fmt.Errorf("failed to FUTEX_WAKE peer Endpoint: %v", err)
	}
	return nil
}

func (ep *Endpoint) futexSwitchFromPeer() error {
	for {
		switch cs := atomic.LoadUint32(ep.connState()); cs {
		case ep.activeState:
			return nil
		case ep.inactiveState:
			// Continue to FUTEX_WAIT.
		default:
			return fmt.Errorf("unexpected connection state before FUTEX_WAIT: %v", cs)
		}
		if ep.isShutdownLocally() {
			return shutdownError{}
		}
		if err := ep.futexWaitConnState(ep.inactiveState); err != nil {
			return fmt.Errorf("failed to FUTEX_WAIT for peer Endpoint: %v", err)
		}
	}
}

func (ep *Endpoint) futexWakeConnState(numThreads int32) error {
	if _, _, e := syscall.RawSyscall(syscall.SYS_FUTEX, ep.packet, linux.FUTEX_WAKE, uintptr(numThreads)); e != 0 {
		return e
	}
	return nil
}

func (ep *Endpoint) futexWaitConnState(curState uint32) error {
	_, _, e := syscall.Syscall6(syscall.SYS_FUTEX, ep.packet, linux.FUTEX_WAIT, uintptr(curState), 0, 0, 0)
	if e != 0 && e != syscall.EAGAIN && e != syscall.EINTR {
		return e
	}
	return nil
}

func yieldThread() {
	syscall.Syscall(syscall.SYS_SCHED_YIELD, 0, 0, 0)
	// The thread we're trying to yield to may be waiting for a Go runtime P.
	// runtime.Gosched() will hand off ours if necessary.
	runtime.Gosched()
}
