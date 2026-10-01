# Data Processing Time

This document explains when newly collected data becomes available, and the difference between the near-real-time "Quick Access" preview and the authoritative "Final" data.

## Two Stages of Data

| Stage       | Covers                                | Refresh            | Use For                          |
| :---------- | :------------------------------------ | :------------------ | :-------------------------------- |
| Quick Access | The current, still-in-progress UTC day | Every 4 hours        | A near-real-time preview          |
| Final        | The previous, fully completed UTC day  | Once a day           | The authoritative record          |

## Why Data Cannot Be Finalized Immediately

Care Active Watches, Stations, and the mobile app queue and forward activity, location, and sensor data to the cloud. Because of this queueing design, packets for a given moment do not necessarily arrive at the cloud in that same moment — some arrive late and need to be resequenced. To produce accurate, consistent reports, the cloud must wait until an entire UTC day has completely elapsed before it can process that day as final.

## Final Data

- A UTC calendar day (00:00-23:59 UTC) is only finalized after it has fully ended.
- The final processing run happens roughly 13 hours into the following UTC day (around 13:20 UTC).
  - Example: data collected on 2026-09-30 (UTC) is finalized and available around 2026-10-01, 13:20 UTC.
- This run produces the authoritative version of that day's data, including:
  - Daily RTLS-Activity reports (RSSI, CVL, resampled CVL, XYZ) — see [Daily RTLS-Activity Report](daily-rtls-activity-report.md)
  - Daily post-processed CSV reports (location, pedometer, motion) — see [Daily Post-Processed Data](daily-pp-csv-report.md)
  - KML location files (the only location-related output not already available via Quick Access — see below)
  - [Maintenance reports](maintenance-reports.md) (station and device status)
  - Billing data
- Once generated, final data will not change unless Care Active support performs a manual reprocessing.

## Quick Access Data

- While a UTC day is still in progress, a preview of that day's data is generated automatically every 4 hours, at approximately 00:10, 04:10, 08:10, 12:10, 16:10, and 20:10 UTC.
- Quick Access covers the RTLS-Activity reports (RSSI, CVL, resampled CVL, XYZ) and the post-processed location, pedometer, and motion CSV reports, all for the current, in-progress day.
  - This includes the location CSV report (GPS coordinates and the other location fields), but **not** the KML location file — see the note below.
- Because the day has not ended, Quick Access data is necessarily incomplete: hours that have not yet happened are not represented, and figures from earlier in the day may still shift slightly as more data continues to arrive and is resequenced.
- Quick Access is meant to give researchers a near-real-time look at recent activity. It is not a final record. A few outputs are **not** part of the Quick Access cycle and are only produced by the once-daily Final run:
  - KML location files (map-viewing format)
  - Maintenance reports (station and device status)
  - Billing data
- The following day's Final run fully supersedes any Quick Access snapshots for that day.

## Schedule Summary (UTC)

| Time (UTC)                                        | Event                                                          |
| :-------------------------------------------------- | :--------------------------------------------------------------- |
| 00:10, 04:10, 08:10, 12:10, 16:10, 20:10             | Quick Access preview refreshed for the current, in-progress day |
| ~13:20 (the next day)                                | Final data for the previous (now complete) day is generated     |

## Revision History

| Document Revision | Revision Date | Description     | Note |
| :----------------: | :-----------: | ---------------- | ---- |
|        1.0         |  2026-10-01   | Initial version  |      |
|        1.1         |  2026-10-01   | Clarified that Quick Access includes location/pedo/motion CSV reports, not just RTLS-Activity; only KML, maintenance, and billing are Final-only |      |
