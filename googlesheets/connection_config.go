package googlesheets

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type googleSheetsConfig struct {
	Credentials           *string  `cty:"credentials"`
	ImpersonatedUserEmail *string  `cty:"impersonated_user_email"`
	TokenPath             *string  `cty:"token_path"`
	SpreadsheetId         *string  `cty:"spreadsheet_id"`
	Sheets                []string `cty:"sheets"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"credentials": {
		Type: schema.TypeString,
	},
	"impersonated_user_email": {
		Type: schema.TypeString,
	},
	"token_path": {
		Type: schema.TypeString,
	},
	"spreadsheet_id": {
		Type: schema.TypeString,
	},
	"sheets": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
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
