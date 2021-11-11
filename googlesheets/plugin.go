package googlesheets

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
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

	// Add static tables
	tables["googlesheets_cell"] = tableGooglesheetsCell(ctx)
	tables["googlesheets_sheet"] = tableGooglesheetsSheet(ctx)

	// Add dynamic tables

	// Get the list of sheets to be retrieved from the spreadsheet
	googlesheetConfig := GetConfig(p.Connection)

	if len(googlesheetConfig.Sheets) == 0 {
		return tables, nil
	}

	// Create a map of headers along with correspomding sheet name
	spreadsheetHeadersMap, err := getSpreadsheetHeadersMap(ctx, p, googlesheetConfig.Sheets)
	if err != nil {
		return tables, nil
	}

	for k, v := range spreadsheetHeadersMap {
		spreadsheetHeaders := v

		// Create columns
		cols := []*plugin.Column{}
		for col_index, col := range spreadsheetHeaders {
			cols = append(cols, &plugin.Column{Name: col, Type: proto.ColumnType_STRING, Transform: transform.FromField(col), Description: fmt.Sprintf("Field %d.", col_index)})
		}

		// Create table definition
		tables[k] = &plugin.Table{
			Name:        k,
			Description: fmt.Sprintf("Retrieves data from %s.", k),
			List: &plugin.ListConfig{
				Hydrate: listSpreadsheetWithPath(ctx, p, k),
			},
			Columns: cols,
		}
	}

	return tables, nil
}
