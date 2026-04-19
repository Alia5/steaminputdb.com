package appinfo

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alia5/steaminputdb.com/steam/client"
	"github.com/Alia5/steaminputdb.com/steam/vdf"
)

type infoVDF struct {
	AppInfo Info `json:"appinfo"`
}

func Get(ctx context.Context, c client.Client, appIDs ...uint32) ([]Info, error) {
	apps := make([]*client.CMsgClientPICSProductInfoRequest_AppInfo, len(appIDs))
	for i := range appIDs {
		apps[i] = &client.CMsgClientPICSProductInfoRequest_AppInfo{Appid: &appIDs[i]}
	}
	var resp client.CMsgClientPICSProductInfoResponse
	err := c.SendMessage(ctx,
		client.EMsg_k_EMsgClientPICSProductInfoRequest,
		client.EMsg_k_EMsgClientPICSProductInfoResponse,
		&client.CMsgClientPICSProductInfoRequest{Apps: apps},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	result := make([]Info, 0, len(resp.GetApps()))
	for _, app := range resp.GetApps() {
		var wrapper infoVDF
		err := vdf.Unmarshal(strings.TrimRight(string(app.GetBuffer()), "\x00"), &wrapper)
		if err != nil {
			return nil, fmt.Errorf("parse VDF for app %d: %w", app.GetAppid(), err)
		}
		wrapper.AppInfo.AppID = app.GetAppid()
		result = append(result, wrapper.AppInfo)
	}
	return result, nil
}
