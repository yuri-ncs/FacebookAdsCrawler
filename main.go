package main

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"teste123/database"
	"teste123/proxy"
	"teste123/req"
)

func main() {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatalf("failed to initialize Jaeger exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			"", // schemaURL
			attribute.String("service.name", "Ads-Crawler"),
		)),
	)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	// Create a tracer
	tracer := otel.Tracer("AdsCrawler")
	meter := otel.Meter("Ads-Crawler")

	cronCounter, err := meter.Int64Counter(
		"total.cron", // Metric name
		metric.WithDescription("Total cron"),
	)

	requestCounter, err := meter.Int64Counter(
		"total.request", // Metric name
		metric.WithDescription("Total Request"),
	)
	_, span := tracer.Start(context.Background(), "DB-Connect")

	db, err := database.DatabaseOpen()
	if err != nil {
		fmt.Println(err)
		return
	}

	span.End()

	fmt.Println("Database connected")

	proxy.Initialize()
	fmt.Println("Initialized the proxy's")

	c := cron.New()

	fmt.Println("Cron job started")

	// Run every 4 hours
	c.AddFunc(
		"*/60 * * * *", func() {

			fmt.Println("Running cron job")
			_, mainSpan := tracer.Start(context.Background(), "Cronjob")
			cronCounter.Add(context.Background(), 1)

			_, keywordGetSpan := tracer.Start(context.Background(), "KeywordGetInfo")
			rows, err := req.GetAllDataFromKeywordTable(db)
			if err != nil {
				fmt.Println(err)
				return
			}
			keywordGetSpan.End()

			for i, row := range rows {
				_, processSpan := tracer.Start(context.Background(), "ProcessTime")
				_, requestSpan := tracer.Start(context.Background(), "RequestTime")
				requestCounter.Add(context.Background(), 1)

				url, err := req.MakeUrl(row.KeyWord)
				if err != nil {
					fmt.Println(err)
					return
				}

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

				requestSpan.End()
				_, saveDbSpan := tracer.Start(context.Background(), "SaveKeywordTime")

				// Imprimir a estrutura parseada
				//fmt.Printf("Ar: %d\n", data.Ar)
				fmt.Printf("%s - Payload TotalCounts: %d %d\n", row.KeyWord, data.Payload.TotalCount, i)

				search := database.SearchHistory{
					KeyWordId:   row.ID,
					GroupId:     row.GroupId,
					SearchCount: uint(data.Payload.TotalCount),
				}

				err = db.Create(&search).Error
				if err != nil {
					fmt.Errorf("error saving data to database: %v", err)
				}

				if err != nil {
					fmt.Println(err)
				}
				saveDbSpan.End()
				processSpan.End()
			}
			fmt.Println("end")
			mainSpan.End()
		},
	)

	c.Start()

	select {}
}
