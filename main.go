package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type RunResult struct {
	start  time.Time
	end    time.Time
	stdout string
	stderr string
	err    error
}

func nextBoundary(now time.Time, interval time.Duration) (truncated time.Time) {
	truncated = now.Truncate(interval)
	if truncated.Equal(now) {
		return now
	}
	truncated = truncated.Add(interval)
	return
}

func formatCountdown(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	secs := int(d.Seconds())
	h := secs / 3600
	m := (secs % 3600) / 60
	s := secs % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func awaitWithCountdown(ctx context.Context, until time.Time) bool {
	ticket := time.NewTicker(1 * time.Second)
	defer ticket.Stop()

	for {
		remaining := time.Until(until)
		if remaining <= 0 {
			return true
		}

		fmt.Printf("\rNext run at %s (in %s)", until.Format(time.Kitchen), formatCountdown(remaining))
		select {
		case <-ticket.C:
		case <-ctx.Done():
			fmt.Println("\n\nShutting down gracefully.")
			return false
		}
	}
}

func buildShellCommand(cmdString string) (name string, args []string) {
	if runtime.GOOS == "windows" {
		return "cmd.exe", []string{"/C", cmdString}
	}

	return "sh", []string{"-c", cmdString}
}

func runCommand(ctx context.Context, cmdString string) RunResult {
	start := time.Now()

	name, args := buildShellCommand(cmdString)
	cmd := exec.CommandContext(ctx, name, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	end := time.Now()

	return RunResult{
		start:  start,
		end:    end,
		stdout: outBuf.String(),
		stderr: errBuf.String(),
		err:    err,
	}
}

func init() {
	fmt.Println("\nScheduler v0.0.0")
}

func main() {
	var cmdString string
	var intervalMin int
	var stdOut bool

	flag.StringVar(&cmdString, "cmd", "", "Command to execute (run through OS shell). e.g., \"nxp your-command --flag\"")
	flag.IntVar(&intervalMin, "interval-min", 0, "Interval in minutes for all-clock aligned execution")
	flag.BoolVar(&stdOut, "stdout", false, "Print stdout from command execution")
	flag.Parse()

	cmdString = strings.TrimSpace(cmdString)
	if cmdString == "" {
		fmt.Fprintf(os.Stderr, "Error: -cmd is required")
		flag.Usage()
		os.Exit(2)
	}

	if intervalMin <= 0 {
		fmt.Fprintln(os.Stderr, "Error: -interval-min must be > 0")
		flag.Usage()
		os.Exit(2)
	}

	interval := time.Duration(intervalMin) * time.Minute

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Interval: %s, Command: %s\n", interval, cmdString)

	for {
		now := time.Now()
		next := nextBoundary(now, interval)
		if !awaitWithCountdown(ctx, next) {
			return
		}

		res := runCommand(ctx, cmdString)

		start := res.start.Format(time.Kitchen)
		end := res.end.Format(time.Kitchen)
		duration := res.end.Sub(res.start).Round(time.Millisecond)

		fmt.Printf("\r[started: %s] [ended: %s][duration: %s]\n", start, end, duration)

		if stdOut {
			fmt.Printf("STDOUT:\n%s", res.stdout)
		}

		if res.stderr != "" {
			fmt.Printf("STDERR:\n%s", res.stderr)
		}

		if res.err != nil {
			fmt.Printf("Command failed (will continue execution):\n%+v", res.err)
		}
	}
}
