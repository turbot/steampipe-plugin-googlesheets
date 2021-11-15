# Table: googlesheets_cell

Retrieves information of a non-empty cell in a given sheet.

## Examples

### List all cells within a specific range

```sql
select
  sheet_name,
  cell_address,
  value
from
  googlesheets_cell
where
  sheet_name = 'My Users'
  and ranges = 'B1:C2';
```

### List all cells with hyperlink information

```sql
select
  sheet_name,
  cell_address,
  value,
  hyperlink
from
  googlesheets_cell
where
  sheet_name = 'My Users'
  and hyperlink is not null;
```

### List all cells with formula

```sql
select
  sheet_name,
  cell_address,
  value,
  formula
from
  googlesheets_cell
where
  sheet_name = 'My Users'
  and formula is not null;
```

### List cells with formula parse error

```sql
select
  sheet_name,
  cell_address,
  value,
  formula
from
  googlesheets_cell
where
  sheet_name = 'My Users'
  and formula is not null
  and value in ('#N/A', '#DIV/0!', '#VALUE!', '#REF!', '#NAME?', '#NUM!', '#ERROR!', '#NULL!');
```
