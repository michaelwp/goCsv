/**************************************
Author: MichaelWP
@202407
v2.0.0
**************************************/

package gocsv

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// Request is data to process
type Request struct {
	WindowSize int64
	FilePath   string
	Headers    []string
	Data       [][]string
}

type newRequest struct {
	Input, Output chan []string
	File          *os.File
	WindowSize    int64
}

type GoCsv struct {
	WindowSize   int64
	Input        chan []string
	Output       chan []string
	File         *os.File
	WindowValues [][]string
	Position     int64
	Ctx          context.Context
}

func newGoCsv(ctx context.Context, request *newRequest) *GoCsv {
	return &GoCsv{
		Input:        request.Input,
		Output:       request.Output,
		File:         request.File,
		WindowSize:   request.WindowSize,
		WindowValues: make([][]string, request.WindowSize),
		Ctx:          ctx,
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
			gc.Position = (gc.Position + 1) % gc.WindowSize
			if gc.Position == 0 {
				gc.Output <- gc.GenerateRowBatch()
			}
		}
	}
}

func (gc *GoCsv) GenerateRowBatch() []string {
	rowBatch := make([]string, gc.WindowSize)
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
func Generate(ctx context.Context, request *Request) error {
	input := make(chan []string)
	output := make(chan []string)

	file, err := os.Create(request.FilePath + ".csv")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	defer file.Close()

	goCsv := newGoCsv(ctx, &newRequest{
		Input:      input,
		Output:     output,
		File:       file,
		WindowSize: request.WindowSize,
	})

	if err := goCsv.AddHeader(request.Headers); err != nil {
		return fmt.Errorf("error add header: %v", err)
	}

	go goCsv.Worker()
	go goCsv.InputData(request.Data)

	if err := goCsv.AddData(); err != nil {
		goCsv.Cancel()
		return err
	}

	return nil
}
