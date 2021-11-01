package googlesheets

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type csvConfig struct {
	Credentials           *string  `cty:"credentials"`
	ImpersonatedUserEmail *string  `cty:"impersonated_user_email"`
	TokenPath             *string  `cty:"token_path"`
	SheetId               *string  `cty:"sheet_id"`
	Ranges                []string `cty:"ranges"`
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
	"sheet_id": {
		Type: schema.TypeString,
	},
	"ranges": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
}

func ConfigInstance() interface{} {
	return &csvConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) csvConfig {
	if connection == nil || connection.Config == nil {
		return csvConfig{}
	}
	config, _ := connection.Config.(csvConfig)
	return config
}
