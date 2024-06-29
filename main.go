package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"os"
	"teste123/database"
	"teste123/req"
	"time"
)

func main() {

	db, err := database.DatabaseOpen()
	if err != nil {
		fmt.Println(err)
		return
	}

	c := cron.New()

	fmt.Println("Cron job started")

	period := os.Getenv("CRON_PERIOD")

	if period == "" {
		period = "0 0 */4 * * *"
	}

	now := time.Now()
	fmt.Println("Hora atual:", now)

	// remove first 0 from period
	periodWithoutFirstZero := period[1:]

	schedule, err := cron.ParseStandard(periodWithoutFirstZero)
	if err != nil {
		fmt.Println("Erro ao analisar o cron schedule:", err)
		return
	}

	// Lista as próximas 10 execuções
	for i := 0; i < 10; i++ {
		now = schedule.Next(now)
		fmt.Println("Próxima execução:", now)
	}

	// Run every 4 hours
	c.AddFunc(
		period, func() {

			fmt.Println("Running cron job")

			rows, err := req.GetAllDataFromKeywordTable(db)
			if err != nil {
				fmt.Println(err)
				return
			}

			for i, row := range rows {

				fmt.Println("Row: ", row.KeyWord, i)

				if row.KeyWord == "" {
					fmt.Println("Keyword is empty")
					continue
				}

				url, err := req.MakeUrl(row.KeyWord)
				if err != nil {
					fmt.Println(err)
					continue
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
