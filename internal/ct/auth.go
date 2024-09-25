package ct

import (
	"time"
)

type Auth struct {
	BaseModel
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string `gorm:"default:user"`
	Sessions []Session
	Verified bool
}

type Session struct {
	BaseModel
	UserAgent    string
	RefreshToken string `gorm:"uniqueIndex"`
	AuthID       string
	Provider     string
}

func (Auth) TableName() string {
	return "_auth"
}

func (Session) TableName() string {
	return "_sessions"
}

func (OTP) TableName() string {
	return "_otps"
}

func (OauthConfig) TableName() string {
	return "_oauth_config"
}
func (SMTPConfig) TableName() string {
	return "_smtp_config"
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SMTPConfig struct {
	BaseModel
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	FromName string `json:"from_name"`
}

type OTP struct {
	UpdatedAt time.Time
	Email     string `json:"email" gorm:"primarykey"`
	OTP       string `json:"otp"`
}

type OauthConfig struct {
	BaseModel
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Provider     string `json:"provider" gorm:"uniqueIndex"`
}
