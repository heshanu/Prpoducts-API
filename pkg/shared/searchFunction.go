package shared

import (
	"strings"
	"sync"

	"github.com/heshanu/go-service/pkg/models"
)

var (
	wg sync.WaitGroup
)

func SearchBooks(productsList []models.Product, keyword string) []models.Product {
	var results []models.Product
	lowerKeyword := strings.ToLower(keyword)

	for _, product := range productsList {
		if strings.Contains(strings.ToLower(product.Name), lowerKeyword) ||
			strings.Contains(strings.ToLower(product.Description), lowerKeyword) {
			results = append(results, product)
		}
	}
	return results
}

func SearchBooksConcurrently(productsList []models.Product, keyword string) []models.Product {
	//channel is like one thread
	resultsChan := make(chan []models.Product)

	// Split the books into chunks
	chunkSize := len(productsList) / 4
	if chunkSize == 0 {
		chunkSize = 1
	}

	// Process each chunk concurrently
	//search each part seperatly like java threads
	for i := 0; i < len(productsList); i += chunkSize {
		end := i + chunkSize
		if end > len(productsList) {
			end = len(productsList)
		}

		wg.Add(1) // Increment the WaitGroup counter
		go func(chunk []models.Product) {
			defer wg.Done()                            // Decrement the WaitGroup counter when done
			resultsChan <- SearchBooks(chunk, keyword) // Send results to the channel
		}(productsList[i:end]) // Pass the chunk to the goroutine
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Merge results from all goroutines
	var results []models.Product
	for chunkResults := range resultsChan {
		results = append(results, chunkResults...)
	}

	//return product slice
	return results
}
