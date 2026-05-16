//go:build !windows

package steam

import "context"

func getCefDebugPortWindows(_ context.Context) uint16 {
	return 8080
}
