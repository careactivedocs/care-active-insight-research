# Daily Source Data

Daily source data is the data tat is post-processed on. Most information has been converted in the post-processed data folders. This source data is provided in case some further information or data validation is needed.

## File Locations

| Type                          | data_type | File Location                                                                                                               |
| :---------------------------- | :-------- | :-------------------------------------------------------------------------------------------------------------------------- |
| [RTLS](#rtls)                 | 1         | exported_data_v2/{collector}/source/{family_account}/rtls/{MAC}/{YYYYMM}/{YYYY-MM-DD}_{MAC}_{device_name}_rtls_json         |
| [Pedometer](#pedometer)       | 2         | exported_data_v2/{collector}/source/{family_account}/pedo/{MAC}/{YYYYMM}/{YYYY-MM-DD}_{MAC}_{device_name}_pedo_json         |
| [Location-GPS](#location-gps) | 3         | exported_data_v2/{collector}/source/{family_account}/location/{MAC}/{YYYYMM}/{YYYY-MM-DD}_{MAC}_{device_name}_location_json |
| [Motion](#motion)             | 4         | exported_data_v2/{collector}/source/{family_account}/motion/{MAC}/{YYYYMM}/{YYYY-MM-DD}_{MAC}_{device_name}_motion_json     |

## RTLS

RTLS data is a combination of watch G-Sensor readings, and station RSSI.

| field          | descriptions                                                                                            |
| :------------- | :------------------------------------------------------------------------------------------------------ |
| id             | The unique id of the record. RTLS is not stored in DB , this data can be ignored.                       |
| family_id      | The family account email that this record belongs to.                                                   |
| collector_id   | The collector that this record belongs to.                                                              |
| device_id      | The serial number of this Care Active Watch.                                                            |
| device_name    | The name of this watch on the app.                                                                      |
| device_mac     | The MAC address of this watch.                                                                          |
| scanned_at_ms  | Time that this RTLS packet was scanned. If there are multiple sources, the lowest number is used.       |
| scanned_at_utc | Human format of scanned_at_ms.                                                                          |
| ttl            | This time-to-live field can be ignored in RTLS data type.                                               |
| data_type      | 1 for RTLS.                                                                                             |
| created_at     | The EPOCH-ms of when this record was created in the processing cache (cloud time of data consolidation) |

### RTLS Payload

Data/Information of this RTLS packet.

| field             | descriptions                                                                                    |
| :---------------- | :---------------------------------------------------------------------------------------------- |
| epoch_byte        | This is time_label in [Care Active Watch Advertising](./README.md#document-of-processing-tools) |
| event_seq         | Sequence ID of this RTLS packet                                                                 |
| motion            | Sampled G-Sensor data.                                                                          |
| ref_scanned_at_ms | The "scanned_at_ms" of the first received one by the cloud.                                     |
| ref_station_id    | The "station_id" of the first received one by the cloud.                                        |

In the motion field, if there is only two sets of motion data, it means the watch is in inactive mode. If in active mode, there are more than two sets of motion data. The last set of motion data is ACVL.

#### station_receive_param

All the station information that scanned this RTLS packet.

| field            | descriptions                                                                                 |
| :--------------- | :------------------------------------------------------------------------------------------- |
| create_at_ms     | Time of the cloud when this packet was received by the cloud. (cloud time)                   |
| name             | The station name of this received station.                                                   |
| raw_sample_at_ms | This value was converted from scanned_at_ms of the Station with the epoch_byte of the watch. |
| room             | Room type of the station                                                                     |
| rssi             | The RSSI value of this scan                                                                  |
| scanned_at_ms    | Time of the Station when scanned this RTLS packet                                            |
| station_id       | The serial number of the station that scanned this RTLS packet.                              |

It is possible that the same station receives the same RTLS packet multiple times. If that happens, there will be more station records with the same station_id in the field of station_receive_param.

#### raw_sample_at_ms

This field is obsolete in Data Format Version 2. Ignore this one.

### Sample RTLS Data

```json
{
    "id": "74322ad5-29ba-4b16-984a-edd31305f937",
    "family_id": "qa+demo@mytracmo.com",
    "collector_id": "qa0000000001",
    "device_id": "1D6C62CA",
    "device_name": "Band_7347",
    "device_mac": "80:6F:B0:7B:73:47",
    "scanned_at_ms": 1726095672378,
    "scanned_at_utc": "2024-09-11T23:01:12Z",
    "ttl": 1726268792,
    "data_type": 1,
    "payload": {
        "epoch_byte": 34,
        "event_seq": 29256,
        "motion": [
            {
                "VectX": 0,
                "VectY": 1,
                "VectZ": 64
            }
        ],
        "ref_scanned_at_ms": 1726095672387,
        "ref_station_id": "T4-A4CF129EC02C",
        "station_receive_param": [
            {
                "created_at_ms": 1726095673357,
                "name": "Qos1-10s",
                "raw_sample_at_ms": 1726095650000,
                "room": {
                    "room_type": 6
                },
                "rssi": -55,
                "scanned_at_ms": 1726095672387,
                "station_id": "T4-A4CF129EC02C"
            },
            {
                "created_at_ms": 1726095675322,
                "name": "DBclean",
                "raw_sample_at_ms": 1726095650000,
                "room": {
                    "room_type": 8
                },
                "rssi": -46,
                "scanned_at_ms": 1726095672378,
                "station_id": "T4-A4CF129EC14C"
            },
            {
                "created_at_ms": 1726095675502,
                "name": "DBclean",
                "raw_sample_at_ms": 1726095650000,
                "room": {
                    "room_type": 8
                },
                "rssi": -47,
                "scanned_at_ms": 1726095673508,
                "station_id": "T4-A4CF129EC14C"
            }
        ]
    },
    "created_at": 1726095992436
}
```

## Pedometer

There are two sources for pedometer data: one is the station, and the other is the mobile device. The "sender_id" field indicates the source. If the source is a mobile device, the pedometer data originated from the Pedo-Beacon. In the Pedo-Beacon, the pedometer reading is only the latest one. If the source is a station, the pedometer readings include the latest reading for today, yesterday, and the day before yesterday.

| field          | descriptions                                                               |
| :------------- | :------------------------------------------------------------------------- |
| id             | The unique id of the record in the activity_data database.                 |
| family_id      | The family account email that this record belongs to.                      |
| collector_id   | The collector that this record belongs to.                                 |
| device_id      | The serial number of this Care Active Watch.                               |
| device_name    | The name of this watch on the app.                                         |
| device_mac     | The MAC address of this watch.                                             |
| scanned_at_ms  | Time that this pedometer packet was scanned.                               |
| scanned_at_utc | Human format of scanned_at_ms.                                             |
| ttl            | time-to-live in activity_data database                                     |
| data_type      | 2 for pedometer                                                            |
| created_at     | The EPOCH-ms of when this record was created in the activity_data database |

### Pedo Payload

Data/Information of this pedometer packet.

| field      | descriptions                                                                                                                             |
| :--------- | :--------------------------------------------------------------------------------------------------------------------------------------- |
| pedometer  | The pedometer readings of today's, yesterday's, and the day before yesterday. If the source is mobile, there is only today's.            |
| sender_id  | The serial number of the sender station or the mobile unique id.                                                                         |
| t_mark_day | Days since January 1st on the watch. This could be wrong if the time of the watch needs to be adjusted.                                  |
| t_mark_rnd | A random number from the watch. This random number is updated per power-cycle. This helps to tell whether there is a power cycle or not. |

### Sample Pedo Data

```json
[
    {
        "id": "74dbfc25-f4a6-44a8-bd1e-2840adc2543e",
        "family_id": "samson@qblinks.com",
        "collector_id": "qa0000000001",
        "device_id": "4C126812",
        "device_name": "G3MR-Samson",
        "device_mac": "80:6F:B0:7B:93:0E",
        "scanned_at_ms": 1726073822000,
        "scanned_at_utc": "2024-09-11T16:57:02Z",
        "ttl": 1726246624,
        "data_type": 2,
        "payload": {
            "pedometer": [
                129,
                4808,
                246
            ],
            "sender_id": "T4-4827E2C70984",
            "t_mark_day": 256,
            "t_mark_rnd": 117
        },
        "created_at": 1726073824091
    },
    {
        "id": "9d217d9e-c2f9-43b1-8b2e-8806b61c0946",
        "family_id": "samson@qblinks.com",
        "collector_id": "qa0000000001",
        "device_id": "4C126812",
        "device_name": "G3MR-Samson",
        "device_mac": "80:6F:B0:7B:93:0E",
        "scanned_at_ms": 1726074011991,
        "scanned_at_utc": "2024-09-11T17:00:11Z",
        "ttl": 1726246813,
        "data_type": 2,
        "payload": {
            "pedometer": [
                254
            ],
            "sender_id": "71B1657E-1881-47C6-B097-D0E044742C5B",
            "t_mark_day": 256,
            "t_mark_rnd": 117
        },
        "created_at": 1726074013497
    }
]
```

## Location GPS

The location data is available only when the following conditions apply.

- Collector GPS feature is enabled.
- Care Active App is installed in the participant's mobile device.
- The app is granted with the location permission.
- The app is logged-in with the participant's family account.
- The above mobile device is nearby the Care Active Watch.
- There is no Care Active Station nearby. This is to avoid the disclosure of the location information of any station installation.

If all the above conditions are met, the GPS location is uploaded to the cloud every 10 minutes.

| field          | descriptions                                                              |
| :------------- | :------------------------------------------------------------------------ |
| id             | The unique id of the record in the activity_data database.                |
| family_id      | The family account email that this record belongs to.                     |
| collector_id   | The collector that this record belongs to.                                |
| device_id      | The serial number of this Care Active Watch.                              |
| device_name    | The name of this watch on the app.                                        |
| device_mac     | The MAC address of this watch.                                            |
| scanned_at_ms  | Time that this pedometer packet was scanned.                              |
| scanned_at_utc | Human format of scanned_at_ms.                                            |
| ttl            | time-to-live in activity_data database                                    |
| data_type      | 3 for location                                                            |
| created_at     | The EPOCH-ms of when this record was createdin the activity_data database |

### Location Payload

Data/ GPS Information of this location packet.

| field            | descriptions                                                                          |
| :--------------- | :------------------------------------------------------------------------------------ |
| created_at       | The EPOCH-ms of when the record was created in the location_log database (cloud time) |
| created_at_ttl   | Time-to-live in the location_log database                                             |
| device_name      | The name of this watch in the location_log database                                   |
| device_photo     | The image of the this watch in the location_log database                              |
| gps_accuracy     | GPS accuracy reported from the mobile                                                 |
| gps_latitude     | GPS latitude                                                                          |
| gps_longitude    | GPS longitude                                                                         |
| id               | The unique id of the record in the location_log database                              |
| reason_code      | Reason Code reference in the location_log database                                    |
| reason_data      | Associate reason data in the location_log database                                    |
| scaned_at        | The EPOCH-ms of when the Pedo-Beacon was received on the mobile (mobile time)         |
| sender_device_id | The mobile unique id                                                                  |
| target_device_id | The serial number of the Watch                                                        |
| user_id          | Cognito unique user id in the location_log database                                   |

### Sample Location Data

```json
[
    {
        "id": "584d0a27-e181-4f68-bf60-e0cd5b001884",
        "family_id": "samson@qblinks.com",
        "collector_id": "qa0000000001",
        "device_id": "4C126812",
        "device_name": "G3MR-Samson",
        "device_mac": "80:6F:B0:7B:93:0E",
        "scanned_at_ms": 1726042576755,
        "scanned_at_utc": "2024-09-11T08:16:16Z",
        "ttl": 1726215378,
        "data_type": 3,
        "payload": {
            "created_at": 1726042578256,
            "created_at_ttl": 1731226578,
            "device_name": "G3MR-Samson",
            "device_photo": "https://portal.careactive.ai/img/sha256/2c5184c94be6f0137497cd788cc001ee14b9cafa7b819d2e98c6a058fca5cd5d.png",
            "gps_accuracy": 35,
            "gps_latitude": "23.610442333808496",
            "gps_longitude": "121.5287793121825",
            "id": "689a2cea-fa0d-40f0-8ba0-aa71dc07a73a",
            "reason_code": 11,
            "reason_data": {
                "phone_name": "O2LAND"
            },
            "scaned_at": 1726042576755,
            "sender_device_id": "71B1657E-1881-47C6-B097-D0E044742C5B",
            "target_device_id": "4C126812",
            "user_id": "samson@qblinks.com:tracmo"
        },
        "created_at": 1726042578272
    }
]
```

## Motion

Motion data was from the Care Active Motion sensor. Only one record of the motion event will be stored. If there are multiple stations that detect the same motion event, they will all send the event to the cloud. But the cloud will keep only one of them.

| field          | descriptions                                                               |
| :------------- | :------------------------------------------------------------------------- |
| id             | The unique id of the record in the activity_data database.                 |
| family_id      | The family account email that this record belongs to.                      |
| collector_id   | The collector that this record belongs to.                                 |
| device_id      | The serial number of this Care Active Motion Sensor.                       |
| device_name    | The name of this motion sensor on the app.                                 |
| device_mac     | The MAC address of this motion sensor.                                     |
| scanned_at_ms  | Time that this motion activity was detected.                               |
| scanned_at_utc | Human format of scanned_at_ms.                                             |
| ttl            | time-to-live in activity_data database                                     |
| data_type      | 4 for motion                                                               |
| created_at     | The EPOCH-ms of when this record was created in the activity_data database |

### Motion Payload

Data/ GPS Information of motion activity packet.

| field          | descriptions                              |
| :------------- | :---------------------------------------- |
| battery        | Battery level of the motion sensor        |
| motion_cvl     | Motion level of this activity event       |
| motion_cvl_max | Max motion level of the recent 30 seconds |
| room           | Room type                                 |
| rssi           | RSSI of this motion event packet          |
| station_id     | The serial number of the sender station   |

#### motion_cvl

CVL is Collected Vector Length to measure the motion intensity. Please refer to [CVL Equation](./post-process-reports.md#cvl-data-report) for more details.

### Sample Motion Data

```json
[    {
        "id": "1a798851-a8b1-4fed-b8a8-1ff09c4b3aeb",
        "family_id": "samson@qblinks.com",
        "collector_id": "qa0000000001",
        "device_id": "D922A666",
        "device_name": "Mom’s TV Remote",
        "device_mac": "C4:64:E3:A9:AB:A4",
        "scanned_at_ms": 1726064141000,
        "scanned_at_utc": "2024-09-11T14:15:41Z",
        "ttl": 1726236941,
        "data_type": 4,
        "payload": {
            "battery": 100,
            "motion_cvl": 22,
            "motion_cvl_max": 22,
            "room": {
                "room_type": 7
            },
            "rssi": -79,
            "station_id": "T4-2462ABC18B94"
        },
        "created_at": 1726064141651
    }
]
```

## Revision History

| Document Revision | Revision Date | Description           | Note |
| :---------------: | :-----------: | --------------------- | ---- |
|        1.0        |  2024-09-12   | Init Version          |      |
|        2.0        |  2025-03-11   | Data Format Version 2 |      |
