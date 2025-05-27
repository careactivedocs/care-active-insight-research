# Device Name Change Logs

Device name change logs keep tracking the device name changes in case the family account user change the name of the device and the researcher is not aware of it.

## File Location

`exported_data/{collector}/{family_account}/logs`

## Filename Format

{MAC}.log

## File Content

The time when this name was used. Device name came from the collected source data and processed daily.

| column | description                             |
| :----- | :-------------------------------------- |
| 1      | UTC time when the names are logged      |
| 2      | Device name on the app (user changable) |
| 3      | Device alias name on CA Insight         |
| 4      | <first log> or <post processing time>   |

```<post processing time>``` means the time was the time from the post processed report. This usually happens when the device has been removed from the account.

## Sample Content

```text
2025-02-18T23:01:10Z	g3mr-samson samson
2025-02-18T14:35:45Z	g3mr-samson samson
2025-02-17T22:28:24Z	g3mr-TPU    samson
2025-02-17T23:00:19Z	g3mr-TPU    samson
2025-02-17T05:44:30Z	g3mr-TPU    samson
2025-02-16T22:59:55Z	g3mr-TPU    samson
2025-03-10T14:56:44Z	g3mr-TPU	samson	<first log>
```
