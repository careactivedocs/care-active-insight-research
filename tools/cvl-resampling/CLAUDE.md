# CVL-Resampling Tool Guidelines

## Build Commands
- Build for all platforms: `make all`
- Build for Linux: `make build-linux`
- Build for macOS (Intel): `make build-macos`
- Build for macOS (ARM): `make build-macos-arm`
- Clean build artifacts: `make clean`
- Run the tool: `go run ./cvl-resampling/*.go -input [input.csv] -output [output.csv]`
- Build manually: `go build -o cvl-resampling ./cvl-resampling`

## Go Style Guidelines
- Format code with: `go fmt ./...`
- Verify code quality: `go vet ./...`
- Prefer descriptive variable names and camelCase for variables/functions 
- Each function should have a descriptive comment
- Return errors instead of handling them internally when appropriate
- Use struct methods for related functionality
- Always check error returns
- Types should be defined in types.go
- Copyright header required on all files
- Process data in modular steps: read → process → write
- Strong error handling with descriptive messages

## Project Summary

### Files and Structure
- `cvl-resampling/main.go`: Command-line flags and main execution flow
- `cvl-resampling/data_processor.go`: Core data processing and interpolation logic
- `cvl-resampling/csv_io.go`: CSV input/output handling
- `cvl-resampling/types.go`: Data structure definitions
- `cvl-resampling/utils.go`: Helper functions

### Key Components
- The application processes time-series RSSI (Received Signal Strength Indicator) data
- Resampling time-series data at 1-second intervals
- Linear interpolation for RSSI values with configurable gap limits
- Independent handling of each Station-MAC's RSSI values
- CSV field order: `family_id, device_name, device_mac, resample_at, resample_at_utc, event_seq, acvl, resample_acvl`

### Key Functions
- `getRssiValueAtTimestamp`: Core interpolation function respecting time gap constraints
- `buildRssiTimeSeriesMap`: Organizes data by Station-MAC for interpolation
- `resampleData`: Main resampling loop for creating 1-second interval data
- `isPointWithinLimit`: Helper function to check if points are within interpolation limit
- `needInterpolation`: Determines if interpolation should be applied

### Implementation Notes
- RSSI values are only interpolated when points are within rssiLimit of each other
- Large time gaps (>15000ms) prevent interpolation for that segment
- Each Station-MAC's RSSI values are handled independently
- Interpolation respects the beginning and end points of the time series
- ACVL (Activity Level) values use exact values from closest event_seq

### Current Status
- Latest work involved completely rewriting the RSSI interpolation algorithm to correctly handle time gaps
- Improved organization and readability through helper functions
- Fixed issues with interpolation across large time gaps

### Next Steps
- Thoroughly test with real data, focusing on interpolation behavior around rssiLimit boundaries
- Test various input data with different gap patterns
- Verify no interpolation occurs for gaps larger than rssiLimit
- Confirm exact values are used when timestamp falls on data point
- Consider visualization of interpolation results to validate algorithm
- Add comprehensive tests to prevent regression

## Original Requirements

