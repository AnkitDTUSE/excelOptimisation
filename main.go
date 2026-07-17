package main

func main() {
	// XlsxWriter()
	BusyWriter()
	// StreamWriter()
	// CsvWriter()
}

// benchmarking
//
// 10K rows --> Copied 10009 rows in 2.7011978s
// 				Saved workbook to test.xlsx in 830.4224ms
// 				Total time: 3.5316202s
//
// 3lkh rows --> Copied 360001 rows in 58.3086173s
//				 Saved workbook to test.xlsx in 21.4936326s
// 				 Total time: 1m19.8022499s
//
// 40K rows --> Copied 40033 rows in 4.9196453s
// 				Saved workbook to test.xlsx in 2.5835634s
// 				Total time: 7.5032087s
//
// 40K rows CSV --> Copied 40000 rows in 2.3202505s
// 					Saved workbook in 3.7260253s
// 					Total time: 6.0462758s
//
// 3lkh rows --> Copied 360001 rows in 57.1621636s
// 				 Saved workbook to test.xlsx in 33.3544289s
// 				 Total time: 1m30.5172252s
// 
// 3lkh rows (streamWriter) --> 360001 new rows (Total rows in file: 360001) in 14.4853067s
// 				 Saved workbook to test.xlsx in 5.4668493s
// 				 Total time: 19.9530858s