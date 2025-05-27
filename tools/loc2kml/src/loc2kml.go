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

// Point represents a geographical point with metadata
type Point struct {
	FamilyAccount  string
	DeviceName     string
	DeviceMAC      string
	ScannedAt      time.Time
	Latitude       float64
	Longitude      float64
	GPSAccuracy    int
	ReasonCode     int
	PhoneName      string
	SenderDeviceID string
}

// getDefaultOutputName returns the default KML output filename based on input
func getDefaultOutputName(inputPath string) string {
	if inputPath == "" {
		return "output.kml"
	}

	// Keep the directory part of the path
	dir := filepath.Dir(inputPath)

	// Get just the filename without extension
	base := filepath.Base(inputPath)
	baseName := strings.TrimSuffix(base, filepath.Ext(base))

	// Combine the original directory with the new filename
	return filepath.Join(dir, baseName+".kml")
}

// formatFamilyAccount transforms email format f1c+nnnnnnnn@careactive.ai to F1Cnnnnnnnn
func formatFamilyAccount(email string) string {
	// Check if the email matches the pattern f1c+nnnnnnnn@careactive.ai
	if strings.HasPrefix(strings.ToLower(email), "f1c+") && strings.HasSuffix(strings.ToLower(email), "@careactive.ai") {
		// Extract the numeric part between f1c+ and @careactive.ai
		numericPart := strings.TrimPrefix(strings.ToLower(email), "f1c+")
		numericPart = strings.TrimSuffix(numericPart, "@careactive.ai")

		// Check if the remaining part is all numeric
		allNumeric := true
		for _, char := range numericPart {
			if char < '0' || char > '9' {
				allNumeric = false
				break
			}
		}

		if allNumeric {
			return "F1C" + numericPart
		}
	}

	// If it doesn't match the pattern, return the original email
	return email
}

func main() {
	// Print copyright message
	fmt.Println("Location CSV to KML converter. Care Active Corp (c) 2025.")

	// Define command line flags
	inputDirPtr := flag.String("dir", "", "Directory containing CSV files")
	defaultOutput := ""
	outputFilePtr := flag.String("output", defaultOutput, "Output KML file path")
	timezonePtr := flag.String("timezone", "UTC", "Timezone for conversion (TZ Database Identifier)")
	pathModePtr := flag.Bool("path", false, "Generate connected paths instead of individual placemarks")

	// Customize usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [options] [csv_files...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  Process all CSVs in a directory:\n")
		fmt.Fprintf(os.Stderr, "    %s -dir=/path/to/csv/files -output=result.kml -timezone=America/Toronto\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Process specific CSV files:\n")
		fmt.Fprintf(os.Stderr, "    %s -output=result.kml -timezone=UTC data1.csv data2.csv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Create path instead of individual points:\n")
		fmt.Fprintf(os.Stderr, "    %s -path -output=path.kml data1.csv\n", os.Args[0])
	}

	// Parse command line flags
	flag.Parse()

	// Check if no arguments provided
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	// If output file is empty, set it based on input
	if *outputFilePtr == "" {
		if *inputDirPtr != "" {
			*outputFilePtr = getDefaultOutputName(*inputDirPtr)
		} else if len(flag.Args()) > 0 {
			*outputFilePtr = getDefaultOutputName(flag.Args()[0])
		} else {
			*outputFilePtr = "output.kml"
		}
	}

	// Load timezone location
	loc, err := time.LoadLocation(*timezonePtr)
	if err != nil {
		fmt.Printf("Error: Invalid timezone '%s': %v\n", *timezonePtr, err)
		os.Exit(1)
	}

	fmt.Println("Timezone = " + *timezonePtr)

	// Initialize points collection
	var points []Point

	// Check if directory mode or specific files mode
	if *inputDirPtr != "" {
		// Directory mode - process all CSV files in the directory
		dirPoints, err := processCSVDirectory(*inputDirPtr, loc)
		if err != nil {
			fmt.Printf("Error processing CSV directory: %v\n", err)
			os.Exit(1)
		}
		points = dirPoints
	} else {
		// Specific files mode - process files from command line arguments
		remainingArgs := flag.Args()
		if len(remainingArgs) == 0 {
			fmt.Println("Error: No input files provided. Use -dir for directory mode or specify CSV files as arguments")
			flag.Usage()
			os.Exit(1)
		}

		// Process each specified file
		for _, filePath := range remainingArgs {
			if !strings.HasSuffix(strings.ToLower(filePath), ".csv") {
				fmt.Printf("Warning: Skipping non-CSV file: %s\n", filePath)
				continue
			}

			fmt.Printf("Processing file: %s\n", filePath)
			filePoints, err := processCSVFile(filePath, loc)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", filePath, err)
				continue
			}

			points = append(points, filePoints...)
		}
	}

	// Verify we have points to process
	if len(points) == 0 {
		fmt.Println("Error: No valid points found in any CSV files")
		os.Exit(1)
	}

	// Generate KML file from points
	err = generateKML(points, *outputFilePtr, *pathModePtr)
	if err != nil {
		fmt.Printf("Error generating KML file: %v\n", err)
		os.Exit(1)
	}

	if *pathModePtr {
		fmt.Printf("Successfully converted %d points to path in KML file: %s\n", len(points), *outputFilePtr)
	} else {
		fmt.Printf("Successfully converted %d points to KML file: %s\n", len(points), *outputFilePtr)
	}
}

