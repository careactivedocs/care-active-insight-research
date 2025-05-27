#
# This program was AI-assisted code generation.
#
# Copyright (c) 2025
# Care Active Corp. ("CA").
# All rights reserved.
#
# The information contained herein is confidential and proprietary to
# CA. Use of this information by anyone other than authorized employees
# of CA is granted only under a written non-disclosure agreement,
# expressly prescribing the scope and manner of such use.
#

"""  
CSV Time Interval Aggregation Tool

This program allows users to input a time interval (in seconds) and
aggregates the resampled CSV data according to that interval. The result includes:
- Start and end acvl values for each interval
- Sum of acvl increments within the interval (aggregate_acvl_increment)
- Maximum value of each RSSI field within the interval
- User and device information (family_id, device_name, device_mac)

Two operation modes are supported:
1. Single file mode: Process a single CSV file
2. Directory mode: Process all CSV files in a directory (batch processing)

Usage:
    # Single file mode:
    python time-aggregator.py --file input.csv interval [--output output.csv] [--verify]

    # Directory mode:
    python time-aggregator.py --dir input_directory interval [--output-dir output_directory] [--verify]

Parameters:
    --file: Path to a single CSV file to process
    --dir: Directory containing multiple CSV files to process in batch
    interval: Time interval in seconds for aggregation
    --output: Output file name (for single file mode)
    --output-dir: Output directory (for directory mode)
    --verify: Verify interval validity against data time span
"""

import csv
import sys
import os
import datetime
import argparse
import pandas as pd
import numpy as np
from pathlib import Path
import glob


# Add room type mapping
ROOM_TYPE_MAP = {
    "00": "Unknown",
    "01": "MainBedroom",
    "02": "Bedroom1",
    "03": "Bedroom2",
    "04": "LivingRoom",
    "05": "FamilyRoom",
    "06": "RecreationRoom",
    "07": "Kitchen",
    "08": "DiningRoom",
    "09": "Den/Office",
    "10": "MasterBathroom",
    "11": "Bathroom",
    "12": "Garage",
    "13": "Patio",
    "14": "Entryway",
    "15": "Other"
}


