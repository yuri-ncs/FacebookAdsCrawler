package main

import (
	"fmt"
	"github.com/abx-software/spyron-ads-crawler/database"
	"github.com/abx-software/spyron-ads-crawler/jobs"
	"github.com/abx-software/spyron-ads-crawler/req"
	"github.com/abx-software/spyron-ads-crawler/setups"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	setups.SetupEnv()
	db, err := database.DatabaseOpen()
	if err != nil {
		fmt.Println(err)
		return
	}

	scraper := jobs.NewScraper(db)

	c := cron.New()

	fmt.Println("Cron job started")

	period := os.Getenv("CRON_PERIOD")

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
	r := gin.Default()
	r.POST("/scrape/:id", scrapeKeyWord(scraper))

	port := os.Getenv("HTTP_SERVER_PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	fmt.Println("Starting server at port :" + port)
	if err := r.Run(":" + port); err != nil {
		fmt.Println(err)
	}
}

func scrapeKeyWord(scraper *jobs.Scraper) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}
		intId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should be an integer"})
			return
		}

		groupId := c.Query("group-id")
		var uintGroupId *uint
		intGroupId, err := strconv.Atoi(groupId)
		if err != nil && groupId != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "group ID should be an integer"})
			return
		}

		if intGroupId != 0 {
			parsed := uint(intGroupId)
			uintGroupId = &parsed
		}

		keyWord, err := req.GetKeywordById(scraper.Database, uint(intId), uintGroupId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cannot find keyword"})
			return
		}

		searchHistory := scraper.ScrapeOne(keyWord)
		fmt.Printf("%s - %s scraped, found %d", keyWord.Name, keyWord.KeyWord, searchHistory.SearchCount)

		c.JSON(http.StatusOK, searchHistory)
	}
}
