package busyexcelwriter
package busyexcelwriter

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// XlsHeader describes the label and source field for a sheet column.
type XlsHeader struct {
	Name  string
	Field string
}

// ExcelDataSet holds a header block together with rows for the top section.
type ExcelDataSet struct {
	Headers []XlsHeader
	Data    []map[string]interface{}
}

// ExcelTopHeader defines the optional top heading layout.
type ExcelTopHeader struct {
	TopHeader               []string
	MergeCells              []int
	FontSize                []int
	Datasets                []ExcelDataSet
	BoldHeaderData          bool
	NeedNoSpace             bool
	DataMoreThanDefinedLimit bool
}

const (
	rowsLimit         = 10000
	rowsLimitFiveLakh = 500000
)

// GenerateExcelfileWithExtraHeadingHeader creates an Excel workbook containing
// the main tabular data and the optional top heading blocks.
func GenerateExcelfileWithExtraHeadingHeader(data, totalData []map[string]interface{}, headers, totalHeaders []XlsHeader, sellerID int, totalPlaceHolder map[string]string, fPath string, excelTopHeader ExcelTopHeader) (string, error) {
	_ = sellerID
	_ = totalHeaders
	_ = totalData
	_ = totalPlaceHolder

	sTime := time.Now()
	f := excelize.NewFile()
	sheet := "Sheet1"
	SetXlsDataWithHeader(f, sheet, headers, totalHeaders, data, totalData, totalPlaceHolder, excelTopHeader)

	basePath := filepath.Clean(fPath)
	if basePath == "." || basePath == "" {
		basePath = "."
	}
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return "", err
	}

	file := filepath.Join(basePath, fmt.Sprintf("%d.xlsx", sTime.Unix()))
	if err := f.SaveAs(file); err != nil {
		return "", err
	}
	return file, nil
}

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

// setCellStyleNData applies styling and writes a value into the workbook.
func setCellStyleNData(f *excelize.File, sheet string, value interface{}, r int, c int, style int, colMax map[string]float64) {
	cell, _ := CoordinatesToCellName(c, r)
	col, _ := ColumnNumberToName(c)

	if value != nil {
		typeOf := reflect.TypeOf(value)
		if typeOf != nil && typeOf.Kind() == reflect.String {
			colMax[col] = math.Max(colMax[col], float64(len(value.(string))))
		}
	}
	colMax[col] = math.Max(colMax[col], 15)
	colMax[col] = math.Min(colMax[col], 25)
	_ = f.SetCellStyle(sheet, cell, cell, style)
	_ = f.SetCellValue(sheet, cell, value)
}

// CoordinatesToCellName converts a column and row index into a cell name.
func CoordinatesToCellName(col, row int) (string, error) {
	if col < 1 || row < 1 {
		return "", fmt.Errorf("invalid cell coordinates [%d, %d]", col, row)
	}
	colname, err := ColumnNumberToName(col)
	if err != nil {
		return "", fmt.Errorf("invalid cell coordinates [%d, %d]: %v", col, row, err)
	}
	return fmt.Sprintf("%s%d", colname, row), nil
}

// ColumnNumberToName converts a 1-based column number into Excel notation.
func ColumnNumberToName(num int) (string, error) {
	if num < 1 {
		return "", fmt.Errorf("incorrect column number %d", num)
	}
	var col string
	for num > 0 {
		col = string((num-1)%26+65) + col
		num = (num - 1) / 26
	}
	return col, nil
}

// GetAlphabetByNumber converts a 1-based column number into Excel notation.
func GetAlphabetByNumber(number int) string {
	result := ""
	for number > 0 {
		remainder := (number - 1) % 26
		result = string('A'+rune(remainder)) + result
		number = (number - 1) / 26
	}
	return result
}

// IfToA converts a value into a string representation.
func IfToA(value interface{}) string {
	if value == nil {
		return ""
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	}
	return fmt.Sprintf("%v", value)
}

// MakeDir ensures a directory exists for the file path.
func MakeDir(file string) (err error) {
	dir := filepath.Dir(file)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}
