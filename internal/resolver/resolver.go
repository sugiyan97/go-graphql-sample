package resolver

import (
	"go-graphql-sample/internal/database"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *gorm.DB
}

// NewResolver は新しいResolverインスタンスを作成します
func NewResolver() *Resolver {
	return &Resolver{
		DB: database.DB,
	}
}
