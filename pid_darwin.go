//go:build darwin

package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/unix"
)

func FindPIDs(target string, exact bool) ([]string, error) {
	kprocs, err := unix.SysctlKinfoProcSlice("kern.proc.all")
	if err != nil {
		return nil, fmt.Errorf("sysctl: %w", err)
	}

	var pids []string

	for _, k := range kprocs {
		pid := int(k.Proc.P_pid)
		if pid <= 0 {
			continue
		}

		raw, err := unix.SysctlRaw("kern.procargs2", pid)
		if err != nil || len(raw) < 4 {
			continue
		}

		argv := raw[4:]
		parts := bytes.Split(argv, []byte{0})
		if len(parts) < 1 {
			continue
		}

		exe := filepath.Base(string(parts[0]))

		if exe == target {
			pids = append(pids, strconv.Itoa(pid))
			continue
		}

		if exact {
			for _, a := range parts[1:] {
				if filepath.Base(string(a)) == target {
					pids = append(pids, strconv.Itoa(pid))
					break
				}
			}
		}
	}

	return pids, nil
}
