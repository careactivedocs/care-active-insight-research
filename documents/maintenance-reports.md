# Maintenance Reports

The maintenance reports are generated daily at 00:35 UTC. This report helps identify the status of all registered devices, including their connectivity and battery levels. This single report contains information about all the devices under the same collector ID for easy reference. There are two types of reports, station report and sensor device report.

## File Location

| type                 | File Location                                                               |
| :------------------- | :-------------------------------------------------------------------------- |
| Station Report       | exported_data/{collector}/maintenance/{YYYYMM}/{YYYY}-{MM}-{DD}_station.csv |
| Sensor Device Report | exported_data/{collector}/maintenance/{YYYYMM}/{YYYY}-{MM}-{DD}_device.csv  |

## Station Daily Report

File name: `{date}_station.csv`

This report aims to furnish essential information about each station. Researchers can know the basic information of each station and evaluate the reliability of their network connectivity.


| Field               | Note                                                                                           |
| :------------------ | :--------------------------------------------------------------------------------------------- |
| family_id           | The family account this station was registered on                                              |
| family_name         | Name of the family on CA Insight                                                               |
| station_mac         | The MAC address of the station                                                                 |
| bind_time           | UTC time that station was added to this family account                                         |
| station_name        | The name of the station on the app                                                             |
| station_name_alias  | The name of the station on CA Insight                                                          |
| device_model        | device model of the station. EX: T4 or T5                                                      |
| room_type           | [room type](#room-type) of the station                                                         |
| wifi_rssi           | wifi signal rssi of the station (last record of the day)                                       |
| status              | online or offline. If the station is offline, the `last_connected_time` needs to be indicated. |
| last_connected_time | the UTC time of the last connection                                                            |
| alive_counts        | The counts of the alive etm events in the past day. alive etm is reported every 7.5 minutes.   |
| disconnect_counts   | The counts of the lwt etm events in the past day (count of disconnections)                     |

### Room Type

| Code | Type of Room    |
| :--- | :-------------- |
| 0    | Unknown         |
| 1    | Main Bedroom    |
| 2    | Bedroom 1       |
| 3    | Bedroom 2       |
| 4    | Living Room     |
| 5    | Family Room     |
| 6    | Recreation Room |
| 7    | Kitchen         |
| 8    | Dining Room     |
| 9    | Den/Office      |
| 10   | Master Bathroom |
| 11   | Bathroom        |
| 12   | Garage          |
| 13   | Patio           |
| 14   | Entryway        |
| 15   | Other           |

### Sample Station Report

```text
family_id,family_name,station_mac,bind_time,station_name,device_model,room_type,wifi_rssi,status,last_connected_time,alive_counts,disconnect_counts
F1C96654979,Samson Chen,74:4D:BD:2A:4C:34,2025-02-06T05:33:05Z,1F Guest Room,t5,2,-65,online,2025-03-10T15:10:08Z,192,0
```

## Sensor Device Daily Report

File name: `{date}_device.csv`

The second report focuses on checking the device's battery capacity and the volume of data it generates. Researchers can use this report to track the battery levels of devices and monitor the quantity of data produced.

| Field              | Note                                                                                             |
| :----------------- | :----------------------------------------------------------------------------------------------- |
| family_id          | The family account this device was registered on                                                 |
| family_name        | Name of the family on CA Insight                                                                 |
| collection         | "OFF" if this family account's data collection is disabled                                       |
| device_mac         | The MAC address of the device                                                                    |
| bind_time          | UTC time that the device was added to this family account                                        |
| last_activity_time | Detected last activity time of the device                                                        |
| device_name        | The name of the device on the app                                                                |
| device_name_alias  | The name of the device on CA Insight                                                             |
| device_model       | device model of the device. EX: g3 or g3mr                                                       |
| battery            | battery level of the device                                                                      |
| data_type          | [Data type](#data-types) of the collected data (the same device may have multiple types of data) |
| total_records      | total number of collected data records of certain data type                                      |
| active_counts      | total number of the active rtls data entries (acvl was changed)                                  |
| inactive_counts    | total number of the inactive rtls data entries                                                   |

### Data Types

| Type       | Description                                                   |
| :--------- | :------------------------------------------------------------ |
| location   | GPS coordinates if the GPS feature is enabled for this device |
| motion     | Motion sensor                                                 |
| pedo       | Pedometer readings of Care Watch                              |
| rtls       | Activity and indoor location source data                      |
| undetected | Refer to [Undetected Data Type](#undetected-data-type)        |

#### Undetected Data Type

If there is a device registered on the family account but no data generated from this device for the entire day, the data type will be marked as "undetected". There are a few reasons that can cause this situation:

- There is no nearby online station around the registered device.
- The registered device runs out of battery.
- Data collection is disabled for this family account.
- The family account of this device has been moved to another property.

### Sample Sensor Device Report

```text
family_id,family_name,collection,device_mac,bind_time,last_activity_time,device_name,device_model,battery,data_type,total_records,active_counts,inactive_counts
F1C84864183,Samson Chen,,80:6F:B0:7B:93:0E,2025-02-06T09:47:06Z,2025-04-18T04:00:41Z,g3mr-samson:TPU,g3mr,15,rtls,21787,7064,14723
```

## Revision History

| Document Revision | Revision Date | Description                                                  | Note        |
| :---------------: | :-----------: | ------------------------------------------------------------ | ----------- |
|        0.1        |  2023-08-03   | initial version                                              | first draft |
|        0.2        |  2023-09-25   | add monthly report                                           |             |
|        1.0        |  2023-10-12   | add the storage path of the files                            |             |
|        1.1        |  2023-11-01   | add new fields                                               |             |
|        1.2        |  2024-03-06   | add first_assigned_at field into monthly reports for billing |             |
|        1.3        |  2024-06-12   | fixed monthly reports to match the order of CSV fields       |             |
|        2.0        |  2025-03-11   | version 2 post-processing tool                               |             |

