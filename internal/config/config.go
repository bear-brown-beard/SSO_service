package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv string `env:"APP_ENV" default:"dev"`
	Debug  bool   `env:"DEBUG" default:"true"`

	Mysql MysqlConfig

	RedisHost     string `env:"REDIS_HOST" required:"true"`
	RedisPort     string `env:"REDIS_PORT" required:"true"`
	RedisPassword string `env:"REDIS_PASSWORD" required:"false"`
	RedisDB       int    `env:"REDIS_DB" required:"true"`

	ServerPort string `env:"SERVER_PORT" required:"true"`

	CORS CORSConfig

	JWT JWTConfig

	SMSCLogin    string `env:"SMSC_LOGIN" required:"true"`
	SMSCPassword string `env:"SMSC_PASSWORD" required:"true"`

	Mindbox MindboxConfig

	TLS TLSConfig
}

type MysqlConfig struct {
	DBName            string        `env:"DB_NAME" required:"true"`
	DBUser            string        `env:"DB_USER" required:"true"`
	DBPassword        string        `env:"DB_PASSWORD" required:"true"`
	DBHost            string        `env:"DB_HOST" required:"true"`
	DBPort            string        `env:"DB_PORT" required:"true"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" default:"25"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" default:"25"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET_KEY" required:"true"`
}

type CORSConfig struct {
	AllowedOrigins   []string `env:"CORS_ALLOWED_ORIGINS" env-default:"http://localhost:3000,http://localhost:8080"`
	AllowedMethods   []string `env:"CORS_ALLOWED_METHODS" env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders   []string `env:"CORS_ALLOWED_HEADERS" env-default:"Accept,Authorization,Content-Type,X-DeviceUUID,X-Platform"`
	ExposedHeaders   []string `env:"CORS_EXPOSED_HEADERS" env-default:"Link"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" env-default:"true"`
	MaxAge           int      `env:"CORS_MAX_AGE" env-default:"300"`
}

type TLSConfig struct {
	SkipVerify bool          `env:"TLS_SKIP_VERIFY" env-default:"false"`
	Timeout    time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`
}

func Load() *Config {
	cfg := &Config{}
	path := "./.env"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fromEnv(cfg)
	}
	cfg = fromFile(path, cfg)
	return cfg
}

func fromFile(path string, cfg *Config) *Config {
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic(err)
	}
	return cfg
}

func fromEnv(cfg *Config) *Config {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(err)
	}
	return cfg
}
