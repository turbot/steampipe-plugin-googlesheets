# Table: googlesheets_csv

Query data from Google Sheets. A table is automatically created to represent each
sheet mentioned in the config using `sheets` in a specific `spreadsheet_id`.

For instance, if you have configured [Google Sheets Plugin Examples](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4/edit#gid=0) spreadsheet, and `sheets` is configured with following sheets:

- Books
- Dashboard
- Employee
- Marks
- Students

This plugin will create 2 tables:

- Books
- Employee
- Marks
- Students

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

- If the sheet is missing some values in the first row(header row), the table will use the column name based on column index, i.e. `A, B, AA`.

For instance, `Employee` sheet has missing values in column 6. `.inspect` of that table will be following:

```shell
.inspect "Employee"

+---------------+------+-------------+
| column        | type | description |
+---------------+------+-------------+
| Birthday      | text | Field 8.    |
| Contact       | text | Field 3.    |
| Contact [E]   | text | Field 4.    |
| Duration      | text | Field 6.    |
| Employee ID   | text | Field 0.    |
| Employee Name | text | Field 1.    |
| H             | text | Field 7.    |
| Joining Date  | text | Field 5.    |
| Profile Image | text | Field 2.    |
+---------------+------+-------------+
```

- If the sheet has more than one columns with same name, the table will update the column name by adding the corresponding column index.

For instance, `Employee` sheet have columns with same name in column 4 and column 5. `.inspect` of that table will be following:

```shell
.inspect "Employee"

+---------------+------+-------------+
| column        | type | description |
+---------------+------+-------------+
| Birthday      | text | Field 8.    |
| Contact       | text | Field 3.    |
| Contact [E]   | text | Field 4.    |
| Duration      | text | Field 6.    |
| Employee ID   | text | Field 0.    |
| Employee Name | text | Field 1.    |
| H             | text | Field 7.    |
| Joining Date  | text | Field 5.    |
| Profile Image | text | Field 2.    |
+---------------+------+-------------+
```

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

Given the sheet `Students`, the query uses identifier quotes:

```sql
select
  *
from
  "Students";
```

### Query specific columns

Columns are always in text form when read from the Google sheet. The column names come from the first row of the sheet.

**NOTE:**

- If your column names are complex, use identifier quotes:

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
