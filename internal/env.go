package internal

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

func EnvOptions(prefix string, spec interface{}) error {
	prefix = strings.Replace(strings.ToUpper(prefix), "-", "_", -1)
	return envconfig.Process(prefix, spec)
}
