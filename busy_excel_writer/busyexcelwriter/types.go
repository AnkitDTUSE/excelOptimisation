package busyexcelwriter

// XlsHeader represents a single column header in the Excel file.
type XlsHeader struct {
	Name  string
	Field string
}

// Dataset represents a dataset to be written in the top header section.
type Dataset struct {
	Headers []XlsHeader
	Data    []map[string]any
}

// ExcelTopHeader defines metadata and additional layout for the top of the Excel sheet.
type ExcelTopHeader struct {
	TopHeader                []string
	MergeCells               []int
	FontSize                 []int
	BoldHeaderData           bool
	Datasets                 []Dataset
	NeedNoSpace              bool
	DataMoreThanDefinedLimit bool
}
