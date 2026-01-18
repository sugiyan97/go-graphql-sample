package resolver

import (
	"context"
	"testing"

	"go-graphql-sample/internal/gql/generated"
	"go-graphql-sample/internal/model"

	"gorm.io/gorm"
)

// TestCreateUser はCreateUserミューテーションのテストです
func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	// テスト実行
	input := generated.CreateUserInput{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user, err := mutationResolver.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	// 結果の検証
	if user == nil {
		t.Fatal("CreateUser() returned nil")
	}
	if user.ID == "" {
		t.Error("CreateUser() ID is empty")
	}
	if user.Name != "John Doe" {
		t.Errorf("CreateUser() Name = %v, want John Doe", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("CreateUser() Email = %v, want john@example.com", user.Email)
	}

	// データベースに保存されているか確認
	var savedUser model.User
	if err := db.Where("id = ?", user.ID).First(&savedUser).Error; err != nil {
		t.Fatalf("Failed to find created user: %v", err)
	}
	if savedUser.Name != "John Doe" {
		t.Errorf("Saved user Name = %v, want John Doe", savedUser.Name)
	}
}

// TestCreateUser_DuplicateEmail は重複メールアドレスでのCreateUserのテストです
func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	// 最初のユーザーを作成
	input1 := generated.CreateUserInput{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	_, err := mutationResolver.CreateUser(ctx, input1)
	if err != nil {
		t.Fatalf("CreateUser() first user error = %v", err)
	}

	// 同じメールアドレスでユーザーを作成（エラーが期待される）
	input2 := generated.CreateUserInput{
		Name:  "Jane Doe",
		Email: "john@example.com",
	}
	_, err = mutationResolver.CreateUser(ctx, input2)
	if err == nil {
		t.Error("CreateUser() expected error for duplicate email, got nil")
	}
}

// TestCreateUser_EmptyName は空の名前でのCreateUserのテストです
func TestCreateUser_EmptyName(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	input := generated.CreateUserInput{
		Name:  "",
		Email: "john@example.com",
	}

	_, err := mutationResolver.CreateUser(ctx, input)
	if err == nil {
		t.Error("CreateUser() expected error for empty name, got nil")
	}
}

// TestUpdateUser はUpdateUserミューテーションのテストです
func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	// テストデータの準備
	user := createTestUser(t, db, "John Doe", "john@example.com")

	// テスト実行
	newName := "Jane Doe"
	input := generated.UpdateUserInput{
		Name: &newName,
	}

	updatedUser, err := mutationResolver.UpdateUser(ctx, user.ID, input)
	if err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}

	// 結果の検証
	if updatedUser == nil {
		t.Fatal("UpdateUser() returned nil")
	}
	if updatedUser.Name != "Jane Doe" {
		t.Errorf("UpdateUser() Name = %v, want Jane Doe", updatedUser.Name)
	}
	if updatedUser.Email != "john@example.com" {
		t.Errorf("UpdateUser() Email = %v, want john@example.com", updatedUser.Email)
	}

	// データベースに反映されているか確認
	var savedUser model.User
	if err := db.Where("id = ?", user.ID).First(&savedUser).Error; err != nil {
		t.Fatalf("Failed to find updated user: %v", err)
	}
	if savedUser.Name != "Jane Doe" {
		t.Errorf("Saved user Name = %v, want Jane Doe", savedUser.Name)
	}
}

// TestUpdateUser_NotFound は存在しないユーザーの更新のテストです
func TestUpdateUser_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	nonExistentID := "non-existent-id"
	newName := "Jane Doe"
	input := generated.UpdateUserInput{
		Name: &newName,
	}

	_, err := mutationResolver.UpdateUser(ctx, nonExistentID, input)
	if err == nil {
		t.Error("UpdateUser() expected error for non-existent user, got nil")
	}
}

// TestDeleteUser はDeleteUserミューテーションのテストです
func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	// テストデータの準備
	user := createTestUser(t, db, "John Doe", "john@example.com")

	// テスト実行
	result, err := mutationResolver.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("DeleteUser() error = %v", err)
	}

	// 結果の検証
	if !result {
		t.Error("DeleteUser() returned false, want true")
	}

	// データベースから削除されているか確認
	var deletedUser model.User
	err = db.Where("id = ?", user.ID).First(&deletedUser).Error
	if err != gorm.ErrRecordNotFound {
		t.Errorf("DeleteUser() user still exists in database: %v", err)
	}
}

// TestDeleteUser_NotFound は存在しないユーザーの削除のテストです
func TestDeleteUser_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	nonExistentID := "non-existent-id"
	_, err := mutationResolver.DeleteUser(ctx, nonExistentID)
	if err == nil {
		t.Error("DeleteUser() expected error for non-existent user, got nil")
	}
}
