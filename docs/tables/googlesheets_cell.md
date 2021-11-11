# Table: googlesheets_cell

Retrieves information of a non-empty cell in a given sheet.

## Examples

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
