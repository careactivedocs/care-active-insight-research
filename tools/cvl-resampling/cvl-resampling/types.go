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

// This file contains type definitions that might be shared across multiple files.

// CSVHeader represents column indices in the input CSV for quick access.
type CSVHeader struct {
	UserEmail       int
	DeviceName      int
	DeviceMAC       int
	SampleAtMs      int
	SampleAtUTC     int
	SampleAtSource  int
	SampleAtStation int
	EventSeq        int
	SeqDiff         int
	ACVL            int
	CVLIncrement    int
	XYZCount        int
	RefStationID    int
	RefScannedAtMs  int
	RSSIIndices     map[string]int // Maps RSSI column names to their indices
}

// NewCSVHeader creates a new CSVHeader instance from the headers slice.
func NewCSVHeader(headers []string) CSVHeader {
	header := CSVHeader{
		RSSIIndices: make(map[string]int),
	}

	for i, h := range headers {
		switch h {
		case "family_id":
			header.UserEmail = i
		case "device_name":
			header.DeviceName = i
		case "device_mac":
			header.DeviceMAC = i
		case "sample_at_ms":
			header.SampleAtMs = i
		case "sample_at_utc":
			header.SampleAtUTC = i
		case "sample_at_source":
			header.SampleAtSource = i
		case "sample_at_station":
			header.SampleAtStation = i
		case "event_seq":
			header.EventSeq = i
		case "seq_diff":
			header.SeqDiff = i
		case "acvl":
			header.ACVL = i
		case "cvl_increment":
			header.CVLIncrement = i
		case "xyz_count":
			header.XYZCount = i
		case "ref_station_id":
			header.RefStationID = i
		case "ref_scanned_at_ms":
			header.RefScannedAtMs = i
		default:
			// Check if it's an RSSI column
			if len(h) > 5 && h[len(h)-5:] == "_rssi" {
				header.RSSIIndices[h] = i
			}
		}
	}

	return header
}
