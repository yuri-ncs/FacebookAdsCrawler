package config

type RequestConfig struct {
	Count          string `json:"count"`
	Countries      string `json:"countries"`
	SearchType     string `json:"search_type"`
	BaseURL        string `json:"base_url"`
	AcceptLanguage string `json:"accept_language"`
	ContentType    string `json:"content_type"`
	Method         string `json:"method"`
	Accept         string `json:"accept"`
}

func LoadRequestConfig() (RequestConfig, error) {

	return RequestConfig{
		Count:          "30",
		Countries:      "BR",
		SearchType:     "keyword_exact_phrase",
		BaseURL:        "https://www.facebook.com/ads/library/async/search_ads/?q=",
		AcceptLanguage: "en-US,en;q=0.9,pt;q=0.8",
		ContentType:    "application/x-www-form-urlencoded",
		Method:         "POST",
		Accept:         "*/*",
	}, nil
}
