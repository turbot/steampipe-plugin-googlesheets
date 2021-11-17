# Table: googlesheets_sheet

Retrieves information of a sheet in the given spreadsheet.

## Examples

### Basic info

```sql
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet;
```

### Get information about a specific sheet

```sql
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet
where
  title = 'Students';
```

### List all hidden sheets

```sql
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet
where
  hidden;
```

### List sheets with merge cells

```sql
select
  title,
  sheet_id,
  spreadsheet_id,
  merges
from
  googlesheets_sheet;
```
