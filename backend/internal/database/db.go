package database

import (
	"art-design-pro-api/internal/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connections = make(map[string]*gorm.DB)

// Init 初始化所有配置的数据库连接
func Init() error {
	logMode := logger.Info
	if config.AppConfig.Server.Mode == "release" {
		logMode = logger.Error
	}

	for name, dbConfig := range config.AppConfig.DB {
		dsn := dbConfig.GetDSN()
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})
		if err != nil {
			return fmt.Errorf("failed to connect database '%s': %w", name, err)
		}

		// 配置连接池
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdle)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpen)
		// 设置连接最大生命周期，避免复用被 MySQL/防火墙断开的失效连接
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)

		connections[name] = db
		log.Printf("✅ Database '%s' connected successfully\n", name)
	}

	return nil
}

// Get 获取指定名称的数据库连接，默认返回 "default"
func Get(name ...string) *gorm.DB {
	connName := "default"
	if len(name) > 0 && name[0] != "" {
		connName = name[0]
	}
	return connections[connName]
}

// Close 关闭所有数据库连接
func Close() error {
	for name, db := range connections {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("⚠️  Failed to get underlying DB for '%s': %v\n", name, err)
				continue
			}
			if err := sqlDB.Close(); err != nil {
				log.Printf("⚠️  Failed to close database '%s': %v\n", name, err)
			}
		}
	}
	return nil
}
