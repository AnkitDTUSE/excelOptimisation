package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func copyTxtToXlsx() {
	startTime := time.Now()

	inputPath := "mock.data"
	outputPath := "test.xlsx"

	// Open input text file
	inFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to open input file '%s': %v\nNote: Please make sure the file exists or change the input path in the script.", inputPath, err)
	}
	defer inFile.Close()

	// Create new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing excel file: %v", err)
		}
	}()

	sheetName := "Sheet1"
	// Ensure the default sheet exists or create it
	index, err := f.GetSheetIndex(sheetName)
	if err != nil || index < 0 {
		index, err = f.NewSheet(sheetName)
		if err != nil {
			log.Fatalf("failed to create sheet: %v", err)
		}
	}

	// Create StreamWriter for memory-efficient writing of large datasets
	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		log.Fatalf("failed to create stream writer: %v", err)
	}

	scanner := bufio.NewScanner(inFile)
	rowCount := 0

	fmt.Println("Starting to copy data from text to Excel...")

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue // skip empty lines
		}

		// Detect delimiter: tab-separated or comma-separated
		var row []string
		if strings.Contains(line, "\t") {
			row = strings.Split(line, "\t")
		} else {
			row = strings.Split(line, ",")
		}

		// Convert []string to []interface{}
		rowData := make([]interface{}, len(row))
		for i, val := range row {
			rowData[i] = val
		}

		// Calculate cell coordinates (e.g. A1, A2...)
		cellName, err := excelize.CoordinatesToCellName(1, rowCount+1)
		if err != nil {
			log.Fatalf("failed to calculate cell coordinates for row %d: %v", rowCount+1, err)
		}

		if err := streamWriter.SetRow(cellName, rowData); err != nil {
			log.Fatalf("failed to write row %d: %v", rowCount+1, err)
		}

		rowCount++
		if rowCount%100000 == 0 {
			fmt.Printf("Processed %d rows...\n", rowCount)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading input file: %v", err)
	}

	// Flush and finalize stream writer
	if err := streamWriter.Flush(); err != nil {
		log.Fatalf("failed to flush stream writer: %v", err)
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Save the Excel file
	fmt.Printf("Saving Excel file to %s...\n", outputPath)
	if err := f.SaveAs(outputPath); err != nil {
		log.Fatalf("failed to save Excel file: %v", err)
	}

	fmt.Printf("Completed! Successfully copied %d rows to %s in %v\n", rowCount, outputPath, time.Since(startTime))
}
