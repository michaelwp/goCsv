package main

import (
	"fmt"
	"log"
	"strconv"

	gocsv "github.com/michaelwp/go-csv"
)

func main() {
	filePath := "output.csv"
	headers := []string{"No", "Name", "Age", "Occupation"}
	data := GenerateSampleData(1000000)

	err := gocsv.Generate(filePath, headers, data)
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
