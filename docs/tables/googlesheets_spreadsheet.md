# Table: googlesheets_spreadsheet

Retrieves information of the given spreadsheet.

## Examples

### Basic info

```sql
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

```sql
select
  name as spreadsheet_name,
  last_modifying_user ->> 'displayName' as user_display_name,
  last_modifying_user ->> 'emailAddress' as user_email_address,
  modified_time
from
  googlesheets_spreadsheet;
```

### Get sharing info with permissions

```sql
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
