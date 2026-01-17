# Step 1: プロジェクト初期化とgqlgenセットアップ

## 学習目標

このStepでは、Goプロジェクトの初期化とgqlgenの設定を行います。以下のことを学びます：

- Goモジュールの初期化
- gqlgenの設定ファイル（gqlgen.yml）の理解
- プロジェクト構造の作成
- 必要な依存関係の追加方法

## 前提知識

- Goの基本的な理解
- Goモジュールの概念

## 実装内容

### 1. Goモジュールの初期化

Goモジュールは、Goプロジェクトの依存関係管理の仕組みです。`go.mod`ファイルで管理されます。

既に`go.mod`が存在する場合は、以下の内容になっています：

```go
module go-graphql-sample

go 1.21.13
```

もし`go.mod`が存在しない場合は、以下のコマンドで初期化します：

```bash
go mod init go-graphql-sample
```

### 2. gqlgen設定ファイルの作成

`gqlgen.yml`は、gqlgenがGraphQLスキーマからGoコードを生成する際の設定ファイルです。

主な設定項目：

- **schema**: GraphQLスキーマファイルの場所
- **exec**: 生成されるGraphQL実行コードの設定
- **model**: モデル生成の設定
- **resolver**: リゾルバーファイルの生成設定
- **models**: カスタムモデルのマッピング

#### 設定の詳細

```yaml
schema:
  - internal/gql/schema.graphqls  # スキーマファイルのパス

exec:
  filename: internal/gql/generated/generated.go  # 生成される実行コード
  package: generated  # パッケージ名

model:
  filename: internal/gql/generated/models_gen.go  # 生成されるモデルコード
  package: generated

resolver:
  layout: follow-schema  # スキーマに従ったレイアウト
  dir: internal/resolver  # リゾルバーディレクトリ
  package: resolver
  filename_template: "{name}.go"  # ファイル名のテンプレート

models:
  User:
    model: go-graphql-sample/internal/model.User  # カスタムモデルのマッピング
```

### 3. プロジェクト構造の作成

以下のディレクトリ構造を作成します：

```
go-graphql-sample/
├── cmd/
│   └── server/          # サーバーのエントリーポイント
├── internal/
│   ├── gql/             # GraphQLスキーマと生成コード
│   ├── model/           # データモデル
│   ├── resolver/        # リゾルバー実装
│   └── database/        # データベース接続
```

各ディレクトリの役割：

- **cmd/server/**: アプリケーションのエントリーポイント（main.go）
- **internal/gql/**: GraphQLスキーマ定義とgqlgenが生成するコード
- **internal/model/**: データベースモデル（GORMモデル）
- **internal/resolver/**: GraphQLリゾルバーの実装
- **internal/database/**: データベース接続と初期化

### 4. 依存関係の追加

このプロジェクトで使用する主要な依存関係は以下の通りです：

- **gqlgen**: GraphQLコード生成ツール
- **GORM**: ORMライブラリ
- **SQLiteドライバ**: SQLiteデータベース接続

依存関係は、実際にコードで使用する際に自動的に追加されますが、事前に追加する場合は以下のコマンドを実行します：

```bash
# gqlgen関連
go get github.com/99designs/gqlgen
go get github.com/99designs/gqlgen/graphql
go get github.com/99designs/gqlgen/graphql/handler
go get github.com/99designs/gqlgen/graphql/playground

# GORMとSQLite
go get gorm.io/gorm
go get gorm.io/driver/sqlite

# HTTPサーバー（標準ライブラリを使用する場合は不要）
go get github.com/gorilla/mux  # オプション
```

## 動作確認

### ディレクトリ構造の確認

以下のコマンドで、ディレクトリ構造が正しく作成されているか確認します：

```bash
tree -L 3 -I '.git|node_modules' .
```

または：

```bash
find . -type d -not -path '*/\.*' | sort
```

### gqlgenの確認

gqlgenが正しくインストールされているか確認します：

```bash
gqlgen version
```

### go.modの確認

`go.mod`ファイルが存在し、正しいモジュール名が設定されているか確認します：

```bash
cat go.mod
```

## トラブルシューティング

### gqlgenが見つからない

- `postCreateCommand`が正常に実行されたか確認
- 手動でインストール: `go install github.com/99designs/gqlgen@latest`
- PATHに`/home/vscode/go/bin`が含まれているか確認: `echo $PATH`

### go.modが存在しない

- `go mod init go-graphql-sample`を実行
- DevContainerの`postCreateCommand`が実行されたか確認

### ディレクトリが作成されない

- 権限を確認: `ls -la`
- 手動で作成: `mkdir -p cmd/server internal/gql internal/model internal/resolver internal/database`

## 次のステップ

プロジェクトの初期化が完了したら、[Step 2: GraphQLスキーマ定義とコード生成](step-2-schema.md) に進みましょう。

## 補足説明

### Goモジュールとは

Goモジュールは、Go 1.11以降で導入された依存関係管理の仕組みです。`go.mod`ファイルでモジュールの依存関係を管理し、`go.sum`ファイルで依存関係のチェックサムを管理します。

### gqlgenの役割

gqlgenは、GraphQLスキーマ（`.graphqls`ファイル）から、以下のGoコードを自動生成します：

- **GraphQL実行コード**: クエリ/ミューテーションの実行に必要なコード
- **モデルコード**: GraphQL型に対応するGo構造体
- **リゾルバーインターフェース**: 実装すべきリゾルバーメソッドの定義

これにより、型安全なGraphQLサーバーを実装できます。

### internalパッケージ

`internal/`ディレクトリは、Goの特別なディレクトリ名です。このディレクトリ内のパッケージは、同じモジュール内からのみインポートできます。これにより、内部実装を外部から隠蔽できます。

### プロジェクト構造のベストプラクティス

- **cmd/**: アプリケーションのエントリーポイントを配置
- **internal/**: 内部パッケージを配置（外部からインポート不可）
- **pkg/**: 外部から使用可能なパッケージを配置（このプロジェクトでは使用しない）
- **api/**: API定義を配置（このプロジェクトでは使用しない）
