package busyexcelwriter

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// SetXlsDataWithHeader writes headers, data rows, totals, and optional top sections.
func SetXlsDataWithHeader(f *excelize.File, sheet string,
	header, headersForTotalValues []XlsHeader, data, totalData []map[string]interface{}, totalPlaceHolder map[string]string, excelTopHeader ExcelTopHeader, rowStartIndex ...int) {
	var index int
	colMax := make(map[string]float64)
	rowIndex := 1
	colIndex := 1
	if len(rowStartIndex) > 0 {
		rowIndex = rowStartIndex[0]
	}

	if len(excelTopHeader.TopHeader) > 0 {
		for idx := range excelTopHeader.TopHeader {
			f.InsertRow(sheet, idx)
			f.MergeCell(sheet, "A"+strconv.Itoa(rowIndex), GetAlphabetByNumber(excelTopHeader.MergeCells[idx])+strconv.Itoa(rowIndex))
			f.SetCellValue(sheet, "A"+strconv.Itoa(rowIndex), excelTopHeader.TopHeader[idx])
			headerStyle, _ := f.NewStyle(fmt.Sprintf(`{
				"font":{"bold":true,"size": %d ,"name":"Arial"},
				"alignment":{"horizontal":"centre","vertical":"center"}
				}`, excelTopHeader.FontSize[idx]))
			f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIndex), GetAlphabetByNumber(excelTopHeader.MergeCells[idx])+strconv.Itoa(rowIndex), headerStyle)
			rowIndex++
		}
	}

	hStyle, _ := f.NewStyle(`
	{
		"font":{"bold":true,"size":12,"color":"#000000"},
		"fill":{"type":"gradient","color":["#c7c5c5","#adacac"],"shading":1},
		"alignment":{"horizontal":"left","vertical":"center","wrap_text":true}
	}`)

	cStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left","vertical":"center","wrap_text":true}, "number_format": 49}`)

	boldStyle, _ := f.NewStyle(`
	{
		"font":{"bold":true,"size":12,"color":"#000000"},
		"alignment":{"horizontal":"left","vertical":"center","wrap_text":true}
	}`)

	if len(excelTopHeader.Datasets) > 0 {
		for _, dataset := range excelTopHeader.Datasets {
			for r, cH := range dataset.Headers {
				if excelTopHeader.BoldHeaderData {
					setCellStyleNData(f, sheet, cH.Name, r+1, colIndex, boldStyle, colMax)
				} else {
					setCellStyleNData(f, sheet, cH.Name, r+1, colIndex, hStyle, colMax)
				}
				if (r + 1) > rowIndex {
					rowIndex = r + 1
				}
			}

			colIndex++
			for _, cData := range dataset.Data {
				for r, cH := range dataset.Headers {
					if value, exists := cData[cH.Field]; exists {
						if excelTopHeader.BoldHeaderData {
							setCellStyleNData(f, sheet, value, r+1, colIndex, boldStyle, colMax)
						} else {
							setCellStyleNData(f, sheet, value, r+1, colIndex, cStyle, colMax)
						}
						if (r + 1) > rowIndex {
							rowIndex = r + 1
						}
					}
				}
				colIndex++
			}
			if !excelTopHeader.NeedNoSpace {
				colIndex += 1
			}
		}
		if excelTopHeader.NeedNoSpace {
			rowIndex += 1
		} else {
			rowIndex += 2
		}
	}

	cTotalStyle, _ := f.NewStyle(`
		{
			"font":{"bold":true,"size":12,"color":"#000000"},
			"alignment":{"horizontal":"left","vertical":"center","wrap_text":true}, "number_format": 49
		}`)

	if len(rowStartIndex) > 0 {
		rowIndex = rowStartIndex[0]
	}

	for c, cH := range header {
		setCellStyleNData(f, sheet, cH.Name, rowIndex, c+1, hStyle, colMax)
	}
	row := 0
	dataPerSheet := rowsLimit
	if excelTopHeader.DataMoreThanDefinedLimit {
		dataPerSheet = rowsLimitFiveLakh
	}

	for _, cData := range data {
		if (row)%dataPerSheet == 0 {
			sheet = "Sheet" + strconv.Itoa(index+1)
			index = f.NewSheet(sheet)
			for c, cH := range header {
				setCellStyleNData(f, sheet, cH.Name, rowIndex, c+1, hStyle, colMax)
			}
			row = 0
		}
		for c, cH := range header {
			if value, exists := cData[cH.Field]; exists {
				setCellStyleNData(f, sheet, value, row+rowIndex+1, c+1, cStyle, colMax)
			}
		}
		row++
	}

	var space int
	if excelTopHeader.NeedNoSpace {
		space = 1
	} else {
		space = 2
	}
	if len(totalData) > 0 && row > 0 {
		for c, cH := range headersForTotalValues {
			if value, exists := totalData[0][cH.Field]; exists {
				setCellStyleNData(f, sheet, value, row+rowIndex+space, c+1, cTotalStyle, colMax)
			} else {
				if cH.Name == totalPlaceHolder["ReferringColumn"] {
					setCellStyleNData(f, sheet, totalPlaceHolder["TotalColumnName"], row+rowIndex+space, c+1, cTotalStyle, colMax)
				}
			}
		}
	}
	f.SetActiveSheet(index)

	for col, size := range colMax {
		f.SetColWidth(sheet, col, col, size)
	}
}
