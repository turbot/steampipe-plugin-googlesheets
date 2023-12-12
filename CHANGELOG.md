## v0.7.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#65](https://github.com/turbot/steampipe-plugin-googlesheets/pull/65))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#65](https://github.com/turbot/steampipe-plugin-googlesheets/pull/65))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-googlesheets/blob/main/docs/LICENSE). ([#65](https://github.com/turbot/steampipe-plugin-googlesheets/pull/65))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#64](https://github.com/turbot/steampipe-plugin-googlesheets/pull/64))

## v0.6.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#45](https://github.com/turbot/steampipe-plugin-googlesheets/pull/45))

## v0.6.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#41](https://github.com/turbot/steampipe-plugin-googlesheets/pull/41))
- Recompiled plugin with Go version `1.21`. ([#41](https://github.com/turbot/steampipe-plugin-googlesheets/pull/41))

## v0.5.0 [2023-08-02]

_Enhancements_

- The `sheets` config argument now allows for wildcard-based searches of sheets, enabling the creation of dynamic tables as needed. By default, no dynamic tables will be created if `sheets` argument is empty or not set.  ([#32](https://github.com/turbot/steampipe-plugin-googlesheets/pull/32))

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
