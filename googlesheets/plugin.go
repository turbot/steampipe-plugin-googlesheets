package googlesheets

import (
	"context"
	"fmt"

	"github.com/turbot/go-kit/helpers"
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

	// Get the list of sheets to be retrieved from the spreadsheet
	googlesheetConfig := GetConfig(p.Connection)

	// Get the headers along with sheet name
	availableSheets, err := getSpreadsheets(ctx, p)
	if err != nil {
		return tables, nil
	}

	// Return all valid sheets
	var validSheets []string
	for _, i := range googlesheetConfig.Sheets {
		if helpers.StringSliceContains(availableSheets, i) {
			validSheets = append(validSheets, i)
		}
	}

	// Get the headers along with sheet name
	headersMap, err := getSpreadsheetHeaders(ctx, p, validSheets)
	if err != nil {
		return tables, nil
	}

	// Create tablemap for all the available sheets
	for _, sheetName := range googlesheetConfig.Sheets {
		// Extract headers from the map
		spreadsheetHeaders := headersMap[sheetName]

		// Create columns
		cols := []*plugin.Column{}
		for idx, j := range spreadsheetHeaders {
			cols = append(cols, &plugin.Column{Name: j, Type: proto.ColumnType_STRING, Transform: transform.FromField(j), Description: fmt.Sprintf("Field %d.", idx)})
		}

		// Create table definition
		tables[sheetName] = &plugin.Table{
			Name:        sheetName,
			Description: fmt.Sprintf("Retrieves data from %s.", sheetName),
			List: &plugin.ListConfig{
				Hydrate: listSpreadsheetWithPath(ctx, p, sheetName),
			},
			Columns: cols,
		}
	}

	return tables, nil
}
