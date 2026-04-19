package steam

import (
	"log/slog"
	"net/url"
	"strings"

	"github.com/Alia5/steaminputdb.com/api/steam/appinfo"
	"github.com/Alia5/steaminputdb.com/api/steam/filedetails"
	"github.com/Alia5/steaminputdb.com/api/steam/login"
	"github.com/Alia5/steaminputdb.com/api/steam/user"
	"github.com/Alia5/steaminputdb.com/db"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(a huma.API, dal db.DAL) {
	login.RegisterRoutes(a)
	user.RegisterRoutes(a)
	appinfo.RegisterRoute(a, dal)
	filedetails.RegisterRoute(a)
}

func OpenIDAuthorizationURL(callbackAddress string) string {
	scheme, rest, _ := strings.Cut(callbackAddress, "://")
	host, _, _ := strings.Cut(rest, "/")

	realm := scheme + "://" + host
	slog.Debug("realm", "realm", realm)

	params := url.Values{}
	params.Add("openid.mode", "checkid_setup")
	params.Add("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Add("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Add("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Add("openid.return_to", callbackAddress)
	params.Add("openid.realm", realm)

	return "https://steamcommunity.com/openid/login?" + params.Encode()
}
