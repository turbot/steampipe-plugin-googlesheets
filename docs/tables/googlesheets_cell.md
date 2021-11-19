# Table: googlesheets_cell

Query cell data from sheets in a Google Sheets spreadsheet. Cells that have no
data, i.e., no value or formula, will not be returned.

When specifying the `range` key qual, you can use use [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1) or [R1C1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-2).

## Examples

### Query cells in all sheets

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
```

### Query cells in a specific sheet

You can query all cells from a specific sheet using the `range` or `sheet_name` column:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students';
```

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  sheet_name = 'Students';
```

### Query a range of cells

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!B1:C2';
```

### Query a specific cell

You can query a specific cell's information using the `range` column:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!A2';
```

Or with the `row` and `col` columns:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  sheet_name = 'Students'
  and row = 2
  and col = 'A';
```

### Query cells in a row or column

Similar to the examples above, you can also query a specific row or column using the `range` column:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!1:1';
```

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!A:A';
```

Or using the `row` and `col` columns:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  sheet_name = 'Students'
  and row = 1;
```

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  sheet_name = 'Students'
  and col = 'A';
```

### List cells with hyperlink information

```sql
select
  sheet_name,
  cell,
  value,
  hyperlink
from
  googlesheets_cell
where
  sheet_name = 'Students'
  and hyperlink is not null;
```

### List cells with a formula

```sql
select
  sheet_name,
  cell,
  value,
  formula
from
  googlesheets_cell
where
  sheet_name = 'Employees'
  and formula is not null;
```

### List cells with formula parse errors

```sql
select
  sheet_name,
  cell,
  value,
  formula
from
  googlesheets_cell
where
  sheet_name = 'Employees'
  and formula is not null
  and value in ('#N/A', '#DIV/0!', '#VALUE!', '#REF!', '#NAME?', '#NUM!', '#ERROR!', '#NULL!');
```

### Create a pivot table using cells within a specific range

You can use the cell data to build a more readable table using pivot tables:

```sql
with cells as (
  select
    *
  from
    googlesheets_cell
  where
    range = 'Marks!A2:C9'
  order by row, col
),
pivot_cells as (
  select
    row,
    max(case when col = 'A' then value else null end) as name,
    max(case when col = 'B' then value else null end) as exam,
    max(case when col = 'C' then value else null end) as score
  from
    cells
  group by row
),
pivot_marks as (
  select
    name,
    max(case when exam = '1' then score else null end) as exam_1,
    max(case when exam = '2' then score else null end) as exam_2,
    max(case when exam = '3' then score else null end) as exam_3,
    max(case when exam = '4' then score else null end) as exam_4
  from
    pivot_cells
  group by
    name
)
select
  *
from
  pivot_marks;
```
