package configs

func InitializeConfig() *Config {
	conf, err := NewConfig()
	if err != nil {
		panic("error loading application config")
	}

	return conf
}
