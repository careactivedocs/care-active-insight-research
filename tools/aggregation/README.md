# Time Interval Aggregation Tool

## Overview

The Time Interval Aggregation Tool is designed to process resampled rtls-cvl data. This tool aggregates per-second resampled data into larger time intervals as specified by the user. During aggregation, it identifies the station with the strongest signal within each time interval to provide location information. The aggregated output significantly reduces data volume, making it more suitable for graphical visualization and analysis of movement patterns over time.

## Features

- **Flexible Time Intervals**: Aggregate data by any valid time interval (in seconds)
- **Comprehensive Data Aggregation**:
  - Captures start and end values for each interval
  - Sums numerical increments within each interval
  - Records maximum RSSI values

## Requirements

- Python 3.11.5
- Dependencies:
  - pandas
  - numpy
  - pathlib


## Usage

### Single File Mode

Process a single CSV file with a specified time interval:

```bash
python time-aggregator.py --file input.csv interval [--output output.csv] [--verify]
```

### Directory Mode (Batch Processing)

Process all CSV files in a directory:

```bash
python time-aggregator.py --dir input_directory interval [--output-dir output_directory] [--verify]
```

### Parameters

- `--file`: Path to a single CSV file to process
- `--dir`: Directory containing multiple CSV files to process in batch
- `interval`: Time interval in seconds for aggregation
- `--output`: Output file name (for single file mode)
- `--output-dir`: Output directory (for directory mode)
- `--verify`: Verify interval validity against data time span

## Input File Requirements

The tool expects CSV files with the following required columns:
- `resample_at`: Timestamp (UNIX format)
- `resample_at_utc`: Human-readable UTC timestamp
- `resample_acvl`: Activity level value
- `resample_acvl_increment`: Incremental activity value

Optional but recommended columns:
- `family_id`: User identifier
- `device_name`: Device name
- `device_mac`: Device MAC address
- Any columns ending with `_rssi`: Signal strength indicators

## Interval Rules

For optimal results, use these time intervals:

- **≤ 60 seconds**: Must evenly divide a minute (1, 2, 3, 4, 5, 6, 10, 12, 15, 20, 30, 60)
- **> 60 seconds**: Must be multiples of 30 seconds (90, 120, 150, 180, etc.)


## Output Format

The tool generates CSV files with the following columns:

- **User/Device Info**: `family_id`, `device_name`, `device_mac`
- **Time Information**:
  - `aggregate_time_group`: Unique identifier for each time group
  - `interval_seconds`: The interval used for aggregation
  - `start_time`, `end_time`, `aggregate_at`: UNIX timestamps
  - `start_time_utc`, `end_time_utc`, `aggregate_at_utc`: Human-readable UTC timestamps
- **Data Values**:
  - `start_acvl`, `end_acvl`: Activity values at interval boundaries
  - `aggregate_acvl_increment`: Sum of incremental activity within interval
  - `max_rssi_station`: Station with strongest signal during interval
  - All original RSSI columns with maximum values

## Examples

### Basic Usage

```bash
# Process a single file using a 30-second interval
python time-aggregator.py --file data.csv 30

# Process all files in a directory using a 5-minute interval
python time-aggregator.py --dir data_folder 300
```


## License

Copyright (c) 2025 Care Active Corp. ("CA"). All rights reserved.

The information contained herein is confidential and proprietary to CA. Use of this information by anyone other than authorized employees of CA is granted only under a written non-disclosure agreement, expressly prescribing the scope and manner of such use.