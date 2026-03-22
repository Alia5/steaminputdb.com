package steam

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func steamPath() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open Steam registry key: %w", err)
	}
	defer k.Close()

	p, _, err := k.GetStringValue("SteamPath")
	if err != nil {
		return "", fmt.Errorf("failed to read SteamPath from registry: %w", err)
	}

	return p, nil
}
