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
	"time"
)

// linearInterpolate performs linear interpolation between two points.
// x1, y1: first point
// x2, y2: second point
// x: the x-coordinate to interpolate at
// Returns the interpolated y-coordinate
func linearInterpolate(x1, y1, x2, y2, x float64) float64 {
	// If x1 and x2 are the same, return y1 to avoid division by zero
	if x1 == x2 {
		return y1
	}

	// Linear interpolation formula: y = y1 + (x - x1) * (y2 - y1) / (x2 - x1)
	return y1 + (x-x1)*(y2-y1)/(x2-x1)
}

// formatTimestamp converts a Unix timestamp (seconds) to a human-readable string in UTC.
func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01-02 15:04:05")
}

// timestampToSecond converts a millisecond timestamp to a second timestamp.
func timestampToSecond(timestamp int64) int64 {
	return timestamp / 1000
}

// findNearestRecord finds the index of the record with the timestamp
// closest to the target timestamp.
func findNearestRecord(records []Record, targetTimestamp int64) int {
	if len(records) == 0 {
		return -1
	}

	nearestIdx := 0
	minDiff := abs(records[0].SampleAtMs - targetTimestamp)

	for i := 1; i < len(records); i++ {
		diff := abs(records[i].SampleAtMs - targetTimestamp)
		if diff < minDiff {
			minDiff = diff
			nearestIdx = i
		}
	}

	return nearestIdx
}

// abs returns the absolute value of an int64.
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
