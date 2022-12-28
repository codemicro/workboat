package paths

import (
	"fmt"
	"github.com/codemicro/workboat/workboat/config"
	"strings"
)

const (
	Auth = "/auth"

	API = "/api"

	APIAuth         = API + "/auth"
	APIAuthNewLogin = APIAuth + "/newLogin"

	AuthOauth         = Auth + "/oauth"
	AuthOauthOutbound = AuthOauth + "/outbound"
	AuthOauthInbound  = AuthOauth + "/inbound"

	Install              = API + "/install"
	InstallGetRepository = Install + "/getRepositories"
	InstallDoInstall     = Install + "/doInstall"

	WebhookInbound = API + "/inboundWebhook"
)

func JoinDomainAndPath(domain, path string) string {
	prepend := domain
	if !strings.HasPrefix(path, "/") {
		prepend += "/"
	}

	return prepend + path
}

func Make(path string, replacements ...any) string {
	x := strings.Split(path, "/")
	for i, p := range x {
		if strings.HasPrefix(p, ":") {
			x[i] = "%s"
		}
	}

	return JoinDomainAndPath(config.HTTP.ExternalURL, fmt.Sprintf(strings.Join(x, "/"), replacements...))
}
