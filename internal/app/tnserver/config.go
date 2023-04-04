package tnserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	DatabaseURL string `toml:"database_url"`
	SessionKey string `toml:"session_key"`
	KafkaAddr string `toml:"kafka_addr"`
}

func NewConfig() *Config {
	cfg := &Config{
		BindAddr: ":8080",
	}
	return cfg
}