package config

type Database struct {
	Name string
	User string
	Pass string
}

type Jwt struct {
	SecretKey []byte
}

type Config struct {
	Database Database
	Jwt      Jwt
}

func LoadConfigs() Config {
	return Config{Database{"blockcoin", "blockcoin", "bl0ckc01n"}, Jwt{[]byte("414ca94841798c6a0fed5fe3a959f8f3")}}
}
