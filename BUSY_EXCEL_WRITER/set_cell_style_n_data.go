package busyexcelwriter

import (
	"math"
	"reflect"

	"github.com/xuri/excelize/v2"
)

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
