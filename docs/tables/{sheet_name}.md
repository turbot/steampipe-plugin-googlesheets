# Table: {sheet_name}

Query cell data from sheets in a Google Sheets spreadsheet in CSV format. A
table is automatically created to represent each sheet mentioned in the
configured `sheets`.

For instance, if using the following configuration:

```
connection "googlesheets" {
  plugin = "googlesheets"

  # Google Sheets Plugin - Sample School Data
  spreadsheet_id = "11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4"
  sheets         = ["Dashboard", "Students", "Books", "Marks", "Employees"]

  token_path = "/Users/myuser/gcloud/application_default_credentials.json"
}
```

This plugin will create 4 tables:

- Books
- Employees
- Marks
- Students

For the tables that were created, these can be queried directly:

```sql
select
  *
from
  "Students";
```

Each of these tables will have the same column structure as the Google Sheet
they were created from and all column values are returned as text data type.

Note: A table is not created for the `Dashboard` sheet as it does not have any
data in cell `A1`. For more information on how tables are created, please see [Table Restrictions and Notes](#table-restrictions-and-notes).

All examples below can be used with the [Google Sheets Plugin - Sample School
Data](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4)
spreadsheet, which is a public spreadsheet maintained by the Steampipe team.

## Examples

### Inspect the table structure

Assuming your connection is called `googlesheets` (the default), list all tables with:

```shell
.inspect googlesheets
+--------------------------+------------------------------------------------------------+
| table                    | description                                                |
+--------------------------+------------------------------------------------------------+
| Books                    | Retrieves data from Books.                                 |
| Employees                | Retrieves data from Employees.                             |
| Marks                    | Retrieves data from Marks.                                 |
| Students                 | Retrieves data from Students.                              |
| googlesheets_cell        | Retrieve information of cells of a sheet in a spreadsheet. |
| googlesheets_sheet       | Retrieve the sheet in a given spreadsheet.                 |
| googlesheets_spreadsheet | Retrieve the metadata of given spreadsheet.                |
+--------------------------+------------------------------------------------------------+
```

To get details for a specific table, inspect it by name:

```shell
.inspect "Students"
+--------------------------+------+-------------+
| column                   | type | description |
+--------------------------+------+-------------+
| Class Level              | text | Field 2.    |
| Extracurricular Activity | text | Field 5.    |
| GPA                      | text | Field 7.    |
| Home State               | text | Field 3.    |
| ID                       | text | Field 1.    |
| Major                    | text | Field 4.    |
| Mentor                   | text | Field 6.    |
| Student Name             | text | Field 0.    |
+--------------------------+------+-------------+
```

### Query a sheet

Given the sheet `Students`, the query requires identifier quotes:

```sql
select
  *
from
  "Students";
```

### Query specific columns

Columns are always in text form when read from Google Sheets. The column names
come from the first row of the sheet.

If your column names are complex, use identifier quotes:

```sql
select
  "Student Name",
  "Major"
from
  "Students";
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  "Name" as book_name,
  "Author" as author,
  "Issued By" as issued_by,
  case
    when "Issue Date" <> '' then "Issue Date"::timestamptz
  end as issued_at,
  "Verified"::boolean as verified
from
  "Books";
```

## Table Restrictions and Notes

- CSV tables will only be created for sheets that have data in cell `A1`.
- If a sheet's header row is missing some values, the table will use the column index for the column name.
- If a sheet's header row has more than one column with same name, column indexes will be appended onto the end of duplicate columns.
- If a sheet's header row has vertically merged cells, the table will use the merged cell's value for all affected cells and apply duplicate protection.

For instance, the `Employees` table has the following header row:

| A           | B             | C             | D       | E       | F + G    | H | I            | J             |
|-------------|---------------|---------------|---------|---------|----------|---|--------------|---------------|
| Employee ID | Employee Name | Profile Image | Contact | Contact | Birthday |   | Joining Date | Days Employed |

Running `.inspect "Employees"` then results in:

```shell
.inspect "Employees"

+---------------+------+-------------+
| column        | type | description |
+---------------+------+-------------+
| Birthday [F]  | text | Field 5.    |
| Birthday [G]  | text | Field 6.    |
| Contact       | text | Field 3.    |
| Contact [E]   | text | Field 4.    |
| Days Employed | text | Field 9.    |
| Employee ID   | text | Field 0.    |
| Employee Name | text | Field 1.    |
| H             | text | Field 7.    |
| Joining Date  | text | Field 8.    |
| Profile Image | text | Field 2.    |
+---------------+------+-------------+
```
