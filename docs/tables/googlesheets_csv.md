# Table: googlesheets_csv

Query data from Google Sheets. A table is automatically created to represent each
sheet mentioned in the config using `sheets` in a specific `spreadsheet_id`.

For instance, if you have configured [Google Sheets Plugin Examples](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4/edit#gid=0) spreadsheet, and `sheets` is configured with following sheets:

- Books
- Dashboard
- Employees
- Marks
- Students

This plugin will create 2 tables:

- Books
- Employees
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

- This table always checks for data in `A1` cell in the actual sheet to assume it as a CSV; otherwise plugin will **skip** that sheet.

- If the sheet is missing some values in the first row(header row), the table will use the column name based on column index, i.e. `A, B, AA`.

  For instance, `Employees` sheet has missing values in A8. `.inspect` of that table will be following:

  ```shell
  .inspect "Employees"

  +---------------+------+-------------+
  | column        | type | description |
  +---------------+------+-------------+
  | Birthday [F]  | text | Field 5.    |
  | Birthday [G]  | text | Field 6.    |
  | Contact       | text | Field 3.    |
  | Contact [E]   | text | Field 4.    |
  | Duration      | text | Field 8.    |
  | Employee ID   | text | Field 0.    |
  | Employee Name | text | Field 1.    |
  | J             | text | Field 9.    |
  | Joining Date  | text | Field 7.    |
  | Profile Image | text | Field 2.    |
  +---------------+------+-------------+
  ```

- If the sheet has more than one columns with same name, the table will update the column name by adding the corresponding column index.

  For instance, `Employees` sheet have columns with same name in D1 and E1. `.inspect` of that table will be following:

  ```shell
  .inspect "Employees"

  +---------------+------+-------------+
  | column        | type | description |
  +---------------+------+-------------+
  | Birthday [F]  | text | Field 5.    |
  | Birthday [G]  | text | Field 6.    |
  | Contact       | text | Field 3.    |
  | Contact [E]   | text | Field 4.    |
  | Duration      | text | Field 8.    |
  | Employee ID   | text | Field 0.    |
  | Employee Name | text | Field 1.    |
  | J             | text | Field 9.    |
  | Joining Date  | text | Field 7.    |
  | Profile Image | text | Field 2.    |
  +---------------+------+-------------+
  ```

- If the sheet has vertically merged cells, the table will update the column name by adding the corresponding column index along with the cell data to differentiate the columns.

  For instance, `Employees` sheet has header where F1 and G1 are merged together. `.inspect` of that table will be following:

  ```shell
  .inspect "Employees"

  +---------------+------+-------------+
  | column        | type | description |
  +---------------+------+-------------+
  | Birthday [F]  | text | Field 5.    |
  | Birthday [G]  | text | Field 6.    |
  | Contact       | text | Field 3.    |
  | Contact [E]   | text | Field 4.    |
  | Duration      | text | Field 8.    |
  | Employee ID   | text | Field 0.    |
  | Employee Name | text | Field 1.    |
  | J             | text | Field 9.    |
  | Joining Date  | text | Field 7.    |
  | Profile Image | text | Field 2.    |
  +---------------+------+-------------+
  ```

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
| CGPA                     | text | Field 8.    |
| Class Level              | text | Field 3.    |
| Extracurricular Activity | text | Field 6.    |
| Gender                   | text | Field 2.    |
| Home State               | text | Field 4.    |
| ID                       | text | Field 1.    |
| Major                    | text | Field 5.    |
| Mentor                   | text | Field 7.    |
| Percentage               | text | Field 9.    |
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