// processCSVDirectory processes all CSV files in the specified directory
func processCSVDirectory(dirPath string, loc *time.Location) ([]Point, error) {
	var allPoints []Point

	// Get all files in the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Process each CSV file
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			filePath := filepath.Join(dirPath, file.Name())
			fmt.Printf("Processing file: %s\n", filePath)

			filePoints, err := processCSVFile(filePath, loc)
			if err != nil {
				fmt.Printf("Warning: Error processing file %s: %v, skipping...\n", file.Name(), err)
				continue
			}

			allPoints = append(allPoints, filePoints...)
		}
	}

	if len(allPoints) == 0 {
		return nil, fmt.Errorf("no valid points found in any CSV files")
	}

	return allPoints, nil
}

// processCSVFile processes a single CSV file and returns the points
func processCSVFile(filePath string, loc *time.Location) ([]Point, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Create a map of column indices
	columnMap := make(map[string]int)
	for i, column := range header {
		columnMap[column] = i
	}

	// Check for required columns
	requiredColumns := []string{
		"family_id", "device_name", "device_mac", "scanned_at_ms",
		"gps_latitude", "gps_longitude", "gps_accuracy", "reason_code",
		"phone_name", "sender_device_id",
	}

	for _, column := range requiredColumns {
		if _, ok := columnMap[column]; !ok {
			return nil, fmt.Errorf("required column '%s' not found in CSV", column)
		}
	}

	// Process each row
	var points []Point
	lineNum := 1 // Start from 1 to account for header

	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: Error reading line %d: %v, skipping...\n", lineNum, err)
			continue
		}

		// Parse timestamp
		scannedAtMs, err := strconv.ParseInt(record[columnMap["scanned_at_ms"]], 10, 64)
		if err != nil {
			fmt.Printf("Warning: Invalid timestamp at line %d: %v, skipping...\n", lineNum, err)
			continue
		}

		// Convert milliseconds to time.Time
		scannedAt := time.Unix(scannedAtMs/1000, (scannedAtMs%1000)*1000000).In(loc)

		// Parse latitude and longitude
		latitude, err := strconv.ParseFloat(record[columnMap["gps_latitude"]], 64)
		if err != nil {
			fmt.Printf("Warning: Invalid latitude at line %d: %v, skipping...\n", lineNum, err)
			continue
		}

		longitude, err := strconv.ParseFloat(record[columnMap["gps_longitude"]], 64)
		if err != nil {
			fmt.Printf("Warning: Invalid longitude at line %d: %v, skipping...\n", lineNum, err)
			continue
		}

		// Parse GPS accuracy and reason code
		gpsAccuracy, err := strconv.Atoi(record[columnMap["gps_accuracy"]])
		if err != nil {
			fmt.Printf("Warning: Invalid GPS accuracy at line %d: %v, skipping...\n", lineNum, err)
			continue
		}

		reasonCode, err := strconv.Atoi(record[columnMap["reason_code"]])
		if err != nil {
			fmt.Printf("Warning: Invalid reason code at line %d: %v, skipping...\n", lineNum, err)
			continue
		}

		// Create a new point
		point := Point{
			FamilyAccount:  formatFamilyAccount(record[columnMap["family_id"]]),
			DeviceName:     record[columnMap["device_name"]],
			DeviceMAC:      record[columnMap["device_mac"]],
			ScannedAt:      scannedAt,
			Latitude:       latitude,
			Longitude:      longitude,
			GPSAccuracy:    gpsAccuracy,
			ReasonCode:     reasonCode,
			PhoneName:      record[columnMap["phone_name"]],
			SenderDeviceID: record[columnMap["sender_device_id"]],
		}

		points = append(points, point)
	}

	return points, nil
}

