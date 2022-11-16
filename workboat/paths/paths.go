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

	Install                 = "/install"
	InstallSelectRepository = "/install/selectRepository"
	InstallDoInstall        = "/install/doInstall"
)

func Make(path string, replacements ...any) string {
	x := strings.Split(path, "/")
	for i, p := range x {
		if strings.HasPrefix(p, ":") {
			x[i] = "%s"
		}
	}

	prepend := config.HTTP.ExternalURL
	if !strings.HasPrefix(path, "/") {
		prepend += "/"
	}

	return prepend + fmt.Sprintf(strings.Join(x, "/"), replacements...)
}
