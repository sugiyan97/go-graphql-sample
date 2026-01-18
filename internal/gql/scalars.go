package gql

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// Time はGraphQLのTimeスカラー型を実装します
// このファイルは、gqlgenが生成するコードで使用されます

// MarshalTime はtime.TimeをGraphQLの文字列に変換します
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.MarshalTime(t)
}

// UnmarshalTime はGraphQLの文字列をtime.Timeに変換します
func UnmarshalTime(v interface{}) (time.Time, error) {
	return graphql.UnmarshalTime(v)
}
