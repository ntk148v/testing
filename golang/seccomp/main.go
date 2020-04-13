// Example in https://sysdig.com/blog/selinux-seccomp-falco-technical-discussion/
// Convert from C to Golang
package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	f, err := os.OpenFile("output.txt", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Println("Calling prctl() to send seccomp strict mode...")
	if err = unix.Prctl(unix.PR_SET_SECCOMP, unix.SECCOMP_MODE_STRICT, 0, 0, 0); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing to an already open file...")
	if _, err = f.WriteString("test"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Trying to open file for reading...")
	f, err = os.OpenFile("output.txt", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("You will not see this message. The process will be killed first")
}
