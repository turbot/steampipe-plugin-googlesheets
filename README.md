![image](https://hub.steampipe.io/images/plugins/turbot/googlesheets-social-graphic.png)

# Google Sheets Plugin for Steampipe

Use SQL to query data from Google Sheets.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/googlesheets)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/googlesheets/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-googlesheets/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install googlesheets
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/googlesheets#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/googlesheets#configuration).

Run a query:

```sql
select
  first_name,
  last_name
from
  "My Users";
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
