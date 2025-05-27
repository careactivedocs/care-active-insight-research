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
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// DataProcessor handles the resampling logic for the input data.
type DataProcessor struct {
	rssiColumns []string
	gapLimit    int64
	rssiLimit   int64
	verbose     bool // Controls progress reporting
}

// NewDataProcessor creates a new DataProcessor instance.
func NewDataProcessor(rssiColumns []string, gapLimit, rssiLimit int64, verbose bool) (*DataProcessor, error) {
	// Validate that rssiLimit >= gapLimit
	if rssiLimit < gapLimit {
		return nil, fmt.Errorf("rssiLimit (%d) must be greater than or equal to gapLimit (%d)", rssiLimit, gapLimit)
	}

	return &DataProcessor{
		rssiColumns: rssiColumns,
		gapLimit:    gapLimit,
		rssiLimit:   rssiLimit,
		verbose:     verbose,
	}, nil
}

// ProcessData resamples the input data according to the specified rules.
func (p *DataProcessor) ProcessData(records [][]string, headers []string) ([][]string, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no records to process")
	}

	// Create index maps for easy access to column indices
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	// Ensure required columns exist
	requiredColumns := []string{
		"family_id", "device_name", "device_mac", "sample_at_ms", "event_seq", "acvl",
	}
	for _, col := range requiredColumns {
		if _, exists := headerMap[col]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in input CSV", col)
		}
	}

	// Parse records
	parsedRecords, err := p.parseRecords(records, headerMap)
	if err != nil {
		return nil, err
	}

	// Resample data
	resampledData, err := p.resampleData(parsedRecords, headerMap)
	if err != nil {
		return nil, err
	}

	// Format the resampled data
	return p.formatResampledData(resampledData, headerMap), nil
}

// Record represents a parsed record from the input CSV.
type Record struct {
	UserEmail  string
	DeviceName string
	DeviceMAC  string
	SampleAtMs int64
	EventSeq   int
	ACVL       int
	RSSIValues map[string]float64
}

// ResampledPoint represents a resampled data point.
type ResampledPoint struct {
	ResampleAt            int64
	UserEmail             string
	DeviceName            string
	DeviceMAC             string
	EventSeq              int
	ACVL                  int // Original ACVL value from the corresponding event_seq
	ResampleACVL          float64
	ResampleACVLIncrement float64
	ResampledRSSI         map[string]float64
}

// TimeDataPoint holds a timestamp and its associated value
type TimeDataPoint struct {
	Timestamp int64
	Value     float64
}

// parseRecords converts the string records to structured data.
func (p *DataProcessor) parseRecords(records [][]string, headerMap map[string]int) ([]Record, error) {
	var parsedRecords []Record

	for i, record := range records {
		// Skip records that don't have enough fields
		if len(record) < len(headerMap) {
			continue
		}

		// Parse record
		sampleAtMs, err := strconv.ParseInt(record[headerMap["sample_at_ms"]], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid sample_at_ms at row %d: %v", i+2, err)
		}

		eventSeq, err := strconv.Atoi(record[headerMap["event_seq"]])
		if err != nil {
			return nil, fmt.Errorf("invalid event_seq at row %d: %v", i+2, err)
		}

		acvl, err := strconv.Atoi(record[headerMap["acvl"]])
		if err != nil {
			return nil, fmt.Errorf("invalid acvl at row %d: %v", i+2, err)
		}

		// Parse RSSI values
		rssiValues := make(map[string]float64)

		for _, rssiCol := range p.rssiColumns {
			if headerMap[rssiCol] >= len(record) {
				continue
			}

			// Skip empty or invalid RSSI values
			if record[headerMap[rssiCol]] == "" {
				continue
			}

			rssiVal, err := strconv.ParseFloat(record[headerMap[rssiCol]], 64)
			if err != nil {
				continue
			}

			rssiValues[rssiCol] = rssiVal
		}

		parsedRecords = append(parsedRecords, Record{
			UserEmail:  record[headerMap["family_id"]],
			DeviceName: record[headerMap["device_name"]],
			DeviceMAC:  record[headerMap["device_mac"]],
			SampleAtMs: sampleAtMs,
			EventSeq:   eventSeq,
			ACVL:       acvl,
			RSSIValues: rssiValues,
		})
	}

	return parsedRecords, nil
}

