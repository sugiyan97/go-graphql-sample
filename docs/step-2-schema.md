# Step 2: GraphQLスキーマ定義とコード生成

## 学習目標

このStepでは、GraphQLスキーマを定義し、gqlgenを使ってGoコードを生成します。以下のことを学びます：

- GraphQLスキーマの書き方
- 型定義（Type、Input、Scalar）
- QueryとMutationの定義
- gqlgenによるコード生成
- 生成されたコードの理解

## 前提知識

- GraphQLの基本的な概念（Type、Query、Mutation）
- Goの基本的な理解

## 実装内容

### 1. GraphQLスキーマの定義

`internal/gql/schema.graphqls`にGraphQLスキーマを定義します。

#### スキーマの構成

```graphql
# スカラー型の定義
scalar Time

# User型の定義
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: Time!
}

# クエリの定義
type Query {
  users: [User!]!
  user(id: ID!): User
}

# ミューテーションの定義
type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
}

# 入力型の定義
input CreateUserInput {
  name: String!
  email: String!
}

input UpdateUserInput {
  name: String
  email: String
}
```

#### 型の説明

- **User型**: ユーザー情報を表す型
  - `id`: ユーザーID（必須、ID型）
  - `name`: ユーザー名（必須、String型）
  - `email`: メールアドレス（必須、String型）
  - `createdAt`: 作成日時（必須、Time型）

- **Query型**: データ取得用のクエリ
  - `users`: 全ユーザー一覧を取得（戻り値はUserの配列、必須）
  - `user`: 指定IDのユーザーを取得（戻り値はUser、オプション）

- **Mutation型**: データ変更用のミューテーション
  - `createUser`: ユーザーを作成（入力はCreateUserInput、戻り値はUser）
  - `updateUser`: ユーザーを更新（IDとUpdateUserInput、戻り値はUser）
  - `deleteUser`: ユーザーを削除（ID、戻り値はBoolean）

- **Input型**: ミューテーションの入力に使用
  - `CreateUserInput`: ユーザー作成用（name、emailは必須）
  - `UpdateUserInput`: ユーザー更新用（name、emailはオプション）

#### GraphQLの型システム

- **!**: 必須フィールド（Non-Null）
- **[]**: 配列
- **ID**: 一意の識別子（Stringのエイリアス）
- **String**: 文字列
- **Boolean**: 真偽値
- **Scalar**: カスタムスカラー型（Time）

### 2. カスタムスカラー型の実装

GraphQLの標準スカラー型には`Time`がないため、カスタムスカラー型として実装します。

`internal/gql/scalars.go`でTimeスカラー型を実装します：

```go
package gql

import (
	"time"
	"github.com/99designs/gqlgen/graphql"
)

// MarshalTime はtime.TimeをGraphQLの文字列に変換します
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.MarshalTime(t)
}

// UnmarshalTime はGraphQLの文字列をtime.Timeに変換します
func UnmarshalTime(v interface{}) (time.Time, error) {
	return graphql.UnmarshalTime(v)
}
```

gqlgenは、これらの関数を自動的に検出して使用します。

### 3. コード生成の実行

gqlgenを使って、GraphQLスキーマからGoコードを生成します。

#### 生成コマンド

```bash
# 方法1: go generateを使用（推奨）
go generate ./...

# 方法2: gqlgenコマンドを直接実行
gqlgen generate

# 方法3: 特定のディレクトリから実行
cd internal/gql && go generate
```

#### 生成されるファイル

コード生成により、以下のファイルが作成されます：

- `internal/gql/generated/generated.go`: GraphQL実行コード
- `internal/gql/generated/models_gen.go`: モデルコード（Input型など）
- `internal/resolver/resolver.go`: リゾルバーインターフェース
- `internal/resolver/query.go`: クエリリゾルバーのスタブ
- `internal/resolver/mutation.go`: ミューテーションリゾルバーのスタブ

#### go:generateディレクティブ

`internal/gql/generate.go`に`go:generate`ディレクティブを追加することで、`go generate ./...`でコード生成を実行できます：

