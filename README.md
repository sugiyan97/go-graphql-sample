# go-graphql-sample

GolangでGraphQLを学習するためのサンプルプロジェクトです。gqlgenを使用して、段階的にGraphQLの理解を深められる構成になっています。

## 参考

> https://zenn.dev/hsaki/books/golang-graphql

## 学習の進め方

このプロジェクトは、8つのStepに分かれており、各Stepを順番に進めることでGraphQLの基礎から実践まで学べます。

### Step 0: 環境構築

DevContainer環境のセットアップ
- [ドキュメント](docs/step-0-setup.md)

### Step 1: プロジェクト初期化とgqlgenセットアップ

Goモジュールの初期化とgqlgenの設定
- [ドキュメント](docs/step-1-initialization.md)

### Step 2: GraphQLスキーマ定義とコード生成

GraphQLスキーマの作成とgqlgenによるコード生成
- [ドキュメント](docs/step-2-schema.md)

### Step 3: データベース設定とモデル定義

SQLiteデータベースの設定とデータモデルの定義
- [ドキュメント](docs/step-3-database.md)

### Step 4: リゾルバー実装（Query）

クエリリゾルバーの実装とデータ取得
- [ドキュメント](docs/step-4-query-resolver.md)

### Step 5: リゾルバー実装（Mutation）

ミューテーションリゾルバーの実装とデータ操作
- [ドキュメント](docs/step-5-mutation-resolver.md)

### Step 6: サーバー起動と動作確認

HTTPサーバーの実装とGraphQL Playgroundでの動作確認
- [ドキュメント](docs/step-6-server.md)

### Step 7: テストの実装

GraphQLクエリ/ミューテーションのテスト
- [ドキュメント](docs/step-7-testing.md)

## クイックスタート

1. **DevContainerで開く**
   - VS Codeでこのプロジェクトを開く
   - コマンドパレット（`Cmd+Shift+P` / `Ctrl+Shift+P`）を開く
   - 「Dev Containers: Reopen in Container」を選択

2. **Step 0から開始**
   - [Step 0のドキュメント](docs/step-0-setup.md)を参照して環境構築を完了

3. **各Stepを順番に進める**
   - 各Stepのドキュメントを参照しながら実装
   - 各Step完了後に動作確認を行う

## プロジェクト構造

```
go-graphql-sample/
├── .devcontainer/          # DevContainer設定
├── cmd/
│   └── server/            # サーバーエントリーポイント
├── internal/
│   ├── gql/               # GraphQLスキーマと生成コード
│   ├── model/             # データモデル
│   ├── resolver/          # リゾルバー実装
│   └── database/          # データベース接続
├── docs/                  # 学習ドキュメント
└── README.md
```

## 使用技術

- **Go 1.21+**: プログラミング言語
- **gqlgen**: GraphQLコード生成ツール
- **GORM**: ORMライブラリ
- **SQLite**: データベース
- **DevContainers**: 開発環境の統一

## 学習ポイント

1. **GraphQLスキーマ定義**: schema.graphqlsでの型定義方法
2. **gqlgenコード生成**: スキーマからGoコードを自動生成する仕組み
3. **リゾルバー実装**: ビジネスロジックの実装方法
4. **データベース連携**: GORMを使ったCRUD操作
5. **テスト**: GraphQLクエリ/ミューテーションのテスト方法

## ライセンス

このプロジェクトのライセンスについては、[LICENSE](LICENSE)ファイルを参照してください。