// createEventSeqMap builds a map of event_seq to its corresponding ACVL value for accurate reference
func createEventSeqMap(records []Record) map[int]int {
	eventSeqToACVL := make(map[int]int)

	// Map each event_seq to its ACVL value
	for _, record := range records {
		eventSeqToACVL[record.EventSeq] = record.ACVL
	}

	return eventSeqToACVL
}

// buildRssiTimeSeriesMap creates a map of RSSI time series for each Station-MAC
// The map is structured as: stationMAC -> ordered list of (timestamp, value) points
func buildRssiTimeSeriesMap(records []Record, rssiColumns []string) map[string][]TimeDataPoint {
	// Initialize map for each RSSI column
	rssiTimeSeriesMap := make(map[string][]TimeDataPoint)

	for _, rssiCol := range rssiColumns {
		var timeSeries []TimeDataPoint

		// Collect all timestamps and values for this RSSI column
		for _, record := range records {
			if val, exists := record.RSSIValues[rssiCol]; exists {
				timeSeries = append(timeSeries, TimeDataPoint{
					Timestamp: record.SampleAtMs,
					Value:     val,
				})
			}
		}

		// Sort time series by timestamp
		sort.Slice(timeSeries, func(i, j int) bool {
			return timeSeries[i].Timestamp < timeSeries[j].Timestamp
		})

		rssiTimeSeriesMap[rssiCol] = timeSeries
	}

	return rssiTimeSeriesMap
}

// interpolateRssi performs linear interpolation between two RSSI points
func interpolateRssi(p1, p2 TimeDataPoint, targetTime int64) float64 {
	// If timestamps are the same, avoid division by zero
	if p1.Timestamp == p2.Timestamp {
		return p1.Value
	}

	ratio := float64(targetTime-p1.Timestamp) / float64(p2.Timestamp-p1.Timestamp)
	return p1.Value + ratio*(p2.Value-p1.Value)
}

// needInterpolation checks if we can interpolate RSSI values between two timestamps
// It returns false if the gap is too large (exceeding rssiLimit)
func needInterpolation(t1, t2 int64, rssiLimit int64) bool {
	return abs(t2-t1) <= rssiLimit
}

// isPointWithinLimit checks if a timestamp is within rssiLimit of a target timestamp
func isPointWithinLimit(point, target int64, rssiLimit int64) bool {
	return abs(point-target) <= rssiLimit
}

// getRssiValueAtTimestamp returns the interpolated RSSI value at a given timestamp
// based on the time series data, considering the rssiLimit constraint
func getRssiValueAtTimestamp(timeSeries []TimeDataPoint, timestamp int64, rssiLimit int64) (float64, bool) {
	if len(timeSeries) == 0 {
		return 0, false
	}

	// Check if any point falls in the same second as the target timestamp
	targetSecond := timestamp / 1000

	for _, point := range timeSeries {
		// If any point is in the same second, use it regardless of other conditions
		if point.Timestamp/1000 == targetSecond {
			return point.Value, true
		}
	}

	// Handle case with only one data point
	if len(timeSeries) == 1 {
		// For a single RSSI value, use it if it's within rssiLimit
		if isPointWithinLimit(timeSeries[0].Timestamp, timestamp, rssiLimit) {
			return timeSeries[0].Value, true
		}
		return 0, false
	}

	// Find the points before and after the target timestamp
	var beforeIdx, afterIdx = -1, -1

	for i, point := range timeSeries {
		if point.Timestamp <= timestamp {
			beforeIdx = i
		} else {
			afterIdx = i
			break
		}
	}

	// Case 1: Target is before the first data point
	if beforeIdx == -1 {
		// Only use the first point if it's within rssiLimit
		if afterIdx >= 0 && isPointWithinLimit(timeSeries[afterIdx].Timestamp, timestamp, rssiLimit) {
			return timeSeries[afterIdx].Value, true
		}
		return 0, false
	}

	// Case 2: Target is after the last data point
	if afterIdx == -1 {
		// Only use the last point if it's within rssiLimit
		if isPointWithinLimit(timeSeries[beforeIdx].Timestamp, timestamp, rssiLimit) {
			return timeSeries[beforeIdx].Value, true
		}
		return 0, false
	}

	// Case 3: Target is between two data points
	// Check first if the gap is too large between these points
	timeDiff := timeSeries[afterIdx].Timestamp - timeSeries[beforeIdx].Timestamp
	if timeDiff >= rssiLimit {
		// Gap is too large or equal to rssiLimit - no interpolation allowed at all
		// Per requirement, if gaps >= rssiLimit, no interpolation should be done
		return 0, false
	}

	// Case 4: Target is between two close enough points - interpolate
	return interpolateRssi(timeSeries[beforeIdx], timeSeries[afterIdx], timestamp), true
}

