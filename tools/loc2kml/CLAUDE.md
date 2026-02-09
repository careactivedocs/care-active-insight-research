# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

loc2kml is a Go CLI tool that converts Care Active location CSV reports into KML files for viewing in Google Earth. It is part of a larger monorepo at `care-active-insight-research/tools/loc2kml`. The tool is primarily used to regenerate KML files with different timezone settings from location CSV data.

## Build Commands

The project uses a Makefile with cross-compilation targets. There is no `go.mod`; it's a single-file Go program using only the standard library.

```bash
# Build all platforms (linux, macOS Intel, macOS ARM, Windows)
make all

# Build for a specific platform
make build-macos-arm    # macOS ARM (Apple Silicon)
make build-macos        # macOS Intel
make build-linux        # Linux amd64
make build-windows      # Windows amd64

# Clean build outputs
make clean
```

Binaries go to `bin/<platform>/`. There are no tests or lint targets.

To build and run manually:
```bash
go build -o loc2kml ./src/loc2kml.go
./loc2kml -dir=. -output=result.kml -timezone=America/Toronto
```

## Architecture

The entire tool is a single Go source file: `src/loc2kml.go`. It has no external dependencies.

**Data flow:** CSV input → parse rows into `Point` structs → generate KML XML output via string formatting (no XML library).

**Key structures and functions:**
- `Point` struct: holds parsed location data (coordinates, timestamps, device info, optional battery/refresh fields)
- `processCSVFile`: parses CSV using a column-name-to-index map, making it resilient to column ordering. Required columns: `family_id`, `device_name`, `device_mac`, `scanned_at_ms`, `gps_latitude`, `gps_longitude`, `gps_accuracy`, `phone_name`, `sender_device_id`. Optional columns: `sender_battery`, `loc_refresh_at_ms`
- `generateKML`: writes KML XML directly via `fmt.Sprintf` string templates. Supports two modes: individual placemarks (default) and connected path mode (`-path` flag)
- `formatFamilyAccount`: transforms `f1c+NNNN@careactive.ai` email format into `F1CNNNN` display format

**Two input modes:** directory mode (`-dir`) processes all CSVs in a folder; file mode accepts CSV paths as positional arguments.

**No GPS data is not an error:** when CSVs contain no valid GPS points (common with GeoUpdate reports where Full GPS Log is disabled), the tool exits gracefully without generating a KML file.

## CSV Format

`sample_input.csv` has Full GPS Log data (produces KML). `sample_input_2.csv` has non-Full-GPS-Log data (no KML output, not an error).
