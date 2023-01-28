package config

import (
	"flag"
	"os"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type Config struct {
	Database string `toml:"database"`
	HostHTTP string `toml:"http_host"`
	PortHTTP string `toml:"http_port"`
	HostGRPC string `toml:"grpc_host"`
	PortGRPC string `toml:"grpc_port"`
	PostgresConnectionString string `toml:"postgres_connection_string"`
}

func LoadConfig() (*Config, error) {
	var configFile string
	flag.StringVar(&configFile, "conf", "config.toml", "toml file with configs")
	flag.Parse()
	
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	decoder := toml.NewDecoder(f).Strict(true)
	if decoder == nil {
		return nil, errors.New("couldn't create decoder")
	}

	conf := &Config{}
	if err := decoder.Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

