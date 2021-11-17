# Table: googlesheets_csv

Query data from Google Sheets. A table is automatically created to represent each
sheet mentioned in the config using `sheets` in a specific `spreadsheet_id`.

For instance, if you have configured [Google Sheets Plugin Examples](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4/edit#gid=0) spreadsheet, and `sheets` is configured with following sheets:

- Students
- Marks

This plugin will create 2 tables:

- Students
- Marks

Which you can then query directly:

```sql
select
  *
from
  "Students";
```

Each of these tables will have the same column structure as the Google Sheet they were
created from and all column values are returned as text data type.

**NOTE:**

- This table always checks for data in `A1` cell in the actual sheet to assume it as a CSV; otherwise plugin will skip that sheet.

## Examples

### Inspect the table structure

Assuming your connection is called `googlesheets` (the default), list all tables with:

```shell
.inspect googlesheets
+----------+-------------------------------+
| table    | description                   |
+----------+-------------------------------+
| Students | Retrieves data from Students. |
| Marks    | Retrieves data from Marks.    |
+----------+-------------------------------+
```

To get details for a specific table, inspect it by name:

```shell
.inspect "Students"
+--------------------------+------+-------------+
| column                   | type | description |
+--------------------------+------+-------------+
| Class Level              | text | Field 2.    |
| Extracurricular Activity | text | Field 5.    |
| Gender                   | text | Field 1.    |
| HOD                      | text | Field 6.    |
| Home State               | text | Field 3.    |
| Major                    | text | Field 4.    |
| Student Name             | text | Field 0.    |
+--------------------------+------+-------------+
```

### Query a sheet

Given the sheet `Students`, the query is:

```sql
select
  *
from
  "Students";
```

### Query a complex file name

Given the sheet `AWS Resource Types`, the query uses identifier quotes:

```sql
select
  *
from
  "AWS Resource Types";
```

### Query specific columns

Columns are always in text form when read from the Google sheet. The column names come from the first row of the sheet.

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
  "Student Name",
  "Major"
from
  "Students";
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  "Student Name" as student_name,
  cast("Percentage" as double precision) as percentage
from
  "Students"
where
  cast("Percentage" as double precision) > 85;
```

### Modifying header if have missing value

If the sheet has some missing values in the first row(header row), the table will add custom name based on column index, i.e. `1-A, 2-B, 27-AA`.

For instance, `Custom` sheet has missing values in column 2 and column 4. `.inspect` of that table will be following:

```shell
.inspect "Custom"

+-----------+------+-------------+
| column    | type | description |
+-----------+------+-------------+
| B         | text | Field 1.    |
| Col 1     | text | Field 0.    |
| Col 3     | text | Field 2.    |
| Col 3 [F] | text | Field 5.    |
| Col5      | text | Field 4.    |
| D         | text | Field 3.    |
+-----------+------+-------------+
```

### Modifying header if have columns with duplicate name

If the sheet has a column with duplicate name, the table will update the column name by adding the corresponding column index.

For instance, `Custom` sheet have columns with same name in column 3 and column 6. `.inspect` of that table will be following:

```shell
.inspect "Custom"

+-----------+------+-------------+
| column    | type | description |
+-----------+------+-------------+
| B         | text | Field 1.    |
| Col 1     | text | Field 0.    |
| Col 3     | text | Field 2.    |
| Col 3 [F] | text | Field 5.    |
| Col5      | text | Field 4.    |
| D         | text | Field 3.    |
+-----------+------+-------------+
```

### Querying a specific row using pivot table

```sql
with marks as (
  select
    "Student Name",
    "Exam",
    "Score"
  from
    "Marks"
  order by 1,2
),
pivot_marks as (
  select
    "Student Name" as student_name,
    max(case when "Exam" = '1' then "Score" else null end) as exam_1,
    max(case when "Exam" = '2' then "Score" else null end) as exam_2,
    max(case when "Exam" = '3' then "Score" else null end) as exam_3,
    max(case when "Exam" = '4' then "Score" else null end) as exam_4
  from
    marks
  group by
    "Student Name"
)
select
  *
from
  pivot_marks
where
  student_name = 'Bob';
```
