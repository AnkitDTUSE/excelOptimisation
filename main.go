package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	srcFile, _ := excelize.OpenFile("MOCK_TEST_10K.xlsx")

	defer srcFile.Close()

	destFile, _ := excelize.OpenFile("test.xlsx")

	defer destFile.Close()

	srcRows, _ := srcFile.GetRows("Sheet1")

	destRows, _ := destFile.GetRows("Sheet1")

	startRow := len(destRows) + 1

	for i, row := range srcRows {
		currentRow := startRow + i

		for colInd, cellValue := range row {

			colName, _ := excelize.ColumnNumberToName(colInd + 1)

			cellCoordinates := fmt.Sprintf("%s%d", colName, currentRow)

			if err := destFile.SetCellValue("Sheet1", cellCoordinates, cellValue); err != nil {
				log.Fatal(err)
			}

		}

	}

	if err:= destFile.Save();err!=nil{
		log.Fatal(err)
	}

}
