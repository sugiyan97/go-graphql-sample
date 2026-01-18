package resolver

import (
	"os"
	"testing"
	"time"

	"go-graphql-sample/internal/model"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB はテスト用のデータベースをセットアップします
func setupTestDB(t *testing.T) *gorm.DB {
	// テスト用のデータベースファイル名
	dbPath := "test.db"

	// 既存のテストデータベースを削除
	os.Remove(dbPath)

	// データベース接続
	config := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// マイグレーション
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// teardownTestDB はテスト用のデータベースをクリーンアップします
func teardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Failed to get database connection: %v", err)
		return
	}
	sqlDB.Close()
	os.Remove("test.db")
}

// createTestUser はテスト用のユーザーを作成します
func createTestUser(t *testing.T, db *gorm.DB, name, email string) *model.User {
	now := time.Now()
	user := &model.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}
