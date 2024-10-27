package oauth

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RsaPrivateKey           string `required:"true" split_words:"true"`
	RsaPrivateKeyPassphrase string `required:"true" split_words:"true"`
	Issuer                  string `required:"true"`
}

func NewConfig() *Config {
	conf := &Config{}
	err := envconfig.Process("oauth", conf)
	if err != nil {
		panic(err)
	}
	return conf
}
