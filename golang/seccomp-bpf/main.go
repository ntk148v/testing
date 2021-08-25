// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Implements seccomp filters using a whitelist approach.
package main

import (
	"fmt"
	"syscall"

	libseccomp "github.com/seccomp/libseccomp-golang"
)

// whitelist - syscalls contains the names of the syscalls
// that we want to our process to have access to.
func whiteList(syscalls []string) {
	// Apply a "deny all" filter
	filter, err := libseccomp.NewFilter(libseccomp.ActErrno.SetReturnCode(int16(syscall.EPERM)))
	if err != nil {
		fmt.Printf("Error creating filter: %s\n", err)
	}
	// Add elements to whitelist
	for _, element := range syscalls {
		fmt.Printf("[+] Whitelisting: %s\n", element)
		syscallID, err := libseccomp.GetSyscallFromName(element)
		if err != nil {
			panic(err)
		}
		filter.AddRule(syscallID, libseccomp.ActAllow)
	}
	// Load the filter which applies the filter we just created
	filter.Load()
}

func main() {
	// A string array contains the names of syscalls extracted
	// from `strace` output.
	var syscalls = []string{
		"write", "mmap", "rt_sigaction", "rt_sigprocmask",
		"clone", "execve", "fcntl", "sigaltstack", "arch_prctl",
		"gettid", "gettid", "futex", "sched_getaffinity",
		"mkdirat", "readlinkat", "exit_group",
	}

	whiteList(syscalls)
	err := syscall.Mkdir("/tmp/moo", 0755)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("I just created a folder")
	}
	// Attempt to execute a shell command
	// Force remove the /tmp/moo first
	err = syscall.Rmdir("/tmp/moo")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("I just removed a folder")
	}
}
