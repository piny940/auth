package oauth

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RsaPrivateKey           string `required:"true" split_words:"true"`
	RsaPrivateKeyPassPhrase string `required:"true" split_words:"true"`
	Issuer                  string `required:"true"`
}

var config = &Config{}

func Init() {
	err := envconfig.Process("oauth", config)
	if err != nil {
		panic(err)
	}
}
