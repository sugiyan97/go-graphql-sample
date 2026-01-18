package resolver

import (
	"context"
	"testing"

	"go-graphql-sample/internal/model"

	"github.com/google/uuid"
)

// TestUsers はUsersクエリのテストです
func TestUsers(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	queryResolver := &queryResolver{resolver}
	ctx := context.Background()

	// テストデータの準備
	user1 := createTestUser(t, db, "John Doe", "john@example.com")
	user2 := createTestUser(t, db, "Jane Doe", "jane@example.com")

	// テスト実行
	users, err := queryResolver.Users(ctx)
	if err != nil {
		t.Fatalf("Users() error = %v", err)
	}

	// 結果の検証
	if len(users) != 2 {
		t.Errorf("Users() returned %d users, want 2", len(users))
	}

	// ユーザーが正しく取得されているか確認
	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	if userMap[user1.ID] == nil {
		t.Error("Users() did not return user1")
	}
	if userMap[user2.ID] == nil {
		t.Error("Users() did not return user2")
	}

	// ユーザー情報の検証
	if userMap[user1.ID].Name != "John Doe" {
		t.Errorf("Users() user1.Name = %v, want John Doe", userMap[user1.ID].Name)
	}
	if userMap[user1.ID].Email != "john@example.com" {
		t.Errorf("Users() user1.Email = %v, want john@example.com", userMap[user1.ID].Email)
	}
}

// TestUsers_Empty は空のデータベースでのUsersクエリのテストです
func TestUsers_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	queryResolver := &queryResolver{resolver}
	ctx := context.Background()

	// テスト実行
	users, err := queryResolver.Users(ctx)
	if err != nil {
		t.Fatalf("Users() error = %v", err)
	}

	// 結果の検証
	if len(users) != 0 {
		t.Errorf("Users() returned %d users, want 0", len(users))
	}
}

// TestUser はUserクエリのテストです
func TestUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	queryResolver := &queryResolver{resolver}
	ctx := context.Background()

	// テストデータの準備
	user := createTestUser(t, db, "John Doe", "john@example.com")

	// テスト実行
	result, err := queryResolver.User(ctx, user.ID)
	if err != nil {
		t.Fatalf("User() error = %v", err)
	}

	// 結果の検証
	if result == nil {
		t.Fatal("User() returned nil")
	}
	if result.ID != user.ID {
		t.Errorf("User() ID = %v, want %v", result.ID, user.ID)
	}
	if result.Name != "John Doe" {
		t.Errorf("User() Name = %v, want John Doe", result.Name)
	}
	if result.Email != "john@example.com" {
		t.Errorf("User() Email = %v, want john@example.com", result.Email)
	}
}

// TestUser_NotFound は存在しないユーザーを取得する場合のテストです
func TestUser_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	queryResolver := &queryResolver{resolver}
	ctx := context.Background()

	// 存在しないIDでテスト実行
	nonExistentID := uuid.New().String()
	result, err := queryResolver.User(ctx, nonExistentID)

	// エラーが返されることを確認
	if err == nil {
		t.Error("User() expected error, got nil")
	}
	if result != nil {
		t.Errorf("User() expected nil, got %v", result)
	}
}
