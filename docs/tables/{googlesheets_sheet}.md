# Table: {googlesheet_sheetname}

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

## Examples

### Inspect the table structure

Assuming your connection is called `googlesheets` (the default), list all tables with:

```sql
.inspect googlesheets
+-------------+----------------------------------+
| table       | description                      |
+-------------+----------------------------------+
| My Products | Retrieves data from My Products. |
| My Users    | Retrieves data from My Users.    |
+-------------+----------------------------------+
```

To get details for a specific table, inspect it by name:

```sql
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
  iage > 25
```
