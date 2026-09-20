package appinfo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/Alia5/steaminputdb.com/steam/client"
	clientappinfo "github.com/Alia5/steaminputdb.com/steam/client/appinfo"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/types"
)

func fetchFromSteamClient(ctx context.Context, sc client.Client, appID uint32) (*clientappinfo.Info, error) {
	infos, err := clientappinfo.Get(ctx, sc, appID)
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, fmt.Errorf("no PICS info for app %d", appID)
	}
	return &infos[0], nil
}

func enrichModelFromPICS(appInfo *models.AppInfo, info *clientappinfo.Info) {
	cs := &models.AppControllerSupport{AppID: appInfo.AppID}
	switch info.Common.ControllerSupport {
	case "full":
		cs.SupportLevel = new(types.ControllerSupportLevelFull)
	case "partial":
		cs.SupportLevel = new(types.ControllerSupportLevelPartial)
	default:
		cs.SupportLevel = new(types.ControllerSupportLevelNone)
	}

	if appInfo.Name == "" && info.Common.Name != "" {
		appInfo.Name = info.Common.Name
	}

	cat := info.Common.Category
	check := func(id int64) *bool {
		_, ok := cat[fmt.Sprintf("category_%d", id)]
		return &ok
	}
	cs.DS4WiredSupport = check(int64(clientappinfo.CategoryPS4Wired))
	cs.DS4WirelessSupport = check(int64(clientappinfo.CategoryPS4Bluetooth))
	cs.DS5WiredSupport = check(int64(clientappinfo.CategoryPS5Wired))
	cs.DS5WirelessSupport = check(int64(clientappinfo.CategoryPS5Bluetooth))
	cs.SteamInputAPISupport = check(int64(clientappinfo.CategorySteamInputAPI))
	appInfo.ControllerSupport = cs

	appInfo.OfficialConfigs = make([]*models.OfficialSteamInputConfig, 0)
	type configCandidate struct {
		config    *models.OfficialSteamInputConfig
		isDefault bool
	}
	best := make(map[steamtypes.ControllerType]configCandidate)
	seen := make(map[steamtypes.ControllerType]bool)
	for idStr, detail := range info.Config.SteamControllerConfigDetails {
		cfg := parseConfig(appInfo.AppID, idStr, detail)
		if cfg == nil {
			continue
		}
		seen[cfg.ControllerType] = true
		isDefault := hasDefaultBranch(detail.EnabledBranches)
		if prev, ok := best[cfg.ControllerType]; !ok || (isDefault && !prev.isDefault) {
			best[cfg.ControllerType] = configCandidate{config: cfg, isDefault: isDefault}
		}
	}
	for idStr, detail := range info.Config.SteamControllerTouchConfigDetails {
		ct := steamtypes.ControllerType(detail.ControllerType)
		if seen[ct] {
			continue
		}
		cfg := parseConfig(appInfo.AppID, idStr, detail)
		if cfg == nil {
			continue
		}
		isDefault := hasDefaultBranch(detail.EnabledBranches)
		if prev, ok := best[cfg.ControllerType]; !ok || (isDefault && !prev.isDefault) {
			best[cfg.ControllerType] = configCandidate{config: cfg, isDefault: isDefault}
		}
	}
	for _, c := range best {
		appInfo.OfficialConfigs = append(appInfo.OfficialConfigs, c.config)
	}
}

func parseConfig(appID uint32, idStr string, detail clientappinfo.ControllerConfigDetail) *models.OfficialSteamInputConfig {
	configID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil
	}
	return &models.OfficialSteamInputConfig{
		AppID:          appID,
		ControllerType: steamtypes.ControllerType(detail.ControllerType),
		ConfigID:       configID,
	}
}

func hasDefaultBranch(enabledBranches string) bool {
	for _, b := range strings.Split(enabledBranches, ",") {
		if strings.TrimSpace(b) == "default" {
			return true
		}
	}
	return false
}
