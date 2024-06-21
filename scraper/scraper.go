package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func InitializeScraper() {
	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create a timeout to avoid hanging indefinitely
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// URL to visit
	url := "https://www.facebook.com/ads/library/?active_status=all&ad_type=all&country=BR&q=&sort_data[direction]=desc&sort_data[mode]=relevancy_monthly_grouped&search_type=keyword_unordered&media_type=all"
	// Set the latest user agent string
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"

	// Set options for the headless browser
	var res string
	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(userAgent),
	}
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Run the tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second), // Adjust sleep time based on how long the content takes to load
		chromedp.WaitVisible(`div.x8t9es0.x1uxerd5.xrohxju.x108nfp6.xq9mrsl.x1h4wwuj.x117nqv4.xeuugli`, chromedp.ByQuery),
		chromedp.Text(`div.x8t9es0.x1uxerd5.xrohxju.x108nfp6.xq9mrsl.x1h4wwuj.x117nqv4.xeuugli`, &res, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print the extracted text
	fmt.Println("Text found:", res)
}
