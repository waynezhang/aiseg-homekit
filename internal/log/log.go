package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetVerbose(verbose bool) {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func E(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func D(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func F(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func Dump(v interface{}) {
	D("Dump object -> %s", v)
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
