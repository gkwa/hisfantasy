package hisfantasy

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	LogFormat string `long:"log-format" choice:"text" choice:"json" default:"text" description:"Log format"`
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose debug information, each -v bumps log level"`
	logLevel  slog.Level
	Dir       string `short:"d" long:"dir" default:"." description:"Directory to search for *.code-workspace files"`
}

func Execute() int {
	if err := parseFlags(); err != nil {
		return 1
	}

	if err := setLogLevel(); err != nil {
		return 1
	}

	if err := setupLogger(); err != nil {
		return 1
	}

	if err := run(); err != nil {
		slog.Error("run failed", "error", err)
		return 1
	}

	return 0
}

func parseFlags() error {
	_, err := flags.Parse(&opts)
	return err
}

func run() error {
	matches, err := filepath.Glob(filepath.Join(opts.Dir, "*.code-workspace"))
	if err != nil {
		return fmt.Errorf("glob failed: %w", err)
	}

	if len(matches) != 1 {
		return fmt.Errorf("found %d code-workspace files in %s", len(matches), opts.Dir)
	}

	workspacePath := matches[0]

	err = runCommand("code", workspacePath)
	if err != nil {
		return fmt.Errorf("code failed: %w", err)
	}

	return nil
}

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Debug("running command", "command", cmd.String())
	return cmd.Run()
}
