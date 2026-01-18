# Step 3: データベース設定とモデル定義

## 学習目標

このStepでは、SQLiteデータベースの設定とUserモデルの定義を行います。以下のことを学びます：

- GORMの基本的な使い方
- SQLiteデータベースへの接続
- データモデルの定義方法
- データベースマイグレーション
- gqlgenとデータモデルの連携

## 前提知識

- データベースの基本的な概念
- ORMの概念
- Goの構造体とタグ

## 実装内容

### 1. Userモデルの定義

`internal/model/user.go`でUserモデルを定義します。

#### モデルの構造

```go
type User struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
}
```

#### フィールドの説明

- **ID**: ユーザーID（UUID形式、主キー）
- **Name**: ユーザー名（必須）
- **Email**: メールアドレス（必須、一意制約）
- **CreatedAt**: 作成日時（自動設定）
- **UpdatedAt**: 更新日時（自動更新）

#### GORMタグの説明

GORMタグは、データベースのカラム定義を指定します：

- `gorm:"type:varchar(36);primaryKey"`: カラム型と主キー指定
- `gorm:"not null"`: NOT NULL制約
- `gorm:"uniqueIndex"`: 一意制約
- `json:"id"`: JSONシリアライゼーション時のフィールド名

#### TableNameメソッド

`TableName()`メソッドで、テーブル名を明示的に指定できます：

```go
func (User) TableName() string {
	return "users"
}
```

指定しない場合、GORMは構造体名を複数形にしてテーブル名とします（`User` → `users`）。

### 2. データベース接続の設定

`internal/database/database.go`でデータベース接続を管理します。

#### データベースの初期化

```go
func InitDB() error {
	dbPath := "data.db"
	if path := os.Getenv("DB_PATH"); path != "" {
		dbPath = path
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
```

#### 設定の説明

- **データベースパス**: デフォルトは`data.db`、環境変数`DB_PATH`で変更可能
- **GORM設定**: ログレベルを`Info`に設定（SQLクエリを表示）
- **自動マイグレーション**: テーブルを自動作成

#### マイグレーション

`AutoMigrate()`メソッドで、モデルに対応するテーブルを自動作成します：

```go
func autoMigrate() error {
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("failed to migrate User model: %w", err)
	}
	return nil
}
```

`AutoMigrate()`は以下の処理を行います：

- テーブルが存在しない場合は作成
- カラムが存在しない場合は追加
- インデックスを追加
- **既存のカラムやデータは削除しない**

### 3. gqlgen.ymlの更新

`gqlgen.yml`で、GraphQLの`User`型をGoの`model.User`にマッピングします：

```yaml
models:
  User:
    model: go-graphql-sample/internal/model.User
```

これにより、gqlgenはGraphQLの`User`型に対して、`internal/model.User`構造体を使用します。

### 4. 依存関係の追加

GORMとSQLiteドライバを追加する必要があります：

```bash
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go mod tidy
```

## 動作確認

### 依存関係のインストール

```bash
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go mod tidy
```

### データベース接続のテスト

簡単なテストコードでデータベース接続を確認できます：

```go
package main

import (
	"fmt"
	"log"
	"go-graphql-sample/internal/database"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	fmt.Println("Database connection successful!")
}
```

実行：

```bash
go run cmd/server/main.go
```

### データベースファイルの確認

`data.db`ファイルが作成されているか確認します：

```bash
ls -la data.db
```

### SQLiteでデータベースを確認

```bash
sqlite3 data.db

# テーブル一覧を表示
.tables

# usersテーブルの構造を確認
.schema users

# 終了
.quit
```

### gqlgenコードの再生成

Userモデルをマッピングしたので、gqlgenのコードを再生成します：

```bash
go generate ./...
```

これにより、GraphQLの`User`型が`model.User`を使用するように更新されます。

## トラブルシューティング

### データベースファイルが作成されない

- 書き込み権限を確認: `ls -la .`
- データベースパスを確認: `echo $DB_PATH`
- エラーメッセージを確認

### マイグレーションエラー

- GORMのバージョンを確認: `go list -m gorm.io/gorm`
- モデルの定義を確認（タグの構文エラーなど）
- 既存のデータベースファイルを削除して再作成: `rm data.db`

### gqlgenのコード生成エラー

- `gqlgen.yml`の`models`設定を確認
- `model.User`が正しくインポートできるか確認
- `go generate ./...`を再実行

### インポートエラー

- `go mod tidy`を実行して依存関係を整理
- `go get`で必要なパッケージをインストール

## 次のステップ

データベース設定とモデル定義が完了したら、[Step 4: リゾルバー実装（Query）](step-4-query-resolver.md) に進みましょう。

## 補足説明

### GORMとは

GORMは、Go言語用のORM（Object-Relational Mapping）ライブラリです。以下の機能を提供します：

- **自動マイグレーション**: モデルからテーブルを自動作成
- **クエリビルダー**: 型安全なクエリ構築
- **リレーション**: テーブル間の関係を定義
- **フック**: 保存前後の処理を定義

### SQLiteについて

SQLiteは、軽量なファイルベースのデータベースです。以下の特徴があります：

- **ファイルベース**: サーバー不要、単一ファイルで管理
- **軽量**: 小規模なアプリケーションに適している
- **学習用に最適**: セットアップが簡単

本番環境では、PostgreSQLやMySQLなどの本格的なデータベースを使用することを推奨します。

### データベース設計のベストプラクティス

1. **主キー**: 必ず主キーを設定（UUID推奨）
2. **インデックス**: 検索頻度の高いカラムにインデックスを設定
3. **一意制約**: 重複を許さない値には一意制約を設定
4. **NOT NULL**: 必須フィールドにはNOT NULL制約を設定
5. **タイムスタンプ**: 作成日時・更新日時を記録

### 環境変数による設定

データベースパスを環境変数で変更できます：

```bash
export DB_PATH=/path/to/database.db
go run cmd/server/main.go
```

本番環境では、環境変数や設定ファイルでデータベース接続情報を管理します。

### マイグレーションの注意点

`AutoMigrate()`は開発環境向けの機能です。本番環境では、以下のようなマイグレーションツールを使用することを推奨します：

- **golang-migrate**: データベースマイグレーションツール
- **sql-migrate**: SQLファイルベースのマイグレーション

### UUIDの生成

UserモデルのIDはUUIDを使用します。UUIDを生成するには、`github.com/google/uuid`パッケージを使用します：

```go
import "github.com/google/uuid"

id := uuid.New().String()
```

Step 4でリゾルバーを実装する際に使用します。
