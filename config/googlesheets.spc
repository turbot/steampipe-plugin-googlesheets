connection "googlesheets" {
  plugin = "googlesheets"

  # The spreadsheet ID can be found in the spreadsheet's URL, e.g., https://docs.google.com/spreadsheets/d/11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4
  # spreadsheet_id = "11iXfj-RHpFsil7_hNK-oQjCqmBLlDfCvju2AOF-ieb4"

  # If no sheets are specified, then all sheets will be retrieved
  # sheets = ["Dashboard", "Students", "Books", "Marks", "Employees"]

  # You may connect to Google Sheet using more than one option:

  # 1. To authenticate using OAuth 2.0, specify a client secret file
  # `token_path` - The path to a JSON credential file that contains Google application credentials.
  # If `token_path` is not specified in a connection, credentials will be loaded from:
  #   - The path specified in the `GOOGLE_APPLICATION_CREDENTIALS` environment variable, if set; otherwise
  #   - The standard location (`~/.config/gcloud/application_default_credentials.json`)
  # token_path = "/Users/myuser/.config/gcloud/application_default_credentials.json"

  # 2. To authenticate using domain-wide delegation, specify a service account credential file and the user email for impersonation
  # `credentials` - Either the path to a JSON credential file that contains Google application credentials, or the contents of a service account key file in JSON format.
  # credentials = "/path/to/my/creds.json"

  # `impersonated_user_email` - The email (string) of the user which should be impersonated. Needs permissions to access the Admin APIs.
  # impersonated_user_email = "username@domain.com"
}
