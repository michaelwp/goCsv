# goCsv

The `goCsv` package provides an easy-to-use interface for generating CSV files with concurrent processing. This package is particularly useful for handling large datasets and performing stream processing.

## Installation

To install the package, use `go get`:

```bash
go get github.com/michaelwp/goCsv/v2
```

## Usage

### Generate Function

The `Generate` function creates a CSV file with the given headers and data.

#### Example:

```go
package main

import (
	"context"
	"log"
	"strconv"
	"time"

	gocsv "github.com/michaelwp/goCsv/v2"
)

func main() {
	windowSize := int64(100)
	filePath := "output"
	headers := []string{"No", "First Name", "Last Name", "Gender", "Age", "Occupation", "Address", "Remark"}
	rowNumbers := 100000000
	data := GenerateSampleData(rowNumbers)
	second := 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), second)
	defer cancel()

	err := gocsv.Generate(ctx, &gocsv.Request{
		WindowSize: windowSize,
		FilePath:   filePath,
		Headers:    headers,
		Data:       data,
	})
	if err != nil {
		log.Println("failed to create csv", err)
		return
	}

	log.Printf("successfully in generating %s.csv with %d data less than %v",
		filePath, rowNumbers, second)
}

func GenerateSampleData(rowNumber int) [][]string {
	var data = make([][]string, rowNumber)

	for i := 0; i < rowNumber; i++ {
		noStr := strconv.Itoa(i + 1)
		data[i] = []string{
			noStr,
			"John",
			"Smith",
			"Male",
			"30",
			"Engineer",
			"South Tangerang, Banten, Indonesia",
			"Has over than 15 years of work experience",
		}
	}

	return data
}
```
