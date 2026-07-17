package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	b "github.com/AnkitDTUSE/excelOptimisation/busy_excel_writer/busyexcelwriter"
)

type InputPayload struct {
	Data             []map[string]interface{} `json:"data"`
	TotalData        []map[string]interface{} `json:"totalData"`
	Headers          []b.XlsHeader            `json:"headers"`
	TotalHeaders     []b.XlsHeader            `json:"totalHeaders"`
	SellerID         int                      `json:"sellerID"`
	TotalPlaceHolder map[string]string        `json:"totalPlaceHolder"`
	FPath            string                   `json:"fPath"`
	ExcelTopHeader   b.ExcelTopHeader         `json:"excelTopHeader"`
}

func BusyWriter() {
	bytes, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatalf("failed to read input.json: %v", err)
	}

	var payload InputPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		log.Fatalf("failed to unmarshal input.json: %v", err)
	}

	file, err := b.GenerateExcelfileWithExtraHeadingHeader(
		payload.Data,
		payload.TotalData,
		payload.Headers,
		payload.TotalHeaders,
		payload.SellerID,
		payload.TotalPlaceHolder,
		payload.FPath,
		payload.ExcelTopHeader,
	)
	if err != nil {
		log.Fatalf("failed to generate excel file: %v", err)
	}

	fmt.Println("wrote:", file)
}