package helpers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

func RunPortForward(ctx context.Context, name string, args ...string) {
	cmd := exec.CommandContext(ctx, name, args...)
	// TODO: should output to slog instead of streaming directly to console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("Start port forwarding", "cmd", cmd.String())

	if err := cmd.Start(); err != nil {
		slog.Error("Error starting port forwarding", "err", err)
		return
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			slog.Error("Command stopped", "cmd", cmd.String(), "error", err)
		} else {
			slog.Info("Command finished", "cmd", cmd.String())
		}
	}()
}

func RunCommand(name string, args ...string) ([]byte, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		slog.Error(fmt.Sprintf("%s not found in PATH", name))
		return nil, err
	}
	slog.Info(fmt.Sprintf("Using %s at %s", name, path))

	cmd := exec.Command(name, args...)
	slog.Info(cmd.String())

	out, err := cmd.CombinedOutput()
	return out, err
}
