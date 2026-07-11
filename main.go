package main

import (
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	sheet := "test"
	index, err := f.NewSheet(sheet)

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"ID", "Name", "Department", "Salary"}
	rows := [][]any{
		{101, "Alice Smith", "Engineering", 95000},
		{102, "Bob Jones", "Marketing", 72000},
		{103, "Charlie Brown", "Design", 80000},
	}

	for colIndex, header := range headers {
		cellName, err := excelize.CoordinatesToCellName(colIndex+1, 1)
		if err != nil {
			log.Fatal(err)
		}
		f.SetCellValue(sheet, cellName, header)
	}

	for rowIndex, rowData := range rows {
		for colIndex, value := range rowData {
			cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				log.Fatal(err)
			}
			f.SetCellValue(sheet, cellName, value)
		}
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs("Employees.xlsx"); err != nil {
		log.Fatal(err)
	}

	copyTxtToXlsx()
}
