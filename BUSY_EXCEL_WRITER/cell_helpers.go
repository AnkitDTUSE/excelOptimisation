package busyexcelwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

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
