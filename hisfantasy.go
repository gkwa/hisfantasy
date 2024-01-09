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
	DryRun    bool     `short:"n" long:"dry-run" description:"Dry run: print the command instead of executing it"`
	Dirs      []string `short:"d" long:"dir" default:"." description:"Directories to search for *.code-workspace files"`
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

func getWorkspacePathForDir(dir string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.code-workspace"))
	if err != nil {
		return "", fmt.Errorf("glob failed: %w", err)
	}

	if len(matches) != 1 {
		return "", fmt.Errorf("found %d code-workspace files in %s", len(matches), dir)
	}

	return matches[0], nil
}

func runCommandForDirs(dryRun bool, dirs ...string) error {
	cmd, err := buildCommand("code", dirs...)
	if err != nil {
		return fmt.Errorf("buildCommand failed: %w", err)
	}

	if dryRun {
		slog.Debug("command", "dry-run", true, "command", cmd.String())
	} else {
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("code failed: %w", err)
		}
	}

	return nil
}

func run() error {
	slog.Debug("running", "directory list requested", opts.Dirs)

	err := runCommandForDirs(opts.DryRun, opts.Dirs...)
	if err != nil {
		return fmt.Errorf("runCommandForDir failed: %w", err)
	}

	return nil
}

func buildCommand(command string, dirs ...string) (*exec.Cmd, error) {
	workspacePaths := make([]string, len(dirs))
	for _, dir := range dirs {
		workspacePath, err := getWorkspacePathForDir(dir)
		if err != nil {
			return nil, fmt.Errorf("getWorkspacePathForDir failed: %w", err)
		}
		workspacePaths = append(workspacePaths, workspacePath)
	}

	cmd := exec.Command(command, workspacePaths...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
