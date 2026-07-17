package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

func CsvWriter() {

	startTime := time.Now()

	csvFile, err := os.Open("mock_data.csv")
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	csvRows, err := reader.ReadAll() // Returns a [][]string
	if err != nil {
		log.Fatalf("Error reading CSV data: %v", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	for rowIndex, row := range csvRows {
		excelRow := rowIndex + 1

		for colIndex, cellValue := range row {
			colName, err := excelize.ColumnNumberToName(colIndex + 1)
			if err != nil {
				log.Fatalf("Error converting column number: %v", err)
			}

			cellCoordinate := fmt.Sprintf("%s%d", colName, excelRow)

			err = f.SetCellValue("Sheet1", cellCoordinate, cellValue)
			if err != nil {
				log.Fatalf("Error setting cell value at %s: %v", cellCoordinate, err)
			}
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Copied %d rows in %v\n", len(csvRows), elapsed)

	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Fatalf("Error saving Excel file: %v", err)
	}
	saveDuration := time.Since(startTime) - elapsed

	fmt.Printf("Saved workbook in %v\n", saveDuration)
	fmt.Printf("Total time: %v\n", time.Since(startTime))
}
