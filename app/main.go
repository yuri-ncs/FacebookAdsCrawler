package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"teste123/config"
	"teste123/database"
	"teste123/req"
)

func main() {

	requestConfig, err := config.LoadRequestConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := database.DatabaseOpen()

	c := cron.New()

	c.AddFunc(
		"@every 1h", func() {

			rows, err := req.GetAllDataFromKeywordTable(db)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Running cron job")

			for i, row := range rows {

				url := req.MakeUrl(row.KeyWord, requestConfig)

				res, err := req.MakeRequest(url, requestConfig)

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
