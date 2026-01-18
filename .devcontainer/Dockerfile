FROM golang:1.21-bookworm

# 必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    git \
    curl \
    sqlite3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 作業ディレクトリの設定
WORKDIR /workspace

# 環境変数の設定
ENV GO111MODULE=on \
    CGO_ENABLED=1

# ユーザーの作成（vscodeユーザー）
RUN groupadd -r vscode && useradd -r -g vscode -m -s /bin/bash vscode

# ワークスペースの所有権を変更
RUN chown -R vscode:vscode /workspace

USER vscode

# Goのパス設定
ENV PATH=$PATH:/home/vscode/go/bin
