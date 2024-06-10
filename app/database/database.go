package database

import (
	"gorm.io/gorm"
	"time"
)

type KeyWord struct {
	ID       uint   `gorm:"primarykey" json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	KeyWord  string `gorm:"index:idx_keyword_group_id,unique;" json:"keyWord,omitempty"`
	URL      string `json:"url"`
	GroupId  uint   `gorm:"index:idx_keyword_group_id,unique;" json:"groupId,omitempty"`
	IsActive bool   `json:"isActive,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

type Payload struct {
	TotalCount int `json:"totalCount"`
}

type Data struct {
	Ar      int     `json:"__ar"`
	Payload Payload `json:"payload"`
}

type SearchHistory struct {
	ID uint `gorm:"primarykey" json:"id,omitempty"`

	KeyWordId   uint `json:"keyWordId,omitempty"`
	GroupId     uint `gorm:"index:idx_keyword_group_id;" json:"groupId,omitempty"`
	SearchCount uint `json:"searchCount,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
