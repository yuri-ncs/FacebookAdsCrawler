package req

import (
"strings"
	// "encoding/json"
	"fmt"
	"io"
	"net/http"
	
)

type FacebookAdsResponse struct {
  Ar              int     `json:"__ar"`
  Payload         Payload `json:"payload"`
  ForwardCursor   string  `json:"forwardCursor"`
  BackwardCursor  string  `json:"backwardCursor"`
  TotalCount      int     `json:"totalCount"`
  CollationToken  string  `json:"collationToken"`
}

// Payload struct represents the nested "payload" object
type Payload struct {
  IsResultComplete bool      `json:"isResultComplete"`
  Results          []interface{} `json:"json:\"results\"omitempty"` // Handle empty array gracefully
  Pageresults       []interface{} `json:"pageresults"`
}

func CreateReqBody(){

  url := "https://www.facebook.com/ads/library/async/search_ads/?q=redbull&v=c7e686&session_id=b237bdc5-3bc8-457c-82a5-c31ae867ba3b&countries[0]=BR"
  method := "POST"

  payload := strings.NewReader("__a=1&__req=2&__hs=19863.BP%3ADEFAULT.2.0..0.0&dpr=1&__rev=1013633365&__s=kyn74h%3Avh3h0m%3Awnnhig&__hsi=7371222370560544025&__dyn=7xeUmxa3-Q8zo5ObwKBAgc9o9E6u5U4e1FxebzEdF8ixy7EiwvoWdwJwCwfW7oqx60Vo1upEK12wvk1bwbG78b87C2m3K2y11wBz81s8hwGwQwoE2LwBgao884y0Mo6i588Egz898mwkE-U6-3e4Ueo2sxOXwJwKwHxaaws8nwhE2Lxiaw4JwJwSyES0gq0K-1LwqobU2cwmo6O1Fw44wt87u&lsd=AVo1jvVNRwY&jazoest=2989&__spin_r=1013633365&__spin_b=trunk&__spin_t=1716246449&__jssesw=1")

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

		

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
  req.Header.Add("Accept", "*/*")
  req.Header.Add("Accept-Language", "en-US,en;q=0.5")
  req.Header.Add("Accept-Encoding", "gzip, deflate, br")
  req.Header.Add("Referer", "https://www.facebook.com/ads/library/?active_status=all&ad_type=all&country=BR&q=software%20house&search_type=keyword_unordered&media_type=all")
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("X-FB-LSD", "AVo1jvVNRwY")
  req.Header.Add("X-ASBD-ID", "129477")
  req.Header.Add("Origin", "https://www.facebook.com")
  req.Header.Add("Alt-Used", "www.facebook.com")
  req.Header.Add("Connection", "keep-alive")
  req.Header.Add("Cookie", "fr=0gzYY3whcf2Mjtqpi..BmS8jM..AAA.0.0.BmS8jc.AWVCqQ52gc0; usida=eyJ2ZXIiOjEsImlkIjoiQXNkdDFhendiNGhrMCIsInRpbWUiOjE3MTYyNDI2NjV9; wd=1420x963; datr=zMhLZnFgITumABuAYLMUj9Bm; ps_l=1; ps_n=1")
  req.Header.Add("Sec-Fetch-Dest", "empty")
  req.Header.Add("Sec-Fetch-Mode", "cors")
  req.Header.Add("Sec-Fetch-Site", "same-origin")
  req.Header.Add("TE", "trailers")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }



		modifiedString := strings.Replace(string(body[:]), "for (;;);", "", -1)
	
		string := strings.Clone(modifiedString)

		fmt.Println(string)

		var fbres FacebookAdsResponse

		// json.Unmarshal([]byte(modifiedString), &fbres)
		


		
  fmt.Print(fbres)
}