- In the example CSV in the knowledge base, the fields before `ref_scanned_at_ms` are fixed and remain the same regardless of which device generated the data. However, the fields after `01_{Station_MAC}_rssi` will vary depending on different devices.
- This CSV serves as input data.
- There is another input called `-gaplimit`, which is a numerical value measured in milliseconds (ms), referred to as the gap limit, which will be used later.
- Another input is `-rssilimit`, also a numerical value measured in milliseconds (ms), called the rssi limit, which will be used later.
- `header_ai.go` contains the header to be used for all the produced Golang code.
- The tool is developed in Golang and adheres to the coding style found at this link: https://google.github.io/styleguide/go/
- This tool will output another CSV format data, resampling the acvl and all rssi values to a time series of every 1 second.
- The original time sequence is `sample_at_ms`, which is in EPOCH format.
- The output data needs to retain the following fields: `family_id`, `device_name`, `device_mac`, filled with their original values.
- Add `resample_at`, which is in EPOCH format to the second.
- Add `resample_at_utc`, which is the human-readable format of `resample_at`, in UTC timezone, so no additional timezone indication is needed.
- The `event_seq` field is the value from the closest original `sample_at_ms` data row before resampling, used for future cross-validation.
- The following fields are not needed in the new data table: `sample_at_source`, `sample_at_station`, `sample_at_ms`, `sample_at_utc`, `cvl_increment`, `xyz_count`, `ref_station_id`, `ref_scanned_at_ms`, `seq_diff`.
- Values that need to be resampled are `acvl` and all `rssi` values.
- Add `resample_acvl`, which is the resampled `acvl` value.
- The `acvl` value in the new table is filled with the `acvl` value from the original row of that `event_seq`.
- All rssi values are resampled values.
- The rules for resampling `acvl` are as follows:
1. If the `sample_at_ms` of the next row exceeds the gap limit, skip the interpolation or extrapolation for this segment.
2. The `acvl` value should only increase or remain the same over time, never decrease. If it decreases, skip interpolation or extrapolation for this segment and omit it directly.
3. Each time series must have a `resample_acvl` value.
- The rules for resampling rssi are as follows:
1. If the `sample_at_ms` of the next row exceeds the gap limit, skip the interpolation or extrapolation for this segment.
2. The format for rssi value fields is explained as follows: `{sequence number}_{Station-MAC}_rssi`. If there is no value in this field in the table, it means that the Station-MAC did not receive any signal.
3. If a specific Station-MAC has no value for longer than the rssi limit time, it indicates that the Station truly did not receive a signal rather than simply losing packets. Therefore, do not perform any interpolation or extrapolation for this Station-MAC afterward; otherwise, fill in the middle time segments through interpolation or extrapolation.
- Use linear interpolation or extrapolation.

  <analysis>
  Let me go through the conversation and identify key points:

  1. **Primary Request:**
     - The initial request was to analyze the codebase and create a `CLAUDE.md` file containing build/lint/test commands and code style guidelines.
     - Later, issues with the codebase around RSSI data interpolation were identified.

  2. **Key Technical Concepts:**
     - Go programming for data processing
     - CSV data reading and writing
     - Data resampling and interpolation algorithms
     - RSSI signal handling and time-series data processing
     - Time gap constraints in interpolation
     - Makefile configuration for cross-platform Go builds

  3. **Files and Code Sections:**
     - Main files:
       - data_processor.go (extensive modifications)
       - main.go
       - csv_io.go
       - types.go
       - utils.go
       - Makefile
     - Key functions:
       - getRssiValueAtTimestamp (most significant changes)
       - resampleData
       - parseRecords
       - formatResampledData
       - buildRssiTimeSeriesMap (newly added)
       - interpolateRssi
     - New helper functions:
       - isPointWithinLimit
       - needInterpolation

  4. **Problem Solving:**
     - Fixed Makefile that was incorrect (paths, binary names)
     - Fixed unused imports in Go files
     - Fixed RSSI interpolation logic for respecting rssiLimit
     - Addressed issues with data field ordering
     - Implemented correct handling of large time gaps between RSSI values

  5. **Pending Tasks:**
     - Testing the new RSSI interpolation logic thoroughly
     - Potential optimizations for performance
     - Further edge case handling

  6. **Current Work:**
     - Just completed a major refactoring of the RSSI interpolation logic
     - Added helper functions to make the code more readable and maintainable
     - Fixed specific issues with timestamp gap handling

  7. **Next Step:**
     - Testing the current solution with additional data points
     - Potential refinements to the interpolation algorithm
     - Code organization and documentation improvements

  The work has been substantial, particularly around the RSSI interpolation algorithm, which needed to be completely rewritten to correctly handle gaps between data points.
  </analysis>

  <summary>
  1. Primary Request and Intent:
     The user initially requested to create a CLAUDE.md file containing build commands and code style guidelines for a Go-based CSV data processing tool. The underlying intent was to provide documentation for agentic coding systems operating in 
  this repository. As the conversation progressed, the focus shifted to fixing significant functional issues in the data interpolation algorithm, specifically around handling RSSI data with time gaps.

  2. Key Technical Concepts:
     - Go language data processing and CSV file handling
     - Linear interpolation for time-series data
     - RSSI (Received Signal Strength Indicator) data processing with Station-MACs
     - Time-based data resampling at fixed intervals (1-second)
     - Gap handling with configurable thresholds (gapLimit and rssiLimit)
     - Cross-platform Go compilation (Linux, macOS Intel, macOS ARM)
     - ACVL (Activity Level) value processing and interpolation
     - Makefile-based build system for Go applications

  3. Files and Code Sections:
     - cvl-resampling/data_processor.go: Major refactoring of data processing logic
       - getRssiValueAtTimestamp(): Core function for determining RSSI values
       - resampleData(): Resamples time-series data at 1-second intervals
       - buildRssiTimeSeriesMap(): New function to organize data by Station-MAC
     - cvl-resampling/main.go: Command-line flags and main execution flow
     - cvl-resampling/csv_io.go: CSV input/output handling, header formatting
     - cvl-resampling/types.go: Data structure definitions
     - cvl-resampling/utils.go: Helper functions like abs() and interpolation
     - Makefile: Build configurations for Linux, macOS Intel, and macOS ARM
     - CLAUDE.md: Documentation file with build commands and style guidelines

  4. Problem Solving:
     - Fixed incorrect Makefile that referenced the wrong source files and binary name
     - Removed unused imports from Go files causing build errors
     - Reordered CSV field output to match exact requirements (family_id, device_name, device_mac, resample_at, resample_at_utc, event_seq, acvl, resample_acvl)
     - Completely rewrote RSSI interpolation logic to correctly handle:
       - Large time gaps (>15000ms) between data points
       - Isolated data points without interpolation partners
       - Independent handling of each Station-MAC's RSSI values
       - Proper boundary conditions at the start and end of the time series
     - Implemented correct mapping of event_seq to original ACVL values

  5. Pending Tasks:
     - Comprehensive testing of RSSI interpolation with more edge cases
     - Testing with the full sample dataset to verify all interpolation rules 
     - Performance optimization for larger datasets
     - Additional error handling improvements
     - More thorough documentation of the algorithm's behavior

  6. Current Work:
     The most recent work involved completely rewriting the RSSI interpolation algorithm to correctly handle time gaps between data points. The key changes include:
     - Adding helper functions isPointWithinLimit and needInterpolation to improve code clarity
     - Implementing a more sophisticated gap detection logic in getRssiValueAtTimestamp
     - Ensuring that RSSI values are only interpolated when points are within rssiLimit of each other
     - Creating a more robust time series handling approach with separate data structures
     - Adding precise checks to prevent interpolation across large time gaps

  7. Next Step Recommendation:
     The most logical next step is to thoroughly test the updated implementation with real data, particularly focusing on the interpolation behavior around the boundaries of the rssiLimit value. This should include:
     1. Testing with a variety of input data covering different gap patterns
     2. Verifying that no interpolation occurs for Station-MACs with gaps larger than rssiLimit
     3. Confirming that when a timestamp falls exactly on a data point, the exact values are used
     4. Creating visualization of the interpolation results to visually validate the algorithm behavior
     5. Adding comprehensive automated tests to prevent regression in this complex interpolation logic
  </summary>.
  