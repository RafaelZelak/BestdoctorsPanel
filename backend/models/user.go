package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"unique;not null" json:"username"`
	Email        string     `gorm:"unique;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	FullName     string     `json:"full_name"`
	Role         string     `gorm:"default:user" json:"role"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
