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
