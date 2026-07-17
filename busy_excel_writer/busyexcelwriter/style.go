package busyexcelwriter

import (
	"math"
	"reflect"

	"github.com/xuri/excelize/v2"
)

// setCellStyleNData sets formatting style and value for a single cell and updates the colMax widths.
func setCellStyleNData(f *excelize.File, sheet string, value interface{}, r int, c int, style int, colMax map[string]float64) {
	cell, _ := CoordinatesToCellName(c, r)
	col, _ := ColumnNumberToName(c)

	typeOf := reflect.TypeOf(value)
	if typeOf != nil && typeOf.Kind() == reflect.String {
		colMax[col] = math.Max(colMax[col], float64(len(value.(string))))
	}
	colMax[col] = math.Max(colMax[col], 15)
	colMax[col] = math.Min(colMax[col], 25)
	f.SetCellStyle(sheet, cell, cell, style)
	f.SetCellValue(sheet, cell, value)
}
