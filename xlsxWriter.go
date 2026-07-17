package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

func XlsxWriter() {

	startTime := time.Now()

	srcPath := "MOCK_DATA.xlsx"
	srcFile, err := excelize.OpenFile(srcPath)
	if err != nil {
		log.Fatalf("failed to open source workbook %q: %v", srcPath, err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			log.Printf("failed to close source workbook: %v", err)
		}
	}()

	destPath := "test.xlsx"
	var destFile *excelize.File

	if _, err := os.Stat(destPath); err == nil {
		destFile, err = excelize.OpenFile(destPath)
		if err != nil {
			log.Printf("destination workbook %q is not a valid Excel file; recreating it: %v", destPath, err)
			if removeErr := os.Remove(destPath); removeErr != nil {
				log.Fatalf("failed to remove invalid destination workbook %q: %v", destPath, removeErr)
			}
			destFile = excelize.NewFile()
		}
	} else if os.IsNotExist(err) {
		destFile = excelize.NewFile()
	} else {
		log.Fatalf("failed to inspect destination workbook %q: %v", destPath, err)
	}

	defer func() {
		if err := destFile.Close(); err != nil {
			log.Printf("failed to close destination workbook: %v", err)
		}
	}()

	sheetName := "data"
	if len(srcFile.GetSheetList()) > 0 {
		sheetName = srcFile.GetSheetList()[0]
	}

	srcRows, err := srcFile.GetRows(sheetName)
	if err != nil {
		log.Fatalf("failed to read rows from source workbook %q sheet %q: %v", srcPath, sheetName, err)
	}

	if index, err := destFile.GetSheetIndex(sheetName); err != nil || index < 0 {
		if _, err := destFile.NewSheet(sheetName); err != nil {
			log.Fatalf("failed to create destination sheet %q: %v", sheetName, err)
		}
	}

	destRows, err := destFile.GetRows(sheetName)
	if err != nil {
		log.Fatalf("failed to read rows from destination workbook %q sheet %q: %v", destPath, sheetName, err)
	}

	startRow := len(destRows) + 1

	for i, row := range srcRows {
		currentRow := startRow + i

		for colInd, cellValue := range row {
			colName, err := excelize.ColumnNumberToName(colInd + 1)
			if err != nil {
				log.Fatalf("failed to convert column index %d to name: %v", colInd+1, err)
			}

			cellCoordinates := fmt.Sprintf("%s%d", colName, currentRow)

			if err := destFile.SetCellValue(sheetName, cellCoordinates, cellValue); err != nil {
				log.Fatal(err)
			}
		}
	}

	copyDuration := time.Since(startTime)
	fmt.Printf("Copied %d rows in %v\n", len(srcRows), copyDuration)

	if err := destFile.SaveAs(destPath); err != nil {
		log.Fatal(err)
	}

	saveDuration := time.Since(startTime) - copyDuration
	fmt.Printf("Saved workbook to %s in %v\n", destPath, saveDuration)
	fmt.Printf("Total time: %v\n", time.Since(startTime))
}