def aggregate_by_time_interval(input_file, interval_seconds, output_file=None):
    """
    Aggregate CSV data based on specified time interval

    Parameters:
        input_file (str): Path to the input CSV file
        interval_seconds (int): Time interval in seconds
        output_file (str, optional): Path to the output file

    Returns:
        str: Path to the output file
    """
    print(f"Reading file: {input_file}")
    print(f"Using time interval: {interval_seconds} seconds")

    # Use Pandas to read CSV
    try:
        df = pd.read_csv(input_file)
        print(f"Successfully read CSV, {len(df)} rows")
    except Exception as e:
        print(f"Error: Cannot read CSV file - {e}")
        sys.exit(1)

    # Ensure required columns exist
    required_columns = ["resample_at", "resample_at_utc", "resample_acvl", "resample_acvl_increment"]
    for column in required_columns:
        if column not in df.columns:
            print(f"Error: CSV file is missing required column: {column}")
            sys.exit(1)
    
    # Check for user and device information columns
    identity_columns = ["family_id", "device_name", "device_mac"]
    missing_identity_columns = [col for col in identity_columns if col not in df.columns]
    if missing_identity_columns:
        print(f"Warning: The following identity columns are missing: {', '.join(missing_identity_columns)}")
        # Add missing columns with empty values
        for col in missing_identity_columns:
            df[col] = ""

    # Identify all RSSI columns
    rssi_columns = [col for col in df.columns if col.endswith("_rssi")]
    if not rssi_columns:
        print("Warning: No RSSI columns found (should end with _rssi)")

    print(f"Found {len(rssi_columns)} RSSI columns")
    
    # Extract room number from RSSI column names
    # Format is expected to be like: "01_744DBD2A408C_rssi"
    rssi_room_map = {}
    for col in rssi_columns:
        parts = col.split('_')
        if len(parts) >= 3:  # Ensure format is correct
            room = parts[0]
            mac = parts[1]
            rssi_room_map[col] = f"{room}_{mac}"

    # Create time groups
    # Use different logic for intervals > 60 seconds
    
    if interval_seconds <= 60:
        # For intervals <= 60 seconds - group by minute alignment
        # Calculate seconds position within each minute
        df['seconds_in_minute'] = df['resample_at'] % 60
        
        # Align to the start of the minute (0th second)
        df['aligned_time'] = df['resample_at'] - df['seconds_in_minute']
        
        # Calculate which interval group within the minute
        df['interval_group_in_minute'] = (df['seconds_in_minute'] // interval_seconds).astype(int)
        
        # Combine aligned_time and interval_group_in_minute to create unique time group identifier
        df['time_group'] = df['aligned_time'].astype(str) + '_' + df['interval_group_in_minute'].astype(str)
    else:
        # For intervals > 60 seconds - group directly by UNIX timestamp
        # Align to the floor division of timestamp by interval
        df['time_group'] = (df['resample_at'] // interval_seconds).astype(int)
        
    # Clean up temporary columns
    for col in ['seconds_in_minute', 'aligned_time', 'interval_group_in_minute']:
        if col in df.columns:
            df = df.drop(columns=[col])

    # Set aggregate_at before aggregation
    # Use original timestamp as aggregate_at
    df['aggregate_at'] = df['resample_at']
    
    # Aggregation functions
    agg_functions = {
        'resample_at': ['first', 'last'],  # First and last timestamp in group
        'resample_at_utc': ['first', 'last'],  # First and last UTC time in group
        'resample_acvl': ['first', 'last'],  # First and last acvl in group
        'resample_acvl_increment': 'sum',  # Sum of all acvl_increments in group
        'aggregate_at': 'last',  # Use the last timestamp in group as aggregate_at
        'family_id': 'first',  # Keep family_id in output
        'device_name': 'first',  # Keep device_name in output
        'device_mac': 'first'   # Keep device_mac in output
    }

    # Add max aggregation for each RSSI column
    for col in rssi_columns:
        agg_functions[col] = 'max'

    # Perform aggregation
    print("Performing time interval aggregation...")
    grouped_df = df.groupby('time_group').agg(agg_functions)

    # Flatten multi-level column names
    flattened_columns = []
    for col in grouped_df.columns:
        if col[0].endswith('_rssi'):
            flattened_columns.append(col[0])
        else:
            flattened_columns.append('_'.join(col).strip())
    grouped_df.columns = flattened_columns

    # Add max_rssi_station, max_rssi_room_type_id and max_rssi_room_type_name columns
    def get_max_rssi_data(row):
        max_rssi_value = -np.inf
        max_rssi_column = None
        for col in rssi_columns:
            if col in row and row[col] > max_rssi_value:
                max_rssi_value = row[col]
                max_rssi_column = col
        
        # Return empty strings when no valid RSSI value is found
        if max_rssi_value == -np.inf:
            return pd.Series({
                'max_rssi_room_type_id': '',
                'max_rssi_room_type_name': '',
                'max_rssi_station': ''
            })
        
        # Process when RSSI value is found
        if max_rssi_column in rssi_room_map:
            room_mac = rssi_room_map[max_rssi_column]
            parts = room_mac.split('_')
            if len(parts) >= 2:
                room_id = parts[0]
                return pd.Series({
                    'max_rssi_room_type_id': room_id,
                    'max_rssi_room_type_name': ROOM_TYPE_MAP.get(room_id, "Unknown"),
                    'max_rssi_station': parts[1]
                })
        
        # Return empty strings when room_mac format cannot be parsed
        return pd.Series({
            'max_rssi_room_type_id': '',
            'max_rssi_room_type_name': '',
            'max_rssi_station': ''
        })

    max_rssi_data = grouped_df.apply(get_max_rssi_data, axis=1)
    grouped_df['max_rssi_room_type_id'] = max_rssi_data['max_rssi_room_type_id']
    grouped_df['max_rssi_room_type_name'] = max_rssi_data['max_rssi_room_type_name']
    grouped_df['max_rssi_station'] = max_rssi_data['max_rssi_station']

    # Rename columns to more friendly format
    column_mapping = {
        'resample_at_first': 'start_time',
        'resample_at_last': 'end_time',
        'resample_at_utc_first': 'start_time_utc',
        'resample_at_utc_last': 'end_time_utc',
        'resample_acvl_first': 'start_acvl',
        'resample_acvl_last': 'end_acvl',
        'resample_acvl_increment_sum': 'aggregate_acvl_increment',
        'aggregate_at_last': 'aggregate_at',  # Changed from aggregated_at_last to aggregate_at_last
        'family_id_first': 'family_id',      # Keep original column name for family_id
        'device_name_first': 'device_name',    # Keep original column name for device_name
        'device_mac_first': 'device_mac'       # Keep original column name for device_mac
    }

    # Rename columns
    grouped_df = grouped_df.rename(columns=column_mapping)
    
    # Create aggregate_at_utc column
    # Get date format template from start_time_utc
    # Then generate UTC time string using aggregate_at timestamp
    
    # First, find the format of resample_at_utc
    # Assume format is consistent, we get it from the first row
    sample_utc = df['resample_at_utc'].iloc[0] if not df.empty else ''
    
    # Create aggregate_at_utc column
    def create_utc_time(timestamp):
        try:
            # Convert timestamp to UTC datetime
            dt = datetime.datetime.utcfromtimestamp(timestamp)  # Changed to utcfromtimestamp
            # Try to match the original UTC time format
            if '/' in sample_utc:  # Like 2025/03/19 07:28:45
                return dt.strftime('%Y/%m/%d %H:%M:%S')
            else:  # Other possible formats
                return dt.strftime('%Y-%m-%d %H:%M:%S')
        except Exception as e:
            print(f"Warning: Cannot create UTC time - {e}")
            return None
    
    # 由於我們重命名了 'aggregated_at_last' 為 'aggregate_at'，需要使用這個新列名
    grouped_df['aggregate_at_utc'] = grouped_df['aggregate_at'].apply(create_utc_time)
    
    # Add interval column
    grouped_df['interval_seconds'] = interval_seconds

    # Reset index to make time_group a regular column
    grouped_df = grouped_df.reset_index()
    
    # Rename time_group column to aggregate_time_group
    if 'time_group' in grouped_df.columns:
        grouped_df = grouped_df.rename(columns={'time_group': 'aggregate_time_group'})
        
    # Sort by start_time to ensure output data is in chronological order
    grouped_df = grouped_df.sort_values('start_time')
    
    # Reorder columns according to specified order
    ordered_columns = [
        'family_id', 'device_name', 'device_mac',
        'aggregate_time_group', 'interval_seconds',
        'start_time', 'end_time', 'aggregate_at',
        'start_time_utc', 'end_time_utc', 'aggregate_at_utc',
        'start_acvl', 'end_acvl', 'aggregate_acvl_increment', 'max_rssi_room_type_id', 'max_rssi_room_type_name', 'max_rssi_station'
    ]
    
    # Add RSSI columns to the end (preserving their order)
    rssi_columns = [col for col in grouped_df.columns if col.endswith('_rssi')]
    ordered_columns.extend(rssi_columns)
    
    # Reorder using only columns that exist in the DataFrame
    existing_ordered_columns = [col for col in ordered_columns if col in grouped_df.columns]
    grouped_df = grouped_df[existing_ordered_columns]

    # If no output file is specified, create a default name without any special directory
    if output_file is None:
        timestamp = datetime.datetime.now().strftime("%Y%m%d%H%M%S")
        file_base = os.path.basename(input_file).split('.')[0]
        output_file = f"aggregate_{file_base}_{interval_seconds}sec_{timestamp}.csv"

    # Save to CSV
    grouped_df.to_csv(output_file, index=False)
    print(f"Aggregation complete! Processed {len(df)} original records, aggregated to {len(grouped_df)} records")
    print(f"Aggregation ratio: {len(df) / len(grouped_df):.2f}:1")
    print(f"Output file: {output_file}")

    # Display some statistics
    print("\nStatistics:")
    print(f"First record time: {grouped_df['start_time_utc'].iloc[0]}")
    print(f"Last record time: {grouped_df['end_time_utc'].iloc[-1]}")

    return output_file


def process_directory(directory_path, interval_seconds, output_dir=None, verify=False):
    """Process all CSV files in a directory
    
    Parameters:
        directory_path (str): Path to directory containing CSV files
        interval_seconds (int): Time interval in seconds
        output_dir (str, optional): Directory to save output files
        verify (bool): Whether to verify interval validity
        
    Returns:
        list: List of output file paths generated
    """
    # Check if the directory exists
    if not os.path.exists(directory_path) or not os.path.isdir(directory_path):
        print(f"Error: Directory not found: {directory_path}")
        sys.exit(1)
    
    # Find all CSV files in the directory
    csv_files = glob.glob(os.path.join(directory_path, "*.csv"))
    if not csv_files:
        print(f"Error: No CSV files found in directory: {directory_path}")
        sys.exit(1)
    
    print(f"Found {len(csv_files)} CSV files in: {directory_path}")
    
    # Create output directory if not specified, use specified naming format
    if not output_dir:
        base_name = os.path.basename(os.path.normpath(directory_path))
        output_dir = os.path.join(directory_path, f"aggregated_{interval_seconds}s_{base_name}")
        
    # Create output directory if it doesn't exist
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
        print(f"Created output directory: {output_dir}")
    
    output_files = []
    
    # Process each CSV file
    for i, csv_file in enumerate(csv_files):
        print(f"\nProcessing file {i+1}/{len(csv_files)}: {os.path.basename(csv_file)}")
        
        # Generate output filename
        base_filename = os.path.basename(csv_file).split('.')[0]
        timestamp = datetime.datetime.now().strftime("%Y%m%d%H%M%S")
        output_file = os.path.join(output_dir, f"aggregate_{base_filename}_{interval_seconds}sec_{timestamp}.csv")
        
        # Verify interval validity if requested
        if verify:
            total_seconds = get_total_time_span(csv_file)
            if total_seconds % interval_seconds != 0:
                print(f"Warning: {interval_seconds} seconds cannot evenly divide the total time span of {total_seconds} seconds")
                print(f"Suggested intervals: {suggest_intervals(total_seconds)}")
                
                # Skip file if verification fails 
                print(f"Skipping file: {os.path.basename(csv_file)}")
                continue
        
        # Process file
        try:
            result_file = aggregate_by_time_interval(csv_file, interval_seconds, output_file)
            output_files.append(result_file)
        except Exception as e:
            print(f"Error processing file {csv_file}: {str(e)}")
            continue
    
    print(f"\nBatch processing complete. Processed {len(output_files)} out of {len(csv_files)} files.")
    print(f"Output files saved to: {output_dir}")
    return output_files


def main():
    """Main function that processes command line arguments and starts the aggregation process"""
    parser = argparse.ArgumentParser(description='Aggregate CSV data based on time intervals')
    parser.add_argument('--file', dest='input_file', help='Input CSV file')
    parser.add_argument('--dir', dest='input_dir', help='Directory containing multiple CSV files')
    parser.add_argument('interval', type=int, help='Time interval (seconds)')
    parser.add_argument('--output', dest='output_file', help='Output file (for single file mode)')
    parser.add_argument('--output-dir', dest='output_dir', help='Output directory (for directory mode)')
    parser.add_argument('--verify', action='store_true', help='Verify interval validity')
    
    args = parser.parse_args()
    
    # Check that either --file or --dir is provided
    if not args.input_file and not args.input_dir:
        print("Error: Either --file or --dir must be specified")
        parser.print_help()
        sys.exit(1)
    
    # Create list of common intervals
    # Values that can evenly divide a minute (less than or equal to 60)
    valid_minute_divisible_intervals = [1, 2, 3, 4, 5, 6, 10, 12, 15, 20, 30, 60]
    
    # Multiples of minutes (greater than 60)
    valid_minute_multiple_intervals = [
        90, 120, 150, 180, 210, 240, 270, 300, 330, 360, 390, 420, 450, 480, 510, 540, 570, 600,
        900, 1200, 1500, 1800, 2100, 2400, 2700, 3000, 3300, 3600, 7200, 10800, 14400, 21600, 43200, 86400  # Including multiples of hours and days
    ]
    
    valid_intervals = valid_minute_divisible_intervals + valid_minute_multiple_intervals
    
    # Validate interval value
    if args.interval <= 0:
        print("Error: Interval must be a positive integer")
        sys.exit(1)
    
    # Verify if the interval meets the rules
    if args.interval <= 60:
        # Intervals less than or equal to 60 must evenly divide 60
        if 60 % args.interval != 0:
            print(f"Error: Interval {args.interval} seconds cannot evenly divide 60 seconds.")
            print(f"Please use one of these values that can evenly divide a minute: {', '.join(map(str, valid_minute_divisible_intervals))}")
            sys.exit(1)
    else:
        # Intervals greater than 60 must be multiples of 30
        if args.interval % 30 != 0:
            print(f"Error: Interval {args.interval} seconds is not a multiple of 30.")
            print(f"Please use one of these values:")
            print(f"- Values that can evenly divide a minute: {', '.join(map(str, valid_minute_divisible_intervals))}")
            print(f"- Multiples of 30 seconds: {', '.join(map(str, [i for i in valid_minute_multiple_intervals if i % 30 == 0]))}")
            sys.exit(1)
    
    # Process based on mode (single file or directory)
    if args.input_file:
        # Single file mode
        
        # Check if input file exists
        if not os.path.exists(args.input_file):
            print(f"Error: Input file not found: {args.input_file}")
            sys.exit(1)
        
        # If interval validity verification is requested
        if args.verify:
            total_seconds = get_total_time_span(args.input_file)
            if total_seconds % args.interval != 0:
                print(f"Warning: {args.interval} seconds cannot evenly divide the total time span of {total_seconds} seconds")
                print(f"Suggested intervals: {suggest_intervals(total_seconds)}")
                
                # Ask whether to continue
                response = input("Do you still want to continue? (y/n): ")
                if response.lower() != 'y':
                    sys.exit(0)
        
        # Execute aggregation for a single file
        aggregate_by_time_interval(args.input_file, args.interval, args.output_file)
    
    else:
        # Directory mode
        process_directory(args.input_dir, args.interval, args.output_dir, args.verify)


def get_total_time_span(file_path):
    """
    Get the total time span of CSV data (in seconds)
    
    Parameters:
        file_path (str): Path to the CSV file
        
    Returns:
        int: Total time span (seconds)
    """
    try:
        df = pd.read_csv(file_path)
        min_time = df['resample_at'].min()
        max_time = df['resample_at'].max()
        return max_time - min_time
    except Exception as e:
        print(f"Warning: Cannot determine time span - {e}")
        return 0


def suggest_intervals(total_seconds):
    """
    Suggest intervals that can evenly divide the total time
    
    Parameters:
        total_seconds (int): Total time span (seconds)
        
    Returns:
        list: List of suggested intervals
    """
    # Common intervals
    common_intervals = [1, 2, 3, 4, 5, 6, 10, 12, 15, 20, 30, 60, 120, 300, 600, 900, 1800, 3600]
    
    # Find intervals that can evenly divide the total time
    valid_intervals = []
    for interval in common_intervals:
        if total_seconds % interval == 0:
            valid_intervals.append(interval)
    
    return valid_intervals


if __name__ == "__main__":
    main()
