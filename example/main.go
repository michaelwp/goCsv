package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	gocsv "github.com/michaelwp/go-csv"
)

func main() {
	filePath := "output"
	headers := []string{"No", "Name", "Age", "Occupation"}
	data := GenerateSampleData(1000000)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := gocsv.Generate(ctx, filePath, headers, data)
	if err != nil {
		log.Println("failed to create csv", err)
		return
	}

	fmt.Println("csv generated successfully")
}

func GenerateSampleData(rowNumber int) [][]string {
	var data = make([][]string, rowNumber)

	for i := 0; i < rowNumber; i++ {
		noStr := strconv.Itoa(i + 1)
		data[i] = []string{noStr, "John Doe", "30", "Engineer"}
	}

	return data
}
