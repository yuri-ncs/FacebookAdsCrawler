package main

import (
	"fmt"
	"teste123/req"
)

func main() {

	searchquotes := []string{
		"software house",
		"encapsulado",
		"SASS",
		"notebook",
		"apple",
		"pampers",
	}

	for i, searchquote := range searchquotes {
		go func(number int, searchQuote string) {

			url := req.MakeUrl(searchQuote)

			res, err := req.MakeRequest(url)

			if err != nil {
				fmt.Println(err)
				return
			}

			data, err := req.ParseResponse(res)
			if err != nil {
				fmt.Println(err)
				return

			}

			// Imprimir a estrutura parseada
			//fmt.Printf("Ar: %d\n", data.Ar)
			fmt.Printf("%s - Payload TotalCounts: %d %d\n", searchQuote, data.Payload.TotalCount, number)
		}(i, searchquote)
	}
	select {}
}
