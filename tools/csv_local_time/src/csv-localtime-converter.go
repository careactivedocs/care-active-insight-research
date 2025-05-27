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
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Print program banner
	fmt.Println("CSV Local Time Converter. Care Active (c) 2025.")

	// Define command line flags
	inputFile := flag.String("input", "", "CSV file to be processed")
	inputDir := flag.String("dir", "", "Directory containing CSV files")
	outputFile := flag.String("output", "", "Output CSV file path")
	timezone := flag.String("tz", "", "Timezone for conversion (TZ Database Identifier)")
	flag.Parse()

	// Check if both input file and directory are provided
	if *inputFile != "" && *inputDir != "" {
		fmt.Println("Warning: Both -input and -dir provided. Using -input file only.")
	}

	// Get the location from timezone
	loc, err := time.LoadLocation(*timezone)
	if err != nil {
		fmt.Printf("Error: Invalid timezone '%s': %v\n", *timezone, err)
		os.Exit(1)
	}

	// timezone is required
	if *timezone == "" {
		fmt.Println("Error: Timezone is required. Use -tz flag.")
		os.Exit(1)
	}

	// Process based on provided options
	if *inputFile != "" {
		// Process single file
		outputPath := determineOutputPath(*inputFile, *outputFile)
		err = processCSVFile(*inputFile, outputPath, loc, *timezone)
		if err != nil {
			fmt.Printf("Error processing file: %v\n", err)
			os.Exit(1)
		}
	} else if *inputDir != "" {
		// Process all CSV files in directory
		err = processDirectory(*inputDir, *outputFile, loc, *timezone)
		if err != nil {
			fmt.Printf("Error processing directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Process from standard input
		outputPath := determineOutputPath("stdin.csv", *outputFile)
		err = processStdin(outputPath, loc, *timezone)
		if err != nil {
			fmt.Printf("Error processing standard input: %v\n", err)
			os.Exit(1)
		}
	}
}

// determineOutputPath generates the appropriate output path based on input and output options
func determineOutputPath(inputPath string, outputPath string) string {
	if outputPath != "" {
		return outputPath
	}

	// If input is from stdin and no output specified, use stdout
	if inputPath == "stdin.csv" {
		return "stdout"
	}

	dir, file := filepath.Split(inputPath)
	return filepath.Join(dir, "local_"+file)
}

// processDirectory processes all CSV files in the given directory
func processDirectory(dirPath string, outputDir string, loc *time.Location, tzName string) error {
	// Check if directory exists
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("error accessing directory: %v", err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("provided path is not a directory: %s", dirPath)
	}

	// Create output directory if specified and doesn't exist
	if outputDir != "" {
		if err := createDirectoryIfNotExists(outputDir); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Read directory contents
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory contents: %v", err)
	}

	// Process each CSV file
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			inputPath := filepath.Join(dirPath, file.Name())
			var outputPath string

			if outputDir != "" {
				// If output directory is specified, use it
				outputPath = filepath.Join(outputDir, "local_"+file.Name())
			} else {
				// Otherwise, place output file next to the input file
				outputPath = filepath.Join(dirPath, "local_"+file.Name())
			}

			err := processCSVFile(inputPath, outputPath, loc, tzName)
			if err != nil {
				fmt.Printf("Warning: Failed to process %s: %v\n", inputPath, err)
				continue
			}
		}
	}
	return nil
}

// createDirectoryIfNotExists creates a directory if it doesn't exist
func createDirectoryIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("Creating output directory: %s\n", dirPath)
		return os.MkdirAll(dirPath, 0755)
	} else if err != nil {
		return err
	}
	return nil
}

// processStdin processes CSV data from standard input
func processStdin(outputPath string, loc *time.Location, tzName string) error {
	fmt.Println("Reading CSV data from standard input...")

	var writer *csv.Writer
	var outputFile *os.File

	// Check if we should write to stdout
	if outputPath == "stdout" {
		writer = csv.NewWriter(os.Stdout)
	} else {
		// Open the output file
		var err error
		outputFile, err = os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer outputFile.Close()
		writer = csv.NewWriter(outputFile)
	}

	// Create CSV reader
	reader := csv.NewReader(os.Stdin)
	defer writer.Flush()

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	// Process the CSV data
	err = processCSVData(reader, writer, header, loc, tzName, "stdin")
	if err != nil {
		return err
	}

	// Only print success message if not writing to stdout
	if outputPath != "stdout" {
		fmt.Printf("Successfully processed file. Output written to '%s'\n", outputPath)
	}

	return nil
}

// processCSVFile processes a single CSV file with the given path
func processCSVFile(inputPath, outputPath string, loc *time.Location, tzName string) error {
	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Create CSV reader and writer
	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	// Process the CSV data
	err = processCSVData(reader, writer, header, loc, tzName, inputPath)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully processed file. Output written to '%s'\n", outputPath)
	return nil
}

