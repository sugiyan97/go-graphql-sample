# Step 5: リゾルバー実装（Mutation）

## 学習目標

このStepでは、GraphQLのミューテーションリゾルバーを実装します。以下のことを学びます：

- ミューテーションリゾルバーの実装方法
- データの作成・更新・削除
- 入力値のバリデーション
- エラーハンドリング
- UUIDの生成

## 前提知識

- GraphQLのミューテーションの概念
- GORMの基本的な使い方
- Goのエラーハンドリング
- ポインタ型の理解

## 実装内容

### 1. CreateUserリゾルバー

新しいユーザーを作成するミューテーションです。

#### 実装

```go
func (r *mutationResolver) CreateUser(ctx context.Context, input generated.CreateUserInput) (*model.User, error) {
	// 入力値のバリデーション
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// メールアドレスの重複チェック
	var existingUser model.User
	if err := r.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("user with email %s already exists", input.Email)
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	// 新しいユーザーを作成
	now := time.Now()
	user := &model.User{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// データベースに保存
	if err := r.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
```

#### 実装の説明

1. **入力値のバリデーション**: NameとEmailが空でないかチェック
2. **重複チェック**: メールアドレスが既に存在しないか確認
3. **UUID生成**: `uuid.New().String()`で一意のIDを生成
4. **タイムスタンプ設定**: CreatedAtとUpdatedAtを現在時刻に設定
5. **データベース保存**: `r.DB.Create(user)`でデータベースに保存
6. **結果の返却**: 作成されたユーザーを返す

#### ポイント

- **バリデーション**: サーバー側で入力値を検証
- **一意制約**: メールアドレスの重複を防ぐ
- **UUID**: 一意のIDを自動生成

### 2. UpdateUserリゾルバー

既存のユーザーを更新するミューテーションです。

#### 実装

```go
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input generated.UpdateUserInput) (*model.User, error) {
	// ユーザーを取得
	var user model.User
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// メールアドレスの重複チェック（変更される場合）
	if input.Email != nil && *input.Email != user.Email {
		var existingUser model.User
		if err := r.DB.Where("email = ? AND id != ?", *input.Email, id).First(&existingUser).Error; err == nil {
			return nil, fmt.Errorf("user with email %s already exists", *input.Email)
		} else if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
	}

	// フィールドを更新
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	user.UpdatedAt = time.Now()

	// データベースを更新
	if err := r.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}
```

#### 実装の説明

1. **ユーザー取得**: IDでユーザーを取得（存在しない場合はエラー）
2. **重複チェック**: メールアドレスが変更される場合、重複をチェック
3. **部分更新**: 入力値がnilでない場合のみ更新（部分更新をサポート）
4. **タイムスタンプ更新**: UpdatedAtを現在時刻に更新
5. **データベース更新**: `r.DB.Save(&user)`でデータベースを更新
6. **結果の返却**: 更新されたユーザーを返す

#### ポイント

- **部分更新**: オプショナルなフィールドは、nilチェックで部分更新を実現
- **ポインタ型**: `*string`型でnilを表現し、更新するかどうかを判断
- **重複チェック**: メールアドレス変更時のみ重複をチェック

### 3. DeleteUserリゾルバー

ユーザーを削除するミューテーションです。

#### 実装

```go
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	// ユーザーが存在するか確認
	var user model.User
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("user with id %s not found", id)
		}
		return false, fmt.Errorf("failed to fetch user: %w", err)
	}

	// ユーザーを削除
	if err := r.DB.Delete(&user).Error; err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	return true, nil
}
```

#### 実装の説明

1. **存在確認**: ユーザーが存在するか確認（存在しない場合はエラー）
2. **削除実行**: `r.DB.Delete(&user)`でデータベースから削除
3. **結果の返却**: 成功時は`true`を返す

#### ポイント

- **論理削除**: GORMはデフォルトで論理削除をサポート（`DeletedAt`フィールドが必要）
- **物理削除**: この実装では物理削除（完全に削除）を行います

### 4. UUIDの生成

UUID（Universally Unique Identifier）は、一意のIDを生成するための標準です。

#### 使用方法

```go
import "github.com/google/uuid"

id := uuid.New().String()
```

#### UUIDの特徴

