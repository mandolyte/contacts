# Process Steps

*Unfortunately, I did not keep good notes on how to import TP data into Google Contacts. Sigh...*

## Step 1. Prepare Data

- Download the entire TP people spreadsheet.
    - Sign into TP
    - Select People, then [search builder] Saved
    - Select `NCCC Import Past and Present`
    - Download by clicking download, then "Export Excel" option "Standard"
    - Select option "Individual" (which should be the default)
- Load into Google Drive and convert to a Google Sheet.
- Download as a CSV file.
    - Rename to be like `2024-01-15_People.csv`
    - Copy CSV file into this folder.

## Step 2. Manipulate data for Google Contacts

In this step, just run `tp2gc.go`by running the script:

```
sh run.sh 2024-01-15
```

## Step 3. Import into Google Account

Using my alternate account, import into Google Contacts.