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
	"teste123/config"
	"teste123/database"
)

/*
MakeRequest is a function that makes a request to the given URL and a given configuration.
*/
func MakeRequest(urlString string, config config.RequestConfig) (*http.Response, error) {
	reader := "__a=1&" + "lsd=" + os.Getenv("X_FB_LSD")
	data := strings.NewReader(reader)

	servers := []struct {
		IP   string
		Port string
	}{
		{"45.248.55.75", "6661"},
		{"192.53.142.178", "5875"},
		{"208.73.42.46", "9057"},
		{"192.46.185.172", "5862"},
		{"130.180.235.96", "5816"},
		{"63.246.130.233", "6434"},
		{"130.180.234.239", "7462"},
		{"63.246.130.53", "6254"},
		{"130.180.235.70", "5790"},
		{"45.196.61.145", "6183"},
	}

	actualServer := servers[2]

	// Defina a URL do proxy com autenticação
	proxyURL, err := url.Parse("http://euhjdfiy:qiy3k4qsrdcw@" + actualServer.IP + ":" + actualServer.Port)
	if err != nil {
		fmt.Println("Erro ao analisar a URL do proxy:", err)
		panic("proxy")
	}

	// Crie um transporte HTTP personalizado
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	fmt.Println(transport)

	// Crie um cliente HTTP usando o transporte personalizado
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(config.Method, urlString, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	req.Header.Add("User-Agent", os.Getenv("USER_AGENT"))
	req.Header.Add("X-FB-LSD", os.Getenv("X_FB_LSD"))
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a requisição: %v", err)
	}

	return res, nil
}

/*
ParseResponse is a function that parses the response from the request.
*/
func ParseResponse(res *http.Response) (database.Data, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return database.Data{}, fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}
	defer res.Body.Close()

	// Remover a parte `for (;;);` para obter o JSON válido
	cleanJsonString := strings.TrimPrefix(string(body), "for (;;);")
	//fmt.Println("Cleaned JSON string:", cleanJsonString)

	var data database.Data

	// Parsear o JSON
	err = json.Unmarshal([]byte(cleanJsonString), &data)
	if err != nil {
		fmt.Println(cleanJsonString)
		return data, fmt.Errorf("erro ao fazer o parsing do JSON: %v", err)

	}

	return data, nil
}

/*
MakeUrl is a function that makes the URL to be used in the request.
*/
func MakeUrl(search string, config config.RequestConfig) string {
	searchQuote := strings.Replace(strings.ToLower(search), " ", "%20", -1)

	endUrl := "&countries[0]=" + config.Countries

	finalUrl := config.BaseURL + searchQuote + endUrl
	fmt.Println(finalUrl)
	return finalUrl

}

// this function below is used when the pc is blocked from fb LOL
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
