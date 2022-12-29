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
	Host string
	Port int
}{
	Host: asString(withDefault("http.host", "0.0.0.0")),
	Port: asInt(withDefault("http.port", 8080)),
}

var Gitea = struct {
	BaseURL     string
	AccessToken string
}{
	BaseURL:     strings.TrimSuffix(asString(required("gitea.baseURL")), "/"),
	AccessToken: asString(required("gitea.accessToken")),
}
