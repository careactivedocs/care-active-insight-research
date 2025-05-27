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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Command line flags
var (
	inputFile  = flag.String("input", "", "Input CSV file path (required)")
	outputFile = flag.String("output", "", "Output CSV file path (optional, defaults to 'resampled_' + input filename)")
	gapLimit   = flag.Int64("gaplimit", 30000, "Gap limit in milliseconds")
	rssiLimit  = flag.Int64("rssilimit", 30000, "RSSI limit in milliseconds")
	verbose    = flag.Bool("verbose", false, "Enable verbose progress reporting (default: false)")
)

func main() {
	// Print copyright message
	fmt.Println("RTLS CVL Resampling Tool. Care Active Corp (c) 2025.")

	// Parse command line flags
	flag.Parse()

	// Validate required flags
	if *inputFile == "" {
		flag.Usage()
		log.Fatal("Error: -input flag is required")
	}

	// Generate output file path if not provided
	if *outputFile == "" {
		dir, filename := filepath.Split(*inputFile)
		*outputFile = filepath.Join(dir, "resampled_"+filename)
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(*outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	startTime := time.Now()

	// Read the input CSV
	fmt.Println("Reading input CSV file...")
	records, headers, err := readCSV(*inputFile)
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}
	fmt.Printf("Read %d records from input file\n", len(records))

	// Validate that no sample_at_ms is zero
	fmt.Println("Validating sample_at_ms values...")
	if err := validateSampleTimes(records, headers); err != nil {
		log.Fatalf("Validation error: %v", err)
		os.Exit(1)
	}

	// Extract RSSI column names
	rssiColumns := extractRSSIColumns(headers)
	fmt.Printf("Found %d RSSI columns\n", len(rssiColumns))

	// Process the data
	fmt.Println("Creating data processor...")
	processor, err := NewDataProcessor(rssiColumns, *gapLimit, *rssiLimit, *verbose)
	if err != nil {
		log.Fatalf("Error creating data processor: %v", err)
	}

	fmt.Println("Processing data...")
	// Check for backward sample times before processing
	checkBackwardSampleTimes(records, headers, *verbose)
	resampledData, err := processor.ProcessData(records, headers)
	if err != nil {
		log.Fatalf("Error processing data: %v", err)
	}
	fmt.Printf("Generated %d resampled records\n", len(resampledData))

	// Write the output CSV
	fmt.Println("Writing output CSV file...")
	outputHeaders := buildOutputHeaders(headers, rssiColumns)
	err = writeCSV(*outputFile, outputHeaders, resampledData)
	if err != nil {
		log.Fatalf("Error writing output CSV: %v", err)
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Successfully processed data and wrote to %s (took %v)\n", *outputFile, elapsedTime)
}

// validateSampleTimes checks if any sample_at_ms value is 0 and returns an error if found
func validateSampleTimes(records [][]string, headers []string) error {
	// Find the index of sample_at_ms column
	sampleAtMsIdx := -1
	for i, header := range headers {
		if header == "sample_at_ms" {
			sampleAtMsIdx = i
			break
		}
	}

	if sampleAtMsIdx == -1 {
		return fmt.Errorf("sample_at_ms column not found in the CSV headers")
	}

	// Collect all records with zero sample_at_ms
	var zeroRows []int
	for i, record := range records {
		if len(record) <= sampleAtMsIdx {
			continue // Skip malformed records
		}

		if record[sampleAtMsIdx] == "0" {
			zeroRows = append(zeroRows, i+1) // Add 1 to get 1-based row number
		}
	}

	// If any zero values found, report all of them and return error
	if len(zeroRows) > 0 {
		fmt.Printf("ERROR: Found %d row(s) with sample_at_ms value of 0\n", len(zeroRows))

		// Only print detailed row information in verbose mode
		if *verbose {
			for _, row := range zeroRows {
				fmt.Printf("  - Row %d\n", row)
			}
		}

		return fmt.Errorf("%d row(s) have sample_at_ms value of 0", len(zeroRows))
	}

	return nil
}

// checkBackwardSampleTimes checks if sample_at_ms times ever go backward
// and prints a warning message if they do, without terminating the program
func checkBackwardSampleTimes(records [][]string, headers []string, verbose bool) {
	// Find the index of sample_at_ms column
	sampleAtMsIdx := -1
	for i, header := range headers {
		if header == "sample_at_ms" {
			sampleAtMsIdx = i
			break
		}
	}

	if sampleAtMsIdx == -1 {
		fmt.Println("WARNING: sample_at_ms column not found, cannot check for backward times")
		return
	}

	// Track backward time rows
	var backwardRows []int
	var lastTime int64 = -1

	for i, record := range records {
		if len(record) <= sampleAtMsIdx {
			continue // Skip malformed records
		}

		// Parse the sample_at_ms value
		var currentTime int64
		_, err := fmt.Sscanf(record[sampleAtMsIdx], "%d", &currentTime)
		if err != nil {
			continue // Skip unparseable values
		}

		// Check if this time is less than the previous time
		if lastTime > 0 && currentTime < lastTime {
			backwardRows = append(backwardRows, i+1) // Add 1 to get 1-based row number
		}

		lastTime = currentTime
	}

	// If any backward times found, report them
	if len(backwardRows) > 0 {
		fmt.Printf("WARNING: Found %d row(s) with backward sample_at_ms times\n", len(backwardRows))

		// Only print detailed row information in verbose mode
		if verbose {
			fmt.Println("Rows with backward sample_at_ms times:")
			for _, row := range backwardRows {
				fmt.Printf("  - Row %d\n", row)
			}
		}
	}
}
