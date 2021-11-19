# Table: googlesheets_cell

Retrieves information of a non-empty cell in a given sheet.

## Examples

### List all cells within a specific range

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

### List all cells in a column

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

### List all cells in a row

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

### Get a specific cell

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

### Get a specific cell using `row` and `col`

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

### List all cells with hyperlink information

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

### List all cells with formula

```sql
select
  sheet_name,
  cell,
  value,
  formula
from
  googlesheets_cell
where
  sheet_name = 'Students'
  and formula is not null;
```

### List cells with formula parse error

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
