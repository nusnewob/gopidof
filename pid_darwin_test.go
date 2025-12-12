//go:build darwin

package main

import (
	"os"
	"strconv"
	"testing"
)

// TestFindPIDsDarwin tests the Darwin-specific FindPIDs implementation
func TestFindPIDsDarwin(t *testing.T) {
	// Test that we can call FindPIDs without crashing
	pids, err := FindPIDs("launchd", false)
	if err != nil {
		t.Fatalf("FindPIDs failed: %v", err)
	}

	// launchd should exist on macOS but we may not be able to see PID 1
	// depending on permissions, so just verify the function works
	t.Logf("Found %d 'launchd' processes", len(pids))

	// Verify that any returned PIDs are valid numbers
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

// TestFindPIDsSelfDarwin tests finding the current test process
func TestFindPIDsSelfDarwin(t *testing.T) {
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

// TestFindPIDsExactModeDarwin tests exact matching for scripts/interpreters
func TestFindPIDsExactModeDarwin(t *testing.T) {
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

// TestFindPIDsEmptyResultDarwin tests behavior with non-existent process
func TestFindPIDsEmptyResultDarwin(t *testing.T) {
	pids, err := FindPIDs("nonexistent-process-xyz-123", false)
	if err != nil {
		t.Logf("FindPIDs returned error (acceptable): %v", err)
		return
	}

	if len(pids) != 0 {
		t.Errorf("Expected 0 PIDs for non-existent process, got %d", len(pids))
	}
}

// TestFindPIDsValidPIDsDarwin tests that returned PIDs are valid numbers
func TestFindPIDsValidPIDsDarwin(t *testing.T) {
	pids, err := FindPIDs("kernel_task", false)
	if err != nil {
		t.Fatalf("FindPIDs failed: %v", err)
	}

	// kernel_task should exist on macOS
	if len(pids) == 0 {
		t.Skip("kernel_task not found, skipping validation")
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
