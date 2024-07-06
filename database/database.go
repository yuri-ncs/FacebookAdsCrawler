package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
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

func DatabaseOpen() (*gorm.DB, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_NAME")

	//dsn := "host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable"
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + DBName + " port=" + port

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&KeyWord{}, &SearchHistory{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func PopulateKeyword(db *gorm.DB) error {
	db.Create(&KeyWord{Name: "teste 1", KeyWord: "Software House", URL: "", GroupId: 1, IsActive: true})
	db.Create(&KeyWord{Name: "teste 2", KeyWord: "Apple", URL: "", GroupId: 1, IsActive: true})
	db.Create(&KeyWord{Name: "teste 3", KeyWord: "Samsung", URL: "", GroupId: 2, IsActive: true})
	db.Create(&KeyWord{Name: "teste 4", KeyWord: "Huawei", URL: "", GroupId: 3, IsActive: true})
	db.Create(&KeyWord{Name: "teste 5", KeyWord: "Dell", URL: "", GroupId: 5, IsActive: true})

	return nil
}
