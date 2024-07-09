package req

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"teste123/database"
	"teste123/proxy"
)

type AdLibraryQuery struct {
	ActiveStatus  string
	AdType        string
	Country       string
	ViewAllPageID string
	SearchType    string
	MediaType     string
}

/*
MakeRequest is a function that makes a request to the given URL and a given configuration.
*/
func MakeRequest(urlString string) (string, error) {

	for i := 0; i < proxy.GetIpCount(); i++ {
		fmt.Printf("Using [%s] to make the request.\n", proxy.GetCurrentIp())
		client, err := proxy.GetClient()
		if err != nil {
			return "", fmt.Errorf("erro: %v", err)
		}

		reader := "__a=1&lsd=" + os.Getenv("X_FB_LSD")
		data := strings.NewReader(reader)

		req, err := http.NewRequest("POST", urlString, data)
		if err != nil {
			return "", fmt.Errorf("erro ao criar a requisição: %v", err)
		}

		req.Header.Add("User-Agent", os.Getenv("USER_AGENT"))
		req.Header.Add("X-FB-LSD", os.Getenv("X_FB_LSD"))
		req.Header.Add("Sec-Fetch-Site", "same-origin")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("erro ao fazer a requisição: %v", err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
		}
		defer res.Body.Close()

		// Remover a parte `for (;;);` para obter o JSON válido
		if !strings.Contains(string(body), "totalCount") {
			fmt.Printf("Failed using [%s], trying another...\n", proxy.GetCurrentIp())

			proxy.ChangeProxy()
			continue
		}

		return string(body), nil
	}
	return "", fmt.Errorf("tried every proxy")
}

/*
ParseResponse is a function that parses the response from the request.
*/
func ParseResponse(res string) (database.Data, error) {

	// Remover a parte `for (;;);` para obter o JSON válido
	cleanJsonString := strings.TrimPrefix(res, "for (;;);")
	//fmt.Println("Cleaned JSON string:", cleanJsonString)

	var data database.Data

	// Parsear o JSON
	err := json.Unmarshal([]byte(cleanJsonString), &data)
	if err != nil {
		fmt.Println(cleanJsonString)
		return data, fmt.Errorf("erro ao fazer o parsing do JSON: %v", err)
	}
	return data, nil
}

/*
MakeUrl is a function that makes the URL to be used in the request. deprecated
*/
func MakeUrl(rawUrl string) (string, error) {

	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	queryParams := parsedURL.Query()

	adLibraryQuery := AdLibraryQuery{
		ActiveStatus:  queryParams.Get("active_status"),
		AdType:        queryParams.Get("ad_type"),
		Country:       queryParams.Get("country"),
		ViewAllPageID: queryParams.Get("view_all_page_id"),
		SearchType:    queryParams.Get("search_type"),
		MediaType:     queryParams.Get("media_type"),
	}

	// Check for invalid PageID (already present)
	if adLibraryQuery.ViewAllPageID == "0" || adLibraryQuery.ViewAllPageID == "" {
		fmt.Println("Warning: Page Id is 0")
		return "", fmt.Errorf("Invalid PageID: 0")
	}

	// Check for empty variables

	for _, field := range []string{adLibraryQuery.ActiveStatus, adLibraryQuery.AdType, adLibraryQuery.Country, adLibraryQuery.SearchType, adLibraryQuery.MediaType} {
		if field == "" {
			switch field {
			case adLibraryQuery.ActiveStatus:
				adLibraryQuery.ActiveStatus = "all"
			case adLibraryQuery.Country:
				fmt.Println("Warning: Country is empty, using BR as default")
				adLibraryQuery.Country = "BR"
			case adLibraryQuery.SearchType:
				fmt.Println("Warning: Search Type is empty")
				adLibraryQuery.SearchType = "page"
			}
			break
		}
	}

	baseurl := "https://www.facebook.com/ads/library/async/search_ads/?"

	finalUrl := baseurl + "active_status=" + url.QueryEscape(adLibraryQuery.ActiveStatus) + "&countries[0]=" + url.QueryEscape(adLibraryQuery.Country) + "&view_all_page_id=" + url.QueryEscape(adLibraryQuery.ViewAllPageID) + "&search_type=" + url.QueryEscape(adLibraryQuery.SearchType)
	return finalUrl, nil
}

// OpenFile this function below is used when the pc is blocked from fb LOL
func OpenFile() (*os.File, error) {
	file, err := os.Open("/home/yuri/Documents/FacebookAdsCrawler/req/fb.txt")
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}

	return file, nil
}

func ReadDataFromFile(file *os.File, db *gorm.DB) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	defer file.Close()

	jsonContent := strings.TrimPrefix(string(data), "for (;;);")

	var dado database.Data

	err = json.Unmarshal([]byte(jsonContent), &dado)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}
	return SaveDataInDb(dado, db)
}

func SaveDataInDb(dado database.Data, db *gorm.DB) error {

	return nil

}

func GetAllDataFromKeywordTable(db *gorm.DB) ([]database.KeyWord, error) {
	var keyw []database.KeyWord

	err := db.Where("is_active = ?", true).Where("deleted_at IS NULL").Find(&keyw).Error
	if err != nil {
		return nil, fmt.Errorf("error getting data from database: %v", err)
	}

	return keyw, nil
}
