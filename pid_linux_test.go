//go:build linux

package main

import (
	"os"
	"strconv"
	"testing"
)

// TestFindPIDsLinux tests the Linux-specific FindPIDs implementation
func TestFindPIDsLinux(t *testing.T) {
	// Test that we can call FindPIDs without crashing
	pids, err := FindPIDs("systemd", false)
	if err != nil {
		t.Fatalf("FindPIDs failed: %v", err)
	}

	// systemd should exist on most modern Linux systems (PID 1)
	// If not systemd, might be init
	if len(pids) == 0 {
		// Try init instead
		pids, err = FindPIDs("init", false)
		if err != nil {
			t.Fatalf("FindPIDs failed: %v", err)
		}
	}

	if len(pids) == 0 {
		t.Error("Expected to find systemd or init process on Linux")
	}

	// Verify that PID 1 is in the results
	found := false
	for _, pid := range pids {
		if pid == "1" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find init system with PID 1")
	}
}

// TestFindPIDsSelfLinux tests finding the current test process
func TestFindPIDsSelfLinux(t *testing.T) {
	self := os.Getpid()

	// Try to find this test binary - it might be named "main.test" or similar
	pids, err := FindPIDs("main.test", false)
	if err != nil {
		t.Fatalf("FindPIDs failed: %v", err)
	}

	// Check if our PID is in the results
	selfStr := strconv.Itoa(self)
	found := false
	for _, pid := range pids {
		if pid == selfStr {
			found = true
			break
		}
	}

	if found {
		t.Logf("Successfully found test process with PID %s", selfStr)
	} else {
		t.Logf("Test process PID %s not found (this may be expected depending on the binary name)", selfStr)
	}
}

// TestFindPIDsExactModeLinux tests exact matching for scripts/interpreters
func TestFindPIDsExactModeLinux(t *testing.T) {
	// Test with exact mode enabled
	pids, err := FindPIDs("sh", true)
	if err != nil {
		t.Fatalf("FindPIDs with exact mode failed: %v", err)
	}

	// We might find shell processes
	t.Logf("Found %d 'sh' processes with exact mode", len(pids))

	// Test with exact mode disabled
	pidsNonExact, err := FindPIDs("sh", false)
	if err != nil {
		t.Fatalf("FindPIDs without exact mode failed: %v", err)
	}

	t.Logf("Found %d 'sh' processes without exact mode", len(pidsNonExact))
}

// TestFindPIDsEmptyResultLinux tests behavior with non-existent process
func TestFindPIDsEmptyResultLinux(t *testing.T) {
	pids, err := FindPIDs("nonexistent-process-xyz-123", false)
	if err != nil {
		t.Logf("FindPIDs returned error (acceptable): %v", err)
		return
	}

	if len(pids) != 0 {
		t.Errorf("Expected 0 PIDs for non-existent process, got %d", len(pids))
	}
}

// TestFindPIDsValidPIDsLinux tests that returned PIDs are valid numbers
func TestFindPIDsValidPIDsLinux(t *testing.T) {
	pids, err := FindPIDs("bash", false)
	if err != nil {
		t.Fatalf("FindPIDs failed: %v", err)
	}

	if len(pids) == 0 {
		t.Skip("No bash processes found, skipping validation")
	}

	for _, pid := range pids {
		num, err := strconv.Atoi(pid)
		if err != nil {
			t.Errorf("PID '%s' is not a valid number: %v", pid, err)
		}

		if num <= 0 {
			t.Errorf("PID %d should be positive", num)
		}
	}
}

// TestFindPIDsKernelThreadFiltering tests that kernel threads are filtered
func TestFindPIDsKernelThreadFiltering(t *testing.T) {
	// Kernel threads have names in brackets like [kthreadd]
	pids, err := FindPIDs("[kthreadd]", false)
	if err != nil {
		t.Fatalf("FindPIDs failed: %v", err)
	}

	// Should not find kernel threads as they're filtered out
	if len(pids) != 0 {
		t.Errorf("Expected 0 PIDs for kernel thread [kthreadd], got %d (kernel threads should be filtered)", len(pids))
	}
}

// TestFindPIDsProcFS tests /proc filesystem access
func TestFindPIDsProcFS(t *testing.T) {
	// Verify /proc exists
	if _, err := os.Stat("/proc"); os.IsNotExist(err) {
		t.Fatal("/proc filesystem does not exist")
	}

	// Verify we can read /proc
	entries, err := os.ReadDir("/proc")
	if err != nil {
		t.Fatalf("Cannot read /proc: %v", err)
	}

	// Count numeric entries (PIDs)
	pidCount := 0
	for _, e := range entries {
		if _, err := strconv.Atoi(e.Name()); err == nil {
			pidCount++
		}
	}

	if pidCount == 0 {
		t.Error("Expected to find at least one PID directory in /proc")
	}

	t.Logf("Found %d PID directories in /proc", pidCount)
}
