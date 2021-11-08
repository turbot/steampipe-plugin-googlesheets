package googlesheets

import (
	"context"
	"fmt"
	"strings"

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

	// Get spreadsheet details
	spreadsheetData, err := getSpreadsheetData(ctx, p, validSheets)
	if err != nil {
		return tables, nil
	}

	// Create tablemap for all the available sheets
	for _, sheetName := range googlesheetConfig.Sheets {
		for _, data := range spreadsheetData {
			// Return if empty sheet
			if len(data.Values) == 0 {
				continue
			}

			// Return if first row is empty
			if len(data.Values[0]) == 0 {
				continue
			}

			// Return if A1 cell is empty
			if len(data.Values[0][0].(string)) == 0 {
				continue
			}

			str := strings.Split(data.Range, "!")[0]
			var spreadsheetHeaders []string

			// API wraps the range inside quotes if sheet name contains whitespaces
			// For example, if the sheet name is 'Sheet 1', then range comes as "'Sheet 1'!A1:Z1"
			if len(str) > 0 && str[0] == '\'' {
				str = str[1 : len(str)-1]
			}

			if sheetName == str {
				for idx, i := range data.Values[0] {
					if len(i.(string)) == 0 {
						columnName := intToLetters(idx + 1) // since index in for is zero-based
						spreadsheetHeaders = append(spreadsheetHeaders, columnName)
					} else {
						spreadsheetHeaders = append(spreadsheetHeaders, i.(string))
					}
				}

				// Create columns
				cols := []*plugin.Column{}
				for idx, j := range spreadsheetHeaders {
					// If no value passed as header use `?column?` as column name
					if len(j) == 0 {
						j = "?column?"
					}
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
		}
	}

	// Add other static tables
	tables["googlesheets_sheet"] = tableGooglesheetsSheet(ctx)

	return tables, nil
}
