package busyexcelwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/AnkitDTUSE/excelOptimisation/busy_excel_writer/perror"
)

// IfToA converts any interface value (especially float64 or nil) to its string representation.
func IfToA(value any) string {
	if value == nil {
		return ""
	} else if reflect.TypeOf(value).Kind() == reflect.Float64 {
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	}
	return fmt.Sprintf("%v", value)
}

// MakeDir creates the directories necessary to write a file to the specified path if they do not exist.
func MakeDir(file string) (err error) {
	dir := filepath.Dir(file)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			err = perror.MiscError(err, "Dir Create Error")
		}
	}
	return
}
