//go:build windows

package steam

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

func getCefDebugPortWindows(ctx context.Context) uint16 {
	const defaultPort uint16 = 8080
	const prefix = "--remote-debugging-port="

	procs, err := process.ProcessesWithContext(ctx)
	if err != nil {
		slog.Error("Failed to enumerate processes", "error", err)
		return defaultPort
	}

	for _, p := range procs {
		name, err := p.NameWithContext(ctx)
		if err != nil {
			continue
		}
		if !strings.EqualFold(name, "steamwebhelper.exe") {
			continue
		}

		args, err := p.CmdlineSliceWithContext(ctx)
		if err != nil || len(args) == 0 {
			cmdline, cmdlineErr := p.CmdlineWithContext(ctx)
			if cmdlineErr != nil || cmdline == "" {
				continue
			}
			args = strings.Fields(cmdline)
		}

		for _, arg := range args {
			if !strings.HasPrefix(arg, prefix) {
				continue
			}

			portStr := strings.TrimPrefix(arg, prefix)
			port, parseErr := strconv.Atoi(portStr)
			if parseErr != nil || port <= 0 || port > 65535 {
				continue
			}

			return uint16(port)
		}
	}

	return defaultPort
}
