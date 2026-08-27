package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     map[string]*DBConfig // 支持多个命名数据库连接
	JWT    JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
	MaxIdle  int // 最大空闲连接数
	MaxOpen  int // 最大打开连接数
}

type JWTConfig struct {
	Secret             string
	ExpireHours        int
	RefreshExpireHours int
}

var AppConfig *Config

// Load 加载配置：优先从 .env 读取，否则从环境变量读取
func Load() error {
	// 尝试加载 .env 文件（如果不存在不报错，直接用环境变量）
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		DB: map[string]*DBConfig{
			"default": {
				Host:     getEnv("DB_HOST", "127.0.0.1"),
				Port:     getEnv("DB_PORT", "3306"),
				User:     getEnv("DB_USER", "root"),
				Password: getEnv("DB_PASSWORD", "root"),
				DBName:   getEnv("DB_NAME", "golang_admin_best"),
				Charset:  getEnv("DB_CHARSET", "utf8mb4"),
				MaxIdle:  getEnvAsInt("DB_MAX_IDLE", 10),
				MaxOpen:  getEnvAsInt("DB_MAX_OPEN", 100),
			},
			// 可以添加更多连接，例如：
			// "slave": {...},
			// "tenant1": {...},
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "art-design-pro-secret-change-in-production"),
			ExpireHours:        getEnvAsInt("JWT_EXPIRE_HOURS", 24),
			RefreshExpireHours: getEnvAsInt("JWT_REFRESH_EXPIRE_HOURS", 168),
		},
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}

// GetDSNWithoutDB 获取不指定数据库的连接字符串（用于建库）
func (c *DBConfig) GetDSNWithoutDB() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Charset)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
