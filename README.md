# go-csv

The `go-csv` package provides an easy-to-use interface for generating CSV files with concurrent processing. This package is particularly useful for handling large datasets and performing stream processing.

## Installation

To install the package, use `go get`:

```bash
go get github.com/michaelwp/go-csv
```

## Usage

### Generate Function

The `Generate` function creates a CSV file with the given headers and data.

#### Example:

```go
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
