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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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
	// Check here: https://en.wikipedia.org/wiki/Linux_namespaces#Namespace_kinds
	// CLONE_NEWUTS namespace isolates hostname
	// CLONE_NEWPID namespace isolates processes
	// CLONE_NEWNS namespace isolates mounts
	// CLONE_NEWIPC namespace isolates interprocess communication (IPC)
	// CLONE_NEWNET namespace isolates network
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	// Run child using namespaces. The command provided will be executed inside that.
	checkErr(cmd.Run())
}

func child() {
	fmt.Printf("running child %v as PID %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create cgroup
	cg()

	// Get homedir
	dir := os.TempDir()
	rootFsDir := path.Join(dir, "rootfs")

	// Change root filesystem
	// Have to create a rootfs by following this, don't know why
	// others not face the same problem.
	// https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
	checkErr(os.MkdirAll(rootFsDir, 0700))
	checkErr(syscall.Sethostname([]byte("kontainer")))
	checkErr(syscall.Mount("/", rootFsDir, "", syscall.MS_BIND, ""))
	// oldRootFsDir := path.Join(rootFsDir, "oldrootfs")
	// checkErr(os.MkdirAll(oldRootFsDir, 0700))
	// checkErr(syscall.PivotRoot("/", oldRootFsDir))
	// checkErr(os.Chdir("/"))

	checkErr(syscall.Chroot(rootFsDir))
	checkErr(os.Chdir("/"))
	checkErr(syscall.Mount("proc", "proc", "proc", 0, ""))

	checkErr(cmd.Run())

	// Cleanup mount
	checkErr(syscall.Unmount("proc", 0))
}

// cg - demonstrate cgroup
func cg() {
	// cgroup location in Ubuntu
	cgroups := "/sys/fs/cgroup"
	mem := filepath.Join(cgroups, "memory")
	kontainerMem := filepath.Join(mem, "kontainer")
	_ = os.Mkdir(kontainerMem, 0755)
	// Limit memory to 1mb
	checkErr(ioutil.WriteFile(filepath.Join(kontainerMem, "memory.limit_in_bytes"), []byte("999424"), 0700))

	pids := filepath.Join(cgroups, "pids")
	kontainerPids := filepath.Join(pids, "kontainer")
	_ = os.Mkdir(kontainerPids, 0755)
	checkErr(ioutil.WriteFile(filepath.Join(kontainerPids, "pids.max"), []byte("20"), 0700))
	// Cleanup cgroup when it is not being used (after the container exits)
	checkErr(ioutil.WriteFile(filepath.Join(kontainerPids, "notify_on_release"), []byte("1"), 0700))
	// Apply this and any child process in this cgroup
	checkErr(ioutil.WriteFile(filepath.Join(kontainerPids, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

// checkErr - simply catch error and panic
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
