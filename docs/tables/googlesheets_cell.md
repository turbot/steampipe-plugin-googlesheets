# Table: googlesheets_cell

Query cell data from sheets in a Google Sheets spreadsheet. Cells that have no
data, i.e., no value or formula, will not be returned.

All examples below can be used with the [Google Sheets Plugin - Sample School
Data](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4)
spreadsheet, which is a public spreadsheet maintained by the Steampipe team.

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

You can query all cells from a specific sheet using the `sheet_name` column:

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

Using the [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1), you can query specific cells using the `range` column:

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

You can query a specific cell with the `sheet_name`, `row`, and `col` columns:

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

Or you can use [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1) with the `range` column:

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

### Query cells in a row

Similar to the examples above, you can also query a specific row using the `sheet_name` and `row` columns:

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

Or by using [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1) with the `range` column:

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

### Query cells in a column

Specific columns can also be queried using the `sheet_name` and `col` columns:

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

Or with [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1) and the `range` column:

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

## Advanced examples

### Query cells in a specific sheet using `range`

In [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1), just the sheet name can be passed in as the `range` to return all cells from that sheet:

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

### Query cells using [R1C1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-2)

In addition to [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1), the `range` column also supports [R1C1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-2).

For instance, to get the first five cells in the first column:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!R1C1:R5C1';
```
