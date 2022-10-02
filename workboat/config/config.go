package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Logger.With().Stack().Logger()
}

var Debug = struct {
	Enabled bool
}{
	Enabled: asBool(get("debug.enable")),
}

var HTTP = struct {
	Host string
	Port int
}{
	Host: asString(withDefault("http.host", "0.0.0.0")),
	Port: asInt(withDefault("http.port", 8080)),
}

var Database = struct {
	Filename string
}{
	Filename: asString(withDefault("db.filename", "database.db")),
}
