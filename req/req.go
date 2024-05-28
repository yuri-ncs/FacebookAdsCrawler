package req

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Payload struct {
	TotalCount int `json:"totalCount"`
}

type Data struct {
	Ar      int     `json:"__ar"`
	Payload Payload `json:"payload"`
}

func MakeRequest(url string) (*http.Response, error) {
	method := "POST"
	payload := strings.NewReader("__aaid=0&__user=100006600013999&__a=1&__req=2&__hs=19870.BP%3ADEFAULT.2.0..0.0&dpr=2&__ccg=EXCELLENT&__rev=1013776097&__s=lhv54i%3Afh9gzd%3Ac704cs&__hsi=7373505932978384251&__dyn=7xe6Eiw_K9zo5ObwKBAgc9o2exu13wqojyUW3qi4EoxW4E7SewXwCwfW7oqx60Vo1upEK12wvk1bwbG78b87C2m3K2y11wBw5Zx62G3i1ywdl0Fw4Hwp8kwyx2cwAxq1yK1LwPxe3C0D8sKUbobEaUiyE725U4q0N8G0iS2S3qazo11E2XU6-1FwLw8O1pwr86C0gi1QwtU&__csr=&fb_dtsg=NAcNtXwpvpsQtRrysnMoxbQhOMpLs2_WJHEsk6a7v1Fv_pfKv2VJO5w%3A45%3A1700793422&jazoest=25818&lsd=qo0BXBOiVCYnSbBj9qT9LH&__spin_r=1013776097&__spin_b=trunk&__spin_t=1716778132&__jssesw=1")

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
		"m_ls=%7B%22c%22%3A%7B%221%22%3A%22HCwAABbKwg0WlvqmuAgTBRbe4vOdk70tAA%22%2C%222%22%3A%22GSwVQBxMAAAWABb25t3GDBYAABV-HEwAABYAFoLn3cYMFgAAFigA%22%7D%2C%22d%22%3A%2204ada221-ff8c-4e34-858e-2290f771195c%22%2C%22s%22%3A%220%22%2C%22u%22%3A%22raydvk%22%7D; datr=agjkZMjFmfE67fxpgGx8H5sT; sb=8MQyZVYLK-FBdlP_IwJlMLQ6; c_user=100006600013999; usida=eyJ2ZXIiOjEsImlkIjoiQXNjMjE0bmhicWJydyIsInRpbWUiOjE3MTMzMDMwOTV9; ps_n=1; ps_l=1; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1713971290512%2C%22v%22%3A1%7D; dpr=1; xs=45%3AT_J_yRyfD1fbzA%3A2%3A1700793422%3A-1%3A10129%3A%3AAcV9pzWt35cKJ9r_ucSkBtF14Q_G_ewlWJZ4iTxw45w; fr=1XEVkxoj0gd0ovo1Q.AWU9GGO6GVPnI3diPd_Y_r7gsHM.BmU_SK..AAA.0.0.BmU_SK.AWVi7rhqAZc; wd=977x867",
	)
	req.Header.Add("dnt", "1")
	req.Header.Add("origin", "https://www.facebook.com")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add(
		"referer",
		"https://www.facebook.com/ads/library/?active_status=all&ad_type=all&country=BR&q=%22software%20house%22&search_type=keyword_exact_phrase&media_type=all",
	)
	req.Header.Add("sec-ch-prefers-color-scheme", "dark")
	req.Header.Add("sec-ch-ua", `"Chromium";v="125", "Not.A/Brand";v="24"`)
	req.Header.Add("sec-ch-ua-full-version-list", `"Chromium";v="125.0.6422.112", "Not.A/Brand";v="24.0.0.0"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-model", `""`)
	req.Header.Add("sec-ch-ua-platform", "macOS")
	req.Header.Add("sec-ch-ua-platform-version", "14.4.1")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	)
	req.Header.Add("x-asbd-id", "129477")
	req.Header.Add("x-fb-lsd", "qo0BXBOiVCYnSbBj9qT9LH")

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

	return baseurl + searchquote + "%22&v=550c25&session_id=8cd2ae4f-4ec2-4757-a6e5-416ff50ae254&count=30&active_status=all&ad_type=all&countries[0]=BR&media_type=all&search_type=keyword_exact_phrase"
}
