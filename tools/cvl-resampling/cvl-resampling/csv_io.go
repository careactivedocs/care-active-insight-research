/*
 * Copyright (c) 2025
 * Care Active Corp. ("CA").
 * All rights reserved.
 *
 * The information contained herein is confidential and proprietary to
 * CA. Use of this information by anyone other than authorized employees
 * of CA is granted only under a written non-disclosure agreement,
 * expressly prescribing the scope and manner of such use.
 */

//
// This program was AI-assisted code generation.
//

package main

import (
	"encoding/csv"
	"os"
	"strings"
)

// readCSV reads a CSV file and returns the records, headers, and any error.
func readCSV(filePath string) ([][]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(records) == 0 {
		return nil, nil, nil
	}

	headers := records[0]
	return records[1:], headers, nil
}

// writeCSV writes data to a CSV file and returns any error.
func writeCSV(filePath string, headers []string, data [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data
	return writer.WriteAll(data)
}

// extractRSSIColumns returns a slice of RSSI column names from the headers.
func extractRSSIColumns(headers []string) []string {
	var rssiColumns []string
	for _, header := range headers {
		if strings.Contains(header, "_rssi") {
			rssiColumns = append(rssiColumns, header)
		}
	}
	return rssiColumns
}

// buildOutputHeaders creates the headers for the output CSV.
func buildOutputHeaders(inputHeaders []string, rssiColumns []string) []string {
	// Define headers in the exact order specified
	orderedHeaders := []string{
		"family_id",
		"device_name",
		"device_mac",
		"resample_at",
		"resample_at_utc",
		"event_seq",
		"acvl",
		"resample_acvl",
		"resample_acvl_increment",
	}

	// Add the RSSI columns at the end
	outputHeaders := append(orderedHeaders, rssiColumns...)

	return outputHeaders
}
