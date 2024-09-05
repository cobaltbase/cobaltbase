package ct

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
	DeviceName   string
	RefreshToken string `gorm:"unique"`
	AuthID       string
	Provider     string
}

func (Auth) TableName() string {
	return "_auth"
}

func (Session) TableName() string {
	return "_sessions"
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