// resampleData performs the resampling of the parsed records.
func (p *DataProcessor) resampleData(records []Record, headerMap map[string]int) ([]ResampledPoint, error) {
	if len(records) == 0 {
		return nil, nil
	}

	var resampledPoints []ResampledPoint
	var lastResampleACVL float64 = 0 // Track last ACVL value for increment calculation
	var isFirstPoint bool = true     // Flag to identify the first point

	// Create a map of event_seq to ACVL for accurate reference
	eventSeqToACVL := createEventSeqMap(records)

	// Build RSSI time series map for each Station-MAC
	rssiTimeSeriesMap := buildRssiTimeSeriesMap(records, p.rssiColumns)

	// Sort records by time if needed
	// For simplicity, we assume records are already sorted by sample_at_ms

	// Find the earliest and latest timestamps
	startTime := records[0].SampleAtMs
	endTime := records[len(records)-1].SampleAtMs

	// Convert to seconds and make sure we're on second boundaries
	startTimeSec := (startTime / 1000) * 1000
	endTimeSec := ((endTime / 1000) + 1) * 1000

	// Create resampled points at 1-second intervals
	for timestamp := startTimeSec; timestamp <= endTimeSec; timestamp += 1000 {
		// Find the records before and after this timestamp for ACVL interpolation
		var beforeIdx, afterIdx int = -1, -1
		for i, record := range records {
			if record.SampleAtMs <= timestamp {
				beforeIdx = i
			}
			if record.SampleAtMs > timestamp && afterIdx == -1 {
				afterIdx = i
				break
			}
		}

		// Skip if no data before this timestamp
		if beforeIdx == -1 {
			continue
		}

		before := records[beforeIdx]

		// If we're exactly on a data point, use its values directly
		if before.SampleAtMs == timestamp {
			// Calculate ACVL increment
			acvlValue := float64(before.ACVL)
			acvlIncrement := 0.0

			// For first point, increment is always 0
			if !isFirstPoint {
				acvlIncrement = acvlValue - lastResampleACVL
			} else {
				isFirstPoint = false
			}

			lastResampleACVL = acvlValue

			resampledPoint := ResampledPoint{
				ResampleAt:            timestamp / 1000, // Convert to seconds
				UserEmail:             before.UserEmail,
				DeviceName:            before.DeviceName,
				DeviceMAC:             before.DeviceMAC,
				EventSeq:              before.EventSeq,
				ACVL:                  before.ACVL, // Use the exact ACVL value from the original record
				ResampleACVL:          acvlValue,
				ResampleACVLIncrement: acvlIncrement,
				ResampledRSSI:         make(map[string]float64),
			}

			// Copy RSSI values directly
			for rssiCol, value := range before.RSSIValues {
				resampledPoint.ResampledRSSI[rssiCol] = value
			}

			resampledPoints = append(resampledPoints, resampledPoint)
			continue
		}

		// Skip if no data after this timestamp
		if afterIdx == -1 {
			continue
		}

		after := records[afterIdx]

		// Check if the gap is too large
		if after.SampleAtMs-before.SampleAtMs > p.gapLimit {
			continue
		}

		// Check if ACVL decreased (it should only increase or stay the same)
		if after.ACVL < before.ACVL {
			continue
		}

		// Calculate time ratio for interpolation
		timeRatio := float64(timestamp-before.SampleAtMs) / float64(after.SampleAtMs-before.SampleAtMs)

		// Interpolate ACVL
		interpolatedACVL := float64(before.ACVL) + timeRatio*float64(after.ACVL-before.ACVL)

		// Determine which event_seq to use based on proximity
		var eventSeq int
		if (timestamp - before.SampleAtMs) < (after.SampleAtMs - timestamp) {
			eventSeq = before.EventSeq
		} else {
			eventSeq = after.EventSeq
		}

		// Get the exact ACVL value for this event_seq from the original data
		originalACVL := eventSeqToACVL[eventSeq]

		// Calculate ACVL increment
		acvlIncrement := 0.0

		// For first point, increment is always 0
		if !isFirstPoint {
			acvlIncrement = interpolatedACVL - lastResampleACVL
		} else {
			isFirstPoint = false
		}

		// lastResampleACVL shall be the rouding value of 2-digit interpolatedACVL
		lastResampleACVL = math.Round(interpolatedACVL*100) / 100

		resampledPoint := ResampledPoint{
			ResampleAt:            timestamp / 1000, // Convert to seconds
			UserEmail:             before.UserEmail,
			DeviceName:            before.DeviceName,
			DeviceMAC:             before.DeviceMAC,
			EventSeq:              eventSeq,     // Use the closest event_seq
			ACVL:                  originalACVL, // Use the exact ACVL for this event_seq from original data
			ResampleACVL:          interpolatedACVL,
			ResampleACVLIncrement: acvlIncrement,
			ResampledRSSI:         make(map[string]float64),
		}

		// Handle each RSSI column independently using the time series data
		for rssiCol, timeSeries := range rssiTimeSeriesMap {
			if rssiValue, exists := getRssiValueAtTimestamp(timeSeries, timestamp, p.rssiLimit); exists {
				resampledPoint.ResampledRSSI[rssiCol] = rssiValue
			}
			// If no valid interpolation, we leave this Station-MAC empty
		}

		resampledPoints = append(resampledPoints, resampledPoint)
	}

	return resampledPoints, nil
}

