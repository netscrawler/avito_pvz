package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ENV  string     `yaml:"env"`
	DB   DBConfig   `yaml:"db"`
	GRPC GRPCServer `yaml:"grpcServer"`
	HTTP HTTPServer `yaml:"httpServer"`
	JWT  JWT        `yaml:"jwt"`
}

type JWT struct {
	SecretKey string        `yaml:"secretKey"`
	Expire    time.Duration `yaml:"expire"`
}

type DBConfig struct {
	Type                string        `yaml:"type"                env-default:"postgres"`
	Port                int           `yaml:"port"                env-default:"5432"`
	Host                string        `yaml:"host"                env-default:"localhost"`
	User                string        `yaml:"user"                env-default:"user"`
	Password            string        `yaml:"password"            env-default:"password"`
	Name                string        `yaml:"name"                env-default:"postgres"`
	SSLMode             string        `yaml:"sslMode"             env-default:"false"`
	PoolMaxConn         int           `yaml:"poolMaxConn"         env-default:"10"`
	PoolMaxConnLifetime time.Duration `yaml:"poolMaxConnLifetime" env-default:"1h30m"`
}

type GRPCServer struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"     env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"     env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (db *DBConfig) DSN() string {
	encodedPassword := url.QueryEscape(db.Password)

	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_max_conn_lifetime=%s",
		db.Type,
		db.User,
		encodedPassword,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
		db.PoolMaxConn,
		db.PoolMaxConnLifetime.String(),
	)
}
