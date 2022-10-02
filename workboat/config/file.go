package config

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const configFileName = "config.yml"

var (
	rawConfigFileContents map[string]any
	lastKey               string
)

func mustLoadConfigFile() {
	if err := loadConfigFile(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func loadConfigFile() error {
	if rawConfigFileContents != nil {
		return nil
	}

	fcont, err := os.ReadFile(configFileName)
	if err != nil {
		return errors.Wrap(err, "failed to load config file")
	}
	rawConfigFileContents = make(map[string]any)
	if err := yaml.Unmarshal(fcont, &rawConfigFileContents); err != nil {
		return errors.Wrap(err, "could not unmarshal config file")
	}

	return nil
}

func Reload() error {
	return loadConfigFile()
}

type optionalItem struct {
	item  any
	found bool
}

var indexedPartRegexp = regexp.MustCompile(`(?m)([a-zA-Z]+)(?:\[(\d+)\])?`)

func get(key string) optionalItem {
	// http[2].bananas
	mustLoadConfigFile()
	lastKey = key

	parts := strings.Split(key, ".")
	var cursor any = rawConfigFileContents
	for _, part := range parts {
		components := indexedPartRegexp.FindStringSubmatch(part)
		key := components[1]
		index, _ := strconv.ParseInt(components[2], 10, 32)
		isIndexed := components[2] != ""

		item, found := cursor.(map[string]any)[key]
		if !found {
			return optionalItem{nil, false}
		}

		if isIndexed {
			arr, conversionOk := item.([]any)
			if !conversionOk {
				log.Fatal().Msgf("attempted to index non-indexable config item %s", key)
			}
			cursor = arr[index]
		} else {
			cursor = item
		}
	}
	return optionalItem{cursor, true}
}

func required(key string) optionalItem {
	opt := get(key)
	if !opt.found {
		log.Fatal().Msgf("required key %s not found in config file", lastKey)
	}
	return opt
}

func withDefault(key string, defaultValue any) optionalItem {
	opt := get(key)
	if !opt.found {
		return optionalItem{item: defaultValue, found: true}
	}
	return opt
}

func asInt(x optionalItem) int {
	if !x.found {
		return 0
	}
	return x.item.(int)
}

func asString(x optionalItem) string {
	if !x.found {
		return ""
	}
	return x.item.(string)
}

func asBool(x optionalItem) bool {
	if !x.found {
		return false
	}
	return x.item.(bool)
}
