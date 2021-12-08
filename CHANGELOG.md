## v0.1.0 [2021-12-08]

_Enhancements_

- Recompile plugn with Go version 1.17 ([#11](https://github.com/turbot/steampipe-plugin-googlesheets/pull/11))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#10](https://github.com/turbot/steampipe-plugin-googlesheets/pull/10))
- `docs/index.md` now includes Google Drive and Sheets API scope information ([#9](https://github.com/turbot/steampipe-plugin-googlesheets/pull/9))

## v0.0.2 [2021-11-22]

_Bug fixes_

- Fixed: Improve error message if `spreadsheet_id` config arg is `nil` when running a query ([#7](https://github.com/turbot/steampipe-plugin-googlesheets/pull/7))

## v0.0.1 [2021-11-21]

_What's new?_

- New tables added
  - [googlesheets_cell](https://hub.steampipe.io/plugins/turbot/googlesheets/tables/googlesheets_cell)
  - [googlesheets_sheet](https://hub.steampipe.io/plugins/turbot/googlesheets/tables/googlesheets_sheet)
  - [googlesheets_spreadsheet](https://hub.steampipe.io/plugins/turbot/googlesheets/tables/googlesheets_spreadsheet)
  - [{sheet_name}](https://hub.steampipe.io/plugins/turbot/googlesheets/tables/{sheet_name})
