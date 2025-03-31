package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Product represents the input product structure
type Product struct {
	ProductTitle  string  `json:"product_title"`
	IsPlainYogurt bool    `json:"is_plain_yogurt"`
	Confidence    float64 `json:"confidence"`
	IsNonfat      bool    `json:"is_nonfat"`
}

// TransformedProduct represents the output product structure
type TransformedProduct struct {
	ProductTitle string  `json:"product_title"`
	IsYogurt     bool    `json:"is_yogurt"`
	IsPlain      bool    `json:"is_plain"`
	IsNonfat     bool    `json:"is_nonfat"`
	Confidence   float64 `json:"confidence"`
}

// isYogurtProduct determines if a product is a yogurt based on its title
func isYogurtProduct(title string) bool {
	title = strings.ToLower(title)

	// Exclusion terms that indicate the product is not a yogurt
	exclusionTerms := []string{
		"kefir",
		"alternative",
		"dairy-free",
		"dairy free",
		"non-dairy",
		"almondmilk",
		"cashewmilk",
		"coconut",
		"powder puff",
		"cushion puff",
		"drinkable",
		"strainer",
		"base",
		"starter culture",
		"extract",
	}

	// Check for yogurt-related terms
	containsYogurtTerm := strings.Contains(title, "yogurt") ||
		strings.Contains(title, "skyr") ||
		strings.Contains(title, "yoghurt")

	// Check for exclusion terms
	containsExclusionTerm := false
	for _, term := range exclusionTerms {
		if strings.Contains(title, term) {
			containsExclusionTerm = true
			break
		}
	}

	// Return true only if it contains a yogurt term and no exclusion terms
	return containsYogurtTerm && !containsExclusionTerm
}

func main() {
	// Check command line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./transform-yogurt input.json output.json")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read the input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Parse the JSON data
	var products []Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed %d products\n", len(products))

	// Transform the data
	transformedProducts := make([]TransformedProduct, 0, len(products))
	yogurtCount := 0
	plainCount := 0

	for _, product := range products {
		isYogurt := isYogurtProduct(product.ProductTitle)
		isPlain := false
		if isYogurt {
			isPlain = product.IsPlainYogurt
		}

		// Update counts
		if isYogurt {
			yogurtCount++
		}
		if isPlain {
			plainCount++
		}

		// Create transformed product
		transformedProduct := TransformedProduct{
			ProductTitle: product.ProductTitle,
			IsYogurt:     isYogurt,
			IsPlain:      isPlain,
			IsNonfat:     product.IsNonfat,
			Confidence:   product.Confidence,
		}

		transformedProducts = append(transformedProducts, transformedProduct)
	}

	// Output statistics
	fmt.Printf("Total yogurt products: %d out of %d\n", yogurtCount, len(products))
	fmt.Printf("Total plain products: %d out of %d\n", plainCount, len(products))

	// Write the transformed data to the output file
	outputData, err := json.MarshalIndent(transformedProducts, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON output: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, outputData, 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Transformation complete! Data written to %s\n", outputFile)
}