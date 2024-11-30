package middleware

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Issuer        string `required:"true" envconfig:"OAUTH_ISSUER"`
	RsaPublicKey  string `required:"true" envconfig:"OAUTH_RSA_PUBLIC_KEY"`
	SessionSecret string `required:"true" envconfig:"SESSION_SECRET"`
}

func NewConfig() *Config {
	conf := &Config{}
	err := envconfig.Process("middleware", conf)
	if err != nil {
		panic(err)
	}
	return conf
}