// generateKML generates a KML file from the provided points
func generateKML(points []Point, outputPath string, pathMode bool) error {
	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write KML header
	_, err = file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Care Active Scan Locations</name>
    <description>GPS locations from scan data</description>
    
`)
	if err != nil {
		return fmt.Errorf("failed to write KML header: %w", err)
	}

	// Add styles based on the mode
	if pathMode {
		// Add both path style and point styles for path mode
		_, err = file.WriteString(`    <!-- Style for path lines -->
    <Style id="pathStyle">
      <LineStyle>
        <color>ffff6633</color>
        <width>4</width>
      </LineStyle>
    </Style>
    
    <!-- Styles for points in path mode -->
    <Style id="reasonCode0">
      <IconStyle>
        <color>ffffffff</color>
        <scale>0.7</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/red-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
    <Style id="reasonCode1">
      <IconStyle>
        <color>ff00ff00</color>
        <scale>0.7</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/grn-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
    <Style id="reasonCode2">
      <IconStyle>
        <color>ff00ffff</color>
        <scale>0.7</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/ylw-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
    <Style id="reasonCodeDefault">
      <IconStyle>
        <color>ffffff00</color>
        <scale>0.7</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/wht-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
`)
	} else {
		// Add point styles for placemark mode
		_, err = file.WriteString(`    <!-- Styles for different reason codes -->
    <Style id="reasonCode0">
      <IconStyle>
        <color>ffffffff</color>
        <scale>1.0</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/red-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
    <Style id="reasonCode1">
      <IconStyle>
        <color>ff00ff00</color>
        <scale>1.0</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/grn-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
    <Style id="reasonCode2">
      <IconStyle>
        <color>ff00ffff</color>
        <scale>1.0</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/ylw-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
    <Style id="reasonCodeDefault">
      <IconStyle>
        <color>ffffff00</color>
        <scale>1.0</scale>
        <Icon>
          <href>http://maps.google.com/mapfiles/kml/paddle/wht-circle.png</href>
        </Icon>
      </IconStyle>
    </Style>
`)
	}

	if err != nil {
		return fmt.Errorf("failed to write KML styles: %w", err)
	}

	if pathMode {
		// Generate paths connecting points in sequence
		if len(points) == 1 {
			// If there's only one point in path mode, only create the placemark, skip the path
			fmt.Println("Only one point found, skipping path generation")
		} else if len(points) < 1 {
			return fmt.Errorf("at least 1 point is required")
		} else {
			// First, create path segments between consecutive points if we have more than one point
			for i := 0; i < len(points)-1; i++ {
				startPoint := points[i]
				endPoint := points[i+1]

				pathSegment := fmt.Sprintf(`    <Placemark>
      <name>Path %d-%d</name>
      <description>
        <![CDATA[
        <strong>From:</strong> %s at %s<br/>
        <strong>To:</strong> %s at %s<br/>
        <strong>Family:</strong> %s<br/>
        ]]>
      </description>
      <styleUrl>#pathStyle</styleUrl>
      <LineString>
        <extrude>1</extrude>
        <tessellate>1</tessellate>
        <coordinates>
          %f,%f,0
          %f,%f,0
        </coordinates>
      </LineString>
    </Placemark>
`,
					i, i+1,
					startPoint.DeviceName, startPoint.ScannedAt.Format(time.RFC3339),
					endPoint.DeviceName, endPoint.ScannedAt.Format(time.RFC3339),
					startPoint.FamilyAccount,
					startPoint.Longitude, startPoint.Latitude,
					endPoint.Longitude, endPoint.Latitude,
				)

				_, err = file.WriteString(pathSegment)
				if err != nil {
					return fmt.Errorf("failed to write path segment: %w", err)
				}
			}
		}

		// Always add individual points along the path (even if there's only one)
		for i, point := range points {
			styleID := "reasonCodeDefault"
			switch point.ReasonCode {
			case 0:
				styleID = "reasonCode0"
			case 1:
				styleID = "reasonCode1"
			case 2:
				styleID = "reasonCode2"
			}

			placemark := fmt.Sprintf(`    <Placemark>
      <name>%s (Point %d)</name>
      <description>
        <![CDATA[
        <strong>Family:</strong> %s<br/>
        <strong>Device:</strong> %s (%s)<br/>
        <strong>Scanned At:</strong> %s<br/>
        <strong>Longtitude:</strong> %f<br/>
        <strong>Latitude:</strong> %f<br/>
        <strong>GPS Accuracy:</strong> %d meters<br/>
        <strong>Phone:</strong> %s<br/>
        ]]>
      </description>
      <styleUrl>#%s</styleUrl>
      <Point>
        <coordinates>%f,%f,0</coordinates>
      </Point>
      <TimeStamp>
        <when>%s</when>
      </TimeStamp>
    </Placemark>
`,
				point.DeviceName, i,
				point.FamilyAccount,
				point.DeviceName,
				point.DeviceMAC,
				point.ScannedAt.Format(time.RFC3339),
				point.Longitude,
				point.Latitude,
				point.GPSAccuracy,
				point.PhoneName,
				styleID,
				point.Longitude,
				point.Latitude,
				point.ScannedAt.Format(time.RFC3339),
			)

			_, err = file.WriteString(placemark)
			if err != nil {
				return fmt.Errorf("failed to write placemark: %w", err)
			}
		}
	} else {
		// Write each point as a placemark (original behavior)
		for _, point := range points {
			styleID := "reasonCodeDefault"
			switch point.ReasonCode {
			case 0:
				styleID = "reasonCode0"
			case 1:
				styleID = "reasonCode1"
			case 2:
				styleID = "reasonCode2"
			}

			placemark := fmt.Sprintf(`    <Placemark>
      <name>%s</name>
      <description>
        <![CDATA[
        <strong>Family:</strong> %s<br/>
        <strong>Device:</strong> %s (%s)<br/>
        <strong>Scanned At:</strong> %s<br/>
        <strong>Longtitude:</strong> %f<br/>
        <strong>Latitude:</strong> %f<br/>
        <strong>GPS Accuracy:</strong> %d meters<br/>
        <strong>Phone:</strong> %s<br/>
        ]]>
      </description>
      <styleUrl>#%s</styleUrl>
      <Point>
        <coordinates>%f,%f,0</coordinates>
      </Point>
      <TimeStamp>
        <when>%s</when>
      </TimeStamp>
    </Placemark>
`,
				point.DeviceName,
				point.FamilyAccount,
				point.DeviceName,
				point.DeviceMAC,
				point.ScannedAt.Format(time.RFC3339),
				point.Longitude,
				point.Latitude,
				point.GPSAccuracy,
				point.PhoneName,
				styleID,
				point.Longitude,
				point.Latitude,
				point.ScannedAt.Format(time.RFC3339),
			)

			_, err = file.WriteString(placemark)
			if err != nil {
				return fmt.Errorf("failed to write placemark: %w", err)
			}
		}
	}

	// Write KML footer
	_, err = file.WriteString(`  </Document>
</kml>
`)
	if err != nil {
		return fmt.Errorf("failed to write KML footer: %w", err)
	}

	return nil
}
