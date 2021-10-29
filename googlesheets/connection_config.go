package googlesheets

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type csvConfig struct {
	SheetId   *string  `cty:"sheet_id"`
	Ranges     []string  `cty:"ranges"`
}

var ConfigSchema = map[string]*schema.Attribute{
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
