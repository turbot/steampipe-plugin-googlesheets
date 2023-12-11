---
title: "Steampipe Table: googlesheets_spreadsheet - Query Google Sheets Spreadsheets using SQL"
description: "Allows users to query Google Sheets Spreadsheets, specifically extracting data from spreadsheets, including metadata, content, and sharing permissions."
---

# Table: googlesheets_spreadsheet - Query Google Sheets Spreadsheets using SQL

Google Sheets is a web-based application that allows users to create, update, and modify spreadsheets and share the data live online. The service is part of the Google Suite of applications and is designed to be a free, web-based alternative to traditional spreadsheet software such as Microsoft Excel. Google Sheets allows multiple users to collaborate in real-time, providing a platform for data manipulation and analysis.

## Table Usage Guide

The `googlesheets_spreadsheet` table provides insights into the spreadsheets within Google Sheets. As a data analyst, explore spreadsheet-specific details through this table, including metadata, content, and sharing permissions. Utilize it to uncover information about spreadsheets, such as those shared with specific users, the content of the spreadsheets, and the metadata associated with each spreadsheet.

## Examples

### Basic info
Explore which Google Sheets spreadsheets have been recently modified, owned, shared, starred, or linked on the web. This can help in managing access and tracking changes to important documents.

```sql+postgres
select
  name,
  modified_time,
  owned_by_me,
  shared,
  starred,
  web_view_link
from
  googlesheets_spreadsheet;
```

```sql+sqlite
select
  name,
  modified_time,
  owned_by_me,
  shared,
  starred,
  web_view_link
from
  googlesheets_spreadsheet;
```

### Get information about last modifying user
Explore which user last modified a Google Sheets spreadsheet and when the modification occurred. This can be useful for tracking changes and maintaining accountability in collaborative environments.

```sql+postgres
select
  name as spreadsheet_name,
  last_modifying_user ->> 'displayName' as user_display_name,
  last_modifying_user ->> 'emailAddress' as user_email_address,
  modified_time
from
  googlesheets_spreadsheet;
```

```sql+sqlite
select
  name as spreadsheet_name,
  json_extract(last_modifying_user, '$.displayName') as user_display_name,
  json_extract(last_modifying_user, '$.emailAddress') as user_email_address,
  modified_time
from
  googlesheets_spreadsheet;
```

### Check if current user has capability to edit the spreadsheet
Explore the access rights of the current user to understand if they have the necessary permissions to edit a Google Sheets spreadsheet. This can be useful in managing user permissions and ensuring data integrity.

```sql+postgres
select
  name as spreadsheet_name,
  web_view_link,
  case
    when capabilities -> 'canEdit' is null then false
    else (capabilities ->> 'canEdit')::boolean
  end as can_edit
from
  googlesheets_spreadsheet;
```

```sql+sqlite
select
  name as spreadsheet_name,
  web_view_link,
  case
    when json_extract(capabilities, '$.canEdit') is null then 0
    else json_extract(capabilities, '$.canEdit')
  end as can_edit
from
  googlesheets_spreadsheet;
```

### Get sharing info with permissions
Explore the sharing details and permissions of your Google Sheets spreadsheets. This query can help you understand who has access to your documents, their roles, and whether they can discover files, providing a comprehensive view of your document's security.

```sql+postgres
select
  name,
  permission ->> 'type' as grantee_type,
  permission ->> 'displayName' as grantee_display_name,
  permission ->> 'emailAddress' as email_address,
  permission ->> 'domain' as domain_address,
  permission ->> 'role' as role,
  permission ->> 'allowFileDiscovery' as allow_file_discovery
from
  googlesheets_spreadsheet,
  jsonb_array_elements(permissions) as permission;
```

```sql+sqlite
select
  name,
  json_extract(permission.value, '$.type') as grantee_type,
  json_extract(permission.value, '$.displayName') as grantee_display_name,
  json_extract(permission.value, '$.emailAddress') as email_address,
  json_extract(permission.value, '$.domain') as domain_address,
  json_extract(permission.value, '$.role') as role,
  json_extract(permission.value, '$.allowFileDiscovery') as allow_file_discovery
from
  googlesheets_spreadsheet,
  json_each(permissions) as permission;
```