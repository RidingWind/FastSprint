package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaginationQuery struct {
	Page     int `form:"page,default=1" json:"page"`
	PageSize int `form:"page_size,default=10" json:"page_size"`
}

type SortQuery struct {
	SortBy string `form:"sort_by" json:"sort_by"`
	Order  string `form:"order,default=desc" json:"order"`
}
