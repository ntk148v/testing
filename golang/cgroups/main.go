package main

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/unix"
)

const unifiedMountpoint = "/sys/fs/cgroup"

// CGMode is the cgroups mode of the host system
type CGMode int

const (
	// Unavailable cgroup mountpoint
	Unavailable CGMode = iota
	// Legacy cgroups v1
	Legacy
	// Hybrid with cgroups v1 and v2 controllers mounted
	Hybrid
	// Unified with only cgroups v2 mounted
	Unified
)

func (c CGMode) String() string {
	switch c {
	case Legacy:
		return "Legacy"
	case Hybrid:
		return "Hybrid"
	case Unified:
		return "Unified"
	default:
		return "Unavailable"
	}
}

// Mode returns the cgroups mode running on the host
func Mode() (CGMode, error) {
	var (
		st     unix.Statfs_t
		cgMode CGMode
	)
	if err := unix.Statfs(unifiedMountpoint, &st); err != nil {
		return Unavailable, err
	}
	switch st.Type {
	case unix.CGROUP2_SUPER_MAGIC:
		cgMode = Unified
	default:
		cgMode = Legacy
		if err := unix.Statfs(filepath.Join(unifiedMountpoint, "unified"), &st); err != nil {
			return Unavailable, err
		}
		if st.Type == unix.CGROUP2_SUPER_MAGIC {
			cgMode = Hybrid
		}
	}
	return cgMode, nil
}

func main() {
	cgMode, err := Mode()
	if err != nil {
		panic(err)
	}
	fmt.Println("CGroup version:", cgMode)
}
