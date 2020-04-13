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
