package busyexcelwriter

import (
	"strconv"

	"github.com/AnkitDTUSE/excelOptimisation/busy_excel_writer/constant"
	"github.com/xuri/excelize/v2"
)

// SetXlsDataWithHeader populates an excelize File sheet with top headers, main data headers,
// row data, and total summary rows, dynamically managing cell sizing and multi-sheet limits.
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
			_ = f.InsertRows(sheet, rowIndex, 1)
			f.MergeCell(sheet, "A"+strconv.Itoa(rowIndex), GetAlphabetByNumber(excelTopHeader.MergeCells[idx])+strconv.Itoa(rowIndex))
			f.SetCellValue(sheet, "A"+strconv.Itoa(rowIndex), excelTopHeader.TopHeader[idx])
			headerStyle, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Bold:   true,
					Size:   float64(excelTopHeader.FontSize[idx]),
					Family: "Arial",
				},
				Alignment: &excelize.Alignment{
					Horizontal: "center",
					Vertical:   "center",
				},
			})
			f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIndex), GetAlphabetByNumber(excelTopHeader.MergeCells[idx])+strconv.Itoa(rowIndex), headerStyle)

			rowIndex++
		}
	}

	hStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#000000",
		},
		Fill: excelize.Fill{
			Type:    "gradient",
			Color:   []string{"#c7c5c5", "#adacac"},
			Shading: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	cStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		NumFmt: 49,
	})

	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
	})

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

	cTotalStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		NumFmt: 49,
	})

	if len(rowStartIndex) > 0 {
		rowIndex = rowStartIndex[0]
	}

	// //setting header
	for c, cH := range header {
		setCellStyleNData(f, sheet, cH.Name, rowIndex, c+1, hStyle, colMax)
	}
	row := 0
	//setting data

	dataPerSheet := constant.ROWS_LIMIT

	if excelTopHeader.DataMoreThanDefinedLimit {
		dataPerSheet = constant.ROWS_LIMIT_FIVE_LAKH
	}

	for _, cData := range data {
		if (row)%dataPerSheet == 0 {
			sheet = "Sheet" + IfToA(index+1)
			index, _ = f.NewSheet(sheet)
			//setting header
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

	var Space int
	if excelTopHeader.NeedNoSpace {
		Space = 1
	} else {
		Space = 2
	}
	if len(totalData) > 0 && row > 0 {
		for c, cH := range headersForTotalValues {
			if value, exists := totalData[0][cH.Field]; exists {
				setCellStyleNData(f, sheet, value, row+rowIndex+Space, c+1, cTotalStyle, colMax)
			} else {
				if cH.Name == totalPlaceHolder["ReferringColumn"] {
					setCellStyleNData(f, sheet, totalPlaceHolder["TotalColumnName"], row+rowIndex+Space, c+1, cTotalStyle, colMax)
				}
			}
		}
	}
	f.SetActiveSheet(index)

	for col, size := range colMax {
		f.SetColWidth(sheet, col, col, size)
	}
}
