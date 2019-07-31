iRail to CSV
============
iRail to CSV is a server that wraps iRail API calls to reply in CSV format.
CSV is chosen as it is the lightest way to send information, this server is meant to serve ultra low memory environments (2019 standards).

## What it is not
This API does not 1-to-1 map all API data to CSV, it will only map the essentials and will use text for complex data (eg connections).
