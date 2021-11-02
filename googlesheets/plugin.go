package googlesheets

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// Plugin creates this (googlesheets) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-googlesheets",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMapFunc:     PluginTables,
	}
	return p
}

func PluginTables(ctx context.Context, p *plugin.Plugin) (map[string]*plugin.Table, error) {
	// Initialize tables
	tables := map[string]*plugin.Table{}

	// Get the list of sheets to be retrieved from the spreadsheet
	googlesheetConfig := GetConfig(p.Connection)

	// Create tablemap for all the available sheets
	for _, sheetName := range googlesheetConfig.Sheets {
		tableCtx := context.WithValue(ctx, "sheet", sheetName)
		tableSchema := tableSpreadsheets(tableCtx, p)
		if tableSchema != nil {
			tables[sheetName] = tableSpreadsheets(tableCtx, p)
		}
	}

	return tables, nil
}
