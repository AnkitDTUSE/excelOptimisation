package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

func StreamWriter() {
	startTime := time.Now()

	srcPath := "MOCK_DATA.xlsx"
	srcFile, err := excelize.OpenFile(srcPath)

	if err != nil {
		log.Fatalf("failed to open source workbook %q: %v", srcPath, err)
	}
	defer srcFile.Close()

	destPath := "test.xlsx"
	sheetName := "data"
	if len(srcFile.GetSheetList()) > 0 {
		sheetName = srcFile.GetSheetList()[0]
	}

	outFile := excelize.NewFile()
	defer outFile.Close()

	if sheetName != "Sheet1" {
		outFile.NewSheet(sheetName)
	}

	sw, err := outFile.NewStreamWriter(sheetName) // make a newStreamWriter

	if err != nil {
		log.Fatalf("failed to create stream writer: %v", err)
	}

	currentRow := 1
	actualCopiedRows := 0

	if _, err := os.Stat(destPath); err == nil {

		destFile, err := excelize.OpenFile(destPath)
		if err == nil {
			destRowsIterator, err := destFile.Rows(sheetName)

			if err == nil {
				for destRowsIterator.Next() {

					row, _ := destRowsIterator.Columns()
					rowVals := make([]interface{}, len(row))
					for i, v := range row {
						rowVals[i] = v
					}
					cellAxis, _ := excelize.CoordinatesToCellName(1, currentRow)
					if err := sw.SetRow(cellAxis, rowVals); err != nil {
						log.Fatalf("failed to stream target row: %v", err)
					}
					currentRow++

				}
				destRowsIterator.Close()
			}
			destFile.Close()
		}
	}

	srcRowsIterator, err := srcFile.Rows(sheetName)
	if err != nil {
		log.Fatalf("failed to create source rows iterator: %v", err)
	}
	defer srcRowsIterator.Close()

	for srcRowsIterator.Next() {
		row, err := srcRowsIterator.Columns()
		if err != nil {
			log.Fatalf("failed to read source row stream: %v", err)
		}

		// Map string elements into interface slots required by SetRow
		rowVals := make([]interface{}, len(row))
		for i, v := range row {
			rowVals[i] = v
		}

		// Grab row coordinate string (e.g., "A15")
		cellAxis, err := excelize.CoordinatesToCellName(1, currentRow)
		if err != nil {
			log.Fatalf("failed to generate coordinate: %v", err)
		}

		// Inject the entire structural row block at once
		if err := sw.SetRow(cellAxis, rowVals); err != nil {
			log.Fatalf("failed to stream-write row: %v", err)
		}

		currentRow++
		actualCopiedRows++
	}

	// 6. CRITICAL TWEAK: Finalize open tags inside the XML stream structure
	if err := sw.Flush(); err != nil {
		log.Fatalf("failed to flush stream writer: %v", err)
	}

	copyDuration := time.Since(startTime)
	fmt.Printf("Streamed %d new rows (Total rows in file: %d) in %v\n", actualCopiedRows, currentRow-1, copyDuration)

	// 7. Save the optimized package out to disk
	if err := outFile.SaveAs(destPath); err != nil {
		log.Fatal(err)
	}

	saveDuration := time.Since(startTime) - copyDuration
	fmt.Printf("Saved workbook to %s in %v\n", destPath, saveDuration)
	fmt.Printf("Total time: %v\n", time.Since(startTime))
}
