# Step 4: リゾルバー実装（Query）

## 学習目標

このStepでは、GraphQLのクエリリゾルバーを実装します。以下のことを学びます：

- リゾルバーの基本構造
- クエリリゾルバーの実装方法
- データベースからのデータ取得
- エラーハンドリング
- コンテキストの使用

## 前提知識

- GraphQLのクエリの概念
- GORMの基本的な使い方
- Goのエラーハンドリング

## 実装内容

### 1. Resolver構造体の更新

`internal/resolver/resolver.go`で、Resolver構造体にデータベース接続を追加します。

#### Resolver構造体

```go
type Resolver struct {
	DB *gorm.DB
}

func NewResolver() *Resolver {
	return &Resolver{
		DB: database.DB,
	}
}
```

#### 説明

- **DBフィールド**: GORMのデータベース接続を保持
- **NewResolver関数**: Resolverインスタンスを作成するファクトリー関数
- **依存性注入**: データベース接続をResolverに注入

### 2. クエリリゾルバーの実装

`internal/resolver/schema.go`で、クエリリゾルバーを実装します。

#### Usersリゾルバー

全ユーザー一覧を取得するリゾルバーです：

```go
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	
	// データベースから全ユーザーを取得
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	
	return users, nil
}
```

#### 実装の説明

1. **変数宣言**: `var users []*model.User`でユーザー配列を宣言
2. **データベースクエリ**: `r.DB.Find(&users)`で全ユーザーを取得
3. **エラーハンドリング**: エラーが発生した場合はエラーを返す
4. **結果の返却**: 取得したユーザー配列を返す

#### Userリゾルバー

指定IDのユーザーを取得するリゾルバーです：

```go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	
	// データベースから指定IDのユーザーを取得
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	
	return &user, nil
}
```

#### 実装の説明

1. **変数宣言**: `var user model.User`でユーザー変数を宣言
2. **条件付きクエリ**: `r.DB.Where("id = ?", id).First(&user)`でIDに一致するユーザーを取得
3. **エラーハンドリング**:
   - レコードが見つからない場合: `gorm.ErrRecordNotFound`をチェックして適切なエラーメッセージを返す
   - その他のエラー: 一般的なエラーメッセージを返す
4. **結果の返却**: 取得したユーザーへのポインタを返す

### 3. GORMクエリの基本

#### Findメソッド

複数のレコードを取得します：

```go
var users []*model.User
r.DB.Find(&users)
```

#### Firstメソッド

最初の1件を取得します（見つからない場合はエラー）：

```go
var user model.User
r.DB.First(&user, "id = ?", id)
```

#### Whereメソッド

条件を指定してクエリを構築します：

```go
r.DB.Where("id = ?", id).First(&user)
```

`?`はプレースホルダーで、SQLインジェクションを防ぎます。

### 4. エラーハンドリング

#### エラーの種類

- **gorm.ErrRecordNotFound**: レコードが見つからない場合
- **その他のエラー**: データベース接続エラーなど

#### エラーメッセージの返却

GraphQLでは、エラーを返す際に適切なメッセージを返すことが重要です：

```go
if err == gorm.ErrRecordNotFound {
	return nil, fmt.Errorf("user with id %s not found", id)
}
```

### 5. コンテキストの使用

`context.Context`は、リクエストのライフサイクルを管理するために使用されます。将来的に以下のような用途で使用できます：

- **認証情報の取得**: ユーザー情報をコンテキストから取得
- **タイムアウト制御**: リクエストのタイムアウトを設定
- **トレーシング**: リクエストの追跡情報を保持

現時点では使用しませんが、関数シグネチャには含めます。

## 動作確認

### サーバーの起動

サーバーを起動する前に、データベースにテストデータを追加する必要があります。Step 5でミューテーションを実装するまで、手動でデータを追加するか、Step 6でサーバーを起動してからミューテーションでデータを作成します。

### GraphQLクエリのテスト

サーバーが起動したら、GraphQL Playgroundで以下のクエリをテストできます：

#### 全ユーザー取得

```graphql
query {
  users {
    id
    name
    email
    createdAt
  }
}
```

#### 特定ユーザー取得

```graphql
query {
  user(id: "ユーザーID") {
    id
    name
    email
    createdAt
  }
}
```

### エラーの確認

存在しないIDでクエリを実行すると、エラーメッセージが返されます：

```graphql
query {
  user(id: "non-existent-id") {
    id
    name
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

### データベース接続エラー

- `database.InitDB()`が実行されているか確認
- `database.DB`がnilでないか確認
- データベースファイルが存在するか確認

### クエリが空の結果を返す

- データベースにデータが存在するか確認
- `sqlite3 data.db "SELECT * FROM users;"`で確認
- ミューテーションでデータを作成（Step 5）

### コンパイルエラー

- `gorm.io/gorm`がインストールされているか確認: `go mod tidy`
- インポートパスが正しいか確認
- `go generate ./...`を実行してコードを再生成

### リゾルバーが呼ばれない

- `NewResolver()`が正しく呼ばれているか確認（Step 6で実装）
- GraphQLスキーマとリゾルバーのシグネチャが一致しているか確認

## 次のステップ

クエリリゾルバーの実装が完了したら、[Step 5: リゾルバー実装（Mutation）](step-5-mutation-resolver.md) に進みましょう。

## 補足説明

### リゾルバーの役割

リゾルバーは、GraphQLクエリ/ミューテーションを実際のデータ操作に変換する関数です：

1. **クエリの受信**: GraphQLクエリを受け取る
2. **データ取得**: データベースからデータを取得
3. **データ変換**: 必要に応じてデータを変換
4. **結果の返却**: GraphQLレスポンスを返す

### GORMのクエリビルダー

GORMは、型安全なクエリを構築するためのメソッドチェーンを提供します：

```go
// 条件を追加
r.DB.Where("name = ?", "John")

// ソート
r.DB.Order("created_at DESC")

// 制限
r.DB.Limit(10)

// オフセット
r.DB.Offset(20)

// 実行
r.DB.Find(&users)
```

### エラーハンドリングのベストプラクティス

1. **具体的なエラーメッセージ**: ユーザーが理解できるメッセージを返す
2. **エラーの種類を区別**: レコードが見つからない場合とその他のエラーを区別
3. **エラーのラッピング**: `fmt.Errorf`でエラーをラップしてコンテキストを追加

### パフォーマンスの考慮

- **N+1問題**: 関連データを取得する際は、`Preload`を使用
- **ページネーション**: 大量のデータを取得する場合は、`Limit`と`Offset`を使用
- **インデックス**: 検索頻度の高いカラムにインデックスを設定

### コンテキストの活用

将来的に、コンテキストから以下の情報を取得できます：

```go
// 認証情報の取得
userID := ctx.Value("userID").(string)

// タイムアウトの設定
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

Step 6でサーバーを実装する際に、これらの機能を追加できます。
