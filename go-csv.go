/**************************************
Author: MichaelWP
@202407
**************************************/

package go_csv

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

const WindowSize = 5

type GoCsv struct {
	Input        chan []string
	Output       chan []string
	File         *os.File
	WindowValues [int64(WindowSize)][]string
	Position     int
	Ctx          context.Context
}

func NewGoCsv(ctx context.Context, input, output chan []string, file *os.File) *GoCsv {
	return &GoCsv{
		Input:  input,
		Output: output,
		File:   file,
		Ctx:    ctx,
	}
}

func (gc *GoCsv) Worker() {
	for {
		select {
		case <-gc.Ctx.Done():
			close(gc.Output)
			return
		case val, ok := <-gc.Input:
			if !ok {
				close(gc.Output)
				return
			}

			gc.WindowValues[gc.Position] = val
			gc.Position = (gc.Position + 1) % WindowSize
			if gc.Position == 0 {
				gc.Output <- gc.GenerateRowBatch()
			}
		}
	}
}

func (gc *GoCsv) GenerateRowBatch() []string {
	rowBatch := make([]string, WindowSize)
	for i, val := range gc.WindowValues {
		rowBatch[i] = strings.Join(val, ",")
	}

	return rowBatch
}

func (gc *GoCsv) InputData(data [][]string) {
	defer close(gc.Input)
	for _, d := range data {
		gc.Input <- d
	}
}

func (gc *GoCsv) AddData() error {
	for line := range gc.Output {
		dataLine := strings.Join(line, "\n")
		if _, err := gc.File.WriteString(dataLine + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func (gc *GoCsv) AddHeader(headers []string) error {
	headerLine := strings.Join(headers, ",")
	if _, err := gc.File.WriteString(headerLine + "\n"); err != nil {
		return err
	}

	return nil
}

func (gc *GoCsv) Cancel() {
	ctx, cancel := context.WithTimeout(gc.Ctx, 5*time.Second)
	defer func() {
		fmt.Println("process cancelled")
		cancel()
	}()

	gc.Ctx = ctx
}

// Generate function to generate the csv file
func Generate(ctx context.Context, filePath string, headers []string, data [][]string) error {
	input := make(chan []string)
	output := make(chan []string)

	file, err := os.Create(filePath + ".csv")
	if err != nil {
		return fmt.Errorf("error creating file:", err)
	}

	defer file.Close()

	goCsv := NewGoCsv(ctx, input, output, file)

	if err := goCsv.AddHeader(headers); err != nil {
		return fmt.Errorf("error add header:", err)
	}

	go goCsv.Worker()
	go goCsv.InputData(data)

	if err := goCsv.AddData(); err != nil {
		goCsv.Cancel()
		return err
	}

	return nil
}
