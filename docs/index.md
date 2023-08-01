---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/googlesheets.svg"
brand_color: "#1EA362"
display_name: "Google Sheets"
short_name: "googlesheets"
description: "Steampipe plugin for query data from Google Sheets."
og_description: "Query Google Sheets with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/googlesheets-social-graphic.png"
---

# Google Sheets + Steampipe

[Google Sheets](https://www.google.com/sheets/about) is an online spreadsheet app that lets you create and format spreadsheets and work with other people.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

Query all the sheets in your spreadsheet:

```sql
select
  title
from
  googlesheets_sheet;
```

```sh
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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/googlesheets/tables)**

## Get started

### Install

Download and install the latest Google Sheets plugin:

```shell
steampipe plugin install googlesheets
```

### Credentials

| Item        | Description |
| :---------- | :---------- |
| APIs | 1. Go to the [Google API Console](https://console.cloud.google.com/apis/dashboard). <br/> 2. Select the project that contains your credentials. <br/> 3. Click `Enable APIs and Services`. <br/> 4. Enable both the `Google Drive API` and `Google Sheets API`.
| Credentials | 1. To use **OAuth client**, configure your [credentials](#authenticate-using-oauth-client).<br />2. To use **domain-wide delegation**, generate your [service account and credentials](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#create_the_service_account_and_credentials) and [delegate domain-wide authority to your service account](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#delegate_domain-wide_authority_to_your_service_account). Use `https://www.googleapis.com/auth/drive.readonly` and `https://www.googleapis.com/auth/spreadsheets.readonly` OAuth 2.0 scopes. |
| Radius      | Each connection represents a single Google spreadsheet. |
| Resolution  | 1. Credentials from the JSON file specified by the `credentials` parameter in your Steampipe config.<br />2. Credentials from the JSON file specified by the `token_path` parameter in your Steampipe config.<br />3. Credentials from the default json file location (`~/.config/gcloud/application_default_credentials.json`). |

### Configuration

Installing the latest googlesheets plugin will create a config file (`~/.steampipe/config/googlesheets.spc`) with a single connection named `googlesheets`:

```hcl
connection "googlesheets" {
  plugin = "googlesheets"

  # The spreadsheet ID can be found in the spreadsheet's URL, e.g., https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4
  # spreadsheet_id = "11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4"

  # List of sheets that will be created as dynamic tables.
  # No dynamic tables will be created if this arg is empty or not set.
  # Wildcard based searches are supported.

  # For example:
  #  - "*" matches all sheets
  #  - "Student*" matches all sheets starting with "Student"
  #  - "Books" matches a sheet named "Books"

  # Defaults to all sheets
  # sheets = ["*"]

  # You may connect to Google Sheet using more than one option:

  # 1. To authenticate using domain-wide delegation, specify a service account credential file and the user email for impersonation
  # `credentials` - Either the path to a JSON credential file that contains Google application credentials,
  # or the contents of a service account key file in JSON format. If `credentials` is not specified in a connection,
  # credentials will be loaded from:
  #   - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
  #   - The standard location (`~/.config/gcloud/application_default_credentials.json`)
  #   - The path specified for the credentials.json file ("/path/to/my/creds.json")
  # credentials = "~/.config/gcloud/application_default_credentials.json"

  # `impersonated_user_email` - The email (string) of the user which should be impersonated. Needs permissions to access the Admin APIs.
  # impersonated_user_email = "username@domain.com"

  # 2. To authenticate using OAuth 2.0, specify a client secret file
  # `token_path` - The path to a JSON credential file that contains Google application credentials.
  # If `token_path` is not specified in a connection, credentials will be loaded from:
  #   - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
  #   - The standard location (`~/.config/gcloud/application_default_credentials.json`)
  # token_path = "/Users/myuser/.config/gcloud/application_default_credentials.json"
}
```

The default values for the `spreadsheet_id` and `sheets` arguments are
configured to be used with our [Google Sheets Plugin - Sample School
Data](https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4)
spreadsheet, which is a public spreadsheet maintained by the Steampipe team.
All of our examples in this document, and in our table documents, are based on
this sample spreadsheet.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-googlesheets
- Community: [Slack Channel](https://steampipe.io/community/join)

## Advanced configuration options

### Authenticate using OAuth client

You can use client secret credentials to protect the user's data by only granting tokens to authorized requestors. Use following steps to configure credentials:

- [Configure the OAuth consent screen](https://developers.google.com/workspace/guides/configure-oauth-consent).
- [Create an OAuth client ID credential](https://developers.google.com/workspace/guides/create-credentials#desktop-app) with the application type `Desktop app`, and download the client secret JSON file.
- Wherever you have the [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) installed, run the following command with the correct client secret JSON file parameters:

  ```sh
  gcloud auth application-default login \
    --client-id-file=client_secret.json \
    --scopes="https://www.googleapis.com/auth/drive.readonly,https://www.googleapis.com/auth/spreadsheets.readonly"
  ```

- In the browser window that just opened, authenticate as the user you would like to make the API calls through.
- Review the output for the location of the **Application Default Credentials** file, which usually appears following the text `Credentials saved to file:`.
- Set the **Application Default Credentials** filepath in the Steampipe config `token_path` or in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.
