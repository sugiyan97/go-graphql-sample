# Step 0: 環境構築

## 学習目標

このStepでは、GraphQL学習用の開発環境をDevContainersで構築します。以下のことを学びます：

- DevContainersの基本概念
- Go開発環境のセットアップ
- 必要なツールのインストール

## 前提知識

- Dockerの基本的な理解
- VS Code（またはDevContainers対応エディタ）の使用経験

## 実装内容

### 1. DevContainer設定ファイル

`.devcontainer/devcontainer.json` は、DevContainerの動作を定義する設定ファイルです。

主な設定項目：

- **name**: コンテナの名前
- **build**: Dockerfileの場所とビルドコンテキスト
- **features**: 追加機能（Git、GitHub CLIなど）
- **customizations**: VS Codeの拡張機能や設定
- **forwardPorts**: ポート転送設定（8080番ポートを転送）
- **postCreateCommand**: コンテナ作成後に実行するコマンド（依存関係のダウンロードとgqlgenのインストール）

### 2. Dockerfile

`.devcontainer/Dockerfile` は、開発環境のDockerイメージを定義します。

含まれる内容：

- **ベースイメージ**: `golang:1.21-bookworm`（Go 1.21がインストールされたDebian Bookworm）
- **必要なパッケージ**: Git、curl、SQLite3など
- **ユーザー設定**: vscodeユーザーの作成
- **環境変数**: Go関連の環境変数設定

### 3. .gitignore

`.gitignore` は、Gitで管理しないファイルを指定します。

除外される主なファイル：

- バイナリファイル（.exe, .dllなど）
- テストカバレッジファイル
- IDE設定ファイル
- データベースファイル
- 生成されたコード（gqlgenで生成されるファイル）

## 動作確認

### DevContainerで開く

1. VS Codeでこのプロジェクトを開く
2. コマンドパレット（`Cmd+Shift+P` / `Ctrl+Shift+P`）を開く
3. 「Dev Containers: Reopen in Container」を選択
4. コンテナのビルドと起動を待つ（初回は数分かかります）

### 環境の確認

コンテナ内で以下のコマンドを実行して、環境が正しくセットアップされているか確認します：

```bash
# Goのバージョン確認
go version

# gqlgenのインストール確認
gqlgen version

# SQLiteの確認
sqlite3 --version

# 作業ディレクトリの確認
pwd
```

すべてのコマンドが正常に実行されれば、環境構築は完了です。

## トラブルシューティング

### コンテナが起動しない

- Dockerが起動しているか確認
- `.devcontainer/devcontainer.json` の構文エラーを確認
- VS CodeのDevContainers拡張機能がインストールされているか確認

### gqlgenが見つからない

- `postCreateCommand` が正常に実行されたか確認
- 手動でインストール: `go install github.com/99designs/gqlgen@latest`
- PATHに `/home/vscode/go/bin` が含まれているか確認

### ポート転送ができない

- VS Codeのポート転送設定を確認
- ファイアウォール設定を確認

## 次のステップ

環境構築が完了したら、[Step 1: プロジェクト初期化とgqlgenセットアップ](step-1-initialization.md) に進みましょう。

## 補足説明

### DevContainersとは

DevContainersは、開発環境をDockerコンテナとして定義し、チーム全体で同じ環境を共有できる仕組みです。以下のメリットがあります：

- **環境の統一**: 全員が同じ開発環境を使用
- **セットアップの簡素化**: 新メンバーもすぐに開発を開始できる
- **依存関係の分離**: ホストマシンの環境に影響を与えない

### 使用するツール

- **Go 1.21**: GraphQLサーバーの実装に使用
- **gqlgen**: GraphQLスキーマからGoコードを生成するツール
- **SQLite3**: 学習用の軽量データベース
- **GORM**: GoのORMライブラリ（Step 3で使用）
