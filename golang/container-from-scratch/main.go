package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker run <container> cmd args
// go run main.go run cmd args
func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("Unknown command. Use run <command_name>, like `run /bin/bash` or `run echo hello`")
	}
}

func parent() {
	fmt.Printf("running parent %v as PID %d\n", os.Args[2:], os.Getpid())

	// /proc/self/exe - a special file containing an in-memory image of the current executable.
	// In other words, we re-run ourselves, but passing childs as the first agrument.
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Cloneflags is only available in Linux
	// CLONE_NEWUTS namespace isolates hostname
	// CLONE_NEWPID namespace isolates processes
	// CLONE_NEWNS namespace isolates mounts
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	// Run child using namespaces. The command provided will be executed inside that.
	must(cmd.Run())
}

func child() {
	fmt.Printf("running child %v as PID %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create cgroup
	cg()

	// Change root filesystem
	// Have to create a rootfs by following this, don't know why
	// others not face the same problem.
	// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/root/containerfs"))
	must(os.Chdir("/"))
	// must(syscall.Mount("proc", "proc", "proc", 0, ""))
	// must(syscall.Mount("something", "mytemp", "tmpfs", 0, ""))

	must(cmd.Run())

	// Cleanup mount
	// must(syscall.Unmount("proc", 0))
	// must(syscall.Unmount("mytemp", 0))
}

// cg - demonstrate cgroup
func cg() {
	// cgourp location in Ubuntu
	cgroups := "/sys/fs/cgroup"
	mem := filepath.Join(cgroups, "memory")
	kontainer := filepath.Join(mem, "kontainer")
	_ = os.Mkdir(kontainer, 0755)
	// Limit memory to 1mb
	must(ioutil.WriteFile(filepath.Join(kontainer, "memory.limit_in_bytes"), []byte("999424"), 0700))
	// Cleanup cgroup when it is not being used
	must(ioutil.WriteFile(filepath.Join(kontainer, "notify_on_release"), []byte("1"), 0700))

	pid := strconv.Itoa(os.Getpid())
	// Apply this and any child process in this cgroup
	must(ioutil.WriteFile(filepath.Join(kontainer, "cgroup.procs"), []byte(pid), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
