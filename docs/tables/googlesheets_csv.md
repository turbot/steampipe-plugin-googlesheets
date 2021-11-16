# Table: googlesheets_csv

Query data from Google Sheets. A table is automatically created to represent each
sheet mentioned in the config using `sheets` in a specific `spreadsheet_id`.

For instance, if `sheets` is configured with following sheets:

- My Products
- My Users

This plugin will create 2 tables:

- My Products
- My Users

Which you can then query directly:

```sql
select
  *
from
  "My Users";
```

Each of these tables will have the same column structure as the Google Sheet they were
created from and all column values are returned as text data type.

If the actual sheet is missing some values in row 1 (header row), the table will assume `A, B, AA` as column name.

Check out [Steampipe for Google Sheet Example](https://docs.google.com/spreadsheets/d/13AdC3MKVg2zQj2OJAzQ3NN6HqaPuYhN3Aagr7OmRn9c/edit#gid=354856876) spreadsheet for more examples.

**NOTE:**

- This table always checks for data in `A1` cell in the actual sheet to assume it as a CSV; otherwise plugin will skip that sheet.

## Examples

### Inspect the table structure

Assuming your connection is called `googlesheets` (the default), list all tables with:

```shell
.inspect googlesheets
+-------------+----------------------------------+
| table       | description                      |
+-------------+----------------------------------+
| My Products | Retrieves data from My Products. |
| My Users    | Retrieves data from My Users.    |
+-------------+----------------------------------+
```

To get details for a specific table, inspect it by name:

```shell
.inspect "My Users"
+------------+------+-------------+
| column     | type | description |
+------------+------+-------------+
| email      | text | Field 2.    |
| first_name | text | Field 0.    |
| last_name  | text | Field 1.    |
+------------+------+-------------+
```

### Query a sheet

Given the sheet `users`, the query is:

```sql
select
  *
from
  users;
```

### Query a complex file name

Given the sheet `My complex sheet-name`, the query uses identifier quotes:

```sql
select
  *
from
  "My complex sheet-name";
```

### Query specific columns

Columns are always in text form when read from the Google sheet. The column names come from the first row of the file.

```sql
select
  first_name,
  last_name
from
  "My Users";
```

If your column names are complex, use identifier quotes:

```sql
select
  "First Name",
  "Last Name"
from
  "My Users";
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  first_name,
  age::int as iage
from
  "My Users"
where
  iage > 25;
```
