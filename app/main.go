package main

import (
	"fmt"
	"github.com/robfig/cron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"teste123/database"
	"teste123/req"
)

func main() {

	dsn := "host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable"

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&database.KeyWord{}, &database.SearchHistory{})

	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()

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

				search := database.SearchHistory{
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