// processCSVData handles the actual CSV processing logic
func processCSVData(reader *csv.Reader, writer *csv.Writer, header []string, loc *time.Location, tzName, sourceName string) error {
	// Find timestamp column index with updated priority order
	timeColumnIndex, timeColumnName, err := findTimeColumnIndex(header)
	if err != nil {
		return err
	}

	// Find sample_at_utc column index to match format
	utcFormatIndex := findColumnIndex(header, "sample_at_utc")

	// Write the new header (local_time + timezone + original headers)
	newHeader := append([]string{"local_time", "local_timezone"}, header...)
	if err := writer.Write(newHeader); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	// Read the first data row to determine UTC format if available
	var firstRecord []string
	var utcFormat string

	firstRecord, err = reader.Read()
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading first record: %v", err)
	}

	// Determine UTC format from sample_at_utc if available
	if err != io.EOF && utcFormatIndex >= 0 && utcFormatIndex < len(firstRecord) {
		utcFormat = determineTimeFormat(firstRecord[utcFormatIndex])
	}

	// Process the first row if we read it
	if err != io.EOF {
		// Convert timestamp to local time
		localTime, err := convertTimestampToLocalTime(firstRecord[timeColumnIndex], loc, utcFormat)
		if err != nil {
			fmt.Printf("Warning: Failed to convert timestamp for row: %v\n", err)
			localTime = "INVALID_TIMESTAMP"
		}

		// Create new row with local_time and timezone added
		newRow := append([]string{localTime, tzName}, firstRecord...)
		if err := writer.Write(newRow); err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
	}

	// Process remaining rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading record: %v", err)
		}

		// Convert timestamp to local time
		localTime, err := convertTimestampToLocalTime(record[timeColumnIndex], loc, utcFormat)
		if err != nil {
			fmt.Printf("Warning: Failed to convert timestamp for row: %v\n", err)
			localTime = "INVALID_TIMESTAMP"
		}

		// Create new row with local_time and timezone added
		newRow := append([]string{localTime, tzName}, record...)
		if err := writer.Write(newRow); err != nil {
			return fmt.Errorf("error writing record: %v", err)
		}
	}

	fmt.Printf("Used time column: %s (index: %d)\n", timeColumnName, timeColumnIndex)
	fmt.Printf("Local timezone: %s\n", tzName)
	return nil
}

// findTimeColumnIndex locates the timestamp column in the CSV header with updated priority
func findTimeColumnIndex(header []string) (int, string, error) {
	// Updated priority order for timestamp columns
	timestampColumns := []string{"aggregate_at", "resample_at", "resample_at_ms", "sample_at_ms", "scanned_at_ms", "created_at_ms"}

	for _, colName := range timestampColumns {
		for i, h := range header {
			if h == colName {
				return i, colName, nil
			}
		}
	}

	return -1, "", fmt.Errorf("no timestamp column found (looking for: %v)", timestampColumns)
}

// findColumnIndex finds the index of a specific column name
func findColumnIndex(header []string, columnName string) int {
	for i, h := range header {
		if h == columnName {
			return i
		}
	}
	return -1
}

// determineTimeFormat attempts to identify the format from a UTC time string
func determineTimeFormat(utcTimeStr string) string {
	// Default format if we can't determine
	defaultFormat := "2006-01-02T15:04:05.000Z"

	// Common UTC time formats
	formats := []string{
		"2006-01-02T15:04:05.000Z",       // ISO8601 with milliseconds
		"2006-01-02T15:04:05Z",           // ISO8601 without milliseconds
		"2006-01-02 15:04:05",            // Simple datetime
		"2006-01-02 15:04:05.000",        // Simple datetime with milliseconds
		"2006-01-02T15:04:05.000",        // ISO8601 without Z
		"2006-01-02T15:04:05",            // ISO8601 without Z or milliseconds
		"Mon Jan 02 15:04:05 -0700 2006", // RFC822 with timezone and year
	}

	for _, format := range formats {
		_, err := time.Parse(format, utcTimeStr)
		if err == nil {
			return format
		}
	}

	return defaultFormat
}

// convertTimestampToLocalTime converts a timestamp to a human-readable local time
// now matching the format of sample_at_utc if available
func convertTimestampToLocalTime(timestampStr string, loc *time.Location, utcFormat string) (string, error) {
	timestampStr = strings.TrimSpace(timestampStr)
	var utcTime time.Time

	// Handle resample_at_utc which uses seconds instead of milliseconds
	if strings.Contains(timestampStr, "T") || strings.Contains(timestampStr, " ") {
		// This appears to be a formatted UTC time string, not a timestamp
		for _, format := range []string{
			"2006-01-02T15:04:05.000Z",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02 15:04:05.000",
		} {
			parsedTime, err := time.Parse(format, timestampStr)
			if err == nil {
				utcTime = parsedTime
				break
			}
		}

		// If we couldn't parse the time string
		if utcTime.IsZero() {
			return "", fmt.Errorf("invalid time format: %s", timestampStr)
		}
	} else {
		// Parse the timestamp
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid timestamp format: %v", err)
		}

		// Convert seconds to milliseconds for resample_at (if seconds)
		// We assume a timestamp before year 2001 is in seconds, not milliseconds
		if timestamp < 1000000000000 {
			timestamp *= 1000 // Convert seconds to milliseconds
		}

		// Convert milliseconds to seconds and nanoseconds
		seconds := timestamp / 1000
		nanoseconds := (timestamp % 1000) * 1000000

		// Create time object in UTC
		utcTime = time.Unix(seconds, nanoseconds).UTC()
	}

	// Convert to local time
	localTime := utcTime.In(loc)

	// Format based on the UTC format if available, otherwise use ISO8601
	var formattedTime string
	if utcFormat != "" {
		// Try to use the same format as sample_at_utc
		formattedTime = localTime.Format(utcFormat)
	} else {
		// Default format as ISO8601 without the Z
		formattedTime = localTime.Format("2006/01/02 15:04:05.000")
	}

	// Remove the 'Z' suffix if present
	formattedTime = strings.TrimSuffix(formattedTime, "Z")

	return formattedTime, nil
}
