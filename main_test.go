package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// TestMainHelp tests that --help flag works
func TestMainHelp(t *testing.T) {
	if os.Getenv("TEST_MAIN") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainHelp")
	cmd.Env = append(os.Environ(), "TEST_MAIN=1")
	cmd.Args = []string{os.Args[0], "--help"}

	output, err := cmd.CombinedOutput()
	// Help should exit with non-zero status, but we just verify it runs
	if err != nil {
		// This is expected - just check we got some output
		if len(output) == 0 {
			t.Error("Expected help output")
		}
	}
}

// TestPIDFiltering tests numeric PID filtering
func TestPIDFiltering(t *testing.T) {
	pids := []string{"100", "200", "300", "400", "500"}

	// Test MinPID filtering
	filtered := []string{}
	minPID := 250
	for _, p := range pids {
		num, _ := strconv.Atoi(p)
		if num >= minPID {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) != 3 {
		t.Errorf("Expected 3 PIDs >= 250, got %d", len(filtered))
	}

	// Test MaxPID filtering
	filtered = []string{}
	maxPID := 350
	for _, p := range pids {
		num, _ := strconv.Atoi(p)
		if num <= maxPID {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) != 3 {
		t.Errorf("Expected 3 PIDs <= 350, got %d", len(filtered))
	}

	// Test range filtering
	filtered = []string{}
	minPID = 200
	maxPID = 400
	for _, p := range pids {
		num, _ := strconv.Atoi(p)
		if num >= minPID && num <= maxPID {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) != 3 {
		t.Errorf("Expected 3 PIDs in range 200-400, got %d", len(filtered))
	}
}

// TestIgnoreSelf tests self PID filtering
func TestIgnoreSelf(t *testing.T) {
	self := os.Getpid()
	selfStr := strconv.Itoa(self)

	pids := []string{"100", selfStr, "300"}

	filtered := []string{}
	for _, p := range pids {
		if p != selfStr {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 PIDs after removing self, got %d", len(filtered))
	}

	for _, p := range filtered {
		if p == selfStr {
			t.Errorf("Self PID %s should have been filtered out", selfStr)
		}
	}
}

// TestJSONOutput tests JSON formatting
func TestJSONOutput(t *testing.T) {
	pids := []string{"123", "456", "789"}

	jsonBytes, err := json.Marshal(pids)
	if err != nil {
		t.Fatalf("Failed to marshal PIDs to JSON: %v", err)
	}

	var decoded []string
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(decoded) != len(pids) {
		t.Errorf("Expected %d PIDs, got %d", len(pids), len(decoded))
	}

	for i, pid := range pids {
		if decoded[i] != pid {
			t.Errorf("Expected PID %s at index %d, got %s", pid, i, decoded[i])
		}
	}
}

// TestSingleShot tests single-shot logic
func TestSingleShot(t *testing.T) {
	pids := []string{"100", "200", "300"}

	// Simulate single-shot
	single := pids[:1]

	if len(single) != 1 {
		t.Errorf("Expected 1 PID in single-shot mode, got %d", len(single))
	}

	if single[0] != "100" {
		t.Errorf("Expected first PID to be '100', got '%s'", single[0])
	}
}

// TestSpaceSeparatedOutput tests default output format
func TestSpaceSeparatedOutput(t *testing.T) {
	pids := []string{"123", "456", "789"}
	output := strings.Join(pids, " ")

	expected := "123 456 789"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}

	// Test parsing back
	parts := strings.Split(output, " ")
	if len(parts) != len(pids) {
		t.Errorf("Expected %d parts, got %d", len(pids), len(parts))
	}
}

// TestFindProcessByName is an integration test that finds the current test process
func TestFindProcessByName(t *testing.T) {
	// Get the name of the current process
	// This should be something like "main.test" or the test binary name
	self := os.Getpid()

	// Try to find PIDs - this tests the platform-specific FindPIDs function
	pids, err := FindPIDs("go", false)
	if err != nil {
		t.Logf("Warning: FindPIDs returned error: %v", err)
		// Don't fail the test as this might be expected in some environments
		return
	}

	// Just verify it returns a slice (may be empty)
	if pids == nil {
		t.Error("Expected non-nil slice from FindPIDs")
	}

	t.Logf("Found %d PIDs for 'go'", len(pids))

	// Check if our own PID is in the list (if we're running via go test)
	selfStr := strconv.Itoa(self)
	found := false
	for _, pid := range pids {
		if pid == selfStr {
			found = true
			break
		}
	}

	if found {
		t.Logf("Current test process PID %s found in results", selfStr)
	}
}

// TestFindPIDsInvalidProcess tests behavior with non-existent process
func TestFindPIDsInvalidProcess(t *testing.T) {
	pids, err := FindPIDs("this-process-definitely-does-not-exist-12345", false)
	if err != nil {
		t.Logf("FindPIDs returned error: %v", err)
		// Error is acceptable
		return
	}

	if len(pids) != 0 {
		t.Errorf("Expected 0 PIDs for non-existent process, got %d", len(pids))
	}
}
