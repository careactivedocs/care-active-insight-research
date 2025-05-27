package main

import (
	"fmt"
)

// linearInterpolation performs linear interpolation or extrapolation calculation
// parameters valueA and valueB can be negative
func linearInterpolation(epochMsA int64, valueA int64, epochMsB int64, valueB int64, targetMs int64) float64 {
	// Check if denominator is zero (when time points are the same)
	if epochMsA == epochMsB {
		// If time points are identical, return one of the values
		return float64(valueA)
	}

	// Calculate slope
	slope := float64(valueB-valueA) / float64(epochMsB-epochMsA)

	// Calculate time difference between target time and epochMsA
	deltaTime := float64(targetMs - epochMsA)

	// Linear interpolation or extrapolation calculation
	interpolatedValue := float64(valueA) + slope*deltaTime

	return interpolatedValue
}

func main() {
	var epochMsA, valueA, epochMsB, valueB, targetS, targetMs int64

	fmt.Println("Please enter the following values:")
	fmt.Print("epoch_ms_a: ")
	fmt.Scan(&epochMsA)

	fmt.Print("value_a (can be negative): ")
	fmt.Scan(&valueA)

	fmt.Print("epoch_ms_b: ")
	fmt.Scan(&epochMsB)

	fmt.Print("value_b (can be negative): ")
	fmt.Scan(&valueB)

	// Loop to continuously ask for target_ms until 0 is entered
	for {
		fmt.Print("target in second (enter 0 to exit): ")
		fmt.Scan(&targetS)

		// Exit condition
		if targetS == 0 {
			fmt.Println("Program terminated.")
			break
		}

		targetMs = targetS * 1000 // Convert seconds to milliseconds
		fmt.Printf("target_ms: %d\n", targetMs)

		result := linearInterpolation(epochMsA, valueA, epochMsB, valueB, targetMs)

		fmt.Println("\nCalculation Result:")

		// Determine if it's interpolation or extrapolation
		if targetMs >= epochMsA && targetMs <= epochMsB || targetMs >= epochMsB && targetMs <= epochMsA {
			fmt.Printf("Interpolation result: %.6f\n", result)
		} else {
			fmt.Printf("Extrapolation result: %.6f\n", result)
		}

		fmt.Printf("Value at time point %d: %.6f\n", targetMs, result)

		// Display calculation process for verification
		fmt.Println("\nCalculation Process:")
		fmt.Printf("Time ratio: (%.0f - %.0f) / (%.0f - %.0f) = %.6f\n",
			float64(targetMs), float64(epochMsA), float64(epochMsB), float64(epochMsA),
			float64(targetMs-epochMsA)/float64(epochMsB-epochMsA))
		fmt.Printf("Value difference: %.0f - %.0f = %.0f\n", float64(valueB), float64(valueA), float64(valueB-valueA))
		fmt.Printf("Formula: %.0f + %.6f * %.0f = %.6f\n",
			float64(valueA),
			float64(valueB-valueA)/float64(epochMsB-epochMsA),
			float64(targetMs-epochMsA),
			result)

		fmt.Println("\n--------------------------------------------------")
	}
}
