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

	// ポート番号の取得（環境変数から取得、デフォルトは8080）
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
