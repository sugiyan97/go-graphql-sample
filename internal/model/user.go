package model

import (
	"time"
)

// User はユーザー情報を表すデータモデルです
// GORMのモデルとして使用されます
type User struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`               // ユーザーID（UUID）
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`              // ユーザー名
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"` // メールアドレス（一意）
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`                           // 作成日時
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`                           // 更新日時
}

// TableName はテーブル名を指定します
func (User) TableName() string {
	return "users"
}
