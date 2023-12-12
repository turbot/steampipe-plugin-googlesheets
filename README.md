![image](https://hub.steampipe.io/images/plugins/turbot/googlesheets-social-graphic.png)

# Google Sheets Plugin for Steampipe

Use SQL to query data from Google Sheets.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/googlesheets)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/googlesheets/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-googlesheets/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install googlesheets
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/googlesheets#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/googlesheets#configuration).

Run steampipe:

```shell
steampipe query
```

Query all the sheets in your spreadsheet:

```sql
select
  title
from
  googlesheets_sheet;
```

```
+-----------+
| title     |
+-----------+
| Marks     |
| Students  |
| Dashboard |
| Books     |
| Employees |
+-----------+
```

A table is automatically created to represent each sheet mentioned in the
configured `sheets` config argument.

For instance, if using the following configuration:

```hcl
connection "googlesheets" {
  plugin = "googlesheets"

  # Google Sheets Plugin - Sample School Data
  spreadsheet_id = "11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4"
  sheets         = ["Dashboard", "Students", "Books", "Marks", "Employees"]
}
```

To get details for a specific table, inspect it by name:

```shell
.inspect "Students"
+--------------------------+------+-------------+
| column                   | type | description |
+--------------------------+------+-------------+
| Class Level              | text | Field 2.    |
| Extracurricular Activity | text | Field 5.    |
| GPA                      | text | Field 7.    |
| Home State               | text | Field 3.    |
| ID                       | text | Field 1.    |
| Major                    | text | Field 4.    |
| Mentor                   | text | Field 6.    |
| Student Name             | text | Field 0.    |
+--------------------------+------+-------------+
```

You can query data from the `Students` sheet using the sheet's column names:

```sql
select
  "Student Name",
  "Class Level",
  "GPA"
from
  "Students";
```

```sh
+--------------+-------------+------+
| Student Name | Class Level | GPA  |
+--------------+-------------+------+
| Andrew       | Freshman    | 3.50 |
| Alexandra    | Senior      | 3.20 |
| Karen        | Sophomore   | 3.90 |
| Edward       | Junior      | 3.40 |
| Ellen        | Freshman    | 2.90 |
| Becky        | Sophomore   | 3.10 |
+--------------+-------------+------+
```

You can also query data from a sheet by cell address:

```sql
select
  cell,
  value
from
  googlesheets_cell
where
  sheet_name = 'Students';
```

```sh
+------+--------------+
| cell | value        |
+------+--------------+
| A1   | Student Name |
| A2   | Alexandra    |
| A3   | Andrew       |
| A4   | Anna         |
| A5   | Becky        |
| B1   | ID           |
| B2   | 1            |
| B3   | 2            |
| B4   | 3            |
| B5   | 4            |
+------|--------------+
```

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/index) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs//steampipe_sqlite/index) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/index) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-googlesheets.git
cd steampipe-plugin-googlesheets
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/googlesheets.spc
```

Try it!

```shell
steampipe query
> .inspect googlesheets
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Google Sheets Plugin](https://github.com/turbot/steampipe-plugin-googlesheets/labels/help%20wanted)
