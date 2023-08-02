## v0.5.0 [2023-08-02]

_Enhancements_

- The `sheets` config argument now allows for wildcard-based searches of sheets, enabling the creation of dynamic tables as needed. By default no dynamic tables will be created if `sheets` argument is empty or not set.  ([#32](https://github.com/turbot/steampipe-plugin-googlesheets/pull/32))

## v0.4.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. This update significantly lowers the plugin initialization time of dynamic plugins by avoiding recursing into child folders when not necessary. ([#29](https://github.com/turbot/steampipe-plugin-googlesheets/pull/29))

## v0.3.0 [2023-03-22]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#26](https://github.com/turbot/steampipe-plugin-googlesheets/pull/26))
- Recompiled plugin with Go version `1.19`. ([#26](https://github.com/turbot/steampipe-plugin-googlesheets/pull/26))

## v0.2.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#22](https://github.com/turbot/steampipe-plugin-googlesheets/pull/22))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#20](https://github.com/turbot/steampipe-plugin-googlesheets/pull/20))

## v0.1.2 [2022-04-14]

_Bug fixes_

- Fixed links in documentation for configuring OAuth client authentication.

## v0.1.1 [2022-01-31]

_Bug fixes_

- Fixed: Credentials in the `credentials` argument now take precedence over those in the `token_path` argument ([#15](https://github.com/turbot/steampipe-plugin-googlesheets/pull/15))

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