```go
//go:build ignore
// +build ignore

package main

//go:generate go run github.com/99designs/gqlgen generate
```

### 4. 生成されたコードの理解

#### generated.go

GraphQLクエリ/ミューテーションの実行に必要なコードが含まれます。通常は直接編集しません。

#### models_gen.go

Input型など、GraphQLスキーマから自動生成されるGo構造体が含まれます：

```go
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
```

#### resolver.go

実装すべきリゾルバーインターフェースが定義されます：

```go
type ResolverRoot interface {
	Query() QueryResolver
	Mutation() MutationResolver
}

type QueryResolver interface {
	Users(ctx context.Context) ([]*model.User, error)
	User(ctx context.Context, id string) (*model.User, error)
}

type MutationResolver interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*model.User, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}
```

#### schema.go

リゾルバーのスタブ実装が生成されます。ここに実際のビジネスロジックを実装します（Step 4、5で実装）。

## 動作確認

### スキーマの構文チェック

スキーマファイルに構文エラーがないか確認します：

```bash
# gqlgenでバリデーション（コード生成時に自動的にチェックされる）
gqlgen generate
```

### コード生成の実行

```bash
# コード生成を実行
go generate ./...
```

エラーがなければ、以下のファイルが生成されます：

- `internal/gql/generated/generated.go`
- `internal/gql/generated/models_gen.go`
- `internal/resolver/resolver.go`
- `internal/resolver/schema.go` (QueryとMutationのリゾルバーが含まれます)

### 生成されたファイルの確認

```bash
# 生成されたファイルを確認
ls -la internal/gql/generated/
ls -la internal/resolver/
```

### ビルドの確認

生成されたコードが正しくコンパイルできるか確認します：

```bash
go build ./...
```

## トラブルシューティング

### スキーマの構文エラー

- スキーマファイルの構文を確認
- `gqlgen generate`を実行してエラーメッセージを確認
- GraphQLの型システム（!、[]など）を確認

### Timeスカラー型のエラー

- `internal/gql/scalars.go`が正しく実装されているか確認
- `MarshalTime`と`UnmarshalTime`関数が存在するか確認

### 生成されたコードのコンパイルエラー

- `go.mod`に必要な依存関係が追加されているか確認
- `go mod tidy`を実行して依存関係を整理
- gqlgenのバージョンを確認

### リゾルバーファイルが生成されない

- `gqlgen.yml`の設定を確認
- `resolver.dir`が正しく設定されているか確認
- スキーマにQueryやMutationが定義されているか確認

## 次のステップ

スキーマ定義とコード生成が完了したら、[Step 3: データベース設定とモデル定義](step-3-database.md) に進みましょう。

## 補足説明

### GraphQLスキーマのベストプラクティス

1. **型名は大文字で始める**: `User`、`CreateUserInput`など
2. **フィールド名はcamelCase**: `createdAt`、`userId`など
3. **必須フィールドには!を使用**: `name: String!`
4. **配列は明示的に定義**: `[User!]!`（要素も配列も必須）

### gqlgenのコード生成の仕組み

1. **スキーマ解析**: `.graphqls`ファイルを解析
2. **型マッピング**: GraphQL型をGo型にマッピング
3. **コード生成**: テンプレートを使ってGoコードを生成
4. **カスタマイズ**: `gqlgen.yml`で生成内容をカスタマイズ

### スキーマの分割

大きなスキーマは複数のファイルに分割できます：

```yaml
schema:
  - internal/gql/schema.graphqls
  - internal/gql/user.graphqls
  - internal/gql/post.graphqls
```

gqlgenは、すべてのスキーマファイルを結合して処理します。

### カスタムモデルの使用

`gqlgen.yml`で、GraphQL型を既存のGo構造体にマッピングできます：

```yaml
models:
  User:
    model: go-graphql-sample/internal/model.User
```

これにより、GraphQLの`User`型が`internal/model.User`構造体を使用します（Step 3で実装）。
