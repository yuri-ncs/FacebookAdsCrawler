package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type KeyWord struct {
	ID uint `gorm:"primarykey" json:"id,omitempty"`

	Name     string `json:"name"`
	URL      string `json:"url"`
	KeyWord  string `gorm:"index:idx_keyword_group_id_active,unique;" json:"keyWord"`
	GroupId  uint   `gorm:"index:idx_keyword_group_id_active,unique;" json:"groupId"`
	IsActive *bool  `gorm:"index:idx_keyword_group_id_active,unique;" json:"isActive"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

type SearchHistory struct {
	ID uint `gorm:"primarykey" json:"id,omitempty"`

	KeyWordId uint    `json:"keyWordId"`
	KeyWord   KeyWord `gorm:"foreignKey:key_word_id;references:ID" json:"keyWord"`

	GroupId     uint `json:"groupId"`
	SearchCount uint `json:"searchCount"`

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

func DatabaseOpen() (*gorm.DB, error) {

	host := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")
	DBName := os.Getenv("POSTGRES_DBNAME")
	user := os.Getenv("POSTGRES_USER")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	timezone := os.Getenv("POSTGRES_TIMEZONE")

	//dsn := "host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable"
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + DBName + " port=" + port + " TimeZone=" + timezone + " sslmode=" + sslmode

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func PopulateKeyword(db *gorm.DB) error {
	True := true

	db.Create(&KeyWord{Name: "teste 1", KeyWord: "Software House", URL: "", GroupId: 1, IsActive: &True})
	db.Create(&KeyWord{Name: "teste 2", KeyWord: "Apple", URL: "", GroupId: 1, IsActive: &True})
	db.Create(&KeyWord{Name: "teste 3", KeyWord: "Samsung", URL: "", GroupId: 2, IsActive: &True})
	db.Create(&KeyWord{Name: "teste 4", KeyWord: "Huawei", URL: "", GroupId: 3, IsActive: &True})
	db.Create(&KeyWord{Name: "teste 5", KeyWord: "Dell", URL: "", GroupId: 5, IsActive: &True})

	return nil
}
