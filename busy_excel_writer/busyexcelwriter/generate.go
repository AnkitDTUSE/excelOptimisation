package busyexcelwriter

import (
	"fmt"
	"time"

	"github.com/AnkitDTUSE/excelOptimisation/busy_excel_writer/perror"
	"github.com/xuri/excelize/v2"
)

// GenerateExcelfileWithExtraHeadingHeader creates an Excel file with top headers, main headers,
// data rows, total rows, and returns the file path or an error.
func GenerateExcelfileWithExtraHeadingHeader(data, totalData []map[string]any, headers, totalHeaders []XlsHeader, sellerID int, totalPlaceHolder map[string]string, fPath string, excelTopHeader ExcelTopHeader) (
	file string, err error) {
	sTime := time.Now()
	f := excelize.NewFile()
	sheet := "Sheet1"

	SetXlsDataWithHeader(f, sheet, headers, totalHeaders, data, totalData, totalPlaceHolder, excelTopHeader)
	
	file = fmt.Sprintf("./%v%v.xlsx", fPath, sTime.Unix())

	elapsed := time.Since(sTime)
	fmt.Printf("%v time to copy all the rows\n",elapsed)
	
	saveT1 := time.Now()
	if err = MakeDir(file); err == nil {
		if err = f.SaveAs(file); err != nil {
			err = perror.MiscError(err)
		}
	}

	saveT2 := time.Since(saveT1)

	fmt.Printf("save time %v\n",saveT2)
	return file,nil
	
}
