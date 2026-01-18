# Step 7: テストの実装

## 学習目標

このStepでは、GraphQLクエリ/ミューテーションのテストを実装します。以下のことを学びます：

- Goのテストの書き方
- テスト用データベースのセットアップ
- クエリリゾルバーのテスト
- ミューテーションリゾルバーのテスト
- エラーケースのテスト

## 前提知識

- Goのテストの基本的な理解
- testingパッケージの使い方
- テーブル駆動テストの概念

## 実装内容

### 1. テスト環境のセットアップ

`internal/resolver/resolver_test.go`で、テスト用のデータベースセットアップ関数を実装します。

#### セットアップ関数

```go
// setupTestDB はテスト用のデータベースをセットアップします
func setupTestDB(t *testing.T) *gorm.DB {
	dbPath := "test.db"
	os.Remove(dbPath)  // 既存のテストデータベースを削除

	config := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

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
```

#### ヘルパー関数

```go
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
```

### 2. クエリのテスト

`internal/resolver/query_test.go`で、クエリリゾルバーのテストを実装します。

#### Usersクエリのテスト

```go
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
}
```

#### Userクエリのテスト

```go
func TestUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	queryResolver := &queryResolver{resolver}
	ctx := context.Background()

	user := createTestUser(t, db, "John Doe", "john@example.com")

	result, err := queryResolver.User(ctx, user.ID)
	if err != nil {
		t.Fatalf("User() error = %v", err)
	}

	if result.ID != user.ID {
		t.Errorf("User() ID = %v, want %v", result.ID, user.ID)
	}
}
```

#### エラーケースのテスト

```go
func TestUser_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	queryResolver := &queryResolver{resolver}
	ctx := context.Background()

	nonExistentID := uuid.New().String()
	result, err := queryResolver.User(ctx, nonExistentID)

	if err == nil {
		t.Error("User() expected error, got nil")
	}
	if result != nil {
		t.Errorf("User() expected nil, got %v", result)
	}
}
```

### 3. ミューテーションのテスト

`internal/resolver/mutation_test.go`で、ミューテーションリゾルバーのテストを実装します。

#### CreateUserのテスト

```go
func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	input := generated.CreateUserInput{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user, err := mutationResolver.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("CreateUser() Name = %v, want John Doe", user.Name)
	}
}
```

#### UpdateUserのテスト

```go
func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	user := createTestUser(t, db, "John Doe", "john@example.com")

	newName := "Jane Doe"
	input := generated.UpdateUserInput{
		Name: &newName,
	}

	updatedUser, err := mutationResolver.UpdateUser(ctx, user.ID, input)
	if err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}

	if updatedUser.Name != "Jane Doe" {
		t.Errorf("UpdateUser() Name = %v, want Jane Doe", updatedUser.Name)
	}
}
```

#### DeleteUserのテスト

```go
func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	resolver := &Resolver{DB: db}
	mutationResolver := &mutationResolver{resolver}
	ctx := context.Background()

	user := createTestUser(t, db, "John Doe", "john@example.com")

	result, err := mutationResolver.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("DeleteUser() error = %v", err)
	}

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
```

### 4. テストの実行

#### すべてのテストを実行

```bash
go test ./...
```

#### 特定のパッケージのテストを実行

```bash
go test ./internal/resolver
```

#### 詳細な出力でテストを実行

```bash
go test -v ./internal/resolver
```

#### カバレッジを確認

```bash
go test -cover ./internal/resolver
```

#### カバレッジレポートを生成

```bash
go test -coverprofile=coverage.out ./internal/resolver
go tool cover -html=coverage.out
```

## 動作確認

### テストの実行

```bash
# すべてのテストを実行
go test ./...

# 特定のパッケージのテストを実行
go test ./internal/resolver

# 詳細な出力で実行
go test -v ./internal/resolver
```

### テスト結果の確認

テストが成功すると、以下のような出力が表示されます：

```
PASS
ok      go-graphql-sample/internal/resolver    0.123s
```

### カバレッジの確認

```bash
go test -cover ./internal/resolver
```

出力例：

```
PASS
coverage: 85.7% of statements
ok      go-graphql-sample/internal/resolver    0.123s
```

## トラブルシューティング

### テストが失敗する

- テストデータベースファイルが残っていないか確認: `ls -la test.db`
- データベース接続エラーを確認
- テストのロジックを確認

### テストが並列実行できない

- 各テストで独立したデータベースを使用する
- テストデータベースファイル名を一意にする（例: `test_<timestamp>.db`）

### カバレッジが低い

- エラーケースのテストを追加
- エッジケースのテストを追加
- すべての分岐をテストする

### テストが遅い

- テストデータベースをメモリ内に作成（`:memory:`）
- 不要なテストデータを削除
- テストを並列実行（`t.Parallel()`）

## 次のステップ

テストの実装が完了したら、すべてのStepが完了です！お疲れ様でした。

## 補足説明

### Goのテストの基本

- **テストファイル**: `*_test.go`という命名規則
- **テスト関数**: `Test`で始まる関数名
- **testing.T**: テストの状態を管理する構造体

### テーブル駆動テスト

複数のテストケースを効率的に実行する方法：

```go
func TestUsers_MultipleCases(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*gorm.DB)
		want    int
		wantErr bool
	}{
		{
			name: "empty database",
			setup: func(db *gorm.DB) {},
			want: 0,
		},
		{
			name: "one user",
			setup: func(db *gorm.DB) {
				createTestUser(t, db, "John", "john@example.com")
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer teardownTestDB(t, db)
			tt.setup(db)
			// テスト実行
		})
	}
}
```

### テストのベストプラクティス

1. **独立性**: 各テストは独立して実行できる
2. **再現性**: 同じ条件で常に同じ結果が得られる
3. **高速性**: テストは高速に実行される
4. **明確性**: テストの意図が明確である

### モックの使用

外部依存（データベースなど）をモック化する場合：

```go
type MockDB struct {
	users []*model.User
}

func (m *MockDB) Find(dest interface{}) *gorm.DB {
	// モック実装
	return &gorm.DB{}
}
```

ただし、このプロジェクトでは実際のデータベースを使用してテストしています。

### 統合テストとユニットテスト

- **ユニットテスト**: 個別の関数やメソッドをテスト
- **統合テスト**: 複数のコンポーネントを組み合わせてテスト

このプロジェクトのテストは、リゾルバーとデータベースを組み合わせた統合テストです。
