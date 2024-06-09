package req

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Payload struct {
	TotalCount int `json:"totalCount"`
}

type Data struct {
	Ar      int     `json:"__ar"`
	Payload Payload `json:"payload"`
}

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

type SearchHistory struct {
	ID uint `gorm:"primarykey" json:"id,omitempty"`

	KeyWordId   uint `json:"keyWordId,omitempty"`
	GroupId     uint `gorm:"index:idx_keyword_group_id,unique;" json:"groupId,omitempty"`
	SearchCount uint `json:"searchCount,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func MakeRequest(url string) (*http.Response, error) {
	method := "POST"
	payload := strings.NewReader("__a=1&fb_dtsg=NAcN8gzvxP02iYCYKB2u8DH4XGGSlImW_3_t2m-QqJuBVTENYAtnD0g%3A45%3A1700793422")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,pt;q=0.8")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add(
		"cookie",
		"m_ls=%7B%22c%22%3A%7B%221%22%3A%22HCwAABbKwg0WlvqmuAgTBRbe4vOdk70tAA%22%2C%222%22%3A%22GSwVQBxMAAAWABb25t3GDBYAABV-HEwAABYAFoLn3cYMFgAAFigA%22%7D%2C%22d%22%3A%2204ada221-ff8c-4e34-858e-2290f771195c%22%2C%22s%22%3A%220%22%2C%22u%22%3A%22raydvk%22%7D; datr=agjkZMjFmfE67fxpgGx8H5sT; sb=8MQyZVYLK-FBdlP_IwJlMLQ6; c_user=100006600013999; ps_n=1; ps_l=1; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1716952782453%2C%22v%22%3A1%7D; xs=45%3AT_J_yRyfD1fbzA%3A2%3A1700793422%3A-1%3A10129%3A%3AAcWcg8lWgsb8GNTu0KEWT-sQgaRHHVlFTwKskf2FFSI; fr=1HTMXeLSSjpjSH1jR.AWVDExA0RUJiymUTauK4FqiKQbQ.BmYQ6K..AAA.0.0.BmYQ6K.AWVZcDSxEcM; usida=eyJ2ZXIiOjEsImlkIjoiQXNldHBpbWw5cDdlZiIsInRpbWUiOjE3MTc5NTM2NDZ9; wd=1472x857",
	)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a requisição: %v", err)
	}

	return res, nil
}

func ParseResponse(res *http.Response) (Data, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Data{}, fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}
	defer res.Body.Close()

	// Remover a parte `for (;;);` para obter o JSON válido
	cleanJsonString := strings.TrimPrefix(string(body), "for (;;);")

	// Definir a estrutura de dados
	var data Data

	// Parsear o JSON
	err = json.Unmarshal([]byte(cleanJsonString), &data)
	if err != nil {
		return data, fmt.Errorf("erro ao fazer o parsing do JSON: %v", err)
	}

	return data, nil
}

func MakeUrl(search string) string {
	baseurl := "https://www.facebook.com/ads/library/async/search_ads/?q=%22"

	searchquote := strings.Replace(search, " ", "%20", -1)

	return baseurl + searchquote + "%22&count=30&countries%5C[0%5C]=BR&search_type=keyword_exact_phrase"

}

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

	var dado Data

	json.Unmarshal([]byte(jsonContent), &dado)

	return SaveDataInDb(dado, db)
}

func SaveDataInDb(dado Data, db *gorm.DB) error {

	err := db.AutoMigrate(&SearchHistory{})
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	var rows []KeyWord

	rows, err = GetAllDataFromKeywordTable(db)

	for i := 0; i < len(rows); i++ {

		save := SearchHistory{
			KeyWordId:   rows[i].ID,
			GroupId:     rows[i].GroupId,
			SearchCount: uint(dado.Payload.TotalCount),
		}
		err := db.Create(&save).Error
		if err != nil {
			fmt.Errorf("error saving data to database: %v", err)
		}

	}

	if err != nil {
		return fmt.Errorf("error saving data to database: %v", err)
	}

	return nil

}

func GetAllDataFromKeywordTable(db *gorm.DB) ([]KeyWord, error) {
	var keyw []KeyWord

	err := db.Where("is_active = ?", true).Where("deleted_at IS NULL").Find(&keyw).Error
	if err != nil {
		return nil, fmt.Errorf("error getting data from database: %v", err)
	}

	fmt.Println(keyw)
	return keyw, nil
}
