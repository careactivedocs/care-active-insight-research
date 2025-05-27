# CSV Local Time Converter

## Overview

CSV Local Time Converter is a command-line tool designed to process CSV files containing UTC timestamps and add human-readable local time columns. This tool is particularly useful for analyzing time-series data across different time zones.

The converter identifies timestamp columns in CSV files, converts the Unix millisecond timestamps to human-readable datetime strings in the specified timezone, and outputs a new CSV with additional `local_time` and `local_timezone` columns.

## Features

- Automatically detects timestamp columns with priority order: `aggregate_at` > `resample_at` > `resample_at_ms` > `sample_at_ms` > `scanned_at_ms` > `created_at_ms`
- Converts UTC millisecond timestamps to any specified timezone
- Formats local time to match the format of `sample_at_utc` when available
- Adds both `local_time` and `local_timezone` columns to the output
- Preserves all original data
- Supports standard IANA TZ identifiers (e.g., "Asia/Taipei", "America/Toronto")
- Simple command-line interface
- Processes files from various sources: single file, directory, or standard input

## Usage

```bash
./csv-localtime-converter -input <input_file.csv> -tz <timezone>
```

or

```bash
./csv-localtime-converter -dir <directory_with_csvs> -tz <timezone>
```

or process from standard input:

```bash
cat input.csv | ./csv-localtime-converter -tz <timezone>
```

### Command Line Arguments

- `-input`: Path to a single input CSV file
- `-dir`: Path to a directory containing CSV files to process
- `-output`: Path for the output file (optional)
- `-tz` (required): Timezone identifier
  - Examples: "Asia/Taipei", "America/Toronto", "Europe/London"

#### Timezone

The location code of timezone is defined by IANA. Refer to the IANA tz database for the available options.

[List of tz database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)

### Examples

Process a single file:
```bash
./csv-localtime-converter -input sensor_data.csv -tz Asia/Taipei
```

Process all CSV files in a directory:
```bash
./csv-localtime-converter -dir /path/to/csv/files -tz America/New_York
```

Process a file and specify output location:
```bash
./csv-localtime-converter -input sensor_data.csv -output converted_data.csv -tz Europe/London
```

Process from standard input:
```bash
cat sensor_data.csv | ./csv-localtime-converter -tz Asia/Tokyo
```

## Output Format

The tool creates a new CSV file with:
- Filename prefixed with "local_" (e.g., "local_sensor_data.csv") unless an output path is specified
- A new `local_time` column as the first column
  - Format matches `sample_at_utc` if available, without the 'Z' suffix
  - Otherwise uses ISO8601 format without the 'Z' suffix
- A new `local_timezone` column as the second column
  - Contains the timezone identifier used for conversion
- All original columns preserved in their original order

## Supported CSV Formats

The tool can process various CSV formats with the following timestamp columns (checked in priority order):
1. `aggregate_at` (highest priority) - UTC time in seconds since the Unix epoch 
2. `resample_at` - UTC time in seconds since the Unix epoch
3. `resample_at_ms` - UTC time in milliseconds since the Unix epoch
4. `sample_at_ms` - UTC time in milliseconds since the Unix epoch
5. `scanned_at_ms` - UTC time in milliseconds since the Unix epoch
6. `created_at_ms` (lowest priority) - UTC time in milliseconds since the Unix epoch

Timestamp values can be:
- Formatted UTC time strings
- Seconds since the Unix epoch
- Milliseconds since the Unix epoch

## License

Copyright (c) 2025 Care Active Corp. ("CA"). All rights reserved.

The information contained herein is confidential and proprietary to CA. Use of this information by anyone other than authorized employees of CA is granted only under a written non-disclosure agreement, expressly prescribing the scope and manner of such use.
