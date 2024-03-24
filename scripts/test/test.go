package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Product struct represents the structure of each product
type Product struct {
	Name           string
	CategoriesName string
	CategoriesID   int64
}

func main() {
	// Open the CSV file
	file, err := os.Open("../../files/name.csv")
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader with comma as the field delimiter
	reader := csv.NewReader(file)
	reader.Comma = ';'
	// Read all records from CSV
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the CSV:", err)
		return
	}

	// Initialize slices
	var products []Product

	// Create a map to store category names and their corresponding IDs
	categoryIDMap := make(map[string]int)

	// Parse CSV records into Product struct and populate slices
	for i, record := range records {
		if i == 1 {
			continue
		}
		cID, _ := strconv.ParseInt(record[0], 10, 64)
		product := Product{
			Name:           record[2], // Assuming the name is in the first column
			CategoriesName: record[1], // Assuming the category name is in the second column
			CategoriesID:   cID,
		}

		// Append to slices
		products = append(products, product)

	}

	// Print the parsed data or use it as needed
	fmt.Println("Parsed Products:")
	for _, p := range products {
		fmt.Printf("%+v\n", p)
	}

	fmt.Println("\nCategories IDs:")
	fmt.Println(categoryIDMap)
}
