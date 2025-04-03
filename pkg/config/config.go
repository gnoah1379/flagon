package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

var config Config

type Config struct {
	Log      Log
	Server   Server
	Database Database
	Auth     Authentication
	Cache    Cache
}

func (conf Config) Validate() error {
	return nil
}

func GetConfig() Config {
	return config
}

func LoadConfig(filePath string) error {
	viper.SetEnvPrefix("FLAGON")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	if filePath != "" {
		viper.SetConfigFile(filePath)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	for _, key := range viper.AllKeys() {
		viper.SetDefault(key, viper.Get(key))
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
		return err
	}
	return nil
}

func init() {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.file", "stderr")
	viper.SetDefault("log.addSource", true)

	viper.SetDefault("http.host", "0.0.0.0")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.enableTLS", false)
	viper.SetDefault("http.certFile", "")
	viper.SetDefault("http.keyFile", "")

	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.host", "")
	viper.SetDefault("database.port", "")
	viper.SetDefault("database.username", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "flagon.sqlite")
	viper.SetDefault("database.sslmode", "")
	viper.SetDefault("database.maxConnLifetime", 0)
	viper.SetDefault("database.maxConnIdleTime", 0)
	viper.SetDefault("database.maxOpenConns", 0)
	viper.SetDefault("database.maxIdleConns", 0)

	viper.SetDefault("auth.secret", "top-secret")
	viper.SetDefault("auth.accessTokenLifetime", 0)
	viper.SetDefault("auth.refreshTokenLifetime", 0)

	viper.SetDefault("cache.addr", "localhost:6379")
	viper.SetDefault("cache.db", 0)
	viper.SetDefault("cache.password", "")
}

type Log struct {
	Level     string
	Format    string
	File      string
	AddSource bool
}
type Server struct {
	Host      string
	Port      int
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

type Database struct {
	Driver          string
	Host            string
	Port            uint16
	Username        string
	Password        string
	Database        string
	SSLMode         string
	MaxConns        int
	MaxIdleConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type Authentication struct {
	Secret               string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

type Cache struct {
	Addr     string
	DB       int
	Password string
}
