package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	JWT      JWTConfig
	CORS     CORSConfig
	MinIO    MinIOConfig
	Admin    AdminConfig
}

type AdminConfig struct {
	Email            string
	Password         string
	AllowSelfRegister bool
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	PublicURL string
}

type AppConfig struct {
	Env  string
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSL      string
}

type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

type CORSConfig struct {
	Origins string
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_SSL", "disable")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	viper.SetDefault("CORS_ORIGINS", "http://localhost:3000")
	viper.SetDefault("ALLOW_SELF_REGISTER", true)
	viper.SetDefault("MINIO_BUCKET", "kimsha-images")
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("MINIO_PUBLIC_URL", "http://localhost:9000")

	expiry, _ := time.ParseDuration(viper.GetString("JWT_EXPIRY"))
	refreshExpiry, _ := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRY"))

	return &Config{
		App: AppConfig{
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
		},
		Postgres: PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSL:      viper.GetString("POSTGRES_SSL"),
		},
		JWT: JWTConfig{
			Secret:        viper.GetString("JWT_SECRET"),
			Expiry:        expiry,
			RefreshExpiry: refreshExpiry,
		},
		CORS: CORSConfig{
			Origins: viper.GetString("CORS_ORIGINS"),
		},
		MinIO: MinIOConfig{
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
			PublicURL: viper.GetString("MINIO_PUBLIC_URL"),
		},
		Admin: AdminConfig{
			Email:             viper.GetString("ADMIN_EMAIL"),
			Password:          viper.GetString("ADMIN_PASSWORD"),
			AllowSelfRegister: viper.GetBool("ALLOW_SELF_REGISTER"),
		},
	}, nil
}
