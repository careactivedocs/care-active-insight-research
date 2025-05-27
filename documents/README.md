# G3MR Researcher's Kits

This tool kits are for the researchers to handle the G3MR data including the Activity, Indoor Location, GPS, and data of Care Active sensors.

## Data Path

| File Tree               | Type of Data                                   |
| :---------------------- | :--------------------------------------------- |
| exported_data_v2        | ROOT                                           |
| ../{collector}          | All data of {collector}                        |
| ../../maintenance       | Maintenance data of all devices of {collector} |
| ../../{family_account}  | All data of {family_account}                   |
| ../../../logs           | Device name change logs                        |
| ../../../post-processed | Daily post-processed data                      |
| ../../../../location    | GPS location records in CSV                    |
| ../../../../motion      | Motion sensor activity records in CSV          |
| ../../../../pedo        | Pedometer readings in CSV                      |
| ../../../../rtls-cvl    | post-processed activity report in CSV          |
| ../../../../rtls-rssi   | Station-watch RSSI table in CSV                |
| ../../../../rtls-xyz    | Original motion raw records in CSV             |
| ../../../YYYYMM         | Data generated in this month                   |
| ../../source            | Aggrgate raw data in JSON format               |

## Data Format of Collected Data

| Tools/Documents                                             | Descriptions                                       |
| :---------------------------------------------------------- | :------------------------------------------------- |
| [Maintenance data](./maintenance-reports.md)                | Daily list of the status of all registered devices |
| [Device name logs](./device-name-logs.md)                   | Device name change logs                            |
| [Daily rtls-activity data](./daily-rtls-activity-report.md) | Daily RTLS Activity data                           |
| [Daily post-processed data](./daily-pp-csv-report.md)       | Daily non-RTLS data                                |
| [Daily source data](daily-source-data.md)                   | Data format descriptions of the daily raw data     |

## Document of Processing Tools

| Tools/Documents                                                                                  | Descriptions                                                                            |
| :----------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------- |
| [Care Active Insight](./ca-insight.md)                                                           | CA Insight is the web portal for the researchers to manage and view the device and data |
| [Care Active Watch Advertising](./Care%20Active%20Watch%20Advertising%20Formats%20G3MR%20v2.pdf) | The basic concepts of how Care Active Watch generates the Location-Activity data        |
| [csv-local-time](../tools/csv_local_time/README.md)                                              | Convert UTC timestamp in CSV file to local timestamp                                    |
| [loc2kml](../tools/loc2kml/README.md)                                                            | Convert location data to KML for map viewing                                            |
| [cvl-resample](../tools/cvl-resampling/README.md)                                                | Resample RSSI-CVL to be per-second basis. This is the same tool running in the cloud.   |
| [cvl-aggregation](../tools/aggregation/README.md)                                                | This tool aggregates per-second resampled data into larger time intervals               |
| [sample date](https://github.com/careactivedocs/g3mr_kits/raw/main/sample_data.zip)              | Sample data                                                                             |

### Resampling Tool

RSSI-CVL per-second resampling will be performed by the cloud automatically when the data is post-processed. The cvl-resample tool in this repository is for reference in case the data needs to be verified. 

### Broadcasting Data Types

Once the Care Active Watch is paired, it begins broadcasting. There are four different types of broadcasting messages that can be identified by the Care Active Stations. The broadcasting continues incessantly, regardless of whether a station or mobile device is nearby. Most of the time, the watch is broadcasting the RTLS messages. Occasionally, the watch will switch to broadcasting other types of messages for brief periods, as outlined in the table below.

| Type          | Packets per Duration | Duration   | How Often        |
| :------------ | :------------------- | :--------- | :--------------- |
| RTLS/Active   | 2                    | 3 seconds  |                  |
| RTLS/Inactive | 1                    | 3 seconds  |                  |
| Pedometer     | 15                   | 15 seconds | Every 53 minutes |
| Pedo-Beacon   | 15                   | 15 seconds | Every 10 minutes |

#### RTLS

RTLS messages contains the G-Sensor readings. When there are motion activities, it is in active mode. If not, it is in inactive mode.

#### Pedometer

Pedometer type of messages are the pedometer readings for the station. This message also triggers the time-sync when power cycle or every 9 days.

#### Pedo-Beacon

Pedo-Beacon message is an Apple iBeacon format message that embeds the pedometer readings. If a nearby mobile device has the Care Active App installed with location permission granted, the pedo-beacon message can trigger the mobile app to report the GPS location to the Care Active Cloud. If the GPS feature is enabled, the GPS data will be logged in the cloud. Otherwise, this GPS information will be discarded.

### Local Time Conversion

If the CSV filename starting with "local_", it means the timestamp in the CSV file is already converted to local time. The timezone used for conversion is the timezone setting when you request the data on CA Insight. If the CSV filename does not start with "local_", it means the timestamp in the CSV file is in UTC format. You can use the [csv-local-time](../tools/csv_local_time/README.md) tool to convert UTC timestamp to local time. Two additinal CSV fields will be added to the converted CSV file:

| Field          | Description                                                             |
| :------------- | :---------------------------------------------------------------------- |
| local_time     | The local time corresponding to the UTC timestamp in the CSV file       |
| local_timezone | The original timezone setting when the data was requested on CA Insight |
