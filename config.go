package booksdb

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigtoml"
)

type Config struct {
	Port       int    `default:"1111" usage:"just give a number"`
	MongoDBURL string `default:"mongodb://localhost:27017" usage:"pass a URI"`
}

func GetConfig() Config {
	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{ //nolint:exhaustivestruct
		SkipDefaults: true,
		SkipFiles:    true,
		SkipEnv:      true,
		SkipFlags:    true,
		EnvPrefix:    "APP",
		FlagPrefix:   "app",
		Files:        []string{"/var/opt/booksdb/config.toml", "config.toml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".toml": aconfigtoml.New(),
		},
	})

	// flagSet := loader.Flags()

	if err := loader.Load(); err != nil {
		panic(err)
	}

	return cfg
}
