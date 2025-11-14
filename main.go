package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Single     bool `short:"s" long:"single-shot" description:"Return only one PID"`
	Kill       bool `short:"k" long:"kill" description:"Send SIGTERM to matched processes"`
	Exact      bool `short:"x" long:"exact" description:"Match exact command name including scripts"`
	IgnoreSelf bool `short:"e" long:"ignore-self" description:"Exclude pidof from results"`
	JSON       bool `short:"j" long:"json" description:"Output as JSON array"`
	MinPID     int  `long:"min-pid" default:"0" description:"Only include PIDs >= this value"`
	MaxPID     int  `long:"max-pid" default:"0" description:"Only include PIDs <= this value (0 = no limit)"`
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()

	if err != nil {
		os.Exit(2)
	}

	if len(args) != 1 {
		fmt.Println("Usage: pidof [OPTIONS] <process-name>")
		os.Exit(2)
	}

	target := args[0]

	pids, err := FindPIDs(target, opts.Exact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// numeric filtering
	filtered := []string{}
	for _, p := range pids {
		num, _ := strconv.Atoi(p)
		if opts.MinPID > 0 && num < opts.MinPID {
			continue
		}
		if opts.MaxPID > 0 && num > opts.MaxPID {
			continue
		}
		filtered = append(filtered, p)
	}
	pids = filtered

	// ignore self
	if opts.IgnoreSelf {
		self := os.Getpid()
		tmp := []string{}
		for _, p := range pids {
			if strconv.Itoa(self) != p {
				tmp = append(tmp, p)
			}
		}
		pids = tmp
	}

	if len(pids) == 0 {
		os.Exit(1)
	}

	sort.Strings(pids)

	// --single-shot
	if opts.Single {
		pids = pids[:1]
	}

	// --kill
	if opts.Kill {
		for _, p := range pids {
			pid, _ := strconv.Atoi(p)
			proc, err := os.FindProcess(pid)
			if err == nil {
				_ = proc.Signal(os.Interrupt)
			}
		}
	}

	// JSON output
	if opts.JSON {
		json.NewEncoder(os.Stdout).Encode(pids)
		return
	}

	// Default: space-separated PIDs
	fmt.Println(strings.Join(pids, " "))
}