- **一意性**: 衝突の可能性が極めて低い
- **標準化**: RFC 4122で標準化されている
- **分散生成**: 複数のシステムで生成しても衝突しない

### 5. ポインタ型の使用

`UpdateUserInput`では、オプショナルなフィールドにポインタ型（`*string`）を使用しています。

#### ポインタ型の利点

- **nilチェック**: `nil`で値が指定されていないことを表現
- **部分更新**: 更新するフィールドのみを指定可能

#### 使用例

```go
// 名前のみ更新
input := generated.UpdateUserInput{
	Name: &name,  // ポインタを渡す
	Email: nil,   // nilで更新しない
}

// 値の取得
if input.Name != nil {
	user.Name = *input.Name  // ポインタの値を取得
}
```

## 動作確認

### サーバーの起動

Step 6でサーバーを実装するまで、動作確認はできませんが、実装の確認はできます。

### GraphQLミューテーションのテスト

サーバーが起動したら、GraphQL Playgroundで以下のミューテーションをテストできます。

#### ユーザー作成

```graphql
mutation {
  createUser(input: {
    name: "John Doe"
    email: "john@example.com"
  }) {
    id
    name
    email
    createdAt
  }
}
```

#### ユーザー更新

```graphql
mutation {
  updateUser(id: "ユーザーID", input: {
    name: "Jane Doe"
  }) {
    id
    name
    email
    updatedAt
  }
}
```

#### ユーザー削除

```graphql
mutation {
  deleteUser(id: "ユーザーID")
}
```

### エラーの確認

#### 重複メールアドレス

```graphql
mutation {
  createUser(input: {
    name: "John Doe"
    email: "existing@example.com"
  }) {
    id
  }
}
```

期待されるエラー：

```json
{
  "errors": [
    {
      "message": "user with email existing@example.com already exists"
    }
  ]
}
```

#### 存在しないユーザーの更新

```graphql
mutation {
  updateUser(id: "non-existent-id", input: {
    name: "New Name"
  }) {
    id
  }
}
```

期待されるエラー：

```json
{
  "errors": [
    {
      "message": "user with id non-existent-id not found"
    }
  ]
}
```

## トラブルシューティング

### UUIDパッケージが見つからない

- `go get github.com/google/uuid`を実行
- `go mod tidy`で依存関係を整理

### 重複チェックが機能しない

- データベースの一意制約を確認
- クエリの条件を確認（`id != ?`が正しく設定されているか）

### 部分更新が機能しない

- ポインタ型のnilチェックを確認
- 入力値が正しく渡されているか確認

### タイムスタンプが更新されない

- `UpdatedAt`フィールドが正しく設定されているか確認
- GORMの`Save`メソッドが正しく呼ばれているか確認

## 次のステップ

ミューテーションリゾルバーの実装が完了したら、[Step 6: サーバー起動と動作確認](step-6-server.md) に進みましょう。

## 補足説明

### バリデーションのベストプラクティス

1. **サーバー側バリデーション**: クライアント側のバリデーションに依存しない
2. **具体的なエラーメッセージ**: ユーザーが理解できるメッセージを返す
3. **早期リターン**: エラーが発生したらすぐに返す

### GORMのCreate、Save、Deleteメソッド

- **Create**: 新しいレコードを作成
- **Save**: レコードを保存（存在しない場合は作成、存在する場合は更新）
- **Delete**: レコードを削除

### 論理削除と物理削除

- **論理削除**: `DeletedAt`フィールドを設定して、レコードを「削除済み」としてマーク
- **物理削除**: データベースから完全に削除

この実装では物理削除を使用していますが、本番環境では論理削除を検討してください。

### トランザクション

複数のデータベース操作を行う場合、トランザクションを使用します：

```go
err := r.DB.Transaction(func(tx *gorm.DB) error {
	// 複数の操作
	if err := tx.Create(&user1).Error; err != nil {
		return err
	}
	if err := tx.Create(&user2).Error; err != nil {
		return err
	}
	return nil
})
```

### 入力値のサニタイズ

本番環境では、入力値をサニタイズ（無害化）することを推奨します：

- **SQLインジェクション**: GORMのプレースホルダーで防ぐ
- **XSS攻撃**: HTMLエスケープを実装
- **不正な文字**: バリデーションで制限
