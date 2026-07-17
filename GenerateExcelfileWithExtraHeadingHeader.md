# `GenerateExcelfileWithExtraHeadingHeader` — functions

The main function and every function it uses.

---

### `GenerateExcelfileWithExtraHeadingHeader`

```go
func GenerateExcelfileWithExtraHeadingHeader(data, totalData []map[string]interface{}, headers, totalHeaders []XlsHeader, sellerID int, totalPlaceHolder map[string]string, fPath string, excelTopHeader ExcelTopHeader) (
	file string, err error) {
	sTime := time.Now()
	f := excelize.NewFile()
	sheet := "Sheet1"
	SetXlsDataWithHeader(f, sheet, headers, totalHeaders, data, totalData, totalPlaceHolder, excelTopHeader)
	fileBasePath := conf.String("path.report_export_path", "/seller/")
	filePath := fileBasePath + fPath
	file = fmt.Sprintf("%v%v.xlsx", filePath, sTime.Unix())
	if err = ally.MakeDir(file); err == nil {
		if err = f.SaveAs(file); err != nil {
			err = perror.MiscError(err)
		}
	}
	return
}
```

---

### `SetXlsDataWithHeader`

```go
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
			sheet = "Sheet" + ally.IfToA(index+1)
			index = f.NewSheet(sheet)
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
```

---

### `setCellStyleNData`

```go
func setCellStyleNData(f *excelize.File, sheet string, value interface{}, r int, c int, style int, colMax map[string]float64) {
	cell, _ := CoordinatesToCellName(c, r)
	col, _ := ColumnNumberToName(c)

	typeOf := reflect.TypeOf(value)
	if typeOf.Kind() == reflect.String {
		colMax[col] = math.Max(colMax[col], float64(len(value.(string))))
	}
	colMax[col] = math.Max(colMax[col], 15)
	colMax[col] = math.Min(colMax[col], 25)
	f.SetCellStyle(sheet, cell, cell, style)
	f.SetCellValue(sheet, cell, value)
}
```

---

### `CoordinatesToCellName` 

```go
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
```

---

### `ColumnNumberToName` 

```go
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
```

---

### `GetAlphabetByNumber` 
```go
func GetAlphabetByNumber(number int) string {
	result := ""
	for number > 0 {
		remainder := (number - 1) % 26
		result = string('A'+rune(remainder)) + result
		number = (number - 1) / 26
	}
	return result
}
```

---

### `ally.IfToA`

```go
func IfToA(value interface{}) string {
	if value == nil {
		return ""
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	}
	return fmt.Sprintf("%v", value)
}
```

---

### `ally.MakeDir` 

```go
func MakeDir(file string) (err error) {
	dir := filepath.Dir(file)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			err = perror.MiscError(err, "Dir Create Error")
		}
	}
	return
}
```