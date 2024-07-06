package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"teste123/database"
	"teste123/proxy"
	"teste123/req"
)

func main() {

	db, err := database.DatabaseOpen()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Database connected")

	proxy.Initialize()
	fmt.Println("Initialized the proxy's")

	c := cron.New()

	fmt.Println("Cron job started")

	// Run every 4 hours
	c.AddFunc(
		"* * * * *", func() {

			fmt.Println("Running cron job")

			rows, err := req.GetAllDataFromKeywordTable(db)
			if err != nil {
				fmt.Println(err)
				return
			}

			for i, row := range rows {

				url, err := req.MakeUrl(row.KeyWord)
				if err != nil {
					fmt.Println(err)
					return
				}

				res, err := req.MakeRequest(url)

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
