package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pdf "github.com/chinmay-sawant/gopdfsuit-client"
)

const (
	baseURL    = "http://localhost:8080"
	sampleJSON = "sample.json"
	outputPath = "output_from_reader.pdf"
)

func main() {
	// Create a client with options
	client := pdf.NewClient(
		baseURL,
		pdf.WithTimeout(60*time.Second),
		pdf.WithMaxRetries(3),
		pdf.WithHeader("Authorization", "Bearer your-token-here"),
	)

	ctx := context.Background()

	fmt.Println("=== Reading from JSON file and sending to endpoint ===")

	// 1. Read document from JSON file
	doc, err := client.ReadFromFile(ctx, sampleJSON)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	fmt.Printf("Successfully loaded document: %s\n", doc.Title.Text)
	fmt.Printf("Document contains %d tables\n", len(doc.Tables))

	// 2. Send the document to the endpoint and save the result
	fmt.Printf("Sending document to %s...\n", baseURL)
	err = client.SendAndSave(ctx, doc, outputPath)
	if err != nil {
		// Handle connection errors gracefully for the sample
		fmt.Printf("Note: Request failed (%v).\n", err)
		fmt.Println("Make sure the PDF generation service is running at", baseURL)

		// For demonstration, we'll just print what would have happened
		if os.Getenv("DEBUG") == "true" {
			panic(err)
		}
		return
	}

	fmt.Printf("Success! PDF saved to: %s\n", outputPath)
}
