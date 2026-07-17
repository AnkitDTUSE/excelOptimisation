package busyexcelwriter
package busyexcelwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
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
