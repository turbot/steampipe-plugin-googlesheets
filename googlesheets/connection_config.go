package googlesheets

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type googleSheetsConfig struct {
	Credentials           *string  `hcl:"credentials"`
	ImpersonatedUserEmail *string  `hcl:"impersonated_user_email"`
	TokenPath             *string  `hcl:"token_path"`
	SpreadsheetId         *string  `hcl:"spreadsheet_id"`
	Sheets                []string `hcl:"sheets,optional"`
}

func ConfigInstance() interface{} {
	return &googleSheetsConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) googleSheetsConfig {
	if connection == nil || connection.Config == nil {
		return googleSheetsConfig{}
	}
	config, _ := connection.Config.(googleSheetsConfig)
	return config
}
