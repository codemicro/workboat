package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"strings"
)

func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Logger.With().Stack().Logger()
}

var Debug = struct {
	Enabled bool
}{
	Enabled: asBool(get("debug.enabled")),
}

var HTTP = struct {
	Host        string
	Port        int
	ExternalURL string
}{
	Host:        asString(withDefault("http.host", "0.0.0.0")),
	Port:        asInt(withDefault("http.port", 8080)),
	ExternalURL: strings.TrimSuffix(asString(required("http.externalURL")), "/"),
}

var Database = struct {
	Filename string
}{
	Filename: asString(withDefault("db.filename", "database.db")),
}

var Gitea = struct {
	BaseURL           string
	OauthClientID     string
	OauthClientSecret string
	AccessToken       string
}{
	BaseURL:           strings.TrimSuffix(asString(required("gitea.baseURL")), "/"),
	OauthClientID:     asString(required("gitea.oauth.clientID")),
	OauthClientSecret: asString(required("gitea.oauth.clientSecret")),
	AccessToken:       asString(required("gitea.accessToken")),
}
