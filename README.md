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

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-googlesheets/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Google Sheets Plugin](https://github.com/turbot/steampipe-plugin-googlesheets/labels/help%20wanted)
