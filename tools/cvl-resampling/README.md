# CVL Resampling Tool

A tool that resamples RSSI and ACVL time-series data into 1-second intervals using linear interpolation. It ensures full 24-hour coverage by creating empty records for gaps between resampling sections.

## Purpose

The CVL Resampling Tool processes raw time-series data collected from RTLS (Real-Time Location System) sensors and resamples it to generate a uniform 1-second interval dataset. This tool is particularly useful for:

- Converting irregularly sampled RSSI data into regular time intervals
- Applying linear interpolation to fill gaps in sensor data
- Ensuring data consistency for downstream analysis tools
- Processing activity level values alongside RSSI measurements

## Usage

### Basic Usage

```bash
./cvl-resampling -input [input.csv]
```

### Command Line Parameters

- `-input` (required): Path to the input CSV file containing raw sensor data
- `-output` (optional): Path to the output CSV file. If not specified, defaults to `resampled_[input_filename].csv`
- `-gaplimit` (optional): Maximum time gap in milliseconds allowed for interpolation (default: 30000)
- `-rssilimit` (optional): Maximum time gap in milliseconds allowed for RSSI signal interpolation (default: 30000)
- `-verbose` (optional): Enable verbose output for debugging purposes

#### gaplimit and rssilimit

Considering the non-RTLS broadcasting, with 15 seconds of pedometer and pedo-beacon data for each 10 and 53 minute interval respectively, the gap limitation is better than these two non-RTLS broadcasting methods. The current default values for gaplimit and rssilimit are both 30000 ms.

### Example

```bash
./cvl-resampling -input sample_rtls_cvl.csv -output resampled_data.csv -gaplimit 10000 -rssilimit 8000
```

## Key Features

- Resamples time-series data at consistent 1-second intervals
- Applies linear interpolation for RSSI values with configurable gap limits
- Handles each Station-MAC's RSSI values independently
- Preserves original user identification and device information
- Includes human-readable timestamps alongside epoch times

## Input Data Requirements

The input CSV file should contain:
- User identification fields (`family_id`, `device_name`, `device_mac`)
- Timestamp fields (`sample_at_ms` in EPOCH format)
- ACVL (Activity Level) values 
- RSSI values in the format `{sequence_number}_{Station-MAC}_rssi`
- All `sample_at_ms` values must not be zero

## Output Data Structure

The output CSV file contains:
1. `family_id`, `device_name`, `device_mac` (preserved from input)
2. `resample_at` (timestamp in EPOCH format)
3. `resample_at_utc` (human-readable UTC timestamp)
4. `event_seq` (sequence ID from original data)
5. `acvl` (original activity level value)
6. `resample_acvl` (resampled activity level)
7. `resample_acvl_increment` (increment of the resampled activity level)
8. RSSI values for each Station-MAC (resampled)

### -0.00 of resample_acvl_increment

If "-0.00" is used, it means this is the first row of the resampling section after the gaplimit is reached, including the very first row of the entire resampling data.

## Restrictions and Limitations

1. **Interpolation Gap Limits**:
   - No interpolation is performed across time gaps larger than `gaplimit` (default: 15000ms)
   - Station-MACs with no signal for longer than `rssilimit` are considered to have no signal rather than missing packets

2. **ACVL Value Rules**:
   - ACVL values must never decrease over time
   - Each time series point must have a `resample_acvl` value
   - No interpolation is performed if the next row exceeds the gap limit

3. **RSSI Value Rules**:
   - No interpolation is performed across time gaps larger than the gap limit
   - Each Station-MAC is handled independently
   - If a Station-MAC has no value for longer than the RSSI limit, no further interpolation is performed

## Interpolation Verification Tool

linear-interpolation.go is a small tool to help to verify the interpolation results. A Golang interpreter is needed to run it.

## License

Copyright (c) 2025 Care Active Corp. All rights reserved.