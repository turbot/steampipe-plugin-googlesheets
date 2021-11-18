---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/googlesheets.svg"
brand_color: "#0F9D58"
display_name: "Google Sheets"
short_name: "googlesheets"
description: "Steampipe plugin for query data from Google Sheets."
og_description: "Query Google Sheets with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/googlesheets-social-graphic.png"
---

# Google Sheets + Steampipe

[Google Sheets](https://www.google.com/sheets/about) is an online spreadsheet app that lets you create and format spreadsheets and work with other people.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  sheet_name,
  cell,
  value
from
  googlesheets_cell
where
  range = 'Students!1:1';
```

```sh
+------------+------+--------------------------+
| sheet_name | cell | value                    |
+------------+------+--------------------------+
| Students   | E1   | Home State               |
| Students   | F1   | Major                    |
| Students   | B1   | ID                       |
| Students   | H1   | Mentor                   |
| Students   | D1   | Class Level              |
| Students   | C1   | Gender                   |
| Students   | I1   | CGPA                     |
| Students   | G1   | Extracurricular Activity |
| Students   | J1   | Percentage               |
| Students   | A1   | Student Name             |
+------------+------+--------------------------+
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
| Credentials | 1. To use **domain-wide delegation**, generate your [service account and credentials](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#create_the_service_account_and_credentials) and [delegate domain-wide authority to your service account](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#delegate_domain-wide_authority_to_your_service_account). Use `https://www.googleapis.com/auth/drive.readonly` OAuth 2.0 scope, so that the service account can access Google Sheet service.<br />2. To use **OAuth client**, configure your [credentials](#authenticate-using-oauth-client). |
| Radius      | Each connection represents a single Google spreadsheet. |
| Resolution  | 1. Credentials from the JSON file specified by the `credentials` parameter in your Steampipe config.<br />2. Credentials from the JSON file specified by the `token_path` parameter in your Steampipe config.<br />3. Credentials from the default json file location (`~/.config/gcloud/application_default_credentials.json`). |

### Configuration

Installing the latest googlesheets plugin will create a config file (`~/.steampipe/config/googlesheets.spc`) with a single connection named `googlesheets`:

```hcl
connection "googlesheets" {
  plugin = "googlesheets"

  # You may connect to Google Sheet using more than one option:
  # 1. To authenticate using domain-wide delegation, specify a service account credential file and the user email for impersonation
  # `credentials` - Either the path to a JSON credential file that contains Google application credentials, or the contents of a service account key file in JSON format.
  # credentials = "/path/to/my/creds.json"
  
  # `impersonated_user_email` - The email (string) of the user which should be impersonated. Needs permissions to access the Admin APIs. 
  # impersonated_user_email = "username@domain.com"

  # 2. To authenticate using OAuth 2.0, specify a client secret file
  # `token_path` - The path to a JSON credential file that contains Google application credentials. 
  # If `token_path` is not specified in a connection, credentials will be loaded from:
      - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
      - The standard location (`~/.config/gcloud/application_default_credentials.json`)
  # token_path = "~/.config/gcloud/application_default_credentials.json"

  # spreadsheet_id = "11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4"
  # sheets         = ["Dashboard", "Students", "Books", "Marks", "Employee"]
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-googlesheets
- Community: [Slack Channel](https://steampipe.io/community/join)

## Advanced configuration options

### Authenticate using OAuth client

You can use client secret credentials to protect the user's data by only granting tokens to authorized requestors. Use following steps to configure credentials:

- [Configure the OAuth consent screen](https://developers.google.com/workspace/guides/create-credentials#configure_the_oauth_consent_screen).
- [Create an OAuth client ID credential](https://developers.google.com/workspace/guides/create-credentials#create_a_oauth_client_id_credential) with the application type `Desktop app`, and download the client secret JSON file.
- Wherever you have the [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) installed, run the following command with the correct client secret JSON file parameters:

  ```sh
  gcloud auth application-default login \
    --client-id-file=client_secret.json \
    --scopes="https://www.googleapis.com/auth/drive.readonly"
  ```

- In the browser window that just opened, authenticate as the user you would like to make the API calls through.
- Review the output for the location of the **Application Default Credentials** file, which usually appears following the text `Credentials saved to file:`.
- Set the **Application Default Credentials** filepath in the Steampipe config `token_path` or in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable.
