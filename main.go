package main

import (
	"fmt"
	"github.com/robfig/cron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"teste123/req"
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

func main() {

	dsn := "host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable"

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&KeyWord{}, &req.SearchHistory{})

	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
	//
	//file, _ := req.OpenFile()
	//
	//data := req.ReadDataFromFile(file, db)
	//
	//if data != nil {
	//	log.Fatal(data)
	//}

	//fmt.Println("Data saved successfully to the database.")

	c := cron.New()

	c.AddFunc(
		"@every 4h", func() {

			rows, err := req.GetAllDataFromKeywordTable(db)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Running cron job")

			for i, row := range rows {

				url := req.MakeUrl(row.KeyWord)

				res, err := req.MakeRequest(url)

				fmt.Println(res.Body)

				if err != nil {
					fmt.Println(err)
					return
				}

				data, err := req.ParseResponse(res)
				if err != nil {
					fmt.Println(err)
					return

				}

				// Imprimir a estrutura parseada
				//fmt.Printf("Ar: %d\n", data.Ar)
				fmt.Printf("%s - Payload TotalCounts: %d %d\n", row.KeyWord, data.Payload.TotalCount, i)

				search := req.SearchHistory{
					KeyWordId:   row.ID,
					GroupId:     row.GroupId,
					SearchCount: uint(data.Payload.TotalCount),
				}

				err = db.Create(&search).Error
				if err != nil {
					fmt.Errorf("error saving data to database: %v", err)
				}

				if err != nil {
					fmt.Println(err)
				}
			}

		},
	)

	c.Start()

	select {}
}
