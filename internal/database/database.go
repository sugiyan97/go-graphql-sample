package database

import (
	"fmt"
	"os"

	"go-graphql-sample/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB はデータベース接続のグローバル変数です
	DB *gorm.DB
)

// InitDB はデータベース接続を初期化します
func InitDB() error {
	// データベースファイルのパス
	dbPath := "data.db"
	if path := os.Getenv("DB_PATH"); path != "" {
		dbPath = path
	}

	// GORMの設定
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// SQLiteデータベースに接続
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// マイグレーション（テーブル作成）
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// autoMigrate はデータベースのマイグレーションを実行します
func autoMigrate() error {
	// Userテーブルを作成
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("failed to migrate User model: %w", err)
	}

	return nil
}

// CloseDB はデータベース接続を閉じます
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
