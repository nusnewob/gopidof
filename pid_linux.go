//go:build linux

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func FindPIDs(target string, exact bool) ([]string, error) {
	var pids []string
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		pid := e.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		commBytes, err := os.ReadFile(filepath.Join("/proc", pid, "comm"))
		if err != nil {
			continue
		}
		comm := strings.TrimSpace(string(commBytes))

		// skip kernel threads
		if strings.HasPrefix(comm, "[") && strings.HasSuffix(comm, "]") {
			continue
		}

		// normal match
		if comm == target {
			pids = append(pids, pid)
			continue
		}

		if exact {
			// script/interpreter match via cmdline
			cmdline, err := os.ReadFile(filepath.Join("/proc", pid, "cmdline"))
			if err == nil && len(cmdline) > 0 {
				parts := bytes.Split(cmdline, []byte{0})
				for _, a := range parts {
					if filepath.Base(string(a)) == target {
						pids = append(pids, pid)
						break
					}
				}
			}
		}
	}

	return pids, nil
}
