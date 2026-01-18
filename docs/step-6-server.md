# Step 6: サーバー起動と動作確認

## 学習目標

このStepでは、HTTPサーバーを実装し、GraphQLエンドポイントとGraphQL Playgroundを設定します。以下のことを学びます：

- HTTPサーバーの実装
- GraphQLハンドラーの設定
- GraphQL Playgroundの設定
- サーバーの起動と動作確認
- GraphQLクエリ/ミューテーションの実行

## 前提知識

- HTTPサーバーの基本的な概念
- GraphQLの基本的な理解
- クライアント-サーバー通信

## 実装内容

### 1. HTTPサーバーの実装

`cmd/server/main.go`でHTTPサーバーを実装します。

#### 実装

```go
package main

import (
	"log"
	"net/http"
	"os"

	"go-graphql-sample/internal/database"
	"go-graphql-sample/internal/gql/generated"
	"go-graphql-sample/internal/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	// データベース接続の初期化
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// ポート番号の取得
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Resolverの作成
	resolver := resolver.NewResolver()

	// GraphQLサーバーの設定
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)

	// ルーティングの設定
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	// サーバーの起動
	log.Printf("Server is running on http://localhost:%s", port)
	log.Printf("GraphQL Playground: http://localhost:%s", port)
	log.Printf("GraphQL Endpoint: http://localhost:%s/query", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

#### 実装の説明

1. **データベース初期化**: `database.InitDB()`でデータベース接続を初期化
2. **ポート設定**: 環境変数`PORT`からポート番号を取得（デフォルトは8080）
3. **Resolver作成**: `resolver.NewResolver()`でResolverインスタンスを作成
4. **GraphQLサーバー設定**: `handler.NewDefaultServer()`でGraphQLハンドラーを作成
5. **ルーティング設定**:
   - `/`: GraphQL Playground（開発用のUI）
   - `/query`: GraphQLエンドポイント
6. **サーバー起動**: `http.ListenAndServe()`でサーバーを起動

### 2. GraphQLハンドラー

gqlgenの`handler.NewDefaultServer()`は、GraphQLリクエストを処理するHTTPハンドラーを作成します。

#### 設定

```go
srv := handler.NewDefaultServer(
	generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}),
)
```

- **NewExecutableSchema**: GraphQLスキーマとリゾルバーを組み合わせて実行可能なスキーマを作成
- **Config**: リゾルバー、ディレクティブ、複雑度制限などの設定

### 3. GraphQL Playground

GraphQL Playgroundは、GraphQLクエリ/ミューテーションを実行できる対話型のUIです。

#### 設定

```go
http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
```

- **第一引数**: Playgroundのタイトル
- **第二引数**: GraphQLエンドポイントのパス

### 4. 環境変数による設定

ポート番号を環境変数で変更できます：

```bash
export PORT=3000
go run cmd/server/main.go
```

## 動作確認

### サーバーの起動

```bash
go run cmd/server/main.go
```

以下のメッセージが表示されれば成功です：

```
Server is running on http://localhost:8080
GraphQL Playground: http://localhost:8080
GraphQL Endpoint: http://localhost:8080/query
```

### GraphQL Playgroundでの動作確認

1. ブラウザで `http://localhost:8080` を開く
2. GraphQL Playgroundが表示される
3. 左側のエディタでクエリ/ミューテーションを記述
4. 実行ボタンをクリックして結果を確認

### クエリの実行

#### ユーザー一覧取得

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

### ミューテーションの実行

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

#### 存在しないユーザーを取得

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
      "message": "user with id non-existent-id not found",
      "path": ["user"]
    }
  ],
  "data": null
}
```

#### 重複メールアドレスでユーザー作成

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
      "message": "user with email existing@example.com already exists",
      "path": ["createUser"]
    }
  ],
  "data": null
}
```

### curlでの動作確認

GraphQL Playgroundを使わずに、curlで直接リクエストを送信することもできます：

```bash
# クエリの実行
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { users { id name email } }"
  }'

# ミューテーションの実行
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createUser(input: { name: \"John Doe\" email: \"john@example.com\" }) { id name email } }"
  }'
```

## トラブルシューティング

### サーバーが起動しない

- ポート8080が既に使用されているか確認: `lsof -i :8080`
- 別のポートを指定: `PORT=3000 go run cmd/server/main.go`
- エラーメッセージを確認

### データベース接続エラー

- データベースファイルが作成されているか確認: `ls -la data.db`
- データベースの権限を確認
- `database.InitDB()`のエラーメッセージを確認

### GraphQL Playgroundが表示されない

- ブラウザで `http://localhost:8080` にアクセス
- サーバーが正しく起動しているか確認
- ポート番号が正しいか確認

### クエリ/ミューテーションが実行されない

- GraphQLスキーマが正しく定義されているか確認
- リゾルバーが正しく実装されているか確認
- エラーメッセージを確認（Playgroundの右側に表示）

### データが取得できない

- データベースにデータが存在するか確認: `sqlite3 data.db "SELECT * FROM users;"`
- ミューテーションでデータを作成
- リゾルバーの実装を確認

## 次のステップ

サーバーの起動と動作確認が完了したら、[Step 7: テストの実装](step-7-testing.md) に進みましょう。

## 補足説明

### GraphQL Playgroundの機能

- **スキーマ閲覧**: 右側の「Schema」タブでスキーマを確認
- **クエリ実行**: 左側のエディタでクエリ/ミューテーションを記述して実行
- **結果表示**: 右側に結果が表示される
- **履歴**: 実行したクエリの履歴を確認

### HTTPメソッド

GraphQLは通常、POSTメソッドでリクエストを送信します：

- **GET**: クエリのみ（推奨されない）
- **POST**: クエリとミューテーション（推奨）

### リクエスト形式

GraphQLリクエストは以下の形式で送信されます：

```json
{
  "query": "query { users { id name } }",
  "variables": {},
  "operationName": null
}
```

- **query**: GraphQLクエリ文字列
- **variables**: 変数（オプション）
- **operationName**: 操作名（オプション）

### レスポンス形式

GraphQLレスポンスは以下の形式で返されます：

```json
{
  "data": {
    "users": [
      {
        "id": "1",
        "name": "John Doe",
        "email": "john@example.com"
      }
    ]
  },
  "errors": []
}
```

- **data**: クエリの結果
- **errors**: エラーの配列（エラーがない場合は空配列）

### 本番環境での注意点

- **GraphQL Playgroundを無効化**: 本番環境ではPlaygroundを無効化することを推奨
- **CORS設定**: 必要に応じてCORSを設定
- **認証・認可**: 認証トークンの検証を実装
- **レート制限**: リクエストのレート制限を実装
- **ログ記録**: リクエストとエラーをログに記録

### パフォーマンスの最適化

- **データローダー**: N+1問題を解決するためにデータローダーを使用
- **クエリの複雑度制限**: 過度に複雑なクエリを制限
- **キャッシング**: 頻繁にアクセスされるデータをキャッシュ
