package jobs

import (
	"fmt"
	"github.com/abx-software/spyron-ads-crawler/database"
	"github.com/abx-software/spyron-ads-crawler/req"
	"gorm.io/gorm"
)

type Scraper struct {
	Database *gorm.DB
}

func NewScraper(database *gorm.DB) *Scraper {
	return &Scraper{Database: database}
}

func (scraper *Scraper) ScrapeAll() {
	rows, err := req.GetAllDataFromKeywordTable(scraper.Database)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, keyword := range rows {
		fmt.Println("Row: ", keyword.KeyWord, i)
		if keyword.KeyWord == "" {
			fmt.Println("Keyword is empty")
			continue
		}

		scraper.ScrapeOne(keyword)
	}
}

func (scraper *Scraper) ScrapeOne(keyword database.KeyWord) *database.SearchHistory {
	url, err := req.MakeUrl(keyword.KeyWord)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	res, err := req.MakeRequest(url)
	if err != nil {

		fmt.Println(err)
		return nil
	}

	data, err := req.ParseResponse(res)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	searchHistory := scraper.createSearchHistory(keyword, data)
	return &searchHistory
}

func (scraper *Scraper) createSearchHistory(keyWord database.KeyWord, facebookData database.Data) database.SearchHistory {
	fmt.Printf("%s - Payload TotalCounts: %d %d\n", keyWord.KeyWord, facebookData.Payload.TotalCount)

	search := database.SearchHistory{
		KeyWordId:   keyWord.ID,
		GroupId:     keyWord.GroupId,
		SearchCount: uint(facebookData.Payload.TotalCount),
	}

	err := scraper.Database.Create(&search).Error
	if err != nil {
		fmt.Errorf("error saving data to Database: %v", err)
	}

	if err != nil {
		fmt.Println(err)
	}

	return search
}
