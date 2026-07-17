```go
func main() {
	// (1) Columns: Name = shown in the sheet, Field = key read from each row.
	headers := []XlsHeader{
		{Name: "Parent Subscription ID", Field: "parent_subscription_id"},
		{Name: "Recom Subscription No.", Field: "recom_subscription_no"},
		{Name: "Purchase Date", Field: "purchase_date"},
		{Name: "Last Transaction Date", Field: "last_transaction_date"},
		{Name: "Recom Expiry Date", Field: "recom_expiry_date"},
		{Name: "Owner Partner", Field: "owner_partner"},
		{Name: "Order Pack", Field: "order_pack"},
		{Name: "Channel", Field: "channel"},
		{Name: "Total Orders", Field: "total_orders"},
		{Name: "Orders Consumed", Field: "orders_consumed"},
		{Name: "Last Used Date", Field: "last_used_date"},
	}

	// (2) The rows (keys MUST equal the Field values above — see Section 5).
	data := []map[string]interface{}{
		{
			"parent_subscription_id": "BD-2023-000817", "recom_subscription_no": "RC-2024-004531",
			"purchase_date": "15-Jan-2025", "last_transaction_date": "10-Jul-2026",
			"recom_expiry_date": "14-Jan-2027", "owner_partner": "Acme Distributors Pvt Ltd",
			"order_pack": 5000, "channel": "Online", "total_orders": 5000,
			"orders_consumed": 1234, "last_used_date": "10-Jul-2026",
		},
		{
			"parent_subscription_id": "BO-2022-114290", "recom_subscription_no": "RC-2023-009980",
			"purchase_date": "02-Mar-2024", "last_transaction_date": "28-Jun-2026",
			"recom_expiry_date": "01-Mar-2026", "owner_partner": "Nova Retail Solutions",
			"order_pack": 10000, "channel": "Offline", "total_orders": 10000,
			"orders_consumed": 9876, "last_used_date": "28-Jun-2026",
		},
		{
			"parent_subscription_id": "BD-2021-556301", "recom_subscription_no": "RC-2025-001204",
			"purchase_date": "20-Nov-2025", "last_transaction_date": "",
			"recom_expiry_date": "19-Nov-2026", "owner_partner": "",
			"order_pack": 0, "channel": "", "total_orders": 2500,
			"orders_consumed": 0, "last_used_date": "",
		},
	}

	// (3) Totals row: values keyed by Field. Columns with no value stay blank,
	//     except the one whose Name == ReferringColumn, which shows TotalColumnName.
	totalHeaders := headers
	totalData := []map[string]interface{}{
		{"total_orders": 17500, "orders_consumed": 11110},
	}
	totalPlaceHolder := map[string]string{
		"ReferringColumn": "Owner Partner", // print the label under "Owner Partner"
		"TotalColumnName": "TOTAL",
	}

	// (4) The banner on top (this is the "ExtraHeading" feature).
	topHeader := ExcelTopHeader{
		TopHeader:  []string{"Recom Order Consumption Report"},
		MergeCells: []int{11}, // merge across all 11 columns (A..K)
		FontSize:   []int{16},
	}

	// (5) Generate.
	file, err := GenerateExcelfileWithExtraHeadingHeader(
		data, totalData, headers, totalHeaders,
		0,                // sellerID (unused)
		totalPlaceHolder, // totals label config
		"recom_report",   // output file name (becomes recom_report.xlsx)
		topHeader,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote:", file)
}
```