// generateEmptyRecord creates an empty record for the given parameters.
func (p *DataProcessor) generateEmptyRecord(firstRow ResampledPoint, emptyTimer int64) []string {
	record := make([]string, 9+len(p.rssiColumns))
	record[0] = firstRow.UserEmail  // family_id
	record[1] = firstRow.DeviceName // device_name
	record[2] = firstRow.DeviceMAC  // device_mac
	record[3] = strconv.FormatInt(emptyTimer, 10)
	t := time.Unix(emptyTimer, 0).UTC()
	record[4] = t.Format("2006/01/02 15:04:05")
	record[5] = strconv.Itoa(firstRow.EventSeq)                        // event_seq
	record[6] = strconv.Itoa(firstRow.ACVL)                            // acvl
	record[7] = strconv.FormatFloat(firstRow.ResampleACVL, 'f', 2, 64) // resample_acvl
	record[8] = "-0.00"                                                // resample_acvl_increment

	// RSSI values
	for j := 0; j < len(p.rssiColumns); j++ {
		record[9+j] = "" // Empty for no value
	}

	return record
}

// formatResampledData converts the resampled points to string records for CSV output.
func (p *DataProcessor) formatResampledData(resampledPoints []ResampledPoint, headerMap map[string]int) [][]string {
	var formattedData [][]string
	var lastResampleAt int64 = 0 // Track last resample timestamp

	if len(resampledPoints) == 0 {
		return formattedData
	}

	// Progress reporting
	if p.verbose {
		fmt.Printf("Formatting %d resampled points...\n", len(resampledPoints))
	}

	firstRow := resampledPoints[0]
	firstResampleAt := firstRow.ResampleAt

	// restart the resample tracking
	lastResampleAt = firstResampleAt

	// Always fill the entire day with empty records
	// convert to Unix timestamp
	fraTime := time.Unix(firstResampleAt, 0).UTC()

	// get the start of the day
	fraTime = time.Date(fraTime.Year(), fraTime.Month(), fraTime.Day(), 0, 0, 0, 0, time.UTC)

	// get the EPOCH value int64
	firstResampleStart := fraTime.Unix()

	// Progress reporting for filling empty records
	if p.verbose && firstResampleAt > firstResampleStart {
		fmt.Printf("Filling %d empty records from start of day...\n", firstResampleAt-firstResampleStart)
	}

	// generate empty records from 00:00:00 of the day to the first resample point
	if firstResampleAt > 0 && firstResampleStart < firstResampleAt {
		for emptyTimer := firstResampleStart; emptyTimer < firstResampleAt; emptyTimer++ {
			formattedData = append(formattedData, p.generateEmptyRecord(firstRow, emptyTimer))
			lastResampleAt = emptyTimer
		}
	}

	// Process actual resampled points
	pointCount := len(resampledPoints)
	for i, point := range resampledPoints {
		// Progress reporting every 10% of records
		if p.verbose && pointCount > 100 && i%(pointCount/10) == 0 {
			fmt.Printf("Processing point %d of %d (%.0f%%)...\n", i, pointCount, float64(i)/float64(pointCount)*100)
		}

		// Create a formatted record
		record := make([]string, 9+len(p.rssiColumns))

		// Format data in the exact order specified:
		// family_id, device_name, device_mac, resample_at, resample_at_utc, event_seq, acvl, resample_acvl, resample_acvl_increment
		record[0] = point.UserEmail
		record[1] = point.DeviceName
		record[2] = point.DeviceMAC

		// Timestamp columns
		record[3] = strconv.FormatInt(point.ResampleAt, 10)

		// Human-readable timestamp
		t := time.Unix(point.ResampleAt, 0).UTC()
		record[4] = t.Format("2006/01/02 15:04:05")

		// Event sequence and ACVL
		record[5] = strconv.Itoa(point.EventSeq)
		record[6] = strconv.Itoa(point.ACVL)

		// Resampled ACVL and increment
		record[7] = strconv.FormatFloat(point.ResampleACVL, 'f', 2, 64)

		// check if this is the first point of a new resampling segment
		if point.ResampleAt == (lastResampleAt + 1) {
			record[8] = strconv.FormatFloat(point.ResampleACVLIncrement, 'f', 2, 64)
		} else {
			// indicate this is the first row of the new resampling segment
			record[8] = "-0.00"

			// Always fill all gaps between records
			// connect the time with empty records before adding the new point
			for emptyTimer := lastResampleAt + 1; emptyTimer < point.ResampleAt; emptyTimer++ {
				formattedData = append(formattedData, p.generateEmptyRecord(point, emptyTimer))
			}
		}

		// track the last resample timestamp
		lastResampleAt = point.ResampleAt

		// RSSI values
		for i, rssiCol := range p.rssiColumns {
			if val, exists := point.ResampledRSSI[rssiCol]; exists {
				record[9+i] = strconv.FormatFloat(val, 'f', 2, 64)
			} else {
				record[9+i] = "" // Empty for no value
			}
		}

		formattedData = append(formattedData, record)
	}

	// Always fill the remaining day with empty records
	// get the EPOCH value of the last second of the day
	lastResampleAtTime := time.Unix(lastResampleAt, 0).UTC()
	lastResampleAtTime = time.Date(lastResampleAtTime.Year(), lastResampleAtTime.Month(), lastResampleAtTime.Day(), 23, 59, 59, 0, time.UTC)
	lastResampleAtEnd := lastResampleAtTime.Unix()

	// Progress reporting
	if p.verbose && lastResampleAtEnd > lastResampleAt {
		fmt.Printf("Filling %d empty records to end of day...\n", lastResampleAtEnd-lastResampleAt)
	}

	// generate empty records from the last resample point to the end of the day
	for emptyTimer := lastResampleAt + 1; emptyTimer <= lastResampleAtEnd; emptyTimer++ {
		formattedData = append(formattedData, p.generateEmptyRecord(resampledPoints[len(resampledPoints)-1], emptyTimer))
	}

	return formattedData
}
