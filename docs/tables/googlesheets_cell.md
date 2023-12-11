---
title: "Steampipe Table: googlesheets_cell - Query Google Sheets Cells using SQL"
description: "Allows users to query Google Sheets Cells, providing detailed information about each cell in a Google Sheet."
---

# Table: googlesheets_cell - Query Google Sheets Cells using SQL

Google Sheets is a web-based spreadsheet program that is part of Google's office suite of web applications. It allows users to create, update, and modify spreadsheets and share the data live online. The cells in Google Sheets are the smallest and most basic unit of a spreadsheet, storing individual data points that can be manipulated and analyzed.

## Table Usage Guide

The `googlesheets_cell` table offers insights into the data points stored in the cells of a Google Sheet. As a data analyst or data scientist, you can dig into cell-specific details using this table, including the cell's value, format, and associated metadata. Use it to extract and analyze data from Google Sheets, such as cell values, formulas, and formatting details, to facilitate data analysis and reporting.

All examples below can be used with the [Google Sheets Plugin - Sample School
Data](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4)
spreadsheet, which is a public spreadsheet maintained by the Steampipe team.

## Examples

### Query cells in all sheets
Explore the content of all your Google Sheets by identifying the specific cell and its corresponding value. This can be useful for auditing data or tracking changes across multiple spreadsheets.

```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell;
```

```sql+sqlite
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell;
```

### Query cells in a specific sheet
Explore which cells in a specific Google Sheets document contain certain values. This is useful for quickly finding and analyzing data in large spreadsheets, such as a list of students.
You can query all cells from a specific sheet using the `sheet_name` column:


```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  sheet_name = 'Students';
```

```sql+sqlite
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
Explore which values are stored in a specific range of cells in a Google Sheets document. This can be useful for understanding the content and organization of your data without having to manually search through the sheet.
Using the [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1), you can query specific cells using the `range` column:


```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!B1:C2';
```

```sql+sqlite
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
This example illustrates how you can pinpoint specific information within a Google Sheets document. It's particularly useful when you need to quickly access a particular cell's value within a larger spreadsheet, such as a student's information in an educational setting.
You can query a specific cell with the `sheet_name`, `row`, and `col` columns:


```sql+postgres
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

```sql+sqlite
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

```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!A2';
```

```sql+sqlite
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
Determine the content of specific cells within a Google Sheets document. This is useful for quickly accessing and understanding key data without having to manually search through the entire document.
Similar to the examples above, you can also query a specific row using the `sheet_name` and `row` columns:


```sql+postgres
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

```sql+sqlite
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

```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!1:1';
```

```sql+sqlite
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
Determine the areas in which specific data is stored in a Google Sheets document. This is useful in scenarios where you need to identify and analyze the information contained in a particular column of a specific sheet, such as a list of student names.
Specific columns can also be queried using the `sheet_name` and `col` columns:


```sql+postgres
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

```sql+sqlite
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

```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!A:A';
```

```sql+sqlite
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
Identify instances where cells contain hyperlinks in a Google Sheets document, specifically within the 'Students' sheet. This could be useful for quickly locating linked resources or references within a large dataset.

```sql+postgres
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

```sql+sqlite
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
Explore which cells within the 'Employees' sheet contain a formula. This can be useful for identifying calculations or automated data within your spreadsheet.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that contain formula parsing errors within your Google Sheets. This is particularly useful for identifying and rectifying any errors in your data calculations, ensuring the accuracy and integrity of your data.

```sql+postgres
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

```sql+sqlite
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
Explore the specific content within a Google Sheets document by targeting a particular sheet, in this case 'Students'. This is useful for extracting and analyzing data from a specific part of your Google Sheets document without having to sift through irrelevant information.
In [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1), just the sheet name can be passed in as the `range` to return all cells from that sheet:


```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students';
```

```sql+sqlite
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
Explore the values of specific cells in a Google Sheets document, specifically targeting the 'Students' sheet. This is useful for analyzing data within a specified range, such as identifying student information in the first column.
In addition to [A1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-1), the `range` column also supports [R1C1 notation](https://developers.google.com/sheets/api/guides/concepts#expandable-2).

For instance, to get the first five cells in the first column:


```sql+postgres
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!R1C1:R5C1';
```

```sql+sqlite
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!R1C1:R5C1';
```