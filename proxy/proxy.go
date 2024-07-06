package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Proxy struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProxyAddress string `json:"proxy_address"`
	Port         int    `json:"port"`
	Valid        bool   `json:"valid"`
}

type Response struct {
	Count   int     `json:"count"`
	Results []Proxy `json:"results"`
}

var response Response
var select_proxy = 0

func Initialize() {
	// Define the request URL
	url := "https://proxy.webshare.io/api/v2/proxy/list/?mode=direct&page=1&page_size=25"

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the Authorization header
	req.Header.Add("Authorization", "Token "+os.Getenv("PROXY_TOKEN"))

	// Create an HTTP client
	client := &http.Client{}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	Parse(string(body))
	return
}

func Parse(res string) {
	// Parse the JSON response
	err := json.Unmarshal([]byte(res), &response)
	if err != nil {
		log.Fatalf("Error parsing Proxy JSON: %s", err)
	}

	// Print the parsed data
	fmt.Printf("Found [%d] proxy\n", response.Count)

}

func GetClient() (*http.Client, error) {
	proxyURL, err := url.Parse("http://" + os.Getenv("PROXY_USER") + ":" + os.Getenv("PROXY_PASS") + "@" + response.Results[select_proxy].ProxyAddress + ":" + strconv.Itoa(response.Results[select_proxy].Port))
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

	return client, nil
}

func GetIpCount() int {
	return response.Count
}

func ChangeProxy() {
	if select_proxy < 9 {
		select_proxy++
	} else {
		select_proxy = 0
	}
}

func GetCurrentIp() string {
	return response.Results[select_proxy].ProxyAddress
}
