connection "googlesheets" {
  plugin = "googlesheets"

  # You may connect to Google Sheet using more than one option:
  # 1. To authenticate using domain-wide delegation, specify a service account credential file and the user email for impersonation
  # credential_file         = "/path/to/my/creds.json"
  # impersonated_user_email = "username@domain.com"

  # 2. To authenticate using OAuth 2.0, specify a client secret file
  # token_path = "~/.config/gcloud/application_default_credentials.json"

  spreadsheet_id = "1kSJ8giCfSndHQCoVaOCVk1CBAXHelloxmEtda5gvdX0"
  sheets         = ["My Users", "My Products"]
}