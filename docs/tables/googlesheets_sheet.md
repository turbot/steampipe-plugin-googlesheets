---
title: "Steampipe Table: googlesheets_sheet - Query Google Sheets Sheets using SQL"
description: "Allows users to query Sheets in Google Sheets, specifically the metadata and content of each sheet, providing insights into the structure of Google Sheets documents."
---

# Table: googlesheets_sheet - Query Google Sheets Sheets using SQL

Google Sheets is a web-based application that allows users to create, update and modify spreadsheets and share the data live online. The Sheets within Google Sheets are individual tabs of a document, where data is input and organized. Each sheet contains cells, and the data within these cells can be manipulated using formulas, functions, and data validation rules.

## Table Usage Guide

The `googlesheets_sheet` table provides insights into Sheets within Google Sheets. As a data analyst, explore sheet-specific details through this table, including metadata, content, and associated cell data. Utilize it to uncover information about sheets, such as their structure, the data they contain, and the relationships between different data points.

## Examples

### Basic info
Explore which Google Sheets are available in your account and their respective identifiers. This can be beneficial for managing and tracking your documents systematically.

```sql+postgres
select
  title,
  sheet_id,
  spreadsheet_id
from
  googlesheets_sheet;
```

```sql+sqlite
select
  title,
  sheet_id,
  spreadsheet_id
from
  googlesheets_sheet;
```

### Get information about a specific sheet
Explore which Google Sheets contain specific information by identifying instances where the title matches a certain term. This allows you to quickly locate and analyze data within a vast collection of spreadsheets.

```sql+postgres
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet
where
  title = 'Students';
```

```sql+sqlite
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet
where
  title = 'Students';
```

### List hidden sheets
Uncover the details of hidden sheets within your Google Sheets. This query is useful for identifying which sheets are hidden, allowing you to better manage and organize your data.

```sql+postgres
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet
where
  hidden;
```

```sql+sqlite
select
  title,
  sheet_id,
  spreadsheet_id,
  hidden,
  sheet_type
from
  googlesheets_sheet
where
  hidden = 1;
```

### List sheets with protected ranges
Explore which Google Sheets contain protected ranges. This can be useful to identify instances where data is safeguarded, helping to maintain data integrity and control access.

```sql+postgres
select
  title,
  sheet_id,
  spreadsheet_id,
  protected_ranges
from
  googlesheets_sheet
where
  protected_ranges is not null;
```

```sql+sqlite
select
  title,
  sheet_id,
  spreadsheet_id,
  protected_ranges
from
  googlesheets_sheet
where
  protected_ranges is not null;
```