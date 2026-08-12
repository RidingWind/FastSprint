package model

import "github.com/fastsprint/common/models"

type User struct {
	models.BaseModel
	Username    string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email       string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	DisplayName string `gorm:"size:100" json:"display_name"`
	AvatarURL   string `gorm:"size:500" json:"avatar_url"`
	Status      string `gorm:"size:20;default:active" json:"status"`
	Roles       []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "user_service.users"
}

type Role struct {
	models.BaseModel
	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

func (Role) TableName() string {
	return "user_service.roles"
}

type UserRole struct {
	models.BaseModel
	UserID uint `gorm:"not null;uniqueIndex:idx_user_role" json:"user_id"`
	RoleID uint `gorm:"not null;uniqueIndex:idx_user_role" json:"role_id"`
}

func (UserRole) TableName() string {
	return "user_service.user_roles"
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6,max=50"`
	DisplayName string `json:"display_name"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Status      string `json:"status"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

type UserQuery struct {
	models.PaginationQuery
	Keyword string `form:"keyword"`
	Status  string `form:"status"`
}
