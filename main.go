package main

import (
	"encoding/json"
	"fmt"
	"github.com/abx-software/spyron-ads-crawler/database"
	"github.com/abx-software/spyron-ads-crawler/jobs"
	"github.com/abx-software/spyron-ads-crawler/req"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	fmt.Println(123)
	db, err := database.DatabaseOpen()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(123)

	scraper := jobs.NewScraper(db)

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

	// create cron jobs
	c.AddFunc(
		period, func() {
			scraper.ScrapeAll()
		},
	)
	c.Start()
	setupHttpServer(scraper)
	select {}
}

func setupHttpServer(scraper *jobs.Scraper) {
	r := mux.NewRouter()
	r.HandleFunc("/scrape/{id}", scrapeKeyWord(scraper)) // Route with a parameter

	fmt.Println("Starting server at port :" + os.Getenv("HTTP_SERVER_PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("HTTP_SERVER_PORT"), r); err != nil {
		fmt.Println(err)
	}
}

func scrapeKeyWord(scraper *jobs.Scraper) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		intId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "ID should be an integer", http.StatusBadRequest)
			return
		}

		keyWord, err := req.GetKeywordById(scraper.Database, uint(intId))
		if err != nil {
			http.Error(w, "Cannot find keyword", http.StatusNotFound)
			return
		}

		searchHistory := scraper.ScrapeOne(keyWord)
		fmt.Printf("%s - %s scraped, found %d", keyWord.Name, keyWord.KeyWord, searchHistory.SearchCount)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(*searchHistory)
	}
}
